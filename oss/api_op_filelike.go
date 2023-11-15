package oss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type OpenOptions struct {
	Context   context.Context
	Offset    int64
	VersionId *string
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
}

// Open opens the named file for reading.
// If successful, methods on the returned file can be used for reading.
func (c *Client) OpenFile(bucket, key string, optFns ...func(*OpenOptions)) (*ReadOnlyFile, error) {
	options := OpenOptions{
		Context: context.Background(),
		Offset:  0,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	f := &ReadOnlyFile{
		client:  c,
		context: options.Context,

		bucket:    bucket,
		key:       key,
		versionId: options.VersionId,

		offset: options.Offset,
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
		}
	}()

	if offset >= f.sizeInBytes {
		err = io.EOF
		return
	}

	if f.readBufOffset != offset {
		f.readBufOffset = offset
		if f.reader != nil {
			f.reader.Close()
			f.reader = nil
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
