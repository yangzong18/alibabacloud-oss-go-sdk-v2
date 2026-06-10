//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestPublicAccessBlock(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	client := getDefaultClient()

	putResult, err := client.PutPublicAccessBlock(context.TODO(), &PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetPublicAccessBlock(context.TODO(), &GetPublicAccessBlockRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeletePublicAccessBlock(context.TODO(), &DeletePublicAccessBlockRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	putResult, err = noPermClient.PutPublicAccessBlock(context.TODO(), &PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = noPermClient.GetPublicAccessBlock(context.TODO(), &GetPublicAccessBlockRequest{})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = noPermClient.DeletePublicAccessBlock(context.TODO(), &DeletePublicAccessBlockRequest{})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestBucketPublicAccessBlock(t *testing.T) {
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
	putResult, err := client.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetBucketPublicAccessBlock(context.TODO(), &GetBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteBucketPublicAccessBlock(context.TODO(), &DeleteBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	putResult, err = noPermClient.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = noPermClient.GetBucketPublicAccessBlock(context.TODO(), &GetBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = noPermClient.DeleteBucketPublicAccessBlock(context.TODO(), &DeleteBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
