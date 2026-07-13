package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region     string
	bucketName string
	objectName string
	uploadId   string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.StringVar(&uploadId, "upload-id", "", "The upload id.")
}

func main() {
	flag.Parse()
	var (
		parts = oss.UploadParts{
			oss.UploadPart{
				PartNumber: int32(1),
				ETag:       oss.Ptr("part one etag"),
			},
			oss.UploadPart{
				PartNumber: int32(2),
				ETag:       oss.Ptr("part two etag"),
			},
			oss.UploadPart{
				PartNumber: int32(3),
				ETag:       oss.Ptr("part three etag"),
			},
		}
	)
	if len(bucketName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, bucket name required")
	}

	if len(region) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, region required")
	}

	if len(objectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, object name required")
	}

	if len(uploadId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, upload id required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)
	request := &oss.CompleteMultipartUploadRequest{
		Bucket:   oss.Ptr(bucketName),
		Key:      oss.Ptr(objectName),
		UploadId: oss.Ptr(uploadId),
		CompleteMultipartUpload: &oss.CompleteMultipartUpload{
			Parts: parts,
		},
	}
	result, err := client.CompleteMultipartUpload(context.TODO(), request)
	if err != nil {
		log.Fatalf("failed to complete multipart upload %v", err)
	}
	log.Printf("complete multipart upload result:%#v\n", result)
}
