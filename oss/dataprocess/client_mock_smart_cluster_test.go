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

var testMockCreateSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *CreateSmartClusterResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateSmartClusterResponse>
  <ObjectId>cluster-abc123def456</ObjectId>
</CreateSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectId, "cluster-abc123def456")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CreateSmartClusterResponse>
  <ObjectId>cluster-abc123def456</ObjectId>
</CreateSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&notification=%7B%22MNS%22%3A%7B%22TopicName%22%3A%22imm-cluster-notification%22%7D%7D&rules=%5B%7B%22RuleType%22%3A%22keywords%22%2C%22Keywords%22%3A%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules: oss.Ptr((SmartClusterRules{
				Rules: []SmartClusterRule{
					{
						RuleType: oss.Ptr("keywords"),
						Keywords: []string{"car"},
					},
				},
			}).ToParameterValue()),
			Notification: oss.Ptr(SmartClusterNotification{
				MNS: &SmartClusterTopicName{
					TopicName: oss.Ptr("imm-cluster-notification"),
				},
			}.ToParameterValue()),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectId, "cluster-abc123def456")
		},
	},
}

func TestMockCreateSmartCluster_Success(t *testing.T) {
	for _, c := range testMockCreateSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCreateSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CreateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *CreateSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
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
			assert.Equal(t, "/bucket/?action=createSmartCluster&clusterType=knowledge&datasetName=test-dataset&metaQuery&name=your_name&rules=%5B%7B%22RuleType%22%3A+%22keywords%22%2C%22Keywords%22%3A+%5B%22car%22%5D%7D%5D", strUrl)
		},
		&CreateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			Name:        oss.Ptr("your_name"),
			ClusterType: SmartClusterTypeKnowledge,
			Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		},
		func(t *testing.T, o *CreateSmartClusterResult, err error) {
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

func TestMockCreateSmartCluster_Error(t *testing.T) {
	for _, c := range testMockCreateSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CreateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *GetSmartClusterResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<GetSmartClusterResponse>
  <SmartCluster>
    <ObjectId>cluster-abc123def456</ObjectId>
    <ClusterType>figure</ClusterType>
    <Name>face-cluster-alice</Name>
    <Description>this is a demo</Description>
    <Rules>
      <Rule>
        <RuleType>face</RuleType>
        <BaseURIs>oss://examplebucket/refs/alice.jpg</BaseURIs>
        <Sensitivity>0.7</Sensitivity>
      </Rule>
    </Rules>
    <Reason></Reason>
    <Notification>
      <MNS><TopicName>imm-cluster-notification</TopicName></MNS>
    </Notification>
    <CreateTime>2026-05-20T11:00:00.000+08:00</CreateTime>
    <UpdateTime>2026-05-20T11:08:00.000+08:00</UpdateTime>
  </SmartCluster>
</GetSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.SmartCluster.ObjectId, "cluster-abc123def456")
			assert.Equal(t, *o.SmartCluster.ClusterType, "figure")
			assert.Equal(t, *o.SmartCluster.CreateTime, "2026-05-20T11:00:00.000+08:00")
			assert.Equal(t, *o.SmartCluster.UpdateTime, "2026-05-20T11:08:00.000+08:00")
			assert.Equal(t, *o.SmartCluster.Name, "face-cluster-alice")
			assert.Equal(t, *o.SmartCluster.Description, "this is a demo")
			assert.Equal(t, *o.SmartCluster.Rules[0].RuleType, "face")
			assert.Equal(t, o.SmartCluster.Rules[0].BaseURIs[0], "oss://examplebucket/refs/alice.jpg")
			assert.Equal(t, *o.SmartCluster.Rules[0].Sensitivity, 0.7)
			assert.Equal(t, *o.SmartCluster.Notification.MNS.TopicName, "imm-cluster-notification")
		},
	},
}

