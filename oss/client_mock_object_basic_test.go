package oss

import (
	"testing"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`{"filename":"object","size":"6","mimeType":""}`),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
			Body:     strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			jsonData, err := json.Marshal(o.CallbackResult)
			assert.Nil(t, err)
			assert.NotEmpty(t, string(jsonData))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&PutObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			TrafficLimit: int64(100 * 1024 * 8),
			Body:         strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "text/plain", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object.txt"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "my-content-type", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object.txt"),
			Body:        strings.NewReader("hi oss"),
			ContentType: Ptr("my-content-type"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object.txt", r.URL.String())
			assert.Equal(t, "my-content-type", r.Header.Get(HTTPHeaderContentType))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&PutObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object.txt"),
			Body:         strings.NewReader("hi oss"),
			ContentType:  Ptr("my-content-type"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
		},
	},
}

func TestMockPutObject_Success(t *testing.T) {
	for _, c := range testMockPutObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockPutObjectDisableDetectMimeTypeCases {
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

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "8707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_Progress(t *testing.T) {
	for _, c := range testMockPutObjectWithProgressCases {
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
		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss")), n)
	}
}

var testMockPutObjectWithCrcDisableCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "6707180448768400016",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "6707180448768400016")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockPutObject_DisableCRC64(t *testing.T) {
	//Disable
	for _, c := range testMockPutObjectWithCrcDisableCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		n := int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
		}
		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss")), n)

		//Enable Got Fail
		client = NewClient(cfg)
		assert.NotNil(t, c)
		n = int64(0)
		c.Request.ProgressFn = func(increment, transferred, total int64) {
			n = transferred
		}
		c.Request.Body = strings.NewReader("hi oss")
		_, err = client.PutObject(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Equal(t, int64(len("hi oss")), n)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 8707180448768400016, server 6707180448768400016")
	}
}

var testMockPutObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Body:   strings.NewReader("hi oss"),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
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
		203,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "5C3D9175B6FC201293AD****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "8707180448768400016",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>CallbackFailed</Code>
  <Message>Error status : 301.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0007-00000203</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0007-00000203</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Body:     strings.NewReader("hi oss"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(203), serr.StatusCode)
			assert.Equal(t, "CallbackFailed", serr.Code)
			assert.Equal(t, "Error status : 301.", serr.Message)
			assert.Equal(t, "0007-00000203", serr.EC)
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Body:     strings.NewReader("hi oss"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutObject fail")
		},
	},
}

func TestMockPutObject_Error(t *testing.T) {
	for _, c := range testMockPutObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "316181249502703****",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(`hi oss,this is a demo!`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "316181249502703****")
			content, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(content), "hi oss,this is a demo!")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":                      "344606",
			"Last-Modified":                       "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":                   "Normal",
			"Accept-Ranges":                       "bytes",
			"Content-disposition":                 "attachment; filename=testing.txt",
			"Cache-control":                       "no-cache",
			"X-Oss-Storage-Class":                 "Standard",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-tagging-count":                 "2",
			"Content-MD5":                         "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-hash-crc64ecma":                "870718044876840****",
			"x-oss-meta-name":                     "demo",
			"x-oss-meta-email":                    "demo@aliyun.com",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "text")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, o.TaggingCount, int32(2))
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":                      "344606",
			"Last-Modified":                       "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":                   "Normal",
			"Accept-Ranges":                       "bytes",
			"Content-disposition":                 "attachment; filename=testing.txt",
			"Cache-control":                       "no-cache",
			"X-Oss-Storage-Class":                 "Standard",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-tagging-count":                 "2",
			"Content-MD5":                         "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-hash-crc64ecma":                "870718044876840****",
			"x-oss-meta-name":                     "demo",
			"x-oss-meta-email":                    "demo@aliyun.com",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&GetObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			TrafficLimit: int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "text")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, o.TaggingCount, int32(2))
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "image/jpeg",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                 "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":       "344606",
			"Last-Modified":        "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":    "Normal",
			"X-Oss-Storage-Class":  "Standard",
			"x-oss-hash-crc64ecma": "870718044876840****",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?x-oss-process=image%2Fresize%2Cm_fixed%2Cw_100%2Ch_100", r.URL.String())
		},
		&GetObjectRequest{
			Bucket:  Ptr("bucket"),
			Key:     Ptr("object"),
			Process: Ptr("image/resize,m_fixed,w_100,h_100"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "image/jpeg")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "image/jpeg",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                 "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":       "344606",
			"Last-Modified":        "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":    "Normal",
			"X-Oss-Storage-Class":  "Standard",
			"x-oss-hash-crc64ecma": "870718044876840****",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "image/jpeg")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := io.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
}

