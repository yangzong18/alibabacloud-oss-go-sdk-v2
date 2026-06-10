package oss

import (
	"testing"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/retry"
	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, "/bucket/?acl", r.URL.String())
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
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)

		var fns []func(*Options)
		fns = append(fns, func(c *Options) { c.OpReadWriteTimeout = Ptr(1 * time.Second) })
		output, err = client.InvokeOperation(context.TODO(), c.Input, fns...)
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
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}

func TestInvokeOperation_RetryMaxAttempts(t *testing.T) {
	for _, c := range testInvokeOperationRetryMaxAttempts {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		var (
			retryAttemptsClient = 2
			retryAttemptsOp     = 4
		)
		assert.NotEqual(t, retryAttemptsClient, retry.DefaultMaxAttempts)
		assert.NotEqual(t, retryAttemptsOp, retry.DefaultMaxAttempts)

		//default
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		testRetryMaxAttemptsRevc = 0
		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, retry.DefaultMaxAttempts, testRetryMaxAttemptsRevc)

		// overwrite througth client options
		cfg = LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL).
			WithRetryMaxAttempts(retryAttemptsClient)

		client = NewClient(cfg)
		assert.NotNil(t, c)

		testRetryMaxAttemptsRevc = 0
		output, err = client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, retryAttemptsClient, testRetryMaxAttemptsRevc)

		// overwrite througth InvokeOperation options
		testRetryMaxAttemptsRevc = 0
		output, err = client.InvokeOperation(context.TODO(), c.Input, func(o *Options) {
			o.RetryMaxAttempts = Ptr(retryAttemptsOp)
		})
		c.CheckOutputFn(t, output, err)
		assert.Equal(t, retryAttemptsOp, testRetryMaxAttemptsRevc)
	}
}

var testMockUserAgentCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResult, err error)
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
			assert.Equal(t, "/bucket/", r.URL.String())
			assert.Contains(t, r.Header.Get("User-Agent"), "/my-agent")
			assert.True(t, strings.HasSuffix(r.Header.Get("User-Agent"), "/my-agent"))
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
		},
	},
}

func TestMockUserAgentCases(t *testing.T) {
	for _, c := range testMockUserAgentCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL).
			WithUserAgent("my-agent")

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


var testMockGetBucketDataRedundancyTransitionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>751f5243f8ac4ae89f34726534d1****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=751f5243f8ac4ae89f34726534d1%2A%2A%2A%2A", strUrl)
		},
		&GetBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("751f5243f8ac4ae89f34726534d1****"),
		},
		func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")

			assert.Equal(t, *o.BucketDataRedundancyTransition.Bucket, "examplebucket")
			assert.Equal(t, *o.BucketDataRedundancyTransition.TaskId, "751f5243f8ac4ae89f34726534d1****")
			assert.Equal(t, *o.BucketDataRedundancyTransition.Status, "Queueing")
			assert.Equal(t, *o.BucketDataRedundancyTransition.CreateTime, "2023-11-17T09:11:58.000Z")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=909c6c818dd041d1a44e0fdc66aa%2A%2A%2A%2A", strUrl)
		},
		&GetBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("909c6c818dd041d1a44e0fdc66aa****"),
		},
		func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")

			assert.Equal(t, *o.BucketDataRedundancyTransition.Bucket, "examplebucket")
			assert.Equal(t, *o.BucketDataRedundancyTransition.TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.BucketDataRedundancyTransition.Status, "Finished")
			assert.Equal(t, *o.BucketDataRedundancyTransition.CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.BucketDataRedundancyTransition.StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.BucketDataRedundancyTransition.ProcessPercentage, int32(100))
			assert.Equal(t, *o.BucketDataRedundancyTransition.EstimatedRemainingTime, int64(0))
			assert.Equal(t, *o.BucketDataRedundancyTransition.EndTime, "2023-11-18T09:14:39.000Z")
		},
	},
}

