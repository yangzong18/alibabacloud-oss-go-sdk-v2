package oss

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

type uploaderMockTracker struct {
	partNum        int
	saveDate       [][]byte
	checkTime      []time.Time
	timeout        []time.Duration
	uploadPartCnt  int32
	putObjectCnt   int32
	contentType    string
	uploadPartErr  []bool
	InitiateMPErr  bool
	CompleteMPErr  bool
	AbortMPErr     bool
	putObjectErr   bool
	ListPartsErr   bool
	crcPartInvalid []bool
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
				if tracker.crcPartInvalid != nil && tracker.crcPartInvalid[num-1] {
					w.Header().Set(HeaderOssCRC64, "12345")
				} else {
					w.Header().Set(HeaderOssCRC64, crc64ecma)
				}
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
		case "GET":
			if query.Get("uploadId") == "uploadId-1234" {
				// ListParts
				if tracker.ListPartsErr {
					w.Header().Set(HTTPHeaderContentType, "application/xml")
					w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
					w.WriteHeader(403)
					w.Write(errData)
					return
				}

				var buf strings.Builder
				buf.WriteString("<ListPartsResult>")
				buf.WriteString("  <Bucket>bucket</Bucket>")
				buf.WriteString("  <Key>key</Key>")
				buf.WriteString("  <UploadId>uploadId-1234</UploadId>")
				buf.WriteString("  <IsTruncated>false</IsTruncated>")
				for i, d := range tracker.saveDate {
					if d != nil {
						buf.WriteString("  <Part>")
						buf.WriteString(fmt.Sprintf("    <PartNumber>%v</PartNumber>", i+1))
						buf.WriteString("    <LastModified>2012-02-23T07:01:34.000Z</LastModified>")
						buf.WriteString("    <ETag>etag</ETag>")
						buf.WriteString(fmt.Sprintf("    <Size>%v</Size>", len(d)))
						hash := NewCRC64(0)
						hash.Write(d)
						buf.WriteString(fmt.Sprintf("    <HashCrc64ecma>%v</HashCrc64ecma>", fmt.Sprint(hash.Sum64())))
						buf.WriteString("  </Part>")
					}
				}
				buf.WriteString("</ListPartsResult>")

				data := buf.String()
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(data)))
				w.WriteHeader(200)
				w.Write([]byte(data))
			}
		}
	}))
	return server
}

func TestMockUploadSinglePart(t *testing.T) {
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

	result, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadSequential(t *testing.T) {
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

	result, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadParallel(t *testing.T) {
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

	result, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		bytes.NewReader(data))
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

func TestMockUploadArgmentCheck(t *testing.T) {
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
	_, err := u.UploadFrom(context.TODO(), nil, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request")

	_, err = u.UploadFrom(context.TODO(), &UploadRequest{
		Bucket: nil,
		Key:    Ptr("key"),
	}, nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = u.UploadFrom(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    nil,
	}, nil)
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
	assert.Contains(t, err.Error(), "File not exists,")

	// nil body
	_, err = u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "the body is null")
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
	d, err := u.newDelegate(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})

	assert.Nil(t, err)
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Equal(t, "", d.filePath)

	assert.Nil(t, d.partPool)
	assert.Nil(t, d.body)
	assert.NotNil(t, d.client)
	assert.NotNil(t, d.context)

	assert.NotNil(t, d.request)

	// empty body
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})

	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(0), d.totalSize)
	assert.Nil(t, d.body)

	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(3), d.totalSize)
	assert.NotNil(t, d.body)

	// non-empty without seek body
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = newNoSeeker(bytes.NewReader([]byte("123")))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(-1), d.totalSize)
	assert.NotNil(t, d.body)

	//file path check
	var localFile = randStr(8) + ".txt"
	createFile(t, localFile, "12345")
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	f, err := os.Open(localFile)
	assert.Nil(t, err)
	d.body = f
	err = d.applySource()
	f.Close()
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, int64(5), d.totalSize)
	assert.NotNil(t, d.body)
	os.Remove(localFile)

	// options
	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		}, func(uo *UploaderOptions) {
			uo.ParallelNum = 10
			uo.PartSize = 10
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, 10, d.options.ParallelNum)
	assert.Equal(t, int64(10), d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	// non-empty body
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		}, func(uo *UploaderOptions) {
			uo.ParallelNum = 0
			uo.PartSize = 0
		})
	assert.Nil(t, err)
	d.body = bytes.NewReader([]byte("123"))
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, DefaultUploadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultUploadParallel, u.options.ParallelNum)
	assert.Equal(t, DefaultUploadPartSize, u.options.PartSize)

	//adjust partSize
	maxSize := DefaultUploadPartSize * int64(MaxUploadParts*4)
	d, err = u.newDelegate(context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		})
	assert.Nil(t, err)
	d.body = newFakeSeeker(bytes.NewReader([]byte("123")), maxSize)
	err = d.applySource()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), d.readerPos)
	assert.Equal(t, maxSize, d.totalSize)
	assert.Equal(t, DefaultUploadPartSize*5, d.options.PartSize)
}

