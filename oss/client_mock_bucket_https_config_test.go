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

var testMockPutBucketHttpsConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketHttpsConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketHttpsConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<HttpsConfiguration><TLS><Enable>true</Enable><TLSVersion>TLSv1.2</TLSVersion><TLSVersion>TLSv1.3</TLSVersion></TLS></HttpsConfiguration>")
		},
		&PutBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
			HttpsConfiguration: &HttpsConfiguration{
				TLS: &TLS{
					Enable:      Ptr(true),
					TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
				},
			},
		},
		func(t *testing.T, o *PutBucketHttpsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<HttpsConfiguration><TLS><Enable>false</Enable></TLS></HttpsConfiguration>")
		},
		&PutBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
			HttpsConfiguration: &HttpsConfiguration{
				TLS: &TLS{
					Enable: Ptr(false),
				},
			},
		},
		func(t *testing.T, o *PutBucketHttpsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<HttpsConfiguration><TLS><Enable>true</Enable><TLSVersion>TLSv1.2</TLSVersion><TLSVersion>TLSv1.3</TLSVersion></TLS><CipherSuite><Enable>true</Enable><StrongCipherSuite>false</StrongCipherSuite><CustomCipherSuite>ECDHE-ECDSA-AES128-SHA256</CustomCipherSuite><CustomCipherSuite>ECDHE-RSA-AES128-GCM-SHA256</CustomCipherSuite><CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite></CipherSuite></HttpsConfiguration>")
		},
		&PutBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
			HttpsConfiguration: &HttpsConfiguration{
				TLS: &TLS{
					Enable:      Ptr(true),
					TLSVersions: []string{"TLSv1.2", "TLSv1.3"},
				},
				CipherSuite: &CipherSuite{
					Enable:            Ptr(true),
					StrongCipherSuite: Ptr(false),
					CustomCipherSuites: []string{
						"ECDHE-ECDSA-AES128-SHA256", "ECDHE-RSA-AES128-GCM-SHA256", "ECDHE-ECDSA-AES256-CCM8",
					},
					TLS13CustomCipherSuites: []string{
						"ECDHE-ECDSA-AES256-CCM8", "ECDHE-ECDSA-AES256-CCM8", "ECDHE-ECDSA-AES256-CCM8",
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketHttpsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketHttpsConfig_Success(t *testing.T) {
	for _, c := range testMockPutBucketHttpsConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketHttpsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketHttpsConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketHttpsConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketHttpsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<HttpsConfiguration><TLS><Enable>false</Enable></TLS></HttpsConfiguration>")
		},
		&PutBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
			HttpsConfiguration: &HttpsConfiguration{
				TLS: &TLS{
					Enable: Ptr(false),
				},
			},
		},
		func(t *testing.T, o *PutBucketHttpsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?httpsConfig", strUrl)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<HttpsConfiguration><TLS><Enable>false</Enable></TLS></HttpsConfiguration>")
		},
		&PutBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
			HttpsConfiguration: &HttpsConfiguration{
				TLS: &TLS{
					Enable: Ptr(false),
				},
			},
		},
		func(t *testing.T, o *PutBucketHttpsConfigResult, err error) {
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

func TestMockPutBucketHttpsConfig_Error(t *testing.T) {
	for _, c := range testMockPutBucketHttpsConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketHttpsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketHttpsConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketHttpsConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketHttpsConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<HttpsConfiguration>
  <TLS>
    <Enable>true</Enable>
    <TLSVersion>TLSv1.2</TLSVersion>
    <TLSVersion>TLSv1.3</TLSVersion>
  </TLS>
</HttpsConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.True(t, *o.HttpsConfiguration.TLS.Enable)
			assert.Equal(t, o.HttpsConfiguration.TLS.TLSVersions[0], "TLSv1.2")
			assert.Equal(t, o.HttpsConfiguration.TLS.TLSVersions[1], "TLSv1.3")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<HttpsConfiguration>
  <TLS>
    <Enable>false</Enable>
  </TLS>
</HttpsConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.False(t, *o.HttpsConfiguration.TLS.Enable)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<HttpsConfiguration><TLS><Enable>true</Enable><TLSVersion>TLSv1.2</TLSVersion><TLSVersion>TLSv1.3</TLSVersion></TLS><CipherSuite><Enable>true</Enable><StrongCipherSuite>false</StrongCipherSuite><CustomCipherSuite>ECDHE-ECDSA-AES128-SHA256</CustomCipherSuite><CustomCipherSuite>ECDHE-RSA-AES128-GCM-SHA256</CustomCipherSuite><CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite><TLS13CustomCipherSuite>ECDHE-ECDSA-AES256-CCM8</TLS13CustomCipherSuite></CipherSuite></HttpsConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.True(t, *o.HttpsConfiguration.TLS.Enable)
			assert.Equal(t, o.HttpsConfiguration.TLS.TLSVersions[0], "TLSv1.2")
			assert.Equal(t, o.HttpsConfiguration.TLS.TLSVersions[1], "TLSv1.3")
			assert.True(t, *o.HttpsConfiguration.CipherSuite.Enable)
			assert.Equal(t, len(o.HttpsConfiguration.CipherSuite.TLS13CustomCipherSuites), 3)
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.TLS13CustomCipherSuites[0], "ECDHE-ECDSA-AES256-CCM8")
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.TLS13CustomCipherSuites[1], "ECDHE-ECDSA-AES256-CCM8")
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.TLS13CustomCipherSuites[2], "ECDHE-ECDSA-AES256-CCM8")
			assert.Equal(t, len(o.HttpsConfiguration.CipherSuite.CustomCipherSuites), 3)
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.CustomCipherSuites[0], "ECDHE-ECDSA-AES128-SHA256")
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.CustomCipherSuites[1], "ECDHE-RSA-AES128-GCM-SHA256")
			assert.Equal(t, o.HttpsConfiguration.CipherSuite.CustomCipherSuites[2], "ECDHE-ECDSA-AES256-CCM8")
		},
	},
}

func TestMockGetBucketHttpsConfig_Success(t *testing.T) {
	for _, c := range testMockGetBucketHttpsConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketHttpsConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketHttpsConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketHttpsConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketHttpsConfigResult, err error)
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
			assert.Equal(t, "/bucket/?httpsConfig", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?httpsConfig", strUrl)
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?httpsConfig", strUrl)
		},
		&GetBucketHttpsConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketHttpsConfigResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketHttpsConfig fail")
		},
	},
}

func TestMockGetBucketHttpsConfig_Error(t *testing.T) {
	for _, c := range testMockGetBucketHttpsConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketHttpsConfig(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}


