package main

import (
	"context"
	"flag"
	"fmt"
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
	flag.StringVar(&bucketName, "bucket", "", "The `name` of the bucket.")
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

	request := &dataprocess.ListSmartClustersRequest{
		Bucket:      oss.Ptr(bucketName),
		DatasetName: oss.Ptr(datasetName),
	}
	p := client.NewListSmartClustersPaginator(request)

	var i int
	log.Println("Smart Clusters:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		for _, SmartCluster := range page.SmartClusters {
			log.Printf("Object Id:%v, Cluster Type:%v, Name:%v, Description:%v, Create Time:%v, Update Time:%v, Reason:%v\n", oss.ToString(SmartCluster.ObjectId), oss.ToString(SmartCluster.ClusterType), oss.ToString(SmartCluster.Name), oss.ToString(SmartCluster.Description), oss.ToString(SmartCluster.CreateTime), oss.ToString(SmartCluster.UpdateTime), oss.ToString(SmartCluster.Reason))
			for _, rule := range SmartCluster.Rules {
				var sensitivity float64
				if rule.Sensitivity != nil {
					sensitivity = *rule.Sensitivity
				}
				fmt.Printf("Rule Type:%v, Base URIs:%v, Keywords:%v, Sensitivity:%v\n",
					oss.ToString(rule.RuleType),
					rule.BaseURIs,
					rule.Keywords,
					sensitivity,
				)
			}
		}
	}
}
