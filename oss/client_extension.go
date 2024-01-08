package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	MaxUploadParts int32 = 10000

	MinUploadPartSize int64 = 5 * 1024 * 1024

	DefaultUploadPartSize = MinUploadPartSize

	DefaultUploadParallel = 3

	DefaultDownloadPartSize = MinUploadPartSize

	DefaultDownloadParallel = 3

	FilePermMode = os.FileMode(0664) // Default file permission

	TempFileSuffix = ".temp" // Temp file suffix

	CheckpointMagic = "92611BED-89E2-46B6-89E5-72F273D4B0A3"
)

type UploaderOptions struct {
	PartSize int64

	ParallelNum int

	LeavePartsOnError bool

	ClientOptions []func(*Options)
}

type Uploader struct {
	options UploaderOptions
	client  *Client
}

func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader {
	options := UploaderOptions{
		PartSize:          DefaultUploadPartSize,
		ParallelNum:       DefaultUploadParallel,
		LeavePartsOnError: false,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	u := &Uploader{
		client:  c,
		options: options,
	}

	return u
}

type UploadRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// The caching behavior of the web page when the object is downloaded.
	CacheControl *string `input:"header,Cache-Control"`

	// The method that is used to access the object.
	ContentDisposition *string `input:"header,Content-Disposition"`

	// The method that is used to encode the object.
	ContentEncoding *string `input:"header,Content-Encoding"`

	// The size of the data in the HTTP message body. Unit: bytes.
	ContentLength *int64 `input:"header,Content-Length"`

	// The MD5 hash of the object that you want to upload.
	ContentMD5 *string `input:"header,Content-MD5"`

	// A standard MIME type describing the format of the contents.
	ContentType *string `input:"header,Content-Type"`

	// The expiration time of the cache in UTC.
	Expires *string `input:"header,Expires"`

	// Specifies whether the object that is uploaded by calling the PutObject operation
	// overwrites an existing object that has the same name. Valid values: true and false
	ForbidOverwrite *string `input:"header,x-oss-forbid-overwrite"`

	// The encryption method on the server side when an object is created.
	// Valid values: AES256 and KMS
	ServerSideEncryption *string `input:"header,x-oss-server-side-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	// This header is valid only when the x-oss-server-side-encryption header is set to KMS.
	ServerSideDataEncryption *string `input:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `input:"header,x-oss-server-side-encryption-key-id"`

	// The access control list (ACL) of the object.
	Acl ObjectACLType `input:"header,x-oss-object-acl"`

	// The storage class of the object.
	StorageClass StorageClassType `input:"header,x-oss-storage-class"`

	// The metadata of the object that you want to upload.
	Metadata map[string]string `input:"header,x-oss-meta-,usermeta"`

	// The tags that are specified for the object by using a key-value pair.
	// You can specify multiple tags for an object. Example: TagA=A&TagB=B.
	Tagging *string `input:"header,x-oss-tagging"`

	// A callback parameter is a Base64-encoded string that contains multiple fields in the JSON format.
	Callback *string `input:"header,x-oss-callback"`

	// Configure custom parameters by using the callback-var parameter.
	CallbackVar *string `input:"header,x-oss-callback-var"`

	// Specify the speed limit value. The speed limit value ranges from 245760 to 838860800, with a unit of bit/s.
	TrafficLimit int64 `input:"header,x-oss-traffic-limit"`

	RequestCommon
}

type UploadResult struct {
	UploadId *string

	ETag *string

	VersionId *string

	HashCRC64 *string

	ResultCommon
}

type UploadError struct {
	Err      error
	UploadId string
	Path     string
}

func (m *UploadError) Error() string {
	var extra string
	if m.Err != nil {
		extra = fmt.Sprintf(", cause: %s", m.Err.Error())
	}
	return fmt.Sprintf("upload failed, upload id: %s%s", m.UploadId, extra)
}

func (m *UploadError) Unwrap() error {
	return m.Err
}

func (u *Uploader) UploadFrom(ctx context.Context, request *UploadRequest, body io.Reader, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
	// Uploader wrapper
	delegate, err := u.newDelegate(ctx, request, optFns...)
	if err != nil {
		return nil, err
	}

	delegate.body = body
	if err = delegate.applySource(); err != nil {
		return nil, err
	}

	return delegate.upload()
}

