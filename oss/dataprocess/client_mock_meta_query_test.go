package dataprocess

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
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
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic&role=my-role", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><Filters><Filter>Size &gt; 0</Filter><Filter>ContentType = &#39;image/jpeg&#39;</Filter></Filters></MetaQuery>")
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			Role:   oss.Ptr("my-role"),
			MetaQuery: &OpenMetaQuery{
				Filters: &Filters{
					Filter: []string{
						"Size > 0",
						"ContentType = 'image/jpeg'",
					},
				},
			},
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
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=semantic&role=AliyunMetaQueryDefaultRole", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><WorkflowParameters><WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter><WorkflowParameter><Name>VideoInsightEnable</Name><Value>True</Value></WorkflowParameter></WorkflowParameters><NotificationAttributes><Notifications><Notification><MNS>imm-index-notification</MNS></Notification></Notifications><WithFields><WithField>Insights</WithField><WithField>Labels</WithField></WithFields></NotificationAttributes><DatasetConfig><Insights><Language>en</Language></Insights></DatasetConfig><IndexOptions><IgnoreObjectDelete>true</IgnoreObjectDelete></IndexOptions><RouteRule><Type>OSSTag</Type><AutoCreateDataset>true</AutoCreateDataset><OSSTagKey>routing-dataset</OSSTagKey></RouteRule></MetaQuery>")
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Role:   oss.Ptr("AliyunMetaQueryDefaultRole"),
			Mode:   oss.Ptr("semantic"),
			MetaQuery: &OpenMetaQuery{
				WorkflowParameters: &WorkflowParameters{
					WorkflowParameter: []WorkflowParameter{
						{
							Name:  oss.Ptr("ImageInsightEnable"),
							Value: oss.Ptr("True"),
						},
						{
							Name:  oss.Ptr("VideoInsightEnable"),
							Value: oss.Ptr("True"),
						},
					},
				},
				NotificationAttributes: &NotificationAttributes{
					Notifications: &Notifications{
						Notification: []Notification{
							{
								MNS: oss.Ptr("imm-index-notification"),
							},
						},
					},
					WithFields: &WithFields{
						[]string{
							"Insights",
							"Labels",
						},
					},
				},
				IndexOptions: &IndexOptions{
					IgnoreObjectDelete: oss.Ptr(true),
				},
				RouteRule: &RouteRule{
					Type:              oss.Ptr("OSSTag"),
					AutoCreateDataset: oss.Ptr(true),
					OSSTagKey:         oss.Ptr("routing-dataset"),
				},
				DatasetConfig: &DatasetConfig{
					Insights: &InsightsConfig{
						Language: oss.Ptr("en"),
					},
				},
			},
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

		cfg := oss.LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=openMetaQuery&metaQuery&mode=basic", strUrl)
		},
		&OpenMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
		},
		func(t *testing.T, o *OpenMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

		cfg := oss.LoadDefaultConfig().
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
  <Phase>IncrementalScanning</Phase>
  <CreateTime>2026-05-20T08:00:00.000+08:00</CreateTime>
  <UpdateTime>2026-05-20T08:30:00.000+08:00</UpdateTime>
  <MetaQueryMode>semantic</MetaQueryMode>
  <WorkflowParameters>
    <WorkflowParameter>
      <Name>ImageInsightEnable</Name>
      <Value>True</Value>
    </WorkflowParameter>
  </WorkflowParameters>
  <Filters>
    <Filter>Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z</Filter>
    <Filter>Filename prefix (YWEvYmIv)</Filter>
  </Filters>
  <IndexOptions>
    <IgnoreObjectDelete>True</IgnoreObjectDelete>
  </IndexOptions>
  <RouteRule>
    <Type>OSSTag</Type>
    <AutoCreateDataset>True</AutoCreateDataset>
    <OSSTagKey>routing-dataset</OSSTagKey>
  </RouteRule>
  <NotificationAttributes>
    <Notifications>
      <Notification>
        <MNS>imm-index-notification</MNS>
      </Notification>
    </Notifications>
    <WithFields>
      <WithField>Insights</WithField>
    </WithFields>
  </NotificationAttributes>
  <DatasetConfig>
    <Insights>
      <Language>en</Language>
    </Insights>
  </DatasetConfig>
</MetaQueryStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.State, "Running")
			assert.Equal(t, *o.Phase, "IncrementalScanning")
			assert.Equal(t, *o.CreateTime, "2026-05-20T08:00:00.000+08:00")
			assert.Equal(t, *o.UpdateTime, "2026-05-20T08:30:00.000+08:00")
			assert.Equal(t, *o.MetaQueryMode, "semantic")
			assert.Equal(t, *o.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, o.Filters.Filter[0], "Size > 1024,FileModifiedTime > 2025-06-03T09:20:47.999Z")
			assert.Equal(t, o.Filters.Filter[1], "Filename prefix (YWEvYmIv)")
			assert.Equal(t, *o.IndexOptions.IgnoreObjectDelete, true)
			assert.Equal(t, *o.RouteRule.Type, "OSSTag")
			assert.Equal(t, *o.RouteRule.AutoCreateDataset, true)
			assert.Equal(t, *o.RouteRule.OSSTagKey, "routing-dataset")
			assert.Equal(t, *o.NotificationAttributes.Notifications.Notification[0].MNS, "imm-index-notification")
			assert.Equal(t, o.NotificationAttributes.WithFields.WithField[0], "Insights")
			assert.Equal(t, *o.DatasetConfig.Insights.Language, "en")

		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<MetaQueryStatus>
    <State>Running</State>
    <Phase>IncrementalScanning</Phase>
    <CreateTime>2026-06-05T10:49:57.109295646+08:00</CreateTime>
    <UpdateTime>2026-06-05T10:49:59.411578836+08:00</UpdateTime>
    <MetaQueryMode>basic</MetaQueryMode>
    <RouteRule>
        <Type>default</Type>
        <AutoCreateDataset>True</AutoCreateDataset>
    </RouteRule>
    <WorkflowParameters>
    </WorkflowParameters>
    <Filters>
    </Filters>
    <DatasetConfig>
        <ReverseImage>
            <Video>
                <Enable>false</Enable>
            </Video>
            <Image>
                <Enable>false</Enable>
            </Image>
        </ReverseImage>
        <Insights>
            <Language>zh-Hans</Language>
            <Image>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                </Caption>
            </Image>
            <Video>
                <Caption>
                    <Enable>false</Enable>
                    <Prompt>
                    </Prompt>
                    <PersonReference>
                        <Enable>false</Enable>
                    </PersonReference>
                </Caption>
                <Label>
                    <System>
                        <Enable>false</Enable>
                    </System>
                    <UserDefined>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </UserDefined>
                    <Highlight>
                        <Enable>false</Enable>
                        <Labels>
                        </Labels>
                    </Highlight>
                </Label>
                <MultiStream>
                    <Enable>false</Enable>
                </MultiStream>
            </Video>
        </Insights>
        <SmartCluster>
            <Figure>
                <AutoGenerate>false</AutoGenerate>
                <AutoClustering>false</AutoClustering>
                <MinEntityCount>3</MinEntityCount>
                <EnabledFeatures>face</EnabledFeatures>
            </Figure>
        </SmartCluster>
    </DatasetConfig>
