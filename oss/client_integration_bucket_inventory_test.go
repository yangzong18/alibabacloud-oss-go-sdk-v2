//go:build integration

package oss

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBucketInventory(t *testing.T) {
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

	id := "report1" + randStr(6)
	putRequest := &PutBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
		InventoryConfiguration: &InventoryConfiguration{
			Id:        Ptr(id),
			IsEnabled: Ptr(true),
			Filter: &InventoryFilter{
				Prefix:                   Ptr("filterPrefix"),
				LastModifyBeginTimeStamp: Ptr(int64(1637883649)),
				LastModifyEndTimeStamp:   Ptr(int64(1638347592)),
				LowerSizeBound:           Ptr(int64(1024)),
				UpperSizeBound:           Ptr(int64(1048576)),
				StorageClass:             Ptr("Standard,IA"),
			},
			Destination: &InventoryDestination{
				&InventoryOSSBucketDestination{
					Format:    InventoryFormatCSV,
					AccountId: Ptr(accountID_),
					RoleArn:   Ptr("acs:ram::" + accountID_ + ":role/AliyunOSSRole"),
					Bucket:    Ptr("acs:oss:::" + bucketName),
					Prefix:    Ptr("prefix1"),
				},
			},
			Schedule: &InventorySchedule{
				Frequency: InventoryFrequencyDaily,
			},
			IncludedObjectVersions: Ptr("All"),
			OptionalFields: &OptionalFields{
				Fields: []InventoryOptionalFieldType{
					InventoryOptionalFieldSize,
					InventoryOptionalFieldLastModifiedDate,
					InventoryOptionalFieldETag,
					InventoryOptionalFieldStorageClass,
					InventoryOptionalFieldIsMultipartUploaded,
					InventoryOptionalFieldEncryptionStatus,
				},
			},
		},
	}
	putResult, err := client.PutBucketInventory(context.TODO(), putRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, putResult.StatusCode)
	assert.NotEmpty(t, putResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	getRequest := &GetBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
	}
	getResult, err := client.GetBucketInventory(context.TODO(), getRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, getResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	listRequest := &ListBucketInventoryRequest{
		Bucket: Ptr(bucketName),
	}
	listResult, err := client.ListBucketInventory(context.TODO(), listRequest)
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	assert.Equal(t, len(listResult.ListInventoryConfigurationsResult.InventoryConfigurations), 1)
	time.Sleep(1 * time.Second)

	delRequest := &DeleteBucketInventoryRequest{
		Bucket:      Ptr(bucketName),
		InventoryId: Ptr(id),
	}
	delResult, err := client.DeleteBucketInventory(context.TODO(), delRequest)
	assert.Nil(t, err)
	assert.Equal(t, 204, delResult.StatusCode)
	assert.NotEmpty(t, listResult.Headers.Get("X-Oss-Request-Id"))
	time.Sleep(1 * time.Second)

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	putRequest = &PutBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
		InventoryConfiguration: &InventoryConfiguration{
			Id:        Ptr(id),
			IsEnabled: Ptr(true),
			Filter: &InventoryFilter{
				Prefix:                   Ptr("filterPrefix"),
				LastModifyBeginTimeStamp: Ptr(int64(1637883649)),
				LastModifyEndTimeStamp:   Ptr(int64(1638347592)),
				LowerSizeBound:           Ptr(int64(1024)),
				UpperSizeBound:           Ptr(int64(1048576)),
				StorageClass:             Ptr("Standard,IA"),
			},
			Destination: &InventoryDestination{
				&InventoryOSSBucketDestination{
					Format:    InventoryFormatCSV,
					AccountId: Ptr(accountID_),
					RoleArn:   Ptr("acs:ram::" + accountID_ + ":role/AliyunOSSRole"),
					Bucket:    Ptr("acs:oss:::" + bucketName),
					Prefix:    Ptr("prefix1"),
				},
			},
			Schedule: &InventorySchedule{
				Frequency: InventoryFrequencyDaily,
			},
			IncludedObjectVersions: Ptr("All"),
			OptionalFields: &OptionalFields{
				Fields: []InventoryOptionalFieldType{
					InventoryOptionalFieldSize,
					InventoryOptionalFieldLastModifiedDate,
					InventoryOptionalFieldETag,
					InventoryOptionalFieldStorageClass,
					InventoryOptionalFieldIsMultipartUploaded,
					InventoryOptionalFieldEncryptionStatus,
				},
			},
		},
	}
	putResult, err = client.PutBucketInventory(context.TODO(), putRequest)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
	time.Sleep(1 * time.Second)

	getRequest = &GetBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
	}
	getResult, err = client.GetBucketInventory(context.TODO(), getRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	listRequest = &ListBucketInventoryRequest{
		Bucket: Ptr(bucketNameNotExist),
	}
	listResult, err = client.ListBucketInventory(context.TODO(), listRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)

	delRequest = &DeleteBucketInventoryRequest{
		Bucket:      Ptr(bucketNameNotExist),
		InventoryId: Ptr(id),
	}
	delResult, err = client.DeleteBucketInventory(context.TODO(), delRequest)
	assert.NotNil(t, err)
	serr = &ServiceError{}
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