func TestMockGetBucketDataRedundancyTransition_Success(t *testing.T) {
	for _, c := range testMockGetBucketDataRedundancyTransitionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketDataRedundancyTransitionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error)
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
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=909c6c818dd041d1a44e0fdc66aa%2A%2A%2A%2A", strUrl)
		},
		&GetBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("909c6c818dd041d1a44e0fdc66aa****"),
		},
		func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error) {
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
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=909c6c818dd041d1a44e0fdc66aa%2A%2A%2A%2A", strUrl)
		},
		&GetBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("909c6c818dd041d1a44e0fdc66aa****"),
		},
		func(t *testing.T, o *GetBucketDataRedundancyTransitionResult, err error) {
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

func TestMockGetBucketDataRedundancyTransition_Error(t *testing.T) {
	for _, c := range testMockGetBucketDataRedundancyTransitionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateBucketDataRedundancyTransitionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-target-redundancy-type=ZRS", urlStr)
		},
		&CreateBucketDataRedundancyTransitionRequest{
			Bucket:               Ptr("bucket"),
			TargetRedundancyType: Ptr("ZRS"),
		},
		func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockCreateBucketDataRedundancyTransition_Success(t *testing.T) {
	for _, c := range testMockCreateBucketDataRedundancyTransitionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateBucketDataRedundancyTransitionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<Error>
  <Code>TooManyCnameToken</Code>
  <Message>You have attempted to create more cname token than allowed.</Message>
  <RequestId>6215FD21DA0E27393F004E9E</RequestId>
  <HostId>127.0.0.1</HostId>
  <Bucket>bucket</Bucket>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-target-redundancy-type=ZRS", urlStr)
		},
		&CreateBucketDataRedundancyTransitionRequest{
			Bucket:               Ptr("bucket"),
			TargetRedundancyType: Ptr("ZRS"),
		},
		func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "TooManyCnameToken", serr.Code)
			assert.Equal(t, "You have attempted to create more cname token than allowed.", serr.Message)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-target-redundancy-type=ZRS", urlStr)
		},
		&CreateBucketDataRedundancyTransitionRequest{
			Bucket:               Ptr("bucket"),
			TargetRedundancyType: Ptr("ZRS"),
		},
		func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-target-redundancy-type=ZRS", urlStr)
		},
		&CreateBucketDataRedundancyTransitionRequest{
			Bucket:               Ptr("bucket"),
			TargetRedundancyType: Ptr("ZRS"),
		},
		func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error) {
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
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-target-redundancy-type=ZRS", urlStr)
		},
		&CreateBucketDataRedundancyTransitionRequest{
			Bucket:               Ptr("bucket"),
			TargetRedundancyType: Ptr("ZRS"),
		},
		func(t *testing.T, o *CreateBucketDataRedundancyTransitionResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute CreateBucketDataRedundancyTransition fail")
		},
	},
}

func TestMockCreateBucketDataRedundancyTransition_Error(t *testing.T) {
	for _, c := range testMockCreateBucketDataRedundancyTransitionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketDataRedundancyTransitionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketDataRedundancyTransitionResult, err error)
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
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=123", strUrl)
		},
		&DeleteBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("123"),
		},
		func(t *testing.T, o *DeleteBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucketDataRedundancyTransition_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketDataRedundancyTransitionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketDataRedundancyTransitionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketDataRedundancyTransitionResult, err error)
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
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=123", strUrl)
		},
		&DeleteBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("123"),
		},
		func(t *testing.T, o *DeleteBucketDataRedundancyTransitionResult, err error) {
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
			assert.Equal(t, "/bucket/?redundancyTransition&x-oss-redundancy-transition-taskid=123", strUrl)
		},
		&DeleteBucketDataRedundancyTransitionRequest{
			Bucket:                     Ptr("bucket"),
			RedundancyTransitionTaskid: Ptr("123"),
		},
		func(t *testing.T, o *DeleteBucketDataRedundancyTransitionResult, err error) {
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

func TestMockDeleteBucketDataRedundancyTransition_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketDataRedundancyTransitionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketDataRedundancyTransitionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition", strUrl)
		},
		&ListBucketDataRedundancyTransitionRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "4be5beb0f74f490186311b268bf6****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Queueing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:11:58.000Z")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition", strUrl)
		},
		&ListBucketDataRedundancyTransitionRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Processing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].ProcessPercentage, int32(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EstimatedRemainingTime, int64(100))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>
</ListBucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?redundancyTransition", strUrl)
		},
		&ListBucketDataRedundancyTransitionRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Finished")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].ProcessPercentage, int32(100))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EstimatedRemainingTime, int64(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].EndTime, "2023-11-18T09:14:39.000Z")
		},
	},
}

