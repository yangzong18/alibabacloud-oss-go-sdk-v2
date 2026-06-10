package dataprocess

import (
	"context"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type DeleteFileMetaRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	DatasetName *string `input:"query,datasetName,required"`

	Uri *string `input:"query,uri,required"`

	oss.RequestCommon
}

type DeleteFileMetaResult struct {
	oss.ResultCommon
}

// DeleteFileMeta Deletes the metadata of a specified object.
func (c *Client) DeleteFileMeta(ctx context.Context, request *DeleteFileMetaRequest, optFns ...func(*oss.Options)) (*DeleteFileMetaResult, error) {
	var err error
	if request == nil {
		request = &DeleteFileMetaRequest{}
	}
	input := &oss.OperationInput{
		OpName: "DeleteFileMeta",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/xml",
		},
		Parameters: map[string]string{
			"metaQuery": "",
			"action":    "deleteFileMeta",
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

	result := &DeleteFileMetaResult{}

	if err = c.client.UnmarshalOutput(result, output); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
