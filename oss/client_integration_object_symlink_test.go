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

func TestPutSymlink(t *testing.T) {
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

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	symlinkName := objectName + "-symlink"
	request := &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err := client.PutSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err = client.PutSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.VersionId)

	_, err = client.PutSymlink(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutSymlinkRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	result, err = client.PutSymlink(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetSymlink(t *testing.T) {
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

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	symlinkName := objectName + "-symlink"
	putSymRequest := &PutSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
		Target: Ptr(objectName),
	}
	_, err = client.PutSymlink(context.TODO(), putSymRequest)
	assert.Nil(t, err)

	request := &GetSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
	}
	result, err := client.GetSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, result.ETag)
	assert.Equal(t, *result.Target, objectName)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetSymlinkRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(symlinkName),
	}
	result, err = client.GetSymlink(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, result.ETag)
	assert.Equal(t, *result.Target, objectName)
	assert.NotEmpty(t, *result.VersionId)

	_, err = client.GetSymlink(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetSymlinkRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(symlinkName),
	}
	result, err = client.GetSymlink(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
