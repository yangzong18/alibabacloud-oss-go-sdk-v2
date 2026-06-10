package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockCreateAccessPointSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *CreateAccessPointResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointResult>
  <AccessPointArn>acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01</AccessPointArn>
  <Alias>ap-01-45ee7945007a2f0bcb595f63e2215cxxxx-ossalias</Alias>
</CreateAccessPointResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointConfiguration><AccessPointName>ap-01</AccessPointName><NetworkOrigin>internet</NetworkOrigin></CreateAccessPointConfiguration>")
		},
		&CreateAccessPointRequest{
			Bucket: Ptr("bucket"),
			CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
				AccessPointName: Ptr("ap-01"),
				NetworkOrigin:   Ptr("internet"),
			},
		},
		func(t *testing.T, o *CreateAccessPointResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccessPointArn, "acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01")
			assert.Equal(t, *o.Alias, "ap-01-45ee7945007a2f0bcb595f63e2215cxxxx-ossalias")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointResult>
  <AccessPointArn>acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01</AccessPointArn>
  <Alias>ap-01-45ee7945007a2f0bcb595f63e2215cxxxx-ossalias</Alias>
</CreateAccessPointResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointConfiguration><AccessPointName>ap-01</AccessPointName><NetworkOrigin>vpc</NetworkOrigin><VpcConfiguration><VpcId>vpc-t4nlw426y44rd3iq4xxxx</VpcId></VpcConfiguration></CreateAccessPointConfiguration>")
		},
		&CreateAccessPointRequest{
			Bucket: Ptr("bucket"),
			CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
				AccessPointName: Ptr("ap-01"),
				NetworkOrigin:   Ptr("vpc"),
				VpcConfiguration: &AccessPointVpcConfiguration{
					VpcId: Ptr("vpc-t4nlw426y44rd3iq4xxxx"),
				},
			},
		},
		func(t *testing.T, o *CreateAccessPointResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccessPointArn, "acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01")
			assert.Equal(t, *o.Alias, "ap-01-45ee7945007a2f0bcb595f63e2215cxxxx-ossalias")
		},
	},
}

func TestMockCreateAccessPoint_Success(t *testing.T) {
	for _, c := range testMockCreateAccessPointSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateAccessPointErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *CreateAccessPointResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointConfiguration><AccessPointName>ap-01</AccessPointName><NetworkOrigin>internet</NetworkOrigin></CreateAccessPointConfiguration>")
		},
		&CreateAccessPointRequest{
			Bucket: Ptr("bucket"),
			CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
				AccessPointName: Ptr("ap-01"),
				NetworkOrigin:   Ptr("internet"),
			},
		},
		func(t *testing.T, o *CreateAccessPointResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointConfiguration><AccessPointName>ap-01</AccessPointName><NetworkOrigin>internet</NetworkOrigin></CreateAccessPointConfiguration>")
		},
		&CreateAccessPointRequest{
			Bucket: Ptr("bucket"),
			CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
				AccessPointName: Ptr("ap-01"),
				NetworkOrigin:   Ptr("internet"),
			},
		},
		func(t *testing.T, o *CreateAccessPointResult, err error) {
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

func TestMockCreateAccessPoint_Error(t *testing.T) {
	for _, c := range testMockCreateAccessPointErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<GetAccessPointResult>
  <AccessPointName>ap-01</AccessPointName>
  <Bucket>oss-example</Bucket>
  <AccountId>111933544165xxxx</AccountId>
  <NetworkOrigin>vpc</NetworkOrigin>
  <VpcConfiguration>
     <VpcId>vpc-t4nlw426y44rd3iq4xxxx</VpcId>
  </VpcConfiguration>
  <AccessPointArn>arn:acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01</AccessPointArn>
  <CreationDate>1626769503</CreationDate>
  <Alias>ap-01-ossalias</Alias>
  <Status>enable</Status>
  <Endpoints>
    <PublicEndpoint>ap-01.oss-cn-hangzhou.oss-accesspoint.aliyuncs.com</PublicEndpoint>
    <InternalEndpoint>ap-01.oss-cn-hangzhou-internal.oss-accesspoint.aliyuncs.com</InternalEndpoint>
  </Endpoints>
  <PublicAccessBlockConfiguration>
    <BlockPublicAccess>true</BlockPublicAccess>
  </PublicAccessBlockConfiguration>
</GetAccessPointResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccessPointName, "ap-01")
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.AccountId, "111933544165xxxx")
			assert.Equal(t, *o.NetworkOrigin, "vpc")
			assert.Equal(t, *o.VpcConfiguration.VpcId, "vpc-t4nlw426y44rd3iq4xxxx")
			assert.Equal(t, *o.AccessPointArn, "arn:acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01")
			assert.Equal(t, *o.CreationDate, "1626769503")
			assert.Equal(t, *o.Alias, "ap-01-ossalias")
			assert.Equal(t, *o.AccessPointStatus, "enable")
			assert.Equal(t, *o.PublicEndpoint, "ap-01.oss-cn-hangzhou.oss-accesspoint.aliyuncs.com")
			assert.Equal(t, *o.InternalEndpoint, "ap-01.oss-cn-hangzhou-internal.oss-accesspoint.aliyuncs.com")
			assert.True(t, *o.PublicAccessBlockConfiguration.BlockPublicAccess)
		},
	},
}

