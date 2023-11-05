package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalOutput_encodetype(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
			<ListBucketResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
			<Name>oss-example</Name>
			<Prefix>hello%20world%21</Prefix>
			<Marker>hello%20</Marker>
			<MaxKeys>100</MaxKeys>
			<Delimiter>hello%20%21world</Delimiter>
			<IsTruncated>false</IsTruncated>
			<EncodingType>url</EncodingType>
			<Contents>
				<Key>fun%2Fmovie%2F001.avi</Key>
				<LastModified>2012-02-24T08:43:07.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
			<Contents>
				<Key>fun%2Fmovie%2F007.avi</Key>
				<LastModified>2012-02-24T08:43:27.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
			<Contents>
				<Key>oss.jpg</Key>
				<LastModified>2012-02-24T06:07:48.000Z</LastModified>
				<ETag>&quot;5B3C1A2E053D763E1B002CC607C5A0FE&quot;</ETag>
				<Type>Normal</Type>
				<Size>344606</Size>
				<StorageClass>Standard</StorageClass>
				<Owner>
					<ID>00220120222</ID>
					<DisplayName>user-example</DisplayName>
				</Owner>
			</Contents>
		</ListBucketResult>`

	// unsupport content-type
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, "hello ", *result.Marker)
	assert.Equal(t, "hello world!", *result.Prefix)
	assert.Equal(t, "hello !world", *result.Delimiter)
	assert.Equal(t, "hello !world", *result.Delimiter)
	assert.Equal(t, "url", *result.EncodingType)
	assert.Equal(t, "oss-example", *result.Name)
	assert.Equal(t, false, result.IsTruncated)
	assert.Nil(t, result.NextMarker)
	assert.Len(t, result.Contents, 3)
	assert.Equal(t, "fun/movie/001.avi", result.Contents[0].Key)
	assert.Equal(t, "\"5B3C1A2E053D763E1B002CC607C5A0FE\"", result.Contents[0].ETag)
	assert.Equal(t, "fun/movie/007.avi", result.Contents[1].Key)
	assert.Equal(t, "oss.jpg", result.Contents[2].Key)
	assert.Len(t, result.CommonPrefixes, 0)
}

func TestUnmarshalOutput_encodetype1(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	body := `<?xml version="1.0" encoding="UTF-8"?>
		<ListBucketResult xmlns="http://doc.oss-cn-hangzhou.aliyuncs.com">
			<Contents>
				<LastModified>2012-02-24T08:43:07.000Z</LastModified>
			</Contents>
		</ListBucketResult>`

	// unsupport content-type
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
}
