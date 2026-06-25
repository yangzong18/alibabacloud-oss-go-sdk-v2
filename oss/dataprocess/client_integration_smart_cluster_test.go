//go:build integration

package dataprocess

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestSmartCluster(t *testing.T) {
	var serr *oss.ServiceError
	var err error
	client := getDefaultClient()

	for {
		_, err = client.OpenMetaQuery(context.TODO(), &OpenMetaQueryRequest{
			Bucket: oss.Ptr(bucket_),
			Mode:   oss.Ptr("basic"),
		})
		if err != nil {
			errors.As(err, &serr)
			if serr.Code == "MetaQueryNotReady" {
				time.Sleep(5 * time.Second)
			} else if serr.Code == "MetaQueryAlreadyExist" {
				break
			}
		} else {
			break
		}
	}
	datasetName := datasetNamePrefix + randStr(5)
	name := "cluster-" + randStr(5)

	_, err = client.CreateDataset(context.TODO(), &CreateDatasetRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
	})
	assert.Nil(t, err)

	createResult, err := client.CreateSmartCluster(context.TODO(), &CreateSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		Name:        oss.Ptr(name),
		ClusterType: SmartClusterTypeKnowledge,
		Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["character","car"]}]`),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, createResult.StatusCode)
	defer func() {
		_, err = client.DeleteSmartCluster(context.TODO(), &DeleteSmartClusterRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(datasetName),
			ObjectId:    createResult.ObjectId,
		})
		assert.Nil(t, err)
		time.Sleep(5 * time.Second)
		for {
			_, err = client.GetSmartCluster(context.TODO(), &GetSmartClusterRequest{
				Bucket:      oss.Ptr(bucket_),
				DatasetName: oss.Ptr(datasetName),
				ObjectId:    createResult.ObjectId,
			})
			if err != nil {
				errors.As(err, &serr)
				if serr.StatusCode == 404 {
					break
				} else {
					time.Sleep(5 * time.Second)
				}
			} else {
				break
			}
		}
		_, err = client.DeleteDataset(context.TODO(), &DeleteDatasetRequest{
			Bucket:      oss.Ptr(bucket_),
			DatasetName: oss.Ptr(datasetName),
		})
		assert.Nil(t, err)
	}()

	getResult, err := client.GetSmartCluster(context.TODO(), &GetSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    createResult.ObjectId,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)

	updateResult, err := client.UpdateSmartCluster(context.TODO(), &UpdateSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    createResult.ObjectId,
		Description: oss.Ptr("this is a demo"),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, updateResult.StatusCode)

	listResult, err := client.ListSmartClusters(context.TODO(), &ListSmartClustersRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, listResult.StatusCode)
	assert.True(t, len(listResult.SmartClusters) > 0)

	invalidClient := getInvalidAkClient()
	_, err = invalidClient.CreateSmartCluster(context.TODO(), &CreateSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		Name:        oss.Ptr(name),
		ClusterType: SmartClusterTypeKnowledge,
		Rules:       oss.Ptr(`[{"RuleType": "keywords","Keywords": ["character","car"]}]`),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.GetSmartCluster(context.TODO(), &GetSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    createResult.ObjectId,
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.UpdateSmartCluster(context.TODO(), &UpdateSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    createResult.ObjectId,
		Description: oss.Ptr("this is a demo"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.ListSmartClusters(context.TODO(), &ListSmartClustersRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.DeleteSmartCluster(context.TODO(), &DeleteSmartClusterRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr(datasetName),
		ObjectId:    createResult.ObjectId,
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
