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
		WithRegion(region).WithEndpoint("http://oss-cn-hangzhou.aliyuncs.com")

	client := dataprocess.NewClient(cfg)

	delResult, err := client.UpdateDataset(context.TODO(), &dataprocess.UpdateDatasetRequest{
		Bucket:      oss.Ptr(bucketName),
		DatasetName: oss.Ptr("test_dataset"),
		Description: oss.Ptr("this is a test"),
		WorkflowParameters: oss.Ptr(`
		[
		  {
		   "Name": "ImageInsightEnable",
			"Value": "True",
			"Description": "The source bucket for data processing"
		  }
		]`),
		DatasetConfig: oss.Ptr(`
		{
		  "Insights": {
			"Language": "zh"
		  }
		}`),
	})
	if err != nil {
		log.Fatalf("failed to update dataset %v", err)
	}

	log.Printf("update dataset result:%#v\n", delResult)
}
