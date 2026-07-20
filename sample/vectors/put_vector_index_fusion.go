package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/vectors"
)

var (
	region     string
	bucketName string
	accountId  string
	indexName  string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the vector bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the vector bucket.")
	flag.StringVar(&accountId, "account-id", "", "The id of vector account.")
	flag.StringVar(&indexName, "index", "", "The name of vector index.")
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

	if len(accountId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, accounId required")
	}

	if len(indexName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, index required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region).WithAccountId(accountId)

	client := vectors.NewVectorsClient(cfg)

	request := &vectors.PutVectorIndexFusionRequest{
		Bucket:    oss.Ptr(bucketName),
		IndexName: oss.Ptr(indexName),
		Mode:      oss.Ptr("fusion"),
		SchemaConfiguration: map[string]any{
			"fields": []any{
				map[string]any{
					"name":           "vector_1",
					"type":           "vector",
					"dataType":       "float32",
					"dimension":      1024,
					"distanceMetric": "euclidean",
				},
				map[string]any{
					"name":           "vector_2",
					"type":           "vector",
					"dataType":       "float32",
					"dimension":      512,
					"distanceMetric": "cosine",
				},
				map[string]any{
					"name":    "timestamps",
					"type":    "long",
					"isArray": true,
				},
				map[string]any{
					"name": "price",
					"type": "double",
				},
				map[string]any{
					"name": "ip",
					"type": "ip",
				},
				map[string]any{
					"name": "location",
					"type": "geoPoint",
				},
				map[string]any{
					"name": "tag",
					"type": "string",
				},
				map[string]any{
					"name":           "user_id",
					"type":           "string",
					"isPartitionKey": true,
				},
				map[string]any{
					"name":    "tags",
					"type":    "string",
					"isArray": true,
				},
				map[string]any{
					"name":       "title_1",
					"type":       "string",
					"exactMatch": true,
					"text": map[string]any{
						"enabled":  true,
						"analyzer": "standard",
						"analyzerParameters": map[string]any{
							"caseSensitive": true,
							"delimitWord":   false,
						},
					},
				},
				map[string]any{
					"name":       "title_2",
					"type":       "string",
					"exactMatch": false,
					"text": map[string]any{
						"enabled":  true,
						"analyzer": "split",
						"analyzerParameters": map[string]any{
							"caseSensitive": true,
							"delimiter":     "",
						},
					},
				},
			},
		},
	}
	result, err := client.PutVectorIndexFusion(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to put vector index fusion%v", err)
	}
	log.Printf("put vector index fusion result:%#v\n", result)
}