func TestMockGetSmartCluster_Success(t *testing.T) {
	for _, c := range testMockGetSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *GetSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
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
			assert.Equal(t, "/bucket/?action=getSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&GetSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *GetSmartClusterResult, err error) {
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

func TestMockGetSmartCluster_Error(t *testing.T) {
	for _, c := range testMockGetSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *UpdateSmartClusterResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><UpdateSmartClusterResponse><ObjectId>cluster-abc123def456</ObjectId></UpdateSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "cluster-abc123def456", *o.ObjectId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?><UpdateSmartClusterResponse><ObjectId>cluster-abc123def456</ObjectId></UpdateSmartClusterResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&notification=%7B%22MNS%22%3A%7B%22TopicName%22%3A%22imm-cluster-notification%22%7D%7D&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules: oss.Ptr(SmartClusterRules{Rules: []SmartClusterRule{
				{
					RuleType:    oss.Ptr("face"),
					Sensitivity: oss.Ptr(0.7),
				},
			}}.ToParameterValue()),
			Notification: oss.Ptr(SmartClusterNotification{
				MNS: &SmartClusterTopicName{
					TopicName: oss.Ptr("imm-cluster-notification"),
				},
			}.ToParameterValue()),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "cluster-abc123def456", *o.ObjectId)
		},
	},
}

func TestMockUpdateSmartCluster_Success(t *testing.T) {
	for _, c := range testMockUpdateSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUpdateSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UpdateSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *UpdateSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
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
			assert.Equal(t, "/bucket/?action=updateSmartCluster&datasetName=your_dataset&description=this+is+a+demo&metaQuery&name=face-cluster-alice&objectId=cluster-abc123def456&rules=%5B%7B%22RuleType%22%3A%22face%22%2C%22Sensitivity%22%3A0.7%7D%5D", strUrl)
		},
		&UpdateSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
			Description: oss.Ptr("this is a demo"),
			Name:        oss.Ptr("face-cluster-alice"),
			Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
		},
		func(t *testing.T, o *UpdateSmartClusterResult, err error) {
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

func TestMockUpdateSmartCluster_Error(t *testing.T) {
	for _, c := range testMockUpdateSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UpdateSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListSmartClustersSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListSmartClustersRequest
	CheckOutputFn  func(t *testing.T, o *ListSmartClustersResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListSmartClustersResponse>
    <SmartClusters>
        <SmartCluster>
            <CreateTime>2026-06-10T09:54:27.484585901+08:00</CreateTime>
            <ObjectId>SmartCluster-cb9f8c95-281f-490b-b677-eca2f6ff0c19</ObjectId>
            <UpdateTime>2026-06-10T09:54:27.484585901+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>demo-2</Name>
            <Rules>
                <Rule>
                    <Keywords>cat</Keywords>
                    <Sensitivity>0.5</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T16:18:20.998039823+08:00</CreateTime>
            <ObjectId>SmartCluster-c30d039b-0b55-4a42-ae19-cae00a53735a</ObjectId>
            <UpdateTime>2026-06-09T18:08:12.90394089+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>new-demo</Name>
            <Description>this is a demo</Description>
            <Rules>
                <Rule>
                    <Keywords>hello</Keywords>
                    <Keywords>world</Keywords>
                    <Sensitivity>0.7</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T15:21:41.288773397+08:00</CreateTime>
            <ObjectId>SmartCluster-fed69d7a-683a-4452-9d70-54148bd458e9</ObjectId>
            <UpdateTime>2026-06-09T15:21:41.288773397+08:00</UpdateTime>
            <ClusterType>knowledge</ClusterType>
            <Name>demo</Name>
            <Rules>
                <Rule>
                    <Keywords>hello</Keywords>
                    <Keywords>world</Keywords>
                    <Sensitivity>0.5</Sensitivity>
                    <RuleType>keywords</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
    </SmartClusters>
</ListSmartClustersResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.SmartClusters), 3)
			assert.Equal(t, *o.SmartClusters[0].Name, "demo-2")
			assert.Equal(t, *o.SmartClusters[0].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[0].CreateTime, "2026-06-10T09:54:27.484585901+08:00")
			assert.Equal(t, *o.SmartClusters[0].UpdateTime, "2026-06-10T09:54:27.484585901+08:00")
			assert.Equal(t, *o.SmartClusters[0].ObjectId, "SmartCluster-cb9f8c95-281f-490b-b677-eca2f6ff0c19")
			assert.Equal(t, *o.SmartClusters[0].Reason, "")
			assert.Equal(t, o.SmartClusters[0].Rules[0].Keywords[0], "cat")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].Sensitivity, float64(0.5))
			assert.Equal(t, *o.SmartClusters[1].Name, "new-demo")
			assert.Equal(t, *o.SmartClusters[1].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[1].CreateTime, "2026-06-09T16:18:20.998039823+08:00")
			assert.Equal(t, *o.SmartClusters[1].UpdateTime, "2026-06-09T18:08:12.90394089+08:00")
			assert.Equal(t, *o.SmartClusters[1].ObjectId, "SmartCluster-c30d039b-0b55-4a42-ae19-cae00a53735a")
			assert.Equal(t, *o.SmartClusters[1].Reason, "")
			assert.Equal(t, o.SmartClusters[1].Rules[0].Keywords[0], "hello")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].Sensitivity, float64(0.7))
			assert.Equal(t, *o.SmartClusters[2].Name, "demo")
			assert.Equal(t, *o.SmartClusters[2].ClusterType, "knowledge")
			assert.Equal(t, *o.SmartClusters[2].CreateTime, "2026-06-09T15:21:41.288773397+08:00")
			assert.Equal(t, *o.SmartClusters[2].UpdateTime, "2026-06-09T15:21:41.288773397+08:00")
			assert.Equal(t, *o.SmartClusters[2].ObjectId, "SmartCluster-fed69d7a-683a-4452-9d70-54148bd458e9")
			assert.Equal(t, *o.SmartClusters[2].Reason, "")
			assert.Equal(t, o.SmartClusters[2].Rules[0].Keywords[0], "hello")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].RuleType, "keywords")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].Sensitivity, float64(0.5))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<ListSmartClustersResponse>
    <SmartClusters>
        <SmartCluster>
            <CreateTime>2026-06-10T09:32:24.217552788+08:00</CreateTime>
            <ObjectId>FigureCluster-cd9da94f-feed-45cb-94df-fdcf3e21d712</ObjectId>
            <UpdateTime>2026-06-10T09:32:25.133798431+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-5</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/alice.jpg</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason></Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:32:46.2522507+08:00</CreateTime>
            <ObjectId>FigureCluster-bd06d605-d469-4339-b256-aa0b832be643</ObjectId>
            <UpdateTime>2026-06-09T18:32:46.2522507+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/OIP-C (1).jpg</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>[InvalidArgument.BaseURIs] The face quality is too low. status: 400, requestId: </Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:30:39.050200483+08:00</CreateTime>
            <ObjectId>FigureCluster-faaef950-81c6-4280-9fa4-bade5e94b3d1</ObjectId>
            <UpdateTime>2026-06-09T18:30:39.050200483+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/local-v1.txt</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>*error.OpError : InvalidArgument | File corrupt.</Reason>
        </SmartCluster>
        <SmartCluster>
            <CreateTime>2026-06-09T18:24:00.460823026+08:00</CreateTime>
            <ObjectId>FigureCluster-af67c554-10a4-4191-91be-5414e431ec43</ObjectId>
            <UpdateTime>2026-06-09T18:24:00.460823026+08:00</UpdateTime>
            <ClusterType>figure</ClusterType>
            <Name>demo-1</Name>
            <Rules>
                <Rule>
                    <BaseURIs>oss://demo-1889/local.txt</BaseURIs>
                    <RuleType>face</RuleType>
                </Rule>
            </Rules>
            <Reason>*error.OpError : InvalidArgument | File does not exist.</Reason>
        </SmartCluster>
    </SmartClusters>
