//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketTags(t *testing.T) {
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

	putRequest := &PutBucketTagsRequest{
		Bucket: Ptr(bucketName),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
					{
						Ptr("key2"),
						Ptr("value2"),
					},
					{
						Ptr("key3"),
						Ptr("value3"),
					},
				},
			},
		},
	}
	putResult, err := client.PutBucketTags(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketTagsRequest{
		Bucket:  Ptr(bucketName),
		Tagging: Ptr("key1,key3"),
	}
	delResult, err := client.DeleteBucketTags(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketTagsRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketTags(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(getResult.Tagging.TagSet.Tags), 1)
	assert.Equal(t, *getResult.Tagging.TagSet.Tags[0].Key, "key2")
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketTagsRequest{
		Bucket: Ptr(bucketNameNotExist),
		Tagging: &Tagging{
			&TagSet{
				[]Tag{
					{
						Ptr("key1"),
						Ptr("value1"),
					},
					{
						Ptr("key2"),
						Ptr("value2"),
					},
					{
						Ptr("key3"),
						Ptr("value3"),
					},
				},
			},
		},
	}
	putResult, err = client.PutBucketTags(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketTagsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketTags(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketTagsRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Tagging: Ptr("key1,key3"),
	}
	delResult, err = client.DeleteBucketTags(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
