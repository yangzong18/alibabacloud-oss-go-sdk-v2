package dataprocess

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockCreateDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *CreateDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh</Language></Insights></DatasetConfig>
</Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=your_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
    <Dataset>
        <DatasetName>test_dataset</DatasetName>
        <WorkflowParameters>
        </WorkflowParameters>
        <CreateTime>2026-06-08T17:36:59.774068044+08:00</CreateTime>
        <UpdateTime>2026-06-08T17:36:59.774068044+08:00</UpdateTime>
        <Description>this is a demo</Description>
        <DatasetMaxFileCount>100000000</DatasetMaxFileCount>
        <DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
        <DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
        <DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
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
                <Language>zh</Language>
                <Image>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
                    </Caption>
                </Image>
                <Video>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
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
    </Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&datasetConfig=%7B%0A%09%09%09%09++%22Insights%22%3A+%7B%0A%09%09%09%09%09%22Language%22%3A+%22zh%22%0A%09%09%09%09++%7D%0A%09%09%09%09%7D&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
			DatasetConfig: oss.Ptr(`{
				  "Insights": {
					"Language": "zh"
				  }
				}`),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Dataset.DatasetName, "test_dataset")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Prompt, "")
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh")

			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
			assert.Equal(t, o.Dataset.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<CreateDatasetResponse>
    <Dataset>
        <DatasetName>test_dataset</DatasetName>
        <WorkflowParameters>
        </WorkflowParameters>
        <CreateTime>2026-06-08T17:36:59.774068044+08:00</CreateTime>
        <UpdateTime>2026-06-08T17:36:59.774068044+08:00</UpdateTime>
        <Description>this is a demo</Description>
        <DatasetMaxFileCount>100000000</DatasetMaxFileCount>
        <DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
        <DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
        <DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
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
                <Language>zh</Language>
                <Image>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
                    </Caption>
                </Image>
                <Video>
                    <Caption>
                        <Enable>false</Enable>
                        <Prompt></Prompt>
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
    </Dataset>
</CreateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createDataset&datasetConfig=%7B%22Insights%22%3A%7B%22Language%22%3A%22zh%22%7D%7D&datasetName=test_dataset&description=this+is+a+demo&metaQuery&workflowParameters=%5B%7B%22Name%22%3A%22VideoInsightEnable%22%2C%22Value%22%3A%22True%22%7D%2C%7B%22Name%22%3A%22ImageInsightEnable%22%2C%22Value%22%3A%22True%22%7D%5D", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
			WorkflowParameters: oss.Ptr(WorkflowParameters{
				WorkflowParameter: []WorkflowParameter{
					{
						Name:  oss.Ptr("VideoInsightEnable"),
						Value: oss.Ptr("True"),
					},
					{
						Name:  oss.Ptr("ImageInsightEnable"),
						Value: oss.Ptr("True"),
					},
				},
			}.ToParameterValue()),
			DatasetConfig: oss.Ptr((&DatasetConfig{
				Insights: &InsightsConfig{
					Language: oss.Ptr("zh"),
				},
			}).ToParameterValue()),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Dataset.DatasetName, "test_dataset")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-06-08T17:36:59.774068044+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Video.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.ReverseImage.Image.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Prompt, "")
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.PersonReference.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh")

			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoGenerate, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.AutoClustering, false)
			assert.Equal(t, *o.Dataset.DatasetConfig.SmartCluster.Figure.MinEntityCount, int64(3))
			assert.Equal(t, o.Dataset.DatasetConfig.SmartCluster.Figure.EnabledFeatures[0], "face")
		},
	},
}

