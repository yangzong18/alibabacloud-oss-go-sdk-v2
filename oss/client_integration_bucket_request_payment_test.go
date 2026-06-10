//go:build integration

package oss

import (
	"bufio"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestPutBucketRequestPayment(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}

	result, err := client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketNameNotExist),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	result, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketRequestPayment(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	request := &GetBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *result.Payer, "BucketOwner")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &GetBucketRequestPaymentRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	result, err = client.GetBucketRequestPayment(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func testPaymentWithRequester(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	_, err = client.PutBucketPublicAccessBlock(context.TODO(), &PutBucketPublicAccessBlockRequest{
		Bucket: Ptr(bucketName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			BlockPublicAccess: Ptr(false),
		},
	})
	policyInfo := `
	{
		"Version":"1",
		"Statement":[
			{
				"Action":[
					"oss:*"
				],
				"Effect":"Allow",
				"Principal":["` + payerUID_ + `"],
				"Resource":["acs:oss:*:*:` + bucketName + `", "acs:oss:*:*:` + bucketName + `/*"]
			}
		]
	}`
	input := &OperationInput{
		OpName: "PutBucketPolicy",
		Bucket: Ptr(bucketName),
		Method: "PUT",
		Parameters: map[string]string{
			"policy": "",
		},
		Body: strings.NewReader(policyInfo),
	}
	input.OpMetadata.Set(signer.SubResource, []string{"policy"})
	_, err = client.InvokeOperation(context.TODO(), input)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	request := &PutBucketRequestPaymentRequest{
		Bucket: Ptr(bucketName),
		PaymentConfiguration: &RequestPaymentConfiguration{
			Payer: Requester,
		},
	}
	_, err = client.PutBucketRequestPayment(context.TODO(), request)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	body := randStr(100)
	creClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider(payerAccessID_, payerAccessKey_))

	objectName := objectNamePrefix + randStr(6)

	putObjReq := &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getObjReq := &GetObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	getObjResult, err := creClient.GetObject(context.TODO(), getObjReq)
	assert.Nil(t, err)
	getObjData, _ := io.ReadAll(getObjResult.Body)
	assert.Equal(t, string(getObjData), body)
	time.Sleep(1 * time.Second)

	objectCopyName := objectName + "-copy"
	copyRequest := &CopyObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		SourceKey:    Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	objectAppendName := objectName + "-append"
	appendRequest := &AppendObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectAppendName),
		Body:         strings.NewReader(body),
		Position:     Ptr(int64(0)),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.AppendObject(context.TODO(), appendRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delObjsRequest := &DeleteMultipleObjectsRequest{
		Bucket:       Ptr(bucketName),
		Objects:      []DeleteObject{{Key: Ptr(objectAppendName)}},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteMultipleObjects(context.TODO(), delObjsRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	headRequest := &HeadObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.HeadObject(context.TODO(), headRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	metaRequest := &GetObjectMetaRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectCopyName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectMeta(context.TODO(), metaRequest)
	assert.Nil(t, err)

	objectRestoreName := objectName + "-restore"
	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectRestoreName),
		Body:         strings.NewReader(body),
		StorageClass: StorageClassColdArchive,
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	restoreRequest := &RestoreObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectRestoreName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putAclRequest := &PutObjectAclRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Acl:          ObjectACLPrivate,
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObjectAcl(context.TODO(), putAclRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getAclRequest := &GetObjectAclRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectAcl(context.TODO(), getAclRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	objectMultiName := objectName + "-multi"
	body = randLowStr(360000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := io.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]
	initRequest := &InitiateMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		RequestPayer: Ptr("requester"),
	}
	initResult, err := creClient.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	contents := []string{string(part1), string(part2), string(part3)}
	var parts []UploadPart
	var wg sync.WaitGroup
	wg.Add(len(contents))
	for i, content1 := range contents {
		partRequest := &UploadPartRequest{
			Bucket:       Ptr(bucketName),
			Key:          Ptr(objectMultiName),
			PartNumber:   int32(i + 1),
			UploadId:     Ptr(*initResult.UploadId),
			Body:         strings.NewReader(content1),
			RequestPayer: Ptr("requester"),
		}
		partResult, err := creClient.UploadPart(context.TODO(), partRequest)
		assert.Nil(t, err)

		part := UploadPart{
			PartNumber: partRequest.PartNumber,
			ETag:       partResult.ETag,
		}
		parts = append(parts, part)
		wg.Done()
	}

	comRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectMultiName),
		UploadId: Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Parts: parts,
		},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.CompleteMultipartUpload(context.TODO(), comRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	initRequest = &InitiateMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		RequestPayer: Ptr("requester"),
	}
	initResult, err = creClient.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	copyMultiRequest := &UploadPartCopyRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		PartNumber:   int32(1),
		UploadId:     Ptr(*initResult.UploadId),
		SourceKey:    Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.UploadPartCopy(context.TODO(), copyMultiRequest)
	assert.Nil(t, err)

	listMultiRequest := &ListMultipartUploadsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListMultipartUploads(context.TODO(), listMultiRequest)

	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listRequest := &ListPartsRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		UploadId:     Ptr(*initResult.UploadId),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListParts(context.TODO(), listRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectMultiName),
		UploadId:     Ptr(*initResult.UploadId),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	symlinkName := objectName + "-symlink"
	putSymRequest := &PutSymlinkRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(symlinkName),
		Target:       Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutSymlink(context.TODO(), putSymRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getSymRequest := &GetSymlinkRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(symlinkName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetSymlink(context.TODO(), getSymRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	putTagRequest := &PutObjectTaggingRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Tagging: &Tagging{
			&TagSet{
				Tags: []Tag{
					{
						Key:   Ptr("k1"),
						Value: Ptr("v1"),
					},
					{
						Key:   Ptr("k2"),
						Value: Ptr("v2"),
					},
				},
			},
		},
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.PutObjectTagging(context.TODO(), putTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	getTagRequest := &GetObjectTaggingRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.GetObjectTagging(context.TODO(), getTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	delTagRequest := &DeleteObjectTaggingRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.DeleteObjectTagging(context.TODO(), delTagRequest)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjReq := &ListObjectsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjects(context.TODO(), listObjReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjReqV2 := &ListObjectsV2Request{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjectsV2(context.TODO(), listObjReqV2)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	listObjVersionReq := &ListObjectVersionsRequest{
		Bucket:       Ptr(bucketName),
		RequestPayer: Ptr("requester"),
	}
	_, err = creClient.ListObjectVersions(context.TODO(), listObjVersionReq)
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	putObjReq = &PutObjectRequest{
		Bucket:       Ptr(bucketName),
		Key:          Ptr(objectName),
		Body:         strings.NewReader(body),
		RequestPayer: Ptr("bucketOwner"),
	}
	_, err = creClient.PutObject(context.TODO(), putObjReq)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Access denied for requester pay bucket")
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "AccessDenied", serr.Code)
	assert.Equal(t, "Access denied for requester pay bucket", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
