package oss

import (
	"context"
	"io"
	"time"
)

type PutObjectRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// The caching behavior of the web page when the object is downloaded.
	CacheControl *string `input:"header,Cache-Control"`

	// The method that is used to access the object.
	ContentDisposition *string `input:"header,Content-Disposition"`

	// The method that is used to encode the object.
	ContentEncoding *string `input:"header,Content-Encoding"`

	// The size of the data in the HTTP message body. Unit: bytes.
	ContentLength *int64 `input:"header,Content-Length"`

	// The MD5 hash of the object that you want to upload.
	ContentMD5 *string `input:"header,Content-MD5"`

	// A standard MIME type describing the format of the contents.
	ContentType *string `input:"header,Content-Type"`

	// The expiration time of the cache in UTC.
	Expires *string `input:"header,Expires"`

	// Specifies whether the object that is uploaded by calling the PutObject operation
	// overwrites an existing object that has the same name. Valid values: true and false
	ForbidOverwrite *string `input:"header,x-oss-forbid-overwrite"`

	// The encryption method on the server side when an object is created.
	// Valid values: AES256 and KMS
	ServerSideEncryption *string `input:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	//This header is valid only when the x-oss-server-side-encryption header is set to KMS.
	ServerSideDataEncryption *string `input:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `input:"header,x-oss-server-side-encryption-key-id"`

	// The access control list (ACL) of the object.
	Acl ObjectACLType `input:"header,x-oss-object-acl"`

	// The storage class of the object.
	StorageClass StorageClassType `input:"header,x-oss-storage-class"`

	// The metadata of the object that you want to upload.
	Metadata map[string]string `input:"header,x-oss-meta-,usermeta"`

	// The tags that are specified for the object by using a key-value pair.
	// You can specify multiple tags for an object. Example: TagA=A&TagB=B.
	Tagging *string `input:"header,x-oss-tagging"`

	RequestCommon
}

type PutObjectResult struct {
	// Content-Md5 for the uploaded object.
	ContentMD5 *string `output:"header,Content-MD5"`

	// Entity tag for the uploaded object.
	ETag *string `output:"header,ETag"`

	// The 64-bit CRC value of the object.
	// This value is calculated based on the ECMA-182 standard.
	HashCRC64 *string `output:"header,x-oss-hash-crc64ecma"`

	// Version of the object.
	VersionId *string `output:"header,x-oss-version-id"`

	ResultCommon
}

// PutObject Uploads a object.
func (c *Client) PutObject(ctx context.Context, request *PutObjectRequest, optFns ...func(*Options)) (*PutObjectResult, error) {
	var err error
	if request == nil {
		request = &PutObjectRequest{}
	}
	input := &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	if err = c.marshalInput(request, input); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &PutObjectResult{}
	if err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type GetObjectRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// If the ETag specified in the request matches the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	IfMatch *string `input:"header,If-Match"`

	// If the ETag specified in the request does not match the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	IfNoneMatch *string `input:"header,If-None-Match"`

	// If the time specified in this header is earlier than the object modified time or is invalid,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfModifiedSince *string `input:"header,If-Modified-Since"`

	// If the time specified in this header is the same as or later than the object modified time,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfUnmodifiedSince *string `input:"header,If-Unmodified-Since"`

	// The content range of the object to be returned.
	// If the value of Range is valid, the total size of the object and the content range are returned.
	// For example, Content-Range: bytes 0~9/44 indicates that the total size of the object is 44 bytes,
	// and the range of data returned is the first 10 bytes.
	// However, if the value of Range is invalid, the entire object is returned,
	// and the response does not include the Content-Range parameter.
	Range *string `input:"header,Range"`

	// Specify standard behaviors to download data by range
	// If the value is "standard", the download behavior is modified when the specified range is not within the valid range.
	// For an object whose size is 1,000 bytes:
	// 1) If you set Range: bytes to 500-2000, the value at the end of the range is invalid.
	// In this case, OSS returns HTTP status code 206 and the data that is within the range of byte 500 to byte 999.
	// 2) If you set Range: bytes to 1000-2000, the value at the start of the range is invalid.
	// In this case, OSS returns HTTP status code 416 and the InvalidRange error code.
	RangeBehavior *string `input:"header,x-oss-range-behavior"`

	// The cache-control header to be returned in the response.
	ResponseCacheControl *string `input:"query,response-cache-control"`

	// The content-disposition header to be returned in the response.
	ResponseContentDisposition *string `input:"query,response-content-disposition"`

	// The content-encoding header to be returned in the response.
	ResponseContentEncoding *string `input:"query,response-content-encoding"`

	// The content-language header to be returned in the response.
	ResponseContentLanguage *string `input:"query,response-content-language"`

	// The content-type header to be returned in the response.
	ResponseContentType *string `input:"query,response-content-type"`

	// The expires header to be returned in the response.
	ResponseExpires *string `input:"query,response-expires"`

	// VersionId used to reference a specific version of the object.
	VersionId *string `input:"query,versionId"`

	RequestCommon
}

