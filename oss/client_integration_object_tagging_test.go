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

func TestPutObjectTagging(t *testing.T) {
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

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	result, err := client.PutObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	versionId := *putObjResult.VersionId
	request = &PutObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
				},
			},
		},
	}
	result, err = client.PutObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.VersionId, versionId)

	_, err = client.PutObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	result, err = client.PutObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectTagging(t *testing.T) {
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

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &GetObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 0)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	_, err = client.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)

	request = &GetObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 2)

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionId := *putObjResult.VersionId

	request = &GetObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Len(t, result.Tags, 0)

	_, err = client.GetObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteObjectTagging(t *testing.T) {
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

	body := randLowStr(100)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
	}
	_, err = client.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)

	request := &DeleteObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.DeleteObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	versionRequest := &PutBucketVersioningRequest{
		Bucket: Ptr(bucketName),
		VersioningConfiguration: &VersioningConfiguration{
			Status: VersionEnabled,
		},
	}
	_, err = client.PutBucketVersioning(context.TODO(), versionRequest)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Body:   strings.NewReader(body),
	}
	putObjResult, err := client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	versionId := *putObjResult.VersionId

	request = &DeleteObjectTaggingRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr(versionId),
	}
	result, err = client.DeleteObjectTagging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	_, err = client.DeleteObjectTagging(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &DeleteObjectTaggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.DeleteObjectTagging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
