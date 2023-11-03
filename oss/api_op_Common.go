package oss

import (
	"context"
)

func (c *Client) InvokeOperation(ctx context.Context, input *OperationInput, optFns ...func(*Options)) (output *OperationOutput, err error) {
	if input == nil {
		return nil, NewErrParamNull("OperationInput")
	}

	if input.Bucket != nil && !isValidBucketName(input.Bucket) {
		return nil, NewErrParamInvalid("OperationInput.Bucket")
	}

	if input.Key != nil && !isValidObjectName(input.Key) {
		return nil, NewErrParamInvalid("OperationInput.Key")
	}

	return c.invokeOperation(ctx, input, optFns)
}
