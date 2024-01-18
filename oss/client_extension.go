package oss

import (
	"context"
	"errors"
	"os"
)

// NewDownloader creates a new Downloader instance to download objects.
func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader {
	return NewDownloader(c, optFns...)
}

// NewUploader creates a new Uploader instance to upload objects.
func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader {
	return NewUploader(c, optFns...)
}

// OpenFile opens the named file for reading.
func (c *Client) OpenFile(ctx context.Context, bucket string, key string, optFns ...func(*OpenOptions)) (*ReadOnlyFile, error) {
	return NewReadOnlyFile(ctx, c, bucket, key, optFns...)
}

// AppendFile opens or creates the named file for appending.
func (c *Client) AppendFile(ctx context.Context, bucket string, key string, optFns ...func(*AppendOptions)) (*AppendOnlyFile, error) {
	return NewAppendFile(ctx, c, bucket, key, optFns...)
}

// IsObjectExist checks if the object exists.
func (c *Client) IsObjectExist(ctx context.Context, bucket string, key string, optFns ...func(*Options)) (bool, error) {
	_, err := c.GetObjectMeta(ctx, &GetObjectMetaRequest{Bucket: Ptr(bucket), Key: Ptr(key)}, optFns...)
	if err == nil {
		return true, nil
	}
	var serr *ServiceError
	errors.As(err, &serr)
	if errors.As(err, &serr) {
		if serr.Code == "NoSuchKey" ||
			// error code not in response header
			(serr.StatusCode == 404 && serr.Code == "BadErrorResponse") {
			return false, nil
		}
	}
	return false, err
}

// IsBucketExist checks if the bucket exists.
func (c *Client) IsBucketExist(ctx context.Context, bucket string, optFns ...func(*Options)) (bool, error) {
	_, err := c.GetBucketAcl(ctx, &GetBucketAclRequest{Bucket: Ptr(bucket)}, optFns...)
	if err == nil {
		return true, nil
	}
	var serr *ServiceError
	if errors.As(err, &serr) {
		if serr.Code == "NoSuchBucket" {
			return false, nil
		}
		return true, nil
	}
	return false, err
}

// PutObjectFromFile creates a new object from the local file.
func (c *Client) PutObjectFromFile(ctx context.Context, request *PutObjectRequest, filePath string, optFns ...func(*Options)) (*PutObjectResult, error) {
	if request == nil {
		return nil, NewErrParamNull("request")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	pRequest := *request
	pRequest.Body = file
	return c.PutObject(ctx, &pRequest, optFns...)
}