func (u *Uploader) UploadFile(ctx context.Context, request *UploadRequest, filePath string, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
	// Uploader wrapper
	delegate, err := u.newDelegate(ctx, request, optFns...)
	if err != nil {
		return nil, err
	}

	// Source
	if err = delegate.checkSource(filePath); err != nil {
		return nil, err
	}

	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return nil, err
	}
	defer func() {
		defer file.Close()
	}()
	delegate.body = file

	if err = delegate.applySource(); err != nil {
		return nil, err
	}

	result, err := delegate.upload()

	//TODO check CRC

	return result, err
}

type uploaderDelegate struct {
	options UploaderOptions
	client  *Client
	context context.Context
	request *UploadRequest

	readerPos int64
	totalSize int64

	filePath string

	body io.Reader

	partPool byteSlicePool
}

func (u *Uploader) newDelegate(ctx context.Context, request *UploadRequest, optFns ...func(*UploaderOptions)) (*uploaderDelegate, error) {
	if request == nil {
		return nil, NewErrParamNull("request")
	}

	if request.Bucket == nil {
		return nil, NewErrParamNull("request.Bucket")
	}

	if request.Key == nil {
		return nil, NewErrParamNull("request.Key")
	}

	d := uploaderDelegate{
		options: u.options,
		client:  u.client,
		context: ctx,
		request: request,
	}

	for _, opt := range optFns {
		opt(&d.options)
	}

	if d.options.ParallelNum <= 0 {
		d.options.ParallelNum = DefaultUploadParallel
	}
	if d.options.PartSize <= 0 {
		d.options.PartSize = DefaultUploadPartSize
	}

	return &d, nil
}

func (u *uploaderDelegate) checkSource(filePath string) error {
	if filePath == "" {
		return NewErrParamRequired("filePath")
	}

	if !FileExists(filePath) {
		return fmt.Errorf("File not exists, %v", filePath)
	}

	u.filePath = filePath

	return nil
}

func (u *uploaderDelegate) applySource() error {
	if u.body == nil {
		return NewErrParamNull("uploader's body is null")
	}

	totalSize := getReaderLen(u.body)

	//Part Size
	partSize := u.options.PartSize
	if totalSize > 0 {
		for totalSize/partSize >= int64(MaxUploadParts) {
			partSize += u.options.PartSize
		}
	}

	u.totalSize = totalSize
	u.options.PartSize = partSize

	return nil
}

func (u *uploaderDelegate) upload() (*UploadResult, error) {
	if u.totalSize >= 0 && u.totalSize < u.options.PartSize {
		return u.singlePart()
	}
	return u.multiPart()
}

func (u *uploaderDelegate) singlePart() (*UploadResult, error) {
	request := &PutObjectRequest{}
	copyRequest(request, u.request)
	request.Body = u.body
	if request.ContentType == nil {
		request.ContentType = u.getContentType()
	}

	result, err := u.client.PutObject(u.context, request, u.options.ClientOptions...)

	if err != nil {
		return nil, u.wrapErr("", err)
	}

	return &UploadResult{
		ETag:         result.ETag,
		VersionId:    result.VersionId,
		HashCRC64:    result.HashCRC64,
		ResultCommon: result.ResultCommon,
	}, nil
}

func (u *uploaderDelegate) nextReader() (io.ReadSeeker, int, func(), error) {
	type readerAtSeeker interface {
		io.ReaderAt
		io.ReadSeeker
	}
	switch r := u.body.(type) {
	case readerAtSeeker:
		var err error

		n := u.options.PartSize
		if u.totalSize >= 0 {
			bytesLeft := u.totalSize - u.readerPos
			if bytesLeft <= u.options.PartSize {
				err = io.EOF
				n = bytesLeft
			}
		}

		reader := io.NewSectionReader(r, u.readerPos, n)
		cleanup := func() {}

		u.readerPos += n

		return reader, int(n), cleanup, err

	default:
		if u.partPool == nil {
			u.partPool = newByteSlicePool(u.options.PartSize)
			u.partPool.ModifyCapacity(u.options.ParallelNum + 1)
		}

		part, err := u.partPool.Get(u.context)
		if err != nil {
			return nil, 0, func() {}, err
		}

		n, err := readFill(r, *part)
		u.readerPos += int64(n)

		cleanup := func() {
			u.partPool.Put(part)
		}

		return bytes.NewReader((*part)[0:n]), n, cleanup, err
	}
}

