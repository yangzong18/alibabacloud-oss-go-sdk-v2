package vectors

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorIndexRequest
	var input *oss.OperationInput
	var err error

	request = &PutVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DataType")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		DataType:  oss.Ptr("string"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Dimension")

	request = &PutVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		DataType:  oss.Ptr("string"),
		IndexName: oss.Ptr("exampleIndex"),
		Dimension: oss.Ptr(128),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, DistanceMetric")

	request = &PutVectorIndexRequest{
		Bucket:         oss.Ptr("oss-demo"),
		DataType:       oss.Ptr("string"),
		Dimension:      oss.Ptr(128),
		DistanceMetric: oss.Ptr("cosine"),
		IndexName:      oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\"}")

	request = &PutVectorIndexRequest{
		Bucket:         oss.Ptr("oss-demo"),
		DataType:       oss.Ptr("string"),
		Dimension:      oss.Ptr(128),
		DistanceMetric: oss.Ptr("cosine"),
		IndexName:      oss.Ptr("exampleIndex"),
		Metadata: map[string]any{
			"nonFilterableMetadataKeys": []string{"foo", "bar"},
		},
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"PutVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"dataType\":\"string\",\"dimension\":128,\"distanceMetric\":\"cosine\",\"indexName\":\"exampleIndex\",\"metadata\":{\"nonFilterableMetadataKeys\":[\"foo\",\"bar\"]}}")
}

