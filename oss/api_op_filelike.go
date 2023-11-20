package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultReadAheadThreshold = int64(20 * 1024 * 1024)
	defaultChunkSize          = int64(8 * 1024 * 1024)
	defaultParallelNum        = 3
)

type OpenOptions struct {
	Context   context.Context
	Offset    int64
	VersionId *string

	EnableParallel     bool
	ParallelNum        int
	ChunkSize          int64
	ReadAheadThreshold int64
}

type ReadOnlyFile struct {
	client  *Client
	context context.Context

	// object info
	bucket    string
	key       string
	versionId *string

	// file info
	sizeInBytes int64
	modTime     string
	etag        string
	headers     http.Header

	// current read position
	offset int64

	// read
	reader        io.ReadCloser
	readBufOffset int64

	// parallel read
	enableParallel     bool
	chunkSize          int64
	parallelNum        int
	readAheadThreshold int64

	asyncReaders  []*AsyncReader
	seqReadAmount int64
	numOOORead    int64 // number of out of order read
}

// Open opens the named file for reading.
// If successful, methods on the returned file can be used for reading.
func (c *Client) OpenFile(bucket, key string, optFns ...func(*OpenOptions)) (*ReadOnlyFile, error) {
	options := OpenOptions{
		Context:            context.Background(),
		Offset:             0,
		EnableParallel:     false,
		ParallelNum:        defaultParallelNum,
		ChunkSize:          defaultChunkSize,
		ReadAheadThreshold: defaultReadAheadThreshold,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	if options.EnableParallel {
		var chunkSize int64
		if options.ChunkSize > 0 {
			chunkSize = (options.ChunkSize + AsyncReadeBufferSize - 1) / AsyncReadeBufferSize * AsyncReadeBufferSize
		} else {
			chunkSize = defaultChunkSize
		}
		options.ChunkSize = chunkSize

		if options.ParallelNum <= 0 {
			options.ParallelNum = defaultParallelNum
		}
	}

	f := &ReadOnlyFile{
		client:  c,
		context: options.Context,

		bucket:    bucket,
		key:       key,
		versionId: options.VersionId,

		offset: options.Offset,

		enableParallel:     options.EnableParallel,
		parallelNum:        options.ParallelNum,
		chunkSize:          options.ChunkSize,
		readAheadThreshold: options.ReadAheadThreshold,
	}

	var query map[string]string
	if f.versionId != nil {
		query["versionId"] = ToString(f.versionId)
	}
	result, err := f.client.InvokeOperation(f.context, &OperationInput{
		OpName:     "HeadObject",
		Method:     "HEAD",
		Bucket:     Ptr(f.bucket),
		Key:        Ptr(f.key),
		Parameters: query,
	})

	if err != nil {
		return nil, err
	}
	if result.Body != nil {
		result.Body.Close()
	}

	//File info
	f.modTime = result.Headers.Get("Last-Modified")
	f.etag = result.Headers.Get("ETag")
	f.headers = result.Headers
	_, f.sizeInBytes = parseOffsetAndSizeFromHeaders(result.Headers)

	if f.sizeInBytes < 0 {
		return nil, fmt.Errorf("file size is invaid, got %v", f.sizeInBytes)
	}

	if f.offset > f.sizeInBytes {
		return nil, fmt.Errorf("offset is unavailable, offset:%v, file size:%v", f.offset, f.sizeInBytes)
	}

	return f, nil
}

// Close closes the File.
func (f *ReadOnlyFile) Close() error {
	if f.reader != nil {
		f.reader.Close()
	}
	f.reader = nil
	return nil
}

// Read reads up to len(b) bytes from the File and stores them in b.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (f *ReadOnlyFile) Read(p []byte) (bytesRead int, err error) {
	defer func() {
		f.offset += int64(bytesRead)
	}()
	nwant := len(p)
	var nread int
	for bytesRead < nwant && err == nil {
		nread, err = f.readInternal(f.offset+int64(bytesRead), p[bytesRead:])
		if nread > 0 {
			bytesRead += nread
		}
	}
	return
}

// Seek sets the offset for the next Read or Write on file to offset, interpreted
// according to whence: 0 means relative to the origin of the file, 1 means
// relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error.
func (f *ReadOnlyFile) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = f.offset + offset
	case io.SeekEnd:
		abs = f.sizeInBytes + offset
	default:
		return 0, fmt.Errorf("Seek: invalid whence")
	}
	if abs < 0 {
		return 0, fmt.Errorf("Seek: negative position")
	}
	if abs > f.sizeInBytes {
		return offset - (abs - f.sizeInBytes), fmt.Errorf("Seek: offset is unavailable")
	}

	f.offset = abs

	return abs, nil
}

type fileInfo struct {
	name    string
	size    int64
	modTime time.Time
	header  http.Header
}

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() os.FileMode  { return os.FileMode(0644) }
func (fi *fileInfo) ModTime() time.Time { return fi.modTime }
func (fi *fileInfo) IsDir() bool        { return false }
func (fi *fileInfo) Sys() any           { return fi.header }

// Stat returns the FileInfo structure describing file.
// If there is an error, it will be of type *PathError.
func (f *ReadOnlyFile) Stat() (os.FileInfo, error) {
	var name string
	if f.versionId != nil {
		name = fmt.Sprintf("oss://%s/%s?%s", f.bucket, f.key, *f.versionId)
	} else {
		name = fmt.Sprintf("oss://%s/%s", f.bucket, f.key)
	}
	mtime, _ := http.ParseTime(f.modTime)
	return &fileInfo{
		name:    name,
		size:    f.sizeInBytes,
		modTime: mtime,
		header:  f.headers,
	}, nil
}

