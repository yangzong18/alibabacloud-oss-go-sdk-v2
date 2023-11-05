package oss

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func testSetupMockServer(t *testing.T, statusCode int, headers map[string]string, body []byte,
	chkfunc func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request
		chkfunc(t, r)

		// header s
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		// status code
		w.WriteHeader(statusCode)

		// body
		w.Write(body)
	}))
}

var testInvokeOperationAnonymousCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<ListAllMyBucketsResult>
			<Owner>
				<ID>512**</ID>
				<DisplayName>51264</DisplayName>
			</Owner>
			<Buckets>
				<Bucket>
				<CreationDate>2014-02-17T18:12:43.000Z</CreationDate>
				<ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
				<IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
				<Location>oss-cn-shanghai</Location>
				<Name>app-base-oss</Name>
				<Region>cn-shanghai</Region>
				<StorageClass>Standard</StorageClass>
				</Bucket>
			</Buckets>				
			</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&OperationInput{
			OpName: "ListBuckets",
			Method: "GET",
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetBucketAcl",
			Bucket: Ptr("bucket"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "654605AA6172673135811AB3", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Sat, 04 Nov 2023 08:49:46 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/key?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetObjectAcl",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<InitiateMultipartUploadResult>
				<Bucket>oss-example</Bucket>
				<Key>key+ 123.data</Key>
				<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
			</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/key%2B%20123/test.data?uploads", r.URL.String())
			assert.Equal(t, "POST", r.Method)
		},
		&OperationInput{
			OpName: "InitiateMultipartUpload",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key+ 123/test.data"),
			Method: "POST",
			Parameters: map[string]string{
				"uploads": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "text/txt",
		},
		[]byte(
			`hello world`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket//subfolder/example.txt?versionId=CAEQNhiBgMDJgZCA0BY%2B123", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Etag1234", r.Header.Get("If-Match"))
		},
		&OperationInput{
			OpName: "GetObject",
			Bucket: Ptr("bucket"),
			Key:    Ptr("/subfolder/example.txt"),
			Method: "GET",
			Headers: map[string]string{
				"If-Match": "Etag1234",
			},
			Parameters: map[string]string{
				"versionId": "CAEQNhiBgMDJgZCA0BY+123",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
}

func TestInvokeOperation_Anonymous(t *testing.T) {
	for _, c := range testInvokeOperationAnonymousCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL).
			WithUsePathStyle(true)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)
	}
}

var testInvokeOperationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput, err error)
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
			assert.Equal(t, "/bucket/test-key.txt", r.URL.String())
		},
		&OperationInput{
			OpName: "PutObject",
			Method: "PUT",
			Bucket: Ptr("bucket"),
			Key:    Ptr("test-key.txt"),
		},
		func(t *testing.T, o *OperationOutput, err error) {
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
			assert.Contains(t, serr.RequestTarget, "/bucket/test-key.txt")
		},
	},
}

func TestInvokeOperation_Error(t *testing.T) {
	for _, c := range testInvokeOperationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL).
			WithUsePathStyle(true)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}