func TestMockGetAccessPoint_Success(t *testing.T) {
	for _, c := range testMockGetAccessPointSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointResult, err error)
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointResult, err error) {
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

func TestMockGetAccessPoint_Error(t *testing.T) {
	for _, c := range testMockGetAccessPointErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointResult, err error)
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteAccessPoint_Success(t *testing.T) {
	for _, c := range testMockDeleteAccessPointSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointResult, err error)
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointResult, err error) {
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

func TestMockDeleteAccessPoint_Error(t *testing.T) {
	for _, c := range testMockDeleteAccessPointErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPoint(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListAccessPointsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListAccessPointsRequest
	CheckOutputFn  func(t *testing.T, o *ListAccessPointsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListAccessPointsResult>
  <IsTruncated>true</IsTruncated>
  <NextContinuationToken>abc</NextContinuationToken>
  <AccountId>111933544165****</AccountId>
  <MaxKeys>3</MaxKeys>
  <AccessPoints>
    <AccessPoint>
      <Bucket>oss-example</Bucket>
      <AccessPointName>ap-01</AccessPointName>
      <Alias>ap-01-ossalias</Alias>
      <NetworkOrigin>vpc</NetworkOrigin>
      <VpcConfiguration>
        <VpcId>vpc-t4nlw426y44rd3iq4****</VpcId>
      </VpcConfiguration>
      <Status>enable</Status>
    </AccessPoint>
    <AccessPoint>
      <Bucket>oss-example</Bucket>
      <AccessPointName>ap-02</AccessPointName>
      <Alias>ap-02-ossalias</Alias>
      <NetworkOrigin>vpc</NetworkOrigin>
      <VpcConfiguration>
        <VpcId>vpc-t4nlw426y44rd3iq2****</VpcId>
      </VpcConfiguration>
      <Status>enable</Status>
    </AccessPoint>
	<AccessPoint>
      <Bucket>oss-example</Bucket>
      <AccessPointName>ap-03</AccessPointName>
      <Alias>ap-03-ossalias</Alias>
      <NetworkOrigin>internet</NetworkOrigin>
      <VpcConfiguration>
        <VpcId></VpcId>
      </VpcConfiguration>
      <Status>enable</Status>
    </AccessPoint>
  </AccessPoints>
</ListAccessPointsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?accessPoint", strUrl)
		},
		&ListAccessPointsRequest{},
		func(t *testing.T, o *ListAccessPointsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccountId, "111933544165****")
			assert.Equal(t, *o.NextContinuationToken, "abc")
			assert.True(t, *o.IsTruncated)
			assert.Equal(t, *o.MaxKeys, int32(3))
			assert.Equal(t, len(o.AccessPoints), 3)
			assert.Equal(t, *o.AccessPoints[0].Bucket, "oss-example")
			assert.Equal(t, *o.AccessPoints[0].AccessPointName, "ap-01")
			assert.Equal(t, *o.AccessPoints[0].Alias, "ap-01-ossalias")
			assert.Equal(t, *o.AccessPoints[0].NetworkOrigin, "vpc")
			assert.Equal(t, *o.AccessPoints[0].VpcConfiguration.VpcId, "vpc-t4nlw426y44rd3iq4****")
			assert.Equal(t, *o.AccessPoints[0].Status, "enable")
			assert.Equal(t, *o.AccessPoints[1].Bucket, "oss-example")
			assert.Equal(t, *o.AccessPoints[1].AccessPointName, "ap-02")
			assert.Equal(t, *o.AccessPoints[1].Alias, "ap-02-ossalias")
			assert.Equal(t, *o.AccessPoints[1].NetworkOrigin, "vpc")
			assert.Equal(t, *o.AccessPoints[1].VpcConfiguration.VpcId, "vpc-t4nlw426y44rd3iq2****")
			assert.Equal(t, *o.AccessPoints[1].Status, "enable")
			assert.Equal(t, *o.AccessPoints[2].Bucket, "oss-example")
			assert.Equal(t, *o.AccessPoints[2].AccessPointName, "ap-03")
			assert.Equal(t, *o.AccessPoints[2].Alias, "ap-03-ossalias")
			assert.Equal(t, *o.AccessPoints[2].NetworkOrigin, "internet")
			assert.Equal(t, *o.AccessPoints[2].VpcConfiguration.VpcId, "")
			assert.Equal(t, *o.AccessPoints[2].Status, "enable")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListAccessPointsResult>
  <IsTruncated>true</IsTruncated>
  <NextContinuationToken>abc</NextContinuationToken>
  <AccountId>111933544165****</AccountId>
  <MaxKeys>2</MaxKeys>
  <AccessPoints>
    <AccessPoint>
      <Bucket>bucket</Bucket>
      <AccessPointName>ap-01</AccessPointName>
      <Alias>ap-01-ossalias</Alias>
      <NetworkOrigin>vpc</NetworkOrigin>
      <VpcConfiguration>
        <VpcId>vpc-t4nlw426y44rd3iq4****</VpcId>
      </VpcConfiguration>
      <Status>enable</Status>
    </AccessPoint>
    <AccessPoint>
      <Bucket>bucket</Bucket>
      <AccessPointName>ap-02</AccessPointName>
      <Alias>ap-02-ossalias</Alias>
      <NetworkOrigin>vpc</NetworkOrigin>
      <VpcConfiguration>
        <VpcId>vpc-t4nlw426y44rd3iq2****</VpcId>
      </VpcConfiguration>
      <Status>enable</Status>
    </AccessPoint>
  </AccessPoints>
</ListAccessPointsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPoint&max-keys=2", strUrl)
		},
		&ListAccessPointsRequest{
			Bucket:  Ptr("bucket"),
			MaxKeys: int64(2),
		},
		func(t *testing.T, o *ListAccessPointsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccountId, "111933544165****")
			assert.Equal(t, *o.NextContinuationToken, "abc")
			assert.True(t, *o.IsTruncated)
			assert.Equal(t, *o.MaxKeys, int32(2))
			assert.Equal(t, len(o.AccessPoints), 2)
			assert.Equal(t, *o.AccessPoints[0].Bucket, "bucket")
			assert.Equal(t, *o.AccessPoints[0].AccessPointName, "ap-01")
			assert.Equal(t, *o.AccessPoints[0].Alias, "ap-01-ossalias")
			assert.Equal(t, *o.AccessPoints[0].NetworkOrigin, "vpc")
			assert.Equal(t, *o.AccessPoints[0].VpcConfiguration.VpcId, "vpc-t4nlw426y44rd3iq4****")
			assert.Equal(t, *o.AccessPoints[0].Status, "enable")
			assert.Equal(t, *o.AccessPoints[1].Bucket, "bucket")
			assert.Equal(t, *o.AccessPoints[1].AccessPointName, "ap-02")
			assert.Equal(t, *o.AccessPoints[1].Alias, "ap-02-ossalias")
			assert.Equal(t, *o.AccessPoints[1].NetworkOrigin, "vpc")
			assert.Equal(t, *o.AccessPoints[1].VpcConfiguration.VpcId, "vpc-t4nlw426y44rd3iq2****")
			assert.Equal(t, *o.AccessPoints[1].Status, "enable")
		},
	},
}

