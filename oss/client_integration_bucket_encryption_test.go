//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketEncryption(t *testing.T) {
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

	putRequest := &PutBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			&ApplyServerSideEncryptionByDefault{
				SSEAlgorithm:      Ptr("KMS"),
				KMSDataEncryption: Ptr("SM4"),
			},
		},
	}
	putResult, err := client.PutBucketEncryption(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketEncryption(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm, "KMS")
	assert.Equal(t, *getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.KMSDataEncryption, "SM4")
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketEncryptionRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketEncryption(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			&ApplyServerSideEncryptionByDefault{
				SSEAlgorithm:      Ptr("KMS"),
				KMSDataEncryption: Ptr("SM4"),
			},
		},
	}
	putResult, err = client.PutBucketEncryption(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketEncryption(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	delRequest = &DeleteBucketEncryptionRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketEncryption(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
