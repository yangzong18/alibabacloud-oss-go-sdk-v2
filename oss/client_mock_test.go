package oss

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/credentials"
	"github.com/stretchr/testify/assert"
)

func sortQuery(r *http.Request) string {
	u := r.URL
	var buf strings.Builder
	keys := make([]string, 0, len(u.Query()))
	for k := range u.Query() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := u.Query()[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			if len(v) > 0 {
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(v))
			}
		}
	}
	u.RawQuery = buf.String()
	return u.String()
}

func testSetupMockServer(t *testing.T, statusCode int, headers map[string]string, body []byte,
	chkfunc func(t *testing.T, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check request
		chkfunc(t, r)

		// header s
		for k, v := range headers {
			w.Header().Set(k, v)
		}

		// status code
		w.WriteHeader(statusCode)

		// body
		w.Write(body)
	}))
}

var testInvokeOperationAnonymousCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "5374A2880232A65C2300****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<ListAllMyBucketsResult>
			<Owner>
				<ID>512**</ID>
				<DisplayName>51264</DisplayName>
			</Owner>
			<Buckets>
				<Bucket>
				<CreationDate>2014-02-17T18:12:43.000Z</CreationDate>
				<ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
				<IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
				<Location>oss-cn-shanghai</Location>
				<Name>app-base-oss</Name>
				<Region>cn-shanghai</Region>
				<StorageClass>Standard</StorageClass>
				</Bucket>
			</Buckets>				
			</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&OperationInput{
			OpName: "ListBuckets",
			Method: "GET",
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5374A2880232A65C2300****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 15 May 2014 11:18:32 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetBucketAcl",
			Bucket: Ptr("bucket"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "654605AA6172673135811AB3", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Sat, 04 Nov 2023 08:49:46 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "654605AA6172673135811AB3",
			"Date":             "Sat, 04 Nov 2023 08:49:46 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<AccessControlPolicy>
				<Owner>
					<ID>12345</ID>
					<DisplayName>12345Name</DisplayName>
				</Owner>
				<AccessControlList>
					<Grant>private</Grant>
				</AccessControlList>
			</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/key?acl", r.URL.String())
		},
		&OperationInput{
			OpName: "GetObjectAcl",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key"),
			Method: "GET",
			Parameters: map[string]string{
				"acl": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<InitiateMultipartUploadResult>
				<Bucket>oss-example</Bucket>
				<Key>key+ 123.data</Key>
				<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
			</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/key%2B%20123/test.data?uploads", r.URL.String())
			assert.Equal(t, "POST", r.Method)
		},
		&OperationInput{
			OpName: "InitiateMultipartUpload",
			Bucket: Ptr("bucket"),
			Key:    Ptr("key+ 123/test.data"),
			Method: "POST",
			Parameters: map[string]string{
				"uploads": "",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
	{
		200,
		map[string]string{
			"Content-Type": "text/txt",
		},
		[]byte(
			`hello world`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket//subfolder/example.txt?versionId=CAEQNhiBgMDJgZCA0BY%2B123", r.URL.String())
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Etag1234", r.Header.Get("If-Match"))
		},
		&OperationInput{
			OpName: "GetObject",
			Bucket: Ptr("bucket"),
			Key:    Ptr("/subfolder/example.txt"),
			Method: "GET",
			Headers: map[string]string{
				"If-Match": "Etag1234",
			},
			Parameters: map[string]string{
				"versionId": "CAEQNhiBgMDJgZCA0BY+123",
			},
		},
		func(t *testing.T, o *OperationOutput) {
		},
	},
}

func TestInvokeOperation_Anonymous(t *testing.T) {
	for _, c := range testInvokeOperationAnonymousCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)

		var fns []func(*Options)
		fns = append(fns, func(c *Options) { c.OpReadWriteTimeout = Ptr(1 * time.Second) })
		output, err = client.InvokeOperation(context.TODO(), c.Input, fns...)
		assert.Nil(t, err)
		c.CheckOutputFn(t, output)
	}
}

var testInvokeOperationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Input          *OperationInput
	CheckOutputFn  func(t *testing.T, o *OperationOutput, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/test-key.txt", r.URL.String())
		},
		&OperationInput{
			OpName: "PutObject",
			Method: "PUT",
			Bucket: Ptr("bucket"),
			Key:    Ptr("test-key.txt"),
		},
		func(t *testing.T, o *OperationOutput, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/bucket/test-key.txt")
		},
	},
}

func TestInvokeOperation_Error(t *testing.T) {
	for _, c := range testInvokeOperationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InvokeOperation(context.TODO(), c.Input)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket", r.URL.String())
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, "<CreateBucketConfiguration><StorageClass>Archive</StorageClass><DataRedundancyType>LRS</DataRedundancyType></CreateBucketConfiguration>", string(requestBody))
		},
		&PutBucketRequest{
			Bucket:          Ptr("bucket"),
			Acl:             BucketACLPrivate,
			ResourceGroupId: Ptr("rg-aek27tc********"),
			CreateBucketConfiguration: &CreateBucketConfiguration{
				StorageClass:       StorageClassArchive,
				DataRedundancyType: DataRedundancyLRS,
			},
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucket_Success(t *testing.T) {
	for _, c := range testMockPutBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>SignatureDoesNotMatch</Code>
				<Message>The request signature we calculated does not match the signature you provided. Check your key and signing method.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "SignatureDoesNotMatch", serr.Code)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
			assert.Contains(t, serr.Message, "The request signature we calculated does not match")
			assert.Contains(t, serr.RequestTarget, "/bucket")
		},
	},
	{
		409,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>BucketAlreadyExists</Code>
				<Message>The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again.</Message>
				<RequestId>6548A043CA31D****</RequestId>
				<EC>0015-00000104</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&PutBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *PutBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "BucketAlreadyExists", serr.Code)
			assert.Equal(t, "0015-00000104", serr.EC)
			assert.Equal(t, "6548A043CA31D****", serr.RequestID)
			assert.Contains(t, serr.Message, "The requested bucket name is not available. The bucket namespace is shared by all users of the system. Please select a different name and try again")
			assert.Contains(t, serr.RequestTarget, "/bucket")
		},
	},
}

func TestMockPutBucket_Error(t *testing.T) {
	for _, c := range testMockPutBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Owner>
    <ID>51264</ID>
    <DisplayName>51264</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <CreationDate>2014-02-17T18:12:43.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-shanghai</Location>
      <Name>app-base-oss</Name>
      <Region>cn-shanghai</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
    <Bucket>
      <CreationDate>2014-02-25T11:21:04.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-hangzhou</Location>
      <Name>mybucket</Name>
      <Region>cn-hangzhou</Region>
      <StorageClass>IA</StorageClass>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		nil,
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Owner.DisplayName, "51264")
			assert.Equal(t, *o.Owner.ID, "51264")
			assert.Equal(t, len(o.Buckets), 2)
			assert.Equal(t, *o.Buckets[0].CreationDate, time.Date(2014, time.February, 17, 18, 12, 43, 0, time.UTC))
			assert.Equal(t, *o.Buckets[0].ExtranetEndpoint, "oss-cn-shanghai.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].IntranetEndpoint, "oss-cn-shanghai-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].Name, "app-base-oss")
			assert.Equal(t, *o.Buckets[0].Region, "cn-shanghai")
			assert.Equal(t, *o.Buckets[0].StorageClass, "Standard")

			assert.Equal(t, *o.Buckets[1].CreationDate, time.Date(2014, time.February, 25, 11, 21, 04, 0, time.UTC))
			assert.Equal(t, *o.Buckets[1].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[1].Name, "mybucket")
			assert.Equal(t, *o.Buckets[1].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets[1].StorageClass, "IA")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Prefix>my</Prefix>
  <Marker>mybucket</Marker>
  <MaxKeys>10</MaxKeys>
  <IsTruncated>true</IsTruncated>
  <NextMarker>mybucket10</NextMarker>
  <Owner>
    <ID>ut_test_put_bucket</ID>
    <DisplayName>ut_test_put_bucket</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <CreationDate>2014-05-14T11:18:32.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-hangzhou.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-hangzhou-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-hangzhou</Location>
      <Name>mybucket01</Name>
      <Region>cn-hangzhou</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F", strUrl)
		},
		&ListBucketsRequest{
			Marker:          Ptr(""),
			MaxKeys:         10,
			Prefix:          Ptr("/"),
			ResourceGroupId: Ptr("rg-aek27tc********"),
		},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Owner.DisplayName, "ut_test_put_bucket")
			assert.Equal(t, *o.Owner.ID, "ut_test_put_bucket")
			assert.Equal(t, *o.Prefix, "my")
			assert.Equal(t, *o.Marker, "mybucket")
			assert.Equal(t, o.MaxKeys, int32(10))
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, *o.NextMarker, "mybucket10")

			assert.Equal(t, len(o.Buckets), 1)
			assert.Equal(t, *o.Buckets[0].CreationDate, time.Date(2014, time.May, 14, 11, 18, 32, 0, time.UTC))
			assert.Equal(t, *o.Buckets[0].ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.Buckets[0].Name, "mybucket01")
			assert.Equal(t, *o.Buckets[0].Region, "cn-hangzhou")
			assert.Equal(t, *o.Buckets[0].StorageClass, "Standard")
		},
	},
}

