package oss

import (
	"testing"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockInitiateMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<Bucket>oss-example</Bucket>
		<Key>multipart.data</Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object.txt?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object.txt"),
			CacheControl:              Ptr("no-cache"),
			ContentDisposition:        Ptr("attachment"),
			ContentEncoding:           Ptr("utf-8"),
			Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:           Ptr("false"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:              StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging: Ptr("TagA=B&TagC=D"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<Bucket>oss-example</Bucket>
		<Key>multipart.data</Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object.txt?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&InitiateMultipartUploadRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object.txt"),
			CacheControl:              Ptr("no-cache"),
			ContentDisposition:        Ptr("attachment"),
			ContentEncoding:           Ptr("utf-8"),
			Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:           Ptr("false"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:              StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging:      Ptr("TagA=B&TagC=D"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<EncodingType>url</EncodingType>
		<Bucket>oss-example</Bucket>
		<Key></Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object.txt?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&InitiateMultipartUploadRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object.txt"),
			CacheControl:              Ptr("no-cache"),
			ContentDisposition:        Ptr("attachment"),
			ContentEncoding:           Ptr("utf-8"),
			Expires:                   Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:           Ptr("false"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:              StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging:      Ptr("TagA=B&TagC=D"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.Key, "")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
}

func TestMockInitiateMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
}

func TestMockInitiateMultipartUpload_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadDisableDetectMimeTypeCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureAutoDetectMimeType
			})
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute InitiateMultipartUpload fail")
		},
	},
}

func TestMockInitiateMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F587****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Pp****",
			"x-oss-hash-crc64ecma": "2060813895736234537",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 2")
			assert.Equal(t, "f811b746eb3e256f97cb3a190d528353", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(2),
			Body:       strings.NewReader("upload part 2"),
			ContentMD5: Ptr("f811b746eb3e256f97cb3a190d528353"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F587****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Pp****")
			assert.Equal(t, *o.HashCRC64, "2060813895736234537")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F587****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Pp****",
			"x-oss-hash-crc64ecma": "2060813895736234537",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 2")
			assert.Equal(t, "f811b746eb3e256f97cb3a190d528353", r.Header.Get("Content-MD5"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&UploadPartRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(2),
			Body:         strings.NewReader("upload part 2"),
			ContentMD5:   Ptr("f811b746eb3e256f97cb3a190d528353"),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F587****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Pp****")
			assert.Equal(t, *o.HashCRC64, "2060813895736234537")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&UploadPartRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			Body:         strings.NewReader("upload part 1"),
			ContentMD5:   Ptr("bce8f3d48247c5d555bb5697bf277b35"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
}

func TestMockUploadPart_Success(t *testing.T) {
	for _, c := range testMockUploadPartSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "6571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "6571598172666981661")
		},
	},
}

func TestMockUploadPart_Progress(t *testing.T) {
	for _, c := range testMockUploadPartWithProgressCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
			//fmt.Printf("got transferred:%v, total:%v\n", transferred, total)
		}
		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("upload part 1")), n)

	}
}

var testMockUploadPartDisableCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "8571598172666981661",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "8571598172666981661")
		},
	},
}

func TestMockUploadPart_DisableCRC64(t *testing.T) {
	for _, c := range testMockUploadPartDisableCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		//Disable
		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)

		client = NewClient(cfg)
		assert.NotNil(t, c)
		c.Request.Body = strings.NewReader("upload part 1")
		_, err = client.UploadPart(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 6571598172666981661, server 8571598172666981661")
	}
}

var testMockUploadPartErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Body:       strings.NewReader("upload part 1"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUploadPart_Error(t *testing.T) {
	for _, c := range testMockUploadPartErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":             "6551DBCF4311A7303980****",
			"Date":                         "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-copy-source-version-id": "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object")+"?versionId=CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")

			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfMatch), "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfNoneMatch), "\"D41D8CD98F00B204E9800998ECF9****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfModifiedSince), "Fri, 13 Nov 2023 14:47:53 GMT")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfUnmodifiedSince), "Fri, 13 Nov 2015 14:47:53 GMT")
		},
		&UploadPartCopyRequest{
			Bucket:            Ptr("bucket"),
			Key:               Ptr("object"),
			UploadId:          Ptr("0004B9895DBBB6EC9"),
			SourceKey:         Ptr("oss-src-object"),
			SourceBucket:      Ptr("oss-src-bucket"),
			PartNumber:        int32(2),
			IfMatch:           Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
			IfNoneMatch:       Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
			IfModifiedSince:   Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
			IfUnmodifiedSince: Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
			SourceVersionId:   Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
			assert.Equal(t, *o.VersionId, "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":             "6551DBCF4311A7303980****",
			"Date":                         "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-copy-source-version-id": "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
			PartNumber:   int32(2),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
			assert.Equal(t, *o.VersionId, "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
		},
	},
}

