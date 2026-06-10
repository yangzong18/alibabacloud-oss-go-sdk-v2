//go:build integration

package oss

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPutObject(t *testing.T) {
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
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err := client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)

	request = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)

	var serr *ServiceError
	request = &PutObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	result, err = client.PutObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	request = &PutObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	//Body is bigger than Content-Length
	request = &PutObjectRequest{
		Bucket:        Ptr(bucketName),
		Key:           Ptr(objectName),
		ContentLength: Ptr(int64(len(content) - 10)),
		Body:          strings.NewReader(content),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), " transport connection broken")

}

func TestGetObject(t *testing.T) {
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
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)
	getRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObject(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.ContentLength, int64(len(content)))

	getRequest = &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.GetObject(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	_, err = client.GetObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	getRequest = &GetObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.GetObject(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCopyObject(t *testing.T) {
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
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	objectCopyName := objectNamePrefix + randLowStr(6) + "copy"
	copyRequest := &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		SourceKey: Ptr(objectCopyName),
	}
	result, err := client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)

	copyRequest = &CopyObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectCopyName),
		SourceKey: Ptr(objectName),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)

	copyRequest = &CopyObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)

	_, err = client.CopyObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	copyRequest = &CopyObjectRequest{
		Bucket:    Ptr(bucketNameNotExist),
		Key:       Ptr(objectCopyName),
		SourceKey: Ptr(objectName),
	}
	_, err = client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

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
		Key:    Ptr(objectName),
	}
	metaResult, err := client.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)
	sourceVersionId := *metaResult.VersionId

	copyRequest = &CopyObjectRequest{
		Bucket:          Ptr(bucketName),
		Key:             Ptr(objectCopyName),
		SourceKey:       Ptr(objectName),
		SourceVersionId: Ptr(sourceVersionId),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, result.VersionId)
	assert.Equal(t, *result.SourceVersionId, sourceVersionId)

	bucketCopyName := bucketNamePrefix + randLowStr(6) + "copy"
	putRequest = &PutBucketRequest{
		Bucket: Ptr(bucketCopyName),
	}
	client = getDefaultClient()
	_, err = client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	copyRequest = &CopyObjectRequest{
		Bucket:       Ptr(bucketCopyName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		SourceBucket: Ptr(bucketName),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)
}

func TestAppendObject(t *testing.T) {
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
	var result *AppendObjectResult
	content := randLowStr(100)
	request := &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition := result.NextPosition
	request = &AppendObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		Position:     Ptr(nextPosition),
		TrafficLimit: int64(100 * 1024 * 8),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*3))
	assert.NotEmpty(t, result.HashCRC64)

	objectName2 := objectName + "-kms-sm4"
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName2),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(int64(0)),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket:                   Ptr(bucketName),
		Key:                      Ptr(objectName2),
		Body:                     strings.NewReader(content),
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
		TrafficLimit:             int64(100 * 1024 * 8),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.ServerSideEncryptionKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	_, err = client.AppendObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "PositionNotEqualToLength", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketName + "-not-exist"
	request = &AppendObjectRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteObject(t *testing.T) {
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
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	exist, err := client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.True(t, exist)

	delRequest := &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectName)
	assert.Nil(t, err)
	assert.False(t, exist)

	objectNameNotExist := objectNamePrefix + randLowStr(6) + "-not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	exist, err = client.IsObjectExist(context.TODO(), bucketName, objectNameNotExist)
	assert.Nil(t, err)
	assert.False(t, exist)

	delRequest = &DeleteObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr("null"),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	_, err = client.DeleteObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectNamePrefix),
	}
	_, err = client.DeleteObject(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteMultipleObjects(t *testing.T) {
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
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest := &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketName),
		Objects: []DeleteObject{{Key: Ptr(objectName)}},
	}
	result, err := client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectName)

	str := "\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"
	objectNameSpecial := objectNamePrefix + randLowStr(6) + str
	content = randLowStr(10)
	request = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr(bucketName),
		Objects:      []DeleteObject{{Key: Ptr(objectNameSpecial)}},
		EncodingType: Ptr("url"),
	}
	result, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectNameSpecial)

	delRequest = &DeleteMultipleObjectsRequest{
		Bucket: Ptr(bucketName),
		Delete: &Delete{Objects: nil},
	}
	for i := 0; i < 10; i++ {
		request.Key = Ptr(objectName + strconv.Itoa(i))
		_, err = client.PutObject(context.TODO(), request)
		assert.Nil(t, err)
		delRequest.Delete.Objects = append(delRequest.Delete.Objects, ObjectIdentifier{Key: Ptr(objectName + strconv.Itoa(i))})
	}
	result, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 10)

	_, err = client.DeleteMultipleObjects(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Objects: []DeleteObject{{Key: Ptr(objectNameSpecial)}},
	}
	_, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestHeadObject(t *testing.T) {
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
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.HeadObject(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ContentMD5)
	assert.NotEmpty(t, *result.ObjectType)
	assert.NotEmpty(t, *result.StorageClass)
	assert.NotEmpty(t, *result.ETag)
	_, err = client.HeadObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.HeadObject(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectMeta(t *testing.T) {
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
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content),
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectMeta(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)

	_, err = client.GetObjectMeta(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &GetObjectMetaRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectMeta(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
