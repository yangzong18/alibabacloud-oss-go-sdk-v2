package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

type uploaderMockTracker struct {
	partNum       int
	saveDate      [][]byte
	checkTime     []time.Time
	timeout       []time.Duration
	uploadPartCnt int32
	putObjectCnt  int32
	contentType   string
	uploadPartErr []bool
	InitiateMPErr bool
	CompleteMPErr bool
	AbortMPErr    bool
	putObjectErr  bool
}

func testSetupUploaderMockServer(t *testing.T, tracker *uploaderMockTracker) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("r.URL :%s\n", r.URL.String())
		errData := []byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>ak</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`)

		query := r.URL.Query()
		switch r.Method {
		case "POST":
			//url := r.URL.String()
			//strings.Contains(url, "/bucket/key?uploads")
			if query.Get("uploads") == "" && query.Get("uploadId") == "" {
				// InitiateMultipartUpload
				if tracker.InitiateMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				sendData := []byte(`
				<InitiateMultipartUploadResult>
					<Bucket>bucket</Bucket>
					<Key>key</Key>
					<UploadId>uploadId-1234</UploadId>
				</InitiateMultipartUploadResult>`)

				tracker.contentType = r.Header.Get(HTTPHeaderContentType)

				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(sendData)))
				w.WriteHeader(200)
				w.Write(sendData)
			} else if query.Get("uploadId") == "uploadId-1234" {
				// strings.Contains(url, "/bucket/key?uploadId=uploadId-1234")
				// CompleteMultipartUpload
				if tracker.CompleteMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				sendData := []byte(`
				<CompleteMultipartUploadResult>
					<EncodingType>url</EncodingType>
					<Location>bucket/key</Location>
					<Bucket>bucket</Bucket>
					<Key>key</Key>
					<ETag>etag</ETag>
			  	</CompleteMultipartUploadResult>`)

				hash := NewCRC64(0)
				mr := NewMultiBytesReader(tracker.saveDate)
				io.Copy(io.MultiWriter(hash), mr)
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(sendData)))
				crc64ecma := fmt.Sprint(hash.Sum64())
				w.Header().Set(HeaderOssCRC64, crc64ecma)
				w.WriteHeader(200)
				w.Write(sendData)
			} else {
				assert.Fail(t, "not support")
			}
		case "PUT":
			in, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			hash := NewCRC64(0)
			hash.Write(in)
			crc64ecma := fmt.Sprint(hash.Sum64())

			md5hash := md5.New()
			md5hash.Write(in)
			etag := fmt.Sprintf("\"%s\"", strings.ToUpper(hex.EncodeToString(md5hash.Sum(nil))))

			if query.Get("uploadId") == "uploadId-1234" {
				// UploadPart
				//in, err := io.ReadAll(r.Body)
				//assert.Nil(t, err)
				num, err := strconv.Atoi(query.Get("partNumber"))
				assert.Nil(t, err)
				assert.LessOrEqual(t, num, tracker.partNum)
				assert.Nil(t, err)
				assert.Equal(t, "uploadId-1234", query.Get("uploadId"))

				//hash := NewCRC64(0)
				//hash.Write(in)
				//crc64ecma := fmt.Sprint(hash.Sum64())
				if tracker.timeout[num-1] > 0 {
					time.Sleep(tracker.timeout[num-1])
				} else {
					time.Sleep(10 * time.Millisecond)
				}

				if tracker.uploadPartErr[num-1] {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				tracker.saveDate[num-1] = in

				// header
				w.Header().Set(HeaderOssCRC64, crc64ecma)
				w.Header().Set(HTTPHeaderETag, etag)

				//status code
				w.WriteHeader(200)

				//body
				w.Write(nil)
				tracker.checkTime[num-1] = time.Now()
				//fmt.Printf("UploadPart done, num :%d, %v\n", num, tracker.checkTime[num-1])
				atomic.AddInt32(&tracker.uploadPartCnt, 1)
			} else if query.Get("uploadId") == "" {
				tracker.contentType = r.Header.Get(HTTPHeaderContentType)

				if tracker.putObjectErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				//PutObject
				w.Header().Set(HeaderOssCRC64, crc64ecma)
				w.Header().Set(HTTPHeaderETag, etag)

				//status code
				w.WriteHeader(200)

				//body
				w.Write(nil)
				tracker.saveDate[0] = in
				tracker.checkTime[0] = time.Now()
				atomic.AddInt32(&tracker.putObjectCnt, 1)
			} else {
				assert.Fail(t, "not support")
			}
		case "DELETE":
			if query.Get("uploadId") == "uploadId-1234" {
				// AbortMultipartUpload
				if tracker.AbortMPErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				w.WriteHeader(204)
				w.Write(nil)
			} else {
				assert.Fail(t, "not support")
			}
		}
	}))
	return server
}

func TestUploadSinglePart(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := client.NewUploader()

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	result, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestUploadSequential(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	result, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, partsNum-1, index)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestUploadParallel(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, 0, index)
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
}

func TestUploadArgmentCheck(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)
	u := client.NewUploader()
	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// upload stream
	_, err := u.Upload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request")

	_, err = u.Upload(context.TODO(), &UploadRequest{
		Bucket: nil,
		Key:    Ptr("key"),
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    nil,
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Key")

	// upload file
	_, err = u.UploadFile(context.TODO(), nil, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request")

	_, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: nil,
		Key:    Ptr("key"),
	}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    nil,
	}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Key")

	//Invalid filePath
	_, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, "#@!Ainvalud-path")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The system cannot find the file specified.")
}

type noSeeker struct {
	io.Reader
}

func newNoSeeker(r io.Reader) noSeeker {
	return noSeeker{r}
}

type fakeSeeker struct {
	r io.Reader
	n int64
	i int64
}

func newFakeSeeker(r io.Reader, n int64) fakeSeeker {
	return fakeSeeker{r: r, n: n, i: 0}
}

func (r fakeSeeker) Read(p []byte) (n int, err error) {
	return r.Read(p)
}

func (r fakeSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.i + offset
	case io.SeekEnd:
		abs = r.n + offset
	default:
		return 0, errors.New("MultiSliceReader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("MultiSliceReader.Seek: negative position")
	}
	r.i = abs
	return abs, nil
}

func createFile(t *testing.T, fileName, content string) {
	fout, err := os.Create(fileName)
	assert.Nil(t, err)
	defer fout.Close()
	_, err = fout.WriteString(content)
	assert.Nil(t, err)
}

func createFileFromByte(t *testing.T, fileName string, content []byte) {
	fout, err := os.Create(fileName)
	assert.Nil(t, err)
	defer fout.Close()
	_, err = fout.Write(content)
	assert.Nil(t, err)
}

func TestUpload_newDelegate(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)
	u := client.NewUploader()
	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// nil body
	d, err := u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, "")

	assert.Nil(t, err)
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Equal(t, "", d.filePath)

	assert.Nil(t, d.partPool)
	assert.Nil(t, d.context)
	assert.Nil(t, d.body)
	assert.Nil(t, d.client)

	assert.NotNil(t, d.request)

	// empty body
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader([]byte{}),
		},
	}, "")

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Nil(t, d.body)

	// non-empty body
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader([]byte("123")),
		},
	}, "")

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(3), d.totalSize)
	assert.Nil(t, d.body)

	// non-empty without seek body
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: newNoSeeker(bytes.NewReader([]byte("123"))),
		},
	}, "")

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(-1), d.totalSize)
	assert.Nil(t, d.body)

	//file path check
	var localFile = randStr(8) + ".txt"
	createFile(t, localFile, "12345")
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(5), d.totalSize)
	assert.Nil(t, d.body)
	os.Remove(localFile)

	// options
	// non-empty body
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader([]byte("123")),
		},
	}, "", func(uo *UploaderOptions) {
		uo.ParallelNum = 10
		uo.PartSize = 10
	})
	assert.Equal(t, 10, d.options.ParallelNum)
	assert.Equal(t, int64(10), d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// non-empty body
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader([]byte("123")),
		},
	}, "", func(uo *UploaderOptions) {
		uo.ParallelNum = 0
		uo.PartSize = 0
	})
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	//adjust partSize
	maxSize := DefaultUploadPartSize * int64(MaxUploadParts*4)
	d, err = u.newDelegate(&UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: newFakeSeeker(bytes.NewReader([]byte("123")), maxSize),
		},
	}, "")
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, maxSize, d.totalSize)
	assert.Equal(t, DefaultUploadPartSize*5, d.options.PartSize)
}

func TestUploadSinglePartFromFile(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + ".txt"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := client.NewUploader()

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	result, err := u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "text/plain", tracker.contentType)
}

func TestUploadSequentialFromFile(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + ".tif"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	result, err := u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, partsNum-1, index)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "image/tiff", tracker.contentType)
}

func TestUploadParallelFromFile(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	localFile := randStr(8) + "-no-surfix"
	createFileFromByte(t, localFile, data)
	defer func() {
		os.Remove(localFile)
	}()

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "uploadId-1234", *result.UploadId)
	assert.Equal(t, dataCrc64ecma, *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	index := 3
	ctime := tracker.checkTime[index]
	for i, t := range tracker.checkTime {
		if t.After(ctime) {
			index = i
			ctime = t
		}
	}
	assert.Equal(t, 0, index)
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "", tracker.contentType)
}

func TestUploadWithNullBody(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond

	result, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Nil(t, result.UploadId)
	assert.Equal(t, "0", *result.HashCRC64)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(all))
	assert.Equal(t, int32(1), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.uploadPartCnt))
	assert.Equal(t, "", tracker.contentType)
}

func TestUploadSinglePartFail(t *testing.T) {
	partSize := DefaultUploadPartSize
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
		putObjectErr:  true,
	}

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := client.NewUploader()

	assert.NotNil(t, u.client)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	_, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestUploadSequentialInitiateMultipartUploadFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
		InitiateMPErr: true,
	}

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 4
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 4, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestUploadSequentialUploadPartFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.uploadPartErr[1] = true

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
}

func TestUploadSequentialCompleteMultipartUploadFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.CompleteMPErr = true

	data := []byte(randStr(length))
	hash := NewCRC64(0)
	hash.Write(data)
	dataCrc64ecma := fmt.Sprint(hash.Sum64())

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)
}

func TestUploadParallelUploadPartFail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:       partsNum,
		saveDate:      make([][]byte, partsNum),
		checkTime:     make([]time.Time, partsNum),
		timeout:       make([]time.Duration, partsNum),
		uploadPartErr: make([]bool, partsNum),
	}
	tracker.uploadPartErr[2] = true

	data := []byte(randStr(length))

	server := testSetupUploaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)

	u := client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 2
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 2, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	_, err := u.Upload(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Body: bytes.NewReader(data),
		},
	})
	assert.NotNil(t, err)
	var uerr *UploadError
	errors.As(err, &uerr)
	assert.NotNil(t, uerr)
	assert.Equal(t, "uploadId-1234", uerr.UploadId)
	assert.Equal(t, "oss://bucket/key", uerr.Path)

	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)

	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.Nil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[5])
}