func TestUnmarshalOutput_PutVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_GetVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *GetVectorIndexRequest
	var input *oss.OperationInput
	var err error

	request = &GetVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &GetVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &GetVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("demo"),
	}
	input = &oss.OperationInput{
		OpName: "GetVectorIndex",
		Method: "POST",
		Parameters: map[string]string{
			"getVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"GetVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["GetVectorIndex"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"demo\"}")
}

func TestUnmarshalOutput_GetVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
   "index": { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "string",
      "metadata": { 
         "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
      "bucketArn": "acs:oss:::test-bucket"
   }
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	//dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Index.CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Index.DataType, "string")
	assert.Equal(t, *result.Index.Dimension, 128)
	assert.Equal(t, *result.Index.DistanceMetric, "string")
	assert.Equal(t, *result.Index.IndexName, "string")
	assert.Len(t, result.Index.Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Index.Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Index.Status, "running")
	assert.Equal(t, *result.Index.BucketArn, "acs:oss:::test-bucket")

	body = `{
  "index": {
    "indexName": "exampleIndex",
    "mode": "fusion",
    "dataType": "string",
    "dimension": 128,
    "status": "running",
    "bucketArn": "acs:oss:::test-bucket",
    "schemaConfiguration": {
      "fields": [
        {
          "dataType": "float32",
          "dimension": 1024,
          "distanceMetric": "euclidean",
          "name": "vector_1",
          "type": "vector"
        },
        {
          "dataType": "float32",
          "dimension": 512,
          "distanceMetric": "cosine",
          "name": "vector_2",
          "type": "vector"
        },
        {
          "isArray": true,
          "name": "timestamps",
          "type": "long"
        },
        {
          "name": "price",
          "type": "double"
        },
        {
          "name": "ip",
          "type": "ip"
        },
        {
          "name": "location",
          "type": "geoPoint"
        },
        {
          "name": "tag",
          "type": "string"
        },
        {
          "isPartitionKey": true,
          "name": "user_id",
          "type": "string"
        },
        {
          "isArray": true,
          "name": "tags",
          "type": "string"
        },
        {
          "exactMatch": true,
          "name": "title_1",
          "text": {
            "analyzer": "standard",
            "analyzerParameters": {
              "caseSensitive": true,
              "delimitWord": false
            },
            "enabled": true
          },
          "type": "string"
        },
        {
          "exactMatch": false,
          "name": "title_2",
          "text": {
            "analyzer": "split",
            "analyzerParameters": {
              "caseSensitive": true,
              "delimiter": " "
            },
            "enabled": true
          },
          "type": "string"
        }
      ]
    }
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &GetVectorIndexResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	//dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.Index.DataType, "string")
	assert.Equal(t, *result.Index.IndexName, "exampleIndex")
	assert.Equal(t, *result.Index.Mode, "fusion")
	assert.Equal(t, *result.Index.Dimension, 128)
	assert.Equal(t, *result.Index.Status, "running")
	assert.Equal(t, *result.Index.BucketArn, "acs:oss:::test-bucket")
	assert.Len(t, result.Index.SchemaConfiguration["fields"], 11)
	if fields, ok := result.Index.SchemaConfiguration["fields"].([]any); ok {
		for i, f := range fields {
			if field, ok := f.(map[string]any); ok {
				if i == 0 {
					assert.Equal(t, field["dataType"].(string), "float32")
					assert.Equal(t, field["dimension"].(float64), float64(1024))
					assert.Equal(t, field["distanceMetric"].(string), "euclidean")
					assert.Equal(t, field["name"].(string), "vector_1")
					assert.Equal(t, field["type"].(string), "vector")
				}
				if i == 1 {
					assert.Equal(t, field["dataType"].(string), "float32")
					assert.Equal(t, field["dimension"].(float64), float64(512))
					assert.Equal(t, field["distanceMetric"].(string), "cosine")
					assert.Equal(t, field["name"].(string), "vector_2")
					assert.Equal(t, field["type"].(string), "vector")
				}
				if i == 2 {
					assert.Equal(t, field["isArray"].(bool), true)
					assert.Equal(t, field["name"].(string), "timestamps")
					assert.Equal(t, field["type"].(string), "long")
				}
				if i == 3 {
					assert.Equal(t, field["name"].(string), "price")
					assert.Equal(t, field["type"].(string), "double")
				}
				if i == 4 {
					assert.Equal(t, field["name"].(string), "ip")
					assert.Equal(t, field["type"].(string), "ip")
				}
				if i == 5 {
					assert.Equal(t, field["name"].(string), "location")
					assert.Equal(t, field["type"].(string), "geoPoint")
				}
				if i == 6 {
					assert.Equal(t, field["name"].(string), "tag")
					assert.Equal(t, field["type"].(string), "string")
				}
				if i == 7 {
					assert.Equal(t, field["isPartitionKey"].(bool), true)
					assert.Equal(t, field["name"].(string), "user_id")
					assert.Equal(t, field["type"].(string), "string")
				}
				if i == 8 {
					assert.Equal(t, field["isArray"].(bool), true)
					assert.Equal(t, field["name"].(string), "tags")
					assert.Equal(t, field["type"].(string), "string")
				}
				if i == 9 {
					assert.Equal(t, field["exactMatch"].(bool), true)
					assert.Equal(t, field["name"].(string), "title_1")
					assert.Equal(t, field["type"].(string), "string")
				}
				if i == 10 {
					assert.Equal(t, field["exactMatch"].(bool), false)
					assert.Equal(t, field["name"].(string), "title_2")
					assert.Equal(t, field["type"].(string), "string")
				}
			}
		}

		output = &oss.OperationOutput{
			StatusCode: 404,
			Status:     "NoSuchBucket",
			Headers: http.Header{
				"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
				"Content-Type":     {"application/json"},
			},
		}
		result = &GetVectorIndexResult{}
		err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
		assert.Nil(t, err)
		assert.Equal(t, result.StatusCode, 404)
		assert.Equal(t, result.Status, "NoSuchBucket")
		assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
		assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
		output = &oss.OperationOutput{
			StatusCode: 400,
			Status:     "InvalidArgument",
			Headers: http.Header{
				"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
				"Content-Type":     {"application/json"},
			},
		}
		result = &GetVectorIndexResult{}
		err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
		assert.Nil(t, err)
		assert.Equal(t, result.StatusCode, 400)
		assert.Equal(t, result.Status, "InvalidArgument")
		assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
		assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

		body = `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
		output = &oss.OperationOutput{
			StatusCode: 403,
			Status:     "AccessDenied",
			Body:       io.NopCloser(bytes.NewReader([]byte(body))),
			Headers: http.Header{
				"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
				"Content-Type":     {"application/json"},
			},
		}
		result = &GetVectorIndexResult{}
		err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
		assert.Nil(t, err)
		assert.Equal(t, result.StatusCode, 403)
		assert.Equal(t, result.Status, "AccessDenied")
		assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
		assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	}
}

func TestMarshalInput_ListVectorIndexes(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *ListVectorIndexesRequest
	var input *oss.OperationInput
	var err error

	request = &ListVectorIndexesRequest{}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &ListVectorIndexesRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")

	request = &ListVectorIndexesRequest{
		Bucket: oss.Ptr("oss-demo"),
		Prefix: oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"prefix\":\"prefix\"}")

	request = &ListVectorIndexesRequest{
		Bucket:     oss.Ptr("oss-demo"),
		MaxResults: 100,
		NextToken:  oss.Ptr("123"),
		Prefix:     oss.Ptr("prefix"),
	}
	input = &oss.OperationInput{
		OpName: "ListVectorIndexes",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"listVectorIndexes": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"ListVectorIndexes"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["ListVectorIndexes"], "")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, input.Headers[oss.HTTPHeaderContentType], contentTypeJSON)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"maxResults\":100,\"nextToken\":\"123\",\"prefix\":\"prefix\"}")
}

