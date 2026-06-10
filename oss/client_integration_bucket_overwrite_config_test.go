//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketOverwriteConfig(t *testing.T) {
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

	putRequest := &PutBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketName),
		OverwriteConfiguration: &OverwriteConfiguration{
			Rules: []OverwriteRule{
				{
					ID:     Ptr("1"),
					Action: Ptr("forbid"),
				},
				{
					ID:     Ptr("2"),
					Action: Ptr("forbid"),
					Prefix: Ptr("pre"),
					Suffix: Ptr(".txt"),
					Principals: &OverwritePrincipals{
						[]string{"1234567890"},
					},
				},
			},
		},
	}
	putResult, err := client.PutBucketOverwriteConfig(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketOverwriteConfig(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketOverwriteConfig(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketOverwriteConfig(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
		OverwriteConfiguration: &OverwriteConfiguration{
			Rules: []OverwriteRule{
				{
					ID:     Ptr("1"),
					Action: Ptr("forbid"),
				},
			},
		},
	}
	serr = &ServiceError{}
	putResult, err = client.PutBucketOverwriteConfig(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "MalformedXML", serr.Code)
	assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
	assert.Equal(t, "0015-00000231", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketOverwriteConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delResult, err = client.DeleteBucketOverwriteConfig(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
