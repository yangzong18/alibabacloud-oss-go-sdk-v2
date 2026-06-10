//go:build integration

package oss

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPutBucketVersioning(t *testing.T) {
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

	request := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err := client.PutBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	request = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err = client.PutBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	_, err = client.PutBucketVersioning(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketNameNotExist),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	result, err = client.PutBucketVersioning(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketVersioning(t *testing.T) {
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

	request := &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Nil(t, result.VersionStatus)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Enabled")

	versionRequest = &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionSuspended,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Suspended")

	_, err = client.GetBucketVersioning(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetBucketVersioningRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	result, err = client.GetBucketVersioning(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListObjectVersions(t *testing.T) {
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

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	request := &GetBucketVersioningRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketVersioning(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionStatus, "Enabled")

	// put object v1
	content1 := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6) + "\v\f\n"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content1),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionIdV1 := putObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(versionIdV1) > 0)

	// put object v2
	content2 := randLowStr(200)
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(content2),
	}
	putObjResult, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionIdV2 := putObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(versionIdV2) > 0)
	assert.NotEqual(t, versionIdV1, versionIdV2)

	delObjRequest := &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	delObjResult, err := client.DeleteObject(context.TODO(), delObjRequest)
	assert.Nil(t, err)
	assert.True(t, delObjResult.DeleteMarker)
	markVersionId := delObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(markVersionId) > 0)

	delObjRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	delObjResult, err = client.DeleteObject(context.TODO(), delObjRequest)
	assert.Nil(t, err)
	assert.True(t, delObjResult.DeleteMarker)
	markVersionIdAgain := delObjResult.Headers.Get("x-oss-version-id")
	assert.True(t, len(markVersionIdAgain) > 0)
	assert.NotEqual(t, markVersionId, markVersionIdAgain)

	versions := &ListObjectVersionsRequest{
		Bucket: Ptr(bucketName),
	}
	versionsResult, err := client.ListObjectVersions(context.TODO(), versions)
	assert.Nil(t, err)
	assert.Len(t, versionsResult.ObjectDeleteMarkers, 2)
	assert.Len(t, versionsResult.ObjectVersions, 2)

	versions = &ListObjectVersionsRequest{
		Bucket: Ptr(bucketName),
		IsMix:  true,
	}
	versionsResult, err = client.ListObjectVersions(context.TODO(), versions)
	assert.Nil(t, err)
	assert.Len(t, versionsResult.ObjectVersionsDeleteMarkers, 4)

	_, err = client.ListObjectVersions(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	versions = &ListObjectVersionsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	versionsResult, err = client.ListObjectVersions(context.TODO(), versions)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