func TestMockListBuckets_Success(t *testing.T) {
	for _, c := range testMockListBucketsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketsRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketsResult, err error)
}{
	{
		403,
		map[string]string{
			"x-oss-request-id": "65467C42E001B4333337****",
			"Date":             "Thu, 15 May 2014 11:18:32 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(
			`<?xml version="1.0" encoding="UTF-8"?>
			<Error>
				<Code>InvalidAccessKeyId</Code>
				<Message>The OSS Access Key Id you provided does not exist in our records.</Message>
				<RequestId>65467C42E001B4333337****</RequestId>
				<SignatureProvided>RizTbeKC/QlwxINq8xEdUPowc84=</SignatureProvided>
				<EC>0002-00000040</EC>
			</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&ListBucketsRequest{},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "InvalidAccessKeyId", serr.Code)
			assert.Equal(t, "The OSS Access Key Id you provided does not exist in our records.", serr.Message)
			assert.Equal(t, "0002-00000040", serr.EC)
			assert.Equal(t, "65467C42E001B4333337****", serr.RequestID)
		},
	},
}

func TestMockListBuckets_Error(t *testing.T) {
	for _, c := range testMockListBucketsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBuckets(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketResult, err error)
}{
	{
		204,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockDeleteBucket_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>BucketNotEmpty</Code>
  <Message>The bucket has objects. Please delete them first.</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000301</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket", r.URL.String())
		},
		&DeleteBucketRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *DeleteBucketResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "BucketNotEmpty", serr.Code)
			assert.Equal(t, "The bucket has objects. Please delete them first.", serr.Message)
			assert.Equal(t, "0015-00000301", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteBucket_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucket(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?encoding-type=url", r.URL.String())
		},
		&ListObjectsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Empty(t, o.Prefix)
			assert.Equal(t, *o.Name, "examplebucket")
			assert.Empty(t, o.Marker)
			assert.Empty(t, o.Delimiter)
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, len(o.Contents), 4)
			assert.Equal(t, *o.Contents[0].Key, "fun/movie/001.avi")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2012, time.February, 24, 8, 43, 27, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[3].Type, "Normal")
			assert.Equal(t, o.Contents[0].Size, int64(344606))
			assert.Equal(t, *o.Contents[1].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[2].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[3].Owner.DisplayName, "user-example")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
  <Prefix>fun</Prefix>
  <Marker>test1.txt</Marker>
  <MaxKeys>3</MaxKeys>
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
        <Key>exampleobject3.txt</Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "fun")
			assert.Equal(t, *o.Marker, "test1.txt")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "exampleobject1.txt")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Empty(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
			assert.Equal(t, *o.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
		},
	},
}

func TestMockListObjects_Success(t *testing.T) {
	for _, c := range testMockListObjectsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?encoding-type=url", r.URL.String())
		},
		&ListObjectsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListObjects_Error(t *testing.T) {
	for _, c := range testMockListObjectsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsV2SuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequestV2
	CheckOutputFn  func(t *testing.T, o *ListObjectsResultV2, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
  <Name>examplebucket</Name>
  <Prefix></Prefix>
  <MaxKeys>3</MaxKeys>
  <Delimiter></Delimiter>
  <IsTruncated>false</IsTruncated>
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
        <Key>exampleobject3.txt</Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Empty(t, o.Prefix)
			assert.Equal(t, *o.Name, "examplebucket")
			assert.Empty(t, o.Delimiter)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "exampleobject1.txt")
			assert.Equal(t, *o.Contents[0].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[0].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[0].Size, int64(344606))
			assert.Equal(t, *o.Contents[0].StorageClass, "ColdArchive")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")

			assert.Equal(t, *o.Contents[1].RestoreInfo, "ongoing-request=\"true\"")
			assert.Equal(t, *o.Contents[2].RestoreInfo, "ongoing-request=\"false\", expiry-date=\"Thu, 24 Sep 2020 12:40:33 GMT\"")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
<Name>examplebucket</Name>
    <Prefix>a/</Prefix>
    <MaxKeys>3</MaxKeys>
	<StartAfter>b</StartAfter>
    <Delimiter>/</Delimiter>
    <EncodingType>url</EncodingType>
    <IsTruncated>false</IsTruncated>
  	<Contents>
        <Key>a/b</Key>
        <LastModified>2020-05-18T05:45:47.000Z</LastModified>
        <ETag>"35A27C2B9EAEEB6F48FD7FB5861D****"</ETag>
		<Type>Normal</Type>
        <Size>25</Size>
        <StorageClass>STANDARD</StorageClass>
		<Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
	</Contents>
  	<Contents>
        <Key>a/b/c</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
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
        <Key>a/b/d</Key>
        <LastModified>2020-06-22T11:42:32.000Z</LastModified>
        <ETag>"5B3C1A2E053D763E1B002CC607C5A0FE1****"</ETag>
        <Type>Normal</Type>
        <Size>344606</Size>
        <StorageClass>Standard</StorageClass>
        <Owner>
            <ID>0022012****</ID>
            <DisplayName>user-example</DisplayName>
        </Owner>
  </Contents>
	<CommonPrefixes>
        <Prefix>a/b/</Prefix>
    </CommonPrefixes>
    <KeyCount>3</KeyCount>
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "a/")
			assert.Equal(t, *o.StartAfter, "b")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxKeys, int32(3))
			assert.Equal(t, o.KeyCount, 3)
			assert.Equal(t, len(o.Contents), 3)
			assert.Equal(t, *o.Contents[0].Key, "a/b")
			assert.Equal(t, *o.Contents[1].LastModified, time.Date(2020, time.June, 22, 11, 42, 32, 0, time.UTC))
			assert.Equal(t, *o.Contents[2].ETag, "\"5B3C1A2E053D763E1B002CC607C5A0FE1****\"")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, o.Contents[1].Size, int64(344606))
			assert.Equal(t, *o.Contents[2].StorageClass, "Standard")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Nil(t, o.Contents[0].RestoreInfo)
			assert.Equal(t, *o.CommonPrefixes[0].Prefix, "a/b/")
		},
	},
}

func TestMockListObjectsV2_Success(t *testing.T) {
	for _, c := range testMockListObjectsV2SuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectsV2(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectsV2ErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectsRequestV2
	CheckOutputFn  func(t *testing.T, o *ListObjectsResultV2, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsRequestV2{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsResultV2, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListObjectsV2_Error(t *testing.T) {
	for _, c := range testMockListObjectsV2ErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectsV2(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInfoSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInfoRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInfoResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketInfo>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketInfo.Name, "oss-example")
			assert.Equal(t, *o.BucketInfo.AccessMonitor, "Enabled")
			assert.Equal(t, *o.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.BucketInfo.StorageClass, "Standard")
			assert.Equal(t, *o.BucketInfo.TransferAcceleration, "Disabled")
			assert.Equal(t, *o.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
			assert.Equal(t, *o.BucketInfo.CrossRegionReplication, "Disabled")
			assert.Equal(t, *o.BucketInfo.ResourceGroupId, "rg-aek27tc********")
			assert.Equal(t, *o.BucketInfo.Owner.ID, "27183473914****")
			assert.Equal(t, *o.BucketInfo.Owner.DisplayName, "username")
			assert.Equal(t, *o.BucketInfo.ACL, "private")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogPrefix, "log/")
			assert.Empty(t, *o.BucketInfo.SseRule.KMSMasterKeyID)
			assert.Equal(t, *o.BucketInfo.SseRule.SSEAlgorithm, "KMS")
			assert.Equal(t, *o.BucketInfo.SseRule.KMSDataEncryption, "SM4")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketInfo>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.BucketInfo.Name, "oss-example")
			assert.Equal(t, *o.BucketInfo.AccessMonitor, "Enabled")
			assert.Equal(t, *o.BucketInfo.ExtranetEndpoint, "oss-cn-hangzhou.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.IntranetEndpoint, "oss-cn-hangzhou-internal.aliyuncs.com")
			assert.Equal(t, *o.BucketInfo.Location, "oss-cn-hangzhou")
			assert.Equal(t, *o.BucketInfo.StorageClass, "Standard")
			assert.Equal(t, *o.BucketInfo.TransferAcceleration, "Disabled")
			assert.Equal(t, *o.BucketInfo.CreationDate, time.Date(2013, time.July, 31, 10, 56, 21, 0, time.UTC))
			assert.Equal(t, *o.BucketInfo.CrossRegionReplication, "Disabled")
			assert.Equal(t, *o.BucketInfo.ResourceGroupId, "rg-aek27tc********")
			assert.Equal(t, *o.BucketInfo.Owner.ID, "27183473914****")
			assert.Equal(t, *o.BucketInfo.Owner.DisplayName, "username")
			assert.Equal(t, *o.BucketInfo.ACL, "private")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogBucket, "examplebucket")
			assert.Equal(t, *o.BucketInfo.BucketPolicy.LogPrefix, "log/")

			assert.Empty(t, o.BucketInfo.SseRule.KMSMasterKeyID)
			assert.Nil(t, o.BucketInfo.SseRule.SSEAlgorithm)
			assert.Nil(t, o.BucketInfo.SseRule.KMSDataEncryption)
		},
	},
}

func TestMockGetBucketInfo_Success(t *testing.T) {
	for _, c := range testMockGetBucketInfoSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInfo(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInfoErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInfoRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInfoResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?bucketInfo", r.URL.String())
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?bucketInfo", strUrl)
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetBucketInfo_Error(t *testing.T) {
	for _, c := range testMockGetBucketInfoErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInfo(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLocationSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLocationResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-hangzhou</LocationConstraint>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.LocationConstraint, "oss-cn-hangzhou")
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<LocationConstraint>oss-cn-chengdu</LocationConstraint>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.LocationConstraint, "oss-cn-chengdu")
		},
	},
}

func TestMockGetBucketLocation_Success(t *testing.T) {
	for _, c := range testMockGetBucketLocationSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketLocationErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketLocationRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketLocationResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?location", r.URL.String())
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetBucketLocation_Error(t *testing.T) {
	for _, c := range testMockGetBucketLocationErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketLocation(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStatSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStatRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStatResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</BucketStat>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?stat", r.URL.String())
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, int64(1600), o.Storage)
			assert.Equal(t, int64(230), o.ObjectCount)
			assert.Equal(t, int64(40), o.MultipartUploadCount)
			assert.Equal(t, int64(4), o.LiveChannelCount)
			assert.Equal(t, int64(1643341269), o.LastModifiedTime)
			assert.Equal(t, int64(430), o.StandardStorage)
			assert.Equal(t, int64(66), o.StandardObjectCount)
			assert.Equal(t, int64(2359296), o.InfrequentAccessStorage)
			assert.Equal(t, int64(360), o.InfrequentAccessRealStorage)
			assert.Equal(t, int64(54), o.InfrequentAccessObjectCount)
			assert.Equal(t, int64(2949120), o.ArchiveStorage)
			assert.Equal(t, int64(450), o.ArchiveRealStorage)
			assert.Equal(t, int64(74), o.ArchiveObjectCount)
			assert.Equal(t, int64(2359296), o.ColdArchiveStorage)
			assert.Equal(t, int64(360), o.ColdArchiveRealStorage)
			assert.Equal(t, int64(36), o.ColdArchiveObjectCount)
		},
	},
}

func TestMockGetBucketStat_Success(t *testing.T) {
	for _, c := range testMockGetBucketStatSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketStat(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketStatErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketStatRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketStatResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?stat", r.URL.String())
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?stat", strUrl)
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetBucketStat_Error(t *testing.T) {
	for _, c := range testMockGetBucketStatErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketStat(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPublicRead), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPublicRead,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPublicReadWrite), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPublicReadWrite,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketAcl_Success(t *testing.T) {
	for _, c := range testMockPutBucketAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?acl", strUrl)
			assert.Equal(t, string(BucketACLPrivate), r.Header.Get("X-Oss-Acl"))
		},
		&PutBucketAclRequest{
			Bucket: Ptr("bucket"),
			Acl:    BucketACLPrivate,
		},
		func(t *testing.T, o *PutBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketAcl_Error(t *testing.T) {
	for _, c := range testMockPutBucketAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read", *o.ACL)
			assert.Equal(t, "0022012****", *o.Owner.ID)
			assert.Equal(t, "user_example", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read-write</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read-write", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
}

func TestMockGetBucketAcl_Success(t *testing.T) {
	for _, c := range testMockGetBucketAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketAclRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?acl", strUrl)
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket?acl", r.URL.String())
		},
		&GetBucketAclRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockGetBucketAcl_Error(t *testing.T) {
	for _, c := range testMockGetBucketAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "316181249502703****",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "316181249502703****")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "870718044876840****",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",

			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "870718044876840****",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`{"filename":"object","size":"6","mimeType":""}`),
		func(t *testing.T, r *http.Request) {
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.NotNil(t, r.Header.Get("x-oss-callback"))
		},
		&PutObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			resultBody, err := ioutil.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(resultBody), `{"filename":"object","size":"6","mimeType":""}`)
		},
	},
}

func TestMockPutObject_Success(t *testing.T) {
	for _, c := range testMockPutObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		203,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "5C3D9175B6FC201293AD****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "870718044876840****",
			"Content-MD5":          "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-version-id":     "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>CallbackFailed</Code>
  <Message>Error status : 301.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0007-00000203</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0007-00000203</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&PutObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
			Callback: Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"http://www.aliyun.com","callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}"}`))),
		},
		func(t *testing.T, o *PutObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(203), serr.StatusCode)
			assert.Equal(t, "CallbackFailed", serr.Code)
			assert.Equal(t, "Error status : 301.", serr.Message)
			assert.Equal(t, "0007-00000203", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockPutObject_Error(t *testing.T) {
	for _, c := range testMockPutObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":         "application/xml",
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"D41D8CD98F00B204E9800998ECF8****\"",
			"x-oss-hash-crc64ecma": "316181249502703****",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7PhC****",
		},
		[]byte(`hi oss,this is a demo!`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7PhC****")
			assert.Equal(t, *o.HashCRC64, "316181249502703****")
			content, err := ioutil.ReadAll(o.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(content), "hi oss,this is a demo!")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"5B3C1A2E05E1B002CC607C****\"",
			"Content-Length":                      "344606",
			"Last-Modified":                       "Fri, 24 Feb 2012 06:07:48 GMT",
			"x-oss-object-type":                   "Normal",
			"Accept-Ranges":                       "bytes",
			"Content-disposition":                 "attachment; filename=testing.txt",
			"Cache-control":                       "no-cache",
			"X-Oss-Storage-Class":                 "Standard",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-tagging-count":                 "2",
			"Content-MD5":                         "si4Nw3Cn9wZ/rPX3XX+j****",
			"x-oss-hash-crc64ecma":                "870718044876840****",
			"x-oss-meta-name":                     "demo",
			"x-oss-meta-email":                    "demo@aliyun.com",
		},
		[]byte(`hi oss`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E05E1B002CC607C****\"")
			assert.Equal(t, *o.LastModified, time.Date(2012, time.February, 24, 6, 7, 48, 0, time.UTC))
			assert.Equal(t, *o.ContentType, "text")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.StorageClass, "Standard")
			content, err := ioutil.ReadAll(o.Body)
			assert.Equal(t, string(content), "hi oss")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, o.TaggingCount, int32(2))
			assert.Equal(t, o.Metadata["name"], "demo")
			assert.Equal(t, o.Metadata["email"], "demo@aliyun.com")
			assert.Equal(t, *o.ContentMD5, "si4Nw3Cn9wZ/rPX3XX+j****")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
		},
	},
}

