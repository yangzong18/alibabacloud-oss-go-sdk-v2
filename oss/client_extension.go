package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"sync/atomic"
)

const (
	MaxUploadParts int32 = 10000

	MinUploadPartSize int64 = 5 * 1024 * 1024

	DefaultUploadPartSize = MinUploadPartSize

	DefaultUploadParallel = 3
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

func (u *Uploader) Upload(ctx context.Context, request *UploadRequest, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
	delegate, err := u.newDelegate(request, "", optFns...)
	if err != nil {
		return nil, err
	}

	delegate.client = u.client
	delegate.context = ctx
	delegate.body = request.Body

	return delegate.upload()
}

func (u *Uploader) UploadFile(ctx context.Context, request *UploadRequest, filePath string, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
	delegate, err := u.newDelegate(request, filePath, optFns...)
	if err != nil {
		return nil, err
	}

	delegate.client = u.client
	delegate.context = ctx
	f, _ := os.Open(filePath)
	delegate.body = f
	if f != nil {
		defer f.Close()
	}

	result, err := delegate.upload()

	//TODO check CRC

	return result, err
}

func (u *Uploader) newDelegate(request *UploadRequest, filePath string, optFns ...func(*UploaderOptions)) (*uploaderDelegate, error) {
	var stat os.FileInfo
	if request == nil {
		return nil, NewErrParamNull("request")
	}

	if request.Bucket == nil {
		return nil, NewErrParamNull("request.Bucket")
	}

	if request.Key == nil {
		return nil, NewErrParamNull("request.Key")
	}

	if filePath != "" {
		fd, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer fd.Close()
		if stat, err = fd.Stat(); err != nil {
			return nil, err
		}
	}

	d := uploaderDelegate{
		options:  u.options,
		request:  request,
		filePath: filePath,
	}

	for _, opt := range optFns {
		opt(&d.options)
	}

	if d.options.ParallelNum == 0 {
		d.options.ParallelNum = DefaultUploadParallel
	}
	if d.options.PartSize == 0 {
		d.options.PartSize = DefaultUploadPartSize
	}

	//Total Size
	totalSize := int64(-1)
	if filePath != "" {
		totalSize = stat.Size()
	} else {
		if request.Body == nil {
			totalSize = 0
		} else {
			totalSize = getReaderLen(request.Body)
		}
	}

	//Part Size
	partSize := d.options.PartSize
	if totalSize > 0 {
		for totalSize/partSize >= int64(MaxUploadParts) {
			partSize += d.options.PartSize
		}
	}

	d.totalSize = totalSize
	d.options.PartSize = partSize

	return &d, nil
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
