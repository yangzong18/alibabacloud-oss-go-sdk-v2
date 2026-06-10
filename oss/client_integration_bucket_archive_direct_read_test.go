//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketArchiveDirectRead(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketName),
		ArchiveDirectReadConfiguration: &ArchiveDirectReadConfiguration{
			Ptr(true),
		},
	}
	putResult, err := client.PutBucketArchiveDirectRead(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketArchiveDirectRead(context.TODO(), getRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.ArchiveDirectReadConfiguration.Enabled)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketArchiveDirectRead(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketArchiveDirectReadRequest{
		Bucket: Ptr(bucketNameNotExist),
		ArchiveDirectReadConfiguration: &ArchiveDirectReadConfiguration{
			Ptr(true),
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketArchiveDirectRead(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