func TestMockGetObject_Success(t *testing.T) {
	for _, c := range testMockGetObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&GetObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *GetObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockGetObject_Error(t *testing.T) {
	for _, c := range testMockGetObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":                 "application/xml",
			"x-oss-request-id":             "534B371674E88A4D8906****",
			"Date":                         "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                         "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma":         "870718044876840****",
			"x-oss-copy-source-version-id": "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****",
			"x-oss-version-id":             "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Source: Ptr("/bucket/copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876840****")
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk3****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":         "text",
			"ETag":                 "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-hash-crc64ecma": "870718044876841****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Source: Ptr("/bucket/copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":                        "text",
			"x-oss-version-id":                    "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****",
			"ETag":                                "\"F2064A169EE92E9775EE5324D0B1****\"",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
			"x-oss-hash-crc64ecma":                "870718044876841****",
			"x-oss-copy-source-version-id":        "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
	<CopyObjectResult>
	 <ETag>"F2064A169EE92E9775EE5324D0B1****"</ETag>
	 <LastModified>2023-02-24T09:41:56.000Z</LastModified>
	</CopyObjectResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Source: Ptr("/bucket/copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"F2064A169EE92E9775EE5324D0B1****\"")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
			assert.Equal(t, *o.HashCRC64, "870718044876841****")
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD0ooWf3hgiIDcyMzYzZTJkZjgwYzRmN2FhNTZkMWZlMGY0YTVj****")
			assert.Equal(t, *o.SourceVersionId, "CAEQHxiBgICDvseg3hgiIGZmOGNjNWJiZDUzNjQxNDM4MWM2NDc1YjhkYTk4****")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
		},
	},
}

func TestMockCopyObject_Success(t *testing.T) {
	for _, c := range testMockCopyObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCopyObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CopyObjectRequest
	CheckOutputFn  func(t *testing.T, o *CopyObjectResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Source: Ptr("/bucket/copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
			assert.Equal(t, "/bucket/copy-object", r.Header.Get("x-oss-copy-source"))
		},
		&CopyObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Source: Ptr("/bucket/copy-object"),
		},
		func(t *testing.T, o *CopyObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockCopyObject_Error(t *testing.T) {
	for _, c := range testMockCopyObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CopyObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":               "application/xml",
			"x-oss-request-id":           "534B371674E88A4D8906****",
			"Date":                       "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-next-append-position": "1717",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss,append object"),
			},
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Nil(t, o.VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":           "6551DBCF4311A7303980****",
			"Date":                       "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":           "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position": "0",
			"x-oss-hash-crc64ecma":       "1474161709526656****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss,append object,this is a demo"),
			},
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(0))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":                    "6551DBCF4311A7303980****",
			"Date":                                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":                    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-next-append-position":          "1717",
			"x-oss-hash-crc64ecma":                "1474161709526656****",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "12f8711f-90df-4e0d-903d-ab972b0f****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "12f8711f-90df-4e0d-903d-ab972b0f****")
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss,append object,this is a demo"),
			},
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("12f8711f-90df-4e0d-903d-ab972b0f****"),
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.Equal(t, *o.HashCRC64, "1474161709526656****")
			assert.Equal(t, o.NextPosition, int64(1717))
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.SSEKMSKeyId, "12f8711f-90df-4e0d-903d-ab972b0f****")
		},
	},
}

