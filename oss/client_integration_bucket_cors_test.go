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

func TestBucketCors(t *testing.T) {
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

	putRequest := &PutBucketCorsRequest{
		Bucket: Ptr(bucketName),
		CORSConfiguration: &CORSConfiguration{
			CORSRules: []CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"PUT", "GET"},
				},
			},
		},
	}
	putResult, err := client.PutBucketCors(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	key := objectNamePrefix + randStr(6)
	_, err = client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(key),
		Body:   strings.NewReader("hi oss"),
	})
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	optionRequest := &OptionObjectRequest{
		Bucket:                     Ptr(bucketName),
		Key:                        Ptr(key),
		Origin:                     Ptr("http://www.example.com"),
		AccessControlRequestMethod: Ptr("PUT"),
	}
	optionResult, err := client.OptionObject(context.TODO(), optionRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, optionResult.StatusCode)
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketCorsRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketCors(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Equal(t, 1, len(getResult.CORSConfiguration.CORSRules))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketCorsRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketCors(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
		CORSConfiguration: &CORSConfiguration{
			CORSRules: []CORSRule{
				{
					AllowedOrigins: []string{"*"},
					AllowedMethods: []string{"PUT", "GET"},
				},
			},
		},
	}
	putResult, err = client.PutBucketCors(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketCors(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketCorsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketCors(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	optionRequest = &OptionObjectRequest{
		Bucket:                     Ptr(bucketNameNotExist),
		Key:                        Ptr(key),
		Origin:                     Ptr("http://www.example.com"),
		AccessControlRequestMethod: Ptr("PUT"),
	}
	optionResult, err = client.OptionObject(context.TODO(), optionRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
