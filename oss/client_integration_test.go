package oss

import (
	"bufio"
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
)

var (
	// Endpoint/ID/Key
	region_    = os.Getenv("OSS_TEST_REGION")
	endpoint_  = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_  = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_ = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")

	instance_ *Client
	testOnce_ sync.Once
)

var (
	bucketNamePrefix = "go-sdk-test-bucket-"
	objectNamePrefix = "go-sdk-test-object-"
	letters          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_)
		instance_ = NewClient(cfg)
	})
	return instance_
}

func getClient(region, endpoint string) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint)
	return NewClient(cfg)
}

func getKmsID() string {
	return ""
}

func randStr(n int) string {
	b := make([]rune, n)
	randMarker := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[randMarker.Intn(len(letters))]
	}
	return string(b)
}

func randLowStr(n int) string {
	return strings.ToLower(randStr(n))
}

func cleanBucket(bucketInfo BucketProperties, t *testing.T) {
	var err error
	assert.NotEmpty(t, *bucketInfo.Name)
	var c *Client
	if strings.Contains(endpoint_, *bucketInfo.ExtranetEndpoint) ||
		strings.Contains(endpoint_, *bucketInfo.IntranetEndpoint) {
		c = getDefaultClient()
	} else {
		c = getClient(*bucketInfo.Region, *bucketInfo.ExtranetEndpoint)
	}
	assert.NotNil(t, c)
	var listRequest *ListObjectsRequest
	var delObjRequest *DeleteObjectRequest
	var lor *ListObjectsResult
	marker := ""
	for {
		listRequest = &ListObjectsRequest{
			Bucket: Ptr(*bucketInfo.Name),
			Marker: Ptr(marker),
		}
		lor, err = c.ListObjects(context.TODO(), listRequest)
		assert.Nil(t, err)
		for _, object := range lor.Contents {
			delObjRequest = &DeleteObjectRequest{
				Bucket: Ptr(*bucketInfo.Name),
				Key:    Ptr(*object.Key),
			}
			_, err = c.DeleteObject(context.TODO(), delObjRequest)
			assert.Nil(t, err)
		}
		if !lor.IsTruncated {
			break
		}
		if lor.NextMarker != nil {
			marker = *lor.NextMarker
		}
	}
	var listUploadRequest *ListMultipartUploadsRequest
	var abortRequest *AbortMultipartUploadRequest
	var lsRes *ListMultipartUploadsResult
	keyMarker := ""
	uploadIdMarker := ""
	for {
		listUploadRequest = &ListMultipartUploadsRequest{
			Bucket:         Ptr(*bucketInfo.Name),
			KeyMarker:      Ptr(keyMarker),
			UploadIdMarker: Ptr(uploadIdMarker),
		}
		lsRes, err = c.ListMultipartUploads(context.TODO(), listUploadRequest)
		assert.Nil(t, err)
		for _, upload := range lsRes.Uploads {
			abortRequest = &AbortMultipartUploadRequest{
				Bucket:   Ptr(*bucketInfo.Name),
				Key:      Ptr(*upload.Key),
				UploadId: Ptr(*upload.UploadId),
			}
			_, err = c.AbortMultipartUpload(context.TODO(), abortRequest)
			assert.Nil(t, err)
		}
		if !lsRes.IsTruncated {
			break
		}
		keyMarker = *lsRes.NextKeyMarker
		uploadIdMarker = *lsRes.NextUploadIdMarker
	}
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(*bucketInfo.Name),
	}
	_, err = c.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func cleanBuckets(prefix string, t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListBucketsRequest{
			Prefix: Ptr(prefix),
		}
		result, err := c.ListBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			cleanBucket(b, t)
		}
	}
}

func before(t *testing.T) func(t *testing.T) {

	//fmt.Println("setup test case")
	return after
}

func after(t *testing.T) {
	cleanBuckets(bucketNamePrefix, t)
	//fmt.Println("teardown  test case")
}