func TestMockAppendObject_Success(t *testing.T) {
	for _, c := range testMockAppendObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAppendObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AppendObjectRequest
	CheckOutputFn  func(t *testing.T, o *AppendObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			strUrl := sortQuery(r)
			assert.Equal(t, strings.NewReader("hi oss,append object"), strings.NewReader(string(requestBody)))
			assert.Equal(t, "/bucket/object?append&position=100", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(100)),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss,append object"),
			},
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>PositionNotEqualToLength</Code>
  <Message>Position is not equal to file length</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0026-00000016</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			requestBody, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, strings.NewReader("hi oss,append object,this is a demo"), strings.NewReader(string(requestBody)))
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?append&position=0", strUrl)
		},
		&AppendObjectRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			Position: Ptr(int64(0)),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss,append object,this is a demo"),
			},
		},
		func(t *testing.T, o *AppendObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "PositionNotEqualToLength", serr.Code)
			assert.Equal(t, "Position is not equal to file length", serr.Message)
			assert.Equal(t, "0026-00000016", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockAppendObject_Error(t *testing.T) {
	for _, c := range testMockAppendObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AppendObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
}{
	{
		204,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionId)
			assert.False(t, o.DeleteMarker)
		},
	},
	{
		204,
		map[string]string{
			"x-oss-request-id":    "6551DBCF4311A7303980****",
			"Date":                "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":    "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****",
			"x-oss-delete-marker": "true",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQHxiBgMD4qOWz3hgiIDUyMWIzNTBjMWM4NjQ5MDJiNTM4YzEwZGQxM2Rk****")
			assert.True(t, o.DeleteMarker)
		},
	},
}

func TestMockDeleteObject_Success(t *testing.T) {
	for _, c := range testMockDeleteObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteObjectRequest
	CheckOutputFn  func(t *testing.T, o *DeleteObjectResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&DeleteObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *DeleteObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockDeleteObject_Error(t *testing.T) {
	for _, c := range testMockDeleteObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>true</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
			Quiet:   true,
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Nil(t, o.DeletedObjects)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		}, []byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), ("<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>"))
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
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
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 2)
			assert.Equal(t, *o.DeletedObjects[0].Key, "key1.txt")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
			assert.Equal(t, *o.DeletedObjects[1].Key, "key2.txt")
			assert.Equal(t, o.DeletedObjects[1].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[1].DeleteMarkerVersionId, "CAEQHxiBgICJld7a3hgiIDJmZGE0OTU5MjMzZDQxNjlhY2NjMmI3YWRkYWI4****")
			assert.Nil(t, o.DeletedObjects[1].VersionId)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<DeleteResult>
  <EncodingType>url</EncodingType>
  <Deleted>
    <Key>go-sdk-v1%01%02%03%04%05%06%07%08%09%0A%0B%0C%0D%0E%0F%10%11%12%13%14%15%16%17%18%19%1A%1B%1C%1D%1E%1F</Key>
    <DeleteMarker>true</DeleteMarker>
    <DeleteMarkerVersionId>CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****</DeleteMarkerVersionId>
  </Deleted>
</DeleteResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>go-sdk-v1&#x01;&#x02;&#x03;&#x04;&#x05;&#x06;&#x07;&#x08;&#x9;&#xA;&#x0B;&#x0C;&#xD;&#x0E;&#x0F;&#x10;&#x11;&#x12;&#x13;&#x14;&#x15;&#x16;&#x17;&#x18;&#x19;&#x1A;&#x1B;&#x1C;&#x1D;&#x1E;&#x1F;</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:       Ptr("bucket"),
			Objects:      []DeleteObject{{Key: Ptr("go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}},
			EncodingType: Ptr("url"),
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Len(t, o.DeletedObjects, 1)
			assert.Equal(t, *o.DeletedObjects[0].Key, "go-sdk-v1\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f")
			assert.Equal(t, o.DeletedObjects[0].DeleteMarker, true)
			assert.Equal(t, *o.DeletedObjects[0].DeleteMarkerVersionId, "CAEQHxiBgMCEld7a3hgiIDYyMmZlNWVhMDU5NDQ3ZTFhODI1ZjZhMTFlMGQz****")
			assert.Nil(t, o.DeletedObjects[0].VersionId)
		},
	},
}

func TestMockDeleteMultipleObjects_Success(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteMultipleObjectsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteMultipleObjectsRequest
	CheckOutputFn  func(t *testing.T, o *DeleteMultipleObjectsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****</VersionId></Object><Object><Key>key2.txt</Key><VersionId>CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****</VersionId></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA4****")}, {Key: Ptr("key2.txt"), VersionId: Ptr("CAEQNRiBgIDyz.6C0BYiIGQ2NWEwNmVhNTA3ZTQ3MzM5ODliYjM1ZTdjYjA5****")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ErrorDetail>the root node is not named Delete.</ErrorDetail>
  <EC>0016-00000608</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0016-00000608</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delete&encoding-type=url", strUrl)
			data, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, string(data), "<Delete><Quiet>false</Quiet><Object><Key>key1.txt</Key></Object><Object><Key>key2.txt</Key></Object></Delete>")
		},
		&DeleteMultipleObjectsRequest{
			Bucket:  Ptr("bucket"),
			Objects: []DeleteObject{{Key: Ptr("key1.txt")}, {Key: Ptr("key2.txt")}},
		},
		func(t *testing.T, o *DeleteMultipleObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "0016-00000608", serr.EC)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
		},
	},
}

func TestMockDeleteMultipleObjects_Error(t *testing.T) {
	for _, c := range testMockDeleteMultipleObjectsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteMultipleObjects(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "6555A936CA31DC333143****",
			"Date":                "Thu, 16 Nov 2023 05:31:34 GMT",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":      "344606",
			"Content-Type":        "image/jpg",
			"ETag":                "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(344606))
			assert.Equal(t, *o.ContentType, "image/jpg")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "5CAC3B40B7AEADE01700****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "text/xml",
			"x-oss-object-type":   "Normal",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id":    "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "text/xml")
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":    "534B371674E88A4D8906****",
			"Date":                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":        "image/jpg",
			"x-oss-object-type":   "Normal",
			"x-oss-restore":       "ongoing-request=\"true\"",
			"x-oss-storage-class": "Archive",
			"Last-Modified":       "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":      "481827",
			"ETag":                "\"A082B659EF78733A5A042FA253B1****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"true\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":                    "534B371674E88A4D8906****",
			"Date":                                "Tue, 04 Dec 2018 15:56:38 GMT",
			"Content-Type":                        "image/jpg",
			"x-oss-object-type":                   "Normal",
			"x-oss-restore":                       "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"",
			"x-oss-storage-class":                 "Archive",
			"x-oss-server-side-encryption":        "KMS",
			"x-oss-server-side-data-encryption":   "SM4",
			"x-oss-server-side-encryption-key-id": "9468da86-3509-4f8d-a61e-6eab1eac****",
			"Content-Length":                      "481827",
			"ETag":                                "\"A082B659EF78733A5A042FA253B1****\"",
			"Last-Modified":                       "Fri, 24 Feb 2023 09:41:59 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ObjectType, "Normal")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, *o.StorageClass, "Archive")
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.ContentType, "image/jpg")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.Restore, "ongoing-request=\"false\", expiry-date=\"Sun, 16 Apr 2017 08:12:33 GMT\"")
			assert.Equal(t, *o.ServerSideEncryption, "KMS")
			assert.Equal(t, *o.ServerSideDataEncryption, "SM4")
			assert.Equal(t, *o.SSEKMSKeyId, "9468da86-3509-4f8d-a61e-6eab1eac****")
		},
	},
}

func TestMockHeadObject_Success(t *testing.T) {
	for _, c := range testMockHeadObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockHeadObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *HeadObjectRequest
	CheckOutputFn  func(t *testing.T, o *HeadObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object", r.URL.String())
		},
		&HeadObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *HeadObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockHeadObject_Error(t *testing.T) {
	for _, c := range testMockHeadObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.HeadObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
			"Last-Modified":    "Fri, 24 Feb 2018 09:41:56 GMT",
			"Content-Length":   "344606",
			"ETag":             "\"fba9dede5f27731c9771645a3986****\"",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"fba9dede5f27731c9771645a3986****\"")
			assert.Equal(t, *o.LastModified, time.Date(2018, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(344606))
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":    "Fri, 24 Feb 2023 09:41:56 GMT",
			"Content-Length":   "481827",
			"ETag":             "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?objectMeta&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
		},
		&GetObjectMetaRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 56, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
			assert.Equal(t, *o.ETag, "\"A082B659EF78733A5A042FA253B1****\"")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id":       "534B371674E88A4D8906****",
			"Date":                   "Tue, 04 Dec 2018 15:56:38 GMT",
			"Last-Modified":          "Fri, 24 Feb 2023 09:41:59 GMT",
			"Content-Length":         "481827",
			"ETag":                   "\"A082B659EF78733A5A042FA253B1****\"",
			"x-oss-last-access-time": "Thu, 14 Oct 2021 11:49:05 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.LastModified, time.Date(2023, time.February, 24, 9, 41, 59, 0, time.UTC))
			assert.Equal(t, o.ContentLength, int64(481827))
			assert.Equal(t, *o.LastAccessTime, time.Date(2021, time.October, 14, 11, 49, 05, 0, time.UTC))
		},
	},
}