func TestMockListBucketDataRedundancyTransition_Success(t *testing.T) {
	for _, c := range testMockListBucketDataRedundancyTransitionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketDataRedundancyTransitionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error)
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
			assert.Equal(t, "/bucket/?redundancyTransition", strUrl)
		},
		&ListBucketDataRedundancyTransitionRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error) {
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
			assert.Equal(t, "/bucket/?redundancyTransition", strUrl)
		},
		&ListBucketDataRedundancyTransitionRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketDataRedundancyTransitionResult, err error) {
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

func TestMockListBucketDataRedundancyTransition_Error(t *testing.T) {
	for _, c := range testMockListBucketDataRedundancyTransitionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListBucketDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListUserDataRedundancyTransitionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListUserDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket1</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket2</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>
<IsTruncated>false</IsTruncated>
<NextContinuationToken></NextContinuationToken>
</ListBucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?redundancyTransition", strUrl)
		},
		&ListUserDataRedundancyTransitionRequest{},
		func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, "", *o.ListBucketDataRedundancyTransition.NextContinuationToken)
			assert.False(t, *o.ListBucketDataRedundancyTransition.IsTruncated)
			assert.Equal(t, 3, len(o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "4be5beb0f74f490186311b268bf6****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Queueing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:11:58.000Z")

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Bucket, "examplebucket1")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Status, "Processing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].ProcessPercentage, int32(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].EstimatedRemainingTime, int64(100))

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Bucket, "examplebucket2")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Status, "Finished")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].ProcessPercentage, int32(100))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EstimatedRemainingTime, int64(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EndTime, "2023-11-18T09:14:39.000Z")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket</Bucket>
  <TaskId>4be5beb0f74f490186311b268bf6****</TaskId>
  <Status>Queueing</Status>
  <CreateTime>2023-11-17T09:11:58.000Z</CreateTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket1</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Processing</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>0</ProcessPercentage>
  <EstimatedRemainingTime>100</EstimatedRemainingTime>
</BucketDataRedundancyTransition>
<BucketDataRedundancyTransition>
  <Bucket>examplebucket2</Bucket>
  <TaskId>909c6c818dd041d1a44e0fdc66aa****</TaskId>
  <Status>Finished</Status>
  <CreateTime>2023-11-17T09:14:39.000Z</CreateTime>
  <StartTime>2023-11-17T09:14:39.000Z</StartTime>
  <ProcessPercentage>100</ProcessPercentage>
  <EstimatedRemainingTime>0</EstimatedRemainingTime>
  <EndTime>2023-11-18T09:14:39.000Z</EndTime>
</BucketDataRedundancyTransition>
<IsTruncated>false</IsTruncated>
<NextContinuationToken></NextContinuationToken>
</ListBucketDataRedundancyTransition>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?continuation-token=123&max-keys=3&redundancyTransition", strUrl)
		},
		&ListUserDataRedundancyTransitionRequest{
			ContinuationToken: Ptr("123"),
			MaxKeys:           int32(3),
		},
		func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")

			assert.Equal(t, "", *o.ListBucketDataRedundancyTransition.NextContinuationToken)
			assert.False(t, *o.ListBucketDataRedundancyTransition.IsTruncated)
			assert.Equal(t, 3, len(o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Bucket, "examplebucket")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].TaskId, "4be5beb0f74f490186311b268bf6****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].Status, "Queueing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[0].CreateTime, "2023-11-17T09:11:58.000Z")

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Bucket, "examplebucket1")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].Status, "Processing")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].ProcessPercentage, int32(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[1].EstimatedRemainingTime, int64(100))

			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Bucket, "examplebucket2")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].TaskId, "909c6c818dd041d1a44e0fdc66aa****")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].Status, "Finished")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].CreateTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].StartTime, "2023-11-17T09:14:39.000Z")
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].ProcessPercentage, int32(100))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EstimatedRemainingTime, int64(0))
			assert.Equal(t, *o.ListBucketDataRedundancyTransition.BucketDataRedundancyTransitions[2].EndTime, "2023-11-18T09:14:39.000Z")
		},
	},
}

