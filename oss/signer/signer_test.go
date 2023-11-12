package signer

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T {
	return &v
}

func TestSigningContext(t *testing.T) {
	r := SigningContext{}
	assert.Empty(t, r.Product)
	assert.Empty(t, r.Region)
	assert.Empty(t, r.Bucket)
	assert.Empty(t, r.Key)
	assert.Nil(t, r.Request)
	assert.Empty(t, r.SubResource)

	assert.Empty(t, r.Credentials)
	assert.Empty(t, r.StringToSign)
	assert.Empty(t, r.SignedHeaders)
	assert.Empty(t, r.Time)
}

func TestNopSigner(t *testing.T) {
	r := NopSigner{}
	assert.Nil(t, r.Sign(context.TODO(), nil))
}

func TestV1AuthHeader(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())

	//case 1
	requst, _ := http.NewRequest("PUT", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com", nil)
	requst.Header = http.Header{}
	requst.Header.Add("Content-MD5", "eB5eJF1ptWaXm4bijSPyxw==")
	requst.Header.Add("Content-Type", "text/html")
	requst.Header.Add("x-oss-meta-author", "alice")
	requst.Header.Add("x-oss-meta-magic", "abracadabra")
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Key:         ptr("nelson"),
		Request:     requst,
		Credentials: &cred,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString := "PUT\neB5eJF1ptWaXm4bijSPyxw==\ntext/html\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\nx-oss-meta-author:alice\nx-oss-meta-magic:abracadabra\n/examplebucket/nelson"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:kSHKmLxlyEAKtZPkJhG9bZb5k7M=", requst.Header.Get("Authorization"))

	//case 2
	requst, _ = http.NewRequest("PUT", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?acl", nil)
	requst.Header = http.Header{}
	requst.Header.Add("Content-MD5", "eB5eJF1ptWaXm4bijSPyxw==")
	requst.Header.Add("Content-Type", "text/html")
	requst.Header.Add("x-oss-meta-author", "alice")
	requst.Header.Add("x-oss-meta-magic", "abracadabra")
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Key:         ptr("nelson"),
		Request:     requst,
		Credentials: &cred,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "PUT\neB5eJF1ptWaXm4bijSPyxw==\ntext/html\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\nx-oss-meta-author:alice\nx-oss-meta-magic:abracadabra\n/examplebucket/nelson?acl"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:/afkugFbmWDQ967j1vr6zygBLQk=", requst.Header.Get("Authorization"))

	//case 3
	requst, _ = http.NewRequest("GET", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?resourceGroup&non-resousce=null", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     requst,
		Credentials: &cred,
		SubResource: []string{"resourceGroup"},
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "GET\n\n\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\n/examplebucket/?resourceGroup"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS ak:vkQmfuUDyi1uDi3bKt67oemssIs=", requst.Header.Get("Authorization"))

	//case 4
	requst, _ = http.NewRequest("GET", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?resourceGroup&acl", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     requst,
		Credentials: &cred,
		SubResource: []string{"resourceGroup"},
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "GET\n\n\nWed, 28 Dec 2022 10:27:41 GMT\nx-oss-date:Wed, 28 Dec 2022 10:27:41 GMT\n/examplebucket/?acl&resourceGroup"
}

func TestV1InvalidArgument(t *testing.T) {
	signer := &SignerV1{}
	signCtx := &SigningContext{}
	err := signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider := credentials.NewStaticCredentialsProvider("", "sk")
	cred, _ := provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Credentials is null or empty")

	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Credentials: &cred,
	}
	err = signer.Sign(context.TODO(), signCtx)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext.Request is null")

	err = signer.Sign(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "SigningContext is null")
}

func TestV1AuthQuery(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var signTime time.Time
	var signer Signer
	var signCtx *SigningContext

	//case 1
	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	requst, _ := http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?versionId=versionId", nil)
	requst.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:43:40 GMT")
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key"),
		Request:         requst,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl := "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?Expires=1699807420&OSSAccessKeyId=ak&Signature=dcLTea%2BYh9ApirQ8o8dOPqtvJXQ%3D&versionId=versionId"
	assert.Equal(t, signUrl, requst.URL.String())

	//case 2
	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())
	requst, _ = http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?versionId=versionId", nil)
	requst.Header = http.Header{}
	signTime, _ = http.ParseTime("Sun, 12 Nov 2023 16:56:44 GMT")
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key+123"),
		Request:         requst,
		Credentials:     &cred,
		Time:            signTime,
		AuthMethodQuery: true,
	}

	signer = &SignerV1{}
	signer.Sign(context.TODO(), signCtx)
	signUrl = "http://bucket.oss-cn-hangzhou.aliyuncs.com/key%2B123?Expires=1699808204&OSSAccessKeyId=ak&Signature=jzKYRrM5y6Br0dRFPaTGOsbrDhY%3D&security-token=token&versionId=versionId"
	assert.Equal(t, signUrl, requst.URL.String())
}