func TestMockGetObjectMeta_Success(t *testing.T) {
	for _, c := range testMockGetObjectMetaSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectMetaErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectMetaRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectMetaResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556E3AED11E553933CCDEDF",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPk5vU3VjaEtleTwvQ29kZT4KICA8TWVzc2FnZT5UaGUgc3BlY2lmaWVkIGtleSBkb2VzIG5vdCBleGlzdC48L01lc3NhZ2U+CiAgPFJlcXVlc3RJZD42NTU2RTNBRUQxMUU1NTM5MzNDQ0RFREY8L1JlcXVlc3RJZD4KICA8SG9zdElkPmRlbW8td2Fsa2VyLTY5NjEub3NzLWNuLWhhbmd6aG91LmFsaXl1bmNzLmNvbTwvSG9zdElkPgogIDxLZXk+d2Fsa2VyMmFzZGFzZGFzZC50eHQ8L0tleT4KICA8RUM+MDAyNi0wMDAwMDAwMTwvRUM+CiAgPFJlY29tbWVuZERvYz5odHRwczovL2FwaS5hbGl5dW4uY29tL3Ryb3VibGVzaG9vdD9xPTAwMjYtMDAwMDAwMDE8L1JlY29tbWVuZERvYz4KPC9FcnJvcj4K",
			"x-oss-ec":         "0026-00000001",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6556E3AED11E553933CCDEDF", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		304,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(304), serr.StatusCode)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E5536368607E8",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-err":        "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz4KPEVycm9yPgogIDxDb2RlPkludmFsaWRUYXJnZXRUeXBlPC9Db2RlPgogIDxNZXNzYWdlPlRoZSBzeW1ib2xpYydzIHRhcmdldCBmaWxlIHR5cGUgaXMgaW52YWxpZDwvTWVzc2FnZT4KICA8UmVxdWVzdElkPjY1NTZGRjVCRDExRTU1MzYzNjg2MDdFODwvUmVxdWVzdElkPgogIDxIb3N0SWQ+ZGVtby13YWxrZXItNjk2MS5vc3MtY24taGFuZ3pob3UuYWxpeXVuY3MuY29tPC9Ib3N0SWQ+CiAgPEVDPjAwMjYtMDAwMDAwMTE8L0VDPgogIDxSZWNvbW1lbmREb2M+aHR0cHM6Ly9hcGkuYWxpeXVuLmNvbS90cm91Ymxlc2hvb3Q/cT0wMDI2LTAwMDAwMDExPC9SZWNvbW1lbmREb2M+CjwvRXJyb3I+Cg==",
			"x-oss-ec":         "0026-00000011",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "HEAD", r.Method)
			assert.Equal(t, "/bucket/object?objectMeta", r.URL.String())
		},
		&GetObjectMetaRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectMetaResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidTargetType", serr.Code)
			assert.Equal(t, "6556FF5BD11E5536368607E8", serr.RequestID)
			assert.Equal(t, "The symbolic's target file type is invalid", serr.Message)
			assert.Equal(t, "0026-00000011", serr.EC)
		},
	},
}

func TestMockGetObjectMeta_Error(t *testing.T) {
	for _, c := range testMockGetObjectMetaErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectMeta(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		202,
		map[string]string{
			"X-Oss-Request-Id": "6555A936CA31DC333143****",
			"Date":             "Thu, 16 Nov 2023 05:31:34 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 202, o.StatusCode)
			assert.Equal(t, "202 Accepted", o.Status)
			assert.Equal(t, "6555A936CA31DC333143****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Thu, 16 Nov 2023 05:31:34 GMT", o.Headers.Get("Date"))

		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "5CAC3B40B7AEADE01700****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
			"x-oss-version-Id": "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?restore&versionId=CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy%2A%2A%2A%2A", strUrl)
			data, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "5CAC3B40B7AEADE01700****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionId, "CAEQNRiBgICb8o6D0BYiIDNlNzk5NGE2M2Y3ZjRhZTViYTAxZGE0ZTEyMWYy****")
		},
	},
	{
		200,
		map[string]string{
			"X-Oss-Request-Id": "534B371674E88A4D8906****",
			"Date":             "Tue, 04 Dec 2018 15:56:38 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())

			data, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(data), "<RestoreRequest><Days>2</Days><JobParameters><Tier>Standard</Tier></JobParameters></RestoreRequest>")
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RestoreRequest: &RestoreRequest{
				Days: int32(2),
				Tier: Ptr("Standard"),
			},
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Tue, 04 Dec 2018 15:56:38 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockRestoreObject_Success(t *testing.T) {
	for _, c := range testMockRestoreObjectSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockRestoreObjectErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *RestoreObjectRequest
	CheckOutputFn  func(t *testing.T, o *RestoreObjectResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6557176CD11E5535303C****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchKey</Code>
		<Message>The specified key does not exist.</Message>
		<RequestId>6557176CD11E5535303C****</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<Key>walker-not-.txt</Key>
		<EC>0026-00000001</EC>
		<RecommendDoc>https://api.aliyun.com/troubleshoot?q=0026-00000001</RecommendDoc>
		</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchKey", serr.Code)
			assert.Equal(t, "6557176CD11E5535303C****", serr.RequestID)
			assert.Equal(t, "The specified key does not exist.", serr.Message)
			assert.Equal(t, "0026-00000001", serr.EC)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6555AC764311A73931E0****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>OperationNotSupported</Code>
  <Message>The operation is not supported for this resource</Message>
  <RequestId>6555AC764311A73931E0****</RequestId>
  <HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
  <Detail>RestoreObject operation does not support this object storage class</Detail>
  <EC>0016-00000702</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "OperationNotSupported", serr.Code)
			assert.Equal(t, "6555AC764311A73931E0****", serr.RequestID)
			assert.Equal(t, "The operation is not supported for this resource", serr.Message)
			assert.Equal(t, "0016-00000702", serr.EC)
		},
	},
	{
		409,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "6556FF5BD11E55363686****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"x-oss-ec":         "0026-00000011",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>RestoreAlreadyInProgress</Code>
  <Message>The restore operation is in progress.</Message>
  <RequestId>6556FF5BD11E55363686****</RequestId>
  <HostId>10.101.XX.XX</HostId>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "/bucket/object?restore", r.URL.String())
		},
		&RestoreObjectRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *RestoreObjectResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(409), serr.StatusCode)
			assert.Equal(t, "RestoreAlreadyInProgress", serr.Code)
			assert.Equal(t, "6556FF5BD11E55363686****", serr.RequestID)
			assert.Equal(t, "The restore operation is in progress.", serr.Message)
		},
	},
}

func TestMockRestoreObject_Error(t *testing.T) {
	for _, c := range testMockRestoreObjectErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.RestoreObject(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPublicRead), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPublicRead,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0%2A%2A%2A%2A", strUrl)
			assert.Equal(t, string(ObjectACLPublicReadWrite), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			Acl:       ObjectACLPublicReadWrite,
			VersionId: Ptr("CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****"),
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get(HTTPHeaderContentType))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get(HeaderOssRequestID))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get(HTTPHeaderDate))

			assert.Equal(t, "CAEQMhiBgIC3rpSD0BYiIDBjYTk5MmIzN2JlNjQxZTFiNGIzM2E3OTliODA0****", *o.VersionId)
		},
	},
}