func TestListBuckets(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketPrefix := bucketNamePrefix + randLowStr(6)
	//TODO
	var bucketName string
	count := 10
	for i := 0; i < count; i++ {
		bucketName = bucketPrefix + strconv.Itoa(i)
		putRequest := &PutBucketRequest{
			Bucket: Ptr(bucketName),
		}

		client := getDefaultClient()
		_, err := client.PutBucket(context.TODO(), putRequest)
		assert.Nil(t, err)
	}

	listRequest := &ListBucketsRequest{
		Prefix: Ptr(bucketPrefix),
	}

	client := getDefaultClient()
	result, err := client.ListBuckets(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, len(result.Buckets), count)
}

func TestPutBucket(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	result, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "200 OK")
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	putRequest = &PutBucketRequest{
		Bucket: Ptr(bucketName),
		CreateBucketConfiguration: &CreateBucketConfiguration{
			StorageClass:       StorageClassStandard,
			DataRedundancyType: DataRedundancyLRS,
		},
	}
	result, err = client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "200 OK")
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)

	delRequest = &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestDeleteBucket(t *testing.T) {
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

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.Status, "204 No Content")
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id") != "", true)

	result, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.NotNil(t, serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, strings.Contains(serr.Message, "not exist"), true)
	assert.Equal(t, serr.RequestID != "", true)
}