type uploaderChunk struct {
	partNum int32
	body    io.ReadSeeker
	cleanup func()
}

func (u *uploaderDelegate) multiPart() (*UploadResult, error) {
	release := func() {
		if u.partPool != nil {
			u.partPool.Close()
		}
	}
	defer release()

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		parts    UploadParts
		errValue atomic.Value
	)

	// Init the multipart
	initRequest := &InitiateMultipartUploadRequest{}
	copyRequest(initRequest, u.request)
	if initRequest.ContentType == nil {
		initRequest.ContentType = u.getContentType()
	}

	initResult, err := u.client.InitiateMultipartUpload(u.context, initRequest, u.options.ClientOptions...)
	if err != nil {
		return nil, u.wrapErr("", err)
	}

	//fmt.Printf("InitiateMultipartUpload result: %#v, %#v\n", initResult, err)

	saveErrFn := func(e error) {
		errValue.Store(e)
	}

	getErrFn := func() error {
		v := errValue.Load()
		if v == nil {
			return nil
		}
		e, _ := v.(error)
		return e
	}

	// readChunk runs in worker goroutines to pull chunks off of the ch channel
	readChunkFn := func(ch chan uploaderChunk) {
		defer wg.Done()
		for {
			data, ok := <-ch
			if !ok {
				break
			}

			if getErrFn() == nil {
				upResult, err := u.client.UploadPart(u.context, &UploadPartRequest{
					Bucket:     u.request.Bucket,
					Key:        u.request.Key,
					UploadId:   initResult.UploadId,
					PartNumber: data.partNum,
					RequestCommon: RequestCommon{
						Body: data.body,
					}}, u.options.ClientOptions...)

				//fmt.Printf("UploadPart result: %#v, %#v\n", upResult, err)

				if err == nil {
					mu.Lock()
					parts = append(parts, UploadPart{ETag: upResult.ETag, PartNumber: data.partNum})
					mu.Unlock()
				} else {
					saveErrFn(err)
				}
			}
			data.cleanup()
		}
	}

	ch := make(chan uploaderChunk, u.options.ParallelNum)
	for i := 0; i < u.options.ParallelNum; i++ {
		wg.Add(1)
		go readChunkFn(ch)
	}

	// Read and queue the parts
	var (
		qnum int32 = 0
		qerr error = nil
	)
	for getErrFn() == nil && qerr == nil {
		var (
			reader       io.ReadSeeker
			nextChunkLen int
			cleanup      func()
		)

		reader, nextChunkLen, cleanup, qerr = u.nextReader()
		// check err
		if (qerr != nil && qerr != io.EOF) ||
			nextChunkLen == 0 {
			cleanup()
			saveErrFn(qerr)
			break
		}
		qnum++
		//fmt.Printf("send chunk: %d\n", qnum)
		ch <- uploaderChunk{body: reader, partNum: qnum, cleanup: cleanup}
	}

	// Close the channel, wait for workers
	close(ch)
	wg.Wait()

	// Complete upload
	var cmResult *CompleteMultipartUploadResult
	if err = getErrFn(); err == nil {
		sort.Sort(parts)
		cmRequest := &CompleteMultipartUploadRequest{}
		copyRequest(cmRequest, u.request)
		cmRequest.UploadId = initResult.UploadId
		cmRequest.CompleteMultipartUpload = &CompleteMultipartUpload{Parts: parts}
		cmResult, err = u.client.CompleteMultipartUpload(u.context, cmRequest, u.options.ClientOptions...)
	}
	//fmt.Printf("CompleteMultipartUpload cmResult: %#v, %#v\n", cmResult, err)

	if err != nil {
		//TODO Abort
		return nil, u.wrapErr(*initResult.UploadId, err)
	}

	return &UploadResult{
		UploadId:     initResult.UploadId,
		ETag:         cmResult.ETag,
		VersionId:    cmResult.VersionId,
		HashCRC64:    cmResult.HashCRC64,
		ResultCommon: cmResult.ResultCommon,
	}, nil
}

