package oss

import (
	"context"
	"errors"
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
	var serr *ServiceError
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
