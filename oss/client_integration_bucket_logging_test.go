//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketLogging(t *testing.T) {
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

	request := &PutBucketLoggingRequest{
		Bucket: Ptr(bucketName),
		BucketLoggingStatus: &BucketLoggingStatus{
			&LoggingEnabled{
				TargetBucket: Ptr(bucketName),
				TargetPrefix: Ptr("TargetPrefix"),
			},
		},
	}
	result, err := client.PutBucketLogging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketLoggingRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err := client.GetBucketLogging(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetBucket, bucketName)
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "TargetPrefix")
	time.Sleep(1 * time.Second)

	request = &PutBucketLoggingRequest{
		Bucket: Ptr(bucketName),
		BucketLoggingStatus: &BucketLoggingStatus{
			&LoggingEnabled{
				TargetBucket: Ptr(bucketName),
				TargetPrefix: Ptr("TargetPrefix"),
				LoggingRole:  Ptr("AliyunOSSLoggingDefaultRole"),
			},
		},
	}
	result, err = client.PutBucketLogging(context.TODO(), request)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.NotEmpty(t, result.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketLoggingRequest{
		Bucket: Ptr(bucketName),
	}
	getResult, err = client.GetBucketLogging(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetBucket, bucketName)
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.TargetPrefix, "TargetPrefix")
	assert.Equal(t, *getResult.BucketLoggingStatus.LoggingEnabled.LoggingRole, "AliyunOSSLoggingDefaultRole")
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketLoggingRequest{
		Bucket: Ptr(bucketName),
	}
	delResult, err := client.DeleteBucketLogging(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.Equal(t, "204 No Content", delResult.Status)
	assert.NotEmpty(t, delResult.Headers.Get("x-oss-request-id"))
	assert.NotEmpty(t, delResult.Headers.Get("Date"))
	time.Sleep(1 * time.Second)

	putUserRequest := &PutUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
		UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
			HeaderSet: &LoggingHeaderSet{
				[]string{"header1", "header2", "header3"},
			},
			ParamSet: &LoggingParamSet{
				[]string{"param1", "param2"},
			},
		},
	}
	putUserResult, err := client.PutUserDefinedLogFieldsConfig(context.TODO(), putUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putUserResult.StatusCode)
	assert.NotEmpty(t, putUserResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getUserRequest := &GetUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	getUserResult, err := client.GetUserDefinedLogFieldsConfig(context.TODO(), getUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getUserResult.StatusCode)
	assert.NotEmpty(t, getUserResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, 3, len(getUserResult.UserDefinedLogFieldsConfiguration.HeaderSet.Headers))
	assert.Equal(t, 2, len(getUserResult.UserDefinedLogFieldsConfiguration.ParamSet.Parameters))
	time.Sleep(1 * time.Second)

	delUserRequest := &DeleteUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketName),
	}
	delUserResult, err := client.DeleteUserDefinedLogFieldsConfig(context.TODO(), delUserRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delUserResult.StatusCode)
	assert.Equal(t, "204 No Content", delUserResult.Status)
	assert.NotEmpty(t, delUserResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &PutBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
		BucketLoggingStatus: &BucketLoggingStatus{
			&LoggingEnabled{
				TargetBucket: Ptr("TargetBucket"),
				TargetPrefix: Ptr("TargetPrefix"),
			},
		},
	}
	result, err = client.PutBucketLogging(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getRequest = &GetBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	getResult, err = client.GetBucketLogging(context.TODO(), getRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketLoggingRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	delResult, err = client.DeleteBucketLogging(context.TODO(), delRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	putUserRequest = &PutUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
		UserDefinedLogFieldsConfiguration: &UserDefinedLogFieldsConfiguration{
			HeaderSet: &LoggingHeaderSet{
				[]string{"header1", "header2", "header3"},
			},
			ParamSet: &LoggingParamSet{
				[]string{"param1", "param2"},
			},
		},
	}
	serr = &ServiceError{}
	putUserResult, err = client.PutUserDefinedLogFieldsConfig(context.TODO(), putUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	getUserRequest = &GetUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	serr = &ServiceError{}
	getUserResult, err = client.GetUserDefinedLogFieldsConfig(context.TODO(), getUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delUserRequest = &DeleteUserDefinedLogFieldsConfigRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	delUserResult, err = client.DeleteUserDefinedLogFieldsConfig(context.TODO(), delUserRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
