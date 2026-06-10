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

var testMockPutBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLoggingResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket><TargetPrefix>TargetPrefix</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
					TargetPrefix: Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
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
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
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
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket><TargetPrefix>TargetPrefix</TargetPrefix><LoggingRole>AliyunOSSLoggingDefaultRole</LoggingRole></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
					TargetPrefix: Ptr("TargetPrefix"),
					LoggingRole:  Ptr("AliyunOSSLoggingDefaultRole"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketLogging_Success(t *testing.T) {
	for _, c := range testMockPutBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketLoggingResult, err error)
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
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket><TargetPrefix>TargetPrefix</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
					TargetPrefix: Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket><TargetPrefix>TargetPrefix</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
					TargetPrefix: Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<BucketLoggingStatus><LoggingEnabled><TargetBucket>TargetBucket</TargetBucket><TargetPrefix>TargetPrefix</TargetPrefix></LoggingEnabled></BucketLoggingStatus>")
		},
		&PutBucketLoggingRequest{
			Bucket: Ptr("bucket"),
			BucketLoggingStatus: &BucketLoggingStatus{
				&LoggingEnabled{
					TargetBucket: Ptr("TargetBucket"),
					TargetPrefix: Ptr("TargetPrefix"),
				},
			},
		},
		func(t *testing.T, o *PutBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketLogging fail")
		},
	},
}

func TestMockPutBucketLogging_Error(t *testing.T) {
	for _, c := range testMockPutBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLoggingResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus>
  <LoggingEnabled>
        <TargetBucket>bucket-log</TargetBucket>
        <TargetPrefix>prefix-access_log</TargetPrefix>
    </LoggingEnabled>
</BucketLoggingStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?logging", r.URL.String())
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetBucket, "bucket-log")
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "prefix-access_log")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus>
  <LoggingEnabled>
        <TargetBucket>bucket-log</TargetBucket>
    </LoggingEnabled>
</BucketLoggingStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?logging", r.URL.String())
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetBucket, "bucket-log")
			assert.Nil(t, o.BucketLoggingStatus.LoggingEnabled.TargetPrefix)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<BucketLoggingStatus>
  <LoggingEnabled>
        <TargetBucket>bucket-log</TargetBucket>
        <TargetPrefix>prefix-access_log</TargetPrefix>
		<LoggingRole>AliyunOSSLoggingDefaultRole</LoggingRole>
    </LoggingEnabled>
</BucketLoggingStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?logging", r.URL.String())
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.LoggingRole, "AliyunOSSLoggingDefaultRole")
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetBucket, "bucket-log")
			assert.Equal(t, *o.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "prefix-access_log")
		},
	},
}

func TestMockGetBucketLogging_Success(t *testing.T) {
	for _, c := range testMockGetBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLoggingResult, err error)
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
			assert.Equal(t, "/bucket/?logging", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
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
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
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
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&GetBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLoggingResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketLogging fail")
		},
	},
}