</MetaQueryStatus>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.State, "Running")
			assert.Equal(t, *o.Phase, "IncrementalScanning")
			assert.Equal(t, *o.CreateTime, "2026-06-05T10:49:57.109295646+08:00")
			assert.Equal(t, *o.UpdateTime, "2026-06-05T10:49:59.411578836+08:00")
			assert.Equal(t, *o.MetaQueryMode, "basic")
			assert.Equal(t, *o.RouteRule.Type, "default")
			assert.Equal(t, *o.RouteRule.AutoCreateDataset, true)
			assert.Equal(t, *o.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Language, "zh-Hans")
			assert.Equal(t, *o.DatasetConfig.Insights.Image.Caption.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Caption.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.System.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
			assert.Equal(t, *o.DatasetConfig.Insights.Video.MultiStream.Enable, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
			assert.Equal(t, *o.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
			assert.Equal(t, o.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")
		},
	},
}

func TestMockGetMetaQueryStatus_Success(t *testing.T) {
	for _, c := range testMockGetMetaQueryStatusSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
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
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=getMetaQueryStatus&metaQuery", strUrl)
		},
		&GetMetaQueryStatusRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *GetMetaQueryStatusResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

		cfg := oss.LoadDefaultConfig().
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
		[]byte(`<?xml version=\"1.0\" encoding=\"UTF-8\"?>
<MetaQuery>
    <NextToken>next-page-token-abc</NextToken>
    <TotalHits>123</TotalHits>
    <Files>
        <File>
            <Filename>photos/sunset.jpg</Filename>
            <Size>2097152</Size>
            <FileModifiedTime>2026-05-19T15:30:00.000+08:00</FileModifiedTime>
            <ContentType>image/jpeg</ContentType>
            <ObjectACL>default</ObjectACL>
            <OSSStorageClass>Standard</OSSStorageClass>
        </File>
        <File>
            <Filename>photos/mountain.png</Filename>
            <Size>5242880</Size>
        </File>
    </Files>
    <Aggregations>
        <Aggregation>
            <Field>Size</Field>
            <Operation>sum</Operation>
            <Value>12345678</Value>
        </Aggregation>
    </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.NextToken, "next-page-token-abc")
			assert.Equal(t, *o.TotalHits, int64(123))
			assert.Equal(t, *o.Files[0].Filename, "photos/sunset.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(2097152))
			assert.Equal(t, *o.Files[0].FileModifiedTime, "2026-05-19T15:30:00.000+08:00")
			assert.Equal(t, *o.Files[0].ContentType, "image/jpeg")
			assert.Equal(t, *o.Files[0].ObjectACL, "default")
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[1].Filename, "photos/mountain.png")
			assert.Equal(t, *o.Files[1].Size, int64(5242880))
			assert.Equal(t, *o.Aggregations[0].Field, "Size")
			assert.Equal(t, *o.Aggregations[0].Operation, "sum")
			assert.Equal(t, *o.Aggregations[0].Value, "12345678")
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
    <TotalHits>50</TotalHits>
    <Aggregations>
        <Aggregation>
            <Field>StorageClass</Field>
            <Operation>group_by</Operation>
            <Groups>
                <Group>
                    <Value>Standard</Value>
                    <Count>30</Count>
                </Group>
                <Group>
                    <Value>IA</Value>
                    <Count>20</Count>
                </Group>
            </Groups>
        </Aggregation>
    </Aggregations>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.TotalHits, int64(50))
			assert.Equal(t, len(o.Aggregations), 1)
			assert.Equal(t, *o.Aggregations[0].Field, "StorageClass")
			assert.Equal(t, *o.Aggregations[0].Operation, "group_by")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Value, "Standard")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[0].Count, int64(30))
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Value, "IA")
			assert.Equal(t, *o.Aggregations[0].AggregationGroups[1].Count, int64(20))
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
    <TotalHits>2</TotalHits>
    <Files>
        <File>
            <Filename>photos/cat-in-living-room.jpg</Filename>
            <Size>3145728</Size>
            <OSSStorageClass>Standard</OSSStorageClass>
            <Labels>
                <Label>
                    <LabelName>cat</LabelName>
                    <LabelConfidence>0.98</LabelConfidence>
                </Label>
            </Labels>
        </File>
        <File>
            <Filename>photos/kitten-sofa.jpg</Filename>
            <Size>2621440</Size>
        </File>
    </Files>
</MetaQuery>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.TotalHits, int64(2))
			assert.Equal(t, *o.Files[0].Filename, "photos/cat-in-living-room.jpg")
			assert.Equal(t, *o.Files[0].Size, int64(3145728))
			assert.Equal(t, *o.Files[0].OSSStorageClass, "Standard")
			assert.Equal(t, *o.Files[0].Labels[0].LabelName, "cat")
			assert.Equal(t, *o.Files[0].Labels[0].LabelConfidence, float64(0.98))
			assert.Equal(t, *o.Files[1].Filename, "photos/kitten-sofa.jpg")
			assert.Equal(t, *o.Files[1].Size, int64(2621440))
		},
	},
}

func TestMockDoMetaQuery_Success(t *testing.T) {
	for _, c := range testMockDoMetaQuerySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=doMetaQuery&metaQuery&mode=basic", strUrl)
			body, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(body), "<MetaQuery><MaxResults>100</MaxResults><Query>{&#34;Field&#34;:&#34;Size&#34;,&#34;Operation&#34;:&#34;gt&#34;,&#34;Value&#34;:&#34;1048576&#34;}</Query><Sort>Size</Sort><Order>desc</Order><Aggregations><Aggregation><Operation>sum</Operation><Field>Size</Field><Groups></Groups></Aggregation></Aggregations></MetaQuery>")
		},
		&DoMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
			Mode:   oss.Ptr("basic"),
			MetaQuery: &DoMetaQuery{
				Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
				Sort:  oss.Ptr("Size"),
				Order: oss.Ptr(MetaQueryOrderDesc),
				Aggregations: &MetaQueryAggregations{
					[]Aggregation{
						{
							Field:     oss.Ptr("Size"),
							Operation: oss.Ptr("sum"),
						},
					},
				},
				MaxResults: oss.Ptr(int64(100)),
			},
		},
		func(t *testing.T, o *DoMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

		cfg := oss.LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
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

		cfg := oss.LoadDefaultConfig().
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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
			assert.Equal(t, "/bucket/?action=closeMetaQuery&metaQuery", strUrl)
		},
		&CloseMetaQueryRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *CloseMetaQueryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *oss.ServiceError
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

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CloseMetaQuery(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
