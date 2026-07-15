//go:build integrationignore

package agentic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestAgenticBucketLifecycle(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	// 1. Create agentic bucket
	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
		CreateAgenticBucketConfiguration: &CreateAgenticBucketConfiguration{
			StorageClass:       oss.StorageClassStandard,
			DataRedundancyType: oss.DataRedundancyLRS,
		},
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// 2. Get agentic bucket
	getResult, err := client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.AgenticBucketInfo)
	assert.Contains(t, oss.ToString(getResult.AgenticBucketInfo.Name), bucket)

	// 3. List agentic buckets, verify the created bucket appears
	found := false
	paginator := client.NewListAgenticBucketsPaginator(&ListAgenticBucketsRequest{})
	for paginator.HasNext() {
		page, err := paginator.NextPage(context.TODO())
		assert.Nil(t, err)
		for _, b := range page.AgenticBuckets {
			if b.Name != nil && strings.Contains(*b.Name, bucket) {
				found = true
			}
		}
	}
	assert.True(t, found, "created agentic bucket should appear in list")

	// 4. Delete agentic bucket
	deleteResult, err := client.DeleteAgenticBucket(context.TODO(), &DeleteAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestAgenticBucketStatus(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	putResult, err := client.PutAgenticBucketStatus(context.TODO(), &PutAgenticBucketStatusRequest{
		Bucket: oss.Ptr(bucket),
		AgenticBucketStatus: &AgenticBucketStatus{
			Status: oss.Ptr("Enabled"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
}

func TestAgenticBucketAcl(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// Put ACL
	putResult, err := client.PutAgenticBucketAcl(context.TODO(), &PutAgenticBucketAclRequest{
		Bucket: oss.Ptr(bucket),
		Acl:    oss.BucketACLPrivate,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)

	// Get ACL
	getResult, err := client.GetAgenticBucketAcl(context.TODO(), &GetAgenticBucketAclRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Equal(t, string(oss.BucketACLPrivate), oss.ToString(getResult.ACL))
}

func TestAgenticBucketEncryption(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// Put encryption
	putResult, err := client.PutAgenticBucketEncryption(context.TODO(), &PutAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket),
		ServerSideEncryptionRule: &ServerSideEncryptionRule{
			ApplyServerSideEncryptionByDefault: &ApplyServerSideEncryptionByDefault{
				SSEAlgorithm: oss.Ptr("AES256"),
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)

	// Get encryption
	getResult, err := client.GetAgenticBucketEncryption(context.TODO(), &GetAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.ServerSideEncryptionRule)
	assert.NotNil(t, getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault)
	assert.Equal(t, "AES256", oss.ToString(getResult.ServerSideEncryptionRule.ApplyServerSideEncryptionByDefault.SSEAlgorithm))

	// Delete encryption
	deleteResult, err := client.DeleteAgenticBucketEncryption(context.TODO(), &DeleteAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestAgenticBucketVersioning(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// Put versioning
	putResult, err := client.PutAgenticBucketVersioning(context.TODO(), &PutAgenticBucketVersioningRequest{
		Bucket: oss.Ptr(bucket),
		VersioningConfiguration: &VersioningConfiguration{
			Status: oss.VersionEnabled,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)

	// Get versioning
	getResult, err := client.GetAgenticBucketVersioning(context.TODO(), &GetAgenticBucketVersioningRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.VersioningConfiguration)
	assert.Equal(t, oss.VersionEnabled, getResult.VersioningConfiguration.Status)
}

func TestAgenticBucketPolicy(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	policy := fmt.Sprintf(`{"Version":"1","Statement":[{"Effect":"Allow","Action":["oss:GetObject"],"Principal":["*"],"Resource":["acs:oss:*:%s:*"]}]}`, accountId_)

	// Put policy
	putResult, err := client.PutAgenticBucketPolicy(context.TODO(), &PutAgenticBucketPolicyRequest{
		Bucket: oss.Ptr(bucket),
		Body:   strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)

	// Get policy
	getResult, err := client.GetAgenticBucketPolicy(context.TODO(), &GetAgenticBucketPolicyRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Contains(t, getResult.Body, "oss:GetObject")

	// Delete policy
	deleteResult, err := client.DeleteAgenticBucketPolicy(context.TODO(), &DeleteAgenticBucketPolicyRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestAgenticBucketPublicAccessBlock(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// Put public access block
	putResult, err := client.PutAgenticBucketPublicAccessBlock(context.TODO(), &PutAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr(bucket),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: oss.Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)

	// Get public access block
	getResult, err := client.GetAgenticBucketPublicAccessBlock(context.TODO(), &GetAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotNil(t, getResult.PublicAccessBlockConfiguration)

	// Delete public access block
	deleteResult, err := client.DeleteAgenticBucketPublicAccessBlock(context.TODO(), &DeleteAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.True(t, deleteResult.StatusCode == 200 || deleteResult.StatusCode == 204)
}

func TestListBucketSpaces(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	listResult, err := client.ListBucketSpaces(context.TODO(), &ListBucketSpacesRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
}

func TestBucketSpaceObjectLifecycle(t *testing.T) {
	skipIfNotConfigured(t)
	client := getDefaultClient()
	bsClient := getBucketSpaceClient()
	bucket := genBucketName()

	createResult, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)

	defer cleanAgenticBucket(bucket)

	// Create a bucket space using the short name.
	putBucketResult, err := bsClient.PutBucket(context.TODO(), &oss.PutBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, putBucketResult.StatusCode)

	defer func() {
		_, _ = bsClient.DeleteBucket(context.TODO(), &oss.DeleteBucketRequest{
			Bucket: oss.Ptr(bucket),
		})
	}()

	// Put an object into the bucket space.
	key := "go-sdk-test-object-" + randStr(6)
	putObjectResult, err := bsClient.PutObject(context.TODO(), &oss.PutObjectRequest{
		Bucket: oss.Ptr(bucket),
		Key:    oss.Ptr(key),
		Body:   strings.NewReader("hello world"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putObjectResult.StatusCode)

	defer func() {
		_, _ = bsClient.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
			Bucket: oss.Ptr(bucket),
			Key:    oss.Ptr(key),
		})
	}()

	// Read the object back.
	getObjectResult, err := bsClient.GetObject(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucket),
		Key:    oss.Ptr(key),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getObjectResult.StatusCode)
	getObjectResult.Body.Close()
}

func TestAgenticBucketServerErrors(t *testing.T) {
	skipIfNotConfigured(t)
	client := getInvalidAkClient()
	bucket := genBucketName()

	var serr *oss.ServiceError

	// Create with invalid AK
	_, err := client.CreateAgenticBucket(context.TODO(), &CreateAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 403, serr.StatusCode)
	assert.NotEmpty(t, serr.RequestID)

	// Get with invalid AK
	serr = nil
	_, err = client.GetAgenticBucket(context.TODO(), &GetAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 403, serr.StatusCode)

	// List with invalid AK
	serr = nil
	_, err = client.ListAgenticBuckets(context.TODO(), &ListAgenticBucketsRequest{})
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &serr))
	assert.Equal(t, 403, serr.StatusCode)
}
