package oss

import (
	"context"
	"fmt"
)

type PaginatorOptions struct {
	// The maximum number of items in the response.
	Limit int32
}

// ListObjectsPaginator is a paginator for ListObjects
type ListObjectsPaginator struct {
	options     PaginatorOptions
	client      *Client
	request     *ListObjectsRequest
	marker      *string
	firstPage   bool
	isTruncated bool
}

func (c *Client) NewListObjectsPaginator(request *ListObjectsRequest, optFns ...func(*PaginatorOptions)) *ListObjectsPaginator {
	if request == nil {
		request = &ListObjectsRequest{}
	}

	options := PaginatorOptions{}
	options.Limit = request.MaxKeys

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListObjectsPaginator{
		options:     options,
		client:      c,
		request:     request,
		marker:      request.Marker,
		firstPage:   true,
		isTruncated: false,
	}
}

// Returns true if there’s a next page.
func (p *ListObjectsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListObjects page.
func (p *ListObjectsPaginator) NextPage(ctx context.Context, optFns ...func(*Options)) (*ListObjectsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.Marker = p.marker

	var limit int32
	if p.options.Limit > 0 {
		limit = p.options.Limit
	}
	request.MaxKeys = limit
	request.EncodingType = Ptr("url")

	result, err := p.client.ListObjects(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = result.IsTruncated
	p.marker = result.NextMarker

	return result, nil
}
