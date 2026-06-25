package dataprocess

import (
	"context"
	"fmt"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

type PaginatorOptions struct {
	// The maximum number of items in the response.
	Limit *int64
}

// ListDatasetsPaginator is a paginator for ListDatasets
type ListDatasetsPaginator struct {
	options     PaginatorOptions
	client      *Client
	request     *ListDatasetsRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *Client) NewListDatasetsPaginator(request *ListDatasetsRequest, optFns ...func(*PaginatorOptions)) *ListDatasetsPaginator {
	if request == nil {
		request = &ListDatasetsRequest{}
	}

	options := PaginatorOptions{}
	options.Limit = request.MaxResults

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListDatasetsPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.NextToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListDatasetsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListDatasets page.
func (p *ListDatasetsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListDatasetsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.NextToken = p.nextToken

	var limit *int64
	if oss.ToInt64(p.options.Limit) > 0 {
		limit = p.options.Limit
	}
	request.MaxResults = limit

	result, err := p.client.ListDatasets(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = result.NextToken != nil
	p.nextToken = result.NextToken

	return result, nil
}

// ListSmartClustersPaginator is a paginator for ListSmartClusters
type ListSmartClustersPaginator struct {
	options     PaginatorOptions
	client      *Client
	request     *ListSmartClustersRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *Client) NewListSmartClustersPaginator(request *ListSmartClustersRequest, optFns ...func(*PaginatorOptions)) *ListSmartClustersPaginator {
	if request == nil {
		request = &ListSmartClustersRequest{}
	}

	options := PaginatorOptions{}
	options.Limit = request.MaxResults

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListSmartClustersPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.NextToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListSmartClustersPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListSmartClusters page.
func (p *ListSmartClustersPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListSmartClustersResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.NextToken = p.nextToken

	var limit *int64
	if oss.ToInt64(p.options.Limit) > 0 {
		limit = p.options.Limit
	}
	request.MaxResults = limit

	result, err := p.client.ListSmartClusters(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = result.NextToken != nil
	p.nextToken = result.NextToken

	return result, nil
}

// ListDataPipelineConfigurationsPaginator is a paginator for ListDataPipelineConfigurations
type ListDataPipelineConfigurationsPaginator struct {
	options     PaginatorOptions
	client      *Client
	request     *ListDataPipelineConfigurationsRequest
	nextToken   *string
	firstPage   bool
	isTruncated bool
}

func (c *Client) NewListDataPipelineConfigurationsPaginator(request *ListDataPipelineConfigurationsRequest, optFns ...func(*PaginatorOptions)) *ListDataPipelineConfigurationsPaginator {
	if request == nil {
		request = &ListDataPipelineConfigurationsRequest{}
	}

	options := PaginatorOptions{}
	options.Limit = request.MaxResults

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListDataPipelineConfigurationsPaginator{
		options:     options,
		client:      c,
		request:     request,
		nextToken:   request.NextToken,
		firstPage:   true,
		isTruncated: false,
	}
}

// HasNext Returns true if there’s a next page.
func (p *ListDataPipelineConfigurationsPaginator) HasNext() bool {
	return p.firstPage || p.isTruncated
}

// NextPage retrieves the next ListDataPipelineConfigurations page.
func (p *ListDataPipelineConfigurationsPaginator) NextPage(ctx context.Context, optFns ...func(*oss.Options)) (*ListDataPipelineConfigurationsResult, error) {
	if !p.HasNext() {
		return nil, fmt.Errorf("no more pages available")
	}

	request := *p.request
	request.NextToken = p.nextToken

	var limit *int64
	if oss.ToInt64(p.options.Limit) > 0 {
		limit = p.options.Limit
	}
	request.MaxResults = limit

	result, err := p.client.ListDataPipelineConfigurations(ctx, &request, optFns...)
	if err != nil {
		return nil, err
	}

	p.firstPage = false
	p.isTruncated = oss.ToString(result.NextToken) != ""
	p.nextToken = result.NextToken

	return result, nil
}