func TestMockUploadPartCopy_Success(t *testing.T) {
	for _, c := range testMockUploadPartCopySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6EC9"),
			PartNumber:   int32(1),
			SourceKey:    Ptr("oss-src-object"),
			SourceBucket: Ptr("oss-src-bucket"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute UploadPartCopy fail")
		},
	},
}

func TestMockUploadPartCopy_Error(t *testing.T) {
	for _, c := range testMockUploadPartCopyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <EncodingType>url</EncodingType>
  <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data</Location>
  <Bucket>oss-example</Bucket>
  <Key>demo%2Fmultipart.data</Key>
  <ETag>"097DE458AD02B5F89F9D0530231876****"</ETag>
</CompleteMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>&#34;8EFDA8BE206636A695359836FE0A****&#34;</ETag></Part><Part><PartNumber>2</PartNumber><ETag>&#34;8C315065167132444177411FDA14****&#34;</ETag></Part><Part><PartNumber>3</PartNumber><ETag>&#34;3349DC700140D7F86A0784842780****&#34;</ETag></Part></CompleteMultipartUpload>`)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Parts: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"097DE458AD02B5F89F9D0530231876****\"")
			assert.Equal(t, *o.Location, "http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data")
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "demo/multipart.data")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":     "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****",
			"Content-Type":         "application/json",
			"x-oss-hash-crc64ecma": "1206617243528768****",
		},
		[]byte(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, "false", r.Header.Get(HeaderOssForbidOverWrite))
			assert.Equal(t, "yes", r.Header.Get("x-oss-complete-all"))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`)), r.Header.Get(HeaderOssCallback))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`)), r.Header.Get(HeaderOssCallbackVar))
		},
		&CompleteMultipartUploadRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			UploadId:        Ptr("0004B9895DBBB6EC9"),
			ForbidOverwrite: Ptr("false"),
			CompleteAll:     Ptr("yes"),
			Callback:        Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))),
			CallbackVar:     Ptr(base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get(HTTPHeaderContentType), "application/json")
			jsonData, _ := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
			assert.Equal(t, *o.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":     "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****",
			"Content-Type":         "application/json",
			"x-oss-hash-crc64ecma": "1206617243528768****",
		},
		[]byte(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, "false", r.Header.Get(HeaderOssForbidOverWrite))
			assert.Equal(t, "yes", r.Header.Get("x-oss-complete-all"))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`)), r.Header.Get(HeaderOssCallback))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`)), r.Header.Get(HeaderOssCallbackVar))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&CompleteMultipartUploadRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			UploadId:        Ptr("0004B9895DBBB6EC9"),
			ForbidOverwrite: Ptr("false"),
			CompleteAll:     Ptr("yes"),
			Callback:        Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))),
			CallbackVar:     Ptr(base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get(HTTPHeaderContentType), "application/json")
			jsonData, _ := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
			assert.Equal(t, *o.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <EncodingType>url</EncodingType>
  <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data</Location>
  <Bucket>oss-example</Bucket>
  <Key></Key>
  <ETag>"097DE458AD02B5F89F9D0530231876****"</ETag>
</CompleteMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>&#34;8EFDA8BE206636A695359836FE0A****&#34;</ETag></Part><Part><PartNumber>2</PartNumber><ETag>&#34;8C315065167132444177411FDA14****&#34;</ETag></Part><Part><PartNumber>3</PartNumber><ETag>&#34;3349DC700140D7F86A0784842780****&#34;</ETag></Part></CompleteMultipartUpload>`)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Parts: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"097DE458AD02B5F89F9D0530231876****\"")
			assert.Equal(t, *o.Location, "http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data")
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "")
		},
	},
}

func TestMockCompleteMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D94CCD11E55313348****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>655D94CCD11E55313348****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0042-00000205</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000205</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "655D94CCD11E55313348****", serr.RequestID)
			assert.Equal(t, "0042-00000205", serr.EC)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D9598CA31DC313626****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>Should not speficy both complete all header and http body.</Message>
  <RequestId>655D9598CA31DC313626****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ArgumentName>x-oss-complete-all</ArgumentName>
  <ArgumentValue>yes</ArgumentValue>
  <EC>0042-00000216</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000216</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object"),
			UploadId:    Ptr("0004B9895DBBB6EC9"),
			CompleteAll: Ptr("yes"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Parts: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
			assert.Equal(t, "655D9598CA31DC313626****", serr.RequestID)
			assert.Equal(t, "0042-00000216", serr.EC)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CompleteMultipartUpload fail")
		},
	},
}

func TestMockCompleteMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AbortMultipartUploadRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B9895DBBB6E"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockAbortMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockAbortMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
    <Bucket>oss-example</Bucket>
    <KeyMarker></KeyMarker>
    <UploadIdMarker></UploadIdMarker>
    <NextKeyMarker>oss.avi</NextKeyMarker>
    <NextUploadIdMarker>0004B99B8E707874FC2D692FA5D77D3F</NextUploadIdMarker>
    <Delimiter></Delimiter>
    <Prefix></Prefix>
    <MaxUploads>1000</MaxUploads>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF518A1FE585B0C9360DC4C8</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>oss.avi</Key>
        <UploadId>0004B99B8E707874FC2D692FA5D7****</UploadId>
        <Initiated>2012-02-23T06:14:27.000Z</Initiated>
    </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "0004B99B8E707874FC2D692FA5D77D3F")
			assert.Equal(t, *o.Delimiter, "")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 3)
			assert.Equal(t, *o.Uploads[0].Key, "multipart.data")
			assert.Equal(t, *o.Uploads[0].UploadId, "0004B999EF518A1FE585B0C9360DC4C8")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2012, time.February, 23, 4, 18, 23, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 10)
			assert.Equal(t, *o.Uploads[0].Key, "demo/gp-\f\n\v")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
			RequestPayer:   Ptr("requester"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 10)
			assert.Equal(t, *o.Uploads[0].Key, "demo/gp-\f\n\v")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key></Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
			RequestPayer:   Ptr("requester"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 1)
			assert.Equal(t, *o.Uploads[0].Key, "")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
}

func TestMockListMultipartUploads_Success(t *testing.T) {
	for _, c := range testMockListMultipartUploadsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockListMultipartUploads_Error(t *testing.T) {
	for _, c := range testMockListMultipartUploadsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>bucket</Bucket>
    <Key>object</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Key, "object")
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(5))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 3)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(6291456))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
  <EncodingType>url</EncodingType>
  <Bucket>bucket</Bucket>
  <Key>demo%2Fgp-%0C%0A%0B</Key>
  <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
  <StorageClass>Standard</StorageClass>
  <PartNumberMarker>0</PartNumberMarker>
  <NextPartNumberMarker>1</NextPartNumberMarker>
  <MaxParts>1000</MaxParts>
  <IsTruncated>false</IsTruncated>
  <Part>
    <PartNumber>1</PartNumber>
    <LastModified>2023-11-22T05:42:34.000Z</LastModified>
    <ETag>"CF3F46D505093571E916FCDD4967****"</ETag>
    <HashCrc64ecma>12066172435287683848</HashCrc64ecma>
    <Size>96316</Size>
  </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/demo/gp-%0C%0A%0B?encoding-type=url&uploadId=D4E111D4EA834F3ABCE4877B2779%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("demo/gp-\f\n\v"),
			UploadId: Ptr("D4E111D4EA834F3ABCE4877B2779****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			key, _ := url.QueryUnescape("demo%2Fgp-%0C%0A%0B")
			assert.Equal(t, *o.Key, key)
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(1))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 1)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"CF3F46D505093571E916FCDD4967****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2023, time.November, 22, 5, 42, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(96316))
			assert.Equal(t, *o.Parts[0].HashCRC64, "12066172435287683848")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>bucket</Bucket>
    <Key>object</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListPartsRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			UploadId:     Ptr("0004B999EF5A239BB9138C6227D6****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Key, "object")
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(5))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 3)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(6291456))
		},
	},
}

func TestMockListParts_Success(t *testing.T) {
	for _, c := range testMockListPartsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListParts fail")
		},
	},
}

func TestMockListParts_Error(t *testing.T) {
	for _, c := range testMockListPartsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


