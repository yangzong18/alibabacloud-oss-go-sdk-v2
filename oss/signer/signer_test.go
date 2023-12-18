package signer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func TestV4AuthHeader(t *testing.T) {
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
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString := "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\n1ef3e8449984a20145d00b2ecb7176be9ca0d082b0b92ee15d06ad05d561908a"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=6b3bb5e2aa90a5655277ae4669c4f08c9cbe653ab3202a6232fa2e132dad67bd", requst.Header.Get("Authorization"))

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
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\naf999f1f855656881dd3fcae98030c82102be740824c8b05b947dad670c34add"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=5cd090f13029eb16905a26599bc57bb80811a671382df8d7361cf286f4b168da", requst.Header.Get("Authorization"))

	//case 3
	requst, _ = http.NewRequest("GET", "http://examplebucket.oss-cn-hangzhou.aliyuncs.com/?resourceGroup&non-resousce=null", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     requst,
		Credentials: &cred,
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
		SubResource: []string{"resourceGroup"},
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\n60cdbe095d466571c6579ba26e59837effc15871bfe8dfc611d9f1ea9d20e9e0"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=9ef4a31dd6c3a1601bde44d22d13b9c6be9cd4c1a6c07e5123c7479ed2ec5715", requst.Header.Get("Authorization"))

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
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)
	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\nc10ddfb39e83bd3cfb6da5152867c809fba16927820169d9fc4f3270f75d34cd"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=abc306ef452ffd62c33a4ee63b6f4b5f0fca52aa5d4cf28239813a82b597ef3e", requst.Header.Get("Authorization"))

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "value")
	requst.Header.Add("content-type", "UNSIGNED-PAYLOAD")

	query, _ := url.ParseQuery(requst.URL.RawQuery)
	query.Add("param1", "value1")
	query.Add("|param1", "value2")
	query.Add("+param1", "value3")
	query.Add("+param1", "value3")
	query.Add("|param1", "value4")
	query.Add("+param2", "")
	query.Add("|param2", "")
	query.Add("param2", "")
	requst.URL.RawQuery = query.Encode()
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     requst,
		Credentials: &cred,
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\nb6f6f9fcf5c216efe5765f5522ae723f4e205ff25f281e5337f4121237d18db9"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=3a9be5c38cbfdade5cdeefd020df92db6b7869aaa520b62dfc05f6ca114a6ca0", requst.Header.Get("Authorization"))

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "value")
	requst.Header.Add("content-type", "UNSIGNED-PAYLOAD")

	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	signCtx = &SigningContext{
		Bucket:            ptr("examplebucket"),
		Request:           requst,
		Credentials:       &cred,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		AdditionalHeaders: []string{"abc", "ZAbc"},
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\nf8d8b7ffb28b1802d1fcdc0914c5309fbd42da3a995629273ec7c297a3a65c02"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,AdditionalHeaders=abc;zabc,Signature=74078299bc82afded3c0a82b1f1f6a9a3a89bd58daf14ca785658deedb6d90ef", requst.Header.Get("Authorization"))

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "value")
	requst.Header.Add("content-type", "UNSIGNED-PAYLOAD")

	query, _ = url.ParseQuery(requst.URL.RawQuery)
	query.Add("param1", "value1")
	query.Add("|param1", "value2")
	query.Add("+param1", "value3")
	query.Add("+param1", "value3")
	query.Add("|param1", "value4")
	query.Add("+param2", "")
	query.Add("|param2", "")
	query.Add("param2", "")
	requst.URL.RawQuery = query.Encode()
	signTime, _ = http.ParseTime("Wed, 28 Dec 2022 10:27:41 GMT")
	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Bucket:      ptr("examplebucket"),
		Request:     requst,
		Credentials: &cred,
		Product:     ptr("oss"),
		Region:      ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	signToString = "OSS4-HMAC-SHA256\nWed, 28 Dec 2022 10:27:41 GMT\n20231218/cn-hangzhou/oss/aliyun_v4_request\nb000442bc07b372f29a28e6b9dc12842725c1209f4368203706ba9a2cb32f467"
	assert.Equal(t, signToString, signCtx.StringToSign)
	assert.Equal(t, signTime, signCtx.Time)
	assert.Equal(t, "OSS4-HMAC-SHA256 Credential=ak/20231218/cn-hangzhou/oss/aliyun_v4_request,Signature=422b7b34d119551cb59ac4e2c432f960a9352b095f8dd8c2a187efac220686d3", requst.Header.Get("Authorization"))
	assert.Equal(t, "token", requst.Header.Get("x-oss-security-token"))
}