func (f *ReadOnlyFile) readInternal(offset int64, p []byte) (bytesRead int, err error) {
	defer func() {
		if bytesRead > 0 {
			f.readBufOffset += int64(bytesRead)
			f.seqReadAmount += int64(bytesRead)
		}
	}()

	if offset >= f.sizeInBytes {
		err = io.EOF
		return
	}

	if f.readBufOffset != offset {
		f.readBufOffset = offset
		f.seqReadAmount = 0

		if f.reader != nil {
			f.reader.Close()
			f.reader = nil
		}

		if f.asyncReaders != nil {
			f.numOOORead++
		}

		for _, ar := range f.asyncReaders {
			ar.Close()
		}
		f.asyncReaders = nil
	}

	if f.enableParallel && f.seqReadAmount >= f.readAheadThreshold && f.numOOORead < 3 {
		//swith to async reader
		if f.reader != nil {
			f.reader.Close()
			f.reader = nil
		}

		err = f.readAhead(offset, len(p))
		if err == nil {
			bytesRead, err = f.readFromReadAhead(offset, p)
			return
		} else {
			// fall back to read serially
			f.seqReadAmount = 0
			for _, ar := range f.asyncReaders {
				ar.Close()
			}
			f.asyncReaders = nil
		}
	}

	bytesRead, err = f.readDirect(offset, p)
	return
}

func (f *ReadOnlyFile) readDirect(offset int64, buf []byte) (bytesRead int, err error) {
	if offset >= f.sizeInBytes {
		return
	}

	if f.reader == nil {
		result, err := f.client.GetObject(f.context, &GetObjectRequest{
			Bucket:        Ptr(f.bucket),
			Key:           Ptr(f.key),
			VersionId:     f.versionId,
			Range:         Ptr(fmt.Sprintf("bytes=%d-", offset)),
			RangeBehavior: Ptr("standard"),
		})
		if err != nil {
			return bytesRead, err
		}

		if err = f.checkFileChanged(offset, result.Headers); err != nil {
			return bytesRead, err
		}

		f.reader = result.Body
	}

	bytesRead, err = f.reader.Read(buf)
	if err != nil {
		f.reader.Close()
		f.reader = nil
		err = nil
	}

	return
}

func (f *ReadOnlyFile) checkFileChanged(offset int64, header http.Header) error {
	modTime := header.Get("Last-Modified")
	etag := header.Get("ETag")
	gotOffset, _ := parseOffsetAndSizeFromHeaders(header)
	if gotOffset != offset {
		return fmt.Errorf("Range get fail, expect offset:%v, got offset:%v", offset, gotOffset)
	}

	if modTime != f.modTime || etag != f.etag {
		return fmt.Errorf("Source file is changed, origin info [%v,%v], new info [%v,%v]",
			f.modTime, f.etag, modTime, etag)
	}
	return nil
}

func (f *ReadOnlyFile) readFromReadAhead(offset int64, buf []byte) (bytesRead int, err error) {
	var nread int
	for len(f.asyncReaders) != 0 {
		readAheadBuf := f.asyncReaders[0]
		nread, err = readAheadBuf.Read(buf)
		bytesRead += nread
		if err != nil {
			if err == io.EOF {
				fmt.Printf("readAheadBuf done\n")
				readAheadBuf.Close()
				f.asyncReaders = f.asyncReaders[1:]
				err = nil
			} else {
				return
			}
		}
		buf = buf[nread:]
		if len(buf) == 0 {
			return
		}
	}

	return
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func (f *ReadOnlyFile) readAhead(offset int64, needAtLeast int) (err error) {
	off := offset
	for _, ar := range f.asyncReaders {
		off = ar.oriHttpRange.Offset + ar.oriHttpRange.Count
	}
	//fmt.Printf("readAhead:offset %v, needAtLeast:%v, off:%v\n", offset, needAtLeast, off)
	for len(f.asyncReaders) < f.parallelNum {
		remaining := f.sizeInBytes - off
		size := minInt64(remaining, f.chunkSize)
		cnt := (size + (AsyncReadeBufferSize - 1)) / AsyncReadeBufferSize
		//fmt.Printf("f.sizeInBytes:%v, off:%v, size:%v, cnt:%v\n", f.sizeInBytes, off, size, cnt)
		if size != 0 {
			getFn := func(ctx context.Context, httpRange HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
				request := &GetObjectRequest{
					Bucket:    Ptr(f.bucket),
					Key:       Ptr(f.key),
					VersionId: f.versionId,
				}
				rangeStr := httpRange.FormatHTTPRange()
				if rangeStr != nil {
					request.Range = rangeStr
					request.RangeBehavior = Ptr("standard")
				}
				result, err := f.client.GetObject(f.context, request)
				if err != nil {
					return nil, 0, "", err
				}
				//fmt.Printf("result.Headers:%#v\n", result.Headers)
				offset, _ = parseOffsetAndSizeFromHeaders(result.Headers)
				return result.Body, offset, result.Headers.Get("ETag"), nil
			}
			ar, err := NewAsyncReader(f.context, getFn, &HTTPRange{off, size}, f.etag, int(cnt))
			if err != nil {
				break
			}
			f.asyncReaders = append(f.asyncReaders, ar)
			off += size
		}

		if size != f.chunkSize {
			break
		}
	}
	return nil
}