func TestMockListAccessPoints_Success(t *testing.T) {
	for _, c := range testMockListAccessPointsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListAccessPoints(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListAccessPointsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListAccessPointsRequest
	CheckOutputFn  func(t *testing.T, o *ListAccessPointsResult, err error)
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
		},
		&ListAccessPointsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListAccessPointsResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPoint", strUrl)
		},
		&ListAccessPointsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListAccessPointsResult, err error) {
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

func TestMockListAccessPoints_Error(t *testing.T) {
	for _, c := range testMockListAccessPointsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListAccessPoints(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPolicyResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`)
		},
		&PutAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
			Body:            strings.NewReader(`{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutAccessPointPolicy_Success(t *testing.T) {
	for _, c := range testMockPutAccessPointPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPolicyResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`)
		},
		&PutAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
			Body:            strings.NewReader(`{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`)
		},
		&PutAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
			Body:            strings.NewReader(`{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyResult, err error) {
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

func TestMockPutAccessPointPolicy_Error(t *testing.T) {
	for _, c := range testMockPutAccessPointPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPolicyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{"Version":"1","Statement":[{"Action":["oss:GetObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, string(o.Body), `{"Version":"1","Statement":[{"Action":["oss:GetObject","oss:GetObject"],"Effect":"Deny","Principal":["27737962156157xxxx"],"Resource":["acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01","acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/ap-01/object/*"]}]}`)
		},
	},
}

func TestMockGetAccessPointPolicy_Success(t *testing.T) {
	for _, c := range testMockGetAccessPointPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPolicyResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&GetAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyResult, err error) {
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

func TestMockGetAccessPointPolicy_Error(t *testing.T) {
	for _, c := range testMockGetAccessPointPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPolicySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPolicyResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteAccessPointPolicy_Success(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPolicySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPolicyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPolicyRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPolicyResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointPolicy", strUrl)
			assert.Equal(t, "ap-01", r.Header.Get("x-oss-access-point-name"))
		},
		&DeleteAccessPointPolicyRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyResult, err error) {
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

func TestMockDeleteAccessPointPolicy_Error(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPolicyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPolicy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPublicAccessBlockSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error)
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<PublicAccessBlockConfiguration><BlockPublicAccess>true</BlockPublicAccess></PublicAccessBlockConfiguration>")
		},
		&PutAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
		func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error) {
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<PublicAccessBlockConfiguration><BlockPublicAccess>false</BlockPublicAccess></PublicAccessBlockConfiguration>")
		},
		&PutAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(false),
			},
		},
		func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutAccessPointPublicAccessBlock_Success(t *testing.T) {
	for _, c := range testMockPutAccessPointPublicAccessBlockSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPublicAccessBlockErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error)
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<PublicAccessBlockConfiguration><BlockPublicAccess>true</BlockPublicAccess></PublicAccessBlockConfiguration>")
		},
		&PutAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
		func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error) {
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<PublicAccessBlockConfiguration><BlockPublicAccess>true</BlockPublicAccess></PublicAccessBlockConfiguration>")
		},
		&PutAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
		func(t *testing.T, o *PutAccessPointPublicAccessBlockResult, err error) {
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

func TestMockPutAccessPointPublicAccessBlock_Error(t *testing.T) {
	for _, c := range testMockPutAccessPointPublicAccessBlockErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPublicAccessBlockSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<PublicAccessBlockConfiguration>
  <BlockPublicAccess>true</BlockPublicAccess>
</PublicAccessBlockConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
		},
		&GetAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.True(t, *o.PublicAccessBlockConfiguration.BlockPublicAccess)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<PublicAccessBlockConfiguration>
  <BlockPublicAccess>false</BlockPublicAccess>
</PublicAccessBlockConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
		},
		&GetAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.False(t, *o.PublicAccessBlockConfiguration.BlockPublicAccess)
		},
	},
}

func TestMockGetAccessPointPublicAccessBlock_Success(t *testing.T) {
	for _, c := range testMockGetAccessPointPublicAccessBlockSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPublicAccessBlockErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error)
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
		},
		&GetAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error) {
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
		},
		&GetAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error) {
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", urlStr)
		},
		&GetAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *GetAccessPointPublicAccessBlockResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetAccessPointPublicAccessBlock fail")
		},
	},
}

func TestMockGetAccessPointPublicAccessBlock_Error(t *testing.T) {
	for _, c := range testMockGetAccessPointPublicAccessBlockErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPublicAccessBlockSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPublicAccessBlockResult, err error)
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", strUrl)
		},
		&DeleteAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *DeleteAccessPointPublicAccessBlockResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteAccessPointPublicAccessBlock_Success(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPublicAccessBlockSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPublicAccessBlockErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPublicAccessBlockRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPublicAccessBlockResult, err error)
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", strUrl)
		},
		&DeleteAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *DeleteAccessPointPublicAccessBlockResult, err error) {
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
			assert.Equal(t, "/bucket/?publicAccessBlock&x-oss-access-point-name=ap", strUrl)
		},
		&DeleteAccessPointPublicAccessBlockRequest{
			Bucket:          Ptr("bucket"),
			AccessPointName: Ptr("ap"),
		},
		func(t *testing.T, o *DeleteAccessPointPublicAccessBlockResult, err error) {
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

func TestMockDeleteAccessPointPublicAccessBlock_Error(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPublicAccessBlockErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPublicAccessBlock(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


var testMockCreateAccessPointForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *CreateAccessPointForObjectProcessResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateAccessPointForObjectProcessResult>
  <AccessPointForObjectProcessArn>acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01</AccessPointForObjectProcessArn>
  <Alias>ap-01-45ee7945007a2f0bcb595f63e2215cxxxx-ossalias</Alias>
</CreateAccessPointForObjectProcessResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointForObjectProcessConfiguration><AccessPointName>ap-01</AccessPointName><ObjectProcessConfiguration><TransformationConfigurations><TransformationConfiguration><Actions><Action>GetObject</Action></Actions><ContentTransformation><FunctionCompute><FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn><FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn></FunctionCompute></ContentTransformation></TransformationConfiguration></TransformationConfigurations><AllowedFeatures><AllowedFeature>GetObject-Range</AllowedFeature></AllowedFeatures></ObjectProcessConfiguration></CreateAccessPointForObjectProcessConfiguration>")
		},
		&CreateAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
				AccessPointName: Ptr("ap-01"),
				ObjectProcessConfiguration: &ObjectProcessConfiguration{
					AllowedFeatures: &ObjectProcessAllowedFeatures{
						[]string{"GetObject-Range"},
					},
					TransformationConfigurations: &TransformationConfigurations{
						[]TransformationConfiguration{
							{
								Actions: &AccessPointActions{
									[]string{"GetObject"},
								},
								ContentTransformation: &ContentTransformation{
									FunctionCompute: &ObjectProcessFunctionCompute{FunctionArn: Ptr("acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01"),
										FunctionAssumeRoleArn: Ptr("acs:ram::111933544165****:role/aliyunfcdefaultrole"),
									},
								},
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *CreateAccessPointForObjectProcessResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.AccessPointForObjectProcessArn, "acs:oss:cn-hangzhou:128364106451xxxx:accesspoint/ap-01")
		},
	},
}

func TestMockCreateAccessPointForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockCreateAccessPointForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateAccessPointForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *CreateAccessPointForObjectProcessResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointForObjectProcessConfiguration><AccessPointName>ap-01</AccessPointName><ObjectProcessConfiguration><TransformationConfigurations><TransformationConfiguration><Actions><Action>GetObject</Action></Actions><ContentTransformation><FunctionCompute><FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn><FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn></FunctionCompute></ContentTransformation></TransformationConfiguration></TransformationConfigurations><AllowedFeatures><AllowedFeature>GetObject-Range</AllowedFeature></AllowedFeatures></ObjectProcessConfiguration></CreateAccessPointForObjectProcessConfiguration>")
		},
		&CreateAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
				AccessPointName: Ptr("ap-01"),
				ObjectProcessConfiguration: &ObjectProcessConfiguration{
					AllowedFeatures: &ObjectProcessAllowedFeatures{
						[]string{"GetObject-Range"},
					},
					TransformationConfigurations: &TransformationConfigurations{
						[]TransformationConfiguration{
							{
								Actions: &AccessPointActions{
									[]string{"GetObject"},
								},
								ContentTransformation: &ContentTransformation{
									FunctionCompute: &ObjectProcessFunctionCompute{FunctionArn: Ptr("acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01"),
										FunctionAssumeRoleArn: Ptr("acs:ram::111933544165****:role/aliyunfcdefaultrole"),
									},
								},
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *CreateAccessPointForObjectProcessResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<CreateAccessPointForObjectProcessConfiguration><AccessPointName>ap-01</AccessPointName><ObjectProcessConfiguration><TransformationConfigurations><TransformationConfiguration><Actions><Action>GetObject</Action></Actions><ContentTransformation><FunctionCompute><FunctionArn>acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01</FunctionArn><FunctionAssumeRoleArn>acs:ram::111933544165****:role/aliyunfcdefaultrole</FunctionAssumeRoleArn></FunctionCompute></ContentTransformation></TransformationConfiguration></TransformationConfigurations><AllowedFeatures><AllowedFeature>GetObject-Range</AllowedFeature></AllowedFeatures></ObjectProcessConfiguration></CreateAccessPointForObjectProcessConfiguration>")
		},
		&CreateAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
				AccessPointName: Ptr("ap-01"),
				ObjectProcessConfiguration: &ObjectProcessConfiguration{
					AllowedFeatures: &ObjectProcessAllowedFeatures{
						[]string{"GetObject-Range"},
					},
					TransformationConfigurations: &TransformationConfigurations{
						[]TransformationConfiguration{
							{
								Actions: &AccessPointActions{
									[]string{"GetObject"},
								},
								ContentTransformation: &ContentTransformation{
									FunctionCompute: &ObjectProcessFunctionCompute{
										FunctionArn:           Ptr("acs:fc:cn-qingdao:111933544165****:services/test-oss-fc.LATEST/functions/fc-01"),
										FunctionAssumeRoleArn: Ptr("acs:ram::111933544165****:role/aliyunfcdefaultrole"),
									},
								},
							},
						},
					},
				},
			},
		},
		func(t *testing.T, o *CreateAccessPointForObjectProcessResult, err error) {
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

func TestMockCreateAccessPointForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockCreateAccessPointForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.CreateAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointForObjectProcessResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<GetAccessPointForObjectProcessResult>
  <AccessPointNameForObjectProcess>fc-ap-01</AccessPointNameForObjectProcess>
  <AccessPointForObjectProcessAlias>fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias</AccessPointForObjectProcessAlias>
  <AccessPointName>ap-01</AccessPointName>
  <AccountId>111933544165****</AccountId>
  <AccessPointForObjectProcessArn>acs:oss:cn-qingdao:11933544165****:accesspointforobjectprocess/fc-ap-01</AccessPointForObjectProcessArn>
  <CreationDate>1626769503</CreationDate>
  <Status>enable</Status>
  <Endpoints>
    <PublicEndpoint>fc-ap-01-111933544165****.oss-cn-qingdao.oss-object-process.aliyuncs.com</PublicEndpoint>
    <InternalEndpoint>fc-ap-01-111933544165****.oss-cn-qingdao-internal.oss-object-process.aliyuncs.com</InternalEndpoint>
  </Endpoints>
  <PublicAccessBlockConfiguration>
    <BlockPublicAccess>true</BlockPublicAccess>
  </PublicAccessBlockConfiguration>
</GetAccessPointForObjectProcessResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointForObjectProcessResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.True(t, *o.PublicAccessBlockConfiguration.BlockPublicAccess)
			assert.Equal(t, *o.AccessPointNameForObjectProcess, "fc-ap-01")
			assert.Equal(t, *o.AccessPointForObjectProcessAlias, "fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias")
			assert.Equal(t, *o.AccessPointName, "ap-01")
			assert.Equal(t, *o.AccountId, "111933544165****")
			assert.Equal(t, *o.AccessPointForObjectProcessArn, "acs:oss:cn-qingdao:11933544165****:accesspointforobjectprocess/fc-ap-01")
			assert.Equal(t, *o.CreationDate, "1626769503")
			assert.Equal(t, *o.AccessPointForObjectProcessStatus, "enable")
			assert.Equal(t, *o.Endpoints.PublicEndpoint, "fc-ap-01-111933544165****.oss-cn-qingdao.oss-object-process.aliyuncs.com")
			assert.Equal(t, *o.Endpoints.InternalEndpoint, "fc-ap-01-111933544165****.oss-cn-qingdao-internal.oss-object-process.aliyuncs.com")
			assert.True(t, *o.PublicAccessBlockConfiguration.BlockPublicAccess)
		},
	},
}

func TestMockGetAccessPointForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockGetAccessPointForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointForObjectProcessResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointForObjectProcessResult, err error) {
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

func TestMockGetAccessPointForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockGetAccessPointForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointForObjectProcessResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteAccessPointForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockDeleteAccessPointForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointForObjectProcessResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointForObjectProcessResult, err error) {
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

func TestMockDeleteAccessPointForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockDeleteAccessPointForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListAccessPointsForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListAccessPointsForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *ListAccessPointsForObjectProcessResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListAccessPointsForObjectProcessResult>
   <IsTruncated>true</IsTruncated>
   <NextContinuationToken>abc</NextContinuationToken>
   <AccountId>111933544165****</AccountId>
   <AccessPointsForObjectProcess>
      <AccessPointForObjectProcess>
          <AccessPointNameForObjectProcess>fc-ap-01</AccessPointNameForObjectProcess>
          <AccessPointForObjectProcessAlias>fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias</AccessPointForObjectProcessAlias>
          <AccessPointName>fc-01</AccessPointName>
          <Status>enable</Status>
      </AccessPointForObjectProcess>
   </AccessPointsForObjectProcess>
</ListAccessPointsForObjectProcessResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/?accessPointForObjectProcess", strUrl)
		},
		&ListAccessPointsForObjectProcessRequest{},
		func(t *testing.T, o *ListAccessPointsForObjectProcessResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.NextContinuationToken, "abc")
			assert.True(t, *o.IsTruncated)
			assert.Equal(t, *o.AccountId, "111933544165****")
			assert.Equal(t, len(o.AccessPointsForObjectProcess.AccessPointForObjectProcesses), 1)
			assert.Equal(t, *o.AccessPointsForObjectProcess.AccessPointForObjectProcesses[0].AccessPointNameForObjectProcess, "fc-ap-01")
			assert.Equal(t, *o.AccessPointsForObjectProcess.AccessPointForObjectProcesses[0].AccessPointNameForObjectProcess, "fc-ap-01")
			assert.Equal(t, *o.AccessPointsForObjectProcess.AccessPointForObjectProcesses[0].AccessPointForObjectProcessAlias, "fc-ap-01-3b00521f653d2b3223680ec39dbbe2****-opapalias")
			assert.Equal(t, *o.AccessPointsForObjectProcess.AccessPointForObjectProcesses[0].AccessPointName, "fc-01")
			assert.Equal(t, *o.AccessPointsForObjectProcess.AccessPointForObjectProcesses[0].Status, "enable")
		},
	},
}

func TestMockListAccessPointsForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockListAccessPointsForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListAccessPointsForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListAccessPointsForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListAccessPointsForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *ListAccessPointsForObjectProcessResult, err error)
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
			assert.Equal(t, "/?accessPointForObjectProcess", strUrl)
		},
		&ListAccessPointsForObjectProcessRequest{},
		func(t *testing.T, o *ListAccessPointsForObjectProcessResult, err error) {
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
			assert.Equal(t, "/?accessPointForObjectProcess", strUrl)
		},
		&ListAccessPointsForObjectProcessRequest{},
		func(t *testing.T, o *ListAccessPointsForObjectProcessResult, err error) {
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

func TestMockListAccessPointsForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockListAccessPointsForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.ListAccessPointsForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPolicyForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPolicyForObjectProcessResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-001/object/*"]}]}`)
		},
		&PutAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			Body:                            strings.NewReader(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-001/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyForObjectProcessResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutAccessPointPolicyForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockPutAccessPointPolicyForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutAccessPointPolicyForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *PutAccessPointPolicyForObjectProcessResult, err error)
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*"]}]}`)
		},
		&PutAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			Body:                            strings.NewReader(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyForObjectProcessResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), `{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*"]}]}`)
		},
		&PutAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
			Body:                            strings.NewReader(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["23721626365841xxxx"],"Resource":["acs:oss:cn-qingdao:111933544165xxxx:accesspointforobjectprocess/fc-ap-01/object/*"]}]}`),
		},
		func(t *testing.T, o *PutAccessPointPolicyForObjectProcessResult, err error) {
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

func TestMockPutAccessPointPolicyForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockPutAccessPointPolicyForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPolicyForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPolicyForObjectProcessResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`{
					   "Version":"1",
					   "Statement":[
					   {
						 "Action":[
						   "oss:PutObject",
						   "oss:GetObject"
						],
						"Effect":"Deny",
						"Principal":["27737962156157xxxx"],
						"Resource":[
						   "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01",
						   "acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01/object/*"
						 ]
					   }
					  ]
					 }`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyForObjectProcessResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, string(o.Body), "{\n\t\t\t\t\t   \"Version\":\"1\",\n\t\t\t\t\t   \"Statement\":[\n\t\t\t\t\t   {\n\t\t\t\t\t\t \"Action\":[\n\t\t\t\t\t\t   \"oss:PutObject\",\n\t\t\t\t\t\t   \"oss:GetObject\"\n\t\t\t\t\t\t],\n\t\t\t\t\t\t\"Effect\":\"Deny\",\n\t\t\t\t\t\t\"Principal\":[\"27737962156157xxxx\"],\n\t\t\t\t\t\t\"Resource\":[\n\t\t\t\t\t\t   \"acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01\",\n\t\t\t\t\t\t   \"acs:oss:cn-hangzhou:111933544165xxxx:accesspoint/$ap-01/object/*\"\n\t\t\t\t\t\t ]\n\t\t\t\t\t   }\n\t\t\t\t\t  ]\n\t\t\t\t\t }")
		},
	},
}

func TestMockGetAccessPointPolicyForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockGetAccessPointPolicyForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetAccessPointPolicyForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *GetAccessPointPolicyForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyForObjectProcessResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&GetAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *GetAccessPointPolicyForObjectProcessResult, err error) {
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

func TestMockGetAccessPointPolicyForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockGetAccessPointPolicyForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.GetAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPolicyForObjectProcessSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPolicyForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyForObjectProcessResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteAccessPointPolicyForObjectProcess_Success(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPolicyForObjectProcessSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteAccessPointPolicyForObjectProcessErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteAccessPointPolicyForObjectProcessRequest
	CheckOutputFn  func(t *testing.T, o *DeleteAccessPointPolicyForObjectProcessResult, err error)
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
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyForObjectProcessResult, err error) {
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
			assert.Equal(t, "/bucket/?accessPointPolicyForObjectProcess", strUrl)
			assert.Equal(t, "fc-ap-01", r.Header.Get("x-oss-access-point-for-object-process-name"))
		},
		&DeleteAccessPointPolicyForObjectProcessRequest{
			Bucket:                          Ptr("bucket"),
			AccessPointForObjectProcessName: Ptr("fc-ap-01"),
		},
		func(t *testing.T, o *DeleteAccessPointPolicyForObjectProcessResult, err error) {
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

func TestMockDeleteAccessPointPolicyForObjectProcess_Error(t *testing.T) {
	for _, c := range testMockDeleteAccessPointPolicyForObjectProcessErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockWriteGetObjectResponseSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *WriteGetObjectResponseRequest
	CheckOutputFn  func(t *testing.T, o *WriteGetObjectResponseResult, err error)
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
			assert.Equal(t, "/?x-oss-write-get-object-response", strUrl)
			assert.Equal(t, "fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com", r.Header.Get("x-oss-request-route"))
			assert.Equal(t, "token", r.Header.Get("x-oss-request-token"))
			assert.Equal(t, "200", r.Header.Get("x-oss-fwd-status"))
		},
		&WriteGetObjectResponseRequest{
			RequestRoute: Ptr("fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
			RequestToken: Ptr("token"),
			FwdStatus:    Ptr("200"),
		},
		func(t *testing.T, o *WriteGetObjectResponseResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockWriteGetObjectResponse_Success(t *testing.T) {
	for _, c := range testMockWriteGetObjectResponseSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.WriteGetObjectResponse(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockWriteGetObjectResponseErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *WriteGetObjectResponseRequest
	CheckOutputFn  func(t *testing.T, o *WriteGetObjectResponseResult, err error)
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
			assert.Equal(t, "/?x-oss-write-get-object-response", strUrl)
			assert.Equal(t, "fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com", r.Header.Get("x-oss-request-route"))
			assert.Equal(t, "token", r.Header.Get("x-oss-request-token"))
			assert.Equal(t, "200", r.Header.Get("x-oss-fwd-status"))
		},
		&WriteGetObjectResponseRequest{
			RequestRoute: Ptr("fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
			RequestToken: Ptr("token"),
			FwdStatus:    Ptr("200"),
		},
		func(t *testing.T, o *WriteGetObjectResponseResult, err error) {
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
			assert.Equal(t, "/?x-oss-write-get-object-response", strUrl)
			assert.Equal(t, "fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com", r.Header.Get("x-oss-request-route"))
			assert.Equal(t, "token", r.Header.Get("x-oss-request-token"))
			assert.Equal(t, "200", r.Header.Get("x-oss-fwd-status"))
		},
		&WriteGetObjectResponseRequest{
			RequestRoute: Ptr("fc-ap-01-***-opap.oss-cn-hangzhou.oss-object-process.aliyuncs.com"),
			RequestToken: Ptr("token"),
			FwdStatus:    Ptr("200"),
		},
		func(t *testing.T, o *WriteGetObjectResponseResult, err error) {
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

func TestMockWriteGetObjectResponse_Error(t *testing.T) {
	for _, c := range testMockWriteGetObjectResponseErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)
		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.WriteGetObjectResponse(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


