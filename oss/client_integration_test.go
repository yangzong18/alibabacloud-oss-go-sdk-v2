//go:build integration

package oss

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var (
	// Endpoint/ID/Key
	region_           = os.Getenv("OSS_TEST_REGION")
	endpoint_         = os.Getenv("OSS_TEST_ENDPOINT")
	accessID_         = os.Getenv("OSS_TEST_ACCESS_KEY_ID")
	accessKey_        = os.Getenv("OSS_TEST_ACCESS_KEY_SECRET")
	ramRoleArn_       = os.Getenv("OSS_TEST_RAM_ROLE_ARN")
	accountID_        = os.Getenv("OSS_TEST_RAM_UID")
	signatureVersion_ = os.Getenv("OSS_TEST_SIGNATURE_VERSION")

	// payer
	payerAccessID_  = os.Getenv("OSS_TEST_PAYER_ACCESS_KEY_ID")
	payerAccessKey_ = os.Getenv("OSS_TEST_PAYER_ACCESS_KEY_SECRET")
	payerUID_       = os.Getenv("OSS_TEST_PAYER_UID")

	// path style
	pathStyleBucket_ = os.Getenv("OSS_TEST_PATHSTYLE_BUCKET")
	pathStyleRegion_ = os.Getenv("OSS_TEST_PATHSTYLE_REGION")

	apEnable = os.Getenv("OSS_TEST_AP_ENABLE")

	instance_ *Client
	testOnce_ sync.Once

	kmdIdMap_ = map[string]string{}
)

var (
	bucketNamePrefix = os.Getenv("OSS_TEST_BUCKET_NAME_PREFIX")
	objectNamePrefix = os.Getenv("OSS_TEST_OBJECT_NAME_PREFIX")
)

func getDefaultClient() *Client {
	testOnce_.Do(func() {
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
			WithRegion(region_).
			WithEndpoint(endpoint_).
			WithSignatureVersion(getSignatrueVersion())

		instance_ = NewClient(cfg)
	})
	return instance_
}

func getClient(region, endpoint string) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessID_, accessKey_)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getClientUseStsToken(region, endpoint string) *Client {
	resp, err := stsAssumeRole(accessID_, accessKey_, ramRoleArn_)
	if err != nil {
		return nil
	}
	accessId := resp.Credentials.AccessKeyId
	accessKey := resp.Credentials.AccessKeySecret
	token := resp.Credentials.SecurityToken
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token)).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getClientUseStsTokenV2(cfg *Config) *Client {
	resp, err := stsAssumeRole(accessID_, accessKey_, ramRoleArn_)
	if err != nil {
		return nil
	}
	accessId := resp.Credentials.AccessKeyId
	accessKey := resp.Credentials.AccessKeySecret
	token := resp.Credentials.SecurityToken
	cfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token))
	/*
		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessId, accessKey, token)).
			WithRegion(region).
			WithEndpoint(endpoint).
			WithSignatureVersion(getSignatrueVersion())
	*/
	return NewClient(cfg)
}

func getClientWithCredentialsProvider(region, endpoint string, cred credentials.CredentialsProvider) *Client {
	cfg := LoadDefaultConfig().
		WithCredentialsProvider(cred).
		WithRegion(region).
		WithEndpoint(endpoint).
		WithSignatureVersion(getSignatrueVersion())

	return NewClient(cfg)
}

func getKmsID(region string) string {
	if id, ok := kmdIdMap_[region]; ok {
		return id
	}

	client := getClient(region, fmt.Sprintf("oss-%s.aliyuncs.com", region))
	bucketName := bucketNamePrefix + randLowStr(6)

	if _, err := client.PutBucket(context.TODO(), &PutBucketRequest{Bucket: Ptr(bucketName)}); err != nil {
		return ""
	}

	kmdId := ""
	if _, err := client.PutObject(context.TODO(), &PutObjectRequest{
		Bucket:               Ptr(bucketName),
		Key:                  Ptr("kms-id"),
		ServerSideEncryption: Ptr("KMS")}); err == nil {

		if result, err := client.HeadObject(context.TODO(), &HeadObjectRequest{
			Bucket: Ptr(bucketName),
			Key:    Ptr("kms-id")}); err == nil {
			kmdId = ToString(result.ServerSideEncryptionKeyId)
			kmdIdMap_[region] = kmdId
		}
	}
	client.DeleteObject(context.TODO(), &DeleteObjectRequest{Bucket: Ptr(bucketName), Key: Ptr("kms-id")})
	client.DeleteBucket(context.TODO(), &DeleteBucketRequest{Bucket: Ptr(bucketName)})
	return kmdId
}