func TestV4InvalidArgument(t *testing.T) {
	signer := &SignerV4{}
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

func TestV4AuthQuery(t *testing.T) {
	var provider credentials.CredentialsProvider
	var cred credentials.Credentials
	var expiration time.Time
	var signer Signer
	var signCtx *SigningContext
	provider = credentials.NewStaticCredentialsProvider("ak", "sk")
	cred, _ = provider.GetCredentials(context.TODO())
	requst, _ := http.NewRequest("GET", "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?versionId=versionId", nil)
	requst.Header = http.Header{}
	expiration = time.Now().Add(1 * time.Hour)
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key"),
		Request:         requst,
		Credentials:     &cred,
		Time:            expiration,
		AuthMethodQuery: true,
		Product:         ptr("oss"),
		Region:          ptr("cn-hangzhou"),
	}
	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	assert.Contains(t, requst.URL.String(), "http://bucket.oss-cn-hangzhou.aliyuncs.com/key?versionId=versionId")
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-date=%v", signCtx.currentTime.UTC().Format("20060102T150405Z")))
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
	assert.Contains(t, requst.URL.String(), "x-oss-signature=")
	credential := fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, requst.URL.String(), "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, requst.URL.String(), "x-oss-signature-version=OSS4-HMAC-SHA256")

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	requst.Header.Add("content-type", "text/plain")

	query, _ := url.ParseQuery(requst.URL.RawQuery)
	query.Add("param1", "value1")
	query.Add("|param1", "value2")
	query.Add("+param1", "value3")
	query.Add("+param1", "value3")
	query.Add("|param1", "value4")
	query.Add("+param2", "")
	query.Add("|param2", "")
	query.Add("param2", "")
	requst.URL.RawQuery = query.Encode()
	expiration = time.Now().Add(1 * time.Hour)
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("1234+-/123/1.txt"),
		Request:         requst,
		Credentials:     &cred,
		Time:            expiration,
		AuthMethodQuery: true,
		Product:         ptr("oss"),
		Region:          ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)
	assert.Contains(t, requst.URL.String(), "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt")
	assert.Contains(t, requst.URL.String(), "x-oss-signature=")
	credential = fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, requst.URL.String(), "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, requst.URL.String(), "param1=value1")
	assert.Contains(t, requst.URL.String(), "%7Cparam2")
	assert.Contains(t, requst.URL.String(), "&param2")

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	requst.Header.Add("content-type", "text/plain")

	signCtx = &SigningContext{
		Bucket:            ptr("bucket"),
		Key:               ptr("key"),
		Request:           requst,
		Credentials:       &cred,
		Time:              expiration,
		AuthMethodQuery:   true,
		Product:           ptr("oss"),
		Region:            ptr("cn-hangzhou"),
		AdditionalHeaders: []string{"abc", "ZAbc"},
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)

	assert.Contains(t, requst.URL.String(), "x-oss-signature=")
	credential = fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-date=%v", signCtx.currentTime.UTC().Format("20060102T150405Z")))
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
	assert.Contains(t, requst.URL.String(), "x-oss-signature=")
	assert.Contains(t, requst.URL.String(), "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, requst.URL.String(), "x-oss-additional-headers=abc%3Bzabc")
	assert.Contains(t, requst.URL.String(), "x-oss-signature-version=OSS4-HMAC-SHA256")
	assert.Contains(t, requst.URL.String(), "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt")

	requst, _ = http.NewRequest("PUT", "http://bucket.oss-cn-hangzhou.aliyuncs.com/1234+-/123/1.txt", nil)
	requst.Header = http.Header{}
	requst.Header.Add("x-oss-date", "Wed, 28 Dec 2022 10:27:41 GMT")
	requst.Header.Add("x-oss-head1", "value")
	requst.Header.Add("abc", "value")
	requst.Header.Add("ZAbc", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("XYZ", "value")
	requst.Header.Add("x-oss-content-sha256", "UNSIGNED-PAYLOAD")
	requst.Header.Add("content-type", "text/plain")

	query, _ = url.ParseQuery(requst.URL.RawQuery)
	query.Add("param1", "value1")
	query.Add("|param1", "value2")
	query.Add("+param1", "value3")
	query.Add("+param1", "value3")
	query.Add("|param1", "value4")
	query.Add("+param2", "")
	query.Add("|param2", "")
	query.Add("param2", "")
	requst.URL.RawQuery = query.Encode()
	provider = credentials.NewStaticCredentialsProvider("ak", "sk", "token")
	cred, _ = provider.GetCredentials(context.TODO())
	signCtx = &SigningContext{
		Bucket:          ptr("bucket"),
		Key:             ptr("key"),
		Request:         requst,
		Credentials:     &cred,
		Time:            expiration,
		AuthMethodQuery: true,
		Product:         ptr("oss"),
		Region:          ptr("cn-hangzhou"),
	}

	signer = &SignerV4{}
	signer.Sign(context.TODO(), signCtx)
	assert.Contains(t, requst.URL.String(), "x-oss-signature=")
	assert.Contains(t, requst.URL.String(), "x-oss-credential="+url.QueryEscape(credential))
	assert.Contains(t, requst.URL.String(), "x-oss-security-token=token")
	assert.Contains(t, requst.URL.String(), "param1=value1")
	assert.Contains(t, requst.URL.String(), "%7Cparam2")
	assert.Contains(t, requst.URL.String(), "&param2")
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-date=%v", signCtx.currentTime.UTC().Format("20060102T150405Z")))
	assert.Contains(t, requst.URL.String(), fmt.Sprintf("x-oss-expires=%v", (1*time.Hour).Seconds()))
	credential = fmt.Sprintf("ak/%v/cn-hangzhou/oss/aliyun_v4_request", expiration.Add(-1*time.Hour).UTC().Format("20060102"))
	assert.Contains(t, requst.URL.String(), "x-oss-signature-version=OSS4-HMAC-SHA256")
}
