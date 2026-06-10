//go:build integration

package oss

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitiateMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, initResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *initResult.Bucket, bucketName)
	assert.Equal(t, *initResult.Key, objectName)
	assert.NotEmpty(t, *initResult.UploadId)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.InitiateMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	initRequest = &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPart(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		Body:         strings.NewReader("upload part 1"),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, partResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *partResult.ETag)
	assert.NotEmpty(t, *partResult.ContentMD5)
	assert.NotEmpty(t, *partResult.HashCRC64)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.UploadPart(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	partRequest = &UploadPartRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		PartNumber:   int32(2),
		UploadId:     Ptr(*initResult.UploadId),
		Body:         strings.NewReader("upload part 2"),
		TrafficLimit: int64(100 * 1024 * 8),
	}

	_, err = client.UploadPart(context.TODO(), partRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000104", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPartCopy(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	body := randLowStr(100000)
	objectSrcName := objectNamePrefix + randLowStr(6) + "src"
	objRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectSrcName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), objRequest)
	assert.Nil(t, err)
	objectDestName := objectNamePrefix + randLowStr(6) + "dest"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectDestName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		SourceKey:    Ptr(objectSrcName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, copyResult.StatusCode)
	assert.NotEmpty(t, copyResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *copyResult.ETag)
	assert.NotEmpty(t, *copyResult.LastModified)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	metaRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectSrcName),
	}
	metaResult, err := client.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)
	sourceVersionId := *metaResult.VersionId

	copyRequest = &UploadPartCopyRequest{
		Bucket:          Ptr(bucketName),
		Key:             Ptr(objectDestName),
		PartNumber:      int32(1),
		UploadId:        Ptr(*initResult.UploadId),
		SourceKey:       Ptr(objectSrcName),
		SourceVersionId: Ptr(sourceVersionId),
	}
	copyResult, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, copyResult.StatusCode)
	assert.NotEmpty(t, copyResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *copyResult.ETag)
	assert.NotEmpty(t, *copyResult.LastModified)
	assert.NotEmpty(t, *copyResult.VersionId)
	assert.Equal(t, *copyResult.VersionId, sourceVersionId)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectDestName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.UploadPartCopy(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		SourceKey:  Ptr(objectSrcName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000311", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCompleteMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	request := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: parts,
		},
	}
	result, err := client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectName)
	getObj := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	getObjresult, err := client.GetObject(context.TODO(), getObj)
	assert.Nil(t, err)
	data, _ := io.ReadAll(getObjresult.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(data), body)

	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectDestName)

	initCopyResult, err = client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)

	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)

	copyPart := UploadPart{
		PartNumber: copyRequest.PartNumber,
		ETag:       copyResult.ETag,
	}
	var serr *ServiceError
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: []UploadPart{
				copyPart,
			},
		},
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 400, serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
	assert.Equal(t, "0042-00000216", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.CompleteMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		Callback:    Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAbortMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	result, err := client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))

	_, err = client.AbortMultipartUpload(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	abortRequest = &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListMultipartUploads(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6) + "\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)

	putObj := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(randLowStr(1000)),
	}

	_, err = client.PutObject(context.TODO(), putObj)
	assert.Nil(t, err)
	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v\n"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}

	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		SourceKey:  Ptr(objectName),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)

	listRequest := &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListMultipartUploads(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Empty(t, *listResult.KeyMarker, bucketName)
	assert.Len(t, listResult.Uploads, 2)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.ListMultipartUploads(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListMultipartUploads(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListParts(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6) + "-\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)

	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part1)),
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part2)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		Body:       strings.NewReader(string(part3)),
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)

	listRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err := client.ListParts(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Equal(t, *listResult.Key, objectName)
	assert.Equal(t, *listResult.UploadId, *initResult.UploadId)
	assert.Equal(t, *listResult.StorageClass, "Standard")
	assert.Equal(t, listResult.IsTruncated, false)
	assert.Equal(t, listResult.PartNumberMarker, int32(0))
	assert.Equal(t, listResult.NextPartNumberMarker, int32(3))
	assert.Equal(t, listResult.MaxParts, int32(1000))
	assert.Len(t, listResult.Parts, count)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	_, err = client.ListParts(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListPartsRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err = client.ListParts(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
