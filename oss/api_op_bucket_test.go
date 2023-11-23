package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

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
	assert.Equal(t, "fun/movie/001.avi", *result.Contents[0].Key)
	assert.Equal(t, "\"5B3C1A2E053D763E1B002CC607C5A0FE\"", *result.Contents[0].ETag)
	assert.Equal(t, "fun/movie/007.avi", *result.Contents[1].Key)
	assert.Equal(t, "oss.jpg", *result.Contents[2].Key)
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

func TestMarshalInput_PutBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketRequest
	var input *OperationInput
	var err error

	request = &PutBucketRequest{}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_PutBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "BucketAlreadyExist",
		Headers: http.Header{
			"Content-Type":     {"application/xml"},
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &PutBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "BucketAlreadyExist")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketRequest{}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteBucketRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucket",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucket(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	output = &OperationOutput{
		StatusCode: 209,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"5C3D9778CC1C2AEDF85B****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &DeleteBucketResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 209)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C3D9778CC1C2AEDF85B****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListObjectsRequest
	var input *OperationInput
	var err error

	request = &ListObjectsRequest{}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListObjectsRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListObjectsRequest{
		Bucket:       Ptr("oss-demo"),
		Delimiter:    Ptr("/"),
		Marker:       Ptr(""),
		MaxKeys:      int32(10),
		Prefix:       Ptr(""),
		EncodingType: Ptr("URL"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
<Prefix></Prefix>
<Marker></Marker>
<MaxKeys>100</MaxKeys>
<Delimiter></Delimiter>
<IsTruncated>false</IsTruncated>
<Contents>
      <Key>fun/movie/001.avi</Key>
      <LastModified>2012-02-24T08:43:07.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>fun/movie/007.avi</Key>
      <LastModified>2012-02-24T08:43:27.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>fun/test.jpg</Key>
      <LastModified>2012-02-24T08:42:32.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
<Contents>
      <Key>oss.jpg</Key>
      <LastModified>2012-02-24T06:07:48.000Z</LastModified>
      <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
      <Type>Normal</Type>
      <Size>344606</Size>
      <StorageClass>Standard</StorageClass>
      <Owner>
          <ID>0022012****</ID>
          <DisplayName>user-example</DisplayName>
      </Owner>
</Contents>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Empty(t, result.Prefix)
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Marker)
	assert.Empty(t, result.Delimiter)
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, len(result.Contents), 4)
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2012, time.February, 24, 8, 43, 27, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
	assert.Equal(t, *result.Contents[3].Type, "Normal")
	assert.Equal(t, result.Contents[0].Size, int64(344606))
	assert.Equal(t, *result.Contents[1].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[2].Owner.ID, "0022012****")
	assert.Equal(t, *result.Contents[3].Owner.DisplayName, "user-example")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
  <Prefix>fun</Prefix>
  <Marker>test1.txt</Marker>
  <MaxKeys>100</MaxKeys>
  <Delimiter>/</Delimiter>
  <IsTruncated>true</IsTruncated>
  <Contents>
        <Key>exampleobject1.txt</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>ColdArchive</StorageClass>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
  <Contents>
        <Key>exampleobject2.txt</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <RestoreInfo>ongoing-request="true"</RestoreInfo>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
  <Contents>
        <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <RestoreInfo>ongoing-request="false", expiry-date="Thu, 24 Sep 2020 12:40:33 GMT"</RestoreInfo>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "fun")
	assert.Equal(t, *result.Marker, "test1.txt")
	assert.Equal(t, *result.Delimiter, "/")
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, len(result.Contents), 3)
	assert.Equal(t, *result.Contents[0].Key, "exampleobject1.txt")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
	assert.Equal(t, *result.Contents[0].Type, "Normal")
	assert.Equal(t, result.Contents[1].Size, int64(344606))
	assert.Equal(t, *result.Contents[2].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[0].Owner.ID, "0022012****")
	assert.Equal(t, *result.Contents[0].Owner.DisplayName, "user-example")
	assert.Empty(t, result.Contents[0].RestoreInfo)
	assert.Equal(t, *result.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
	assert.Equal(t, *result.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "Conflict",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "Conflict")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &ListObjectsResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListObjectsV2(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListObjectsRequestV2
	var input *OperationInput
	var err error

	request = &ListObjectsRequestV2{}
	input = &OperationInput{
		OpName: "ListObjectsV2",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &ListObjectsRequestV2{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "ListObjects",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	request = &ListObjectsRequestV2{
		Bucket:       Ptr("oss-demo"),
		Delimiter:    Ptr("/"),
		StartAfter:   Ptr(""),
		MaxKeys:      int32(10),
		Prefix:       Ptr(""),
		EncodingType: Ptr("URL"),
		FetchOwner:   true,
	}
	input = &OperationInput{
		OpName: "ListObjectsV2",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListObjectsV2(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix></Prefix>
    <MaxKeys>100</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a</Key>
        <LastModified>2020-05-18T05:45:43.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>b</Key>
        <LastModified>2020-05-18T05:45:50.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>b/c</Key>
        <LastModified>2020-05-18T05:45:54.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>bc</Key>
        <LastModified>2020-05-18T05:45:59.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <Contents>
        <Key>c</Key>
        <LastModified>2020-05-18T05:45:57.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>Standard</StorageClass>
    </Contents>
    <KeyCount>6</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Prefix)
	assert.Empty(t, result.Delimiter)
	assert.Empty(t, result.StartAfter)
	assert.Empty(t, result.ContinuationToken)
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 6)
	assert.Equal(t, len(result.Contents), 6)
	assert.Equal(t, *result.Contents[0].Key, "a")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[2].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[3].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "Standard")
	assert.Equal(t, *result.Contents[1].StorageClass, "Standard")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a</Prefix>
    <MaxKeys>100</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a</Key>
        <LastModified>2020-05-18T05:45:43.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <KeyCount>2</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "a")
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, len(result.Contents), 2)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 2)
	assert.Equal(t, *result.Contents[0].Key, "a")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 43, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")

	assert.Equal(t, *result.Contents[1].Key, "a/b")
	assert.Equal(t, *result.Contents[1].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[1].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[1].Size, int64(25))
	assert.Equal(t, *result.Contents[1].StorageClass, "STANDARD")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a/</Prefix>
    <MaxKeys>100</MaxKeys>
    <Delimiter>/</Delimiter>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <CommonPrefixes>
        <Prefix>a/b/</Prefix>
    </CommonPrefixes>
    <KeyCount>2</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Equal(t, *result.Prefix, "a/")
	assert.Equal(t, result.MaxKeys, int32(100))
	assert.Equal(t, len(result.Contents), 1)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, false)
	assert.Equal(t, result.KeyCount, 2)
	assert.Equal(t, *result.Contents[0].Key, "a/b")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 47, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")

	assert.Equal(t, *result.CommonPrefixes[0].Prefix, "a/b/")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix></Prefix>
    <StartAfter>b</StartAfter>
    <MaxKeys>3</MaxKeys>
    <EncodingType>url</EncodingType>
    <IsTruncated>true</IsTruncated>
    <NextContinuationToken>CgJiYw--</NextContinuationToken>
    <Contents>
        <Key>b%2Fc</Key>
        <LastModified>2020-05-18T05:45:54.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <Contents>
        <Key>ba</Key>
        <LastModified>2020-05-18T11:17:58.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <Contents>
        <Key>bc</Key>
        <LastModified>2020-05-18T05:45:59.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1686240967192623</ID>
            <DisplayName>1686240967192623</DisplayName>
        </Owner>
    </Contents>
    <KeyCount>3</KeyCount>
</ListBucketResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListObjectsResultV2{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.Name, "examplebucket")
	assert.Empty(t, result.Prefix)
	assert.Equal(t, *result.StartAfter, "b")
	assert.Equal(t, result.MaxKeys, int32(3))
	assert.Equal(t, len(result.Contents), 3)
	assert.Equal(t, *result.EncodingType, "url")
	assert.Equal(t, result.IsTruncated, true)
	assert.Equal(t, *result.NextContinuationToken, "CgJiYw--")
	assert.Equal(t, result.KeyCount, 3)
	assert.Equal(t, *result.Contents[0].Key, "b/c")
	assert.Equal(t, *result.Contents[0].LastModified, time.Date(2020, time.May, 18, 5, 45, 54, 0, time.UTC))
	assert.Equal(t, *result.Contents[0].ETag, "\"35A27C2B9EAEEB6F48FD7FB5861D****\"")
	assert.Equal(t, result.Contents[0].Size, int64(25))
	assert.Equal(t, *result.Contents[0].StorageClass, "STANDARD")
	assert.Equal(t, *result.Contents[0].Owner.DisplayName, "1686240967192623")
	assert.Equal(t, *result.Contents[0].Owner.ID, "1686240967192623")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &ListObjectsResultV2{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketInfoRequest
	var input *OperationInput
	var err error

	request = &GetBucketInfoRequest{}
	input = &OperationInput{
		OpName: "GetBucketInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketInfoRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketInfo",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"bucketInfo": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketInfo(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")

	assert.Empty(t, result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Nil(t, result.BucketInfo.SseRule.SSEAlgorithm)
	assert.Nil(t, result.BucketInfo.SseRule.KMSDataEncryption)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
	<ServerSideEncryptionRule>
		<SSEAlgorithm>None</SSEAlgorithm>
	</ServerSideEncryptionRule>
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	if result.BucketInfo.SseRule.KMSMasterKeyID != nil && *result.BucketInfo.SseRule.KMSMasterKeyID == "None" {
		*result.BucketInfo.SseRule.KMSMasterKeyID = ""
	}
	if result.BucketInfo.SseRule.SSEAlgorithm != nil && *result.BucketInfo.SseRule.SSEAlgorithm == "None" {
		*result.BucketInfo.SseRule.SSEAlgorithm = ""
	}
	if result.BucketInfo.SseRule.KMSDataEncryption != nil && *result.BucketInfo.SseRule.KMSDataEncryption == "None" {
		*result.BucketInfo.SseRule.KMSDataEncryption = ""
	}
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")
	assert.Empty(t, result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Equal(t, *result.BucketInfo.SseRule.SSEAlgorithm, "")
	assert.Nil(t, result.BucketInfo.SseRule.KMSDataEncryption)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<BucketInfo>
  <Bucket>
    <AccessMonitor>Enabled</AccessMonitor>
    <CreationDate>2013-07-31T10:56:21.000Z</CreationDate>
    <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
    <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
    <Location>oss-cn-hangzhou</Location>
    <StorageClass>Standard</StorageClass>
    <TransferAcceleration>Disabled</TransferAcceleration>
    <CrossRegionReplication>Disabled</CrossRegionReplication>
    <Name>oss-example</Name>
    <ResourceGroupId>rg-aek27tc********</ResourceGroupId>
    <Owner>
      <DisplayName>username</DisplayName>
      <ID>27183473914****</ID>
    </Owner>
    <AccessControlList>
      <Grant>private</Grant>
    </AccessControlList>  
	<ServerSideEncryptionRule>
		<SSEAlgorithm>KMS</SSEAlgorithm>
		<KMSMasterKeyID></KMSMasterKeyID>
		<KMSDataEncryption>SM4</KMSDataEncryption>
	</ServerSideEncryptionRule>
    <BucketPolicy>
      <LogBucket>examplebucket</LogBucket>
      <LogPrefix>log/</LogPrefix>
    </BucketPolicy>
  </Bucket>
</BucketInfo>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInfoResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.BucketInfo.Name, "oss-example")
	assert.Equal(t, *result.BucketInfo.AccessMonitor, "Enabled")
	assert.Equal(t, *result.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
	assert.Equal(t, *result.BucketInfo.Location, "oss-cn-hangzhou")
	assert.Equal(t, *result.BucketInfo.StorageClass, "Standard")
	assert.Equal(t, *result.BucketInfo.TransferAcceleration, "Disabled")
	assert.Equal(t, *result.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
	assert.Equal(t, *result.BucketInfo.CrossRegionReplication, "Disabled")
	assert.Equal(t, *result.BucketInfo.ResourceGroupId, "rg-aek27tc********")
	assert.Equal(t, *result.BucketInfo.Owner.ID, "27183473914****")
	assert.Equal(t, *result.BucketInfo.Owner.DisplayName, "username")
	assert.Equal(t, *result.BucketInfo.ACL, "private")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
	assert.Equal(t, *result.BucketInfo.BucketPolicy.LogPrefix, "log/")
	assert.Empty(t, *result.BucketInfo.SseRule.KMSMasterKeyID)
	assert.Equal(t, *result.BucketInfo.SseRule.SSEAlgorithm, "KMS")
	assert.Equal(t, *result.BucketInfo.SseRule.KMSDataEncryption, "SM4")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",

		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &GetBucketInfoResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketLocation(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketLocationRequest
	var input *OperationInput
	var err error

	request = &GetBucketLocationRequest{}
	input = &OperationInput{
		OpName: "GetBucketLocation",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"location": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketLocationRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketLocation",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"location": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketLocation(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-hangzhou</LocationConstraint>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketLocationResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, *result.LocationConstraint, "oss-cn-hangzhou")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>534B371674E88A4D8906****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &GetBucketLocationResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")

}

func TestMarshalInput_GetBucketStat(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketStatRequest
	var input *OperationInput
	var err error

	request = &GetBucketStatRequest{}
	input = &OperationInput{
		OpName: "GetBucketStat",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"stat": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketStatRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketStat",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"stat": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketStat(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<BucketStat>
  <Storage>1600</Storage>
  <ObjectCount>230</ObjectCount>
  <MultipartUploadCount>40</MultipartUploadCount>
  <LiveChannelCount>4</LiveChannelCount>
  <LastModifiedTime>1643341269</LastModifiedTime>
  <StandardStorage>430</StandardStorage>
  <StandardObjectCount>66</StandardObjectCount>
  <InfrequentAccessStorage>2359296</InfrequentAccessStorage>
  <InfrequentAccessRealStorage>360</InfrequentAccessRealStorage>
  <InfrequentAccessObjectCount>54</InfrequentAccessObjectCount>
  <ArchiveStorage>2949120</ArchiveStorage>
  <ArchiveRealStorage>450</ArchiveRealStorage>
  <ArchiveObjectCount>74</ArchiveObjectCount>
  <ColdArchiveStorage>2359296</ColdArchiveStorage>
  <ColdArchiveRealStorage>360</ColdArchiveRealStorage>
  <ColdArchiveObjectCount>36</ColdArchiveObjectCount>
</BucketStat>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketStatResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Equal(t, int64(1600), result.Storage)
	assert.Equal(t, int64(230), result.ObjectCount)
	assert.Equal(t, int64(40), result.MultipartUploadCount)
	assert.Equal(t, int64(4), result.LiveChannelCount)
	assert.Equal(t, int64(1643341269), result.LastModifiedTime)
	assert.Equal(t, int64(430), result.StandardStorage)
	assert.Equal(t, int64(66), result.StandardObjectCount)
	assert.Equal(t, int64(2359296), result.InfrequentAccessStorage)
	assert.Equal(t, int64(360), result.InfrequentAccessRealStorage)
	assert.Equal(t, int64(54), result.InfrequentAccessObjectCount)
	assert.Equal(t, int64(2949120), result.ArchiveStorage)
	assert.Equal(t, int64(450), result.ArchiveRealStorage)
	assert.Equal(t, int64(74), result.ArchiveObjectCount)
	assert.Equal(t, int64(2359296), result.ColdArchiveStorage)
	assert.Equal(t, int64(360), result.ColdArchiveRealStorage)
	assert.Equal(t, int64(36), result.ColdArchiveObjectCount)

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>534B371674E88A4D8906****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &GetBucketStatResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")

}

func TestMarshalInput_PutBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketAclRequest
	var input *OperationInput
	var err error

	request = &PutBucketAclRequest{}
	input = &OperationInput{
		OpName: "PutBucketAcl",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutBucketAclRequest{
		Bucket: Ptr("oss-demo"),
		Acl:    BucketACLPrivate,
	}
	input = &OperationInput{
		OpName: "PutBucketAcl",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_PutBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &PutBucketAclResult{}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &PutBucketAclResult{}
	err = c.unmarshalOutput(resultErr, output, discardBody)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketAclRequest
	var input *OperationInput
	var err error

	request = &GetBucketAclRequest{}
	input = &OperationInput{
		OpName: "GetBucketAcl",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetBucketAclRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketAcl",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeDefault,
		},
		Parameters: map[string]string{
			"acl": "",
		},
		Bucket: request.Bucket,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketAcl(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result := &GetBucketAclResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	assert.Equal(t, *result.ACL, "public-read")
	assert.Equal(t, *result.Owner.ID, "0022012****")
	assert.Equal(t, *result.Owner.DisplayName, "user_example")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>AccessDenied</Code>
  <Message>AccessDenied</Message>
  <RequestId>568D5566F2D0F89F5C0E****</RequestId>
  <HostId>test.oss.aliyuncs.com</HostId>
</Error>`
	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"568D5566F2D0F89F5C0E****"},
			"Content-Type":     {"application/xml"},
		},
	}
	resultErr := &PutBucketAclResult{}
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, resultErr.StatusCode, 403)
	assert.Equal(t, resultErr.Status, "AccessDenied")
	assert.Equal(t, resultErr.Headers.Get("X-Oss-Request-Id"), "568D5566F2D0F89F5C0E****")
	assert.Equal(t, resultErr.Headers.Get("Content-Type"), "application/xml")
}