type GetObjectResult struct {
	// Size of the body in bytes. -1 indicates that the Content-Length dose not exist.
	ContentLength int64 `output:"header,Content-Length"`

	// The portion of the object returned in the response.
	ContentRange *string `output:"header,Content-Range"`

	// A standard MIME type describing the format of the object data.
	ContentType *string `output:"header,Content-Type"`

	// The entity tag (ETag). An ETag is created when an object is created to identify the content of the object.
	ETag *string `output:"header,ETag"`

	// The time when the returned objects were last modified.
	LastModified *time.Time `output:"header,Last-Modified,time"`

	// The storage class of the object.
	StorageClass *string `output:"header,x-oss-storage-class"`

	// Content-Md5 for the uploaded object.
	ContentMD5 *string `output:"header,Content-MD5"`

	// A map of metadata to store with the object.
	Metadata map[string]string `output:"header,x-oss-meta-,usermeta"`

	// If the requested object is encrypted by using a server-side encryption algorithm based on entropy encoding,
	// OSS automatically decrypts the object and returns the decrypted object after OSS receives the GetObject request.
	// The x-oss-server-side-encryption header is included in the response to indicate
	// the encryption algorithm used to encrypt the object on the server.
	ServerSideEncryption *string `output:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	ServerSideDataEncryption *string `output:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `output:"header,x-oss-server-side-encryption-key-id"`

	// The type of the object.
	ObjectType *string `output:"header,x-oss-object-type"`

	// The position for the next append operation.
	// If the type of the object is Appendable, this header is included in the response.
	NextAppendPosition *string `output:"header,x-oss-next-append-position"`

	// The 64-bit CRC value of the object.
	// This value is calculated based on the ECMA-182 standard.
	HashCRC64 *string `output:"header,x-oss-hash-crc64ecma"`

	// The lifecycle information about the object.
	// If lifecycle rules are configured for the object, this header is included in the response.
	// This header contains the following parameters: expiry-date that indicates the expiration time of the object,
	// and rule-id that indicates the ID of the matched lifecycle rule.
	Expiration *string `output:"header,x-oss-expiration"`

	// The status of the object when you restore an object.
	// If the storage class of the bucket is Archive and a RestoreObject request is submitted,
	Restore *string `output:"header,x-oss-restore"`

	// The result of an event notification that is triggered for the object.
	ProcessStatus *string `output:"header,x-oss-process-status"`

	// The number of tags added to the object.
	// This header is included in the response only when you have read permissions on tags.
	TaggingCount int32 `output:"header,x-oss-tagging-count"`

	// Specifies whether the object retrieved was (true) or was not (false) a Delete  Marker.
	DeleteMarker bool `output:"header,x-oss-delete-marker"`

	// Version of the object.
	VersionId *string `output:"header,x-oss-version-id"`

	// Object data.
	Body io.ReadCloser

	ResultCommon
}

func (c *Client) GetObject(ctx context.Context, request *GetObjectRequest, optFns ...func(*Options)) (*GetObjectResult, error) {
	var err error
	if request == nil {
		request = &GetObjectRequest{}
	}
	input := &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	if err = c.marshalInput(request, input); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &GetObjectResult{
		Body: output.Body,
	}
	if err = c.unmarshalOutput(result, output, unmarshalHeader); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type CopyObjectRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	//The path of the source object.
	Source *string `input:"header,x-oss-copy-source,required"`

	// Specifies whether the object that is uploaded by calling the CopyObject operation
	// overwrites an existing object that has the same name. Valid values: true and false
	ForbidOverwrite *string `input:"header,x-oss-forbid-overwrite"`

	// If the ETag specified in the request matches the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	IfMatch *string `input:"header,x-oss-copy-source-if-match"`

	// If the ETag specified in the request does not match the ETag value of the object,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	IfNoneMatch *string `input:"header,x-oss-copy-source-if-none-match"`

	// If the time specified in this header is earlier than the object modified time or is invalid,
	// the object and 200 OK are returned. Otherwise, 304 Not Modified is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfModifiedSince *string `input:"header,x-oss-copy-source-if-modified-since"`

	// If the time specified in this header is the same as or later than the object modified time,
	// the object and 200 OK are returned. Otherwise, 412 Precondition Failed is returned.
	// The time must be in GMT. Example: Fri, 13 Nov 2015 14:47:53 GMT.
	IfUnmodifiedSince *string `input:"header,x-oss-copy-source-if-unmodified-since"`

	// The method that is used to configure the metadata of the destination object.
	// COPY (default): The metadata of the source object is copied to the destination object.
	// The configurations of the x-oss-server-side-encryption
	// header of the source object are not copied to the destination object.
	// The x-oss-server-side-encryption header in the CopyObject request specifies
	// the method used to encrypt the destination object.
	// REPLACE: The metadata specified in the request is used as the metadata of the destination object.
	MetadataDirective *string `input:"header,x-oss-metadata-directive"`

	// The encryption method on the server side when an object is created.
	// Valid values: AES256 and KMS
	ServerSideEncryption *string `input:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	//This header is valid only when the x-oss-server-side-encryption header is set to KMS.
	ServerSideDataEncryption *string `input:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `input:"header,x-oss-server-side-encryption-key-id"`

	// The access control list (ACL) of the object.
	Acl ObjectACLType `input:"header,x-oss-object-acl"`

	// The storage class of the object.
	StorageClass StorageClassType `input:"header,x-oss-storage-class"`

	// The metadata of the object that you want to upload.
	Metadata map[string]string `input:"header,x-oss-meta-,usermeta"`

	// The tags that are specified for the object by using a key-value pair.
	// You can specify multiple tags for an object. Example: TagA=A&TagB=B.
	Tagging *string `input:"header,x-oss-tagging"`

	// The method that is used to configure tags for the destination object.
	// Valid values: Copy (default): The tags of the source object are copied to the destination object.
	// Replace: The tags specified in the request are configured for the destination object.
	TaggingDirective *string `input:"header,x-oss-tagging-directive"`

	// The version ID of the source object.
	VersionId *string `input:"query,x-oss-copy-source-version-id"`

	RequestCommon
}

