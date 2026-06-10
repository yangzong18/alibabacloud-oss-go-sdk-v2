package dataprocess

import (
	"context"
	"encoding/xml"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type CreateSmartClusterRequest struct {
	Bucket       *string          `input:"host,bucket,required"`
	DatasetName  *string          `input:"query,datasetName,required"`
	Name         *string          `input:"query,name,required"`
	ClusterType  SmartClusterType `input:"query,clusterType,required"`
	Rules        *string          `input:"query,rules,required"`
	Description  *string          `input:"query,description"`
	Notification *string          `input:"query,notification"`
	oss.RequestCommon
}

type CreateSmartClusterResult struct {
	XMLName  xml.Name `xml:"CreateSmartClusterResponse"`
	ObjectId *string  `xml:"ObjectId"`
	oss.ResultCommon
}

func (c *Client) CreateSmartCluster(ctx context.Context, request *CreateSmartClusterRequest, optFns ...func(*oss.Options)) (*CreateSmartClusterResult, error) {
	var err error
	if request == nil {
		request = &CreateSmartClusterRequest{}
	}

	input := &oss.OperationInput{
		OpName: "CreateSmartCluster",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "createSmartCluster",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &CreateSmartClusterResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type GetSmartClusterRequest struct {
	Bucket      *string `input:"host,bucket,required"`
	DatasetName *string `input:"query,datasetName,required"`
	ObjectId    *string `input:"query,objectId,required"`
	oss.RequestCommon
}

type GetSmartClusterResult struct {
	XMLName      xml.Name          `xml:"GetSmartClusterResponse"`
	SmartCluster *SmartClusterInfo `xml:"SmartCluster,omitempty"`
	oss.ResultCommon
}

type SmartClusterInfo struct {
	ObjectId     *string           `xml:"ObjectId"`
	ClusterType  *string           `xml:"ClusterType"`
	Name         *string           `xml:"Name"`
	Description  *string           `xml:"Description"`
	Rules        []Rule            `xml:"Rules>Rule"`
	Reason       *string           `xml:"Reason"`
	Notification *NotificationInfo `xml:"Notification"`
	CreateTime   *string           `xml:"CreateTime"`
	UpdateTime   *string           `xml:"UpdateTime"`
}

type Rule struct {
	RuleType    *string  `xml:"RuleType"`
	BaseURIs    []string `xml:"BaseURIs"`
	Keywords    []string `xml:"Keywords"`
	Sensitivity *float64 `xml:"Sensitivity"`
}

type NotificationInfo struct {
	MNS *TopicName `xml:"MNS"`
}

type TopicName struct {
	TopicName *string `xml:"TopicName"`
}

func (c *Client) GetSmartCluster(ctx context.Context, request *GetSmartClusterRequest, optFns ...func(*oss.Options)) (*GetSmartClusterResult, error) {
	var err error
	if request == nil {
		request = &GetSmartClusterRequest{}
	}

	input := &oss.OperationInput{
		OpName: "GetSmartCluster",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "getSmartCluster",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &GetSmartClusterResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type UpdateSmartClusterRequest struct {
	Bucket       *string `input:"host,bucket,required"`
	DatasetName  *string `input:"query,datasetName,required"`
	ObjectId     *string `input:"query,objectId,required"`
	Name         *string `input:"query,name"`
	Description  *string `input:"query,description"`
	Rules        *string `input:"query,rules"`
	Notification *string `input:"query,notification"`
	oss.RequestCommon
}

type UpdateSmartClusterResult struct {
	oss.ResultCommon
}

func (c *Client) UpdateSmartCluster(ctx context.Context, request *UpdateSmartClusterRequest, optFns ...func(*oss.Options)) (*UpdateSmartClusterResult, error) {
	var err error
	if request == nil {
		request = &UpdateSmartClusterRequest{}
	}

	input := &oss.OperationInput{
		OpName: "UpdateSmartCluster",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "updateSmartCluster",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &UpdateSmartClusterResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type DeleteSmartClusterRequest struct {
	Bucket      *string `input:"host,bucket,required"`
	DatasetName *string `input:"query,datasetName,required"`
	ObjectId    *string `input:"query,objectId,required"`
	oss.RequestCommon
}

type DeleteSmartClusterResult struct {
	oss.ResultCommon
}

func (c *Client) DeleteSmartCluster(ctx context.Context, request *DeleteSmartClusterRequest, optFns ...func(*oss.Options)) (*DeleteSmartClusterResult, error) {
	var err error
	if request == nil {
		request = &DeleteSmartClusterRequest{}
	}

	input := &oss.OperationInput{
		OpName: "DeleteSmartCluster",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "deleteSmartCluster",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &DeleteSmartClusterResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}

type ListSmartClustersRequest struct {
	Bucket      *string          `input:"host,bucket,required"`
	DatasetName *string          `input:"query,datasetName,required"`
	ClusterType SmartClusterType `input:"query,clusterType"`
	MaxResults  *int64           `input:"query,maxResults"`
	RuleTypes   *string          `input:"query,ruleTypes"`
	NextToken   *string          `input:"query,nextToken"`
	oss.RequestCommon
}

type ListSmartClustersResult struct {
	XMLName       xml.Name           `xml:"ListSmartClustersResponse"`
	SmartClusters []SmartClusterInfo `xml:"SmartClusters>SmartCluster"`
	NextToken     *string            `xml:"NextToken"`
	oss.ResultCommon
}

func (c *Client) ListSmartClusters(ctx context.Context, request *ListSmartClustersRequest, optFns ...func(*oss.Options)) (*ListSmartClustersResult, error) {
	var err error
	if request == nil {
		request = &ListSmartClustersRequest{}
	}

	input := &oss.OperationInput{
		OpName: "ListSmartClusters",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "listSmartClusters",
		},
		Bucket: request.Bucket,
	}

	if err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.client.InvokeOperation(ctx, input, optFns...)
	if err != nil {
		return nil, err
	}

	result := &ListSmartClustersResult{}

	if err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	}); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, nil
}
