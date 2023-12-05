package oss

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestPresignPresignOptions(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.Equal(t, expiration, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")

	expires := 50 * time.Minute
	expiration = time.Now().Add(expires)
	result, err = client.Presign(context.TODO(), request, PresignExpires(expires))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.True(t, result.Expiration.Unix()-expiration.Unix() < 2)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", result.Expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
}

func TestPresignWithToken(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
	assert.Contains(t, result.URL, "security-token=token")
}

func TestPresignWithHeader(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	request := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"Content-Type": "application/octet-stream",
			},
		},
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Len(t, result.SignedHeaders, 1)
	assert.Equal(t, "application/octet-stream", result.SignedHeaders["Content-Type"])
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
}

func TestPresignWithQuery(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	reqeust := &GetObjectRequest{
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"x-oss-process": "abc",
			},
		},
	}

	expiration := time.Now().Add(1 * time.Hour)
	result, err := client.Presign(context.TODO(), reqeust, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, fmt.Sprintf("Expires=%v", expiration.Unix()))
	assert.Contains(t, result.URL, "Signature=")
	assert.Contains(t, result.URL, "x-oss-process=abc")
}

func TestPresignOperationInput(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	request := &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: Ptr("bucket"),
		Key:    Ptr("key"),
		Parameters: map[string]string{
			"versionId": "versionId",
		},
	}

	expiration, _ := http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")
	result, err := client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, "Expires=1699807420")
	assert.Contains(t, result.URL, "Signature=dcLTea%2BYh9ApirQ8o8dOPqtvJXQ%3D")

	//token
	cfg = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk", "token")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client = NewClient(cfg)

	request = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: Ptr("bucket"),
		Key:    Ptr("key+123"),
		Parameters: map[string]string{
			"versionId": "versionId",
		},
	}

	expiration, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	result, err = client.Presign(context.TODO(), request, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Empty(t, result.SignedHeaders)
	assert.Contains(t, result.URL, "bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?")
	assert.Contains(t, result.URL, "OSSAccessKeyId=ak")
	assert.Contains(t, result.URL, "Expires=1699808204")
	assert.Contains(t, result.URL, "Signature=jzKYRrM5y6Br0dRFPaTGOsbrDhY%3D")
	assert.Contains(t, result.URL, "security-token=token")
	assert.Contains(t, result.URL, "versionId=versionId")
}

func TestPresignWithError(t *testing.T) {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("ak", "sk")).
		WithRegion("cn-hangzhou").
		WithEndpoint("oss-cn-hangzhou.aliyuncs.com")

	client := NewClient(cfg)

	// unsupport request
	request := &ListObjectsRequest{
		Bucket: Ptr("bucket"),
	}
	_, err := client.Presign(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "request *oss.ListObjectsRequest")

	// request is nil
	_, err = client.Presign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "null field, request")
}