func TestMockCreateDataset_Success(t *testing.T) {
	for _, c := range testMockCreateDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *CreateDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
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
			assert.Equal(t, "/bucket/?action=createDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&CreateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *CreateDatasetResult, err error) {
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

func TestMockCreateDataset_Error(t *testing.T) {
	for _, c := range testMockCreateDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDatasetRequest
	CheckOutputFn  func(t *testing.T, o *GetDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<GetDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
<FileCount>3456</FileCount>
<TotalFileSize>10737418240</TotalFileSize>
</Dataset>
</GetDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=your_dataset&metaQuery&withStatistics=true", strUrl)
		},
		&GetDatasetRequest{
			Bucket:         oss.Ptr("bucket"),
			DatasetName:    oss.Ptr("your_dataset"),
			WithStatistics: oss.Ptr(true),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh-Hans")
			assert.Equal(t, *o.Dataset.FileCount, int64(3456))
			assert.Equal(t, *o.Dataset.TotalFileSize, int64(10737418240))

		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<GetDatasetResponse>
<Dataset>
<DatasetName>photos-2026</DatasetName>
<WorkflowParameters></WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig><Insights><Language>zh-Hans</Language></Insights></DatasetConfig>
</Dataset>
</GetDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Dataset.DatasetName, "photos-2026")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.Equal(t, *o.Dataset.DatasetConfig.Insights.Language, "zh-Hans")
		},
	},
}

