package oss

import (
	"testing"
	"context"
	"errors"
	"html"
	"io"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockOpenMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OpenMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *OpenMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
			Mode:   Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
			Mode:   Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery&mode=semantic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
			Mode:   Ptr("semantic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockOpenMetaQuery_Success(t *testing.T) {
	for _, c := range testMockOpenMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OpenMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockOpenMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *OpenMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *OpenMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=add&metaQuery", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
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

func TestMockOpenMetaQuery_Error(t *testing.T) {
	for _, c := range testMockOpenMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.OpenMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetMetaQueryStatusSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetMetaQueryStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetMetaQueryStatusResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQueryStatus>
  <State>Running</State>
  <Phase>FullScanning</Phase>
  <CreateTime>2021-08-02T10:49:17.289372919+08:00</CreateTime>
  <UpdateTime>2021-08-02T10:49:17.289372919+08:00</UpdateTime>
</MetaQueryStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.MetaQueryStatus.State, "Running")
			assert.Equal(t, *o.MetaQueryStatus.Phase, "FullScanning")
			assert.Equal(t, *o.MetaQueryStatus.CreateTime, "2021-08-02T10:49:17.289372919+08:00")
			assert.Equal(t, *o.MetaQueryStatus.UpdateTime, "2021-08-02T10:49:17.289372919+08:00")
		},
	},
}

