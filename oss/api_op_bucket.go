package oss

import (
	"context"
	"encoding/xml"
	"net/url"
	"strings"
	"time"
)

type PutBucketRequest struct {
	// The name of the bucket to create.
	Bucket *string `input:"host,bucket,required"`

	// The access control list (ACL) of the bucket.
	Acl BucketACLType `input:"header,x-oss-acl"`

	// The ID of the resource group.
	ResourceGroupId *string `input:"header,x-oss-resource-group-id"`

	// The configuration information for the bucket.
	CreateBucketConfiguration *CreateBucketConfiguration `input:"body,CreateBucketConfiguration,xml"`

	RequestCommon
}

type CreateBucketConfiguration struct {
	XMLName xml.Name `xml:"CreateBucketConfiguration"`

	// The storage class of the bucket.
	StorageClass StorageClassType `xml:"StorageClass"`

	// The redundancy type of the bucket.
	DataRedundancyType DataRedundancyType `xml:"DataRedundancyType"`
}

type PutBucketResult struct {
	ResultCommon
}

// Creates a bucket.
func (c *Client) PutBucket(ctx context.Context, request *PutBucketRequest, optFns ...func(*Options)) (*PutBucketResult, error) {
	var err error
	if request == nil {
		request = &PutBucketRequest{}
	}
	input := &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutBucketResult{}

	if err = c.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

// Deletes creates a bucket.
type DeleteBucketRequest struct {
	// The name of the bucket to delete.
	Bucket *string `input:"host,bucket,required"`

	RequestCommon
}

type DeleteBucketResult struct {
	ResultCommon
}

// Deletes a bucket.
func (c *Client) DeleteBucket(ctx context.Context, request *DeleteBucketRequest, optFns ...func(*Options)) (*DeleteBucketResult, error) {
	var err error
	if request == nil {
		request = &DeleteBucketRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteBucket",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &DeleteBucketResult{}
	if err = c.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type ListObjectsRequest struct {
	// The name of the bucket containing the objects
	Bucket *string `input:"host,bucket,required"`

	// The character that is used to group objects by name. If you specify the delimiter parameter in the request,
	// the response contains the CommonPrefixes parameter. The objects whose names contain the same string from
	// the prefix to the next occurrence of the delimiter are grouped as a single result element in CommonPrefixes.
	Delimiter *string `input:"query,delimiter"`

	// The encoding type of the content in the response. Valid value: url
	EncodingType *string `input:"query,encoding-type"`

	// The name of the object after which the ListObjects (GetBucket) operation starts.
	// If this parameter is specified, objects whose names are alphabetically greater than the marker value are returned.
	Marker *string `input:"query,marker"`

	// The maximum number of objects that you want to return. If the list operation cannot be complete at a time
	// because the max-keys parameter is specified, the NextMarker element is included in the response as the marker
	// for the next list operation.
	MaxKeys int32 `input:"query,max-keys"`

	// The prefix that the names of the returned objects must contain.
	Prefix *string `input:"query,prefix"`

	RequestCommon
}

type ListObjectsResult struct {
	// The name of the bucket.
	Name *string `xml:"Name"`

	// The prefix contained in the returned object names.
	Prefix *string `xml:"Prefix"`

	// The name of the object after which the list operation begins.
	Marker *string `xml:"Marker"`

	// The maximum number of returned objects in the response.
	MaxKeys int32 `xml:"MaxKeys"`

	// The character that is used to group objects by name.
	Delimiter *string `xml:"Delimiter"`

	// Indicates whether the returned results are truncated.
	// true indicates that not all results are returned this time.
	// false indicates that all results are returned this time.
	IsTruncated bool `xml:"IsTruncated"`

	// The position from which the next list operation starts.
	NextMarker *string `xml:"NextMarker"`

	// The encoding type of the content in the response.
	EncodingType *string `xml:"EncodingType"`

	// The container that stores the metadata of the returned objects.
	Contents []ObjectProperties `xml:"Contents"`

	// If the Delimiter parameter is specified in the request, the response contains the CommonPrefixes element.
	CommonPrefixes []CommonPrefix `xml:"CommonPrefixes"`

	ResultCommon
}

type ObjectProperties struct {
	// The name of the object.
	Key string `xml:"Key"`

	// The type of the object. Valid values: Normal, Multipart and Appendable
	Type string `xml:"Type"`

	// The size of the returned object. Unit: bytes.
	Size int64 `xml:"Size"`

	// The entity tag (ETag). An ETag is created when an object is created to identify the content of the object.
	ETag string `xml:"ETag"`

	// The time when the returned objects were last modified.
	LastModified time.Time `xml:"LastModified"`

	// The storage class of the object.
	StorageClass string `xml:"StorageClass"`

	// The container that stores information about the bucket owner.
	Owner *Owner `xml:"Owner"`

	// The restoration status of the object.
	RestoreInfo *string `xml:"RestoreInfo"`
}

type Owner struct {
	// The ID of the bucket owner.
	ID string `xml:"ID"`

	// The name of the object owner.
	DisplayName string `xml:"DisplayName"`
}

type CommonPrefix struct {
	// The prefix contained in the returned object names.
	Prefix string `xml:"Prefix"`
}

// Queries information about objects in a bucket.
func (c *Client) ListObjects(ctx context.Context, request *ListObjectsRequest, optFns ...func(*Options)) (*ListObjectsResult, error) {
	var err error
	if request == nil {
		request = &ListObjectsRequest{}
	}
	input := &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &ListObjectsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

// private function
func unmarshalEncodeType(result interface{}, output *OperationOutput) error {
	switch r := result.(type) {
	case *ListObjectsResult:
		if r.EncodingType != nil && strings.EqualFold(*r.EncodingType, "url") {
			fields := []**string{&r.Prefix, &r.Marker, &r.Delimiter, &r.NextMarker}
			var s string
			var err error
			for _, pp := range fields {
				if pp != nil && *pp != nil {
					if s, err = url.QueryUnescape(**pp); err != nil {
						return err
					}
					*pp = Ptr(s)
				}
			}
			for i := 0; i < len(r.Contents); i++ {
				if r.Contents[i].Key, err = url.QueryUnescape(r.Contents[i].Key); err != nil {
					return err
				}
			}
			for i := 0; i < len(r.CommonPrefixes); i++ {
				if r.CommonPrefixes[i].Prefix, err = url.QueryUnescape(r.CommonPrefixes[i].Prefix); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
