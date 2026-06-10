//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketLifecycle(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	createRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), createRequest)
	assert.Nil(t, err)

	putRequest := &PutBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
		LifecycleConfiguration: &LifecycleConfiguration{
			Rules: []LifecycleRule{
				{
					Status: Ptr("Enabled"),
					ID:     Ptr("rule"),
					Prefix: Ptr("log/"),
					Transitions: []LifecycleRuleTransition{
						{
							Days:         Ptr(int32(30)),
							StorageClass: StorageClassIA,
						},
					},
				},
			},
		},
	}

	putResult, err := client.PutBucketLifecycle(context.TODO(), putRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))

	getRequest := &GetBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
	}

	result, err := client.GetBucketLifecycle(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.LifecycleConfiguration)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketLifecycleRequest{
		Bucket: Ptr(bucketName),
	}

	delResult, err := client.DeleteBucketLifecycle(context.TODO(), delRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 204, delResult.StatusCode)
	assert.Equal(t, "204 No Content", delResult.Status)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest.Bucket = Ptr(bucketNameNotExist)
	_, err = client.PutBucketLifecycle(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))

	getRequest.Bucket = Ptr(bucketNameNotExist)
	serr = &ServiceError{}
	_, err = client.GetBucketLifecycle(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest.Bucket = Ptr(bucketNameNotExist)
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketLifecycle(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