func TestUnmarshalOutput_ListVectorIndexes(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	body := `{
  "indexes": [
    { 
      "createTime": "2025-08-02T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo1",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo", "bar"]
      },
      "status": "running",
	  "bucketArn": "acs:oss:::test-bucket"
    },
    { 
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "dataType": "string",
      "dimension": 128,
      "distanceMetric": "string",
      "indexName": "demo2",
      "metadata": { 
        "nonFilterableMetadataKeys": ["foo2", "bar2"]
      },
      "status": "deleting",
	  "bucketArn": "acs:oss:::test-bucket"
    }
  ],
  "nextToken": "123"
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result := &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Indexes), 2)
	assert.Equal(t, *result.Indexes[0].CreateTime, time.Date(2025, time.August, 2, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[0].DataType, "string")
	assert.Equal(t, *result.Indexes[0].Dimension, 128)
	assert.Equal(t, *result.Indexes[0].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[0].IndexName, "demo1")
	assert.Len(t, result.Indexes[0].Metadata["nonFilterableMetadataKeys"], 2)
	if metadataValue, ok := result.Indexes[0].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo")
			assert.Equal(t, keys[1].(string), "bar")
		}
	}
	assert.Equal(t, *result.Indexes[0].Status, "running")
	assert.Equal(t, *result.Indexes[0].BucketArn, "acs:oss:::test-bucket")

	assert.Equal(t, *result.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[1].DataType, "string")
	assert.Equal(t, *result.Indexes[1].Dimension, 128)
	assert.Equal(t, *result.Indexes[1].DistanceMetric, "string")
	assert.Equal(t, *result.Indexes[1].IndexName, "demo2")
	if metadataValue, ok := result.Indexes[1].Metadata["nonFilterableMetadataKeys"]; ok {
		if keys, ok := metadataValue.([]any); ok {
			assert.Equal(t, keys[0].(string), "foo2")
			assert.Equal(t, keys[1].(string), "bar2")
		}
	}
	assert.Equal(t, *result.Indexes[1].BucketArn, "acs:oss:::test-bucket")
	assert.Equal(t, *result.Indexes[1].Status, "deleting")
	assert.Equal(t, *result.NextToken, "123")

	body = `{
  "nextToken": "CAESCG15aC1mESgEeKaGA",
  "indexes": [
    {
      "vectorBucketName": "oss-vector-bucket1",
      "indexName": "index1",
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "dataType": "float32",
      "dimension": 512,
      "distanceMetric": "euclidean",
      "status": "enable",
      "mode": "fusion"
    },
    {
      "vectorBucketName": "oss-vector-bucket1",
      "indexName": "index2",
      "createTime": "2025-08-20T10:49:17.289372919Z",
      "status": "enable",
      "mode": "standard"
    }
  ]
}`
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.Indexes), 2)
	assert.Equal(t, *result.NextToken, "CAESCG15aC1mESgEeKaGA")
	assert.Equal(t, *result.Indexes[0].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[0].IndexName, "index1")
	assert.Equal(t, *result.Indexes[0].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))
	assert.Equal(t, *result.Indexes[0].DataType, "float32")
	assert.Equal(t, *result.Indexes[0].Dimension, 512)
	assert.Equal(t, *result.Indexes[0].DistanceMetric, "euclidean")
	assert.Equal(t, *result.Indexes[0].Status, "enable")
	assert.Equal(t, *result.Indexes[0].Mode, "fusion")
	assert.Equal(t, *result.Indexes[0].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[1].Mode, "standard")
	assert.Equal(t, *result.Indexes[1].IndexName, "index2")
	assert.Equal(t, *result.Indexes[1].VectorBucketName, "oss-vector-bucket1")
	assert.Equal(t, *result.Indexes[1].Status, "enable")
	assert.Equal(t, *result.Indexes[1].CreateTime, time.Date(2025, time.August, 20, 10, 49, 17, 289372919, time.UTC))

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body = `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &ListVectorIndexesResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyJsonStyle)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_DeleteVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *DeleteVectorIndexRequest
	var input *oss.OperationInput
	var err error

	request = &DeleteVectorIndexRequest{}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &DeleteVectorIndexRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &DeleteVectorIndexRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("demo"),
	}
	input = &oss.OperationInput{
		OpName: "DeleteVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"deleteVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"DeleteVectorIndex"})
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Parameters["DeleteVectorIndex"], "")
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"demo\"}")
}