func getSignatrueVersion() SignatureVersionType {
	switch signatureVersion_ {
	case "v1":
		return SignatureVersionV1
	default:
		return SignatureVersionV4
	}
}

func cleanBucket(bucketInfo BucketProperties, t *testing.T) {
	assert.NotEmpty(t, *bucketInfo.Name)
	var c *Client
	if strings.Contains(endpoint_, *bucketInfo.ExtranetEndpoint) ||
		strings.Contains(endpoint_, *bucketInfo.IntranetEndpoint) {
		c = getDefaultClient()
	} else {
		c = getClient(*bucketInfo.Region, *bucketInfo.ExtranetEndpoint)
	}
	assert.NotNil(t, c)
	cleanObjects(c, *bucketInfo.Name, t)
}

func deleteBucket(bucketName string, t *testing.T) {
	assert.NotEmpty(t, bucketName)
	var c *Client
	c = getDefaultClient()
	assert.NotNil(t, c)
	cleanObjects(c, bucketName, t)
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

func cleanObjects(c *Client, bucketName string, t *testing.T) {
	var err error
	var listRequest *ListObjectsRequest
	var delObjRequest *DeleteObjectRequest
	var lor *ListObjectsResult
	marker := ""
	for {
		listRequest = &ListObjectsRequest{
			Bucket: Ptr(bucketName),
			Marker: Ptr(marker),
		}
		lor, err = c.ListObjects(context.TODO(), listRequest)
		assert.Nil(t, err)
		var deleteObjects []DeleteObject
		for _, object := range lor.Contents {
			deleteObjects = append(deleteObjects, DeleteObject{Key: object.Key})
		}
		if len(deleteObjects) > 0 {
			_, err = c.DeleteMultipleObjects(context.TODO(), &DeleteMultipleObjectsRequest{
				Bucket:  Ptr(bucketName),
				Objects: deleteObjects,
			})
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
			Bucket:         Ptr(bucketName),
			KeyMarker:      Ptr(keyMarker),
			UploadIdMarker: Ptr(uploadIdMarker),
		}
		lsRes, err = c.ListMultipartUploads(context.TODO(), listUploadRequest)
		assert.Nil(t, err)
		for _, upload := range lsRes.Uploads {
			abortRequest = &AbortMultipartUploadRequest{
				Bucket:   Ptr(bucketName),
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
	var lsVersionRq *ListObjectVersionsRequest
	var lsVersionRs *ListObjectVersionsResult
	versionKeyMarker := ""
	VersionIdMarker := ""
	for {
		lsVersionRq = &ListObjectVersionsRequest{
			Bucket:          Ptr(bucketName),
			KeyMarker:       Ptr(versionKeyMarker),
			VersionIdMarker: Ptr(VersionIdMarker),
		}
		lsVersionRs, err = c.ListObjectVersions(context.TODO(), lsVersionRq)
		assert.Nil(t, err)
		for _, object := range lsVersionRs.ObjectDeleteMarkers {
			delObjRequest = &DeleteObjectRequest{
				Bucket:    Ptr(bucketName),
				Key:       Ptr(*object.Key),
				VersionId: Ptr(*object.VersionId),
			}
			_, err = c.DeleteObject(context.TODO(), delObjRequest)
			assert.Nil(t, err)
		}
		for _, object := range lsVersionRs.ObjectVersions {
			delObjRequest = &DeleteObjectRequest{
				Bucket:    Ptr(bucketName),
				Key:       Ptr(*object.Key),
				VersionId: Ptr(*object.VersionId),
			}
			_, err = c.DeleteObject(context.TODO(), delObjRequest)
			assert.Nil(t, err)
		}
		if !lsVersionRs.IsTruncated {
			break
		}
		versionKeyMarker = *lsVersionRs.NextKeyMarker
		VersionIdMarker = *lsVersionRs.NextVersionIdMarker
	}
	delRequest := &DeleteBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err = c.DeleteBucket(context.TODO(), delRequest)
	assert.Nil(t, err)
}

func deleteAccessPoint(c *Client, bucketName string, t *testing.T) {
	var err error
	var listRequest *ListAccessPointsRequest
	var lap *ListAccessPointsResult
	var delRequest *DeleteAccessPointRequest
	var token = ""
	var getResult *GetAccessPointResult
	var serr = &ServiceError{}
	for {
		listRequest = &ListAccessPointsRequest{
			Bucket:            Ptr(bucketName),
			ContinuationToken: Ptr(token),
		}
		lap, err = c.ListAccessPoints(context.TODO(), listRequest)
		assert.Nil(t, err)
		if len(lap.AccessPoints) > 0 {
			for _, accessPoint := range lap.AccessPoints {
				switch *accessPoint.Status {
				case "creating":
					time.Sleep(3 * time.Second)
					for {
						getResult, err = c.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
							Bucket:          Ptr(bucketName),
							AccessPointName: accessPoint.AccessPointName,
						})
						assert.Nil(t, err)
						if *getResult.AccessPointStatus != "creating" {
							break
						} else {
							time.Sleep(3 * time.Second)
						}
					}
					delRequest = &DeleteAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					}
					_, err = c.DeleteAccessPoint(context.TODO(), delRequest)
					assert.Nil(t, err)
				case "deleting":
					time.Sleep(5 * time.Second)
				default:
					delRequest = &DeleteAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					}
					_, err = c.DeleteAccessPoint(context.TODO(), delRequest)
					assert.Nil(t, err)
					time.Sleep(3 * time.Second)

				}
				for {
					getResult, err = c.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
						Bucket:          Ptr(bucketName),
						AccessPointName: accessPoint.AccessPointName,
					})
					if err != nil {
						errors.As(err, &serr)
						if serr.StatusCode == 404 && serr.Code == "NoSuchAccessPoint" {
							break
						}
					}
					time.Sleep(3 * time.Second)
				}
			}

			if !*lap.IsTruncated {
				break
			}
			if lap.NextContinuationToken != nil {
				token = *lap.NextContinuationToken
			}
		} else {
			break
		}

	}
}

