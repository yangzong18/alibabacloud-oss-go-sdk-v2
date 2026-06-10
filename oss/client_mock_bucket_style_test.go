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

var testMockPutBucketStyleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutStyleRequest
	CheckOutputFn  func(t *testing.T, o *PutStyleResult, err error)
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
			assert.Equal(t, "/bucket/?category=image&style&styleName=imageStyle", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<Style><Content>image/resize,p_50</Content></Style>")
		},
		&PutStyleRequest{
			Bucket:    Ptr("bucket"),
			Category:  Ptr("image"),
			StyleName: Ptr("imageStyle"),
			Style: &StyleContent{
				Ptr("image/resize,p_50"),
			},
		},
		func(t *testing.T, o *PutStyleResult, err error) {
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
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<Style><Content>image/resize,p_50</Content></Style>")
		},
		&PutStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
			Style: &StyleContent{
				Ptr("image/resize,p_50"),
			},
		},
		func(t *testing.T, o *PutStyleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketStyle_Success(t *testing.T) {
	for _, c := range testMockPutBucketStyleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketStyleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutStyleRequest
	CheckOutputFn  func(t *testing.T, o *PutStyleResult, err error)
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
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<Style><Content>image/resize,p_50</Content></Style>")
		},
		&PutStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
			Style: &StyleContent{
				Ptr("image/resize,p_50"),
			},
		},
		func(t *testing.T, o *PutStyleResult, err error) {
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
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<Style><Content>image/resize,p_50</Content></Style>")
		},
		&PutStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
			Style: &StyleContent{
				Ptr("image/resize,p_50"),
			},
		},
		func(t *testing.T, o *PutStyleResult, err error) {
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

func TestMockPutBucketStyle_Error(t *testing.T) {
	for _, c := range testMockPutBucketStyleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStyleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetStyleRequest
	CheckOutputFn  func(t *testing.T, o *GetStyleResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<Style>
 <Name>imageStyle</Name>
 <Content>image/resize,p_50</Content>
 <Category>image</Category>
 <CreateTime>Wed, 20 May 2020 12:07:15 GMT</CreateTime>
 <LastModifyTime>Wed, 21 May 2020 12:07:15 GMT</LastModifyTime>
</Style>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
		},
		&GetStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
		},
		func(t *testing.T, o *GetStyleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Style.Name, "imageStyle")
			assert.Equal(t, *o.Style.Category, "image")
			assert.Equal(t, *o.Style.Content, "image/resize,p_50")
			assert.Equal(t, *o.Style.LastModifyTime, "Wed, 21 May 2020 12:07:15 GMT")
			assert.Equal(t, *o.Style.CreateTime, "Wed, 20 May 2020 12:07:15 GMT")
		},
	},
}

func TestMockGetBucketStyle_Success(t *testing.T) {
	for _, c := range testMockGetBucketStyleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStyleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetStyleRequest
	CheckOutputFn  func(t *testing.T, o *GetStyleResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
		},
		&GetStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
		},
		func(t *testing.T, o *GetStyleResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
		},
		&GetStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
		},
		func(t *testing.T, o *GetStyleResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?style&styleName=imageStyle", urlStr)
		},
		&GetStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("imageStyle"),
		},
		func(t *testing.T, o *GetStyleResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetStyle fail")
		},
	},
}

