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
	apiKey           string
	dataPipelineName string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&apiKey, "api-key", "", "the Bailian API key.")
	flag.StringVar(&dataPipelineName, "data-pipeline-name", "", "the name of the data pipeline.")
}

func main() {
	flag.Parse()
	if len(apiKey) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, api key required")
	}

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

	result, err := client.PutDataPipelineConfiguration(context.TODO(), &dataprocess.PutDataPipelineConfigurationRequest{
		DataPipelineName: oss.Ptr(dataPipelineName),
		Role:             oss.Ptr("AliyunOSSDataPipelineRole"),
		DataPipelineConfiguration: &dataprocess.DataPipelineConfiguration{
			DataPipelineDescription: oss.Ptr("Vectorize business data using the BERT multimodal model"),
			Sources: []dataprocess.DataPipelineSource{
				{
					InputBucket:    oss.Ptr("bucket"),
					InputDataScope: oss.Ptr("All"),
					FilterConfiguration: &dataprocess.DataPipelineSourceFilterConfiguration{
						PrefixSet:        []string{"prefix1"},
						ObjectMediaTypes: []string{"text"},
					},
				},
			},
			DataPipelineEmbeddingConfiguration: &dataprocess.DataPipelineEmbeddingConfiguration{
				ApiKey:            oss.Ptr(apiKey),
				EmbeddingProvider: oss.Ptr("bailian"),
				FPS:               oss.Ptr(float64(1)),
				Model:             oss.Ptr("qwen2.5-vl-embedding"),
			},
			Destination: &dataprocess.DataPipelineDestination{
				VectorBucketName:    oss.Ptr("my-vector-bucket"),
				VectorIndexNames:    []string{"index"},
				VectorKeyPrefix:     oss.Ptr("prefix"),
				ObjectTagToMetadata: []string{"key1"},
				UsermetaToMetadata:  []string{"x-oss-meta-key1"},
			},
			DataPipelineError: &dataprocess.DataPipelineError{
				ErrorMode:   oss.Ptr("ignoreAndRecord"),
				ErrorBucket: oss.Ptr("my-error-bucket"),
				ErrorPrefix: oss.Ptr("error-output/"),
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to put pipeline configuration %v", err)
	}
	log.Printf("put pipeline configuration result:%#v\n", result)
}