func TestMockGetMetaQueryStatus_Success(t *testing.T) {
	for _, c := range testMockGetMetaQueryStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetMetaQueryStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetMetaQueryStatusErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetMetaQueryStatusRequest
	CheckOutputFn  func(t *testing.T, o *GetMetaQueryStatusResult, err error)
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
			assert.Equal(t, "/bucket/?metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
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
			assert.Equal(t, "/bucket/?metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
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

func TestMockGetMetaQueryStatus_Error(t *testing.T) {
	for _, c := range testMockGetMetaQueryStatusErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetMetaQueryStatus(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken>
  <Files>
    <File>
      <Filename>exampleobject.txt</Filename>
      <Size>120</Size>
      <FileModifiedTime>2021-06-29T15:04:05.000000000Z07:00</FileModifiedTime>
      <OSSObjectType>Normal</OSSObjectType>
      <OSSStorageClass>Standard</OSSStorageClass>
      <ObjectACL>default</ObjectACL>
      <ETag>"fba9dede5f27731c9771645a3986****"</ETag>
      <OSSCRC64>4858A48BD1466884</OSSCRC64>
      <OSSTaggingCount>2</OSSTaggingCount>
      <OSSTagging>
        <Tagging>
          <Key>owner</Key>
          <Value>John</Value>
        </Tagging>
        <Tagging>
          <Key>type</Key>
          <Value>document</Value>
        </Tagging>
      </OSSTagging>
      <OSSUserMeta>
        <UserMeta>
          <Key>x-oss-meta-location</Key>
          <Value>hangzhou</Value>
        </UserMeta>
      </OSSUserMeta>
    </File>
  </Files>
  <Aggregations>
    <Aggregation>
      <Field>Size</Field>
      <Operation>sum</Operation>
      <Value>4859250309</Value>
    </Aggregation>
    <Aggregation>
      <Field>Size</Field>
      <Operation>max</Operation>
      <Value>2235483240</Value>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=query&metaQuery", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>5</MaxResults><Query>{\"Field\": \"Size\",\"Value\": \"1048576\",\"Operation\": \"gt\"}</Query><Sort>Size</Sort><Order>asc</Order><Aggregations><Aggregation><Field>Size</Field><Operation>sum</Operation></Aggregation><Aggregation><Field>Size</Field><Operation>max</Operation></Aggregation></Aggregations><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: Ptr("bucket"),
			MetaQuery: &MetaQuery{
				NextToken:  Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
				MaxResults: Ptr(int64(5)),
				Query:      Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
				Sort:       Ptr("Size"),
				Order:      Ptr(MetaQueryOrderAsc),
				Aggregations: &MetaQueryAggregations{
					[]MetaQueryAggregation{
						{
							Field:     Ptr("Size"),
							Operation: Ptr("sum"),
						},
						{
							Field:     Ptr("Size"),
							Operation: Ptr("max"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].Filename, "exampleobject.txt")
			assert.Equal(t, *o.Files[0].Size, int64(120))
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2021-06-29T15:04:05.000000000Z07:00")
			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, *o.Files[0].ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSTagging[0].Key, "owner")
			assert.Equal(t, *o.Files[0].OSSTagging[0].Value, "John")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Key, "type")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Value, "document")
			assert.Equal(t, len(o.Files[0].OSSUserMeta), 1)
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Key, "x-oss-meta-location")
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Value, "hangzhou")
			assert.Equal(t, len(o.Aggregations), 2)
			assert.Equal(t, *o.Aggregations[0].Field, "Size")
			assert.Equal(t, *o.Aggregations[0].Operation, "sum")
			assert.Equal(t, *o.Aggregations[0].Value, float64(4859250309))
			assert.Equal(t, *o.Aggregations[1].Field, "Size")
			assert.Equal(t, *o.Aggregations[1].Operation, "max")
			assert.Equal(t, *o.Aggregations[1].Value, float64(2235483240))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<MetaQuery>
  <NextToken></NextToken>
  <Aggregations>
    <Aggregation>
      <Field>Size</Field>
      <Operation>sum</Operation>
      <Value>30930054</Value>
    </Aggregation>
    <Aggregation>
      <Field>Size</Field>
      <Operation>group</Operation>
      <Groups>
        <Group>
          <Value>1536000</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>5472362</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>10354204</Value>
          <Count>1</Count>
        </Group>
        <Group>
          <Value>1890304</Value>
          <Count>3</Count>
        </Group>
        <Group>
          <Value>2632192</Value>
          <Count>3</Count>
        </Group>
      </Groups>
    </Aggregation>
  </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=query&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>5</MaxResults><Query>{\"Field\": \"Size\",\"Value\": \"1048576\",\"Operation\": \"gt\"}</Query><Sort>Size</Sort><Order>asc</Order><Aggregations><Aggregation><Field>Size</Field><Operation>sum</Operation></Aggregation><Aggregation><Field>Size</Field><Operation>group</Operation></Aggregation></Aggregations><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: Ptr("bucket"),
			Mode:   Ptr("basic"),
			MetaQuery: &MetaQuery{
				NextToken:  Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
				MaxResults: Ptr(int64(5)),
				Query:      Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
				Sort:       Ptr("Size"),
				Order:      Ptr(MetaQueryOrderAsc),
				Aggregations: &MetaQueryAggregations{
					[]MetaQueryAggregation{
						{
							Field:     Ptr("Size"),
							Operation: Ptr("sum"),
						},
						{
							Field:     Ptr("Size"),
							Operation: Ptr("group"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Aggregations), 2)
			assert.Equal(t, *o.Aggregations[0].Field, "Size")
			assert.Equal(t, *o.Aggregations[0].Operation, "sum")
			assert.Equal(t, *o.Aggregations[0].Value, float64(30930054))
			assert.Equal(t, *o.Aggregations[1].Field, "Size")
			assert.Equal(t, *o.Aggregations[1].Operation, "group")
			assert.Equal(t, len(o.Aggregations[1].Groups.Groups), 5)
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[0].Value, "1536000")
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[0].Count, int64(1))
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[1].Value, "5472362")
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[1].Count, int64(1))
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[2].Value, "10354204")
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[2].Count, int64(1))
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[3].Value, "1890304")
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[3].Count, int64(3))
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[4].Value, "2632192")
			assert.Equal(t, *o.Aggregations[1].Groups.Groups[4].Count, int64(3))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
  <Files>
    <File>
      <URI>oss://bucket/sample-object.jpg</URI>
      <Filename>sample-object.jpg</Filename>
      <Size>1000</Size>
      <ObjectACL>default</ObjectACL>
      <FileModifiedTime>2021-06-29T14:50:14.011643661+08:00</FileModifiedTime>
      <ServerSideEncryption>AES256</ServerSideEncryption>
      <ServerSideEncryptionCustomerAlgorithm>SM4</ServerSideEncryptionCustomerAlgorithm>
      <ETag>"1D9C280A7C4F67F7EF873E28449****"</ETag>
      <OSSCRC64>559890638950338001</OSSCRC64>
      <ProduceTime>2021-06-29T14:50:15.011643661+08:00</ProduceTime>
      <ContentType>image/jpeg</ContentType>
      <MediaType>image</MediaType>
      <LatLong>30.134390,120.074997</LatLong>
      <Title>test</Title>
      <OSSExpiration>2024-12-01T12:00:00.000Z</OSSExpiration>
      <AccessControlAllowOrigin>https://aliyundoc.com</AccessControlAllowOrigin>
      <AccessControlRequestMethod>PUT</AccessControlRequestMethod>
      <ServerSideDataEncryption>SM4</ServerSideDataEncryption>
      <ServerSideEncryptionKeyId>9468da86-3509-4f8d-a61e-6eab1eac****</ServerSideEncryptionKeyId>
      <CacheControl>no-cache</CacheControl>
      <ContentDisposition>attachment; filename=test.jpg</ContentDisposition>
      <ContentEncoding>UTF-8</ContentEncoding>
      <ContentLanguage>zh-CN</ContentLanguage>
      <ImageHeight>500</ImageHeight>
      <ImageWidth>270</ImageWidth>
      <VideoWidth>1080</VideoWidth>
      <VideoHeight>1920</VideoHeight>
      <VideoStreams>
        <VideoStream>
          <CodecName>h264</CodecName>
          <Language>en</Language>
          <Bitrate>5407765</Bitrate>
          <FrameRate>25/1</FrameRate>
          <StartTime>0</StartTime>
          <Duration>22.88</Duration>
          <FrameCount>572</FrameCount>
          <BitDepth>8</BitDepth>
          <PixelFormat>yuv420p</PixelFormat>
          <ColorSpace>bt709</ColorSpace>
          <Height>720</Height>
          <Width>1280</Width>
        </VideoStream>
        <VideoStream>
          <CodecName>h264</CodecName>
          <Language>en</Language>
          <Bitrate>5407765</Bitrate>
          <FrameRate>25/1</FrameRate>
          <StartTime>0</StartTime>
          <Duration>22.88</Duration>
          <FrameCount>572</FrameCount>
          <BitDepth>8</BitDepth>
          <PixelFormat>yuv420p</PixelFormat>
          <ColorSpace>bt709</ColorSpace>
          <Height>720</Height>
          <Width>1280</Width>
        </VideoStream>
      </VideoStreams>
      <AudioStreams>
        <AudioStream>
          <CodecName>aac</CodecName>
          <Bitrate>1048576</Bitrate>
          <SampleRate>48000</SampleRate>
          <StartTime>0.0235</StartTime>
          <Duration>3.690667</Duration>
          <Channels>2</Channels>
          <Language>en</Language>
        </AudioStream>
      </AudioStreams>
      <Subtitles>
        <Subtitle>
          <CodecName>mov_text</CodecName>
          <Language>en</Language>
          <StartTime>0</StartTime>
          <Duration>71.378</Duration>
        </Subtitle>
        <Subtitle>
          <CodecName>mov_text</CodecName>
          <Language>en</Language>
          <StartTime>72</StartTime>
          <Duration>71.378</Duration>
        </Subtitle>
      </Subtitles>
      <Bitrate>5407765</Bitrate>
      <Artist>Jane</Artist>
      <AlbumArtist>Jenny</AlbumArtist>
      <Composer>Jane</Composer>
      <Performer>Jane</Performer>
      <Album>FirstAlbum</Album>
      <Duration>71.378</Duration>
      <Addresses>
        <Address>
          <AddressLine>中国浙江省杭州市余杭区文一西路969号</AddressLine>
          <City>杭州市</City>
          <Country>中国</Country>
          <District>余杭区</District>
          <Language>zh-Hans</Language>
          <Province>浙江省</Province>
          <Township>文一西路</Township>
        </Address>
        <Address>
          <AddressLine>中国浙江省杭州市余杭区文一西路970号</AddressLine>
          <City>杭州市</City>
          <Country>中国</Country>
          <District>余杭区</District>
          <Language>zh-Hans</Language>
          <Province>浙江省</Province>
          <Township>文一西路</Township>
        </Address>
      </Addresses>
      <OSSObjectType>Normal</OSSObjectType>
      <OSSStorageClass>Standard</OSSStorageClass>
      <OSSTaggingCount>2</OSSTaggingCount>
      <OSSTagging>
        <Tagging>
          <Key>key</Key>
          <Value>val</Value>
        </Tagging>
        <Tagging>
          <Key>key2</Key>
          <Value>val2</Value>
        </Tagging>
      </OSSTagging>
      <OSSUserMeta>
        <UserMeta>
          <Key>key</Key>
          <Value>val</Value>
        </UserMeta>
      </OSSUserMeta>
      <Insights>
        <Image>
          <Caption>There stands a person.</Caption>
          <Description>In the picture, there is a person wearing a dark suit jacket with a white shirt underneath. The background is a gradient from light blue to gray</Description>
        </Image>
		<Video>
          <Caption>The video shows two different scenes</Caption>
        </Video>
      </Insights>
    </File>
  </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?comp=query&metaQuery&mode=semantic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>99</MaxResults><Query>Overlook the snow-covered forest</Query><MediaTypes><MediaType>image</MediaType></MediaTypes><SimpleQuery>{\"Operation\":\"gt\", \"Field\": \"Size\", \"Value\": \"30\"}</SimpleQuery></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: Ptr("bucket"),
			Mode:   Ptr("semantic"),
			MetaQuery: &MetaQuery{
				MaxResults: Ptr(int64(99)),
				Query:      Ptr("Overlook the snow-covered forest"),
				MediaTypes: &MetaQueryMediaTypes{
					MediaTypes: []string{"image"},
				},
				SimpleQuery: Ptr(`{"Operation":"gt", "Field": "Size", "Value": "30"}`),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.Files), 1)
			assert.Equal(t, *o.Files[0].URI, "oss://bucket/sample-object.jpg")
			assert.Equal(t, *o.Files[0].Filename, "sample-object.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(1000))
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2021-06-29T14:50:14.011643661+08:00")
			assert.Equal(t, *o.Files[0].ServerSideEncryption, "AES256")
			assert.Equal(t, *o.Files[0].ServerSideEncryptionCustomerAlgorithm, "SM4")
			assert.Equal(t, *o.Files[0].ETag, "\"1D9C280A7C4F67F7EF873E28449****\"")
			assert.Equal(t, *o.Files[0].OSSCRC64, "559890638950338001")
			assert.Equal(t, *o.Files[0].ProduceTime, "2021-06-29T14:50:15.011643661+08:00")
			assert.Equal(t, *o.Files[0].ContentType, "image/jpeg")
			assert.Equal(t, *o.Files[0].MediaType, "image")
			assert.Equal(t, *o.Files[0].LatLong, "30.134390,120.074997")
			assert.Equal(t, *o.Files[0].Title, "test")
			assert.Equal(t, *o.Files[0].OSSExpiration, "2024-12-01T12:00:00.000Z")
			assert.Equal(t, *o.Files[0].AccessControlAllowOrigin, "https://aliyundoc.com")
			assert.Equal(t, *o.Files[0].AccessControlRequestMethod, "PUT")
			assert.Equal(t, *o.Files[0].ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.Files[0].ServerSideEncryptionKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, *o.Files[0].CacheControl, "no-cache")
			assert.Equal(t, *o.Files[0].ContentDisposition, "attachment; filename=test.jpg")
			assert.Equal(t, *o.Files[0].ContentEncoding, "UTF-8")
			assert.Equal(t, *o.Files[0].ContentLanguage, "zh-CN")
			assert.Equal(t, *o.Files[0].ImageHeight, int64(500))
			assert.Equal(t, *o.Files[0].ImageWidth, int64(270))
			assert.Equal(t, *o.Files[0].VideoWidth, int64(1080))
			assert.Equal(t, *o.Files[0].VideoHeight, int64(1920))
			assert.Equal(t, len(o.Files[0].VideoStreams), 2)
			assert.Equal(t, *o.Files[0].VideoStreams[0].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Language, "en")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Bitrate, int64(5407765))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameRate, "25/1")
			assert.Equal(t, *o.Files[0].VideoStreams[0].StartTime, float64(0))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Duration, float64(22.88))
			assert.Equal(t, *o.Files[0].VideoStreams[0].FrameCount, int64(572))
			assert.Equal(t, *o.Files[0].VideoStreams[0].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[0].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[0].ColorSpace, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[0].Height, int64(720))
			assert.Equal(t, *o.Files[0].VideoStreams[0].Width, int64(1280))

			assert.Equal(t, *o.Files[0].VideoStreams[1].CodecName, "h264")
			assert.Equal(t, *o.Files[0].VideoStreams[1].Language, "en")
			assert.Equal(t, *o.Files[0].VideoStreams[1].Bitrate, int64(5407765))
			assert.Equal(t, *o.Files[0].VideoStreams[1].FrameRate, "25/1")
			assert.Equal(t, *o.Files[0].VideoStreams[1].StartTime, float64(0))
			assert.Equal(t, *o.Files[0].VideoStreams[1].Duration, float64(22.88))
			assert.Equal(t, *o.Files[0].VideoStreams[1].FrameCount, int64(572))
			assert.Equal(t, *o.Files[0].VideoStreams[1].BitDepth, int64(8))
			assert.Equal(t, *o.Files[0].VideoStreams[1].PixelFormat, "yuv420p")
			assert.Equal(t, *o.Files[0].VideoStreams[1].ColorSpace, "bt709")
			assert.Equal(t, *o.Files[0].VideoStreams[1].Height, int64(720))
			assert.Equal(t, *o.Files[0].VideoStreams[1].Width, int64(1280))

			assert.Equal(t, len(o.Files[0].AudioStreams), 1)
			assert.Equal(t, *o.Files[0].AudioStreams[0].CodecName, "aac")
			assert.Equal(t, *o.Files[0].AudioStreams[0].Bitrate, int64(1048576))
			assert.Equal(t, *o.Files[0].AudioStreams[0].SampleRate, int64(48000))
			assert.Equal(t, *o.Files[0].AudioStreams[0].StartTime, float64(0.0235))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Duration, float64(3.690667))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Channels, int64(2))
			assert.Equal(t, *o.Files[0].AudioStreams[0].Language, "en")

			assert.Equal(t, len(o.Files[0].Subtitles), 2)
			assert.Equal(t, *o.Files[0].Subtitles[0].CodecName, "mov_text")
			assert.Equal(t, *o.Files[0].Subtitles[0].Language, "en")
			assert.Equal(t, *o.Files[0].Subtitles[0].StartTime, float64(0))
			assert.Equal(t, *o.Files[0].Subtitles[0].Duration, float64(71.378))
			assert.Equal(t, *o.Files[0].Subtitles[1].CodecName, "mov_text")
			assert.Equal(t, *o.Files[0].Subtitles[1].Language, "en")
			assert.Equal(t, *o.Files[0].Subtitles[1].StartTime, float64(72))
			assert.Equal(t, *o.Files[0].Subtitles[1].Duration, float64(71.378))

			assert.Equal(t, *o.Files[0].Bitrate, int64(5407765))
			assert.Equal(t, *o.Files[0].Artist, "Jane")
			assert.Equal(t, *o.Files[0].AlbumArtist, "Jenny")
			assert.Equal(t, *o.Files[0].Composer, "Jane")
			assert.Equal(t, *o.Files[0].Performer, "Jane")
			assert.Equal(t, *o.Files[0].Album, "FirstAlbum")
			assert.Equal(t, *o.Files[0].Duration, float64(71.378))

			assert.Equal(t, len(o.Files[0].Addresses), 2)
			assert.Equal(t, *o.Files[0].Addresses[0].AddressLine, "中国浙江省杭州市余杭区文一西路969号")
			assert.Equal(t, *o.Files[0].Addresses[0].City, "杭州市")
			assert.Equal(t, *o.Files[0].Addresses[0].Country, "中国")
			assert.Equal(t, *o.Files[0].Addresses[0].District, "余杭区")
			assert.Equal(t, *o.Files[0].Addresses[0].Language, "zh-Hans")
			assert.Equal(t, *o.Files[0].Addresses[0].Province, "浙江省")
			assert.Equal(t, *o.Files[0].Addresses[0].Township, "文一西路")

			assert.Equal(t, *o.Files[0].Addresses[1].AddressLine, "中国浙江省杭州市余杭区文一西路970号")
			assert.Equal(t, *o.Files[0].Addresses[1].City, "杭州市")
			assert.Equal(t, *o.Files[0].Addresses[1].Country, "中国")
			assert.Equal(t, *o.Files[0].Addresses[1].District, "余杭区")
			assert.Equal(t, *o.Files[0].Addresses[1].Language, "zh-Hans")
			assert.Equal(t, *o.Files[0].Addresses[1].Province, "浙江省")
			assert.Equal(t, *o.Files[0].Addresses[1].Township, "文一西路")

			assert.Equal(t, *o.Files[0].OSSObjectType, "Normal")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].OSSTaggingCount, int64(2))
			assert.Equal(t, *o.Files[0].OSSTagging[0].Key, "key")
			assert.Equal(t, *o.Files[0].OSSTagging[0].Value, "val")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Key, "key2")
			assert.Equal(t, *o.Files[0].OSSTagging[1].Value, "val2")
			assert.Equal(t, len(o.Files[0].OSSUserMeta), 1)
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Key, "key")
			assert.Equal(t, *o.Files[0].OSSUserMeta[0].Value, "val")
			assert.Equal(t, *o.Files[0].Insights.Image.Caption, "There stands a person.")
			assert.Equal(t, *o.Files[0].Insights.Image.Description, "In the picture, there is a person wearing a dark suit jacket with a white shirt underneath. The background is a gradient from light blue to gray")
			assert.Equal(t, *o.Files[0].Insights.Video.Caption, "The video shows two different scenes")
			assert.Nil(t, o.Files[0].Insights.Video.Description)
		},
	},
}

