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

var testMockPutBucketOverwriteConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketOverwriteConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?overwriteConfig", r.URL.String())
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule></OverwriteConfiguration>")
		},
		&PutBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
			OverwriteConfiguration: &OverwriteConfiguration{
				Rules: []OverwriteRule{
					{
						ID:     Ptr("1"),
						Action: Ptr("forbid"),
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketOverwriteConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?overwriteConfig", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule><Rule><ID>2</ID><Action>forbid</Action><Prefix>pre</Prefix><Suffix>.txt</Suffix><Principals><Principal>1234567890</Principal></Principals></Rule></OverwriteConfiguration>")
		},
		&PutBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
			OverwriteConfiguration: &OverwriteConfiguration{
				Rules: []OverwriteRule{
					{
						ID:     Ptr("1"),
						Action: Ptr("forbid"),
					},
					{
						ID:     Ptr("2"),
						Action: Ptr("forbid"),
						Prefix: Ptr("pre"),
						Suffix: Ptr(".txt"),
						Principals: &OverwritePrincipals{
							[]string{"1234567890"},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketOverwriteConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketOverwriteConfig_Success(t *testing.T) {
	for _, c := range testMockPutBucketOverwriteConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketOverwriteConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketOverwriteConfigResult, err error)
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
			assert.Equal(t, "/bucket/?overwriteConfig", r.URL.String())
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule><Rule><ID>2</ID><Action>forbid</Action><Prefix>pre</Prefix><Suffix>.txt</Suffix><Principals><Principal>1234567890</Principal></Principals></Rule></OverwriteConfiguration>")
		},
		&PutBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
			OverwriteConfiguration: &OverwriteConfiguration{
				Rules: []OverwriteRule{
					{
						ID:     Ptr("1"),
						Action: Ptr("forbid"),
					},
					{
						ID:     Ptr("2"),
						Action: Ptr("forbid"),
						Prefix: Ptr("pre"),
						Suffix: Ptr(".txt"),
						Principals: &OverwritePrincipals{
							[]string{"1234567890"},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketOverwriteConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule><Rule><ID>2</ID><Action>forbid</Action><Prefix>pre</Prefix><Suffix>.txt</Suffix><Principals><Principal>1234567890</Principal></Principals></Rule></OverwriteConfiguration>")
		},
		&PutBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
			OverwriteConfiguration: &OverwriteConfiguration{
				Rules: []OverwriteRule{
					{
						ID:     Ptr("1"),
						Action: Ptr("forbid"),
					},
					{
						ID:     Ptr("2"),
						Action: Ptr("forbid"),
						Prefix: Ptr("pre"),
						Suffix: Ptr(".txt"),
						Principals: &OverwritePrincipals{
							[]string{"1234567890"},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketOverwriteConfigResult, err error) {
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
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
			assert.Equal(t, "PUT", r.Method)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<OverwriteConfiguration><Rule><ID>1</ID><Action>forbid</Action></Rule><Rule><ID>2</ID><Action>forbid</Action><Prefix>pre</Prefix><Suffix>.txt</Suffix><Principals><Principal>1234567890</Principal></Principals></Rule></OverwriteConfiguration>")
		},
		&PutBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
			OverwriteConfiguration: &OverwriteConfiguration{
				Rules: []OverwriteRule{
					{
						ID:     Ptr("1"),
						Action: Ptr("forbid"),
					},
					{
						ID:     Ptr("2"),
						Action: Ptr("forbid"),
						Prefix: Ptr("pre"),
						Suffix: Ptr(".txt"),
						Principals: &OverwritePrincipals{
							[]string{"1234567890"},
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketOverwriteConfigResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute PutBucketOverwriteConfig fail")
		},
	},
}

func TestMockPutBucketOverwriteConfig_Error(t *testing.T) {
	for _, c := range testMockPutBucketOverwriteConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketOverwriteConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketOverwriteConfigResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<OverwriteConfiguration>
  <Rule>
    <ID>1</ID>
    <Action>forbid</Action>
    <Prefix></Prefix>
    <Suffix></Suffix>
    <Principals/>
  </Rule>
  <Rule>
    <ID>2</ID>
    <Action>forbid</Action>
    <Prefix>pre</Prefix>
    <Suffix>.txt</Suffix>
    <Principals>
      <Principal>1234567890</Principal>
    </Principals>
  </Rule>
</OverwriteConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/?overwriteConfig", r.URL.String())
		},
		&GetBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketOverwriteConfigResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.OverwriteConfiguration.Rules), 2)

			assert.Equal(t, *o.OverwriteConfiguration.Rules[0].ID, "1")
			assert.Equal(t, *o.OverwriteConfiguration.Rules[0].Action, "forbid")
			assert.Empty(t, o.OverwriteConfiguration.Rules[0].Prefix)
			assert.Empty(t, o.OverwriteConfiguration.Rules[0].Suffix)
			assert.Equal(t, *o.OverwriteConfiguration.Rules[1].ID, "2")
			assert.Equal(t, *o.OverwriteConfiguration.Rules[1].Action, "forbid")
			assert.Equal(t, *o.OverwriteConfiguration.Rules[1].Prefix, "pre")
			assert.Equal(t, *o.OverwriteConfiguration.Rules[1].Suffix, ".txt")
			assert.Equal(t, o.OverwriteConfiguration.Rules[1].Principals.Principals[0], "1234567890")
		},
	},
}

func TestMockGetBucketOverwriteConfig_Success(t *testing.T) {
	for _, c := range testMockGetBucketOverwriteConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketOverwriteConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketOverwriteConfigResult, err error)
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
			assert.Equal(t, "/bucket/?overwriteConfig", r.URL.String())
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketOverwriteConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
		},
		&GetBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketOverwriteConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
		},
		&GetBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketOverwriteConfigResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketOverwriteConfig fail")
		},
	},
}

func TestMockGetBucketOverwriteConfig_Error(t *testing.T) {
	for _, c := range testMockGetBucketOverwriteConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketOverwriteConfigSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketOverwriteConfigResult, err error)
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
		},
		&DeleteBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketOverwriteConfigResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketOverwriteConfig_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketOverwriteConfigSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketOverwriteConfigErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketOverwriteConfigRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketOverwriteConfigResult, err error)
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
		},
		&DeleteBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketOverwriteConfigResult, err error) {
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
			assert.Equal(t, "/bucket/?overwriteConfig", strUrl)
		},
		&DeleteBucketOverwriteConfigRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketOverwriteConfigResult, err error) {
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

func TestMockDeleteBucketOverwriteConfig_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketOverwriteConfigErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketOverwriteConfig(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


