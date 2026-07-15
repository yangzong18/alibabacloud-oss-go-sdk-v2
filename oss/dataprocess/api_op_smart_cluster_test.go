package dataprocess

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateSmartClusterRequest
	var input *oss.OperationInput
	var err error

	request = &CreateSmartClusterRequest{}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &CreateSmartClusterRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &CreateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Name.")

	request = &CreateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		Name:        oss.Ptr("your_name"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ClusterType.")

	request = &CreateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		Name:        oss.Ptr("your_name"),
		ClusterType: SmartClusterTypeFigure,
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Rules.")

	request = &CreateSmartClusterRequest{
		Bucket:       oss.Ptr("bucket"),
		DatasetName:  oss.Ptr("your_dataset"),
		Name:         oss.Ptr("your_name"),
		ClusterType:  SmartClusterTypeKnowledge,
		Rules:        oss.Ptr(`[{"RuleType": "keywords","Keywords": ["car"]}]`),
		Description:  oss.Ptr("your_description"),
		Notification: oss.Ptr(`{"MNS":{"TopicName":"imm-cluster-notification"}}`),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["name"], "your_name")
	assert.Equal(t, input.Parameters["clusterType"], "knowledge")
	assert.Equal(t, input.Parameters["rules"], "[{\"RuleType\": \"keywords\",\"Keywords\": [\"car\"]}]")
	assert.Equal(t, input.Parameters["description"], "your_description")
	assert.Equal(t, input.Parameters["notification"], `{"MNS":{"TopicName":"imm-cluster-notification"}}`)

	request = &CreateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
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
		Description: oss.Ptr("your_description"),
		Notification: oss.Ptr(SmartClusterNotification{
			MNS: &SmartClusterTopicName{
				TopicName: oss.Ptr("imm-cluster-notification"),
			},
		}.ToParameterValue()),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["name"], "your_name")
	assert.Equal(t, input.Parameters["clusterType"], "knowledge")
	assert.Equal(t, input.Parameters["rules"], "[{\"RuleType\":\"keywords\",\"Keywords\":[\"car\"]}]")
	assert.Equal(t, input.Parameters["description"], "your_description")
	assert.Equal(t, input.Parameters["notification"], `{"MNS":{"TopicName":"imm-cluster-notification"}}`)
}

func TestUnmarshalOutput_CreateSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CreateSmartClusterResponse>
  <ObjectId>cluster-abc123def456</ObjectId>
</CreateSmartClusterResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &CreateSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.ObjectId, "cluster-abc123def456")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &CreateSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetSmartClusterRequest
	var input *oss.OperationInput
	var err error

	request = &GetSmartClusterRequest{}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetSmartClusterRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &GetSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ObjectId.")

	request = &GetSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		ObjectId:    oss.Ptr("cluster-abc123def456"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["objectId"], "cluster-abc123def456")
}

func TestUnmarshalOutput_GetSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
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
</GetSmartClusterResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.SmartCluster.ObjectId, "cluster-abc123def456")
	assert.Equal(t, *result.SmartCluster.ClusterType, "figure")
	assert.Equal(t, *result.SmartCluster.CreateTime, "2026-05-20T11:00:00.000+08:00")
	assert.Equal(t, *result.SmartCluster.UpdateTime, "2026-05-20T11:08:00.000+08:00")
	assert.Equal(t, *result.SmartCluster.Name, "face-cluster-alice")
	assert.Equal(t, *result.SmartCluster.Description, "this is a demo")
	assert.Equal(t, *result.SmartCluster.Rules[0].RuleType, "face")
	assert.Equal(t, result.SmartCluster.Rules[0].BaseURIs[0], "oss://examplebucket/refs/alice.jpg")
	assert.Equal(t, *result.SmartCluster.Rules[0].Sensitivity, 0.7)
	assert.Equal(t, *result.SmartCluster.Notification.MNS.TopicName, "imm-cluster-notification")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_UpdateSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *UpdateSmartClusterRequest
	var input *oss.OperationInput
	var err error

	request = &UpdateSmartClusterRequest{}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &UpdateSmartClusterRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &UpdateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ObjectId.")

	request = &UpdateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		ObjectId:    oss.Ptr("cluster-abc123def456"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["objectId"], "cluster-abc123def456")

	request = &UpdateSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		ObjectId:    oss.Ptr("cluster-abc123def456"),
		Description: oss.Ptr("this is a demo"),
		Name:        oss.Ptr("face-cluster-alice"),
		Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["objectId"], "cluster-abc123def456")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["name"], "face-cluster-alice")
	assert.Equal(t, input.Parameters["rules"], "[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]")

	request = &UpdateSmartClusterRequest{
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
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["objectId"], "cluster-abc123def456")
	assert.Equal(t, input.Parameters["description"], "this is a demo")
	assert.Equal(t, input.Parameters["name"], "face-cluster-alice")
	assert.Equal(t, input.Parameters["rules"], "[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]")
	assert.Equal(t, input.Parameters["notification"], "{\"MNS\":{\"TopicName\":\"imm-cluster-notification\"}}")
}