func TestListObjects(t *testing.T) {
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
	request := &ListObjectsRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.ListObjects(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.Name, bucketName)
	assert.Equal(t, len(result.Contents), 0)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.Marker)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)

	bucketNotExist := bucketNamePrefix + "not-exist" + randLowStr(5)
	request = &ListObjectsRequest{
		Bucket: Ptr(bucketNotExist),
	}
	_, err = client.ListObjects(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestListObjectsV2(t *testing.T) {
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
	request := &ListObjectsRequestV2{
		Bucket: Ptr(bucketName),
	}
	result, err := client.ListObjectsV2(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.Name, bucketName)
	assert.Equal(t, len(result.Contents), 0)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.StartAfter)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)

	bucketNotExist := bucketNamePrefix + "not-exist" + randLowStr(5)
	request = &ListObjectsRequestV2{
		Bucket: Ptr(bucketNotExist),
	}
	_, err = client.ListObjectsV2(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func TestGetBucketInfo(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	getRequest := &GetBucketInfoRequest{
		Bucket: Ptr(bucketName),
	}
	info, err := client.GetBucketInfo(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, *info.BucketInfo.Name, bucketName)
	assert.Equal(t, *info.BucketInfo.AccessMonitor, "Disabled")
	assert.NotEmpty(t, *info.BucketInfo.CreationDate)
	assert.True(t, strings.Contains(*info.BucketInfo.ExtranetEndpoint, ".aliyuncs.com"))
	assert.True(t, strings.Contains(*info.BucketInfo.IntranetEndpoint, "internal.aliyuncs.com"))
	assert.True(t, strings.Contains(*info.BucketInfo.Location, "oss-"))
	assert.True(t, strings.Contains(*info.BucketInfo.StorageClass, "Standard"))
	assert.Equal(t, *info.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *info.BucketInfo.CrossRegionReplication, "Disabled")
	assert.NotEmpty(t, *info.BucketInfo.ResourceGroupId)
	assert.NotEmpty(t, *info.BucketInfo.Owner.DisplayName)
	assert.NotEmpty(t, *info.BucketInfo.Owner.DisplayName)
	assert.Equal(t, *info.BucketInfo.ACL, "private")
	assert.Empty(t, info.BucketInfo.BucketPolicy.LogBucket)
	assert.Empty(t, info.BucketInfo.BucketPolicy.LogPrefix)

	assert.Equal(t, *info.BucketInfo.SseRule.SSEAlgorithm, "")
	assert.Nil(t, info.BucketInfo.SseRule.KMSDataEncryption)
	assert.Nil(t, info.BucketInfo.SseRule.KMSMasterKeyID)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	_, err = client.GetBucketInfo(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketLocation(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	getRequest := &GetBucketLocationRequest{
		Bucket: Ptr(bucketName),
	}
	info, err := client.GetBucketLocation(context.TODO(), getRequest)
	assert.Nil(t, err)

	endpoint := endpoint_
	if strings.HasPrefix(endpoint_, "http://") {
		endpoint = endpoint_[len("http://"):]
	} else if strings.HasPrefix(endpoint_, "https://") {
		endpoint = endpoint_[len("https://"):]
	}
	endpoint = strings.TrimSuffix(endpoint, ".aliyuncs.com")
	assert.Equal(t, *info.LocationConstraint, endpoint)
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	_, err = client.GetBucketLocation(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetBucketStat(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	getRequest := &GetBucketStatRequest{
		Bucket: Ptr(bucketName),
	}
	stat, err := client.GetBucketStat(context.TODO(), getRequest)
	assert.Nil(t, err)

	assert.Equal(t, int64(0), stat.Storage)
	assert.Equal(t, int64(0), stat.ObjectCount)
	assert.Equal(t, int64(0), stat.MultipartUploadCount)
	assert.Equal(t, int64(0), stat.LiveChannelCount)
	assert.Equal(t, int64(0), stat.LastModifiedTime)
	assert.Equal(t, int64(0), stat.StandardStorage)
	assert.Equal(t, int64(0), stat.StandardObjectCount)
	assert.Equal(t, int64(0), stat.InfrequentAccessStorage)
	assert.Equal(t, int64(0), stat.InfrequentAccessRealStorage)
	assert.Equal(t, int64(0), stat.InfrequentAccessObjectCount)
	assert.Equal(t, int64(0), stat.ArchiveStorage)
	assert.Equal(t, int64(0), stat.ArchiveRealStorage)
	assert.Equal(t, int64(0), stat.ArchiveObjectCount)
	assert.Equal(t, int64(0), stat.ColdArchiveStorage)
	assert.Equal(t, int64(0), stat.ColdArchiveRealStorage)
	assert.Equal(t, int64(0), stat.ColdArchiveObjectCount)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	_, err = client.GetBucketStat(context.TODO(), getRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutBucketAcl(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	request := &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
		Acl:    BucketACLPublicRead,
	}
	result, err := client.PutBucketAcl(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	infoRequest := &GetBucketInfoRequest{
		Bucket: Ptr(bucketName),
	}

	info, err := client.GetBucketInfo(context.TODO(), infoRequest)
	assert.Nil(t, err)
	assert.Equal(t, string(BucketACLPublicRead), *info.BucketInfo.ACL)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	request = &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
		Acl:    BucketACLPrivate,
	}
	result, err = client.PutBucketAcl(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))

	info, err = client.GetBucketInfo(context.TODO(), infoRequest)
	assert.Nil(t, err)
	assert.Equal(t, string(BucketACLPrivate), *info.BucketInfo.ACL)

	delRequest = &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	request = &PutBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.PutBucketAcl(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
}

func TestGetBucketAcl(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)

	request := &GetBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	result, err := client.GetBucketAcl(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, BucketACLType(*result.ACL), BucketACLPrivate)
	assert.NotEmpty(t, *result.Owner.ID)
	assert.NotEmpty(t, *result.Owner.DisplayName)

	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = client.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)

	request = &GetBucketAclRequest{
		Bucket: Ptr(bucketName),
	}
	result, err = client.GetBucketAcl(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	result, err := client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	var serr *ServiceError
	request = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	request = &PutObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	result, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	getRequest := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObject(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.NotEmpty(t, *result.ContentMD5)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.ContentLength, int64(len(content)))

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	request = &PutObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCopyObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	objectCopyName := objectNamePrefix + randLowStr(6) + "copy"

	copyRequest := &CopyObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Source: Ptr(objectCopyName),
	}
	result, err := client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(400), serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "Copy Source must mention the source bucket and key: /sourcebucket/sourcekey.", serr.Message)

	source := "/" + bucketName + "/" + objectName
	copyRequest = &CopyObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectCopyName),
		Source: Ptr(source),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	copyRequest = &CopyObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectCopyName),
		Source: Ptr(source),
	}
	_, err = client.CopyObject(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketCopyName := bucketNamePrefix + randLowStr(6) + "copy"
	putRequest = &PutBucketRequest{
		Bucket: Ptr(bucketCopyName),
	}

	client = getDefaultClient()
	_, err = client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	copyRequest = &CopyObjectRequest{
		Bucket: Ptr(bucketCopyName),
		Key:    Ptr(objectCopyName),
		Source: Ptr("/" + bucketName + "/" + objectName),
	}
	result, err = client.CopyObject(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)
	assert.Nil(t, result.VersionId)
}

func TestAppendObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	var result *AppendObjectResult
	content := randLowStr(100)
	request := &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position: Ptr(int64(0)),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.SSEKMSKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition := result.NextPosition

	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position: Ptr(nextPosition),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.SSEKMSKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Nil(t, result.SSEKMSKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*3))
	assert.NotEmpty(t, result.HashCRC64)

	objectName2 := objectName + "-kms-sm4"
	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName2),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position:                 Ptr(int64(0)),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.SSEKMSKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)))
	assert.NotEmpty(t, result.HashCRC64)

	nextPosition = result.NextPosition
	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName2),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position:                 Ptr(nextPosition),
		ServerSideDataEncryption: Ptr("SM4"),
		ServerSideEncryption:     Ptr("KMS"),
	}
	result, err = client.AppendObject(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.NotEmpty(t, result.SSEKMSKeyId)
	assert.Nil(t, result.VersionId)
	assert.Equal(t, result.NextPosition, int64(len(content)*2))
	assert.NotEmpty(t, result.HashCRC64)

	var serr *ServiceError
	request = &AppendObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "PositionNotEqualToLength", serr.Code)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketName + "-not-exist"
	request = &AppendObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		Position: Ptr(int64(0)),
	}
	_, err = client.AppendObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(1000)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest := &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	objectNameNotExist := objectNamePrefix + randLowStr(6) + "-not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	delRequest = &DeleteObjectRequest{
		Bucket:    Ptr(bucketName),
		Key:       Ptr(objectName),
		VersionId: Ptr("null"),
	}
	result, err = client.DeleteObject(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, result.StatusCode)
	assert.Equal(t, "204 No Content", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectNamePrefix),
	}
	_, err = client.DeleteObject(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestDeleteMultipleObjects(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest := &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketName),
		Objects: []DeleteObject{{Key: Ptr(objectName)}},
	}
	result, err := client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectName)

	str := "\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"
	objectNameSpecial := objectNamePrefix + randLowStr(6) + str
	content = randLowStr(10)
	request = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr(bucketName),
		Objects:      []DeleteObject{{Key: Ptr(objectNameSpecial)}},
		EncodingType: Ptr("url"),
	}
	result, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "200 OK", result.Status)
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, result.Headers.Get("Date"))
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, objectNameSpecial)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	delRequest = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr(bucketNameNotExist),
		Objects: []DeleteObject{{Key: Ptr(objectNameSpecial)}},
	}
	_, err = client.DeleteMultipleObjects(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestHeadObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.HeadObject(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ContentMD5)
	assert.NotEmpty(t, *result.ObjectType)
	assert.NotEmpty(t, *result.StorageClass)
	assert.NotEmpty(t, *result.ETag)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &HeadObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.HeadObject(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectMeta(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	headRequest := &GetObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectMeta(context.TODO(), headRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.ContentLength, int64(len(content)))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.LastModified)
	assert.NotEmpty(t, *result.HashCRC64)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	headRequest = &GetObjectMetaRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	result, err = client.GetObjectMeta(context.TODO(), headRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestRestoreObject(t *testing.T) {
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
	objectName := objectNamePrefix + randLowStr(6)
	content := randLowStr(10)
	request := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(content),
		},
		StorageClass: StorageClassColdArchive,
	}
	_, err = client.PutObject(context.TODO(), request)
	assert.Nil(t, err)

	restoreRequest := &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.RestoreObject(context.TODO(), restoreRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 202)
	assert.Equal(t, result.Status, "202 Accepted")
	assert.NotEmpty(t, result.Headers.Get("x-oss-request-id"))

	var serr *ServiceError
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(409), serr.StatusCode)
	assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
	assert.Equal(t, "The restore operation is in progress.", serr.Message)
	assert.NotEmpty(t, serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "not-exist"
	restoreRequest = &RestoreObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.RestoreObject(context.TODO(), restoreRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestPutObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	request := &PutObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicRead,
	}
	result, err := client.PutObjectAcl(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	infoRequest := &HeadObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	_, err = client.HeadObject(context.TODO(), infoRequest)
	assert.Nil(t, err)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	request = &PutObjectAclRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicRead,
	}
	_, err = client.PutObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestGetObjectAcl(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6)
	objectRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		Acl:    ObjectACLPublicReadWrite,
	}
	_, err = client.PutObject(context.TODO(), objectRequest)
	request := &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	result, err := client.GetObjectAcl(context.TODO(), request)
	assert.Nil(t, err)

	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, ObjectACLType(*result.ACL), ObjectACLPublicReadWrite)
	assert.NotEmpty(t, *result.Owner.ID)
	assert.NotEmpty(t, *result.Owner.DisplayName)

	objectNameNotExist := objectName + "-not-exist"
	request = &GetObjectAclRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameNotExist),
	}
	result, err = client.GetObjectAcl(context.TODO(), request)
	assert.NotNil(t, err)
	var serr *ServiceError
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchKey", serr.Code)
	assert.Equal(t, "The specified key does not exist.", serr.Message)
	assert.Equal(t, "0026-00000001", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestInitiateMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, initResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *initResult.Bucket, bucketName)
	assert.Equal(t, *initResult.Key, objectName)
	assert.NotEmpty(t, *initResult.UploadId)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	var serr *ServiceError
	bucketNameNotExist := bucketNamePrefix + randLowStr(6) + "-not-exist"
	initRequest = &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectName),
	}
	_, err = client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPart(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)

	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader("upload part 1"),
		},
	}

	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, initResult.StatusCode)
	assert.NotEmpty(t, partResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *partResult.ETag)
	assert.NotEmpty(t, *partResult.ContentMD5)
	assert.NotEmpty(t, *partResult.HashCRC64)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	var serr *ServiceError
	abortRequest = &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestUploadPartCopy(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)

	body := randLowStr(100000)
	objectSrcName := objectNamePrefix + randLowStr(6) + "src"
	objRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectSrcName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(body),
		},
	}
	_, err = client.PutObject(context.TODO(), objRequest)

	objectDestName := objectNamePrefix + randLowStr(6) + "dest"
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	source := "/" + bucketName + "/" + objectSrcName
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Source:     Ptr(source),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, copyResult.StatusCode)
	assert.NotEmpty(t, copyResult.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *copyResult.ETag)
	assert.NotEmpty(t, *copyResult.LastModified)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectDestName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	var serr *ServiceError
	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		Source:     Ptr(source),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000311", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestCompleteMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)

	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := ioutil.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]
	objectName := objectNamePrefix + randLowStr(6)

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part1)),
		},
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part2)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part3)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	request := &CompleteMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Part: parts,
		},
	}
	result, err := client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectName)

	getObj := &GetObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	getObjresult, err := client.GetObject(context.TODO(), getObj)
	assert.Nil(t, err)
	data, _ := ioutil.ReadAll(getObjresult.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(data), body)

	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}
	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	source := "/" + bucketName + "/" + objectName
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		Source:     Ptr(source),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.Nil(t, err)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))
	assert.NotEmpty(t, *result.ETag)
	assert.NotEmpty(t, *result.Location)
	assert.Equal(t, *result.Bucket, bucketName)
	assert.Equal(t, *result.Key, objectDestName)

	initCopyResult, err = client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	copyRequest = &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		Source:     Ptr(source),
	}
	copyResult, err := client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)
	copyPart := UploadPart{
		PartNumber: copyRequest.PartNumber,
		ETag:       copyResult.ETag,
	}

	var serr *ServiceError
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		CompleteMultipartUpload: &CompleteMultipartUpload{
			Part: []UploadPart{
				copyPart,
			},
		},
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 400, serr.StatusCode)
	assert.Equal(t, "InvalidArgument", serr.Code)
	assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
	assert.Equal(t, "0042-00000216", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	request = &CompleteMultipartUploadRequest{
		Bucket:      Ptr(bucketName),
		Key:         Ptr(objectDestName),
		UploadId:    Ptr(*initCopyResult.UploadId),
		CompleteAll: Ptr("yes"),
		Callback:    Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
	}
	result, err = client.CompleteMultipartUpload(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, 203, serr.StatusCode)
	assert.Equal(t, "CallbackFailed", serr.Code)
	assert.Equal(t, "Error status : 301.", serr.Message)
	assert.Equal(t, "0007-00000203", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAbortMultipartUpload(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6)
	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	result, err := client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.NotEmpty(t, result.Headers.Get(HeaderOssRequestID))

	var serr *ServiceError
	abortRequest = &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchUpload", serr.Code)
	assert.Equal(t, "The specified upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.", serr.Message)
	assert.Equal(t, "0042-00000002", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListMultipartUploads(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6) + "\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := ioutil.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part1)),
		},
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part2)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part3)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)

	putObj := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(randLowStr(1000)),
		},
	}

	_, err = client.PutObject(context.TODO(), putObj)
	assert.Nil(t, err)
	objectDestName := objectNamePrefix + randLowStr(6) + "dest" + "\f\v\n"
	initCopyRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectDestName),
	}

	initCopyResult, err := client.InitiateMultipartUpload(context.TODO(), initCopyRequest)
	assert.Nil(t, err)
	source := "/" + bucketName + "/" + objectName
	copyRequest := &UploadPartCopyRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectDestName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initCopyResult.UploadId),
		Source:     Ptr(source),
	}
	_, err = client.UploadPartCopy(context.TODO(), copyRequest)
	assert.Nil(t, err)

	listRequest := &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListMultipartUploads(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Empty(t, *listResult.KeyMarker, bucketName)
	assert.Len(t, listResult.Uploads, 2)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListMultipartUploadsRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListMultipartUploads(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestListParts(t *testing.T) {
	after := before(t)
	defer after(t)
	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	objectName := objectNamePrefix + randLowStr(6) + "-\v\n\f"
	body := randLowStr(400000)
	reader := strings.NewReader(body)
	bufReader := bufio.NewReader(reader)
	content, err := ioutil.ReadAll(bufReader)
	assert.Nil(t, err)
	count := 3
	partSize := len(content) / count
	part1 := content[:partSize]
	part2 := content[partSize : 2*partSize]
	part3 := content[2*partSize:]

	initRequest := &InitiateMultipartUploadRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectName),
	}
	initResult, err := client.InitiateMultipartUpload(context.TODO(), initRequest)
	assert.Nil(t, err)
	partRequest := &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(1),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part1)),
		},
	}
	var parts []UploadPart
	partResult, err := client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part := UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(2),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part2)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)
	assert.Nil(t, err)
	part = UploadPart{
		PartNumber: partRequest.PartNumber,
		ETag:       partResult.ETag,
	}
	parts = append(parts, part)
	partRequest = &UploadPartRequest{
		Bucket:     Ptr(bucketName),
		Key:        Ptr(objectName),
		PartNumber: int32(3),
		UploadId:   Ptr(*initResult.UploadId),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(string(part3)),
		},
	}
	partResult, err = client.UploadPart(context.TODO(), partRequest)

	listRequest := &ListPartsRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err := client.ListParts(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get(HeaderOssRequestID))
	assert.Equal(t, *listResult.Bucket, bucketName)
	assert.Equal(t, *listResult.Key, objectName)
	assert.Equal(t, *listResult.UploadId, *initResult.UploadId)
	assert.Equal(t, *listResult.StorageClass, "Standard")
	assert.Equal(t, listResult.IsTruncated, false)
	assert.Equal(t, listResult.PartNumberMarker, int32(0))
	assert.Equal(t, listResult.NextPartNumberMarker, int32(3))
	assert.Equal(t, listResult.MaxParts, int32(1000))
	assert.Len(t, listResult.Parts, count)

	abortRequest := &AbortMultipartUploadRequest{
		Bucket:   Ptr(bucketName),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	_, err = client.AbortMultipartUpload(context.TODO(), abortRequest)
	assert.Nil(t, err)

	bucketNameNotExist := bucketName + "-not-exist"
	listRequest = &ListPartsRequest{
		Bucket:   Ptr(bucketNameNotExist),
		Key:      Ptr(objectName),
		UploadId: Ptr(*initResult.UploadId),
	}
	listResult, err = client.ListParts(context.TODO(), listRequest)
	var serr *ServiceError
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
