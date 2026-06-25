//go:build integration

package dataprocess

import (
	"context"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func TestSimpleQueryBasic(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr(queryJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	f := result.Files[0]
	assert.NotNil(t, f.URI, "URI should not be nil")
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.MediaType, "MediaType should not be nil")
	assert.NotNil(t, f.ContentType, "ContentType should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
	assert.True(t, *f.Size > 0, "Size should be > 0")

	// Verify labels parsing
	anyFileHasLabels := false
	for _, file := range result.Files {
		if oss.ToInt64(file.OSSTaggingCount) > 0 {
			assert.Equal(t, *file.OSSTaggingCount, (int64)(len(file.OSSTagging)))
		}

		if len(file.Labels) > 0 {
			anyFileHasLabels = true
			assert.NotNil(t, file.Labels[0].LabelName, "Label.LabelName should not be nil")
			break
		}
	}
	assert.True(t, anyFileHasLabels, "At least one file should have labels")
}

func TestSimpleQueryWithAggregations(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	aggregationsJSON := `[{"Field":"Size","Operation":"sum"}]`

	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:       oss.Ptr(bucket_),
		DatasetName:  oss.Ptr("test-dataset"),
		Query:        oss.Ptr(queryJSON),
		Aggregations: oss.Ptr(aggregationsJSON),
		MaxResults:   oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	// Verify aggregations
	assert.NotNil(t, result.Aggregations, "aggregations should not be nil")
	assert.True(t, len(result.Aggregations) > 0, "aggregations should not be empty")

	agg := result.Aggregations[0]
	assert.Equal(t, "Size", *agg.Field)
	assert.Equal(t, "sum", *agg.Operation)
	assert.NotNil(t, agg.Value, "aggregation value should not be nil")
}

func TestSimpleQueryWithSortAndOrder(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:           oss.Ptr(bucket_),
		DatasetName:      oss.Ptr("test-dataset"),
		Query:            oss.Ptr(queryJSON),
		Sort:             oss.Ptr("Filename"),
		Order:            oss.Ptr("asc"),
		MaxResults:       oss.Ptr(int32(10)),
		WithoutTotalHits: oss.Ptr(false),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify sorted by Filename ascending
	files := result.Files
	for i := 1; i < len(files); i++ {
		prev := *files[i-1].Filename
		curr := *files[i].Filename
		assert.True(t, prev <= curr,
			"files should be sorted by Filename asc: %s <= %s", prev, curr)
	}
}

func TestSimpleQueryWithFields(t *testing.T) {
	client := getDefaultClient()

	queryJSON := `{"Field":"Filename","Value":"test-media","Operation":"prefix"}`
	withFieldsJSON := `["Filename","Size","ContentType"]`

	result, err := client.SimpleQuery(context.TODO(), &SimpleQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr(queryJSON),
		WithFields:  oss.Ptr(withFieldsJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify requested fields are populated
	f := result.Files[0]
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
	assert.NotNil(t, f.ContentType, "ContentType should not be nil")
}

func TestSemanticQueryBasic(t *testing.T) {
	client := getDefaultClient()

	result, err := client.SemanticQuery(context.TODO(), &SemanticQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr("雪景"),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	f := result.Files[0]
	assert.NotNil(t, f.URI, "URI should not be nil")
	assert.NotNil(t, f.Filename, "Filename should not be nil")
	assert.NotNil(t, f.MediaType, "MediaType should not be nil")
	assert.NotNil(t, f.Size, "Size should not be nil")
}

func TestSemanticQueryWithMediaTypes(t *testing.T) {
	client := getDefaultClient()

	mediaTypesJSON := `["image"]`
	withFieldsJSON := `["Filename","Size","MediaType"]`

	result, err := client.SemanticQuery(context.TODO(), &SemanticQueryRequest{
		Bucket:      oss.Ptr(bucket_),
		DatasetName: oss.Ptr("test-dataset"),
		Query:       oss.Ptr("雪景"),
		MediaTypes:  oss.Ptr(mediaTypesJSON),
		WithFields:  oss.Ptr(withFieldsJSON),
		MaxResults:  oss.Ptr(int32(10)),
	})
	dumpErrIfNotNil(err)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)

	assert.NotNil(t, result.Files, "files should not be nil")
	assert.True(t, len(result.Files) > 0, "files should not be empty")

	// Verify all returned files are images
	for _, f := range result.Files {
		assert.Equal(t, "image", *f.MediaType, "MediaType should be image")
	}
}
