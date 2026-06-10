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

var testMockRestoreObjectLegacySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		202,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 202, o.StatusCode)
			assert.Equal(t, "202 Accepted", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?restore&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				Tier: Ptr("Standard"),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				Tier: Ptr("Standard"),
			},
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockRestoreObjectLegacy_Success(t *testing.T) {
	for _, c := range testMockRestoreObjectLegacySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		202,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 202, o.StatusCode)
			assert.Equal(t, "202 Accepted", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?restore&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				JobParameters: &JobParameters{
					Tier: Ptr("Standard"),
				},
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Bluk</Tier></JobParameters></RestoreRequest>")
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				JobParameters: &JobParameters{
					Tier: Ptr("Bluk"),
				},
				Tier: Ptr("Standard"),
			},
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockRestoreObject_Success(t *testing.T) {
	for _, c := range testMockRestoreObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6557176CD11E5535303C****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchKey</Code>
		<Message>The specified key does not exist.</Message>
		<RequestId>6557176CD11E5535303C****</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<Key>walker-not-.txt</Key>
		<EC>0026-00000001</EC>
		<RecommendDoc>https://api.aliyun.com/troubleshoot?q=0026-00000001</RecommendDoc>
		</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6557176CD11E5535303C****", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
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
  <Code>OperationNotSupported</Code>
  <Message>The operation is not supported for this resource</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <Detail>RestoreObject operation does not support this object storage class</Detail>
  <EC>0016-00000702</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "OperationNotSupported", serr.Code)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
			assert.Equal(t, "The operation is not supported for this resource", serr.Message)
			assert.Equal(t, "0016-00000702", serr.EC)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E55363686****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0026-00000011",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>RestoreAlreadyInProgress</Code>
  <Message>The restore operation is in progress.</Message>
  <RequestId>6556FF5BD11E55363686****</RequestId>
  <HostId>10.101.XX.XX</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
			assert.Equal(t, "6556FF5BD11E55363686****", serr.RequestID)
			assert.Equal(t, "The restore operation is in progress.", serr.Message)
		},
	},
}

func TestMockRestoreObject_Error(t *testing.T) {
	for _, c := range testMockRestoreObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


var testMockCleanRestoredObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CleanRestoredObjectRequest
	CheckOutputFn  func(t *testing.T, o *CleanRestoredObjectResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/key?cleanRestoredObject", strUrl)
		},
		&CleanRestoredObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		func(t *testing.T, o *CleanRestoredObjectResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/key?cleanRestoredObject&versionId=version-id", strUrl)
		},
		&CleanRestoredObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("key"),
			VersionId: Ptr("version-id"),
		},
		func(t *testing.T, o *CleanRestoredObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockCleanRestoredObject_Success(t *testing.T) {
	for _, c := range testMockCleanRestoredObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CleanRestoredObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCleanRestoredObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CleanRestoredObjectRequest
	CheckOutputFn  func(t *testing.T, o *CleanRestoredObjectResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/key?cleanRestoredObject", strUrl)
		},
		&CleanRestoredObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		func(t *testing.T, o *CleanRestoredObjectResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/key?cleanRestoredObject", strUrl)
		},
		&CleanRestoredObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		func(t *testing.T, o *CleanRestoredObjectResult, err error) {
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
  <Code>ArchiveRestoreNotFinished</Code>
  <Message>The archive file's restore is not finished.</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0016-00000719</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000719</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/key?cleanRestoredObject", strUrl)
		},
		&CleanRestoredObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
		},
		func(t *testing.T, o *CleanRestoredObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "ArchiveRestoreNotFinished", serr.Code)
			assert.Equal(t, "The archive file's restore is not finished.", serr.Message)
			assert.Equal(t, "0016-00000719", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockCleanRestoredObject_Error(t *testing.T) {
	for _, c := range testMockCleanRestoredObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CleanRestoredObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


var testMockSealAppendObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SealAppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *SealAppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":  "534B371674E88A4D8906****",
			"Date":              "Wed, 07 May 2025 23:00:00 GMT",
			"x-oss-sealed-time": "Wed, 07 May 2025 23:00:00 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?position=10&seal", strUrl)
			assert.Equal(t, contentTypeXML, r.Header.Get(HTTPHeaderContentType))
		},
		&SealAppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(10)),
		},
		func(t *testing.T, o *SealAppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Wed, 07 May 2025 23:00:00 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "Wed, 07 May 2025 23:00:00 GMT", *o.SealedTime)
		},
	},
}

func TestMockSealAppendObject_Success(t *testing.T) {
	for _, c := range testMockSealAppendObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SealAppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockSealAppendObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *SealAppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *SealAppendObjectResult, err error)
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
			assert.Equal(t, "/bucket/object?position=10&seal", strUrl)
		},
		&SealAppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(10)),
		},
		func(t *testing.T, o *SealAppendObjectResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?position=10&seal", strUrl)
		},
		&SealAppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(10)),
		},
		func(t *testing.T, o *SealAppendObjectResult, err error) {
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

func TestMockSealAppendObject_Error(t *testing.T) {
	for _, c := range testMockSealAppendObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.SealAppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


