//go:build integration

package oss

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRestoreObjectLegacy(t *testing.T) {
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
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	restoreRequest := &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 202)
	assert.Equal(t, result.Status, "202 Accepted")
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))

	_, err = client.RestoreObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
	assert.Equal(t, "The restore operation is in progress.", serr.Message)
	assert.NotEmpty(t, serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestRestoreObject(t *testing.T) {
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
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(content),
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	restoreRequest := &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RestoreRequest: &RestoreRequest{
			Days: 1,
			JobParameters: &JobParameters{
				Tier: Ptr("Standard"),
			},
		},
	}
	result, err := client.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 202)
	assert.Equal(t, result.Status, "202 Accepted")
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))

	_, err = client.RestoreObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
	assert.Equal(t, "The restore operation is in progress.", serr.Message)
	assert.NotEmpty(t, serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCleanRestoredObject(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	_, err = client.RestoreObject(context.TODO(), &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	_, err = client.CleanRestoredObject(context.TODO(), &CleanRestoredObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "ArchiveRestoreNotFinished", serr.Code)
	assert.Equal(t, "The archive file's restore is not finished.", serr.Message)
	assert.Equal(t, "0016-00000719", serr.EC)
}

func TestSealAppendObject(t *testing.T) {
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
	var result *SealAppendObjectResult
	content := randLowStr(100)
	appRequest := &AppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Body:     strings.NewReader(content),
		Position: Ptr(int64(0)),
	}
	appResult, err := client.AppendObject(context.TODO(), appRequest)
	assert.Nil(t, err)
	assert.Equal(t, appResult.NextPosition, int64(len(content)))
	assert.NotEmpty(t, appResult.HashCRC64)

	request := &SealAppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Position: Ptr(appResult.NextPosition),
	}
	result, err = client.SealAppendObject(context.TODO(), request)

	var serr *ServiceError
	// TOD remove later
	if err != nil {
		errors.As(err, &serr)
		assert.Equal(t, int(400), serr.StatusCode)
		assert.Equal(t, "OperationNotSupported", serr.Code)
		assert.Equal(t, "SealAppendable is not supported.", serr.Message)
		assert.Equal(t, "0016-00000513", serr.EC)
		return
	}

	assert.Nil(t, err)
	assert.NotEmpty(t, *result.SealedTime)

	request = &SealAppendObjectRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		Position: Ptr(int64(0)),
	}
	_, err = client.SealAppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "PositionNotEqualToLength", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketName + "-not-exist"
	request = &SealAppendObjectRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		Position: Ptr(int64(0)),
	}
	_, err = client.SealAppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
