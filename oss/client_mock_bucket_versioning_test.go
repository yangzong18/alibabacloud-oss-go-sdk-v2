package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionSuspended,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockPutBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
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
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketVersioning fail")
		},
	},
}

func TestMockPutBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockPutBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionStatus)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Enabled</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Enabled")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Suspended</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Suspended")
		},
	},
}

func TestMockGetBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockGetBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
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
			assert.Equal(t, "/bucket/?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketVersioning fail")
		},
	},
}

func TestMockGetBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockGetBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 0)
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
<ListVersionsResult>
    <Name>examplebucket-1250000000</Name>
    <Prefix/>
    <KeyMarker/>
    <VersionIdMarker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Version>
        <Key>example-object-1.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-05T12:03:10.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B669CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-2.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-09T12:03:09.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B002CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-3.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-10T12:03:08.000Z</LastModified>
        <ETag>4B3F1A2E053D763E1B002CC607C5AGTRF****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "examplebucket-1250000000")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(1000))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			assert.Equal(t, *o.ObjectVersions[0].Key, "example-object-1.jpg")
			assert.Empty(t, *o.ObjectVersions[1].VersionId)
			assert.True(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "5B3C1A2E053D763E1B002CC607C5A0FE1****")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(20))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "1250000000")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "1250000000")
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
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2Fgp-</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>5</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgIDAj.jV3xgiIGFjMDI5ZTRmNGNiODQ0NjE4MDFhODM0Y2UxNTI3****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2023-11-22T05:15:05.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:11:25.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgICCuNrV3xgiIDI2YzMyYTBhM2U1ZTQwNjI4OWQ4OTllZGJiNGIz****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:07:37.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-keys=5&prefix=demo%2Fgp-&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			Delimiter:       Ptr("/"),
			MaxKeys:         int32(5),
			Prefix:          Ptr("demo/gp-"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(5))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			key, _ := url.QueryUnescape(*o.ObjectVersions[0].Key)
			assert.Equal(t, *o.ObjectVersions[0].Key, key)
			assert.Equal(t, *o.ObjectVersions[1].VersionId, "CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****")
			assert.False(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "\"29B94424BC241D80B0AF488A4E4B86AF-1\"")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(96316))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "150692521021****")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "150692521021****")
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
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 0)
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
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
   <Version>
        <Key></Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-05T12:03:10.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B669CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <DeleteMarker>
        <Key></Key>
        <VersionId>CAEQMxiBgICAof2D0BYiIDJhMGE3N2M1YTI1NDQzOGY5NTkyNTI3MGYyMzJm****</VersionId>
        <IsLatest>false</IsLatest>
        <LastModified>2019-04-09T07:27:28.000Z</LastModified>
        <Owner>
          <ID>1234512528586****</ID>
          <DisplayName>12345125285864390</DisplayName>
        </Owner>
    </DeleteMarker>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
			RequestPayer:    Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 1)
			assert.Len(t, o.ObjectDeleteMarkers, 1)
			assert.Equal(t, *o.ObjectVersions[0].Key, "")
			assert.Equal(t, *o.ObjectDeleteMarkers[0].Key, "")
		},
	},
}

func TestMockListObjectVersions_Success(t *testing.T) {
	for _, c := range testMockListObjectVersionsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
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
			assert.Equal(t, "/bucket/?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
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
			assert.Equal(t, "/bucket/?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
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
			assert.Equal(t, "/bucket/?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListObjectVersions fail")
		},
	},
}

func TestMockListObjectVersions_Error(t *testing.T) {
	for _, c := range testMockListObjectVersionsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