type CopyObjectResult struct {
	// The 64-bit CRC value of the object.
	// This value is calculated based on the ECMA-182 standard.
	HashCRC64 *string `output:"header,x-oss-hash-crc64ecma"`

	// Version of the object.
	VersionId *string `output:"header,x-oss-version-id"`

	// The version ID of the source object.
	SourceVersionId *string `output:"header,x-oss-copy-source-version-id"`

	// If the requested object is encrypted by using a server-side encryption algorithm based on entropy encoding,
	// OSS automatically decrypts the object and returns the decrypted object after OSS receives the GetObject request.
	// The x-oss-server-side-encryption header is included in the response to indicate
	// the encryption algorithm used to encrypt the object on the server.
	ServerSideEncryption *string `output:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	ServerSideDataEncryption *string `output:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `output:"header,x-oss-server-side-encryption-key-id"`

	// The time when the returned objects were last modified.
	LastModified *time.Time `xml:"LastModified"`

	// The entity tag (ETag). An ETag is created when an object is created to identify the content of the object.
	ETag *string `xml:"ETag"`

	ResultCommon
}

// CopyObject Copies objects within a bucket or between buckets in the same region
func (c *Client) CopyObject(ctx context.Context, request *CopyObjectRequest, optFns ...func(*Options)) (*CopyObjectResult, error) {
	var err error
	if request == nil {
		request = &CopyObjectRequest{}
	}
	input := &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &CopyObjectResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type AppendObjectRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// The position from which the AppendObject operation starts.
	// Each time an AppendObject operation succeeds, the x-oss-next-append-position header is included in
	// the response to specify the position from which the next AppendObject operation starts.
	Position *int64 `input:"query,position,required"`

	// The caching behavior of the web page when the object is downloaded.
	CacheControl *string `input:"header,Cache-Control"`

	// The method that is used to access the object.
	ContentDisposition *string `input:"header,Content-Disposition"`

	// The method that is used to encode the object.
	ContentEncoding *string `input:"header,Content-Encoding"`

	// The size of the data in the HTTP message body. Unit: bytes.
	ContentLength *int64 `input:"header,Content-Length"`

	// The MD5 hash of the object that you want to upload.
	ContentMD5 *string `input:"header,Content-MD5"`

	// The expiration time of the cache in UTC.
	Expires *string `input:"header,Expires"`

	// Specifies whether the object that is uploaded by calling the PutObject operation
	// overwrites an existing object that has the same name. Valid values: true and false
	ForbidOverwrite *string `input:"header,x-oss-forbid-overwrite"`

	// The encryption method on the server side when an object is created.
	// Valid values: AES256 and KMS
	ServerSideEncryption *string `input:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	//This header is valid only when the x-oss-server-side-encryption header is set to KMS.
	ServerSideDataEncryption *string `input:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `input:"header,x-oss-server-side-encryption-key-id"`

	// The access control list (ACL) of the object.
	Acl ObjectACLType `input:"header,x-oss-object-acl"`

	// The storage class of the object.
	StorageClass StorageClassType `input:"header,x-oss-storage-class"`

	// The metadata of the object that you want to upload.
	Metadata map[string]string `input:"header,x-oss-meta-,usermeta"`

	// The tags that are specified for the object by using a key-value pair.
	// You can specify multiple tags for an object. Example: TagA=A&TagB=B.
	Tagging *string `input:"header,x-oss-tagging"`

	RequestCommon
}

