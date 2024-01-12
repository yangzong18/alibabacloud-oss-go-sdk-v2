package oss

import (
	"bytes"
	"context"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

	EnableCheckpoint bool

	CheckpointDir string

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

func (u *Uploader) UploadFrom(ctx context.Context, request *PutObjectRequest, body io.Reader, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
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

func (u *Uploader) UploadFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*UploaderOptions)) (*UploadResult, error) {
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
	if file, err = delegate.openReader(); err != nil {
		return nil, err
	}
	delegate.body = file

	if err = delegate.applySource(); err != nil {
		return nil, err
	}

	if err = delegate.checkCheckpoint(); err != nil {
		return nil, err
	}

	if err = delegate.adjustSource(); err != nil {
		return nil, err
	}

	result, err := delegate.upload()

	return result, delegate.closeReader(file, err)
}

type uploaderDelegate struct {
	options UploaderOptions
	client  *Client
	context context.Context
	request *PutObjectRequest

	body      io.Reader
	readerPos int64
	totalSize int64
	hashCRC64 uint64

	// Source's Info, from file or reader
	filePath string
	fileInfo os.FileInfo

	// for resume upload
	uploadId   string
	partNumber int32

	partPool byteSlicePool

	checkpoint *uploadCheckpoint
}

func (u *Uploader) newDelegate(ctx context.Context, request *PutObjectRequest, optFns ...func(*UploaderOptions)) (*uploaderDelegate, error) {
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

	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	u.filePath = filePath
	u.fileInfo = info

	return nil
}