func (u *uploaderDelegate) getContentType() *string {
	if u.filePath != "" {
		if contentType := TypeByExtension(u.filePath); contentType != "" {
			return Ptr(contentType)
		}
	}
	return nil
}

func (u *uploaderDelegate) wrapErr(uploadId string, err error) error {
	return &UploadError{
		UploadId: uploadId,
		Path:     fmt.Sprintf("oss://%s/%s", *u.request.Bucket, *u.request.Key),
		Err:      err}
}

type DownloaderOptions struct {
	PartSize int64

	ParallelNum int

	EnableCheckpoint bool

	CheckpointDir string

	UseTempFile bool

	ClientOptions []func(*Options)
}

type Downloader struct {
	options DownloaderOptions
	client  *Client
}

func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader {
	options := DownloaderOptions{
		PartSize:    DefaultUploadPartSize,
		ParallelNum: DefaultUploadParallel,
		UseTempFile: true,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	u := &Downloader{
		client:  c,
		options: options,
	}

	return u
}

type DownloadRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// If the ETag specified in the request matches the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	IfMatch *string `input:"header,If-Match"`

	// If the ETag specified in the request does not match the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	IfNoneMatch *string `input:"header,If-None-Match"`

	// If the time specified in this header is earlier than the object modified time or is invalid,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfModifiedSince *string `input:"header,If-Modified-Since"`

	// If the time specified in this header is the same as or later than the object modified time,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfUnmodifiedSince *string `input:"header,If-Unmodified-Since"`

	// The content range of the object to be returned.
	// If the value of Range is valid, the total size of the object and the content range are returned.
	// For example, Content-Range: bytes 0~9/44 indicates that the total size of the object is 44 bytes,
	// and the range of data returned is the first 10 bytes.
	// However, if the value of Range is invalid, the entire object is returned,
	// and the response does not include the Content-Range parameter.
	Range *string `input:"header,Range"`

	// The cache-control header to be returned in the response.
	ResponseCacheControl *string `input:"query,response-cache-control"`

	// The content-disposition header to be returned in the response.
	ResponseContentDisposition *string `input:"query,response-content-disposition"`

	// The content-encoding header to be returned in the response.
	ResponseContentEncoding *string `input:"query,response-content-encoding"`

	// The content-language header to be returned in the response.
	ResponseContentLanguage *string `input:"query,response-content-language"`

	// The content-type header to be returned in the response.
	ResponseContentType *string `input:"query,response-content-type"`

	// The expires header to be returned in the response.
	ResponseExpires *string `input:"query,response-expires"`

	// VersionId used to reference a specific version of the object.
	VersionId *string `input:"query,versionId"`

	// Specify the speed limit value. The speed limit value ranges from 245760 to 838860800, with a unit of bit/s.
	TrafficLimit int64 `input:"header,x-oss-traffic-limit"`
}

type DownloadResult struct {
	Written int64
}

type DownloadError struct {
	Err  error
	Path string
}

func (m *DownloadError) Error() string {
	var extra string
	if m.Err != nil {
		extra = fmt.Sprintf(", cause: %s", m.Err.Error())
	}
	return fmt.Sprintf("download failed, %s", extra)
}

func (m *DownloadError) Unwrap() error {
	return m.Err
}

func (d *Downloader) DownloadFile(ctx context.Context, request *DownloadRequest, filePath string, optFns ...func(*DownloaderOptions)) (result *DownloadResult, err error) {
	// Downloader wrapper
	delegate, err := d.newDelegate(ctx, request, optFns...)
	if err != nil {
		return nil, err
	}

	// Source
	if err = delegate.checkSource(); err != nil {
		return nil, err
	}

	// Destination
	if err = delegate.checkDestination(filePath); err != nil {
		return nil, err
	}

	// Range
	if err = delegate.adjustRange(); err != nil {
		return nil, err
	}

	// Checkpoint
	if err = delegate.checkCheckpoint(); err != nil {
		return nil, err
	}

	// open file to write
	var file *os.File
	file, err = delegate.openWriter()
	if err != nil {
		return nil, err
	}

	// download
	result, err = delegate.download()

	return result, delegate.closeWriter(file, err)
}