func TestUnmarshalOutput_DeleteVectorIndex(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Body:       io.NopCloser(bytes.NewReader([]byte(nil))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &DeleteVectorIndexResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}

func TestMarshalInput_PutVectorIndexFusion(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var request *PutVectorIndexFusionRequest
	var input *oss.OperationInput
	var err error

	request = &PutVectorIndexFusionRequest{}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket")

	request = &PutVectorIndexFusionRequest{
		Bucket: oss.Ptr("oss-demo"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, IndexName")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Mode")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		IndexName: oss.Ptr("exampleIndex"),
		Mode:      oss.Ptr("fusion"),
	}
	input = &oss.OperationInput{
		OpName: "PutVectorIndex",
		Method: "POST",
		Headers: map[string]string{
			oss.HTTPHeaderContentType: contentTypeJSON,
		},
		Parameters: map[string]string{
			"putVectorIndex": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, SchemaConfiguration")

	request = &PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr("oss-demo"),
		Mode:      oss.Ptr("fusion"),
		IndexName: oss.Ptr("exampleIndex"),
		SchemaConfiguration: map[string]any{
			"fields": []any{
				map[string]any{
					"name":           "vector_1",
					"type":           "vector",
					"dataType":       "float32",
					"dimension":      1024,
					"distanceMetric": "euclidean",
				},
				map[string]any{
					"name":           "vector_2",
					"type":           "vector",
					"dataType":       "float32",
					"dimension":      512,
					"distanceMetric": "cosine",
				},
				map[string]any{
					"name":    "timestamps",
					"type":    "long",
					"isArray": true,
				},
				map[string]any{
					"name": "price",
					"type": "double",
				},
				map[string]any{
					"name": "ip",
					"type": "ip",
				},
				map[string]any{
					"name": "location",
					"type": "geoPoint",
				},
				map[string]any{
					"name": "tag",
					"type": "string",
				},
				map[string]any{
					"name":           "user_id",
					"type":           "string",
					"isPartitionKey": true,
				},
				map[string]any{
					"name":    "tags",
					"type":    "string",
					"isArray": true,
				},
				map[string]any{
					"name":       "title_1",
					"type":       "string",
					"exactMatch": true,
					"text": map[string]any{
						"enabled":  true,
						"analyzer": "standard",
						"analyzerParameters": map[string]any{
							"caseSensitive": true,
							"delimitWord":   false,
						},
					},
				},
				map[string]any{
					"name":       "title_2",
					"type":       "string",
					"exactMatch": false,
					"text": map[string]any{
						"enabled":  true,
						"analyzer": "split",
						"analyzerParameters": map[string]any{
							"caseSensitive": true,
							"delimiter":     " ",
						},
					},
				},
			},
		},
	}
	err = c.marshalInputJson(request, input, oss.MarshalUpdateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, input.Method, "POST")
	assert.Equal(t, *input.Bucket, "oss-demo")
	assert.Equal(t, input.Parameters["putVectorIndex"], "")
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "{\"indexName\":\"exampleIndex\",\"mode\":\"fusion\",\"schemaConfiguration\":{\"fields\":[{\"dataType\":\"float32\",\"dimension\":1024,\"distanceMetric\":\"euclidean\",\"name\":\"vector_1\",\"type\":\"vector\"},{\"dataType\":\"float32\",\"dimension\":512,\"distanceMetric\":\"cosine\",\"name\":\"vector_2\",\"type\":\"vector\"},{\"isArray\":true,\"name\":\"timestamps\",\"type\":\"long\"},{\"name\":\"price\",\"type\":\"double\"},{\"name\":\"ip\",\"type\":\"ip\"},{\"name\":\"location\",\"type\":\"geoPoint\"},{\"name\":\"tag\",\"type\":\"string\"},{\"isPartitionKey\":true,\"name\":\"user_id\",\"type\":\"string\"},{\"isArray\":true,\"name\":\"tags\",\"type\":\"string\"},{\"exactMatch\":true,\"name\":\"title_1\",\"text\":{\"analyzer\":\"standard\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimitWord\":false},\"enabled\":true},\"type\":\"string\"},{\"exactMatch\":false,\"name\":\"title_2\",\"text\":{\"analyzer\":\"split\",\"analyzerParameters\":{\"caseSensitive\":true,\"delimiter\":\" \"},\"enabled\":true},\"type\":\"string\"}]}}")
}

func TestUnmarshalOutput_PutVectorIndexFusion(t *testing.T) {
	c := VectorsClient{}
	assert.NotNil(t, c)
	var output *oss.OperationOutput
	var err error
	output = &oss.OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	output = &oss.OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
	output = &oss.OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")

	body := `{
  "Error": {
    "Code": "AccessDenied",
    "Message": "AccessDenied",
    "RequestId": "568D5566F2D0F89F5C0E****",
    "HostId": "test.oss.aliyuncs.com"
  }
}`
	output = &oss.OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/json"},
		},
	}
	result = &PutVectorIndexFusionResult{}
	err = c.unmarshalOutput(result, output, oss.UnmarshalDiscardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/json")
}
