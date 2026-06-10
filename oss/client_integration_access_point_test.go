//go:build integration

package oss

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func TestAccessPoint(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listResult, err := client.ListAccessPoints(context.TODO(), &ListAccessPointsRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	policy := `{"Version":"1","Statement":[{"Action":["oss:PutObject","oss:GetObject"],"Effect":"Deny","Principal":["` + accountID_ + `"],"Resource":["acs:oss:` + region_ + `:` + accountID_ + `:accesspoint/` + accessPointName + `","acs:oss:` + region_ + `:` + accountID_ + `:accesspoint/` + accessPointName + `/object/*"]}]}`
	putPolicyResult, err := client.PutAccessPointPolicy(context.TODO(), &PutAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
		Body:            strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, putPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getPolicyResult, err := client.GetAccessPointPolicy(context.TODO(), &GetAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, getPolicyResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, policy, getPolicyResult.Body)
	time.Sleep(1 * time.Second)

	delPolicyResult, err := client.DeleteAccessPointPolicy(context.TODO(), &DeleteAccessPointPolicyRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPolicyResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	for {
		getResult, err = client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		if *getResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	createResult, err = client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketNameNotExist),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listResult, err = client.ListAccessPoints(context.TODO(), &ListAccessPointsRequest{
		Bucket: Ptr(bucketNameNotExist),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	putPolicyResult, err = client.PutAccessPointPolicy(context.TODO(), &PutAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
		Body:            strings.NewReader(policy),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getPolicyResult, err = client.GetAccessPointPolicy(context.TODO(), &GetAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delPolicyResult, err = client.DeleteAccessPointPolicy(context.TODO(), &DeleteAccessPointPolicyRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestAccessPointPublicAccessBlock(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	client := getDefaultClient()
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	putResult, err := client.PutAccessPointPublicAccessBlock(context.TODO(), &PutAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getResult, err := client.GetAccessPointPublicAccessBlock(context.TODO(), &GetAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.True(t, *getResult.PublicAccessBlockConfiguration.BlockPublicAccess)
	time.Sleep(1 * time.Second)

	delResult, err := client.DeleteAccessPointPublicAccessBlock(context.TODO(), &DeleteAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putResult, err = client.PutAccessPointPublicAccessBlock(context.TODO(), &PutAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
		PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
			Ptr(true),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getResult, err = client.GetAccessPointPublicAccessBlock(context.TODO(), &GetAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delResult, err = client.DeleteAccessPointPublicAccessBlock(context.TODO(), &DeleteAccessPointPublicAccessBlockRequest{
		Bucket:          Ptr(bucketNameNotExist),
		AccessPointName: Ptr(accessPointName),
	})
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	for {
		getPointResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		assert.Nil(t, err)
		if *getPointResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delPointResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPointResult.StatusCode)
	time.Sleep(1 * time.Second)
}

func TestAccessPointForObjectProcess(t *testing.T) {
	if apEnable == "" {
		return
	}
	after := before(t)
	defer after(t)
	defer clearAp(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)
	accessPointName := "ap-01-" + randLowStr(5)
	createResult, err := client.CreateAccessPoint(context.TODO(), &CreateAccessPointRequest{
		Bucket: Ptr(bucketName),
		CreateAccessPointConfiguration: &CreateAccessPointConfiguration{
			AccessPointName: Ptr(accessPointName),
			NetworkOrigin:   Ptr("internet"),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	assert.NotEmpty(t, createResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	objectProcessName := "fc-ap-01-" + randLowStr(5)
	arn := "acs:fc:" + region_ + ":" + accountID_ + ":services/test-oss-fc.LATEST/functions/" + objectProcessName
	roleArn := "acs:ram::" + accountID_ + ":role/aliyunfcdefaultrole"
	createObjectResult, err := client.CreateAccessPointForObjectProcess(context.TODO(), &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
			AccessPointName: Ptr(accessPointName),
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: &ObjectProcessAllowedFeatures{
					[]string{"GetObject-Range"},
				},
				TransformationConfigurations: &TransformationConfigurations{
					[]TransformationConfiguration{
						{
							Actions: &AccessPointActions{
								[]string{"GetObject"},
							},
							ContentTransformation: &ContentTransformation{
								FunctionCompute: &ObjectProcessFunctionCompute{
									FunctionArn:           Ptr(arn),
									FunctionAssumeRoleArn: Ptr(roleArn),
								},
							},
						},
					},
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createObjectResult.StatusCode)
	assert.NotEmpty(t, createObjectResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getResult, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	listResult, err := client.ListAccessPointsForObjectProcess(context.TODO(), &ListAccessPointsForObjectProcessRequest{})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	policy := `{"Version":"1","Statement":[{"Action":["oss:GetObject"],"Effect":"Allow","Principal":["` + accountID_ + `"],"Resource":["acs:oss:` + region_ + `:` + accountID_ + `:accesspointforobjectprocess/` + objectProcessName + `/object/*"]}]}`
	putPolicyResult, err := client.PutAccessPointPolicyForObjectProcess(context.TODO(), &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		Body:                            strings.NewReader(policy),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, putPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getPolicyResult, err := client.GetAccessPointPolicyForObjectProcess(context.TODO(), &GetAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putPolicyResult.StatusCode)
	assert.NotEmpty(t, getPolicyResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, getPolicyResult.Body, policy)
	time.Sleep(1 * time.Second)
	delPolicyResult, err := client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delPolicyResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	putConfigResult, err := client.PutAccessPointConfigForObjectProcess(context.TODO(), &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		PutAccessPointConfigForObjectProcessConfiguration: &PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: &ObjectProcessAllowedFeatures{
					[]string{"GetObject-Range"},
				},
				TransformationConfigurations: &TransformationConfigurations{
					[]TransformationConfiguration{
						{
							Actions: &AccessPointActions{
								[]string{"GetObject"},
							},
							ContentTransformation: &ContentTransformation{
								FunctionCompute: &ObjectProcessFunctionCompute{
									FunctionArn:           Ptr(arn),
									FunctionAssumeRoleArn: Ptr(roleArn),
								},
							},
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, putConfigResult.StatusCode)
	assert.NotEmpty(t, putConfigResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	getConfigResult, err := client.GetAccessPointConfigForObjectProcess(context.TODO(), &GetAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getConfigResult.StatusCode)
	assert.NotEmpty(t, getConfigResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	for {
		getResult, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr(bucketName),
			AccessPointForObjectProcessName: Ptr(objectProcessName),
		})
		assert.Nil(t, err)
		if *getResult.AccessPointForObjectProcessStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	delResult, err := client.DeleteAccessPointForObjectProcess(context.TODO(), &DeleteAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketName),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, delPolicyResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)
	var serr *ServiceError
	for {
		_, err := client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
			Bucket:                          Ptr(bucketName),
			AccessPointForObjectProcessName: Ptr(objectProcessName),
		})
		if err != nil {
			errors.As(err, &serr)
			if serr.StatusCode == 404 && serr.Code == "NoSuchAccessPointForObjectProcess" {
				break
			}
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	for {
		gResult, err := client.GetAccessPoint(context.TODO(), &GetAccessPointRequest{
			Bucket:          Ptr(bucketName),
			AccessPointName: Ptr(accessPointName),
		})
		assert.Nil(t, err)
		if *gResult.AccessPointStatus != "creating" {
			break
		} else {
			time.Sleep(3 * time.Second)
		}
	}
	dResult, err := client.DeleteAccessPoint(context.TODO(), &DeleteAccessPointRequest{
		Bucket:          Ptr(bucketName),
		AccessPointName: Ptr(accessPointName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 204, dResult.StatusCode)
	time.Sleep(1 * time.Second)
	bucketNameNotExist := bucketName + "-not-exist"
	_, err = client.CreateAccessPointForObjectProcess(context.TODO(), &CreateAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		CreateAccessPointForObjectProcessConfiguration: &CreateAccessPointForObjectProcessConfiguration{
			AccessPointName: Ptr(accessPointName),
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: &ObjectProcessAllowedFeatures{
					[]string{"GetObject-Range"},
				},
				TransformationConfigurations: &TransformationConfigurations{
					[]TransformationConfiguration{
						{
							Actions: &AccessPointActions{
								[]string{"GetObject"},
							},
							ContentTransformation: &ContentTransformation{
								FunctionCompute: &ObjectProcessFunctionCompute{
									FunctionArn:           Ptr(arn),
									FunctionAssumeRoleArn: Ptr(roleArn),
								},
							},
						},
					},
				},
			},
		},
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	getResult, err = client.GetAccessPointForObjectProcess(context.TODO(), &GetAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	noPermClient := getClientWithCredentialsProvider(region_, endpoint_,
		credentials.NewStaticCredentialsProvider("ak", "sk"))
	listResult, err = noPermClient.ListAccessPointsForObjectProcess(context.TODO(), &ListAccessPointsForObjectProcessRequest{})
	serr = &ServiceError{}
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.Equal(t, "0002-00000902", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	putPolicyResult, err = client.PutAccessPointPolicyForObjectProcess(context.TODO(), &PutAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		Body:                            strings.NewReader(policy),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	getPolicyResult, err = client.GetAccessPointPolicyForObjectProcess(context.TODO(), &GetAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	delPolicyResult, err = client.DeleteAccessPointPolicyForObjectProcess(context.TODO(), &DeleteAccessPointPolicyForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	putConfigResult, err = client.PutAccessPointConfigForObjectProcess(context.TODO(), &PutAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
		PutAccessPointConfigForObjectProcessConfiguration: &PutAccessPointConfigForObjectProcessConfiguration{
			ObjectProcessConfiguration: &ObjectProcessConfiguration{
				AllowedFeatures: &ObjectProcessAllowedFeatures{
					[]string{"GetObject-Range"},
				},
				TransformationConfigurations: &TransformationConfigurations{
					[]TransformationConfiguration{
						{
							Actions: &AccessPointActions{
								[]string{"GetObject"},
							},
							ContentTransformation: &ContentTransformation{
								FunctionCompute: &ObjectProcessFunctionCompute{
									FunctionArn:           Ptr(arn),
									FunctionAssumeRoleArn: Ptr(roleArn),
								},
							},
						},
					},
				},
			},
			PublicAccessBlockConfiguration: &PublicAccessBlockConfiguration{
				Ptr(true),
			},
		},
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	_, err = client.GetAccessPointConfigForObjectProcess(context.TODO(), &GetAccessPointConfigForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)
	_, err = client.DeleteAccessPointForObjectProcess(context.TODO(), &DeleteAccessPointForObjectProcessRequest{
		Bucket:                          Ptr(bucketNameNotExist),
		AccessPointForObjectProcessName: Ptr(objectProcessName),
	})
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
