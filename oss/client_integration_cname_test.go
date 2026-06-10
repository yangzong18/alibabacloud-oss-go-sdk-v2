//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCname(t *testing.T) {
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

	createResult, err := client.CreateCnameToken(context.TODO(), &CreateCnameTokenRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotNil(t, createResult.CnameToken)
	time.Sleep(1 * time.Second)

	getResult, err := client.GetCnameToken(context.TODO(), &GetCnameTokenRequest{
		Bucket: Ptr(bucketName),
		Cname:  Ptr("example.com"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.NotNil(t, getResult.CnameToken)
	time.Sleep(1 * time.Second)

	listResult, err := client.ListCname(context.TODO(), &ListCnameRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.Cnames), 0)
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteCname(context.TODO(), &DeleteCnameRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.PutCname(context.TODO(), &PutCnameRequest{
		Bucket: Ptr(bucketName),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "NeedVerifyDomainOwnership", serr.Code)
	assert.Equal(t, "0018-00000115", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	createResult, err = client.CreateCnameToken(context.TODO(), &CreateCnameTokenRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.GetCnameToken(context.TODO(), &GetCnameTokenRequest{
		Bucket: Ptr(bucketNameNotExist),
		Cname:  Ptr("example.com"),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.ListCname(context.TODO(), &ListCnameRequest{
		Bucket: Ptr(bucketNameNotExist),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DeleteCname(context.TODO(), &DeleteCnameRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketCnameConfiguration: &BucketCnameConfiguration{
			Domain: Ptr("example.com"),
		},
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
