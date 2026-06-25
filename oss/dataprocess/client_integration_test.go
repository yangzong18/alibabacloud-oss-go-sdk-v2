//go:build integration

package dataprocess

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
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
	bucket_    = os.Getenv("OSS_TEST_DATAPROCESS_BUCKET")

	instance_ *Client
	testOnce_ sync.Once
)

var (
	datasetNamePrefix = "go-sdk-test-ds-"
	letters           = []rune("abcdefghijklmnopqrstuvwxyz")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := oss.LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_)

		instance_ = NewClient(cfg)
	})
	return instance_
}

func getInvalidAkClient() *Client {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider("invalid-ak", "invalid-sk")).
		WithRegion(region_).
		WithEndpoint(endpoint_)

	return NewClient(cfg)
}

func getClient(region, endpoint string) *Client {
	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(oss.SignatureVersionV4)

	return NewClient(cfg)
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func genDatasetName() string {
	return datasetNamePrefix + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-" + randStr(5)
}

func cleanDatasets(prefix string, t *testing.T) {
	c := getDefaultClient()
	request := &ListDatasetsRequest{
		Bucket: oss.Ptr(bucket_),
		Prefix: oss.Ptr(prefix),
	}
	result, err := c.ListDatasets(context.TODO(), request)
	if err != nil {
		return
	}
	for _, ds := range result.Datasets {
		if ds.DatasetName != nil && strings.HasPrefix(*ds.DatasetName, prefix) {
			_, _ = c.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
				Bucket:      oss.Ptr(bucket_),
				DatasetName: ds.DatasetName,
			})
		}
	}
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}