type downloaderDelegate struct {
	options DownloaderOptions
	client  *Client
	context context.Context

	m sync.Mutex

	request DownloadRequest
	w       io.WriterAt
	rstart  int64
	pos     int64
	epos    int64
	written int64

	// Source's Info
	sizeInBytes int64
	etag        string
	modTime     string
	headers     http.Header

	//Destination's Info
	filePath     string
	tempFilePath string

	checkpoint *downloadCheckpoint
}

type downloaderChunk struct {
	w      io.WriterAt
	start  int64
	size   int64
	cur    int64
	rstart int64 //range start
}

type downloadedChunk struct {
	start int64
	size  int64
	crc64 uint64
}

type downloadedChunks []downloadedChunk

func (slice downloadedChunks) Len() int {
	return len(slice)
}
func (slice downloadedChunks) Less(i, j int) bool {
	return slice[i].start < slice[j].start
}
func (slice downloadedChunks) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (c *downloaderChunk) Write(p []byte) (n int, err error) {
	if c.cur >= c.size {
		return 0, io.EOF
	}

	n, err = c.w.WriteAt(p, c.start+c.cur-c.rstart)
	c.cur += int64(n)
	return
}

func (d *Downloader) newDelegate(ctx context.Context, request *DownloadRequest, optFns ...func(*DownloaderOptions)) (*downloaderDelegate, error) {
	if request == nil {
		return nil, NewErrParamNull("request")
	}

	if !isValidBucketName(request.Bucket) {
		return nil, NewErrParamInvalid("request.Bucket")
	}

	if !isValidObjectName(request.Key) {
		return nil, NewErrParamInvalid("request.Key")
	}

	if request.Range != nil && !isValidRange(request.Range) {
		return nil, NewErrParamInvalid("request.Range")
	}

	delegate := downloaderDelegate{
		options: d.options,
		client:  d.client,
		context: ctx,
		request: *request,
	}

	for _, opt := range optFns {
		opt(&delegate.options)
	}

	if delegate.options.ParallelNum <= 0 {
		delegate.options.ParallelNum = DefaultDownloadParallel
	}
	if delegate.options.PartSize <= 0 {
		delegate.options.PartSize = DefaultDownloadPartSize
	}

	return &delegate, nil
}

func (d *downloaderDelegate) checkSource() error {
	var request HeadObjectRequest
	copyRequest(&request, &d.request)
	result, err := d.client.HeadObject(d.context, &request, d.options.ClientOptions...)
	if err != nil {
		return err
	}

	d.sizeInBytes = result.ContentLength
	d.modTime = result.Headers.Get(HTTPHeaderLastModified)
	d.etag = result.Headers.Get(HTTPHeaderETag)
	d.headers = result.Headers

	return nil
}

func (d *downloaderDelegate) checkDestination(filePath string) error {
	if filePath == "" {
		return NewErrParamInvalid("filePath")
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// use temporary file
	tempFilePath := absFilePath
	if d.options.UseTempFile {
		tempFilePath += TempFileSuffix
	}
	d.filePath = absFilePath
	d.tempFilePath = tempFilePath

	// use openfile to check the filepath is valid
	var file *os.File
	if file, err = os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE, FilePermMode); err != nil {
		return err
	}
	file.Close()

	return nil
}

func (d *downloaderDelegate) openWriter() (*os.File, error) {
	file, err := os.OpenFile(d.tempFilePath, os.O_WRONLY|os.O_CREATE, FilePermMode)
	if err != nil {
		return nil, err
	}

	if err = file.Truncate(d.pos - d.rstart); err != nil {
		file.Close()
		return nil, err
	}

	d.w = file
	return file, nil
}