type credentialsForSts struct {
	AccessKeyId     string
	AccessKeySecret string
	Expiration      time.Time
	SecurityToken   string
}

type assumedRoleUserForSts struct {
	Arn           string
	AssumedRoleId string
}

type responseForSts struct {
	Credentials     credentialsForSts
	AssumedRoleUser assumedRoleUserForSts
	RequestId       string
}

func stsAssumeRole(accessKeyId string, accessKeySecret string, roleArn string) (*responseForSts, error) {
	// StsSignVersion sts sign version
	StsSignVersion := "1.0"
	// StsAPIVersion sts api version
	StsAPIVersion := "2015-04-01"
	// StsHost sts host
	StsHost := "https://sts.aliyuncs.com/"
	// TimeFormat time fomrat
	TimeFormat := "2006-01-02T15:04:05Z"
	// RespBodyFormat  respone body format
	RespBodyFormat := "JSON"
	// PercentEncode '/'
	PercentEncode := "%2F"
	// HTTPGet http get method
	HTTPGet := "GET"
	rand.Seed(time.Now().UnixNano())
	uuid := fmt.Sprintf("Nonce-%d", rand.Intn(10000))
	queryStr := "SignatureVersion=" + StsSignVersion
	queryStr += "&Format=" + RespBodyFormat
	queryStr += "&Timestamp=" + url.QueryEscape(time.Now().UTC().Format(TimeFormat))
	queryStr += "&RoleArn=" + url.QueryEscape(roleArn)
	queryStr += "&RoleSessionName=" + "oss_test_sess"
	queryStr += "&AccessKeyId=" + accessKeyId
	queryStr += "&SignatureMethod=HMAC-SHA1"
	queryStr += "&Version=" + StsAPIVersion
	queryStr += "&Action=AssumeRole"
	queryStr += "&SignatureNonce=" + uuid
	queryStr += "&DurationSeconds=" + strconv.FormatInt(3600, 10)

	// Sort query string
	queryParams, err := url.ParseQuery(queryStr)
	if err != nil {
		return nil, err
	}

	strToSign := HTTPGet + "&" + PercentEncode + "&" + url.QueryEscape(queryParams.Encode())

	// Generate signature
	hashSign := hmac.New(sha1.New, []byte(accessKeySecret+"&"))
	hashSign.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(hashSign.Sum(nil))

	// Build url
	assumeURL := StsHost + "?" + queryStr + "&Signature=" + url.QueryEscape(signature)

	// Send Request
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(assumeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	// Handle Response
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	result := responseForSts{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func before(_ *testing.T) func(t *testing.T) {

	//fmt.Println("setup test case")
	return after
}

func after(t *testing.T) {
	cleanBuckets(bucketNamePrefix, t)
	//fmt.Println("teardown  test case")
}

func clearAp(t *testing.T) {
	c := getDefaultClient()
	for {
		request := &ListBucketsRequest{
			Prefix: Ptr(bucketNamePrefix),
		}
		result, err := c.ListBuckets(context.TODO(), request)
		assert.Nil(t, err)
		if len(result.Buckets) == 0 {
			return
		}
		for _, b := range result.Buckets {
			deleteAccessPoint(c, *b.Name, t)
		}
		if !result.IsTruncated {
			break
		}
	}
}

func dumpErrIfNotNil(err error) {
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
	}
}
