package oss

import "context"

func (c *Client) InvokeOperation(ctx context.Context, input *OperationInput, optFns ...func(*Options)) (output *OperationOutput, err error) {
	return c.invokeOperation(ctx, input, optFns)
}
