//go:build integrationignore

package agentic

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region_    = os.Getenv("OSS_TEST_REGION")
	endpoint_  = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_  = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_ = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	accountId_ = os.Getenv("OSS_TEST_ACCOUNT_ID")

	instance_ *AgenticBucketClient
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = "go-sdk-test-ab-"
	letters          = []rune("abcdefghijklmnopqrstuvwxyz")
)

func getDefaultClient() *AgenticBucketClient {
	testOnce_.Do(func() {
		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_).
			WithAccountId(accountId_)

		instance_ = NewAgenticBucketClient(cfg)
	})
	return instance_
}

func getInvalidAkClient() *AgenticBucketClient {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithAccountId(accountId_)

	return NewAgenticBucketClient(cfg)
}

func getBucketSpaceClient() *oss.Client {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithAccountId(accountId_)

	return NewBucketSpaceClient(cfg)
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func genBucketName() string {
	return bucketNamePrefix + randStr(6)
}

// cleanAgenticBucket best-effort deletes an agentic bucket and its properties.
func cleanAgenticBucket(bucket string) {
	c := getDefaultClient()
	_, _ = c.DeleteAgenticBucketPolicy(context.TODO(), &DeleteAgenticBucketPolicyRequest{
		Bucket: oss.Ptr(bucket),
	})
	_, _ = c.DeleteAgenticBucketEncryption(context.TODO(), &DeleteAgenticBucketEncryptionRequest{
		Bucket: oss.Ptr(bucket),
	})
	_, _ = c.DeleteAgenticBucketPublicAccessBlock(context.TODO(), &DeleteAgenticBucketPublicAccessBlockRequest{
		Bucket: oss.Ptr(bucket),
	})
	_, _ = c.DeleteAgenticBucket(context.TODO(), &DeleteAgenticBucketRequest{
		Bucket: oss.Ptr(bucket),
	})
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

func skipIfNotConfigured(t *testing.T) {
	if accountId_ == "" || region_ == "" {
		t.Skip("agentic integration test requires OSS_TEST_ACCOUNT_ID and OSS_TEST_REGION")
	}
}