func TestUnmarshalOutput_UpdateSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	body := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<UpdateSmartClusterResponse>\n  <ObjectId>cluster-abc123def456</ObjectId>\n</UpdateSmartClusterResponse>"
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &UpdateSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.ObjectId, "cluster-abc123def456")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &UpdateSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListSmartClusters(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListSmartClustersRequest
	var input *oss.OperationInput
	var err error

	request = &ListSmartClustersRequest{}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListSmartClustersRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &ListSmartClustersRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		ClusterType: SmartClusterTypeFigure,
		MaxResults:  oss.Ptr(int64(10)),
		RuleTypes:   oss.Ptr(`["face"]`),
		NextToken:   oss.Ptr("nextToken"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["clusterType"], "figure")
	assert.Equal(t, input.Parameters["maxResults"], "10")
	assert.Equal(t, input.Parameters["ruleTypes"], `["face"]`)
	assert.Equal(t, input.Parameters["nextToken"], "nextToken")
}

func TestUnmarshalOutput_ListSmartClusters(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `<ListSmartClustersResponse>
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
</ListSmartClustersResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListSmartClustersResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.SmartClusters), 4)
	assert.Equal(t, *result.SmartClusters[0].Name, "demo-5")
	assert.Equal(t, *result.SmartClusters[0].ClusterType, "figure")
	assert.Equal(t, *result.SmartClusters[0].CreateTime, "2026-06-10T09:32:24.217552788+08:00")
	assert.Equal(t, *result.SmartClusters[0].UpdateTime, "2026-06-10T09:32:25.133798431+08:00")
	assert.Equal(t, *result.SmartClusters[0].ObjectId, "FigureCluster-cd9da94f-feed-45cb-94df-fdcf3e21d712")
	assert.Equal(t, *result.SmartClusters[0].Reason, "")
	assert.Equal(t, result.SmartClusters[0].Rules[0].BaseURIs[0], "oss://demo-1889/alice.jpg")
	assert.Equal(t, *result.SmartClusters[0].Rules[0].RuleType, "face")
	assert.Equal(t, *result.SmartClusters[1].Name, "demo-1")
	assert.Equal(t, *result.SmartClusters[1].ClusterType, "figure")
	assert.Equal(t, *result.SmartClusters[1].CreateTime, "2026-06-09T18:32:46.2522507+08:00")
	assert.Equal(t, *result.SmartClusters[1].UpdateTime, "2026-06-09T18:32:46.2522507+08:00")
	assert.Equal(t, *result.SmartClusters[1].ObjectId, "FigureCluster-bd06d605-d469-4339-b256-aa0b832be643")
	assert.Equal(t, *result.SmartClusters[1].Reason, "[InvalidArgument.BaseURIs] The face quality is too low. status: 400, requestId: ")
	assert.Equal(t, result.SmartClusters[1].Rules[0].BaseURIs[0], "oss://demo-1889/OIP-C (1).jpg")
	assert.Equal(t, *result.SmartClusters[1].Rules[0].RuleType, "face")
	assert.Equal(t, *result.SmartClusters[2].Name, "demo-1")
	assert.Equal(t, *result.SmartClusters[2].ClusterType, "figure")
	assert.Equal(t, *result.SmartClusters[2].CreateTime, "2026-06-09T18:30:39.050200483+08:00")
	assert.Equal(t, *result.SmartClusters[2].UpdateTime, "2026-06-09T18:30:39.050200483+08:00")
	assert.Equal(t, *result.SmartClusters[2].ObjectId, "FigureCluster-faaef950-81c6-4280-9fa4-bade5e94b3d1")
	assert.Equal(t, *result.SmartClusters[2].Reason, "*error.OpError : InvalidArgument | File corrupt.")
	assert.Equal(t, result.SmartClusters[2].Rules[0].BaseURIs[0], "oss://demo-1889/local-v1.txt")
	assert.Equal(t, *result.SmartClusters[2].Rules[0].RuleType, "face")
	assert.Equal(t, *result.SmartClusters[3].Name, "demo-1")
	assert.Equal(t, *result.SmartClusters[3].ClusterType, "figure")
	assert.Equal(t, *result.SmartClusters[3].CreateTime, "2026-06-09T18:24:00.460823026+08:00")
	assert.Equal(t, *result.SmartClusters[3].UpdateTime, "2026-06-09T18:24:00.460823026+08:00")
	assert.Equal(t, *result.SmartClusters[3].ObjectId, "FigureCluster-af67c554-10a4-4191-91be-5414e431ec43")
	assert.Equal(t, *result.SmartClusters[3].Reason, "*error.OpError : InvalidArgument | File does not exist.")
	assert.Equal(t, result.SmartClusters[3].Rules[0].BaseURIs[0], "oss://demo-1889/local.txt")
	assert.Equal(t, *result.SmartClusters[3].Rules[0].RuleType, "face")

	body = `<ListSmartClustersResponse>
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
</ListSmartClustersResponse>`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListSmartClustersResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, len(result.SmartClusters), 3)
	assert.Equal(t, *result.SmartClusters[0].Name, "demo-2")
	assert.Equal(t, *result.SmartClusters[0].ClusterType, "knowledge")
	assert.Equal(t, *result.SmartClusters[0].CreateTime, "2026-06-10T09:54:27.484585901+08:00")
	assert.Equal(t, *result.SmartClusters[0].UpdateTime, "2026-06-10T09:54:27.484585901+08:00")
	assert.Equal(t, *result.SmartClusters[0].ObjectId, "SmartCluster-cb9f8c95-281f-490b-b677-eca2f6ff0c19")
	assert.Equal(t, *result.SmartClusters[0].Reason, "")
	assert.Equal(t, result.SmartClusters[0].Rules[0].Keywords[0], "cat")
	assert.Equal(t, *result.SmartClusters[0].Rules[0].RuleType, "keywords")
	assert.Equal(t, *result.SmartClusters[0].Rules[0].Sensitivity, float64(0.5))
	assert.Equal(t, *result.SmartClusters[1].Name, "new-demo")
	assert.Equal(t, *result.SmartClusters[1].ClusterType, "knowledge")
	assert.Equal(t, *result.SmartClusters[1].CreateTime, "2026-06-09T16:18:20.998039823+08:00")
	assert.Equal(t, *result.SmartClusters[1].UpdateTime, "2026-06-09T18:08:12.90394089+08:00")
	assert.Equal(t, *result.SmartClusters[1].ObjectId, "SmartCluster-c30d039b-0b55-4a42-ae19-cae00a53735a")
	assert.Equal(t, *result.SmartClusters[1].Reason, "")
	assert.Equal(t, result.SmartClusters[1].Rules[0].Keywords[0], "hello")
	assert.Equal(t, *result.SmartClusters[1].Rules[0].RuleType, "keywords")
	assert.Equal(t, *result.SmartClusters[1].Rules[0].Sensitivity, float64(0.7))
	assert.Equal(t, *result.SmartClusters[2].Name, "demo")
	assert.Equal(t, *result.SmartClusters[2].ClusterType, "knowledge")
	assert.Equal(t, *result.SmartClusters[2].CreateTime, "2026-06-09T15:21:41.288773397+08:00")
	assert.Equal(t, *result.SmartClusters[2].UpdateTime, "2026-06-09T15:21:41.288773397+08:00")
	assert.Equal(t, *result.SmartClusters[2].ObjectId, "SmartCluster-fed69d7a-683a-4452-9d70-54148bd458e9")
	assert.Equal(t, *result.SmartClusters[2].Reason, "")
	assert.Equal(t, result.SmartClusters[2].Rules[0].Keywords[0], "hello")
	assert.Equal(t, *result.SmartClusters[2].Rules[0].RuleType, "keywords")
	assert.Equal(t, *result.SmartClusters[2].Rules[0].Sensitivity, float64(0.5))

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListSmartClustersResult{}
	err = c.client.UnmarshalOutput(result, output, func(result interface{}, output *oss.OperationOutput) error {
		if output.Body == nil {
			return nil
		}
		defer output.Body.Close()
		return xml.NewDecoder(output.Body).Decode(result)
	})
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteSmartClusterRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteSmartClusterRequest{}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteSmartClusterRequest{
		Bucket: oss.Ptr("bucket"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DatasetName.")

	request = &DeleteSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, ObjectId.")

	request = &DeleteSmartClusterRequest{
		Bucket:      oss.Ptr("bucket"),
		DatasetName: oss.Ptr("your_dataset"),
		ObjectId:    oss.Ptr("cluster-abc123def456"),
	}
	input = &oss.OperationInput{
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
	err = c.client.MarshalInput(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "bucket")
	assert.Equal(t, input.Parameters["datasetName"], "your_dataset")
	assert.Equal(t, input.Parameters["objectId"], "cluster-abc123def456")
}

func TestUnmarshalOutput_DeleteSmartCluster(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error

	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DeleteSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "Bad Request",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteSmartClusterResult{}
	err = c.client.UnmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "Bad Request")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