type AppendObjectResult struct {
	// Version of the object.
	VersionId *string `output:"header,x-oss-version-id"`

	// The 64-bit CRC value of the object.
	// This value is calculated based on the ECMA-182 standard.
	HashCRC64 *string `output:"header,x-oss-hash-crc64ecma"`

	// The position that must be provided in the next request, which is the current length of the object.
	NextPosition int64 `output:"header,x-oss-next-append-position"`

	// The encryption method on the server side when an object is created.
	// Valid values: AES256 and KMS
	ServerSideEncryption *string `output:"header,x-oss-server-side-encryption"`

	//The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	//This header is valid only when the x-oss-server-side-encryption header is set to KMS.
	ServerSideDataEncryption *string `output:"header,x-oss-server-side-data-encryption"`

	// The ID of the customer master key (CMK) that is managed by Key Management Service (KMS).
	SSEKMSKeyId *string `output:"header,x-oss-server-side-encryption-key-id"`

	// A map of metadata to store with the object.
	Metadata map[string]string `output:"header,x-oss-meta-,usermeta"`

	ResultCommon
}

// AppendObject Uploads an object by appending the object to an existing object.
// Objects created by using the AppendObject operation are appendable objects.
func (c *Client) AppendObject(ctx context.Context, request *AppendObjectRequest, optFns ...func(*Options)) (*AppendObjectResult, error) {
	var err error
	if request == nil {
		request = &AppendObjectRequest{}
	}
	input := &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &AppendObjectResult{}
	if err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteObjectRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The name of the object.
	Key *string `input:"path,key,required"`

	// The version ID of the source object.
	VersionId *string `input:"query,versionId"`

	RequestCommon
}

type DeleteObjectResult struct {
	// Version of the object.
	VersionId *string `output:"header,x-oss-version-id"`

	// Specifies whether the object retrieved was (true) or was not (false) a Delete  Marker.
	DeleteMarker bool `output:"header,x-oss-delete-marker"`

	ResultCommon
}

// DeleteObject Deletes an object.
func (c *Client) DeleteObject(ctx context.Context, request *DeleteObjectRequest, optFns ...func(*Options)) (*DeleteObjectResult, error) {
	var err error
	if request == nil {
		request = &DeleteObjectRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}
	result := &DeleteObjectResult{}
	if err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}

type DeleteMultipleObjectsRequest struct {
	// The name of the bucket.
	Bucket *string `input:"host,bucket,required"`

	// The encoding type of the object names in the response. Valid value: url
	EncodingType *string `input:"query,encoding-type"`

	// The size of the data in the HTTP message body. Unit: bytes.
	ContentLength int64 `input:"header,Content-Length"`

	// The MD5 hash of the object that you want to upload.
	ContentMD5 *string `input:"header,Content-MD5"`

	// The container that stores information about you want to delete objects.
	Objects []DeleteObject

	// Specifies whether to enable the Quiet return mode.
	// The DeleteMultipleObjects operation provides the following return modes: Valid value: true,false
	Quiet bool

	RequestCommon
}

type DeleteObject struct {
	// The name of the object that you want to delete.
	Key *string `xml:"Key"`

	// The version ID of the object that you want to delete.
	VersionId *string `xml:"VersionId"`
}

type DeleteMultipleObjectsResult struct {
	// The container that stores information about the deleted objects.
	DeletedObjects []DeletedInfo `xml:"Deleted"`

	// The encoding type of the name of the deleted object in the response.
	// If encoding-type is specified in the request, the object name is encoded in the returned result.
	EncodingType *string `xml:"EncodingType"`

	ResultCommon
}

type DeletedInfo struct {
	// The name of the deleted object.
	Key *string `xml:"Key"`

	// The version ID of the object that you deleted.
	VersionId *string `xml:"VersionId"`

	// Indicates whether the deleted version is a delete marker.
	DeleteMarker bool `xml:"DeleteMarker"`

	// The version ID of the delete marker.
	DeleteMarkerVersionId *string `xml:"DeleteMarkerVersionId"`
}

// DeleteMultipleObjects Deletes multiple objects from a bucket.
func (c *Client) DeleteMultipleObjects(ctx context.Context, request *DeleteMultipleObjectsRequest, optFns ...func(*Options)) (*DeleteMultipleObjectsResult, error) {
	var err error
	if request == nil {
		request = &DeleteMultipleObjectsRequest{}
	}
	input := &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"delete":        "",
			"encoding-type": "url",
		},
		Bucket: request.Bucket,
	}
	if err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}
	result := &DeleteMultipleObjectsResult{}
	if err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}
	return result, err
}
