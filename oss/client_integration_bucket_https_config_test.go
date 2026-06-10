//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketHttpsConfig(t *testing.T) {
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

	putRequest := &PutBucketHttpsConfigRequest{
		Bucket: Ptr(bucketName),
		HttpsConfiguration: &HttpsConfiguration{
			TLS: &TLS{
				Enable:      Ptr(true),
				TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
			},
			CipherSuite: &CipherSuite{
				Enable:            Ptr(true),
				StrongCipherSuite: Ptr(false),
				CustomCipherSuites: []string{
					"ECDHE-ECDSA-AES128-SHA256", "ECDHE-RSA-AES128-GCM-SHA256", "ECDHE-ECDSA-AES256-CCM8",
				},
				TLS13CustomCipherSuites: []string{
					"TLS_AES_256_GCM_SHA384", "TLS_AES_128_GCM_SHA256", "TLS_CHACHA20_POLY1305_SHA256",
				},
			},
		},
	}
	putResult, err := client.PutBucketHttpsConfig(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketHttpsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketHttpsConfig(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.HttpsConfiguration.TLS.Enable)
	assert.Equal(t, len(getResult.HttpsConfiguration.TLS.TLSVersions), 2)
	assert.True(t, *getResult.HttpsConfiguration.CipherSuite.Enable)
	assert.Equal(t, len(getResult.HttpsConfiguration.CipherSuite.TLS13CustomCipherSuites), 3)
	assert.Equal(t, len(getResult.HttpsConfiguration.CipherSuite.CustomCipherSuites), 3)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	getRequest = &GetBucketHttpsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetBucketHttpsConfig(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	putRequest = &PutBucketHttpsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
		HttpsConfiguration: &HttpsConfiguration{
			TLS: &TLS{
				Enable:      Ptr(true),
				TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
			},
		},
	}
	putResult, err = client.PutBucketHttpsConfig(context.TODO(), putRequest)
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