func TestMockGetObject_Success(t *testing.T) {
	for _, c := range testMockGetObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
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
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
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

func TestMockGetObject_Error(t *testing.T) {
	for _, c := range testMockGetObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":                 "application/xml",
			"x-oss-request-id":             "534B371674E88A4D8906****",
			"Date":                         "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                         "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma":         "870718044876840****",
			"x-oss-copy-source-version-id": "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****",
			"x-oss-version-id":             "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "text",
			"ETag":                 "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma": "870718044876841****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			TrafficLimit:    int64(100 * 1024 * 8),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object?versionId=CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****", r.Header.Get("x-oss-copy-source"))
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&CopyObjectRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			SourceKey:       Ptr("copy-object"),
			TrafficLimit:    int64(100 * 1024 * 8),
			SourceVersionId: Ptr("CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****"),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
}

func TestMockCopyObject_Success(t *testing.T) {
	for _, c := range testMockCopyObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
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
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
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
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			SourceKey: Ptr("copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CopyObject fail")
		},
	},
}

func TestMockCopyObject_Error(t *testing.T) {
	for _, c := range testMockCopyObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":           "6551DBCF4311A7303980****",
			"Date":                       "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":           "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position": "0",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			Body:     strings.NewReader("hi oss,append object,this is a demo"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(0))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object"),
			Position:                  Ptr(int64(100)),
			Body:                      strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
		},
		&AppendObjectRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object"),
			Position:                  Ptr(int64(100)),
			Body:                      strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
			TrafficLimit:              int64(100 * 1024 * 8),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-traffic-limit"), strconv.FormatInt(100*1024*8, 10))
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&AppendObjectRequest{
			Bucket:                    Ptr("bucket"),
			Key:                       Ptr("object"),
			Position:                  Ptr(int64(100)),
			Body:                      strings.NewReader("hi oss,append object,this is a demo"),
			ServerSideEncryption:      Ptr("KMS"),
			ServerSideDataEncryption:  Ptr("SM4"),
			ServerSideEncryptionKeyId: Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
			TrafficLimit:              int64(100 * 1024 * 8),
			RequestPayer:              Ptr("requester"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
}

func TestMockAppendObject_Success(t *testing.T) {
	for _, c := range testMockAppendObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectDisableDetectMimeTypeCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockAppendObject_DisableDetectMimeType(t *testing.T) {
	for _, c := range testMockAppendObjectDisableDetectMimeTypeCases {
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

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectWithProgressCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
}

func TestMockAppendObject_Progress(t *testing.T) {
	for _, c := range testMockAppendObjectWithProgressCases {
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
		}
		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, int64(len("hi oss,append object")), n)
	}
}

var testMockAppendObjectCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "20",
			"x-oss-hash-crc64ecma":       "2313496259928504459",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(0)),
			Body:          strings.NewReader("hi oss,append object"),
			InitHashCRC64: Ptr("0"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(20))
			assert.Equal(t, *o.HashCRC64, "2313496259928504459")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":           "6551DBCF4311A7303980****",
			"Date":                       "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-next-append-position": "35",
			"x-oss-hash-crc64ecma":       "8586970469916596321",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader(",this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=20", strUrl)
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(20)),
			Body:          strings.NewReader(",this is a demo"),
			InitHashCRC64: Ptr("2313496259928504459"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.HashCRC64, "8586970469916596321")
			assert.Equal(t, o.NextPosition, int64(35))
		},
	},
}

func TestMockAppendObject_CRC64(t *testing.T) {
	for _, c := range testMockAppendObjectCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectDisableCRC64Cases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "20",
			"x-oss-hash-crc64ecma":       "4313496259928504459",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
			assert.Equal(t, "application/octet-stream", r.Header.Get(HTTPHeaderContentType))
		},
		&AppendObjectRequest{
			Bucket:        Ptr("bucket"),
			Key:           Ptr("object"),
			Position:      Ptr(int64(0)),
			Body:          strings.NewReader("hi oss,append object"),
			InitHashCRC64: Ptr("0"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(20))
			assert.Equal(t, *o.HashCRC64, "4313496259928504459")
		},
	},
}

func TestMockAppendObject_DisableCRC64(t *testing.T) {
	for _, c := range testMockAppendObjectDisableCRC64Cases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		//Enable, meets error
		client := NewClient(cfg)
		assert.NotNil(t, c)

		_, err := client.AppendObject(context.TODO(), c.Request)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "crc is inconsistent, client 2313496259928504459, server 4313496259928504459")

		// Disable, no error
		client = NewClient(cfg,
			func(o *Options) {
				o.FeatureFlags = o.FeatureFlags & ^FeatureEnableCRC64CheckUpload
			})
		assert.NotNil(t, c)
		c.Request.Body = strings.NewReader("hi oss,append object")
		output, err := client.AppendObject(context.TODO(), c.Request)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output, err)

		// don't set initCRC, no error
		client = NewClient(cfg)
		assert.NotNil(t, c)
		c.Request.InitHashCRC64 = nil
		c.Request.Body = strings.NewReader("hi oss,append object")
		output, err = client.AppendObject(context.TODO(), c.Request)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			strUrl := sortQuery(r)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			Body:     strings.NewReader("hi oss,append object"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
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
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>PositionNotEqualToLength</Code>
  <Message>Position is not equal to file length</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0026-00000016</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			Body:     strings.NewReader("hi oss,append object,this is a demo"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "PositionNotEqualToLength", serr.Code)
			assert.Equal(t, "Position is not equal to file length", serr.Message)
			assert.Equal(t, "0026-00000016", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockAppendObject_Error(t *testing.T) {
	for _, c := range testMockAppendObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
}{
	{
		204,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionId)
			assert.False(t, o.DeleteMarker)
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id":    "6551DBCF4311A7303980****",
			"Date":                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-delete-marker": "true",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.True(t, o.DeleteMarker)
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id":    "6551DBCF4311A7303980****",
			"Date":                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-delete-marker": "true",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.True(t, o.DeleteMarker)
		},
	},
}

