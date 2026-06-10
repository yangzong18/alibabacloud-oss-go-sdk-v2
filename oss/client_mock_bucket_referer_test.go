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

var testMockPutBucketRefererSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRefererRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRefererResult, err error)
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
			assert.Equal(t, "/bucket/?referer", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><AllowTruncateQueryString>true</AllowTruncateQueryString><TruncatePath>true</TruncatePath><RefererList><Referer>http://www.aliyun.com</Referer><Referer>https://www.aliyun.com</Referer><Referer>http://www.*.com</Referer><Referer>https://www.?.aliyuncs.com</Referer></RefererList><RefererBlacklist><Referer>http://www.refuse.com</Referer><Referer>https://*.hack.com</Referer><Referer>http://ban.*.com</Referer><Referer>https://www.?.deny.com</Referer></RefererBlacklist></RefererConfiguration>")
		},
		&PutBucketRefererRequest{
			Bucket: Ptr("bucket"),
			RefererConfiguration: &RefererConfiguration{
				AllowEmptyReferer:        Ptr(false),
				AllowTruncateQueryString: Ptr(true),
				TruncatePath:             Ptr(true),
				RefererList: &RefererList{
					[]string{
						"http://www.aliyun.com", "https://www.aliyun.com", "http://www.*.com", "https://www.?.aliyuncs.com",
					},
				},
				RefererBlacklist: &RefererBlacklist{
					[]string{
						"http://www.refuse.com", "https://*.hack.com", "http://ban.*.com", "https://www.?.deny.com",
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketRefererResult, err error) {
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
			assert.Equal(t, "/bucket/?referer", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><AllowTruncateQueryString>true</AllowTruncateQueryString><TruncatePath>true</TruncatePath><RefererList><Referer>http://www.aliyun.com</Referer><Referer>https://www.aliyun.com</Referer><Referer>http://www.*.com</Referer><Referer>https://www.?.aliyuncs.com</Referer></RefererList></RefererConfiguration>")
		},
		&PutBucketRefererRequest{
			Bucket: Ptr("bucket"),
			RefererConfiguration: &RefererConfiguration{
				AllowEmptyReferer:        Ptr(false),
				AllowTruncateQueryString: Ptr(true),
				TruncatePath:             Ptr(true),
				RefererList: &RefererList{
					[]string{
						"http://www.aliyun.com", "https://www.aliyun.com", "http://www.*.com", "https://www.?.aliyuncs.com",
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketRefererResult, err error) {
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
			assert.Equal(t, "/bucket/?referer", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><RefererList><Referer></Referer></RefererList></RefererConfiguration>")
		},
		&PutBucketRefererRequest{
			Bucket: Ptr("bucket"),
			RefererConfiguration: &RefererConfiguration{
				AllowEmptyReferer: Ptr(false),
				RefererList: &RefererList{
					Referers: []string{""},
				},
			},
		},
		func(t *testing.T, o *PutBucketRefererResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketReferer_Success(t *testing.T) {
	for _, c := range testMockPutBucketRefererSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketReferer(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketRefererErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRefererRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketRefererResult, err error)
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
			assert.Equal(t, "/bucket/?referer", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><RefererList><Referer></Referer></RefererList></RefererConfiguration>")
		},
		&PutBucketRefererRequest{
			Bucket: Ptr("bucket"),
			RefererConfiguration: &RefererConfiguration{
				AllowEmptyReferer: Ptr(false),
				RefererList: &RefererList{
					Referers: []string{""},
				},
			},
		},
		func(t *testing.T, o *PutBucketRefererResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?referer", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RefererConfiguration><AllowEmptyReferer>false</AllowEmptyReferer><RefererList><Referer></Referer></RefererList></RefererConfiguration>")
		},
		&PutBucketRefererRequest{
			Bucket: Ptr("bucket"),
			RefererConfiguration: &RefererConfiguration{
				AllowEmptyReferer: Ptr(false),
				RefererList: &RefererList{
					Referers: []string{""},
				},
			},
		},
		func(t *testing.T, o *PutBucketRefererResult, err error) {
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

func TestMockPutBucketReferer_Error(t *testing.T) {
	for _, c := range testMockPutBucketRefererErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketReferer(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketRefererSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketRefererRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketRefererResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<RefererConfiguration>
  <AllowEmptyReferer>true</AllowEmptyReferer>
</RefererConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?referer", r.URL.String())
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.True(t, *o.RefererConfiguration.AllowEmptyReferer)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<RefererConfiguration>
  <AllowEmptyReferer>false</AllowEmptyReferer>
  <AllowTruncateQueryString>true</AllowTruncateQueryString>
  <TruncatePath>true</TruncatePath>
  <RefererList>
    <Referer>http://www.aliyun.com</Referer>
    <Referer>https://www.aliyun.com</Referer>
    <Referer>http://www.*.com</Referer>
    <Referer>https://www.?.aliyuncs.com</Referer>
  </RefererList>
</RefererConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?referer", r.URL.String())
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.False(t, *o.RefererConfiguration.AllowEmptyReferer)
			assert.True(t, *o.RefererConfiguration.AllowTruncateQueryString)
			assert.True(t, *o.RefererConfiguration.TruncatePath)
			assert.Equal(t, len(o.RefererConfiguration.RefererList.Referers), 4)

			assert.Equal(t, o.RefererConfiguration.RefererList.Referers[3], "https://www.?.aliyuncs.com")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<RefererConfiguration>
  <AllowEmptyReferer>false</AllowEmptyReferer>
  <AllowTruncateQueryString>true</AllowTruncateQueryString>
  <TruncatePath>true</TruncatePath>
  <RefererList>
    <Referer>http://www.aliyun.com</Referer>
    <Referer>https://www.aliyun.com</Referer>
    <Referer>http://www.*.com</Referer>
    <Referer>https://www.?.aliyuncs.com</Referer>
  </RefererList>
  <RefererBlacklist>
    <Referer>http://www.refuse.com</Referer>
    <Referer>https://*.hack.com</Referer>
    <Referer>http://ban.*.com</Referer>
    <Referer>https://www.?.deny.com</Referer>
  </RefererBlacklist>
</RefererConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?referer", r.URL.String())
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.False(t, *o.RefererConfiguration.AllowEmptyReferer)
			assert.True(t, *o.RefererConfiguration.AllowTruncateQueryString)
			assert.True(t, *o.RefererConfiguration.TruncatePath)
			assert.Equal(t, len(o.RefererConfiguration.RefererList.Referers), 4)
			assert.Equal(t, len(o.RefererConfiguration.RefererBlacklist.Referers), 4)
			assert.Equal(t, o.RefererConfiguration.RefererList.Referers[3], "https://www.?.aliyuncs.com")
			assert.Equal(t, o.RefererConfiguration.RefererBlacklist.Referers[2], "http://ban.*.com")
		},
	},
}

func TestMockGetBucketReferer_Success(t *testing.T) {
	for _, c := range testMockGetBucketRefererSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReferer(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketRefererErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketRefererRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketRefererResult, err error)
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
			assert.Equal(t, "/bucket/?referer", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
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
			assert.Equal(t, "/bucket/?referer", strUrl)
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
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
			assert.Equal(t, "/bucket/?referer", strUrl)
		},
		&GetBucketRefererRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketRefererResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketReferer fail")
		},
	},
}

func TestMockGetBucketReferer_Error(t *testing.T) {
	for _, c := range testMockGetBucketRefererErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketReferer(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}


