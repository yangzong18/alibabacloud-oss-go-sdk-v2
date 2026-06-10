package oss

import (
	"testing"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockProcessObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *ProcessObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "bucket": "",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "bucket": "dest-bucket",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=,b_ZGVzdC1idWNrZXQ=", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v,b_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")), base64.URLEncoding.EncodeToString([]byte("dest-bucket")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "dest-bucket")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{
    "bucket": "",
    "fileSize": 3267,
    "object": "dest.jpg",
    "status": "OK"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_ZGVzdC5qcGc=", string(data))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			Process:      Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("dest.jpg")))),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Bucket, "")
			assert.Equal(t, o.FileSize, 3267)
			assert.Equal(t, o.Object, "dest.jpg")
			assert.Equal(t, o.ProcessStatus, "OK")
		},
	},
}

func TestMockProcessObject_Success(t *testing.T) {
	for _, c := range testMockProcessObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockProcessObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *ProcessObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
		},
	},
	{
		400,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>operation not support post: test</Message>
  <RequestId>65467C42E001B4333337****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "operation not support post: test")
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
			assert.Equal(t, "/bucket/object?x-oss-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-process=image/resize,w_100|sys/saveas,o_a2V5LWRlc3QuanBn", string(data))
		},
		&ProcessObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr(fmt.Sprintf("image/resize,w_100|sys/saveas,o_%v", base64.URLEncoding.EncodeToString([]byte("key-dest.jpg")))),
		},
		func(t *testing.T, o *ProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ProcessObject fail")
		},
	},
}

func TestMockProcessObject_Error(t *testing.T) {
	for _, c := range testMockProcessObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAsyncProcessObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AsyncProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *AsyncProcessObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"EventId":"181-1kZUlN60OH4fWOcOjZEnGnG****","RequestId":"1D99637F-F59E-5B41-9200-C4892F52****","TaskId":"MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.EventId, "181-1kZUlN60OH4fWOcOjZEnGnG****")
			assert.Equal(t, o.RequestId, "1D99637F-F59E-5B41-9200-C4892F52****")
			assert.Equal(t, o.TaskId, "MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/json",
		},
		[]byte(`{"EventId":"181-1kZUlN60OH4fWOcOjZEnGnG****","RequestId":"1D99637F-F59E-5B41-9200-C4892F52****","TaskId":"MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.EventId, "181-1kZUlN60OH4fWOcOjZEnGnG****")
			assert.Equal(t, o.RequestId, "1D99637F-F59E-5B41-9200-C4892F52****")
			assert.Equal(t, o.TaskId, "MediaConvert-e4a737df-69e9-4fca-8d9b-17c40ea3****")
		},
	},
}

func TestMockAsyncProcessObject_Success(t *testing.T) {
	for _, c := range testMockAsyncProcessObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AsyncProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAsyncProcessObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AsyncProcessObjectRequest
	CheckOutputFn  func(t *testing.T, o *AsyncProcessObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
		},
	},
	{
		400,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidRequest</Code>
  <Message>no x-oss-async-process parameter found</Message>
  <RequestId>65467C42E001B4333337****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidRequest", serr.Code)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "no x-oss-async-process parameter found")
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
			assert.Equal(t, "/bucket/object?x-oss-async-process", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, "x-oss-async-process=video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1|sys/saveas,b_ZGVzY3QtYnVja2V0,o_ZGVtby5tcDQ", string(data))
		},
		&AsyncProcessObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			AsyncProcess: Ptr(fmt.Sprintf("%s|sys/saveas,b_%v,o_%v", "video/convert,f_avi,vcodec_h265,s_1920x1080,vb_2000000,fps_30,acodec_aac,ab_100000,sn_1", strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("desct-bucket")), "="), strings.TrimRight(base64.URLEncoding.EncodeToString([]byte("demo.mp4")), "="))),
		},
		func(t *testing.T, o *AsyncProcessObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute AsyncProcessObject fail")
		},
	},
}

func TestMockAsyncProcessObject_Error(t *testing.T) {
	for _, c := range testMockAsyncProcessObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AsyncProcessObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