func TestMockDeleteObject_Success(t *testing.T) {
	for _, c := range testMockDeleteObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
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
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
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

func TestMockDeleteObject_Error(t *testing.T) {
	for _, c := range testMockDeleteObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
			Quiet:   true,
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Nil(t, o.DeletedObjects)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		}, []byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), ("<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>"))
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 2)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key></Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.DeletedObjects[0].Key, "")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			EncodingType: Ptr("url"),
			Delete: &Delete{
				Objects: []ObjectIdentifier{
					{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")},
					{Key: Ptr("key2.txt"), VersionId: Ptr("CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")},
				},
			},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 2)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			EncodingType: Ptr("url"),
			Delete: &Delete{
				Quiet: true,
				Objects: []ObjectIdentifier{
					{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")},
					{Key: Ptr("key2.txt"), VersionId: Ptr("CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")},
				},
			},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 2)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
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
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			EncodingType: Ptr("url"),
			Delete: &Delete{
				Quiet: true,
				Objects: []ObjectIdentifier{
					{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")},
				},
			},
			Objects: []DeleteObject{
				{Key: Ptr("key2.txt"), VersionId: Ptr("CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")},
			},
			Quiet: false,
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
		},
	},
}

func TestMockDeleteMultipleObjects_Success(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
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
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ErrorDetail>the root node is not named Delete.</ErrorDetail>
  <EC>0016-00000608</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000608</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "0016-00000608", serr.EC)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
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
			assert.Equal(t, "/bucket/?delete&encoding-type=url", strUrl)
			data, err := io.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute DeleteMultipleObjects fail")
		},
	},
}

func TestMockDeleteMultipleObjects_Error(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "6555A936CA31DC333143****",
			"Date":                "Thu, 16 Nov 2023 05:31:34 GMT",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":      "344606",
			"Content-Type":        "image/jpg",
			"ETag":                "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ContentType, "image/jpg")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "5CAC3B40B7AEADE01700****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "text/xml",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id":    "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "text/xml")
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "534B371674E88A4D8906****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "image/jpg",
			"x-oss-object-type":   "Normal",
			"x-oss-restore":       "ongoing-request=\"true\"",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"true\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":                    "534B371674E88A4D8906****",
			"Date":                                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":                        "image/jpg",
			"x-oss-object-type":                   "Normal",
			"x-oss-restore":                       "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"",
			"x-oss-storage-class":                 "Archive",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "9468da86-3509-4f8d-a61e-6eab1eac****",
			"Content-Length":                      "481827",
			"ETag":                                "\"A082B659EF78733A5A042FA253B1****\"",
			"Last-Modified":                       "Fri, 24 Feb 2023 09:41:59 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":                    "534B371674E88A4D8906****",
			"Date":                                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":                        "image/jpg",
			"x-oss-object-type":                   "Normal",
			"x-oss-restore":                       "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"",
			"x-oss-storage-class":                 "Archive",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "9468da86-3509-4f8d-a61e-6eab1eac****",
			"Content-Length":                      "481827",
			"ETag":                                "\"A082B659EF78733A5A042FA253B1****\"",
			"Last-Modified":                       "Fri, 24 Feb 2023 09:41:59 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&HeadObjectRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryptionKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
		},
	},
}

func TestMockHeadObject_Success(t *testing.T) {
	for _, c := range testMockHeadObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockHeadObject_Error(t *testing.T) {
	for _, c := range testMockHeadObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
			"Last-Modified":    "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":   "344606",
			"ETag":             "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(344606))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":    "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":   "481827",
			"ETag":             "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?objectMeta&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
		},
		&GetObjectMetaRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":       "534B371674E88A4D8906****",
			"Date":                   "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":          "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":         "481827",
			"ETag":                   "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-last-access-time": "Thu, 14 Oct 2021 11:49:05 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":       "534B371674E88A4D8906****",
			"Date":                   "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":          "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":         "481827",
			"ETag":                   "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-last-access-time": "Thu, 14 Oct 2021 11:49:05 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&GetObjectMetaRequest{
			Bucket:       Ptr("bucket"),
			Key:          Ptr("object"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
		},
	},
}

func TestMockGetObjectMeta_Success(t *testing.T) {
	for _, c := range testMockGetObjectMetaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockGetObjectMeta_Error(t *testing.T) {
	for _, c := range testMockGetObjectMetaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


