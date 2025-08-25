package oss

import (
	"context"
)


// GetPublicAccessBlock Queries the Block Public Access configurations of vector.
func (c *VectorsClient) GetPublicAccessBlock(ctx context.Context, request *GetPublicAccessBlockRequest, optFns ...func(*Options)) (*GetPublicAccessBlockResult, error) {
	return c.client.GetPublicAccessBlock(ctx, request, optFns...)
}

// PutPublicAccessBlock Enables or disables Block Public Access for vector.
func (c *VectorsClient) PutPublicAccessBlock(ctx context.Context, request *PutPublicAccessBlockRequest, optFns ...func(*Options)) (*PutPublicAccessBlockResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers[HTTPHeaderContentType] = contentTypeJSON
	return c.client.PutPublicAccessBlock(ctx, request, optFns...)
}

// DeletePublicAccessBlock Deletes the Block Public Access configurations of vector .
func (c *VectorsClient) DeletePublicAccessBlock(ctx context.Context, request *DeletePublicAccessBlockRequest, optFns ...func(*Options)) (*DeletePublicAccessBlockResult, error) {
	return c.client.DeletePublicAccessBlock(ctx, request, optFns...)
}
