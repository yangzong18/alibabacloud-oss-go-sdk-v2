package oss

import (
	"context"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
)

var (
	// Endpoint/ID/Key
	region_    = os.Getenv("OSS_TEST_REGION")
	endpoint_  = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_  = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_ = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")

	instance_ *Client
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = "go-sdk-test-bucket-"
	objectNamePrefix = "go-sdk-test-object-"
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_)
		instance_ = NewClient(cfg)
	})
	return instance_
}

func getClient(region, endpoint string) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint)
	return NewClient(cfg)
}

func getKmsID() string {
	return ""
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func randLowStr(n int) string {
	return strings.ToLower(randStr(n))
}

func cleanBucket(bucketinfo BucketProperties) {
	if bucketinfo.Name == "" {
		return
	}

	var c *Client
	if strings.Contains(endpoint_, bucketinfo.ExtranetEndpoint) ||
		strings.Contains(endpoint_, bucketinfo.IntranetEndpoint) {
		c = getDefaultClient()
	} else {
		c = getClient(*bucketinfo.Region, bucketinfo.ExtranetEndpoint)
	}

	if c == nil {
		return
	}
}

func cleanBuckets(prefix string) {
	c := getDefaultClient()
	for {
		request := &ListBucketsRequest{}
		result, err := c.ListBuckets(context.TODO(), request)
		if err != nil {
			return
		}

		for _, b := range result.Buckets {
			cleanBucket(b)
		}
	}
}

func before(t *testing.T) func(t *testing.T) {
	//fmt.Println("setup test case")

	return after
}

func after(t *testing.T) {
	//fmt.Println("teardown  test case")
}

func TestListBuckets(t *testing.T) {
	after := before(t)
	defer after(t)

	//bucketPrefix := bucketNamePrefix + randLowStr(6)
	//TODO
}