func TestMockGetBucketStyle_Error(t *testing.T) {
	for _, c := range testMockGetBucketStyleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetStyle(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketStyleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListStyleRequest
	CheckOutputFn  func(t *testing.T, o *ListStyleResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<StyleList>
 <Style>
 <Name>imagestyle</Name>
 <Content>image/resize,p_50</Content>
 <Category>image</Category>
 <CreateTime>Wed, 20 May 2020 12:07:15 GMT</CreateTime>
 <LastModifyTime>Wed, 21 May 2020 12:07:15 GMT</LastModifyTime>
 </Style>
 <Style>
 <Name>imagestyle1</Name>
 <Content>image/resize,w_200</Content>
 <Category>image</Category>
 <CreateTime>Wed, 20 May 2020 12:08:04 GMT</CreateTime>
 <LastModifyTime>Wed, 21 May 2020 12:08:04 GMT</LastModifyTime>
 </Style>
 <Style>
 <Name>imagestyle2</Name>
 <Content>image/resize,w_300</Content>
 <Category>image</Category>
 <CreateTime>Fri, 12 Mar 2021 06:19:13 GMT</CreateTime>
 <LastModifyTime>Fri, 13 Mar 2021 06:27:21 GMT</LastModifyTime>
 </Style>
</StyleList>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?style", urlStr)
		},
		&ListStyleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListStyleResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.StyleList.Styles), 3)
			assert.Equal(t, *o.StyleList.Styles[0].Name, "imagestyle")
			assert.Equal(t, *o.StyleList.Styles[0].Category, "image")
			assert.Equal(t, *o.StyleList.Styles[0].Content, "image/resize,p_50")
			assert.Equal(t, *o.StyleList.Styles[0].LastModifyTime, "Wed, 21 May 2020 12:07:15 GMT")
			assert.Equal(t, *o.StyleList.Styles[0].CreateTime, "Wed, 20 May 2020 12:07:15 GMT")
			assert.Equal(t, *o.StyleList.Styles[1].Name, "imagestyle1")
			assert.Equal(t, *o.StyleList.Styles[1].Category, "image")
			assert.Equal(t, *o.StyleList.Styles[1].Content, "image/resize,w_200")
			assert.Equal(t, *o.StyleList.Styles[1].LastModifyTime, "Wed, 21 May 2020 12:08:04 GMT")
			assert.Equal(t, *o.StyleList.Styles[1].CreateTime, "Wed, 20 May 2020 12:08:04 GMT")
			assert.Equal(t, *o.StyleList.Styles[2].Name, "imagestyle2")
			assert.Equal(t, *o.StyleList.Styles[2].Category, "image")
			assert.Equal(t, *o.StyleList.Styles[2].Content, "image/resize,w_300")
			assert.Equal(t, *o.StyleList.Styles[2].LastModifyTime, "Fri, 13 Mar 2021 06:27:21 GMT")
			assert.Equal(t, *o.StyleList.Styles[2].CreateTime, "Fri, 12 Mar 2021 06:19:13 GMT")
		},
	},
}

func TestMockListBucketStyle_Success(t *testing.T) {
	for _, c := range testMockListBucketStyleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketStyleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListStyleRequest
	CheckOutputFn  func(t *testing.T, o *ListStyleResult, err error)
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
			assert.Equal(t, "/bucket/?style", urlStr)
			assert.Equal(t, "GET", r.Method)
		},
		&ListStyleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListStyleResult, err error) {
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
			assert.Equal(t, "/bucket/?style", strUrl)
		},
		&ListStyleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListStyleResult, err error) {
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
			assert.Equal(t, "/bucket/?style", strUrl)
		},
		&ListStyleRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListStyleResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListStyle fail")
		},
	},
}

func TestMockListBucketStyle_Error(t *testing.T) {
	for _, c := range testMockListBucketStyleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListStyle(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketStyleSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteStyleRequest
	CheckOutputFn  func(t *testing.T, o *DeleteStyleResult, err error)
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
			assert.Equal(t, "/bucket/?style&styleName=styleName", strUrl)
		},
		&DeleteStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("styleName"),
		},
		func(t *testing.T, o *DeleteStyleResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketStyle_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketStyleSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketStyleErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteStyleRequest
	CheckOutputFn  func(t *testing.T, o *DeleteStyleResult, err error)
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
			assert.Equal(t, "/bucket/?style&styleName=styleName", strUrl)
		},
		&DeleteStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("styleName"),
		},
		func(t *testing.T, o *DeleteStyleResult, err error) {
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
			assert.Equal(t, "/bucket/?style&styleName=styleName", strUrl)
		},
		&DeleteStyleRequest{
			Bucket:    Ptr("bucket"),
			StyleName: Ptr("styleName"),
		},
		func(t *testing.T, o *DeleteStyleResult, err error) {
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

func TestMockDeleteStyle_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketStyleErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteStyle(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