func TestMockGetBucketLogging_Error(t *testing.T) {
	for _, c := range testMockGetBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLoggingSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLoggingResult, err error)
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
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucketLogging_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketLoggingSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketLoggingErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketLoggingRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketLoggingResult, err error)
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
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
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
			assert.Equal(t, "/bucket/?logging", strUrl)
		},
		&DeleteBucketLoggingRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketLoggingResult, err error) {
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

func TestMockDeleteBucketLogging_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketLoggingErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketLogging(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutUserDefinedLogFieldsConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<UserDefinedLogFieldsConfiguration><HeaderSet><header>header1</header><header>header2</header><header>header3</header></HeaderSet><ParamSet><parameter>param1</parameter><parameter>param2</parameter></ParamSet></UserDefinedLogFieldsConfiguration>")
		},
		&PutUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
			UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
				HeaderSet: &LoggingHeaderSet{
					[]string{"header1", "header2", "header3"},
				},
				ParamSet: &LoggingParamSet{
					[]string{"param1", "param2"},
				},
			},
		},
		func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<UserDefinedLogFieldsConfiguration><HeaderSet><header>header1</header></HeaderSet></UserDefinedLogFieldsConfiguration>")
		},
		&PutUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
			UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
				HeaderSet: &LoggingHeaderSet{
					[]string{"header1"},
				},
			},
		},
		func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutUserDefinedLogFieldsConfig_Success(t *testing.T) {
	for _, c := range testMockPutUserDefinedLogFieldsConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutUserDefinedLogFieldsConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<UserDefinedLogFieldsConfiguration><HeaderSet><header>header1</header><header>header2</header><header>header3</header></HeaderSet><ParamSet><parameter>param1</parameter><parameter>param2</parameter></ParamSet></UserDefinedLogFieldsConfiguration>")
		},
		&PutUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
			UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
				HeaderSet: &LoggingHeaderSet{
					[]string{"header1", "header2", "header3"},
				},
				ParamSet: &LoggingParamSet{
					[]string{"param1", "param2"},
				},
			},
		},
		func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<UserDefinedLogFieldsConfiguration><HeaderSet><header>header1</header><header>header2</header><header>header3</header></HeaderSet><ParamSet><parameter>param1</parameter><parameter>param2</parameter></ParamSet></UserDefinedLogFieldsConfiguration>")
		},
		&PutUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
			UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
				HeaderSet: &LoggingHeaderSet{
					[]string{"header1", "header2", "header3"},
				},
				ParamSet: &LoggingParamSet{
					[]string{"param1", "param2"},
				},
			},
		},
		func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<UserDefinedLogFieldsConfiguration><HeaderSet><header>header1</header><header>header2</header><header>header3</header></HeaderSet><ParamSet><parameter>param1</parameter><parameter>param2</parameter></ParamSet></UserDefinedLogFieldsConfiguration>")
		},
		&PutUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
			UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
				HeaderSet: &LoggingHeaderSet{
					[]string{"header1", "header2", "header3"},
				},
				ParamSet: &LoggingParamSet{
					[]string{"param1", "param2"},
				},
			},
		},
		func(t *testing.T, o *PutUserDefinedLogFieldsConfigResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutUserDefinedLogFieldsConfig fail")
		},
	},
}

func TestMockPutUserDefinedLogFieldsConfig_Error(t *testing.T) {
	for _, c := range testMockPutUserDefinedLogFieldsConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetUserDefinedLogFieldsConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<UserDefinedLogFieldsConfiguration>
	<HeaderSet>
		<header>header1</header>
		<header>header2</header>
		<header>header3</header>
	</HeaderSet>
	<ParamSet>
		<parameter>param1</parameter>
		<parameter>param2</parameter>
	</ParamSet>
</UserDefinedLogFieldsConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
		},
		&GetUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, 3, len(o.UserDefinedLogFieldsConfiguration.HeaderSet.Headers))
			assert.Equal(t, "header3", o.UserDefinedLogFieldsConfiguration.HeaderSet.Headers[2])
			assert.Equal(t, 2, len(o.UserDefinedLogFieldsConfiguration.ParamSet.Parameters))
			assert.Equal(t, "param2", o.UserDefinedLogFieldsConfiguration.ParamSet.Parameters[1])
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<UserDefinedLogFieldsConfiguration>
	<HeaderSet>
		<header>header1</header>
	</HeaderSet>
</UserDefinedLogFieldsConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
		},
		&GetUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, 1, len(o.UserDefinedLogFieldsConfiguration.HeaderSet.Headers))
			assert.Equal(t, "header1", o.UserDefinedLogFieldsConfiguration.HeaderSet.Headers[0])
			assert.Nil(t, o.UserDefinedLogFieldsConfiguration.ParamSet)
		},
	},
}

func TestMockGetUserDefinedLogFieldsConfig_Success(t *testing.T) {
	for _, c := range testMockGetUserDefinedLogFieldsConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetUserDefinedLogFieldsConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&GetUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&GetUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetUserDefinedLogFieldsConfigResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetUserDefinedLogFieldsConfig fail")
		},
	},
}

func TestMockGetUserDefinedLogFieldsConfig_Error(t *testing.T) {
	for _, c := range testMockGetUserDefinedLogFieldsConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteUserDefinedLogFieldsConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&DeleteUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-version-id": "CAEQNRiBgMClj7qD0BYiIDQ5Y2QyMjc3NGZkODRlMTU5M2VkY2U3MWRiNGRh****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&DeleteUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteUserDefinedLogFieldsConfig_Success(t *testing.T) {
	for _, c := range testMockDeleteUserDefinedLogFieldsConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteUserDefinedLogFieldsConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteUserDefinedLogFieldsConfigRequest
	CheckOutputFn  func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&DeleteUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?userDefinedLogFieldsConfig", strUrl)
		},
		&DeleteUserDefinedLogFieldsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteUserDefinedLogFieldsConfigResult, err error) {
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

func TestMockDeleteUserDefinedLogFieldsConfig_Error(t *testing.T) {
	for _, c := range testMockDeleteUserDefinedLogFieldsConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteUserDefinedLogFieldsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


