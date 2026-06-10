//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketStyle(t *testing.T) {
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

	putRequest := &PutStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
		Style: &StyleContent{
			Content: Ptr("image/resize,p_50"),
		},
	}
	putResult, err := client.PutStyle(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
	}
	getResult, err := client.GetStyle(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listRequest := &ListStyleRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListStyle(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.StyleList.Styles), 1)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteStyleRequest{
		Bucket:    Ptr(bucketName),
		StyleName: Ptr("demo"),
	}
	delResult, err := client.DeleteStyle(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
		Style: &StyleContent{
			Content: Ptr("image/resize,p_50"),
		},
	}
	putResult, err = client.PutStyle(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
	}
	getResult, err = client.GetStyle(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listRequest = &ListStyleRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListStyle(context.TODO(), listRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteStyleRequest{
		Bucket:    Ptr(bucketNameNotExist),
		StyleName: Ptr("demo"),
	}
	delResult, err = client.DeleteStyle(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
