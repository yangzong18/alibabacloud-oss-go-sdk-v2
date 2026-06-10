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
	region      string
	bucketName  string
	datasetName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&datasetName, "dataset", "", "The name of the dataset.")
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

	if len(datasetName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, dataset name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	request := &dataprocess.CreateSmartClusterRequest{
		Bucket:      oss.Ptr(bucketName),
		DatasetName: oss.Ptr(datasetName),
		Name:        oss.Ptr("demo"),
		ClusterType: dataprocess.SmartClusterTypeKnowledge,
		Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["character","car"]}]`),
	}
	result, err := client.CreateSmartCluster(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to create smart cluster %v", err)
	}
	log.Printf("create smart cluster result:%#v\n", result)
}
