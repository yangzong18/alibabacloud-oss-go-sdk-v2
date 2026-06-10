package oss

import (
	"testing"
	"context"
	"errors"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutSymlinkSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *PutSymlinkResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("src-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
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
			"Content-Type":     "application/xml",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)

		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("src-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-symlink-target"), "target-object")
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "true")
			assert.Equal(t, r.Header.Get("x-oss-object-acl"), string(ObjectACLPrivate))
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "demo")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
		},
		&PutSymlinkRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			Target:          Ptr("target-object"),
			ForbidOverwrite: Ptr("true"),
			Acl:             ObjectACLPrivate,
			StorageClass:    StorageClassStandard,
			Metadata: map[string]string{
				"name":  "demo",
				"email": "demo@aliyun.com",
			},
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutSymlinkRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Target:       Ptr("src-object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutSymlink_Success(t *testing.T) {
	for _, c := range testMockPutSymlinkSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutSymlinkErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *PutSymlinkResult, err error)
}{
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
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("target-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
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
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&PutSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Target: Ptr("target-object"),
		},
		func(t *testing.T, o *PutSymlinkResult, err error) {
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

func TestMockPutSymlink_Error(t *testing.T) {
	for _, c := range testMockPutSymlinkErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSymlinkSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *GetSymlinkResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":         "application/xml",
			"x-oss-version-id":     "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
			"x-oss-meta-name":      "demo",
			"x-oss-meta-email":     "demo@aliyun.com",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink&versionId=CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh%2A%2A%2A%2A", strUrl)
		},
		&GetSymlinkRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****")
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-symlink-target": "example.jpg",
			"ETag":                 "A797938C31D59EDD08D86188F6D5****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetSymlinkRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Target, "example.jpg")
			assert.Equal(t, *o.ETag, "A797938C31D59EDD08D86188F6D5****")
		},
	},
}

func TestMockGetSymlink_Success(t *testing.T) {
	for _, c := range testMockGetSymlinkSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSymlinkErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSymlinkRequest
	CheckOutputFn  func(t *testing.T, o *GetSymlinkResult, err error)
}{
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
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
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
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?symlink", strUrl)
		},
		&GetSymlinkRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetSymlinkResult, err error) {
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

func TestMockGetSymlink_Error(t *testing.T) {
	for _, c := range testMockGetSymlinkErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSymlink(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


