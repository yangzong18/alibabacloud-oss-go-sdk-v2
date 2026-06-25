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
	region string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
}

func main() {
	flag.Parse()

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := dataprocess.NewClient(cfg)

	p := client.NewListDataPipelineConfigurationsPaginator(&dataprocess.ListDataPipelineConfigurationsRequest{})

	var i int
	log.Println("Data Pipeline Configurations:")
	for p.HasNext() {
		i++

		page, err := p.NextPage(context.TODO())
		if err != nil {
			log.Fatalf("failed to get page %v, %v", i, err)
		}

		for _, data := range page.DataPipelineConfigurations {
			log.Printf("Data Pipeline Name:%v, Data Pipeline Description:%v, Data Pipeline Role:%v, Status:%v, Phase:%v\n", oss.ToString(data.DataPipelineName), oss.ToString(data.DataPipelineDescription), oss.ToString(data.DataPipelineRole), oss.ToString(data.Status), oss.ToString(data.Phase))
		}
	}
}