func (d *downloaderDelegate) closeWriter(file *os.File, err error) error {
	if file != nil {
		file.Close()
	}

	if err != nil {
		if d.checkpoint == nil {
			os.Remove(d.tempFilePath)
		}
	} else {
		if d.tempFilePath != d.filePath {
			err = os.Rename(d.tempFilePath, d.filePath)
		}
		if err == nil && d.checkpoint != nil {
			d.checkpoint.remove()
		}
	}

	d.w = nil
	d.checkpoint = nil

	return err
}

func (d *downloaderDelegate) adjustRange() error {
	d.pos = 0
	d.rstart = 0
	d.epos = d.sizeInBytes
	if d.request.Range != nil {
		httpRange, _ := ParseRange(*d.request.Range)
		if httpRange.Offset >= d.sizeInBytes {
			return fmt.Errorf("invalid range, object size :%v, range: %v", d.sizeInBytes, ToString(d.request.Range))
		}
		d.pos = httpRange.Offset
		d.rstart = d.pos
		if httpRange.Count > 0 {
			d.epos = minInt64(httpRange.Offset+httpRange.Count, d.sizeInBytes)
		}
	}

	return nil
}

func (d *downloaderDelegate) checkCheckpoint() error {
	if d.options.EnableCheckpoint {
		d.checkpoint = newDownloadCheckpoint(&d.request, d.tempFilePath, d.options.CheckpointDir, d.headers, d.options.PartSize)
		if err := d.checkpoint.load(); err != nil {
			return err
		}

		if d.checkpoint.Loaded {
			d.pos = d.checkpoint.Info.Data.DownloadInfo.Offset
			d.written = d.pos - d.rstart
		} else {
			d.checkpoint.Info.Data.DownloadInfo.Offset = d.pos
		}
	}
	return nil
}

func (d *downloaderDelegate) download() (*DownloadResult, error) {
	var (
		wg       sync.WaitGroup
		errValue atomic.Value
		cpCh     chan downloadedChunk
		cpWg     sync.WaitGroup
		cpChunks downloadedChunks
	)

	saveErrFn := func(e error) {
		errValue.Store(e)
	}

	getErrFn := func() error {
		v := errValue.Load()
		if v == nil {
			return nil
		}
		e, _ := v.(error)
		return e
	}

	// writeChunkFn runs in worker goroutines to pull chunks off of the ch channel
	writeChunkFn := func(ch chan downloaderChunk) {
		defer wg.Done()
		for {
			chunk, ok := <-ch
			if !ok {
				break
			}

			if getErrFn() != nil {
				continue
			}

			dchunk, derr := d.downloadChunk(chunk)

			if derr != nil && derr != io.EOF {
				saveErrFn(derr)
			} else {
				// update checkpoint
				if d.checkpoint != nil {
					cpCh <- dchunk
				}
			}
		}
	}

	// checkpointFn runs in worker goroutines to update checkpoint info
	checkpointFn := func(ch chan downloadedChunk) {
		defer cpWg.Done()
		for {
			chunk, ok := <-ch
			if !ok {
				break
			}
			cpChunks = append(cpChunks, chunk)
			sort.Sort(cpChunks)
			newOffset := d.checkpoint.Info.Data.DownloadInfo.Offset
			i := 0
			for ii := range cpChunks {
				if cpChunks[ii].start == newOffset {
					newOffset += cpChunks[ii].size
					i++
				} else {
					break
				}
			}
			if newOffset != d.checkpoint.Info.Data.DownloadInfo.Offset {
				//remove updated chunk in cpChunks
				cpChunks = cpChunks[i:]
				d.checkpoint.Info.Data.DownloadInfo.Offset = newOffset
				d.checkpoint.dump()
			}
		}
	}

	// Start the download workers
	ch := make(chan downloaderChunk, d.options.ParallelNum)
	for i := 0; i < d.options.ParallelNum; i++ {
		wg.Add(1)
		go writeChunkFn(ch)
	}

	// Start checkpoint worker if enable checkpoint
	if d.checkpoint != nil {
		cpCh = make(chan downloadedChunk, 1)
		cpWg.Add(1)
		go checkpointFn(cpCh)
	}

	// Queue the next range of bytes to read.
	for getErrFn() == nil {
		if d.pos >= d.epos {
			break
		}
		size := minInt64(d.epos-d.pos, d.options.PartSize)
		ch <- downloaderChunk{w: d.w, start: d.pos, size: size, rstart: d.rstart}
		d.pos += size
	}

	// Waiting for parts download finished
	close(ch)
	wg.Wait()

	if d.checkpoint != nil {
		close(cpCh)
		cpWg.Wait()
	}

	if err := getErrFn(); err != nil {
		return nil, &DownloadError{
			Err:  err,
			Path: fmt.Sprintf("oss://%s/%s", ToString(d.request.Bucket), ToString(d.request.Key)),
		}
	}

	return &DownloadResult{
		Written: d.written,
	}, nil
}

