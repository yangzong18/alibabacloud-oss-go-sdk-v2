//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketReplication(t *testing.T) {
	after := before(t)
	defer after(t)
	//TODO
	bucketName := bucketNamePrefix + randLowStr(6)
	request := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}
	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	targetBucketName := bucketNamePrefix + "-target-" + randLowStr(6)
	request = &PutBucketRequest{
		Bucket: Ptr(targetBucketName),
	}
	client1 := getClient("cn-beijing", "https://oss-cn-beijing.aliyuncs.com")
	_, err = client1.PutBucket(context.TODO(), request)
	assert.Nil(t, err)

	putRequest := &PutBucketReplicationRequest{
		Bucket: Ptr(bucketName),
		ReplicationConfiguration: &ReplicationConfiguration{
			[]ReplicationRule{
				{
					RTC: &ReplicationTimeControl{
						Status: Ptr("enabled"),
					},
					Destination: &ReplicationDestination{
						Bucket:       Ptr(targetBucketName),
						Location:     Ptr("oss-cn-beijing"),
						TransferType: TransferTypeInternal,
					},
					HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					SourceSelectionCriteria: &ReplicationSourceSelectionCriteria{
						&SseKmsEncryptedObjects{
							Status: StatusEnabled,
						},
					},
				},
			},
		},
	}
	putResult, err := client.PutBucketReplication(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketReplicationRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketReplication(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.Equal(t, 1, len(getResult.ReplicationConfiguration.Rules))
	time.Sleep(1 * time.Second)

	getLocationRequest := &GetBucketReplicationLocationRequest{
		Bucket: Ptr(bucketName),
	}
	getLocationResult, err := client.GetBucketReplicationLocation(context.TODO(), getLocationRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getLocationResult.StatusCode)
	time.Sleep(1 * time.Second)

	getProgressRequest := &GetBucketReplicationProgressRequest{
		Bucket: Ptr(bucketName),
		RuleId: getResult.ReplicationConfiguration.Rules[0].ID,
	}
	getProgressResult, err := client.GetBucketReplicationProgress(context.TODO(), getProgressRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getProgressResult.StatusCode)
	assert.Equal(t, 1, len(getProgressResult.ReplicationProgress.Rules))
	time.Sleep(1 * time.Second)

	rtcRequest := &PutBucketRtcRequest{
		Bucket: Ptr(bucketName),
		RtcConfiguration: &RtcConfiguration{
			RTC: &ReplicationTimeControl{
				Status: Ptr("disabled"),
			},
			ID: getResult.ReplicationConfiguration.Rules[0].ID,
		},
	}
	rtcResult, err := client.PutBucketRtc(context.TODO(), rtcRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, rtcResult.StatusCode)
	assert.NotEmpty(t, rtcResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketReplicationRequest{
		Bucket: Ptr(bucketName),
		ReplicationRules: &ReplicationRules{
			[]string{*getResult.ReplicationConfiguration.Rules[0].ID},
		},
	}
	delResult, err := client.DeleteBucketReplication(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, delResult.StatusCode)
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
		ReplicationConfiguration: &ReplicationConfiguration{
			[]ReplicationRule{
				{
					RTC: &ReplicationTimeControl{
						Status: Ptr("enabled"),
					},
					Destination: &ReplicationDestination{
						Bucket:       Ptr(targetBucketName),
						Location:     Ptr("oss-cn-cn-hangzhou"),
						TransferType: TransferTypeOssAcc,
					},
					HistoricalObjectReplication: HistoricalObjectReplicationEnabled,
					SourceSelectionCriteria: &ReplicationSourceSelectionCriteria{
						&SseKmsEncryptedObjects{
							Status: StatusEnabled,
						},
					},
				},
			},
		},
	}
	putResult, err = client.PutBucketReplication(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	_, err = client.GetBucketReplication(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getLocationRequest = &GetBucketReplicationLocationRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	getLocationResult, err = client.GetBucketReplicationLocation(context.TODO(), getLocationRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getProgressRequest = &GetBucketReplicationProgressRequest{
		Bucket: Ptr(bucketNameNotExist),
		RuleId: getResult.ReplicationConfiguration.Rules[0].ID,
	}
	getProgressResult, err = client.GetBucketReplicationProgress(context.TODO(), getProgressRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	rtcRequest = &PutBucketRtcRequest{
		Bucket: Ptr(bucketNameNotExist),
		RtcConfiguration: &RtcConfiguration{
			RTC: &ReplicationTimeControl{
				Status: Ptr("disabled"),
			},
			ID: getResult.ReplicationConfiguration.Rules[0].ID,
		},
	}
	rtcResult, err = client.PutBucketRtc(context.TODO(), rtcRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketReplicationRequest{
		Bucket: Ptr(bucketNameNotExist),
		ReplicationRules: &ReplicationRules{
			[]string{*getResult.ReplicationConfiguration.Rules[0].ID},
		},
	}
	delResult, err = client.DeleteBucketReplication(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
