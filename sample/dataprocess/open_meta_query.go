package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/dataprocess"
)

var (
	region     string
	bucketName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
}

func main() {
	flag.Parse()
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	request := &dataprocess.OpenMetaQueryRequest{
		Bucket: oss.Ptr(bucketName),
		Mode:   oss.Ptr("basic"),
		MetaQuery: &dataprocess.OpenMetaQuery{
			Filters: &dataprocess.Filters{
				Filter: []string{"Size > 1024"},
			},
		},
	}
	result, err := client.OpenMetaQuery(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to open meta query %v", err)
	}
	log.Printf("open meta query result:%#v\n", result)
}