</ListSmartClustersResponse>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?action=listSmartClusters&clusterType=figure&datasetName=your_dataset&maxResults=10&metaQuery&nextToken=nextToken&ruleTypes=%5B%22face%22%5D", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
			ClusterType: SmartClusterTypeFigure,
			MaxResults:  oss.Ptr(int64(10)),
			RuleTypes:   oss.Ptr(`["face"]`),
			NextToken:   oss.Ptr("nextToken"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.SmartClusters), 4)
			assert.Equal(t, *o.SmartClusters[0].Name, "demo-5")
			assert.Equal(t, *o.SmartClusters[0].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[0].CreateTime, "2026-06-10T09:32:24.217552788+08:00")
			assert.Equal(t, *o.SmartClusters[0].UpdateTime, "2026-06-10T09:32:25.133798431+08:00")
			assert.Equal(t, *o.SmartClusters[0].ObjectId, "FigureCluster-cd9da94f-feed-45cb-94df-fdcf3e21d712")
			assert.Equal(t, *o.SmartClusters[0].Reason, "")
			assert.Equal(t, o.SmartClusters[0].Rules[0].BaseURIs[0], "oss://demo-1889/alice.jpg")
			assert.Equal(t, *o.SmartClusters[0].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[1].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[1].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[1].CreateTime, "2026-06-09T18:32:46.2522507+08:00")
			assert.Equal(t, *o.SmartClusters[1].UpdateTime, "2026-06-09T18:32:46.2522507+08:00")
			assert.Equal(t, *o.SmartClusters[1].ObjectId, "FigureCluster-bd06d605-d469-4339-b256-aa0b832be643")
			assert.Equal(t, *o.SmartClusters[1].Reason, "[InvalidArgument.BaseURIs] The face quality is too low. status: 400, requestId: ")
			assert.Equal(t, o.SmartClusters[1].Rules[0].BaseURIs[0], "oss://demo-1889/OIP-C (1).jpg")
			assert.Equal(t, *o.SmartClusters[1].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[2].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[2].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[2].CreateTime, "2026-06-09T18:30:39.050200483+08:00")
			assert.Equal(t, *o.SmartClusters[2].UpdateTime, "2026-06-09T18:30:39.050200483+08:00")
			assert.Equal(t, *o.SmartClusters[2].ObjectId, "FigureCluster-faaef950-81c6-4280-9fa4-bade5e94b3d1")
			assert.Equal(t, *o.SmartClusters[2].Reason, "*error.OpError : InvalidArgument | File corrupt.")
			assert.Equal(t, o.SmartClusters[2].Rules[0].BaseURIs[0], "oss://demo-1889/local-v1.txt")
			assert.Equal(t, *o.SmartClusters[2].Rules[0].RuleType, "face")
			assert.Equal(t, *o.SmartClusters[3].Name, "demo-1")
			assert.Equal(t, *o.SmartClusters[3].ClusterType, "figure")
			assert.Equal(t, *o.SmartClusters[3].CreateTime, "2026-06-09T18:24:00.460823026+08:00")
			assert.Equal(t, *o.SmartClusters[3].UpdateTime, "2026-06-09T18:24:00.460823026+08:00")
			assert.Equal(t, *o.SmartClusters[3].ObjectId, "FigureCluster-af67c554-10a4-4191-91be-5414e431ec43")
			assert.Equal(t, *o.SmartClusters[3].Reason, "*error.OpError : InvalidArgument | File does not exist.")
			assert.Equal(t, o.SmartClusters[3].Rules[0].BaseURIs[0], "oss://demo-1889/local.txt")
			assert.Equal(t, *o.SmartClusters[3].Rules[0].RuleType, "face")
		},
	},
}