func (u *uploaderDelegate) applySource() error {
	if u.body == nil {
		return NewErrParamNull("the body is null")
	}

	totalSize := GetReaderLen(u.body)

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

func (u *uploaderDelegate) adjustSource() error {
	// resume from upload id
	if u.uploadId != "" {
		// if the body supports seek
		r, ok := u.body.(io.Seeker)
		// not support
		if !ok {
			u.uploadId = ""
			return nil
		}

		// if upload id is valid
		paginator := u.client.NewListPartsPaginator(&ListPartsRequest{
			Bucket:   u.request.Bucket,
			Key:      u.request.Key,
			UploadId: Ptr(u.uploadId),
		})

		// find consecutive sequence from min part number
		var (
			checkPartNumber int32  = 1
			updateCRC64     bool   = u.client.hasFeature(FeatureEnableCRC64CheckUpload)
			hashCRC64       uint64 = 0
		)
	outerLoop:
		for paginator.HasNext() {
			page, err := paginator.NextPage(u.context, u.options.ClientOptions...)
			if err != nil {
				u.uploadId = ""
				return nil
			}
			for _, p := range page.Parts {
				if p.PartNumber != checkPartNumber ||
					p.Size != u.options.PartSize {
					break outerLoop
				}
				checkPartNumber++
				if updateCRC64 && p.HashCRC64 != nil {
					value, _ := strconv.ParseUint(ToString(p.HashCRC64), 10, 64)
					hashCRC64 = CRC64Combine(hashCRC64, value, uint64(p.Size))
				}
			}
		}

		partNumber := checkPartNumber - 1
		newOffset := int64(partNumber) * u.options.PartSize
		if _, err := r.Seek(newOffset, io.SeekStart); err != nil {
			u.uploadId = ""
			return nil
		}
		u.partNumber = partNumber
		u.readerPos = newOffset
		u.hashCRC64 = hashCRC64
	}
	return nil
}

func (d *uploaderDelegate) checkCheckpoint() error {
	if d.options.EnableCheckpoint {
		d.checkpoint = newUploadCheckpoint(d.request, d.filePath, d.options.CheckpointDir, d.fileInfo, d.options.PartSize)
		if err := d.checkpoint.load(); err != nil {
			return err
		}

		if d.checkpoint.Loaded {
			d.uploadId = d.checkpoint.Info.Data.UploadInfo.UploadId
		}
		d.options.LeavePartsOnError = true
	}
	return nil
}

func (d *uploaderDelegate) openReader() (*os.File, error) {
	file, err := os.Open(d.filePath)
	if err != nil {
		return nil, err
	}

	d.body = file
	return file, nil
}

func (d *uploaderDelegate) closeReader(file *os.File, err error) error {
	if file != nil {
		file.Close()
	}

	if d.checkpoint != nil && err == nil {
		d.checkpoint.remove()
	}

	d.body = nil
	d.checkpoint = nil

	return err
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
	size    int
	body    io.ReadSeeker
	cleanup func()
}

type uploadPartCRC struct {
	partNumber int32
	size       int
	hashCRC64  *string
}

type uploadPartCRCs []uploadPartCRC

func (slice uploadPartCRCs) Len() int {
	return len(slice)
}
func (slice uploadPartCRCs) Less(i, j int) bool {
	return slice[i].partNumber < slice[j].partNumber
}
func (slice uploadPartCRCs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (u *uploaderDelegate) multiPart() (*UploadResult, error) {
	release := func() {
		if u.partPool != nil {
			u.partPool.Close()
		}
	}
	defer release()

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		parts     UploadParts
		errValue  atomic.Value
		crcParts  uploadPartCRCs
		enableCRC = u.client.hasFeature(FeatureEnableCRC64CheckUpload)
	)

	// Init the multipart
	uploadId, startPartNum, err := u.getUploadId()
	if err != nil {
		return nil, u.wrapErr("", err)
	}
	//fmt.Printf("getUploadId result: %v, %#v\n", uploadId, err)

	// Update Checkpoint
	if u.checkpoint != nil {
		u.checkpoint.Info.Data.UploadInfo.UploadId = uploadId
		u.checkpoint.dump()
	}

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
				upResult, err := u.client.UploadPart(
					u.context,
					&UploadPartRequest{
						Bucket:     u.request.Bucket,
						Key:        u.request.Key,
						UploadId:   Ptr(uploadId),
						PartNumber: data.partNum,
						Body:       data.body},
					u.options.ClientOptions...)

				//fmt.Printf("UploadPart result: %#v, %#v\n", upResult, err)

				if err == nil {
					mu.Lock()
					parts = append(parts, UploadPart{ETag: upResult.ETag, PartNumber: data.partNum})
					if enableCRC {
						crcParts = append(crcParts,
							uploadPartCRC{partNumber: data.partNum, hashCRC64: upResult.HashCRC64, size: data.size})
					}
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
		qnum int32 = startPartNum
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
		ch <- uploaderChunk{body: reader, partNum: qnum, cleanup: cleanup, size: nextChunkLen}
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
		cmRequest.UploadId = Ptr(uploadId)
		cmRequest.CompleteMultipartUpload = &CompleteMultipartUpload{Parts: parts}
		cmResult, err = u.client.CompleteMultipartUpload(u.context, cmRequest, u.options.ClientOptions...)
	}
	//fmt.Printf("CompleteMultipartUpload cmResult: %#v, %#v\n", cmResult, err)

	if err != nil {
		//Abort
		if !u.options.LeavePartsOnError {
			abortRequest := &AbortMultipartUploadRequest{}
			copyRequest(abortRequest, u.request)
			abortRequest.UploadId = Ptr(uploadId)
			_, _ = u.client.AbortMultipartUpload(u.context, abortRequest, u.options.ClientOptions...)
		}
		return nil, u.wrapErr(uploadId, err)
	}

	if enableCRC {
		caclCRC := fmt.Sprint(u.combineCRC(crcParts))
		if err = checkResponseHeaderCRC64(caclCRC, cmResult.Headers); err != nil {
			return nil, u.wrapErr(uploadId, err)
		}
	}

	return &UploadResult{
		UploadId:     Ptr(uploadId),
		ETag:         cmResult.ETag,
		VersionId:    cmResult.VersionId,
		HashCRC64:    cmResult.HashCRC64,
		ResultCommon: cmResult.ResultCommon,
	}, nil
}

func (u *uploaderDelegate) getUploadId() (uploadId string, startNum int32, err error) {
	if u.uploadId != "" {
		return u.uploadId, u.partNumber, nil
	}

	// if not exist or fail, create a new upload id
	request := &InitiateMultipartUploadRequest{}
	copyRequest(request, u.request)
	if request.ContentType == nil {
		request.ContentType = u.getContentType()
	}

	result, err := u.client.InitiateMultipartUpload(u.context, request, u.options.ClientOptions...)
	if err != nil {
		return "", 0, err
	}

	return *result.UploadId, 0, nil
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

func (u *uploaderDelegate) combineCRC(crcs uploadPartCRCs) uint64 {
	if len(crcs) == 0 {
		return 0
	}
	sort.Sort(crcs)
	crc := u.hashCRC64
	for _, c := range crcs {
		if c.hashCRC64 == nil {
			return 0
		}
		if value, err := strconv.ParseUint(*c.hashCRC64, 10, 64); err == nil {
			crc = CRC64Combine(crc, value, uint64(c.size))
		} else {
			break
		}
	}
	return crc
}

type DownloaderOptions struct {
	PartSize int64

	ParallelNum int

	EnableCheckpoint bool

	CheckpointDir string

	VerifyData bool

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

func (d *Downloader) DownloadSmallFile(ctx context.Context, request *GetObjectRequest, filePath string) (result *DownloadResult, err error) {
	var file *os.File
	if file, err = os.Open(filePath); err != nil {
		return
	}
	defer file.Close()
	p := progressTracker{
		pr: nil, /*request ProgressFn*/
	}
	for i := 0; i < d.client.options.RetryMaxAttempts; i++ {
		var out *GetObjectResult
		if out, err = d.client.GetObject(ctx, request); err != nil {
			return
		}
		p.total = out.ContentLength
		tr := io.TeeReader(out.Body, &p)
		var n int64
		if n, err = io.Copy(file, tr); err == nil {
			return &DownloadResult{Written: n}, nil
		}
		file.Seek(0, io.SeekStart)
		p.Reset()
	}
	return
}

func (d *Downloader) DownloadSmallFile2(ctx context.Context, request *GetObjectRequest, filePath string) (result *DownloadResult, err error) {
	var file *os.File
	requestWarp := *request
	if file, err = os.Open(filePath); err != nil {
		return
	}
	defer file.Close()
	r, _ := ParseRange(ToString(request.Range))
	reader, _ := NewRangeReader(
		ctx,
		func(ctx context.Context, httpRange HTTPRange) (output *ReaderRangeGetOutput, err error) {
			// update range
			requestWarp.Range = nil
			rangeStr := httpRange.FormatHTTPRange()
			requestWarp.RangeBehavior = nil
			if rangeStr != nil {
				requestWarp.Range = rangeStr
				requestWarp.RangeBehavior = Ptr("standard")
			}
			result, err := d.client.GetObject(ctx, &requestWarp, d.options.ClientOptions...)
			if err != nil {
				return nil, err
			}

			return &ReaderRangeGetOutput{
				Body:          result.Body,
				ETag:          result.ETag,
				ContentLength: result.ContentLength,
				ContentRange:  result.ContentRange,
			}, nil
		},
		r,
		"",
	)
	if n, err := io.Copy(file, io.TeeReader(reader, NewProgress(nil, -1))); err == nil {
		return &DownloadResult{Written: n}, nil
	} else {
		return nil, err
	}
}

func (d *Downloader) DownloadFile(ctx context.Context, request *GetObjectRequest, filePath string, optFns ...func(*DownloaderOptions)) (result *DownloadResult, err error) {
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
	var file *os.File
	if file, err = delegate.checkDestination(filePath); err != nil {
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

	// truncate to the right position
	if err = delegate.adjustWriter(file); err != nil {
		return nil, err
	}

	// CRC Part
	delegate.updateCRCFlag()

	// download
	result, err = delegate.download()

	return result, delegate.closeWriter(file, err)
}

type downloaderDelegate struct {
	options DownloaderOptions
	client  *Client
	context context.Context

	m sync.Mutex

	request *GetObjectRequest
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

	//crc
	calcCRC  bool
	checkCRC bool

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

func (d *Downloader) newDelegate(ctx context.Context, request *GetObjectRequest, optFns ...func(*DownloaderOptions)) (*downloaderDelegate, error) {
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
		request: request,
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
	copyRequest(&request, d.request)
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

func (d *downloaderDelegate) checkDestination(filePath string) (*os.File, error) {
	if filePath == "" {
		return nil, NewErrParamInvalid("filePath")
	}
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return file, nil
}

func (d *downloaderDelegate) adjustWriter(file *os.File) error {
	if err := file.Truncate(d.pos - d.rstart); err != nil {
		return err
	}
	d.w = file
	return nil
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
		d.checkpoint = newDownloadCheckpoint(d.request, d.tempFilePath, d.options.CheckpointDir, d.headers, d.options.PartSize)
		d.checkpoint.VerifyData = d.options.VerifyData
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

func (d *downloaderDelegate) updateCRCFlag() error {
	if d.client.hasFeature(FeatureEnableCRC64CheckDownload) {
		d.checkCRC = d.request.Range == nil
		d.calcCRC = (d.checkpoint != nil && d.checkpoint.VerifyData) || d.checkCRC
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
		tracker  bool   = d.calcCRC || d.checkpoint != nil
		tCRC64   uint64 = 0
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
		var hash hash.Hash64
		if d.calcCRC {
			hash = NewCRC64(0)
		}

		for {
			chunk, ok := <-ch
			if !ok {
				break
			}

			if getErrFn() != nil {
				continue
			}

			dchunk, derr := d.downloadChunk(chunk, hash)

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

	// trackerFn runs in worker goroutines to update checkpoint info or calc downloaded crc
	trackerFn := func(ch chan downloadedChunk) {
		defer cpWg.Done()
		var (
			tOffset int64 = 0
		)

		if d.checkpoint != nil {
			tOffset = d.checkpoint.Info.Data.DownloadInfo.Offset
			tCRC64 = d.checkpoint.Info.Data.DownloadInfo.CRC64
		}

		for {
			chunk, ok := <-ch
			if !ok {
				break
			}
			cpChunks = append(cpChunks, chunk)
			sort.Sort(cpChunks)
			newOffset := tOffset
			i := 0
			for ii := range cpChunks {
				if cpChunks[ii].start == newOffset {
					newOffset += cpChunks[ii].size
					i++
				} else {
					break
				}
			}
			if newOffset != tOffset {
				//remove updated chunk in cpChunks
				if d.calcCRC {
					tCRC64 = d.combineCRC(tCRC64, cpChunks[0:i])
				}
				tOffset = newOffset
				cpChunks = cpChunks[i:]
				if d.checkpoint != nil {
					d.checkpoint.Info.Data.DownloadInfo.Offset = tOffset
					d.checkpoint.Info.Data.DownloadInfo.CRC64 = tCRC64
					d.checkpoint.dump()
				}
			}
		}
	}

	// Start the download workers
	ch := make(chan downloaderChunk, d.options.ParallelNum)
	for i := 0; i < d.options.ParallelNum; i++ {
		wg.Add(1)
		go writeChunkFn(ch)
	}

	// Start tracker worker if need track downloaded chunk
	if tracker {
		cpCh = make(chan downloadedChunk, maxInt(3, d.options.ParallelNum))
		cpWg.Add(1)
		go trackerFn(cpCh)
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

	if tracker {
		close(cpCh)
		cpWg.Wait()
	}

	if d.checkCRC {
		if len(cpChunks) > 0 {
			sort.Sort(cpChunks)
		}
		if derr := checkResponseHeaderCRC64(fmt.Sprint(d.combineCRC(tCRC64, cpChunks)), d.headers); derr != nil {
			saveErrFn(derr)
		}
	}

	if err := getErrFn(); err != nil {
		return nil, d.wrapErr(err)
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

func (d *downloaderDelegate) downloadChunk(chunk downloaderChunk, hash hash.Hash64) (downloadedChunk, error) {
	// Get the next byte range of data
	var request GetObjectRequest
	copyRequest(&request, d.request)

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
	defer reader.Close()

	var (
		r     io.Reader = reader
		crc64 uint64    = 0
	)
	if hash != nil {
		hash.Reset()
		r = io.TeeReader(reader, hash)
	}

	n, err := io.Copy(&chunk, r)
	d.incrWritten(n)

	if hash != nil {
		crc64 = hash.Sum64()
	}

	return downloadedChunk{
		start: chunk.start,
		size:  n,
		crc64: crc64,
	}, err
}

func (u *downloaderDelegate) combineCRC(hashCRC uint64, crcs downloadedChunks) uint64 {
	if len(crcs) == 0 {
		return hashCRC
	}
	crc := hashCRC
	for _, c := range crcs {
		crc = CRC64Combine(crc, c.crc64, uint64(c.size))
	}
	return crc
}

func (u *downloaderDelegate) wrapErr(err error) error {
	return &DownloadError{
		Path: fmt.Sprintf("oss://%s/%s", *u.request.Bucket, *u.request.Key),
		Err:  err}
}
