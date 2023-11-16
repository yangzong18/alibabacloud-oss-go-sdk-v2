package oss

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMarshalInput_PutObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutObjectRequest
	var input *OperationInput
	var err error

	request = &PutObjectRequest{}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
		Body:   request.Body,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &PutObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)

	request = &PutObjectRequest{
		Bucket:               Ptr("oss-bucket"),
		Key:                  Ptr("oss-key"),
		CacheControl:         Ptr("no-cache"),
		ContentDisposition:   Ptr("attachment"),
		ContentEncoding:      Ptr("utf-8"),
		ContentMD5:           Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:        Ptr(int64(100)),
		Expires:              Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:      Ptr("true"),
		ServerSideEncryption: Ptr("AES256"),
		Acl:                  ObjectACLPrivate,
		StorageClass:         StorageClassStandard,
		Metadata: map[string]string{
			"location": "demo",
			"user":     "walker",
		},
		Tagging: Ptr("TagA=A&TagB=B"),
	}
	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-user"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-location"], "demo")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "AES256")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "true")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=A&TagB=B")
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)

	body := randLowStr(1000)
	request = &PutObjectRequest{
		Bucket:                   Ptr("oss-bucket"),
		Key:                      Ptr("oss-key"),
		CacheControl:             Ptr("no-cache"),
		ContentDisposition:       Ptr("attachment"),
		ContentEncoding:          Ptr("utf-8"),
		ContentMD5:               Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:            Ptr(int64(100)),
		Expires:                  Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:          Ptr("false"),
		ServerSideEncryption:     Ptr("KMS"),
		ServerSideDataEncryption: Ptr("SM4"),
		SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		Acl:                      ObjectACLPrivate,
		StorageClass:             StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(body),
		},
	}

	input = &OperationInput{
		OpName: "PutObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Nil(t, input.Parameters)
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_PutObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
		},
	}
	result := &PutObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Nil(t, result.VersionId)

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"5C06A3B67B8B5A3DA422****"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma": {"316181249502703****"},
			"Content-MD5":          {"1B2M2Y8AsgTpgAmY7PhC****"},
			"x-oss-version-id":     {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
		},
	}
	result = &PutObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"A797938C31D59EDD08D86188F6D5****\"")
	assert.Equal(t, *result.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
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
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetObjectRequest
	var input *OperationInput
	var err error

	request = &GetObjectRequest{}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &GetObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")

	request = &GetObjectRequest{
		Bucket:                     Ptr("oss-bucket"),
		Key:                        Ptr("oss-key"),
		IfMatch:                    Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:                Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:            Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince:          Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		Range:                      Ptr("bytes 0~9/44"),
		ResponseCacheControl:       Ptr("gzip"),
		ResponseContentDisposition: Ptr("attachment; filename=testing.txt"),
		ResponseContentEncoding:    Ptr("utf-8"),
		ResponseContentLanguage:    Ptr("中文"),
		ResponseContentType:        Ptr("text"),
		ResponseExpires:            Ptr("Fri, 24 Feb 2012 17:00:00 GMT"),
		VersionId:                  Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")

	assert.Equal(t, input.Headers["If-Match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["If-None-Match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["If-Modified-Since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["If-Unmodified-Since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Headers["Range"], "bytes 0~9/44")
	assert.Equal(t, input.Parameters["response-cache-control"], "gzip")
	assert.Equal(t, input.Parameters["response-content-disposition"], "attachment; filename=testing.txt")
	assert.Equal(t, input.Parameters["response-content-encoding"], "utf-8")
	assert.Equal(t, input.Parameters["response-content-language"], "中文")
	assert.Equal(t, input.Parameters["response-expires"], "Fri, 24 Feb 2012 17:00:00 GMT")
	assert.Equal(t, input.Parameters["versionId"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_GetObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":      {"image/jpg"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"Content-Length":    {"344606"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))

	body = randLowStr(34460)
	output = &OperationOutput{
		StatusCode: 206,
		Status:     "Partial Content",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":      {"image/jpg"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":    {"801"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
			"Accept-Ranges":     {"bytes"},
			"Content-Range":     {"bytes 100-900/34460"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 206)
	assert.Equal(t, result.Status, "Partial Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "image/jpg")
	assert.Equal(t, result.ContentLength, int64(801))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.ContentRange, "bytes 100-900/34460")

	body = randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":                        {"text"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":                      {"344606"},
			"Last-Modified":                       {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type":                   {"Normal"},
			"Accept-Ranges":                       {"bytes"},
			"Content-disposition":                 {"attachment; filename=testing.txt"},
			"Cache-control":                       {"no-cache"},
			"X-Oss-Storage-Class":                 {"Standard"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
			"x-oss-tagging-count":                 {"2"},
			"Content-MD5":                         {"si4Nw3Cn9wZ/rPX3XX+j****"},
			"x-oss-hash-crc64ecma":                {"870718044876840****"},
			"x-oss-meta-name":                     {"demo"},
			"x-oss-meta-email":                    {"demo@aliyun.com"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "text")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, *result.StorageClass, "Standard")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
	assert.Equal(t, result.TaggingCount, int32(2))
	assert.Equal(t, result.Metadata["name"], "demo")
	assert.Equal(t, result.Metadata["email"], "demo@aliyun.com")
	assert.Equal(t, *result.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
	assert.Equal(t, *result.HashCRC64, "870718044876840****")
	body = randLowStr(344606)
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":  {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":      {"text"},
			"Date":              {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":              {"\"5B3C1A2E05E1B002CC607C****\""},
			"Content-Length":    {"344606"},
			"Last-Modified":     {"Fri, 24 Feb 2012 06:07:48 GMT"},
			"x-oss-object-type": {"Normal"},
			"x-oss-restore":     {"ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\""},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &GetObjectResult{
		Body: output.Body,
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Cache-control"), "")
	assert.Equal(t, *result.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
	assert.Equal(t, *result.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
	assert.Equal(t, *result.ContentType, "text")
	assert.Equal(t, result.ContentLength, int64(344606))
	assert.Equal(t, *result.ObjectType, "Normal")
	assert.Equal(t, result.Body, io.NopCloser(bytes.NewReader([]byte(body))))
	assert.Equal(t, *result.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalHeader)
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
	err = c.unmarshalOutput(result, output, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_CopyObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CopyObjectRequest
	var input *OperationInput
	var err error

	request = &CopyObjectRequest{}
	input = &OperationInput{
		OpName: "CopyObject",
		Method: "PUT",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &CopyObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-copy-key"),
		Source: Ptr("/oss-bucket/oss-key"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-key")

	request = &CopyObjectRequest{
		Bucket:            Ptr("oss-bucket"),
		Key:               Ptr("oss-copy-key"),
		Source:            Ptr("/oss-bucket/oss-key"),
		IfMatch:           Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:       Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:   Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince: Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		VersionId:         Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-none-match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-modified-since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-unmodified-since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Parameters["x-oss-copy-source-version-id"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-key")
	assert.Nil(t, input.OpMetadata.values)

	request = &CopyObjectRequest{
		Bucket:                   Ptr("oss-copy-bucket"),
		Key:                      Ptr("oss-copy-key"),
		Source:                   Ptr("/oss-bucket/oss-key"),
		IfMatch:                  Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
		IfNoneMatch:              Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
		IfModifiedSince:          Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
		IfUnmodifiedSince:        Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
		VersionId:                Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****"),
		ForbidOverwrite:          Ptr("false"),
		ServerSideEncryption:     Ptr("KMS"),
		ServerSideDataEncryption: Ptr("SM4"),
		SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		MetadataDirective:        Ptr("REPLACE"),
		TaggingDirective:         Ptr("Replace"),
		Acl:                      ObjectACLPrivate,
		StorageClass:             StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
	}
	input = &OperationInput{
		OpName: "GetObject",
		Method: "GET",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-copy-bucket")
	assert.Equal(t, *input.Key, "oss-copy-key")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-match"], "\"D41D8CD98F00B204E9800998ECF8****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-none-match"], "\"D41D8CD98F00B204E9800998ECF9****\"")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-modified-since"], "Fri, 13 Nov 2023 14:47:53 GMT")
	assert.Equal(t, input.Headers["x-oss-copy-source-if-unmodified-since"], "Fri, 13 Nov 2015 14:47:53 GMT")
	assert.Equal(t, input.Parameters["x-oss-copy-source-version-id"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY*****")
	assert.Equal(t, input.Headers["x-oss-copy-source"], "/oss-bucket/oss-key")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Equal(t, input.Headers["x-oss-tagging-directive"], "Replace")
	assert.Equal(t, input.Headers["x-oss-metadata-directive"], "REPLACE")
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_CopyObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
<CopyObjectResult>
  <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
  <LastModified>2018-02-24T09:41:56.000Z</LastModified>
</CopyObjectResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":     {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":         {"image/jpg"},
			"Date":                 {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                 {"\"F2064A169EE92E9775EE5324D0B1****\""},
			"Content-Length":       {"344606"},
			"x-oss-hash-crc64ecma": {"1275300285919610****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result := &CopyObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
	assert.Equal(t, *result.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.HashCRC64, "1275300285919610****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"28f6-15ea-8224-234e-c0ce407****"},
			"Content-Type":                        {"text"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"F2064A169EE92E9775EE5324D0B1****\""},
			"Content-Length":                      {"344606"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
			"x-oss-hash-crc64ecma":                {"870718044876840****"},
			"x-oss-copy-source-version-id":        {"CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****"},
			"x-oss-version-id":                    {"CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &CopyObjectResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
	assert.Equal(t, *result.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
	assert.Equal(t, *result.HashCRC64, "870718044876840****")
	assert.Equal(t, *result.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
	assert.Equal(t, *result.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
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
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader)
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
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	var resultErr PutBucketAclResult
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml, unmarshalHeader)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "deserialization failed, non-pointer passed to Unmarshal")
}

func TestMarshalInput_AppendObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *AppendObjectRequest
	var input *OperationInput
	var err error

	request = &AppendObjectRequest{}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &AppendObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")
	p := int64(0)
	request = &AppendObjectRequest{
		Bucket:               Ptr("oss-bucket"),
		Key:                  Ptr("oss-key"),
		Position:             Ptr(p),
		CacheControl:         Ptr("no-cache"),
		ContentDisposition:   Ptr("attachment"),
		ContentEncoding:      Ptr("gzip"),
		ContentMD5:           Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:        Ptr(int64(100)),
		Expires:              Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:      Ptr("true"),
		ServerSideEncryption: Ptr("AES256"),
		Acl:                  ObjectACLPrivate,
		StorageClass:         StorageClassStandard,
		Metadata: map[string]string{
			"location": "demo",
			"user":     "walker",
		},
		Tagging: Ptr("TagA=A&TagB=B"),
	}
	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.Body)
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-user"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-location"], "demo")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "AES256")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "true")
	assert.Equal(t, input.Headers["Content-Encoding"], "gzip")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=A&TagB=B")
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)

	body := randLowStr(1000)
	request = &AppendObjectRequest{
		Bucket:                   Ptr("oss-bucket"),
		Key:                      Ptr("oss-key"),
		Position:                 Ptr(int64(0)),
		CacheControl:             Ptr("no-cache"),
		ContentDisposition:       Ptr("attachment"),
		ContentEncoding:          Ptr("utf-8"),
		ContentMD5:               Ptr("eB5eJF1ptWaXm4bijSPyxw=="),
		ContentLength:            Ptr(int64(100)),
		Expires:                  Ptr("2022-10-12T00:00:00.000Z"),
		ForbidOverwrite:          Ptr("false"),
		ServerSideEncryption:     Ptr("KMS"),
		ServerSideDataEncryption: Ptr("SM4"),
		SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
		Acl:                      ObjectACLPrivate,
		StorageClass:             StorageClassStandard,
		Metadata: map[string]string{
			"name":  "walker",
			"email": "demo@aliyun.com",
		},
		Tagging: Ptr("TagA=B&TagC=D"),
		RequestCommon: RequestCommon{
			Body: strings.NewReader(body),
		},
	}

	input = &OperationInput{
		OpName:     "AppendObject",
		Method:     "POST",
		Parameters: map[string]string{"append": ""},
		Bucket:     request.Bucket,
		Key:        request.Key,
	}
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)

	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Body, strings.NewReader(body))
	assert.Equal(t, input.Headers["Cache-Control"], "no-cache")
	assert.Equal(t, input.Headers["Content-Disposition"], "attachment")
	assert.Equal(t, input.Headers["x-oss-meta-name"], "walker")
	assert.Equal(t, input.Headers["x-oss-meta-email"], "demo@aliyun.com")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption"], "KMS")
	assert.Equal(t, input.Headers["x-oss-server-side-data-encryption"], "SM4")
	assert.Equal(t, input.Headers["x-oss-server-side-encryption-key-id"], "9468da86-3509-4f8d-a61e-6eab1eac****")
	assert.Equal(t, input.Headers["x-oss-storage-class"], string(StorageClassStandard))
	assert.Equal(t, input.Headers["x-oss-object-acl"], string(ObjectACLPrivate))
	assert.Equal(t, input.Headers["x-oss-forbid-overwrite"], "false")
	assert.Equal(t, input.Headers["Content-Encoding"], "utf-8")
	assert.Equal(t, input.Headers["Content-Length"], "100")
	assert.Equal(t, input.Headers["Content-MD5"], "eB5eJF1ptWaXm4bijSPyxw==")
	assert.Equal(t, input.Headers["Expires"], "2022-10-12T00:00:00.000Z")
	assert.Equal(t, input.Headers["x-oss-tagging"], "TagA=B&TagC=D")
	assert.Empty(t, input.Parameters["append"])
	assert.Equal(t, input.Parameters["position"], strconv.FormatInt(p, 10))
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_AppendObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"D41D8CD98F00B204E9800998ECF8****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-next-append-position": {"0"},
		},
	}
	result := &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, result.NextPosition, int64(0))
	assert.Nil(t, result.VersionId)

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-version-id":           {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position": {"1717"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))

	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":           {"5C06A3B67B8B5A3DA422****"},
			"Date":                       {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                       {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":       {"316181249502703****"},
			"x-oss-version-id":           {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position": {"1717"},
			"x-oss-meta-name":            {"demo"},
			"x-oss-meta-email":           {"demo@aliyun.com"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))
	assert.Equal(t, result.Metadata["name"], "demo")
	assert.Equal(t, result.Metadata["email"], "demo@aliyun.com")
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id":                    {"5C06A3B67B8B5A3DA422****"},
			"Date":                                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"ETag":                                {"\"A797938C31D59EDD08D86188F6D5****\""},
			"x-oss-hash-crc64ecma":                {"316181249502703****"},
			"x-oss-version-id":                    {"CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****"},
			"x-oss-next-append-position":          {"1717"},
			"x-oss-server-side-encryption":        {"KMS"},
			"x-oss-server-side-data-encryption":   {"SM4"},
			"x-oss-server-side-encryption-key-id": {"12f8711f-90df-4e0d-903d-ab972b0f****"},
		},
	}
	result = &AppendObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "5C06A3B67B8B5A3DA422****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.HashCRC64, "316181249502703****")
	assert.Equal(t, *result.VersionId, "CAEQNhiBgMDJgZCA0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY0****")
	assert.Equal(t, result.NextPosition, int64(1717))
	assert.Equal(t, *result.ServerSideDataEncryption, "SM4")
	assert.Equal(t, *result.ServerSideEncryption, "KMS")
	assert.Equal(t, *result.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 409,
		Status:     "ObjectNotAppendable",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 409)
	assert.Equal(t, result.Status, "ObjectNotAppendable")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteObjectRequest
	var input *OperationInput
	var err error

	request = &DeleteObjectRequest{}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteObjectRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteObjectRequest{
		Bucket: Ptr("oss-bucket"),
		Key:    Ptr("oss-key"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Nil(t, input.OpMetadata.values)
	assert.Nil(t, input.Parameters)
	request = &DeleteObjectRequest{
		Bucket:    Ptr("oss-bucket"),
		Key:       Ptr("oss-key"),
		VersionId: Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY****"),
	}
	input = &OperationInput{
		OpName: "DeleteObject",
		Method: "DELETE",
		Bucket: request.Bucket,
		Key:    request.Key,
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, *input.Key, "oss-key")
	assert.Equal(t, input.Parameters["versionId"], "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY****")
	assert.Nil(t, input.OpMetadata.values)
}

func TestUnmarshalOutput_DeleteObject(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"3a8f-2e2d-7965-3ff9-51c875b*****"},
			"Content-Type":     {"image/jpg"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
		},
	}
	result := &DeleteObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "3a8f-2e2d-7965-3ff9-51c875b*****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Nil(t, result.VersionId)
	assert.False(t, result.DeleteMarker)

	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id":    {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":                {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"x-oss-version-id":    {"CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****"},
			"x-oss-delete-marker": {"true"},
		},
	}
	result = &DeleteObjectResult{}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")

	assert.Equal(t, *result.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
	assert.True(t, result.DeleteMarker)

	output = &OperationOutput{
		StatusCode: 403,
		Status:     "AccessDenied",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	err = c.unmarshalOutput(result, output, discardBody, unmarshalHeader)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteMultipleObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteMultipleObjectsRequest
	var input *OperationInput
	var err error

	request = &DeleteMultipleObjectsRequest{}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteMultipleObjectsRequest{
		Bucket: Ptr("oss-bucket"),
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	request = &DeleteMultipleObjectsRequest{
		Bucket:  Ptr("oss-bucket"),
		Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, input.Body, strings.NewReader("<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>"))
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["delete"])
	request = &DeleteMultipleObjectsRequest{
		Bucket:       Ptr("oss-bucket"),
		Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		EncodingType: Ptr("url"),
		Quiet:        true,
	}
	input = &OperationInput{
		OpName: "DeleteMultipleObjects",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{"delete": ""},
		Bucket:     request.Bucket,
	}
	err = c.marshalInput(request, input, marshalDeleteObjects, updateContentMd5)
	assert.Nil(t, err)
	assert.Equal(t, *input.Bucket, "oss-bucket")
	assert.Equal(t, input.Body, strings.NewReader("<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>"))
	assert.Nil(t, input.OpMetadata.values)
	assert.Empty(t, input.Parameters["delete"])
	assert.Equal(t, input.Parameters["encoding-type"], "url")
}

func TestUnmarshalOutput_DeleteMultipleObjects(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555A936CA31DC333143****"},
			"Date":             {"Thu, 16 Nov 2023 05:31:34 GMT"},
		},
	}
	result := &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "6555A936CA31DC333143****")
	assert.Equal(t, result.Headers.Get("Date"), "Thu, 16 Nov 2023 05:31:34 GMT")

	body := `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 2)
	assert.Equal(t, *result.DeletedObjects[0].Key, "key1.txt")
	assert.Equal(t, result.DeletedObjects[0].DeleteMarker, true)
	assert.Equal(t, *result.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
	assert.Nil(t, result.DeletedObjects[0].VersionId)
	assert.Equal(t, *result.DeletedObjects[1].Key, "key2.txt")
	assert.Equal(t, result.DeletedObjects[1].DeleteMarker, true)
	assert.Equal(t, *result.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
	assert.Nil(t, result.DeletedObjects[1].VersionId)

	body = `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>key1.txt</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****</VersionId>
  </Deleted>
  <Deleted>
    <Key>key2.txt</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmY1****</VersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 2)
	assert.Equal(t, *result.DeletedObjects[0].Key, "key1.txt")
	assert.False(t, result.DeletedObjects[0].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[0].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[0].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****")
	assert.Equal(t, *result.DeletedObjects[1].Key, "key2.txt")
	assert.False(t, result.DeletedObjects[1].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[1].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[1].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmY1****")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <VersionId>CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****</VersionId>
  </Deleted>
</DeleteResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"28f6-15ea-8224-234e-c0ce407****"},
			"Date":             {"Tue, 04 Dec 2018 15:56:38 GMT"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	result = &DeleteMultipleObjectsResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "28f6-15ea-8224-234e-c0ce407****")
	assert.Equal(t, result.Headers.Get("Date"), "Tue, 04 Dec 2018 15:56:38 GMT")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	assert.Len(t, result.DeletedObjects, 1)
	assert.Equal(t, *result.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
	assert.False(t, result.DeletedObjects[0].DeleteMarker)
	assert.Nil(t, result.DeletedObjects[0].DeleteMarkerVersionId)
	assert.Equal(t, *result.DeletedObjects[0].VersionId, "CAEQFxiBgIDztZ2IuRgiIDMyNzg1MTY1NWI5NjQyOGJiZWIwOTA0NTI0MmYx****")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "MalformedXML",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "MalformedXML")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	body = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ErrorDetail>the root node is not named Delete.</ErrorDetail>
  <EC>0016-00000608</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000608</RecommendDoc>
</Error>`
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "MalformedXML",
		Headers: http.Header{
			"X-Oss-Request-Id": {"6555AC764311A73931E0****"},
			"Content-Type":     {"application/xml"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(body))),
	}
	var resultErr DeleteMultipleObjectsResult
	err = c.unmarshalOutput(resultErr, output, unmarshalBodyXml, unmarshalHeader, unmarshalEncodeType)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "deserialization failed, non-pointer passed to Unmarshal")
}