func (d *downloaderDelegate) incrWritten(n int64) {
	d.m.Lock()
	defer d.m.Unlock()

	d.written += n
}

func (d *downloaderDelegate) downloadChunk(chunk downloaderChunk) (downloadedChunk, error) {
	// Get the next byte range of data
	var request GetObjectRequest
	copyRequest(&request, &d.request)

	getFn := func(ctx context.Context, httpRange HTTPRange) (output *ReaderRangeGetOutput, err error) {
		// update range
		request.Range = nil
		rangeStr := httpRange.FormatHTTPRange()
		request.RangeBehavior = nil
		if rangeStr != nil {
			request.Range = rangeStr
			request.RangeBehavior = Ptr("standard")
		}

		result, err := d.client.GetObject(ctx, &request, d.options.ClientOptions...)
		if err != nil {
			return nil, err
		}

		return &ReaderRangeGetOutput{
			Body:          result.Body,
			ETag:          result.ETag,
			ContentLength: result.ContentLength,
			ContentRange:  result.ContentRange,
		}, nil
	}

	reader, _ := NewRangeReader(d.context, getFn, &HTTPRange{chunk.start, chunk.size}, d.etag)
	n, err := io.Copy(&chunk, reader)
	d.incrWritten(n)

	return downloadedChunk{
		start: chunk.start,
		size:  n,
		crc64: uint64(0),
	}, err
}

// ----- Concurrent download with chcekpoint  -----

type downloadCheckpoint struct {
	CpDirPath  string // checkpoint dir full path
	CpFilePath string // checkpoint file full path
	VerifyData bool   // verify donwloaded data in FilePath
	Loaded     bool   // If Info.Data.DownloadInfo is loaded from checkpoint

	Info struct { //checkpoint data
		Magic string // Magic
		MD5   string // The Data's MD5
		Data  struct {
			// source
			ObjectInfo struct {
				Name      string // oss://bucket/key
				VersionId string
				Range     string
			}
			ObjectMeta struct {
				Size         int64
				LastModified string
				ETag         string
			}

			// destination
			FilePath string // Local file

			// download info
			PartSize int64

			DownloadInfo struct {
				Offset int64
				CRC64  uint64
			}
		}
	}
}

func newDownloadCheckpoint(request *DownloadRequest, filePath string, baseDir string, header http.Header, partSize int64) *downloadCheckpoint {
	var buf strings.Builder
	name := fmt.Sprintf("%v/%v", ToString(request.Bucket), ToString(request.Key))
	buf.WriteString("oss://" + escapePath(name, false))
	buf.WriteString("\n")
	buf.WriteString(ToString(request.VersionId))
	buf.WriteString("\n")
	buf.WriteString(ToString(request.Range))

	hashmd5 := md5.New()
	hashmd5.Write([]byte(buf.String()))
	srcHash := hex.EncodeToString(hashmd5.Sum(nil))

	absPath, _ := filepath.Abs(filePath)
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))

	var dir string
	if baseDir == "" {
		dir = os.TempDir()
	} else {
		dir = filepath.Dir(baseDir)
	}

	cpFilePath := filepath.Join(dir, fmt.Sprintf("%v-%v.cp", srcHash, destHash))

	cp := &downloadCheckpoint{
		CpFilePath: cpFilePath,
		CpDirPath:  dir,
	}

	objectSize, _ := strconv.ParseInt(header.Get("Content-Length"), 10, 64)

	cp.Info.Magic = CheckpointMagic
	cp.Info.Data.ObjectInfo.Name = "oss://" + name
	cp.Info.Data.ObjectInfo.VersionId = ToString(request.VersionId)
	cp.Info.Data.ObjectInfo.Range = ToString(request.Range)
	cp.Info.Data.ObjectMeta.Size = objectSize
	cp.Info.Data.ObjectMeta.LastModified = header.Get("Last-Modified")
	cp.Info.Data.ObjectMeta.ETag = header.Get("ETag")
	cp.Info.Data.FilePath = filePath
	cp.Info.Data.PartSize = partSize

	return cp
}

