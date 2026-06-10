//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	request := &PutObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPrivate,
	}
	result, err := client.PutObjectAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	infoRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.HeadObject(context.TODO(), infoRequest)
	assert.Nil(t, err)
	_, err = client.PutObjectAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	request = &PutObjectAclRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPrivate,
	}
	_, err = client.PutObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPrivate,
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	assert.Nil(t, err)
	request := &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, ObjectACLType(*result.ACL), ObjectACLPrivate)
	assert.NotEmpty(t, *result.Owner.ID)
	assert.NotEmpty(t, *result.Owner.DisplayName)

	_, err = client.GetObjectAcl(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	objectNameNotExist := objectName + "-not-exist"
	request = &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.GetObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)
	assert.Equal(t, "0026-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