func TestMockPutObjectAcl_Success(t *testing.T) {
	for _, c := range testMockPutObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *PutObjectAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl", strUrl)
			assert.Equal(t, string(ObjectACLPrivate), r.Header.Get(HeaderOssObjectACL))
		},
		&PutObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			Acl:    ObjectACLPrivate,
		},
		func(t *testing.T, o *PutObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutObjectAcl_Error(t *testing.T) {
	for _, c := range testMockPutObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012****</ID>
        <DisplayName>user_example</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read", *o.ACL)
			assert.Equal(t, "0022012****", *o.Owner.ID)
			assert.Equal(t, "user_example", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>public-read-write</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "public-read-write", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
		},
	},
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"X-Oss-Version-Id": "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****",
		},
		[]byte(`<?xml version="1.0" ?>
<AccessControlPolicy>
    <Owner>
        <ID>0022012</ID>
        <DisplayName>0022012</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>private</Grant>
    </AccessControlList>
</AccessControlPolicy>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?acl&versionId=CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk%2A%2A%2A%2A", strUrl)
		},
		&GetObjectAclRequest{
			Bucket:    Ptr("bucket"),
			Key:       Ptr("object"),
			VersionId: Ptr("CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, "private", *o.ACL)
			assert.Equal(t, "0022012", *o.Owner.ID)
			assert.Equal(t, "0022012", *o.Owner.DisplayName)
			assert.Equal(t, "CAEQMhiBgMC1qpSD0BYiIGQ0ZmI5ZDEyYWVkNTQwMjBiNTliY2NjNmY3ZTVk****", *o.VersionId)
		},
	},
}

func TestMockGetObjectAcl_Success(t *testing.T) {
	for _, c := range testMockGetObjectAclSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetObjectAclErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetObjectAclRequest
	CheckOutputFn  func(t *testing.T, o *GetObjectAclResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", strUrl)
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket/object?acl", r.URL.String())
		},
		&GetObjectAclRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *GetObjectAclResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockGetObjectAcl_Error(t *testing.T) {
	for _, c := range testMockGetObjectAclErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetObjectAcl(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<InitiateMultipartUploadResult>
    <Bucket>oss-example</Bucket>
    <Key>multipart.data</Key>
    <UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "6551DBCF4311A7303980****",
			"Date":             "Mon, 13 Nov 2023 08:18:23 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<InitiateMultipartUploadResult>
		<Bucket>oss-example</Bucket>
		<Key>multipart.data</Key>
		<UploadId>0004B9894A22E5B1888A1E29F823****</UploadId>
		</InitiateMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
			assert.Equal(t, r.Header.Get("Cache-Control"), "no-cache")
			assert.Equal(t, r.Header.Get("Content-Disposition"), "attachment")
			assert.Equal(t, r.Header.Get("x-oss-meta-name"), "walker")
			assert.Equal(t, r.Header.Get("x-oss-meta-email"), "demo@aliyun.com")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption"), "KMS")
			assert.Equal(t, r.Header.Get("x-oss-server-side-data-encryption"), "SM4")
			assert.Equal(t, r.Header.Get("x-oss-server-side-encryption-key-id"), "9468da86-3509-4f8d-a61e-6eab1eac****")
			assert.Equal(t, r.Header.Get("x-oss-storage-class"), string(StorageClassStandard))
			assert.Equal(t, r.Header.Get("x-oss-forbid-overwrite"), "false")
			assert.Equal(t, r.Header.Get("Content-Encoding"), "utf-8")
			assert.Equal(t, r.Header.Get("Content-MD5"), "1B2M2Y8AsgTpgAmY7PhCfg==")
			assert.Equal(t, r.Header.Get("Expires"), "2022-10-12T00:00:00.000Z")
			assert.Equal(t, r.Header.Get("x-oss-tagging"), "TagA=B&TagC=D")
		},
		&InitiateMultipartUploadRequest{
			Bucket:                   Ptr("bucket"),
			Key:                      Ptr("object"),
			CacheControl:             Ptr("no-cache"),
			ContentDisposition:       Ptr("attachment"),
			ContentEncoding:          Ptr("utf-8"),
			Expires:                  Ptr("2022-10-12T00:00:00.000Z"),
			ForbidOverwrite:          Ptr("false"),
			ServerSideEncryption:     Ptr("KMS"),
			ServerSideDataEncryption: Ptr("SM4"),
			SSEKMSKeyId:              Ptr("9468da86-3509-4f8d-a61e-6eab1eac****"),
			StorageClass:             StorageClassStandard,
			Metadata: map[string]string{
				"name":  "walker",
				"email": "demo@aliyun.com",
			},
			Tagging: Ptr("TagA=B&TagC=D"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "multipart.data")
			assert.Equal(t, *o.UploadId, "0004B9894A22E5B1888A1E29F823****")
		},
	},
}

func TestMockInitiateMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockInitiateMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *InitiateMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *InitiateMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploads", strUrl)
		},
		&InitiateMultipartUploadRequest{
			Bucket: Ptr("bucket"),
			Key:    Ptr("object"),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("hi oss"),
			},
		},
		func(t *testing.T, o *InitiateMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockInitiateMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockInitiateMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.InitiateMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id":     "534B371674E88A4D8906****",
			"Date":                 "Fri, 24 Feb 2017 03:15:40 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F586****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Ph****",
			"x-oss-hash-crc64ecma": "316181249502703*****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
			assert.Equal(t, "bce8f3d48247c5d555bb5697bf277b35", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("upload part 1"),
			},
			ContentMD5: Ptr("bce8f3d48247c5d555bb5697bf277b35"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F586****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Ph****")
			assert.Equal(t, *o.HashCRC64, "316181249502703*****")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"ETag":                 "\"7265F4D211B56873A381D321F587****\"",
			"Content-MD5":          "1B2M2Y8AsgTpgAmY7Pp****",
			"x-oss-hash-crc64ecma": "316181249502704*****",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 2")
			assert.Equal(t, "f811b746eb3e256f97cb3a190d528353", r.Header.Get("Content-MD5"))
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(2),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("upload part 2"),
			},
			ContentMD5: Ptr("f811b746eb3e256f97cb3a190d528353"),
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"7265F4D211B56873A381D321F587****\"")
			assert.Equal(t, *o.ContentMD5, "1B2M2Y8AsgTpgAmY7Pp****")
			assert.Equal(t, *o.HashCRC64, "316181249502704*****")
		},
	},
}

func TestMockUploadPart_Success(t *testing.T) {
	for _, c := range testMockUploadPartSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("upload part 1"),
			},
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("upload part 1"),
			},
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "upload part 1")
		},
		&UploadPartRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			RequestCommon: RequestCommon{
				Body: strings.NewReader("upload part 1"),
			},
		},
		func(t *testing.T, o *UploadPartResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUploadPart_Error(t *testing.T) {
	for _, c := range testMockUploadPartErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPart(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))
		},
		&UploadPartCopyRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Source:     Ptr("/oss-src-bucket/" + url.QueryEscape("oss-src-object")),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":             "6551DBCF4311A7303980****",
			"Date":                         "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-copy-source-version-id": "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CopyPartResult>
    <LastModified>2014-07-17T06:27:54.000Z</LastModified>
    <ETag>"5B3C1A2E053D763E1B002CC607C5****"</ETag>
</CopyPartResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=2&uploadId=0004B9895DBBB6EC9&versionId=CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY", strUrl)
			assert.Equal(t, r.Header.Get(HeaderOssCopySource), "/oss-src-bucket/"+url.QueryEscape("oss-src-object"))

			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfMatch), "\"D41D8CD98F00B204E9800998ECF8****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfNoneMatch), "\"D41D8CD98F00B204E9800998ECF9****\"")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfModifiedSince), "Fri, 13 Nov 2023 14:47:53 GMT")
			assert.Equal(t, r.Header.Get(HeaderOssCopySourceIfUnmodifiedSince), "Fri, 13 Nov 2015 14:47:53 GMT")
		},
		&UploadPartCopyRequest{
			Bucket:            Ptr("bucket"),
			Key:               Ptr("object"),
			UploadId:          Ptr("0004B9895DBBB6EC9"),
			Source:            Ptr("/oss-src-bucket/" + url.QueryEscape("oss-src-object")),
			PartNumber:        int32(2),
			IfMatch:           Ptr("\"D41D8CD98F00B204E9800998ECF8****\""),
			IfNoneMatch:       Ptr("\"D41D8CD98F00B204E9800998ECF9****\""),
			IfModifiedSince:   Ptr("Fri, 13 Nov 2023 14:47:53 GMT"),
			IfUnmodifiedSince: Ptr("Fri, 13 Nov 2015 14:47:53 GMT"),
			VersionId:         Ptr("CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY"),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"5B3C1A2E053D763E1B002CC607C5****\"")
			assert.Equal(t, *o.LastModified, time.Date(2014, time.July, 17, 6, 27, 54, 0, time.UTC))
			assert.Equal(t, *o.VersionId, "CAEQNhiBgM0BYiIDc4MGZjZGI2OTBjOTRmNTE5NmU5NmFhZjhjYmY")
		},
	},
}

func TestMockUploadPartCopy_Success(t *testing.T) {
	for _, c := range testMockUploadPartCopySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockUploadPartCopyErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *UploadPartCopyRequest
	CheckOutputFn  func(t *testing.T, o *UploadPartCopyResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Source:     Ptr("/oss-src-bucket/" + url.QueryEscape("oss-src-object")),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Source:     Ptr("/oss-src-bucket/" + url.QueryEscape("oss-src-object")),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?partNumber=1&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&UploadPartCopyRequest{
			Bucket:     Ptr("bucket"),
			Key:        Ptr("object"),
			UploadId:   Ptr("0004B9895DBBB6EC9"),
			PartNumber: int32(1),
			Source:     Ptr("/oss-src-bucket/" + url.QueryEscape("oss-src-object")),
		},
		func(t *testing.T, o *UploadPartCopyResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockUploadPartCopy_Error(t *testing.T) {
	for _, c := range testMockUploadPartCopyErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.UploadPartCopy(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<CompleteMultipartUploadResult>
  <EncodingType>url</EncodingType>
  <Location>http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data</Location>
  <Bucket>oss-example</Bucket>
  <Key>demo%2Fmultipart.data</Key>
  <ETag>"097DE458AD02B5F89F9D0530231876****"</ETag>
</CompleteMultipartUploadResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), `<CompleteMultipartUpload><Part><PartNumber>1</PartNumber><ETag>&#34;8EFDA8BE206636A695359836FE0A****&#34;</ETag></Part><Part><PartNumber>2</PartNumber><ETag>&#34;8C315065167132444177411FDA14****&#34;</ETag></Part><Part><PartNumber>3</PartNumber><ETag>&#34;3349DC700140D7F86A0784842780****&#34;</ETag></Part></CompleteMultipartUpload>`)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Part: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.ETag, "\"097DE458AD02B5F89F9D0530231876****\"")
			assert.Equal(t, *o.Location, "http://oss-example.oss-cn-hangzhou.aliyuncs.com/multipart.data")
			assert.Equal(t, *o.EncodingType, "url")
			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.Key, "demo/multipart.data")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id":     "6551DBCF4311A7303980****",
			"Date":                 "Mon, 13 Nov 2023 08:18:23 GMT",
			"x-oss-version-id":     "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****",
			"Content-Type":         "application/json",
			"x-oss-hash-crc64ecma": "1206617243528768****",
		},
		[]byte(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
			assert.Equal(t, "false", r.Header.Get(HeaderOssForbidOverWrite))
			assert.Equal(t, "yes", r.Header.Get("x-oss-complete-all"))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`)), r.Header.Get(HeaderOssCallback))
			assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`)), r.Header.Get(HeaderOssCallbackVar))
		},
		&CompleteMultipartUploadRequest{
			Bucket:          Ptr("bucket"),
			Key:             Ptr("object"),
			UploadId:        Ptr("0004B9895DBBB6EC9"),
			ForbidOverwrite: Ptr("false"),
			CompleteAll:     Ptr("yes"),
			Callback:        Ptr(base64.StdEncoding.EncodeToString([]byte(`{"callbackUrl":"www.aliyuncs.com", "callbackBody":"filename=${object}&size=${size}&mimeType=${mimeType}&x=${x:a}&b=${x:b}"}`))),
			CallbackVar:     Ptr(base64.StdEncoding.EncodeToString([]byte(`{"x:a":"a", "x:b":"b"}`))),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "6551DBCF4311A7303980****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Mon, 13 Nov 2023 08:18:23 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get(HTTPHeaderContentType), "application/json")
			assert.Equal(t, o.Body, io.NopCloser(strings.NewReader(`{"filename":"oss-obj.txt","size":"100","mimeType":"","x":"a","b":"b"}`)))
			assert.Equal(t, *o.VersionId, "CAEQMxiBgMC0vs6D0BYiIGJiZWRjOTRjNTg0NzQ1MTRiN2Y1OTYxMTdkYjQ0****")
		},
	},
}

func TestMockCompleteMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockCompleteMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *CompleteMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *CompleteMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D94CCD11E55313348****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>MalformedXML</Code>
  <Message>The XML you provided was not well-formed or did not validate against our published schema.</Message>
  <RequestId>655D94CCD11E55313348****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <EC>0042-00000205</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000205</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "MalformedXML", serr.Code)
			assert.Equal(t, "The XML you provided was not well-formed or did not validate against our published schema.", serr.Message)
			assert.Equal(t, "655D94CCD11E55313348****", serr.RequestID)
			assert.Equal(t, "0042-00000205", serr.EC)
		},
	},
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "655D9598CA31DC313626****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>Should not speficy both complete all header and http body.</Message>
  <RequestId>655D9598CA31DC313626****</RequestId>
  <HostId>demo-walker-6961.oss-cn-hangzhou.aliyuncs.com</HostId>
  <ArgumentName>x-oss-complete-all</ArgumentName>
  <ArgumentValue>yes</ArgumentValue>
  <EC>0042-00000216</EC>
  <RecommendDoc>https://api.aliyun.com/troubleshoot?q=0042-00000216</RecommendDoc>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:      Ptr("bucket"),
			Key:         Ptr("object"),
			UploadId:    Ptr("0004B9895DBBB6EC9"),
			CompleteAll: Ptr("yes"),
			CompleteMultipartUpload: &CompleteMultipartUpload{
				Part: []UploadPart{
					{PartNumber: int32(3), ETag: Ptr("\"3349DC700140D7F86A0784842780****\"")},
					{PartNumber: int32(1), ETag: Ptr("\"8EFDA8BE206636A695359836FE0A****\"")},
					{PartNumber: int32(2), ETag: Ptr("\"8C315065167132444177411FDA14****\"")},
				},
			},
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "Should not speficy both complete all header and http body.", serr.Message)
			assert.Equal(t, "655D9598CA31DC313626****", serr.RequestID)
			assert.Equal(t, "0042-00000216", serr.EC)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B9895DBBB6EC9", strUrl)
		},
		&CompleteMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6EC9"),
		},
		func(t *testing.T, o *CompleteMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockCompleteMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockCompleteMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.CompleteMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		nil,
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockAbortMultipartUpload_Success(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockAbortMultipartUploadErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *AbortMultipartUploadRequest
	CheckOutputFn  func(t *testing.T, o *AbortMultipartUploadResult, err error)
}{
	{
		400,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>InvalidArgument</Code>
  <Message>no such bucket access control exists</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>***-test.example.com</HostId>
  <ArgumentName>x-oss-acl</ArgumentName>
  <ArgumentValue>error-acl</ArgumentValue>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(400), serr.StatusCode)
			assert.Equal(t, "InvalidArgument", serr.Code)
			assert.Equal(t, "no such bucket access control exists", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "/bucket/object?uploadId=0004B9895DBBB6E", r.URL.String())
		},
		&AbortMultipartUploadRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B9895DBBB6E"),
		},
		func(t *testing.T, o *AbortMultipartUploadResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockAbortMultipartUpload_Error(t *testing.T) {
	for _, c := range testMockAbortMultipartUploadErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.AbortMultipartUpload(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
    <Bucket>oss-example</Bucket>
    <KeyMarker></KeyMarker>
    <UploadIdMarker></UploadIdMarker>
    <NextKeyMarker>oss.avi</NextKeyMarker>
    <NextUploadIdMarker>0004B99B8E707874FC2D692FA5D77D3F</NextUploadIdMarker>
    <Delimiter></Delimiter>
    <Prefix></Prefix>
    <MaxUploads>1000</MaxUploads>
    <IsTruncated>false</IsTruncated>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF518A1FE585B0C9360DC4C8</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>multipart.data</Key>
        <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
        <Initiated>2012-02-23T04:18:23.000Z</Initiated>
    </Upload>
    <Upload>
        <Key>oss.avi</Key>
        <UploadId>0004B99B8E707874FC2D692FA5D7****</UploadId>
        <Initiated>2012-02-23T06:14:27.000Z</Initiated>
    </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "0004B99B8E707874FC2D692FA5D77D3F")
			assert.Equal(t, *o.Delimiter, "")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 3)
			assert.Equal(t, *o.Uploads[0].Key, "multipart.data")
			assert.Equal(t, *o.Uploads[0].UploadId, "0004B999EF518A1FE585B0C9360DC4C8")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2012, time.February, 23, 4, 18, 23, 0, time.UTC))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
  <EncodingType>url</EncodingType>
  <Bucket>oss-example</Bucket>
  <KeyMarker></KeyMarker>
  <UploadIdMarker></UploadIdMarker>
  <NextKeyMarker>oss.avi</NextKeyMarker>
  <NextUploadIdMarker>89F0105AA66942638E35300618DF****</NextUploadIdMarker>
  <Delimiter>/</Delimiter>
  <Prefix>pre</Prefix>
  <MaxUploads>1000</MaxUploads>
  <IsTruncated>false</IsTruncated>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>0214A87687F040F1BA4D83AB17C9****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:45:57.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>3AE2ED7A60E04AFE9A5287055D37****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:03:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>47E0E90F5DCB4AD5B3C4CD886CB0****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:02:11.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>A89E0E28E2E948A1BFF6FD5CDAFF****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T06:57:03.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>B18E1DCDB6964F5CB197F5F6B26A****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:02.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-22T05:42:33.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>5209986C3A96486EA16B9C52C160****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:34:47.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>63B652FA2C1342DCB3CCCC86D748****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:28:46.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>6F67B34BCA3C481F887D73508A07****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:32:12.000Z</Initiated>
  </Upload>
  <Upload>
    <Key>walker-dest.txt</Key>
    <UploadId>89F0105AA66942638E35300618D****</UploadId>
    <StorageClass>Standard</StorageClass>
    <Initiated>2023-11-21T08:37:53.000Z</Initiated>
  </Upload>
</ListMultipartUploadsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-uploads=10&prefix=pre&upload-id-marker&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket:         Ptr("bucket"),
			Delimiter:      Ptr("/"),
			Prefix:         Ptr("pre"),
			EncodingType:   Ptr("url"),
			KeyMarker:      Ptr(""),
			MaxUploads:     int32(10),
			UploadIdMarker: Ptr(""),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Bucket, "oss-example")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.UploadIdMarker, "")
			assert.Equal(t, *o.NextKeyMarker, "oss.avi")
			assert.Equal(t, *o.NextUploadIdMarker, "89F0105AA66942638E35300618DF****")
			assert.Equal(t, *o.Delimiter, "/")
			assert.Equal(t, *o.Prefix, "pre")
			assert.Equal(t, o.MaxUploads, int32(1000))
			assert.Equal(t, o.IsTruncated, false)
			assert.Len(t, o.Uploads, 10)
			assert.Equal(t, *o.Uploads[0].Key, "demo/gp-\f\n\v")
			assert.Equal(t, *o.Uploads[0].UploadId, "0214A87687F040F1BA4D83AB17C9****")
			assert.Equal(t, *o.Uploads[0].Initiated, time.Date(2023, time.November, 22, 5, 45, 57, 0, time.UTC))
		},
	},
}

func TestMockListMultipartUploads_Success(t *testing.T) {
	for _, c := range testMockListMultipartUploadsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListMultipartUploadsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListMultipartUploadsRequest
	CheckOutputFn  func(t *testing.T, o *ListMultipartUploadsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&uploads", strUrl)
		},
		&ListMultipartUploadsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListMultipartUploadsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockListMultipartUploads_Error(t *testing.T) {
	for _, c := range testMockListMultipartUploadsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListMultipartUploads(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
    <Bucket>bucket</Bucket>
    <Key>object</Key>
    <UploadId>0004B999EF5A239BB9138C6227D6****</UploadId>
    <NextPartNumberMarker>5</NextPartNumberMarker>
    <MaxParts>1000</MaxParts>
    <IsTruncated>false</IsTruncated>
    <Part>
        <PartNumber>1</PartNumber>
        <LastModified>2012-02-23T07:01:34.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>2</PartNumber>
        <LastModified>2012-02-23T07:01:12.000Z</LastModified>
        <ETag>"3349DC700140D7F86A0784842780****"</ETag>
        <Size>6291456</Size>
    </Part>
    <Part>
        <PartNumber>5</PartNumber>
        <LastModified>2012-02-23T07:02:03.000Z</LastModified>
        <ETag>"7265F4D211B56873A381D321F586****"</ETag>
        <Size>1024</Size>
    </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			assert.Equal(t, *o.Key, "object")
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(5))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 3)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"3349DC700140D7F86A0784842780****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2012, time.February, 23, 7, 1, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(6291456))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
  <EncodingType>url</EncodingType>
  <Bucket>bucket</Bucket>
  <Key>demo%2Fgp-%0C%0A%0B</Key>
  <UploadId>D4E111D4EA834F3ABCE4877B2779****</UploadId>
  <StorageClass>Standard</StorageClass>
  <PartNumberMarker>0</PartNumberMarker>
  <NextPartNumberMarker>1</NextPartNumberMarker>
  <MaxParts>1000</MaxParts>
  <IsTruncated>false</IsTruncated>
  <Part>
    <PartNumber>1</PartNumber>
    <LastModified>2023-11-22T05:42:34.000Z</LastModified>
    <ETag>"CF3F46D505093571E916FCDD4967****"</ETag>
    <HashCrc64ecma>12066172435287683848</HashCrc64ecma>
    <Size>96316</Size>
  </Part>
</ListPartsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/demo/gp-%0C%0A%0B?encoding-type=url&uploadId=D4E111D4EA834F3ABCE4877B2779%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("demo/gp-\f\n\v"),
			UploadId: Ptr("D4E111D4EA834F3ABCE4877B2779****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.Bucket, "bucket")
			key, _ := url.QueryUnescape("demo%2Fgp-%0C%0A%0B")
			assert.Equal(t, *o.Key, key)
			assert.Empty(t, o.PartNumberMarker)
			assert.Equal(t, o.NextPartNumberMarker, int32(1))
			assert.Equal(t, o.IsTruncated, false)
			assert.Equal(t, o.MaxParts, int32(1000))
			assert.Len(t, o.Parts, 1)
			assert.Equal(t, o.Parts[0].PartNumber, int32(1))
			assert.Equal(t, *o.Parts[0].ETag, "\"CF3F46D505093571E916FCDD4967****\"")
			assert.Equal(t, *o.Parts[0].LastModified, time.Date(2023, time.November, 22, 5, 42, 34, 0, time.UTC))
			assert.Equal(t, o.Parts[0].Size, int64(96316))
			assert.Equal(t, *o.Parts[0].HashCRC64, "12066172435287683848")
		},
	},
}

