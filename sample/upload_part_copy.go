package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

var (
	region        string
	bucketName    string
	objectName    string
	srcBucketName string
	srcObjectName string
	uploadId      string
)

func init() {
	flag.StringVar(&region, "region", "", "The region in which the bucket is located.")
	flag.StringVar(&bucketName, "bucket", "", "The name of the bucket.")
	flag.StringVar(&objectName, "object", "", "The name of the object.")
	flag.StringVar(&srcBucketName, "src-bucket", "", "The name of the source bucket.")
	flag.StringVar(&srcObjectName, "src-object", "", "The name of the source object.")
	flag.StringVar(&uploadId, "upload-id", "", "The upload id.")
}

func main() {
	flag.Parse()
	if len(uploadId) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, upload id required")
	}

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

	if len(srcObjectName) == 0 {
		flag.PrintDefaults()
		log.Fatalf("invalid parameters, the source object name required")
	}

	cfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewEnvironmentVariableCredentialsProvider()).
		WithRegion(region)

	client := oss.NewClient(cfg)

	headResult, err := client.HeadObject(context.TODO(), &oss.HeadObjectRequest{
		Bucket: oss.Ptr(srcBucketName),
		Key:    oss.Ptr(srcObjectName),
	})
	if err != nil {
		log.Fatalf("failed to head object %v", err)
	}
	total := headResult.ContentLength
	partSize := 64 * 1024 * 1024

	var wg sync.WaitGroup
	var parts []oss.UploadPart
	concurrency := 3
	var mu sync.Mutex

	totalParts := (int(total) + partSize - 1) / int(partSize)

	partCh := make(chan int, totalParts)
	for pn := 1; pn <= totalParts; pn++ {
		partCh <- pn
	}
	close(partCh)
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pn := range partCh {
				partRequest := &oss.UploadPartCopyRequest{
					Bucket:       oss.Ptr(bucketName),
					Key:          oss.Ptr(objectName),
					SourceBucket: oss.Ptr(srcBucketName),
					SourceKey:    oss.Ptr(srcObjectName),
					PartNumber:   int32(pn),
					UploadId:     oss.Ptr(uploadId),
					Range:        oss.Ptr(GetPartRange(int(total), partSize, pn)),
				}

				partResult, err := client.UploadPartCopy(context.TODO(), partRequest)
				if err != nil {
					log.Printf("failed to upload part copy %d: %v", pn, err)
					continue
				}

				mu.Lock()
				parts = append(parts, oss.UploadPart{
					PartNumber: partRequest.PartNumber,
					ETag:       partResult.ETag,
				})
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	log.Println("upload part copy success!")
}

func GetPartRange(totalSize, partSize, partNumber int) string {
	start := (partNumber - 1) * partSize
	end := partNumber*partSize - 1
	if end > totalSize-1 {
		end = totalSize - 1
	}
	return fmt.Sprintf("bytes=%d-%d", start, end)
}
