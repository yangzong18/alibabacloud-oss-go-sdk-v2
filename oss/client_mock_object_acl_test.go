package oss

import (
	"testing"
	"context"
	"errors"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPublicRead), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPublicRead,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0%2A%2A%2A%2A", strUrl)
			assert.Equal(t, string(ObjectACLPublicReadWrite), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			Acl:       ObjectACLPublicReadWrite,
			VersionId: Ptr("CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****"),
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get(HTTPHeaderContentType))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get(HeaderOssRequestID))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get(HTTPHeaderDate))

			assert.Equal(t, "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****", *o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0%2A%2A%2A%2A", strUrl)
			assert.Equal(t, string(ObjectACLPublicReadWrite), r.Header.Get(HeaderOssObjectACL))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutObjectAclRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Acl:          ObjectACLPublicReadWrite,
			VersionId:    Ptr("CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get(HTTPHeaderContentType))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get(HeaderOssRequestID))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get(HTTPHeaderDate))

			assert.Equal(t, "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****", *o.VersionId)
		},
	},
}

func TestMockPutObjectAcl_Success(t *testing.T) {
	for _, c := range testMockPutObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
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
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl", strUrl)
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
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
}

func TestMockPutObjectAcl_Error(t *testing.T) {
	for _, c := range testMockPutObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read", *o.ACL)
			assert.Equal(t, "0022012****", *o.Owner.ID)
			assert.Equal(t, "user_example", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read-write</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read-write", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk%2A%2A%2A%2A", strUrl)
		},
		&GetObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
			assert.Equal(t, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****", *o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk%2A%2A%2A%2A", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectAclRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			VersionId:    Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
			assert.Equal(t, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****", *o.VersionId)
		},
	},
}

func TestMockGetObjectAcl_Success(t *testing.T) {
	for _, c := range testMockGetObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", strUrl)
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
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
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
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

func TestMockGetObjectAcl_Error(t *testing.T) {
	for _, c := range testMockGetObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