func TestMockListUserDataRedundancyTransition_Success(t *testing.T) {
	for _, c := range testMockListUserDataRedundancyTransitionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListUserDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListUserDataRedundancyTransitionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListUserDataRedundancyTransitionRequest
	CheckOutputFn  func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error)
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
			assert.Equal(t, "/?redundancyTransition", strUrl)
		},
		&ListUserDataRedundancyTransitionRequest{},
		func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error) {
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
			assert.Equal(t, "/?redundancyTransition", strUrl)
		},
		&ListUserDataRedundancyTransitionRequest{},
		func(t *testing.T, o *ListUserDataRedundancyTransitionResult, err error) {
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

func TestMockListUserDataRedundancyTransition_Error(t *testing.T) {
	for _, c := range testMockListUserDataRedundancyTransitionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListUserDataRedundancyTransition(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDescribeRegionsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DescribeRegionsRequest
	CheckOutputFn  func(t *testing.T, o *DescribeRegionsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RegionInfoList>
  <RegionInfo>
     <Region>oss-cn-hangzhou</Region>
     <InternetEndpoint>oss-cn-hangzhou.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
  <RegionInfo>
     <Region>oss-cn-shanghai</Region>
     <InternetEndpoint>oss-cn-shanghai.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-shanghai-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
</RegionInfoList>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?regions", strUrl)
		},
		&DescribeRegionsRequest{},
		func(t *testing.T, o *DescribeRegionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, len(o.RegionInfoList.RegionInfos), 2)
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].InternetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].InternalEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].Region, "oss-cn-hangzhou")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")

			assert.Equal(t, *o.RegionInfoList.RegionInfos[1].InternetEndpoint, "oss-cn-shanghai.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[1].InternalEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[1].Region, "oss-cn-shanghai")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[1].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<RegionInfoList>
  <RegionInfo>
     <Region>oss-cn-hangzhou</Region>
     <InternetEndpoint>oss-cn-hangzhou.aliyuncs.com</InternetEndpoint>
     <InternalEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</InternalEndpoint>
     <AccelerateEndpoint>oss-accelerate.aliyuncs.com</AccelerateEndpoint>  
  </RegionInfo>
</RegionInfoList>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?regions=oss-cn-hangzhou", strUrl)
		},
		&DescribeRegionsRequest{
			Regions: Ptr("oss-cn-hangzhou"),
		},
		func(t *testing.T, o *DescribeRegionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, len(o.RegionInfoList.RegionInfos), 1)
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].InternetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].InternalEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].Region, "oss-cn-hangzhou")
			assert.Equal(t, *o.RegionInfoList.RegionInfos[0].AccelerateEndpoint, "oss-accelerate.aliyuncs.com")
		},
	},
}

