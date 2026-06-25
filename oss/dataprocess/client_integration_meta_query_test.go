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

func TestMetaQuery(t *testing.T) {
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

	getResult, err := client.GetMetaQueryStatus(context.TODO(), &GetMetaQueryStatusRequest{
		Bucket: oss.Ptr(bucket_),
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, getResult.StatusCode)

	doResult, err := client.DoMetaQuery(context.TODO(), &DoMetaQueryRequest{
		Bucket: oss.Ptr(bucket_),
		Mode:   oss.Ptr("basic"),
		MetaQuery: &DoMetaQuery{
			Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
			Sort:  oss.Ptr("Size"),
			Order: oss.Ptr(MetaQueryOrderDesc),
			Aggregations: &MetaQueryAggregations{
				[]Aggregation{
					{
						Field:     oss.Ptr("Size"),
						Operation: oss.Ptr("sum"),
					},
				},
			},
			MaxResults: oss.Ptr(int64(100)),
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, doResult.StatusCode)

	invalidClient := getInvalidAkClient()
	_, err = invalidClient.OpenMetaQuery(context.TODO(), &OpenMetaQueryRequest{
		Bucket: oss.Ptr(bucket_),
		Mode:   oss.Ptr("basic"),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.GetMetaQueryStatus(context.TODO(), &GetMetaQueryStatusRequest{
		Bucket: oss.Ptr(bucket_),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.DoMetaQuery(context.TODO(), &DoMetaQueryRequest{
		Bucket: oss.Ptr(bucket_),
		Mode:   oss.Ptr("basic"),
		MetaQuery: &DoMetaQuery{
			Query: oss.Ptr(`{"Field":"Size","Operation":"gt","Value":"1048576"}`),
			Sort:  oss.Ptr("Size"),
			Order: oss.Ptr(MetaQueryOrderDesc),
			Aggregations: &MetaQueryAggregations{
				[]Aggregation{
					{
						Field:     oss.Ptr("Size"),
						Operation: oss.Ptr("sum"),
					},
				},
			},
			MaxResults: oss.Ptr(int64(100)),
		},
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)

	serr = &oss.ServiceError{}
	_, err = invalidClient.CloseMetaQuery(context.TODO(), &CloseMetaQueryRequest{
		Bucket: oss.Ptr(bucket_),
	})
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(403), serr.StatusCode)
	assert.Equal(t, "InvalidAccessKeyId", serr.Code)
	assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
	assert.NotEmpty(t, serr.RequestID)
}
