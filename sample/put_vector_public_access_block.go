package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region string
	uid    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&uid, "uid", "", "The id of vector account.")
}

func main() {
	flag.Parse()
	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(uid) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, uid required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithUid(uid)

	client := oss.NewVectorsClient(cfg)

	request := &oss.PutPublicAccessBlockRequest{
		PublicAccessBlockConfiguration: &oss.PublicAccessBlockConfiguration{
			oss.Ptr(true),
		},
	}
	putResult, err := client.PutPublicAccessBlock(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector public access block %v", err)
	}

	log.Printf("put vector public access block result:%#v\n", putResult)
}