func TestMockGetDataset_Success(t *testing.T) {
	for _, c := range testMockGetDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetDatasetRequest
	CheckOutputFn  func(t *testing.T, o *GetDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
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
			assert.Equal(t, "/bucket/?action=getDataset&datasetName=test_dataset&metaQuery", strUrl)
		},
		&GetDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
		},
		func(t *testing.T, o *GetDatasetResult, err error) {
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

func TestMockGetDataset_Error(t *testing.T) {
	for _, c := range testMockGetDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *UpdateDatasetResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<UpdateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig>
<Insights>
	<Language>zh</Language>
	<Image>
		<Caption>
			<Enable>true</Enable>
		</Caption>
	</Image>
	<Video>
		<Caption>
			<Enable>true</Enable>
		</Caption>
		<Label>
			<System>
				<Enable>true</Enable>
			</System>
			<UserDefined>
				<Enable>true</Enable>
			</UserDefined>
			<Highlight>
				<Enable>true</Enable>
			</Highlight>
		</Label>
		<MultiStream>
			<Enable>true</Enable>
		</MultiStream>
	</Video>
</Insights>
</DatasetConfig>
</Dataset>
</UpdateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=updateDataset&datasetConfig=%7B%22Insights%22%3A%7B%22Language%22%3A%22zh%22%2C%22Image%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%2C%22Video%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22Label%22%3A%7B%22System%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22UserDefined%22%3A%7B%22Enable%22%3A%22true%22%7D%2C%22Highlight%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%2C%22MultiStream%22%3A%7B%22Enable%22%3A%22true%22%7D%7D%7D%7D&datasetName=your_dataset&description=this+is+a+demo&metaQuery&workflowParameters=%5B%7B%22Name%22%3A+%22ImageInsightEnable%22%2C%22Value%22%3A+%22True%22%2C%22Description%22%3A+%22The+source+bucket+for+data+processing%22%7D%5D", strUrl)

		},
		&UpdateDatasetRequest{
			Bucket:             oss.Ptr("bucket"),
			DatasetName:        oss.Ptr("your_dataset"),
			Description:        oss.Ptr("this is a demo"),
			WorkflowParameters: oss.Ptr(`[{"Name": "ImageInsightEnable","Value": "True","Description": "The source bucket for data processing"}]`),
			DatasetConfig:      oss.Ptr(`{"Insights":{"Language":"zh","Image":{"Caption":{"Enable":"true"}},"Video":{"Caption":{"Enable":"true"},"Label":{"System":{"Enable":"true"},"UserDefined":{"Enable":"true"},"Highlight":{"Enable":"true"}},"MultiStream":{"Enable":"true"}}}}`),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<UpdateDatasetResponse>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<WorkflowParameters>
      <WorkflowParameter><Name>ImageInsightEnable</Name><Value>True</Value></WorkflowParameter>
</WorkflowParameters>
<CreateTime>2026-04-21T18:17:58.727923181+08:00</CreateTime>
<UpdateTime>2026-04-21T18:17:58.727923181+08:00</UpdateTime>
<Description>this is a demo</Description>
<DatasetMaxFileCount>100000000</DatasetMaxFileCount>
<DatasetMaxEntityCount>10000000000</DatasetMaxEntityCount>
<DatasetMaxRelationCount>100000000000</DatasetMaxRelationCount>
<DatasetMaxTotalFileSize>90000000000000000</DatasetMaxTotalFileSize>
<DatasetConfig>
<Insights>
	<Language>zh</Language>
	<Image>
		<Caption>
			<Enable>true</Enable>
		</Caption>
	</Image>
	<Video>
		<Caption>
			<Enable>true</Enable>
		</Caption>
		<Label>
			<System>
				<Enable>true</Enable>
			</System>
			<UserDefined>
				<Enable>true</Enable>
			</UserDefined>
			<Highlight>
				<Enable>true</Enable>
			</Highlight>
		</Label>
		<MultiStream>
			<Enable>true</Enable>
		</MultiStream>
	</Video>
</Insights>
</DatasetConfig>
</Dataset>
</UpdateDatasetResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=updateDataset&datasetConfig=%7B%22Insights%22%3A%7B%22Language%22%3A%22zh%22%2C%22Image%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3Atrue%7D%7D%2C%22Video%22%3A%7B%22Caption%22%3A%7B%22Enable%22%3Atrue%7D%2C%22Label%22%3A%7B%22System%22%3A%7B%22Enable%22%3Atrue%7D%2C%22UserDefined%22%3A%7B%22Enable%22%3Atrue%2C%22Labels%22%3A%5B%7B%22Name%22%3A%22%E6%9C%89%E4%BA%BA%E6%91%94%E5%80%92%22%2C%22Description%22%3A%22%E7%94%BB%E9%9D%A2%E4%B8%AD%E6%9C%89%E4%BA%BA%E7%94%B1%E7%AB%99%E7%AB%8B%E6%88%96%E8%A1%8C%E8%B5%B0%E5%8F%98%E4%B8%BA%E5%80%92%E5%9C%A8%E5%9C%B0%E9%9D%A2%E7%9A%84%E5%8A%A8%E4%BD%9C%22%7D%5D%7D%7D%2C%22MultiStream%22%3A%7B%22Enable%22%3Atrue%7D%7D%7D%7D&datasetName=your_dataset&description=this+is+a+demo&metaQuery&workflowParameters=%5B%7B%22Name%22%3A%22ImageInsightEnable%22%2C%22Value%22%3A%22True%22%2C%22Description%22%3A%22The+source+bucket+for+data+processing%22%7D%5D", strUrl)

		},
		&UpdateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			Description: oss.Ptr("this is a demo"),
			WorkflowParameters: oss.Ptr(WorkflowParameters{
				WorkflowParameter: []WorkflowParameter{
					{
						Name:        oss.Ptr("ImageInsightEnable"),
						Value:       oss.Ptr("True"),
						Description: oss.Ptr("The source bucket for data processing"),
					},
				},
			}.ToParameterValue()),
			DatasetConfig: oss.Ptr((&DatasetConfig{
				Insights: &InsightsConfig{
					Language: oss.Ptr("zh"),
					Image: &InsightsImage{
						Caption: &InsightsImageCaption{
							Enable: oss.Ptr(true),
						},
					},
					Video: &InsightsVideo{
						Caption: &InsightsVideoCaption{
							Enable: oss.Ptr(true),
						},
						Label: &InsightsVideoLabel{
							System: &InsightsVideoSystem{
								Enable: oss.Ptr(true),
							},
							UserDefined: &InsightsVideoUserDefined{
								Enable: oss.Ptr(true),
								Labels: []LabelItem{
									{
										Name:        oss.Ptr("有人摔倒"),
										Description: oss.Ptr("画面中有人由站立或行走变为倒在地面的动作"),
									},
								},
							},
						},
						MultiStream: &InsightsVideoMultiStream{
							Enable: oss.Ptr(true),
						},
					},
				},
			}).ToParameterValue()),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Dataset.DatasetName, "test-dataset")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Name, "ImageInsightEnable")
			assert.Equal(t, *o.Dataset.WorkflowParameters.WorkflowParameter[0].Value, "True")
			assert.Equal(t, *o.Dataset.CreateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.UpdateTime, "2026-04-21T18:17:58.727923181+08:00")
			assert.Equal(t, *o.Dataset.Description, "this is a demo")
			assert.Equal(t, *o.Dataset.DatasetMaxFileCount, int64(100000000))
			assert.Equal(t, *o.Dataset.DatasetMaxEntityCount, int64(10000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxRelationCount, int64(100000000000))
			assert.Equal(t, *o.Dataset.DatasetMaxTotalFileSize, int64(90000000000000000))
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Image.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Caption.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.System.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.UserDefined.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.Label.Highlight.Enable)
			assert.True(t, *o.Dataset.DatasetConfig.Insights.Video.MultiStream.Enable)
		},
	},
}

