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
	region           string
	dataPipelineName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&dataPipelineName, "data-pipeline-name", "", "the name of the data pipeline.")
}

func main() {
	flag.Parse()
	if len(dataPipelineName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, data pipeline name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	result, err := client.DeleteDataPipelineConfiguration(context.TODO(), &dataprocess.DeleteDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr(dataPipelineName),
	})

	if err != nil {
		log.Fatalf("failed to delete pipeline configuration %v", err)
	}
	log.Printf("delete pipeline configuration result:%#v\n", result)
}