func TestMockUploadSinglePartFromFile(t *testing.T) {
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

func TestMockUploadSequentialFromFile(t *testing.T) {
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

func TestMockUploadParallelFromFile(t *testing.T) {
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
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadWithEmptyBody(t *testing.T) {
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

	result, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(nil))
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
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadSinglePartFail(t *testing.T) {
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

	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadSequentialInitiateMultipartUploadFail(t *testing.T) {
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

	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadSequentialUploadPartFail(t *testing.T) {
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

	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadSequentialCompleteMultipartUploadFail(t *testing.T) {
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

	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploadParallelUploadPartFail(t *testing.T) {
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

	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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

func TestMockUploaderUploadFileEnableCheckpointNotUseCp(t *testing.T) {
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
			uo.CheckpointDir = "."
			uo.EnableCheckpoint = true
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
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploaderUploadFileEnableCheckpointUseCp(t *testing.T) {
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

	localFile := "upload-file-with-cp-no-surfix"
	absPath, _ := filepath.Abs(localFile)
	hashmd5 := md5.New()
	hashmd5.Write([]byte(absPath))
	srcHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFile := srcHash + "-d36fc07f5d963b319b1b48e20a9b8ae9.ucp"

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
			uo.ParallelNum = 5
			uo.PartSize = partSize
			uo.CheckpointDir = "."
			uo.EnableCheckpoint = true
		},
	)
	assert.Equal(t, 5, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)

	// Case 1, fail in part number 4
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond
	tracker.uploadPartErr[3] = true
	os.Remove(cpFile)

	result, err := u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.Nil(t, result)
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
	assert.NotNil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime := time.Now()
	tracker.uploadPartErr[3] = false
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].Before(retryTime))
	assert.True(t, tracker.checkTime[1].Before(retryTime))
	assert.True(t, tracker.checkTime[2].Before(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr := NewMultiBytesReader(tracker.saveDate)
	all, err := io.ReadAll(mr)
	assert.Nil(t, err)

	hashall := NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma := fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(3), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)

	assert.NoFileExists(t, cpFile)

	// Case 2, fail in part number 1
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.timeout[0] = 1 * time.Second
	tracker.timeout[2] = 500 * time.Millisecond
	tracker.uploadPartErr[0] = true
	os.Remove(cpFile)

	result, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.Nil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.NotNil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime = time.Now()
	tracker.uploadPartErr[0] = false
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].After(retryTime))
	assert.True(t, tracker.checkTime[1].After(retryTime))
	assert.True(t, tracker.checkTime[2].After(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr = NewMultiBytesReader(tracker.saveDate)
	all, err = io.ReadAll(mr)
	assert.Nil(t, err)

	hashall = NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma = fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(6), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.NoFileExists(t, cpFile)

	// Case 3, list Parts Fail
	tracker.saveDate = make([][]byte, partsNum)
	tracker.checkTime = make([]time.Time, partsNum)
	tracker.timeout = make([]time.Duration, partsNum)
	tracker.uploadPartErr = make([]bool, partsNum)
	tracker.uploadPartErr[3] = true
	os.Remove(cpFile)

	result, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.NotNil(t, err)
	assert.NotNil(t, tracker.saveDate[0])
	assert.NotNil(t, tracker.saveDate[1])
	assert.NotNil(t, tracker.saveDate[2])
	assert.Nil(t, tracker.saveDate[3])
	assert.NotNil(t, tracker.saveDate[4])

	assert.FileExists(t, cpFile)

	//retry
	time.Sleep(2 * time.Second)
	retryTime = time.Now()
	tracker.uploadPartErr[3] = false
	tracker.ListPartsErr = true
	atomic.StoreInt32(&tracker.uploadPartCnt, 0)

	result, err = u.UploadFile(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, localFile)

	assert.Nil(t, err)
	assert.NotNil(t, result)

	assert.True(t, tracker.checkTime[0].After(retryTime))
	assert.True(t, tracker.checkTime[1].After(retryTime))
	assert.True(t, tracker.checkTime[2].After(retryTime))
	assert.True(t, tracker.checkTime[3].After(retryTime))
	assert.True(t, tracker.checkTime[4].After(retryTime))
	assert.True(t, tracker.checkTime[5].After(retryTime))

	mr = NewMultiBytesReader(tracker.saveDate)
	all, err = io.ReadAll(mr)
	assert.Nil(t, err)

	hashall = NewCRC64(0)
	hashall.Write(all)
	allCrc64ecma = fmt.Sprint(hashall.Sum64())
	assert.Equal(t, dataCrc64ecma, allCrc64ecma)

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(6), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
	assert.NoFileExists(t, cpFile)
}

func TestMockUploadParallelFromStreamWithoutSeeker(t *testing.T) {
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

	file, err := os.Open(localFile)
	assert.Nil(t, err)
	defer file.Close()

	result, err := u.UploadFrom(context.TODO(), &UploadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}, io.LimitReader(file, int64(length)))
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

	assert.Equal(t, int32(0), atomic.LoadInt32(&tracker.putObjectCnt))
	assert.Equal(t, int32(partsNum), atomic.LoadInt32(&tracker.uploadPartCnt))
	//FeatureAutoDetectMimeType is enabled default
	assert.Equal(t, "application/octet-stream", tracker.contentType)
}

func TestMockUploadCRC64Fail(t *testing.T) {
	partSize := int64(100 * 1024)
	length := 5*100*1024 + 123
	partsNum := length/int(partSize) + 1
	tracker := &uploaderMockTracker{
		partNum:        partsNum,
		saveDate:       make([][]byte, partsNum),
		checkTime:      make([]time.Time, partsNum),
		timeout:        make([]time.Duration, partsNum),
		uploadPartErr:  make([]bool, partsNum),
		crcPartInvalid: make([]bool, partsNum),
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
	tracker.crcPartInvalid[2] = true
	_, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "crc is inconsistent")

	//disable crc check
	client = NewClient(cfg,
		func(o *Options) {
			o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
		})

	u = client.NewUploader(
		func(uo *UploaderOptions) {
			uo.ParallelNum = 1
			uo.PartSize = partSize
		},
	)
	assert.Equal(t, 1, u.options.ParallelNum)
	assert.Equal(t, partSize, u.options.PartSize)
	tracker.crcPartInvalid[2] = true
	tracker.saveDate = make([][]byte, partsNum)
	result, err := u.UploadFrom(
		context.TODO(),
		&UploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key")},
		bytes.NewReader(data))
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
}

type downloaderMockTracker struct {
	lastModified string
	data         []byte

	maxRangeCount int64
	getPartCnt    int32

	etagChangeOffset int64
	failPartNum      int32

	partSize     int32
	gotMinOffset int64

	rStart int64
}

func testSetupDownloaderMockServer(t *testing.T, tracker *downloaderMockTracker) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		length := len(tracker.data)
		data := tracker.data
		errData := []byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>ak</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`)

		switch r.Method {
		case "HEAD":
			tracker.gotMinOffset = int64(length)
			// header
			w.Header().Set(HTTPHeaderLastModified, tracker.lastModified)
			w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(length))
			w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
			w.Header().Set(HTTPHeaderContentType, "text/plain")

			//status code
			w.WriteHeader(200)

			//body
			w.Write(nil)
		case "GET":
			// header
			var httpRange *HTTPRange
			if r.Header.Get("Range") != "" {
				httpRange, _ = ParseRange(r.Header.Get("Range"))
			}

			offset := int64(0)
			statusCode := 200
			sendLen := int64(length)
			if httpRange != nil {
				offset = httpRange.Offset
				sendLen = int64(length) - httpRange.Offset
				if httpRange.Count > 0 {
					sendLen = minInt64(httpRange.Count, sendLen)
					tracker.maxRangeCount = maxInt64(httpRange.Count, tracker.maxRangeCount)
				}
				cr := httpContentRange{
					Offset: httpRange.Offset,
					Count:  sendLen,
					Total:  int64(length),
				}
				w.Header().Set("Content-Range", ToString(cr.FormatHTTPContentRange()))
				statusCode = 206
			}

			tracker.gotMinOffset = minInt64(tracker.gotMinOffset, offset)

			if tracker.failPartNum > 0 && (int64(tracker.partSize*tracker.failPartNum)+tracker.rStart) == offset {
				w.Header().Set(HTTPHeaderContentType, "application/xml")
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(len(errData)))
				w.WriteHeader(403)
				w.Write(errData)
			} else {
				w.Header().Set(HTTPHeaderContentLength, fmt.Sprint(sendLen))
				w.Header().Set(HTTPHeaderLastModified, tracker.lastModified)
				if tracker.etagChangeOffset > 0 && offset > 0 && offset > tracker.etagChangeOffset {
					w.Header().Set(HTTPHeaderETag, "2ba9dede5f27731c9771645a3986****")
				} else {
					w.Header().Set(HTTPHeaderETag, "fba9dede5f27731c9771645a3986****")
				}
				w.Header().Set(HTTPHeaderContentType, "text/plain")

				//status code
				w.WriteHeader(statusCode)

				//body
				sendData := data[int(offset):int(offset+sendLen)]
				//fmt.Printf("sendData offset%d, len:%d, total:%d\n", offset, len(sendData), length)
				w.Write(sendData)
			}
		}
	}))
	return server
}

func TestMockDownloaderSingleRead(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()
	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderLoopSingleRead(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for i := 1; i <= 20; i++ {
		if FileExists(localFile) {
			assert.Nil(t, os.Remove(localFile))
		}
		tracker.maxRangeCount = 0
		result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
			func(do *DownloaderOptions) {
				do.ParallelNum = 1
				do.PartSize = int64(i)
			})
		assert.Nil(t, err)
		assert.Equal(t, int64(length), result.Written)
		hash := NewCRC64(0)
		rfile, err := os.Open(localFile)
		assert.Nil(t, err)
		io.Copy(hash, rfile)
		rfile.Close()
		os.Remove(localFile)
		assert.Equal(t, datasum, hash.Sum64())
		assert.Equal(t, int64(i), tracker.maxRangeCount)
	}
}

func TestMockDownloaderLoopSingleReadWithRange(t *testing.T) {
	length := 63
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for rs := 0; rs < 7; rs++ {
		for rcount := 1; rcount < length; rcount++ {
			for i := 1; i <= 3; i++ {
				//fmt.Printf("rs:%v, rcount:%v, i:%v\n", rs, rcount, i)
				if FileExists(localFile) {
					assert.Nil(t, os.Remove(localFile))
				}
				tracker.maxRangeCount = 0
				httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
				result, err := d.DownloadFile(context.TODO(),
					&DownloadRequest{
						Bucket: Ptr("bucket"),
						Key:    Ptr("key"),
						Range:  httpRange.FormatHTTPRange()},
					localFile,
					func(do *DownloaderOptions) {
						do.ParallelNum = 1
						do.PartSize = int64(i)
					})
				assert.Nil(t, err)
				expectLen := minInt64(int64(length-rs), int64(rcount))
				assert.Equal(t, expectLen, result.Written)
				hash := NewCRC64(0)
				rfile, err := os.Open(localFile)
				assert.Nil(t, err)
				io.Copy(hash, rfile)
				rfile.Close()
				//ldata, err := os.ReadFile(localFile)
				//assert.Nil(t, err)
				os.Remove(localFile)
				hdata := NewCRC64(0)
				pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
				hdata.Write(pat)
				//assert.EqualValues(t, ldata, pat)
				assert.Equal(t, hdata.Sum64(), hash.Sum64())
				assert.Equal(t, minInt64(int64(i), expectLen), tracker.maxRangeCount)
			}
		}
	}
}

func TestMockDownloaderParalleRead(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 3
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 3, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()

	io.Copy(hash, rfile)
	assert.Equal(t, datasum, hash.Sum64())
}

func TestMockDownloaderLoopParalleRead(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for i := 1; i <= 20; i++ {
		if FileExists(localFile) {
			assert.Nil(t, os.Remove(localFile))
		}
		tracker.maxRangeCount = 0
		result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
			func(do *DownloaderOptions) {
				do.ParallelNum = 4
				do.PartSize = int64(i)
			})
		assert.Nil(t, err)
		assert.Equal(t, int64(length), result.Written)
		hash := NewCRC64(0)
		rfile, err := os.Open(localFile)
		assert.Nil(t, err)
		io.Copy(hash, rfile)
		assert.Nil(t, rfile.Close())
		assert.Nil(t, os.Remove(localFile))
		assert.Equal(t, datasum, hash.Sum64())
		assert.Equal(t, int64(i), tracker.maxRangeCount)
	}
}

func TestMockDownloaderLoopParalleReadWithRange(t *testing.T) {
	length := 63
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	for rs := 0; rs < 7; rs++ {
		for rcount := 1; rcount < length; rcount++ {
			for i := 1; i <= 3; i++ {
				//fmt.Printf("rs:%v, rcount:%v, i:%v\n", rs, rcount, i)
				if FileExists(localFile) {
					assert.Nil(t, os.Remove(localFile))
				}
				tracker.maxRangeCount = 0
				httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
				result, err := d.DownloadFile(context.TODO(),
					&DownloadRequest{
						Bucket: Ptr("bucket"),
						Key:    Ptr("key"),
						Range:  httpRange.FormatHTTPRange()},
					localFile,
					func(do *DownloaderOptions) {
						do.ParallelNum = 3
						do.PartSize = int64(i)
					})
				assert.Nil(t, err)
				expectLen := minInt64(int64(length-rs), int64(rcount))
				assert.Equal(t, expectLen, result.Written)
				hash := NewCRC64(0)
				rfile, err := os.Open(localFile)
				assert.Nil(t, err)
				io.Copy(hash, rfile)
				rfile.Close()
				//ldata, err := os.ReadFile(localFile)
				//assert.Nil(t, err)
				os.Remove(localFile)
				hdata := NewCRC64(0)
				pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
				hdata.Write(pat)
				//assert.EqualValues(t, ldata, pat)
				assert.Equal(t, hdata.Sum64(), hash.Sum64())
				assert.Equal(t, minInt64(int64(i), expectLen), tracker.maxRangeCount)
			}
		}
	}
}

func TestDownloaderConstruct(t *testing.T) {
	c := &Client{}
	d := c.NewDownloader()
	assert.Equal(t, DefaultDownloadParallel, d.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, d.options.PartSize)
	assert.True(t, d.options.UseTempFile)
	assert.False(t, d.options.EnableCheckpoint)
	assert.Equal(t, "", d.options.CheckpointDir)

	d = c.NewDownloader(func(do *DownloaderOptions) {
		do.CheckpointDir = "dir"
		do.EnableCheckpoint = true
		do.ParallelNum = 1
		do.PartSize = 2
		do.UseTempFile = false
	})
	assert.Equal(t, 1, d.options.ParallelNum)
	assert.Equal(t, int64(2), d.options.PartSize)
	assert.False(t, d.options.UseTempFile)
	assert.True(t, d.options.EnableCheckpoint)
	assert.Equal(t, "dir", d.options.CheckpointDir)
}

func TestDownloaderDelegateConstruct(t *testing.T) {
	c := &Client{}
	d := c.NewDownloader()

	_, err := d.newDelegate(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")

	_, err = d.newDelegate(context.TODO(), &DownloadRequest{})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = d.newDelegate(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket")})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Key")

	delegate, err := d.newDelegate(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)
	assert.True(t, delegate.options.UseTempFile)
	assert.False(t, delegate.options.EnableCheckpoint)
	assert.Empty(t, delegate.options.CheckpointDir)

	delegate, err = d.newDelegate(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = 5
			do.PartSize = 1
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, 5, delegate.options.ParallelNum)
	assert.Equal(t, int64(1), delegate.options.PartSize)

	delegate, err = d.newDelegate(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = 0
			do.PartSize = 0
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)

	delegate, err = d.newDelegate(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")},
		func(do *DownloaderOptions) {
			do.ParallelNum = -1
			do.PartSize = -1
			do.CheckpointDir = "dir"
			do.EnableCheckpoint = true
			do.UseTempFile = false
		})
	assert.Nil(t, err)
	assert.NotNil(t, delegate)
	assert.Equal(t, DefaultDownloadParallel, delegate.options.ParallelNum)
	assert.Equal(t, DefaultDownloadPartSize, delegate.options.PartSize)
	assert.False(t, delegate.options.UseTempFile)
	assert.True(t, delegate.options.EnableCheckpoint)
	assert.Equal(t, "dir", delegate.options.CheckpointDir)
}

func TestDownloaderDownloadFileArgument(t *testing.T) {
	c := &Client{}
	d := c.NewDownloader()

	_, err := d.DownloadFile(context.TODO(), nil, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field")

	_, err = d.DownloadFile(context.TODO(), &DownloadRequest{}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Bucket")

	_, err = d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket")}, "file")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request.Key")

	_, err = d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "operation error HeadObject")
}

func TestMockDownloaderDownloadFileWithoutTempFile(t *testing.T) {
	length := 3*int(DefaultDownloadPartSize) + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"

	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.UseTempFile = false
			do.PartSize = 1024
			do.ParallelNum = 2
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	defer func() {
		rfile.Close()
		os.Remove(localFile)
	}()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, int64(1024), tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileInvalidPartSizeAndParallelNum(t *testing.T) {
	length := int(DefaultDownloadPartSize*2) + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = 0
			do.ParallelNum = 0
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, DefaultDownloadPartSize, tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileFileSizeLessPartSize(t *testing.T) {
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = 0
			do.ParallelNum = 0
		})
	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, int64(length), tracker.maxRangeCount)
}

func TestMockDownloaderDownloadFileFileChange(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified:     gmtTime,
		data:             data,
		etagChangeOffset: 700,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
		})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Source file is changed")
	assert.False(t, FileExists(localFile))
	assert.False(t, FileExists(localFile+TempFileSuffix))
}

func TestMockDownloaderDownloadFileEnableCheckpointNormal(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()

	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := randStr(8) + "-no-surfix"
	defer func() {
		os.Remove(localFile)
	}()

	_, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})
	assert.Nil(t, err)
}

func TestMockDownloaderDownloadFileEnableCheckpoint2(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	datasum := func() uint64 {
		h := NewCRC64(0)
		h.Write(data)
		return h.Sum64()
	}()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		failPartNum:  6,
		partSize:     int32(partSize),
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := "check-point-to-check-no-surfix"
	localFileTmep := "check-point-to-check-no-surfix" + TempFileSuffix
	absPath, _ := filepath.Abs(localFileTmep)
	hashmd5 := md5.New()
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFileName := "ddaf063c8f69766ecc8e4a93b6402e3e-" + destHash + ".dcp"
	defer func() {
		os.Remove(localFile)
	}()
	os.Remove(localFile)
	os.Remove(localFileTmep)
	os.Remove(cpFileName)
	_, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.NotNil(t, err)
	assert.True(t, FileExists(localFileTmep))
	assert.True(t, FileExists(cpFileName))

	//load CheckPointFile
	content, err := os.ReadFile(cpFileName)
	assert.Nil(t, err)
	dcp := downloadCheckpoint{}
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "fba9dede5f27731c9771645a3986****", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, gmtTime, dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(length), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.Range)

	abslocalFileTmep, _ := filepath.Abs(localFileTmep)
	assert.Equal(t, abslocalFileTmep, dcp.Info.Data.FilePath)
	assert.Equal(t, int64(partSize), dcp.Info.Data.PartSize)

	assert.Equal(t, int64(tracker.failPartNum*tracker.partSize), dcp.Info.Data.DownloadInfo.Offset)
	h := NewCRC64(0)
	h.Write(data[0:int(dcp.Info.Data.DownloadInfo.Offset)])
	assert.Equal(t, h.Sum64(), dcp.Info.Data.DownloadInfo.CRC64)

	// resume from checkpoint
	tracker.failPartNum = 0
	result, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
			do.VerifyData = true
		})

	assert.Nil(t, err)
	assert.Equal(t, int64(length), result.Written)

	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	io.Copy(hash, rfile)
	rfile.Close()
	assert.Equal(t, datasum, hash.Sum64())
	assert.Equal(t, dcp.Info.Data.DownloadInfo.Offset, tracker.gotMinOffset)
}

func TestMockDownloaderDownloadFileEnableCheckpointWithRange(t *testing.T) {
	partSize := 128
	length := 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	rs := 5
	rcount := 832
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
		failPartNum:  6,
		partSize:     int32(partSize),
		rStart:       int64(rs),
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader()
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)

	localFile := "check-point-to-check-no-surfix"
	localFileTmep := "check-point-to-check-no-surfix" + TempFileSuffix
	absPath, _ := filepath.Abs(localFileTmep)
	hashmd5 := md5.New()
	hashmd5.Reset()
	hashmd5.Write([]byte(absPath))
	destHash := hex.EncodeToString(hashmd5.Sum(nil))
	cpFileName := "0fbbf3bb7c80debbecb37dca52a646eb-" + destHash + ".dcp"
	defer func() {
		os.Remove(localFile)
	}()
	os.Remove(localFile)
	os.Remove(localFileTmep)
	os.Remove(cpFileName)
	httpRange := HTTPRange{Offset: int64(rs), Count: int64(rcount)}
	_, err := d.DownloadFile(context.TODO(),
		&DownloadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Range:  httpRange.FormatHTTPRange()},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.NotNil(t, err)
	assert.True(t, FileExists(localFileTmep))
	assert.True(t, FileExists(cpFileName))

	//load CheckPointFile
	content, err := os.ReadFile(cpFileName)
	assert.Nil(t, err)
	dcp := downloadCheckpoint{}
	err = json.Unmarshal(content, &dcp.Info)
	assert.Nil(t, err)

	assert.Equal(t, "fba9dede5f27731c9771645a3986****", dcp.Info.Data.ObjectMeta.ETag)
	assert.Equal(t, gmtTime, dcp.Info.Data.ObjectMeta.LastModified)
	assert.Equal(t, int64(length), dcp.Info.Data.ObjectMeta.Size)

	assert.Equal(t, "oss://bucket/key", dcp.Info.Data.ObjectInfo.Name)
	assert.Equal(t, "", dcp.Info.Data.ObjectInfo.VersionId)
	assert.Equal(t, ToString(httpRange.FormatHTTPRange()), dcp.Info.Data.ObjectInfo.Range)

	abslocalFileTmep, _ := filepath.Abs(localFileTmep)
	assert.Equal(t, abslocalFileTmep, dcp.Info.Data.FilePath)
	assert.Equal(t, int64(partSize), dcp.Info.Data.PartSize)

	assert.Equal(t, int64(tracker.failPartNum*tracker.partSize)+int64(rs), dcp.Info.Data.DownloadInfo.Offset)
	assert.Equal(t, uint64(0), dcp.Info.Data.DownloadInfo.CRC64)

	// resume from checkpoint
	tracker.failPartNum = 0
	result, err := d.DownloadFile(context.TODO(),
		&DownloadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Range:  httpRange.FormatHTTPRange()},
		localFile,
		func(do *DownloaderOptions) {
			do.PartSize = int64(partSize)
			do.ParallelNum = 3
			do.CheckpointDir = "."
			do.EnableCheckpoint = true
		})

	assert.Nil(t, err)
	expectLen := minInt64(int64(length-rs), int64(rcount))
	assert.Equal(t, expectLen, result.Written)
	hash := NewCRC64(0)
	rfile, err := os.Open(localFile)
	assert.Nil(t, err)
	io.Copy(hash, rfile)
	rfile.Close()
	os.Remove(localFile)
	hdata := NewCRC64(0)
	pat := data[rs:int(minInt64(int64(rs+rcount), int64(length)))]
	hdata.Write(pat)
	assert.Equal(t, hdata.Sum64(), hash.Sum64())
	assert.Equal(t, dcp.Info.Data.DownloadInfo.Offset, tracker.gotMinOffset)
}

func TestMockDownloaderDownloadWithError(t *testing.T) {
	length := 3*1024*1024 + 1234
	data := []byte(randStr(length))
	gmtTime := getNowGMT()
	tracker := &downloaderMockTracker{
		lastModified: gmtTime,
		data:         data,
	}
	server := testSetupDownloaderMockServer(t, tracker)
	defer server.Close()
	assert.NotNil(t, server)

	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
		WithRegion("cn-hangzhou").
		WithEndpoint(server.URL).
		WithReadWriteTimeout(300 * time.Second)

	client := NewClient(cfg)
	d := client.NewDownloader(func(do *DownloaderOptions) {
		do.ParallelNum = 1
		do.PartSize = 1 * 1024 * 1024
	})
	assert.NotNil(t, d)
	assert.NotNil(t, d.client)
	assert.Equal(t, int64(1*1024*1024), d.options.PartSize)
	assert.Equal(t, 1, d.options.ParallelNum)

	// filePath is invalid
	_, err := d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, filePath")

	localFile := "./no-exist-folder/file-no-surfix"
	_, err = d.DownloadFile(context.TODO(), &DownloadRequest{Bucket: Ptr("bucket"), Key: Ptr("key")}, localFile)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "The system cannot find the path specified")

	// Range is invalid
	localFile = randStr(8) + "-no-surfix"
	_, err = d.DownloadFile(context.TODO(), &DownloadRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Range:  Ptr("invalid range")},
		localFile)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid field, request.Range")
}

func TestDownloadedChunksSort(t *testing.T) {
	chunks := downloadedChunks{}
	chunks = append(chunks, downloadedChunk{start: 0, size: 10})
	chunks = append(chunks, downloadedChunk{start: 10, size: 5})
	chunks = append(chunks, downloadedChunk{start: 15, size: 30})
	chunks = append(chunks, downloadedChunk{start: 45, size: 1})
	sort.Sort(chunks)

	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, int64(0), chunks[0].start)
	assert.Equal(t, int64(10), chunks[1].start)
	assert.Equal(t, int64(15), chunks[2].start)
	assert.Equal(t, int64(45), chunks[3].start)

	chunks = downloadedChunks{}
	chunks = append(chunks, downloadedChunk{start: 10, size: 5})
	chunks = append(chunks, downloadedChunk{start: 0, size: 10})
	chunks = append(chunks, downloadedChunk{start: 45, size: 1})
	chunks = append(chunks, downloadedChunk{start: 15, size: 30})

	assert.Equal(t, 4, len(chunks))
	assert.Equal(t, int64(10), chunks[0].start)
	assert.Equal(t, int64(0), chunks[1].start)
	assert.Equal(t, int64(45), chunks[2].start)
	assert.Equal(t, int64(15), chunks[3].start)

	sort.Sort(chunks)

	assert.Equal(t, int64(0), chunks[0].start)
	assert.Equal(t, int64(10), chunks[1].start)
	assert.Equal(t, int64(15), chunks[2].start)
	assert.Equal(t, int64(45), chunks[3].start)
}