func TestMockUpdateDataset_Success(t *testing.T) {
	for _, c := range testMockUpdateDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateDatasetRequest
	CheckOutputFn  func(t *testing.T, o *UpdateDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=updateDataset&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&UpdateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
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
			assert.Equal(t, "/bucket/?action=updateDataset&datasetName=test_dataset&description=this+is+a+demo&metaQuery", strUrl)
		},
		&UpdateDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test_dataset"),
			Description: oss.Ptr("this is a demo"),
		},
		func(t *testing.T, o *UpdateDatasetResult, err error) {
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

func TestMockUpdateDataset_Error(t *testing.T) {
	for _, c := range testMockUpdateDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDatasetsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDatasetsRequest
	CheckOutputFn  func(t *testing.T, o *ListDatasetsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListDatasetsResponse>
<NextToken>1986505809429276:oss_1234567890_demo-bucket:test-dataset</NextToken>
<Datasets>
<Dataset>
<DatasetName>test-dataset</DatasetName>
<CreateTime>2026-04-22T11:39:28.148283473+08:00</CreateTime>
<UpdateTime>2026-04-22T11:39:28.148283473+08:00</UpdateTime>
</Dataset>
</Datasets>
</ListDatasetsResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, len(o.Datasets), 1)
			assert.Equal(t, *o.NextToken, "1986505809429276:oss_1234567890_demo-bucket:test-dataset")
			assert.Equal(t, *o.Datasets[0].DatasetName, "test-dataset")
			assert.Equal(t, *o.Datasets[0].CreateTime, "2026-04-22T11:39:28.148283473+08:00")
			assert.Equal(t, *o.Datasets[0].UpdateTime, "2026-04-22T11:39:28.148283473+08:00")

		},
	},
}

func TestMockListDatasets_Success(t *testing.T) {
	for _, c := range testMockListDatasetsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDatasets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListDatasetsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListDatasetsRequest
	CheckOutputFn  func(t *testing.T, o *ListDatasetsResult, err error)
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
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
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
			assert.Equal(t, "/bucket/?action=listDatasets&metaQuery", strUrl)
		},
		&ListDatasetsRequest{
			Bucket: oss.Ptr("bucket"),
		},
		func(t *testing.T, o *ListDatasetsResult, err error) {
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

func TestMockListDatasets_Error(t *testing.T) {
	for _, c := range testMockListDatasetsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListDatasets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDatasetSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDatasetRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteDataset_Success(t *testing.T) {
	for _, c := range testMockDeleteDatasetSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteDatasetErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteDatasetRequest
	CheckOutputFn  func(t *testing.T, o *DeleteDatasetResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
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
			assert.Equal(t, "/bucket/?action=deleteDataset&datasetName=test-dataset&metaQuery", strUrl)
		},
		&DeleteDatasetRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
		},
		func(t *testing.T, o *DeleteDatasetResult, err error) {
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

func TestMockDeleteDataset_Error(t *testing.T) {
	for _, c := range testMockDeleteDatasetErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteDataset(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
