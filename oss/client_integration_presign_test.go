//go:build integration

package oss

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestPresign(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	// PutObjRequest
	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	// GetObjRequest
	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	// HeadObjRequest
	headObjRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), headObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(HTTPHeaderContentLength), strconv.Itoa(len(body)))

	// MultiPart
	objectNameMultipart := objectNamePrefix + randLowStr(6) + "-multi-part"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initResult)
	assert.Equal(t, *initResult.Key, objectNameMultipart)
	uploadId := initResult.UploadId

	//UploadPart
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectNameMultipart),
		PartNumber: int32(1),
		UploadId:   uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), partRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	var parts []UploadPart
	uploadResult := &UploadPartResult{}
	err = xml.Unmarshal(data, uploadResult)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       Ptr(resp.Header.Get("ETag")),
	}
	parts = append(parts, part)
	completeRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipart),
		UploadId: uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), completeRequest, PresignExpiration(expiration))
	assert.Nil(t, err)

	//Complete
	upload := CompleteMultipartUpload{
		Parts: parts,
	}
	xmlData, err := xml.Marshal(upload)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(string(xmlData)))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	headObjRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	headResult, err := client.HeadObject(context.TODO(), headObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, headResult.Headers.Get(HTTPHeaderContentLength), strconv.FormatInt(int64(len(body)), 10))
	assert.Equal(t, *headResult.ObjectType, "Multipart")

	// Test Abort
	objectNameMultipartCopy := objectNamePrefix + randLowStr(6) + "-multi-part-copy"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipartCopy),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initCopyRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initCopyResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initCopyResult)
	assert.Equal(t, *initCopyResult.Key, objectNameMultipartCopy)
	copyUploadId := *initCopyResult.UploadId

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), abortRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 204)

	listPartsRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	_, err = client.ListParts(context.TODO(), listPartsRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	deleteBucket(bucketName, t)
}

func TestPresignExtra(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	// PutObjRequest
	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	cfgV1 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV1).WithAdditionalHeaders([]string{"content-length"})

	clientV1 := NewClient(cfgV1)
	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(1 * time.Second)
	result, err = clientV1.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(-10 * time.Second)
	result, err = clientV1.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)

	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 403)

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"content-length": "1000",
			},
		},
	}
	result, err = clientV1.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	cfgV4 := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV4).WithAdditionalHeaders([]string{"content-length"})

	clientV4 := NewClient(cfgV4)
	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(1 * time.Second)
	result, err = clientV4.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	assert.Equal(t, map[string]string(nil), result.SignedHeaders)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ = io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	getObjRequest = &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(8 * 24 * time.Hour)
	result, err = clientV4.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "expires should be not greater than 604800(seven days)")

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"content-length": "1000",
			},
		},
	}
	result, err = clientV4.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"Content-Length": "1000"}, result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader("hi oss"))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 403)

	cfgV4 = LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region_).
		WithEndpoint(endpoint_).
		WithSignatureVersion(SignatureVersionV4).WithAdditionalHeaders([]string{"email", "name"})

	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Headers: map[string]string{
				"email": "demo@aliyun.com",
				"name":  "aliyun",
			},
		},
	}
	clientV4 = NewClient(cfgV4)
	result, err = clientV4.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"Email": "demo@aliyun.com",
		"Name": "aliyun"}, result.SignedHeaders)
	assert.NotEmpty(t, result.Expiration)

	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 400)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)

	header := make(http.Header)
	for key, value := range result.SignedHeaders {
		header[key] = []string{value}
	}
	req.Header = header
	resp, err = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	deleteBucket(bucketName, t)
}

func TestPresignWithStsToken(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	client := getClientUseStsToken(region_, endpoint_)
	assert.NotNil(t, client)

	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := randLowStr(1000)
	objectName := objectNamePrefix + randLowStr(6)
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.Presign(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	req, err := http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	c := &http.Client{}
	resp, err := c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	getObjRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration := time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), getObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	assert.Equal(t, "GET", result.Method)
	assert.NotEmpty(t, result.Expiration)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	data, _ := io.ReadAll(resp.Body)
	assert.Equal(t, string(data), body)

	headObjRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), headObjRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, resp.Header.Get(HTTPHeaderContentLength), fmt.Sprint(len(body)))

	objectNameMultipart := objectNamePrefix + randLowStr(6) + "-multi-part"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initResult)
	assert.Equal(t, *initResult.Key, objectNameMultipart)
	uploadId := initResult.UploadId

	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectNameMultipart),
		PartNumber: int32(1),
		UploadId:   uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), partRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	parts := []UploadPart{}
	uploadResult := &UploadPartResult{}
	err = xml.Unmarshal(data, uploadResult)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       Ptr(resp.Header.Get("ETag")),
	}
	parts = append(parts, part)
	completeRequest := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipart),
		UploadId: uploadId,
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), completeRequest, PresignExpiration(expiration))
	assert.Nil(t, err)

	upload := CompleteMultipartUpload{
		Parts: parts,
	}
	xmlData, err := xml.Marshal(upload)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(string(xmlData)))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)

	headObjRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipart),
	}
	headResult, err := client.HeadObject(context.TODO(), headObjRequest)
	assert.Nil(t, err)
	assert.Equal(t, headResult.Headers.Get(HTTPHeaderContentLength), strconv.FormatInt(int64(len(body)), 10))
	assert.Equal(t, *headResult.ObjectType, "Multipart")

	objectNameMultipartCopy := objectNamePrefix + randLowStr(6) + "-multi-part-copy"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameMultipartCopy),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), initCopyRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, nil)
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 200)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	assert.Nil(t, err)
	initCopyResult := &InitiateMultipartUploadResult{}
	err = xml.Unmarshal(data, initCopyResult)
	assert.Equal(t, *initCopyResult.Key, objectNameMultipartCopy)
	copyUploadId := *initCopyResult.UploadId

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	expiration = time.Now().Add(100 * time.Second)
	result, err = client.Presign(context.TODO(), abortRequest, PresignExpiration(expiration))
	assert.Nil(t, err)
	req, err = http.NewRequest(result.Method, result.URL, strings.NewReader(body))
	assert.Nil(t, err)
	resp, _ = c.Do(req)
	assert.Equal(t, resp.StatusCode, 204)

	listPartsRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectNameMultipartCopy),
		UploadId: Ptr(copyUploadId),
	}
	_, err = client.ListParts(context.TODO(), listPartsRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	cleanObjects(client, bucketName, t)
}
