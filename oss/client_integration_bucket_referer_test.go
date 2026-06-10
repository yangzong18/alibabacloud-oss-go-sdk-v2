//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketReferer(t *testing.T) {
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

	putRequest := &PutBucketRefererRequest{
		Bucket: Ptr(bucketName),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer: Ptr(false),
			RefererList: &RefererList{
				Referers: []string{""},
			},
		},
	}
	putResult, err := client.PutBucketReferer(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketRefererRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketReferer(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.False(t, *getResult.RefererConfiguration.AllowEmptyReferer)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketRefererRequest{
		Bucket: Ptr(bucketNameNotExist),
		RefererConfiguration: &RefererConfiguration{
			AllowEmptyReferer: Ptr(false),
		},
	}
	putResult, err = client.PutBucketReferer(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketRefererRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketReferer(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