func TestMockDescribeRegions_Success(t *testing.T) {
	for _, c := range testMockDescribeRegionsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DescribeRegions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDescribeRegionsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DescribeRegionsRequest
	CheckOutputFn  func(t *testing.T, o *DescribeRegionsResult, err error)
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
			assert.Equal(t, "/?regions", strUrl)
		},
		&DescribeRegionsRequest{},
		func(t *testing.T, o *DescribeRegionsResult, err error) {
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
			assert.Equal(t, "/?regions", strUrl)
		},
		&DescribeRegionsRequest{},
		func(t *testing.T, o *DescribeRegionsResult, err error) {
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

func TestMockDescribeRegions_Error(t *testing.T) {
	for _, c := range testMockDescribeRegionsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DescribeRegions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListCloudBoxesSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListCloudBoxesRequest
	CheckOutputFn  func(t *testing.T, o *ListCloudBoxesResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListCloudBoxResult>
  <Owner>
     <ID>51264</ID>
    <DisplayName>51264</DisplayName>
  </Owner>
  <CloudBoxes>
    <CloudBox>
      <ID>cb-f8z7yvzgwfkl9q0h****</ID>
      <Name>bucket1</Name>
      <Region>cn-shanghai</Region>
      <ControlEndpoint>cb-f8z7yvzgwfkl9q0h****.cn-shanghai.oss-cloudbox-control.aliyuncs.com</ControlEndpoint>
      <DataEndpoint>cb-f8z7yvzgwfkl9q0h****.cn-shanghai.oss-cloudbox.aliyuncs.com</DataEndpoint>
    </CloudBox>
    <CloudBox>
      <ID>cb-f9z7yvzgwfkl9q0h****</ID>
      <Name>bucket2</Name>
      <Region>cn-hangzhou</Region>
      <ControlEndpoint>cb-f9z7yvzgwfkl9q0h****.cn-hangzhou.oss-cloudbox-control.aliyuncs.com</ControlEndpoint>
      <DataEndpoint>cb-f9z7yvzgwfkl9q0h****.cn-hangzhou.oss-cloudbox.aliyuncs.com</DataEndpoint>
    </CloudBox>
  </CloudBoxes>
</ListCloudBoxResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/?cloudboxes", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListCloudBoxesResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Owner.DisplayName, "51264")
			assert.Equal(t, *o.Owner.ID, "51264")
			assert.Equal(t, len(o.CloudBoxes), 2)
			assert.Equal(t, *o.CloudBoxes[0].ID, "cb-f8z7yvzgwfkl9q0h****")
			assert.Equal(t, *o.CloudBoxes[0].ControlEndpoint, "cb-f8z7yvzgwfkl9q0h****.cn-shanghai.oss-cloudbox-control.aliyuncs.com")
			assert.Equal(t, *o.CloudBoxes[0].DataEndpoint, "cb-f8z7yvzgwfkl9q0h****.cn-shanghai.oss-cloudbox.aliyuncs.com")
			assert.Equal(t, *o.CloudBoxes[0].Name, "bucket1")
			assert.Equal(t, *o.CloudBoxes[0].Region, "cn-shanghai")
			assert.Equal(t, *o.CloudBoxes[1].ID, "cb-f9z7yvzgwfkl9q0h****")
			assert.Equal(t, *o.CloudBoxes[1].ControlEndpoint, "cb-f9z7yvzgwfkl9q0h****.cn-hangzhou.oss-cloudbox-control.aliyuncs.com")
			assert.Equal(t, *o.CloudBoxes[1].DataEndpoint, "cb-f9z7yvzgwfkl9q0h****.cn-hangzhou.oss-cloudbox.aliyuncs.com")
			assert.Equal(t, *o.CloudBoxes[1].Name, "bucket2")
			assert.Equal(t, *o.CloudBoxes[1].Region, "cn-hangzhou")
		},
	},
}

func TestMockListCloudBoxes_Success(t *testing.T) {
	for _, c := range testMockListCloudBoxesSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListCloudBoxes(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListCloudBoxesErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListCloudBoxesRequest
	CheckOutputFn  func(t *testing.T, o *ListCloudBoxesResult, err error)
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
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/?cloudboxes", r.URL.String())
		},
		&ListCloudBoxesRequest{},
		func(t *testing.T, o *ListCloudBoxesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "InvalidAccessKeyId", serr.Code)
			assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
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
			assert.Equal(t, "/?cloudboxes", r.URL.String())
		},
		&ListCloudBoxesRequest{},
		func(t *testing.T, o *ListCloudBoxesResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListCloudBoxes fail")
		},
	},
}

func TestMockListCloudBoxes_Error(t *testing.T) {
	for _, c := range testMockListCloudBoxesErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListCloudBoxes(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockReturnsJsonErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLocationResult, err error)
}{
	// Normal case
	{
		404,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
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
	// Bad Xml
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "4C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`
  "Error": {
    "Code": "NoSuchBucket",
    "Message": "The specified bucket does not exist.",
    "RequestId": "5C3D9175B6FC201293AD****",
    "HostId": "test.oss-cn-hangzhou.aliyuncs.com",
    "BucketName": "test",
    "EC": "0015-00000101"
  }
}`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "BadErrorResponse", serr.Code)
			assert.Contains(t, serr.Message, "Failed to parse json from response body due to")
			assert.Equal(t, "", serr.EC)
			assert.Equal(t, "4C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	// Empty Xml
	{
		403,
		map[string]string{
			"Content-Type":     "application/json",
			"x-oss-request-id": "4C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "BadErrorResponse", serr.Code)
			assert.Contains(t, serr.Message, "Failed to parse json from response body due to")
			assert.Equal(t, "", serr.EC)
			assert.Equal(t, "4C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockReturnsXmlError(t *testing.T) {
	for _, c := range testMockReturnsJsonErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