func TestMockListParts_Success(t *testing.T) {
	for _, c := range testMockListPartsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListPartsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListPartsRequest
	CheckOutputFn  func(t *testing.T, o *ListPartsResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/object?encoding-type=url&uploadId=0004B999EF5A239BB9138C6227D6%2A%2A%2A%2A", strUrl)
		},
		&ListPartsRequest{
			Bucket:   Ptr("bucket"),
			Key:      Ptr("object"),
			UploadId: Ptr("0004B999EF5A239BB9138C6227D6****"),
		},
		func(t *testing.T, o *ListPartsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "0015-00000101", serr.EC)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
}

func TestMockListParts_Error(t *testing.T) {
	for _, c := range testMockListPartsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListParts(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Suspended</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionSuspended,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockPutBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketVersioningResult, err error)
}{
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, string(body), "<VersioningConfiguration><Status>Enabled</Status></VersioningConfiguration>")
		},
		&PutBucketVersioningRequest{
			Bucket: Ptr("bucket"),
			VersioningConfiguration: &VersioningConfiguration{
				Status: VersionEnabled,
			},
		},
		func(t *testing.T, o *PutBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockPutBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockPutBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.PutBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Nil(t, o.VersionStatus)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Enabled</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Enabled")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
<Status>Suspended</Status>
</VersioningConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.VersionStatus, "Suspended")
		},
	},
}

func TestMockGetBucketVersioning_Success(t *testing.T) {
	for _, c := range testMockGetBucketVersioningSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketVersioningErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketVersioningRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketVersioningResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "/bucket?versioning", r.URL.String())
		},
		&GetBucketVersioningRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketVersioningResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockGetBucketVersioning_Error(t *testing.T) {
	for _, c := range testMockGetBucketVersioningErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketVersioning(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsSuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2F</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>20</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-keys=20&prefix=demo%2F&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			Delimiter:       Ptr("/"),
			Prefix:          Ptr("demo/"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			MaxKeys:         int32(20),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(20))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 0)
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
    <Name>examplebucket-1250000000</Name>
    <Prefix/>
    <KeyMarker/>
    <VersionIdMarker/>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Version>
        <Key>example-object-1.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-05T12:03:10.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B669CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-2.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-09T12:03:09.000Z</LastModified>
        <ETag>5B3C1A2E053D763E1B002CC607C5A0FE1****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
    <Version>
        <Key>example-object-3.jpg</Key>
        <VersionId/>
        <IsLatest>true</IsLatest>
        <LastModified>2019-08-10T12:03:08.000Z</LastModified>
        <ETag>4B3F1A2E053D763E1B002CC607C5AGTRF****</ETag>
        <Size>20</Size>
        <StorageClass>STANDARD</StorageClass>
        <Owner>
            <ID>1250000000</ID>
            <DisplayName>1250000000</DisplayName>
        </Owner>
    </Version>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "examplebucket-1250000000")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(1000))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			assert.Equal(t, *o.ObjectVersions[0].Key, "example-object-1.jpg")
			assert.Empty(t, *o.ObjectVersions[1].VersionId)
			assert.True(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "5B3C1A2E053D763E1B002CC607C5A0FE1****")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(20))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "1250000000")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "1250000000")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
			"Content-Type":     "application/xml",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<ListVersionsResult>
  <Name>demo-bucket</Name>
  <Prefix>demo%2Fgp-</Prefix>
  <KeyMarker></KeyMarker>
  <VersionIdMarker></VersionIdMarker>
  <MaxKeys>5</MaxKeys>
  <Delimiter>%2F</Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>false</IsTruncated>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgIDAj.jV3xgiIGFjMDI5ZTRmNGNiODQ0NjE4MDFhODM0Y2UxNTI3****</VersionId>
    <IsLatest>true</IsLatest>
    <LastModified>2023-11-22T05:15:05.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:11:25.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
  <Version>
    <Key>demo%2Fgp-%0C%0A%0B</Key>
    <VersionId>CAEQHxiBgICCuNrV3xgiIDI2YzMyYTBhM2U1ZTQwNjI4OWQ4OTllZGJiNGIz****</VersionId>
    <IsLatest>false</IsLatest>
    <LastModified>2023-11-22T05:07:37.000Z</LastModified>
    <ETag>"29B94424BC241D80B0AF488A4E4B86AF-1"</ETag>
    <Type>Multipart</Type>
    <Size>96316</Size>
    <StorageClass>Standard</StorageClass>
    <Owner>
      <ID>150692521021****</ID>
      <DisplayName>150692521021****</DisplayName>
    </Owner>
  </Version>
</ListVersionsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?delimiter=%2F&encoding-type=url&key-marker&max-keys=5&prefix=demo%2Fgp-&version-id-marker&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket:          Ptr("bucket"),
			KeyMarker:       Ptr(""),
			VersionIdMarker: Ptr(""),
			Delimiter:       Ptr("/"),
			MaxKeys:         int32(5),
			Prefix:          Ptr("demo/gp-"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, o.Headers.Get("Content-Type"), "application/xml")
			assert.Equal(t, *o.Name, "demo-bucket")
			prefix, _ := url.QueryUnescape(*o.Prefix)
			assert.Equal(t, *o.Prefix, prefix)
			assert.Equal(t, *o.KeyMarker, "")
			assert.Equal(t, *o.VersionIdMarker, "")
			assert.Equal(t, o.MaxKeys, int32(5))
			assert.False(t, o.IsTruncated)
			assert.Len(t, o.ObjectVersions, 3)
			key, _ := url.QueryUnescape(*o.ObjectVersions[0].Key)
			assert.Equal(t, *o.ObjectVersions[0].Key, key)
			assert.Equal(t, *o.ObjectVersions[1].VersionId, "CAEQHxiBgMDYseHV3xgiIDg2Mzk0Zjg3MjQ0MTRhM2FiMzgxOGY1NjdmN2Rk****")
			assert.False(t, o.ObjectVersions[2].IsLatest)
			assert.NotEmpty(t, *o.ObjectVersions[0].LastModified)
			assert.Equal(t, *o.ObjectVersions[1].ETag, "\"29B94424BC241D80B0AF488A4E4B86AF-1\"")
			assert.Equal(t, o.ObjectVersions[2].Size, int64(96316))
			assert.Equal(t, *o.ObjectVersions[2].Owner.ID, "150692521021****")
			assert.Equal(t, *o.ObjectVersions[2].Owner.DisplayName, "150692521021****")
		},
	},
}

func TestMockListObjectVersions_Success(t *testing.T) {
	for _, c := range testMockListObjectVersionsSuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListObjectVersionsErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListObjectVersionsRequest
	CheckOutputFn  func(t *testing.T, o *ListObjectVersionsResult, err error)
}{
	{
		404,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D9175B6FC201293AD****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>NoSuchBucket</Code>
  <Message>The specified bucket does not exist.</Message>
  <RequestId>5C3D9175B6FC201293AD****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0015-00000101</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(404), serr.StatusCode)
			assert.Equal(t, "NoSuchBucket", serr.Code)
			assert.Equal(t, "The specified bucket does not exist.", serr.Message)
			assert.Equal(t, "5C3D9175B6FC201293AD****", serr.RequestID)
		},
	},
	{
		403,
		map[string]string{
			"Content-Type":     "application/xml",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Error>
  <Code>UserDisable</Code>
  <Message>UserDisable</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000801</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket?encoding-type=url&versions", strUrl)
		},
		&ListObjectVersionsRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectVersionsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "UserDisable", serr.Code)
			assert.Equal(t, "UserDisable", serr.Message)
			assert.Equal(t, "0003-00000801", serr.EC)
			assert.Equal(t, "5C3D8D2A0ACA54D87B43****", serr.RequestID)
		},
	},
}

func TestMockListObjectVersions_Error(t *testing.T) {
	for _, c := range testMockListObjectVersionsErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListObjectVersions(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}