func TestMockListSmartClusters_Success(t *testing.T) {
	for _, c := range testMockListSmartClustersSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListSmartClusters(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListSmartClustersErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListSmartClustersRequest
	CheckOutputFn  func(t *testing.T, o *ListSmartClustersResult, err error)
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
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
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
			assert.Equal(t, "/bucket/?action=listSmartClusters&datasetName=your_dataset&metaQuery", strUrl)
		},
		&ListSmartClustersRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("your_dataset"),
		},
		func(t *testing.T, o *ListSmartClustersResult, err error) {
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

func TestMockListSmartClusters_Error(t *testing.T) {
	for _, c := range testMockListSmartClustersErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListSmartClusters(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteSmartClusterSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *DeleteSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteSmartCluster_Success(t *testing.T) {
	for _, c := range testMockDeleteSmartClusterSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteSmartClusterErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteSmartClusterRequest
	CheckOutputFn  func(t *testing.T, o *DeleteSmartClusterResult, err error)
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
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
			assert.Equal(t, "/bucket/?action=deleteSmartCluster&datasetName=test-dataset&metaQuery&objectId=cluster-abc123def456", strUrl)
		},
		&DeleteSmartClusterRequest{
			Bucket:      oss.Ptr("bucket"),
			DatasetName: oss.Ptr("test-dataset"),
			ObjectId:    oss.Ptr("cluster-abc123def456"),
		},
		func(t *testing.T, o *DeleteSmartClusterResult, err error) {
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

func TestMockDeleteSmartCluster_Error(t *testing.T) {
	for _, c := range testMockDeleteSmartClusterErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteSmartCluster(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
