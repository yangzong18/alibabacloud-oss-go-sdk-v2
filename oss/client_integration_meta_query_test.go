//go:build integration

package oss

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetaQuery(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	bucketNameAiSearch := bucketNamePrefix + randLowStr(6)
	_, err = client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketNameAiSearch),
	})
	assert.Nil(t, err)

	openRequest := &OpenMetaQueryRequest{
		Bucket: Ptr(bucketName),
	}
	openResult, err := client.OpenMetaQuery(context.TODO(), openRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, openResult.StatusCode)
	assert.NotEmpty(t, openResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetMetaQueryStatusRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetMetaQueryStatus(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	time.Sleep(1 * time.Second)

	doRequest := &DoMetaQueryRequest{
		Bucket: Ptr(bucketName),
		MetaQuery: &MetaQuery{
			Query: Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:  Ptr("Size"),
			Order: Ptr(MetaQueryOrderAsc),
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
	}
	doResult, err := client.DoMetaQuery(context.TODO(), doRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, doResult.StatusCode)
	assert.Empty(t, *doResult.NextToken)
	assert.Equal(t, len(doResult.Files), 0)
	assert.Equal(t, len(doResult.Aggregations), 2)
	time.Sleep(1 * time.Second)

	closeRequest := &CloseMetaQueryRequest{
		Bucket: Ptr(bucketName),
	}
	closeResult, err := client.CloseMetaQuery(context.TODO(), closeRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, closeResult.StatusCode)
	time.Sleep(1 * time.Second)

	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameAiSearch),
		Mode:   Ptr("semantic"),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, openResult.StatusCode)
	assert.NotEmpty(t, openResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	doRequest = &DoMetaQueryRequest{
		Bucket: Ptr(bucketNameAiSearch),
		Mode:   Ptr("semantic"),
		MetaQuery: &MetaQuery{
			MaxResults: Ptr(int64(99)),
			Query:      Ptr("Overlook the snow-covered forest"),
			MediaTypes: &MetaQueryMediaTypes{
				MediaTypes: []string{"image"},
			},
			SimpleQuery: Ptr(`{"Operation":"gt", "Field": "Size", "Value": "30"}`),
		},
	}
	doResult, err = client.DoMetaQuery(context.TODO(), doRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, doResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	openRequest = &OpenMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	openResult, err = client.OpenMetaQuery(context.TODO(), openRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetMetaQueryStatusRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getResult, err = client.GetMetaQueryStatus(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	doRequest = &DoMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
		MetaQuery: &MetaQuery{
			Query: Ptr(`{"Field": "Size","Value": "1048576","Operation": "gt"}`),
			Sort:  Ptr("Size"),
			Order: Ptr(MetaQueryOrderAsc),
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
	}
	doResult, err = client.DoMetaQuery(context.TODO(), doRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	closeRequest = &CloseMetaQueryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	closeResult, err = client.CloseMetaQuery(context.TODO(), closeRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDoMetaQueryAction(t *testing.T) {
	after := before(t)
	defer after(t)

	var err error
	bucketName := bucketNamePrefix + randLowStr(6)
	dataSetName := objectNamePrefix + randLowStr(6)
	client := getDefaultClient()
	_, err = client.PutBucket(context.TODO(), &PutBucketRequest{
		Bucket: Ptr(bucketName),
	})
	assert.Nil(t, err)

	var serr *ServiceError
	_, err = client.DoMetaQueryAction(context.TODO(), &DoMetaQueryActionRequest{
		Bucket: Ptr(bucketName),
		Action: Ptr("createDataset"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"datasetName": dataSetName,
			},
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "OperationNotSupported", serr.Code)
	assert.Equal(t, "The operation is not supported for this resource", serr.Message)
	assert.Equal(t, "0037-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DoMetaQueryAction(context.TODO(), &DoMetaQueryActionRequest{
		Bucket: Ptr(bucketName),
		Action: Ptr("createDataset"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"datasetName": dataSetName,
			},
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "OperationNotSupported", serr.Code)
	assert.Equal(t, "The operation is not supported for this resource", serr.Message)
	assert.Equal(t, "0037-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DoMetaQueryAction(context.TODO(), &DoMetaQueryActionRequest{
		Bucket: Ptr(bucketName),
		Action: Ptr("listDatasets"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "OperationNotSupported", serr.Code)
	assert.Equal(t, "The operation is not supported for this resource", serr.Message)
	assert.Equal(t, "0037-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.DoMetaQueryAction(context.TODO(), &DoMetaQueryActionRequest{
		Bucket: Ptr(bucketNameNotExist),
		Action: Ptr("listDatasets"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
}

func TestDoDataPipeLineAction(t *testing.T) {
	after := before(t)
	defer after(t)

	var err error
	endpoint := "https://oss-" + region_ + ".aliyuncs.com"
	client := getClient(region_, endpoint)
	_, err = client.DoDataPipeLineAction(context.TODO(), &DoDataPipeLineActionRequest{
		Action: Ptr("listDataPipelineConfigurations"),
	})
	assert.Nil(t, err)

	var serr *ServiceError
	_, err = client.DoDataPipeLineAction(context.TODO(), &DoDataPipeLineActionRequest{
		Action: Ptr("putDataPipelineConfiguration"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"dataPipelineName": "data-pipeline",
				"role":             "not-exist-role",
			},
		},
		Body: strings.NewReader(`<?xml version="1.0" encoding="UTF-8" ?>
<DataPipelineConfiguration>
    <DataPipelineDescription>Vectorize business data using the BERT multimodal model</DataPipelineDescription>
    <Sources>
        <InputBucket>oss-sdk-test-an</InputBucket>
        <InputDataScope>All</InputDataScope>
        <IgnoreDelete>true</IgnoreDelete>
        <FilterConfiguration>
            <PrefixSet>prefix1/</PrefixSet>
            <ObjectMediaTypes>video</ObjectMediaTypes>
        </FilterConfiguration>
    </Sources>
    <DataPipelineEmbeddingConfiguration>
        <EmbeddingProvider>bailian</EmbeddingProvider>
        <ApiKey>sk-123323423423423423424242425436457657567</ApiKey>
        <Model>qwen2.5-vl-embedding</Model>
        <FPS>1</FPS>
    </DataPipelineEmbeddingConfiguration>
    <Destination>
        <VectorBucketName>my-vector-bucket</VectorBucketName>
        <VectorIndexNames>index</VectorIndexNames>
        <ObjectTagToMetadata>key1</ObjectTagToMetadata>
        <ObjectTagToMetadata>key2</ObjectTagToMetadata>
        <UsermetaToMetadata>x-oss-meta-key1</UsermetaToMetadata>
    </Destination>
    <DataPipelineError>
        <ErrorMode>ignoreAndRecord</ErrorMode>
        <ErrorBucket>my-error-bucket</ErrorBucket>
        <ErrorPrefix>error-output/</ErrorPrefix>
    </DataPipelineError>
</DataPipelineConfiguration>`),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "The ServiceRole (not-exist-role) you provided does not exist.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	_, err = client.DoDataPipeLineAction(context.TODO(), &DoDataPipeLineActionRequest{
		Action: Ptr("getDataPipelineConfiguration"),
		RequestCommon: RequestCommon{
			Parameters: map[string]string{
				"dataPipelineName": "not-exist-data-pipeline",
			},
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchDataPipeline", serr.Code)
	assert.Equal(t, "The specified resource DataPipeline is not found.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