// load checkpoint from local file
func (cp *downloadCheckpoint) load() error {
	if !DirExists(cp.CpDirPath) {
		return fmt.Errorf("Invaid checkpoint dir, %v", cp.CpDirPath)
	}

	if !FileExists(cp.CpFilePath) {
		return nil
	}

	if !cp.valid() {
		cp.remove()
		return nil
	}

	cp.Loaded = true

	return nil
}

func (cp *downloadCheckpoint) valid() bool {
	// Compare the CP's Magic and the MD5
	contents, err := os.ReadFile(cp.CpFilePath)
	if err != nil {
		return false
	}

	dcp := downloadCheckpoint{}

	if err = json.Unmarshal(contents, &dcp.Info); err != nil {
		return false
	}

	js, _ := json.Marshal(dcp.Info.Data)
	sum := md5.Sum(js)
	md5sum := hex.EncodeToString(sum[:])

	if CheckpointMagic != dcp.Info.Magic ||
		md5sum != dcp.Info.MD5 {
		return false
	}

	// compare
	if !reflect.DeepEqual(cp.Info.Data.ObjectInfo, dcp.Info.Data.ObjectInfo) ||
		!reflect.DeepEqual(cp.Info.Data.ObjectMeta, dcp.Info.Data.ObjectMeta) ||
		cp.Info.Data.FilePath != dcp.Info.Data.FilePath ||
		cp.Info.Data.PartSize != dcp.Info.Data.PartSize {
		return false
	}

	// download info
	if dcp.Info.Data.DownloadInfo.Offset < 0 {
		return false
	}

	if dcp.Info.Data.DownloadInfo.Offset == 0 &&
		dcp.Info.Data.DownloadInfo.CRC64 != 0 {
		return false
	}

	rOffset := int64(0)
	if len(cp.Info.Data.ObjectInfo.Range) > 0 {
		if r, err := ParseRange(cp.Info.Data.ObjectInfo.Range); err != nil {
			return false
		} else {
			rOffset = r.Offset
		}
	}

	if dcp.Info.Data.DownloadInfo.Offset < rOffset {
		return false
	}

	remains := (dcp.Info.Data.DownloadInfo.Offset - rOffset) % dcp.Info.Data.PartSize
	if remains != 0 {
		return false
	}

	//valid data
	if cp.VerifyData && dcp.Info.Data.DownloadInfo.CRC64 != 0 {
		if file, err := os.Open(cp.Info.Data.FilePath); err == nil {
			hash := NewCRC64(0)
			limitN := dcp.Info.Data.DownloadInfo.Offset - rOffset
			io.Copy(hash, io.LimitReader(file, limitN))
			file.Close()
			if hash.Sum64() != dcp.Info.Data.DownloadInfo.CRC64 {
				return false
			}
		}
	}

	// update
	cp.Info.Data.DownloadInfo = dcp.Info.Data.DownloadInfo

	return true
}

// dump dumps to file
func (cp *downloadCheckpoint) dump() error {
	// Calculate MD5
	js, _ := json.Marshal(cp.Info.Data)
	sum := md5.Sum(js)
	md5sum := hex.EncodeToString(sum[:])
	cp.Info.MD5 = md5sum

	// Serialize
	js, err := json.Marshal(cp.Info)
	if err != nil {
		return err
	}

	// Dump
	return os.WriteFile(cp.CpFilePath, js, FilePermMode)
}

func (cp *downloadCheckpoint) remove() error {
	return os.Remove(cp.CpFilePath)
}
