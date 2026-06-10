package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketCorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketCorsResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader></CORSRule><CORSRule><AllowedOrigin>http://example.com</AllowedOrigin><AllowedOrigin>http://example.net</AllowedOrigin><AllowedMethod>GET</AllowedMethod><AllowedHeader>Authorization</AllowedHeader><ExposeHeader>x-oss-test</ExposeHeader><ExposeHeader>x-oss-test1</ExposeHeader><MaxAgeSeconds>100</MaxAgeSeconds></CORSRule><ResponseVary>false</ResponseVary></CORSConfiguration>")
		},
		&PutBucketCorsRequest{
			Bucket: Ptr("bucket"),
			CORSConfiguration: &CORSConfiguration{
				CORSRules: []CORSRule{
					{
						AllowedOrigins: []string{"*"},
						AllowedMethods: []string{"PUT", "GET"},
						AllowedHeaders: []string{"Authorization"},
					},
					{
						AllowedOrigins: []string{"http://example.com", "http://example.net"},
						AllowedMethods: []string{"GET"},
						AllowedHeaders: []string{"Authorization"},
						ExposeHeaders:  []string{"x-oss-test", "x-oss-test1"},
						MaxAgeSeconds:  Ptr(int64(100)),
					},
				},
				ResponseVary: Ptr(false),
			},
		},
		func(t *testing.T, o *PutBucketCorsResult, err error) {
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
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod></CORSRule></CORSConfiguration>")
		},
		&PutBucketCorsRequest{
			Bucket: Ptr("bucket"),
			CORSConfiguration: &CORSConfiguration{
				CORSRules: []CORSRule{
					{
						AllowedOrigins: []string{"*"},
						AllowedMethods: []string{"PUT", "GET"},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketCorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketCors_Success(t *testing.T) {
	for _, c := range testMockPutBucketCorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketCors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketCorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketCorsResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod></CORSRule></CORSConfiguration>")
		},
		&PutBucketCorsRequest{
			Bucket: Ptr("bucket"),
			CORSConfiguration: &CORSConfiguration{
				CORSRules: []CORSRule{
					{
						AllowedOrigins: []string{"*"},
						AllowedMethods: []string{"PUT", "GET"},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketCorsResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CORSConfiguration><CORSRule><AllowedOrigin>*</AllowedOrigin><AllowedMethod>PUT</AllowedMethod><AllowedMethod>GET</AllowedMethod></CORSRule></CORSConfiguration>")
		},
		&PutBucketCorsRequest{
			Bucket: Ptr("bucket"),
			CORSConfiguration: &CORSConfiguration{
				CORSRules: []CORSRule{
					{
						AllowedOrigins: []string{"*"},
						AllowedMethods: []string{"PUT", "GET"},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketCorsResult, err error) {
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

func TestMockPutBucketCors_Error(t *testing.T) {
	for _, c := range testMockPutBucketCorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketCors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketCorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketCorsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CORSConfiguration>
    <CORSRule>
      <AllowedOrigin>*</AllowedOrigin>
      <AllowedMethod>PUT</AllowedMethod>
      <AllowedMethod>GET</AllowedMethod>
      <AllowedHeader>Authorization</AllowedHeader>
    </CORSRule>
    <CORSRule>
      <AllowedOrigin>http://example.com</AllowedOrigin>
      <AllowedOrigin>http://example.net</AllowedOrigin>
      <AllowedMethod>GET</AllowedMethod>
      <AllowedHeader>Authorization</AllowedHeader>
      <ExposeHeader>x-oss-test</ExposeHeader>
      <ExposeHeader>x-oss-test1</ExposeHeader>
      <MaxAgeSeconds>100</MaxAgeSeconds>
    </CORSRule>
    <ResponseVary>false</ResponseVary>
</CORSConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
		},
		&GetBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketCorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.CORSConfiguration.CORSRules), 2)
			assert.Equal(t, *o.CORSConfiguration.ResponseVary, false)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedOrigins[0], "*")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[0].AllowedMethods), 2)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedMethods[0], "PUT")
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedMethods[1], "GET")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[0].AllowedHeaders), 1)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedHeaders[0], "Authorization")
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].AllowedOrigins[0], "http://example.com")
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].AllowedOrigins[1], "http://example.net")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[1].AllowedMethods), 1)
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].AllowedMethods[0], "GET")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[1].AllowedHeaders), 1)
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].AllowedHeaders[0], "Authorization")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[1].ExposeHeaders), 2)
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].ExposeHeaders[0], "x-oss-test")
			assert.Equal(t, o.CORSConfiguration.CORSRules[1].ExposeHeaders[1], "x-oss-test1")
			assert.Equal(t, *o.CORSConfiguration.CORSRules[1].MaxAgeSeconds, int64(100))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<CORSConfiguration>
    <CORSRule>
      <AllowedOrigin>*</AllowedOrigin>
      <AllowedMethod>PUT</AllowedMethod>
      <AllowedMethod>GET</AllowedMethod>
      <AllowedHeader>Authorization</AllowedHeader>
    </CORSRule>
</CORSConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
		},
		&GetBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketCorsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.CORSConfiguration.CORSRules), 1)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedOrigins[0], "*")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[0].AllowedMethods), 2)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedMethods[0], "PUT")
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedMethods[1], "GET")
			assert.Equal(t, len(o.CORSConfiguration.CORSRules[0].AllowedHeaders), 1)
			assert.Equal(t, o.CORSConfiguration.CORSRules[0].AllowedHeaders[0], "Authorization")
		},
	},
}

func TestMockGetBucketCors_Success(t *testing.T) {
	for _, c := range testMockGetBucketCorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketCors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketCorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketCorsResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", urlStr)
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketCorsResult, err error) {
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
			assert.Equal(t, "/bucket/?cors", strUrl)
		},
		&GetBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketCorsResult, err error) {
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
			assert.Equal(t, "/bucket/?cors", strUrl)
		},
		&GetBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketCorsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketCors fail")
		},
	},
}

func TestMockGetBucketCors_Error(t *testing.T) {
	for _, c := range testMockGetBucketCorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketCors(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketCorsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketCorsResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", strUrl)
		},
		&DeleteBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketCorsResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketCors_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketCorsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketCors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketCorsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketCorsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketCorsResult, err error)
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
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", strUrl)
		},
		&DeleteBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketCorsResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?cors", strUrl)
		},
		&DeleteBucketCorsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketCorsResult, err error) {
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

func TestMockDeleteBucketCors_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketCorsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketCors(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockOptionObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OptionObjectRequest
	CheckOutputFn  func(t *testing.T, o *OptionObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":              "534B371674E88A4D8906****",
			"Date":                          "Fri, 24 Feb 2017 03:15:40 GMT",
			"Access-Control-Allow-Origin":   "http://www.example.com",
			"Access-Control-Allow-Methods":  "PUT",
			"Access-Control-Expose-Headers": "x-oss-test1,x-oss-test2",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "OPTIONS", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/oss-object", strUrl)
			assert.Equal(t, "http://www.example.com", r.Header.Get("Origin"))
			assert.Equal(t, "x-oss-test1,x-oss-test2", r.Header.Get("Access-Control-Request-Headers"))
			assert.Equal(t, "PUT", r.Header.Get("Access-Control-Request-Method"))
		},
		&OptionObjectRequest{
			Bucket:                      Ptr("bucket"),
			Key:                         Ptr("oss-object"),
			Origin:                      Ptr("http://www.example.com"),
			AccessControlRequestHeaders: Ptr("x-oss-test1,x-oss-test2"),
			AccessControlRequestMethod:  Ptr("PUT"),
		},
		func(t *testing.T, o *OptionObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "http://www.example.com", *o.AccessControlAllowOrigin)
			assert.Equal(t, "x-oss-test1,x-oss-test2", *o.AccessControlExposeHeaders)
			assert.Equal(t, "PUT", *o.AccessControlAllowMethods)
		},
	},
}

func TestMockOptionObjectSuccess(t *testing.T) {
	for _, c := range testMockOptionObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OptionObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockOptionObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OptionObjectRequest
	CheckOutputFn  func(t *testing.T, o *OptionObjectResult, err error)
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
			assert.Equal(t, "OPTIONS", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/oss-object", strUrl)
			assert.Equal(t, "http://www.example.com", r.Header.Get("Origin"))
			assert.Equal(t, "x-oss-test1,x-oss-test2", r.Header.Get("Access-Control-Request-Headers"))
			assert.Equal(t, "PUT", r.Header.Get("Access-Control-Request-Method"))
		},
		&OptionObjectRequest{
			Bucket:                      Ptr("bucket"),
			Key:                         Ptr("oss-object"),
			Origin:                      Ptr("http://www.example.com"),
			AccessControlRequestHeaders: Ptr("x-oss-test1,x-oss-test2"),
			AccessControlRequestMethod:  Ptr("PUT"),
		},
		func(t *testing.T, o *OptionObjectResult, err error) {
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
			assert.Equal(t, "OPTIONS", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/oss-object", strUrl)
			assert.Equal(t, "http://www.example.com", r.Header.Get("Origin"))
			assert.Equal(t, "x-oss-test1,x-oss-test2", r.Header.Get("Access-Control-Request-Headers"))
			assert.Equal(t, "PUT", r.Header.Get("Access-Control-Request-Method"))
		},
		&OptionObjectRequest{
			Bucket:                      Ptr("bucket"),
			Key:                         Ptr("oss-object"),
			Origin:                      Ptr("http://www.example.com"),
			AccessControlRequestHeaders: Ptr("x-oss-test1,x-oss-test2"),
			AccessControlRequestMethod:  Ptr("PUT"),
		},
		func(t *testing.T, o *OptionObjectResult, err error) {
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

func TestMockOptionObject_Error(t *testing.T) {
	for _, c := range testMockOptionObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OptionObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