func TestMockDoMetaQuery_Success(t *testing.T) {
	for _, c := range testMockDoMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?comp=query&metaQuery", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>5</MaxResults><Query>{\"Field\": \"Size\",\"Value\": \"1048576\",\"Operation\": \"gt\"}</Query><Sort>Size</Sort><Order>asc</Order><Aggregations><Aggregation><Field>Size</Field><Operation>sum</Operation></Aggregation><Aggregation><Field>Size</Field><Operation>max</Operation></Aggregation></Aggregations><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: Ptr("bucket"),
			MetaQuery: &MetaQuery{
				NextToken:  Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
				MaxResults: Ptr(int64(5)),
				Query:      Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
				Sort:       Ptr("Size"),
				Order:      Ptr(MetaQueryOrderAsc),
				Aggregations: &MetaQueryAggregations{
					[]MetaQueryAggregation{
						{
							Field:     Ptr("Size"),
							Operation: Ptr("sum"),
						},
						{
							Field:     Ptr("Size"),
							Operation: Ptr("max"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=query&metaQuery", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, html.UnescapeString(string(body)), "<MetaQuery><MaxResults>5</MaxResults><Query>{\"Field\": \"Size\",\"Value\": \"1048576\",\"Operation\": \"gt\"}</Query><Sort>Size</Sort><Order>asc</Order><Aggregations><Aggregation><Field>Size</Field><Operation>sum</Operation></Aggregation><Aggregation><Field>Size</Field><Operation>max</Operation></Aggregation></Aggregations><NextToken>MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****</NextToken></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: Ptr("bucket"),
			MetaQuery: &MetaQuery{
				NextToken:  Ptr("MTIzNDU2Nzg6aW1tdGVzdDpleGFtcGxlYnVja2V0OmRhdGFzZXQwMDE6b3NzOi8vZXhhbXBsZWJ1Y2tldC9zYW1wbGVvYmplY3QxLmpw****"),
				MaxResults: Ptr(int64(5)),
				Query:      Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
				Sort:       Ptr("Size"),
				Order:      Ptr(MetaQueryOrderAsc),
				Aggregations: &MetaQueryAggregations{
					[]MetaQueryAggregation{
						{
							Field:     Ptr("Size"),
							Operation: Ptr("sum"),
						},
						{
							Field:     Ptr("Size"),
							Operation: Ptr("max"),
						},
					},
				},
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
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

func TestMockDoMetaQuery_Error(t *testing.T) {
	for _, c := range testMockDoMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCloseMetaQuerySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CloseMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *CloseMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?comp=delete&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockCloseMetaQuery_Success(t *testing.T) {
	for _, c := range testMockCloseMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CloseMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCloseMetaQueryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CloseMetaQueryRequest
	CheckOutputFn  func(t *testing.T, o *CloseMetaQueryResult, err error)
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
			assert.Equal(t, "/bucket/?comp=delete&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
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
			assert.Equal(t, "/bucket/?comp=delete&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
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

func TestMockCloseMetaQuery_Error(t *testing.T) {
	for _, c := range testMockCloseMetaQueryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CloseMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


var testMockDoMetaQueryActionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryActionRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryActionResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<WorkflowParametersString></WorkflowParametersString>
<TemplateId>Official:OSSBasicMeta</TemplateId>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxBindCount>10</DatasetMaxBindCount>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>
</Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&clusterType=auto&datasetName=your_dataset&description=this+is+a+demo&metaQuery&templateId=Official%3AOSSBasicMeta", urlStr)
		},
		&DoMetaQueryActionRequest{
			Bucket: Ptr("bucket"),
			Action: Ptr("createDataset"),
			RequestCommon: RequestCommon{
				Parameters: map[string]string{
					"datasetName": "your_dataset",
					"description": "this is a demo",
					"templateId":  "Official:OSSBasicMeta",
					"clusterType": "auto",
				},
			},
		},
		func(t *testing.T, o *DoMetaQueryActionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			getBody, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(getBody), "<CreateDatasetResponse>\n<Dataset>\n<DatasetName>test-dataset</DatasetName>\n<WorkflowParameters></WorkflowParameters>\n<WorkflowParametersString></WorkflowParametersString>\n<TemplateId>Official:OSSBasicMeta</TemplateId>\n<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>\n<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>\n<Description>this is a demo</Description>\n<DatasetMaxBindCount>10</DatasetMaxBindCount>\n<DatasetMaxFileCount>100000000</DatasetMaxFileCount>\n<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>\n<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>\n<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>\n<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>\n</Dataset>\n</CreateDatasetResponse>")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListDatasetsResponse><Datasets><Dataset><DatasetName>oss_1760225545084331_demo-1889</DatasetName><WorkflowParameters></WorkflowParameters><WorkflowParametersString></WorkflowParametersString><TemplateId>Official:OSSBasicMeta</TemplateId><CreateTime>2026-04-28T15:03:47.028110649+08:00</CreateTime><UpdateTime>2026-04-28T15:03:47.156498415+08:00</UpdateTime><DatasetMaxBindCount>10</DatasetMaxBindCount><DatasetMaxFileCount>100000000</DatasetMaxFileCount><DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount><DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount><DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize><DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig></Dataset></Datasets></ListDatasetsResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", urlStr)
		},
		&DoMetaQueryActionRequest{
			Bucket: Ptr("bucket"),
			Action: Ptr("listDatasets"),
		},
		func(t *testing.T, o *DoMetaQueryActionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			getBody, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(getBody), "<ListDatasetsResponse><Datasets><Dataset><DatasetName>oss_1760225545084331_demo-1889</DatasetName><WorkflowParameters></WorkflowParameters><WorkflowParametersString></WorkflowParametersString><TemplateId>Official:OSSBasicMeta</TemplateId><CreateTime>2026-04-28T15:03:47.028110649+08:00</CreateTime><UpdateTime>2026-04-28T15:03:47.156498415+08:00</UpdateTime><DatasetMaxBindCount>10</DatasetMaxBindCount><DatasetMaxFileCount>100000000</DatasetMaxFileCount><DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount><DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount><DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize><DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig></Dataset></Datasets></ListDatasetsResponse>")
		},
	},
}

func TestMockDoMetaQueryAction_Success(t *testing.T) {
	for _, c := range testMockDoMetaQueryActionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQueryAction(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoMetaQueryActionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoMetaQueryActionRequest
	CheckOutputFn  func(t *testing.T, o *DoMetaQueryActionResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", urlStr)
		},
		&DoMetaQueryActionRequest{
			Bucket: Ptr("bucket"),
			Action: Ptr("listDatasets"),
		},
		func(t *testing.T, o *DoMetaQueryActionResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", urlStr)
		},
		&DoMetaQueryActionRequest{
			Bucket: Ptr("bucket"),
			Action: Ptr("listDatasets"),
		},
		func(t *testing.T, o *DoMetaQueryActionResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", urlStr)
		},
		&DoMetaQueryActionRequest{
			Bucket: Ptr("bucket"),
			Action: Ptr("listDatasets"),
		},
		func(t *testing.T, o *DoMetaQueryActionResult, err error) {
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

func TestMockDoMetaQueryAction_Error(t *testing.T) {
	for _, c := range testMockDoMetaQueryActionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoMetaQueryAction(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoDataPipeLineActionSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoDataPipeLineActionRequest
	CheckOutputFn  func(t *testing.T, o *DoDataPipeLineActionResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<ListDataPipelineConfigurationsResult>
  <DataPipelineConfigurations>
    <DataPipelineConfiguration>
      <DataPipelineName>my-data-pipeline</DataPipelineName>
      <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
      <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
      <Status>Running</Status>
      <Sources>
          <InputBucket>my-bucket</InputBucket>
          <InputDataScope>All</InputDataScope>
          <IgnoreDelete>true</IgnoreDelete>
          <FilterConfiguration>
              <PrefixSet>prefix1/</PrefixSet>
              <PrefixSet>prefix2/prefix3/</PrefixSet>
              <ObjectMediaTypes>text</ObjectMediaTypes>
              <ObjectMediaTypes>image</ObjectMediaTypes>
              <ObjectMediaTypes>video</ObjectMediaTypes>
          </FilterConfiguration>
      </Sources>
      <DataPipelineEmbeddingConfiguration>
          <EmbeddingProvider>bailian</EmbeddingProvider>
          <ApiKey>xxxx</ApiKey>
          <Model>qwen2.5-vl-embedding</Model>
          <FPS>1</FPS>
      </DataPipelineEmbeddingConfiguration>
      <Destination>
          <VectorBucketName>my-vector-bucket</VectorBucketName>
          <VectorIndexNames>my-index</VectorIndexNames>
          <VectorKeyPrefix></VectorKeyPrefix>
          <ObjectTagToMetadata>key1</ObjectTagToMetadata>
          <ObjectTagToMetadata>key2</ObjectTagToMetadata>
          <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
      </Destination>
      <DataPipelineError>
          <ErrorMode>ignoreAndRecord</ErrorMode>
          <ErrorBucket>my-error-bucket</ErrorBucket>
          <ErrorPrefix>error-output/</ErrorPrefix>
      </DataPipelineError>
      <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
    </DataPipelineConfiguration>
  </DataPipelineConfigurations>
  <NextToken>xxx</NextToken>
</ListDataPipelineConfigurationsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", urlStr)
		},
		&DoDataPipeLineActionRequest{
			Action: Ptr("listDataPipelineConfigurations"),
		},
		func(t *testing.T, o *DoDataPipeLineActionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			getBody, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(getBody), "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<ListDataPipelineConfigurationsResult>\n  <DataPipelineConfigurations>\n    <DataPipelineConfiguration>\n      <DataPipelineName>my-data-pipeline</DataPipelineName>\n      <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>\n      <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>\n      <Status>Running</Status>\n      <Sources>\n          <InputBucket>my-bucket</InputBucket>\n          <InputDataScope>All</InputDataScope>\n          <IgnoreDelete>true</IgnoreDelete>\n          <FilterConfiguration>\n              <PrefixSet>prefix1/</PrefixSet>\n              <PrefixSet>prefix2/prefix3/</PrefixSet>\n              <ObjectMediaTypes>text</ObjectMediaTypes>\n              <ObjectMediaTypes>image</ObjectMediaTypes>\n              <ObjectMediaTypes>video</ObjectMediaTypes>\n          </FilterConfiguration>\n      </Sources>\n      <DataPipelineEmbeddingConfiguration>\n          <EmbeddingProvider>bailian</EmbeddingProvider>\n          <ApiKey>xxxx</ApiKey>\n          <Model>qwen2.5-vl-embedding</Model>\n          <FPS>1</FPS>\n      </DataPipelineEmbeddingConfiguration>\n      <Destination>\n          <VectorBucketName>my-vector-bucket</VectorBucketName>\n          <VectorIndexNames>my-index</VectorIndexNames>\n          <VectorKeyPrefix></VectorKeyPrefix>\n          <ObjectTagToMetadata>key1</ObjectTagToMetadata>\n          <ObjectTagToMetadata>key2</ObjectTagToMetadata>\n          <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>\n      </Destination>\n      <DataPipelineError>\n          <ErrorMode>ignoreAndRecord</ErrorMode>\n          <ErrorBucket>my-error-bucket</ErrorBucket>\n          <ErrorPrefix>error-output/</ErrorPrefix>\n      </DataPipelineError>\n      <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>\n    </DataPipelineConfiguration>\n  </DataPipelineConfigurations>\n  <NextToken>xxx</NextToken>\n</ListDataPipelineConfigurationsResult>")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8" ?>
<DataPipelineConfiguration>
  <DataPipelineName>data-pipeline</DataPipelineName>
  <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>
  <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>
  <Status>Running</Status>
  <Sources>
      <InputBucket>my-bucket</InputBucket>
      <InputDataScope>All</InputDataScope>
      <IgnoreDelete>true</IgnoreDelete>
      <FilterConfiguration>
          <PrefixSet>prefix1/</PrefixSet>
          <PrefixSet>prefix2/prefix3/</PrefixSet>
          <ObjectMediaTypes>text</ObjectMediaTypes>
          <ObjectMediaTypes>image</ObjectMediaTypes>
          <ObjectMediaTypes>video</ObjectMediaTypes>
      </FilterConfiguration>
  </Sources>
  <DataPipelineEmbeddingConfiguration>
      <EmbeddingProvider>bailian</EmbeddingProvider>
      <ApiKey>xxxx</ApiKey>
      <Model>qwen2.5-vl-embedding</Model>
      <FPS>1</FPS>
  </DataPipelineEmbeddingConfiguration>
  <Destination>
      <VectorBucketName>my-vector-bucket</VectorBucketName>
      <VectorIndexNames>my-index</VectorIndexNames>
      <VectorKeyPrefix></VectorKeyPrefix>
      <ObjectTagToMetadata>key1</ObjectTagToMetadata>
      <ObjectTagToMetadata>key2</ObjectTagToMetadata>
      <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
  </Destination>
  <DataPipelineError>
      <ErrorMode>ignoreAndRecord</ErrorMode>
      <ErrorBucket>my-error-bucket</ErrorBucket>
      <ErrorPrefix>error-output/</ErrorPrefix>
  </DataPipelineError>
  <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>
</DataPipelineConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=getDataPipelineConfiguration&dataPipeline&dataPipelineName=data-pipeline", urlStr)
		},
		&DoDataPipeLineActionRequest{
			Action: Ptr("getDataPipelineConfiguration"),
			RequestCommon: RequestCommon{
				Parameters: map[string]string{
					"dataPipelineName": "data-pipeline",
				},
			},
		},
		func(t *testing.T, o *DoDataPipeLineActionResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			getBody, err := io.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(getBody), "<?xml version=\"1.0\" encoding=\"UTF-8\" ?>\n<DataPipelineConfiguration>\n  <DataPipelineName>data-pipeline</DataPipelineName>\n  <DataPipelineDescription>使用百炼多模态模型为业务数据向量化</DataPipelineDescription>\n  <DataPipelineRole>my-data-pipeline-role</DataPipelineRole>\n  <Status>Running</Status>\n  <Sources>\n      <InputBucket>my-bucket</InputBucket>\n      <InputDataScope>All</InputDataScope>\n      <IgnoreDelete>true</IgnoreDelete>\n      <FilterConfiguration>\n          <PrefixSet>prefix1/</PrefixSet>\n          <PrefixSet>prefix2/prefix3/</PrefixSet>\n          <ObjectMediaTypes>text</ObjectMediaTypes>\n          <ObjectMediaTypes>image</ObjectMediaTypes>\n          <ObjectMediaTypes>video</ObjectMediaTypes>\n      </FilterConfiguration>\n  </Sources>\n  <DataPipelineEmbeddingConfiguration>\n      <EmbeddingProvider>bailian</EmbeddingProvider>\n      <ApiKey>xxxx</ApiKey>\n      <Model>qwen2.5-vl-embedding</Model>\n      <FPS>1</FPS>\n  </DataPipelineEmbeddingConfiguration>\n  <Destination>\n      <VectorBucketName>my-vector-bucket</VectorBucketName>\n      <VectorIndexNames>my-index</VectorIndexNames>\n      <VectorKeyPrefix></VectorKeyPrefix>\n      <ObjectTagToMetadata>key1</ObjectTagToMetadata>\n      <ObjectTagToMetadata>key2</ObjectTagToMetadata>\n      <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>\n  </Destination>\n  <DataPipelineError>\n      <ErrorMode>ignoreAndRecord</ErrorMode>\n      <ErrorBucket>my-error-bucket</ErrorBucket>\n      <ErrorPrefix>error-output/</ErrorPrefix>\n  </DataPipelineError>\n  <CreateTime>2021-06-29T14:50:13.011643661+08:00</CreateTime>\n</DataPipelineConfiguration>")
		},
	},
}

func TestMockDoDataPipeLineAction_Success(t *testing.T) {
	for _, c := range testMockDoDataPipeLineActionSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoDataPipeLineAction(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDoDataPipeLineActionErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DoDataPipeLineActionRequest
	CheckOutputFn  func(t *testing.T, o *DoDataPipeLineActionResult, err error)
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", urlStr)
		},
		&DoDataPipeLineActionRequest{
			Action: Ptr("listDataPipelineConfigurations"),
		},
		func(t *testing.T, o *DoDataPipeLineActionResult, err error) {
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
			assert.Equal(t, "POST", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/?action=listDataPipelineConfigurations&dataPipeline", urlStr)
		},
		&DoDataPipeLineActionRequest{
			Action: Ptr("listDataPipelineConfigurations"),
		},
		func(t *testing.T, o *DoDataPipeLineActionResult, err error) {
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

func TestMockDoDataPipeLineAction_Error(t *testing.T) {
	for _, c := range testMockDoDataPipeLineActionErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DoDataPipeLineAction(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

