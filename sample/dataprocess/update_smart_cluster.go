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
	objectId    string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&datasetName, "dataset", "", "The name of the dataset.")
	flag.StringVar(&objectId, "objectId", "", "The id of the object.")
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

	if len(objectId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	result, err := client.UpdateSmartCluster(context.TODO(), &dataprocess.UpdateSmartClusterRequest{
		Bucket:      oss.Ptr(bucketName),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    oss.Ptr(objectId),
		Description: oss.Ptr("this is a demo"),
		Name:        oss.Ptr("face-cluster-alice"),
		Rules:       oss.Ptr("[{\"RuleType\":\"face\",\"Sensitivity\":0.7}]"),
	})
	if err != nil {
		log.Fatalf("failed to update smart cluster %v", err)
	}
	log.Printf("update smart cluster result:%#v\n", result)
}
