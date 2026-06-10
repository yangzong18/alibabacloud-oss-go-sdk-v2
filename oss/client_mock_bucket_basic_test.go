package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/bucket/", r.URL.String())
			requestBody, err := io.ReadAll(r.Body)
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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F&tagging=%22k%22%3A%22v%22", strUrl)
		},
		&ListBucketsRequest{
			Marker:          Ptr(""),
			MaxKeys:         10,
			Prefix:          Ptr("/"),
			ResourceGroupId: Ptr("rg-aek27tc********"),
			Tagging:         Ptr("\"k\":\"v\""),
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
			assert.Equal(t, "/?marker&max-keys=10&prefix=%2F&tag-key=k&tag-value=v", strUrl)
		},
		&ListBucketsRequest{
			Marker:          Ptr(""),
			MaxKeys:         10,
			Prefix:          Ptr("/"),
			ResourceGroupId: Ptr("rg-aek27tc********"),
			TagKey:          Ptr("k"),
			TagValue:        Ptr("v"),
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
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/", r.URL.String())
		},
		&ListBucketsRequest{},
		func(t *testing.T, o *ListBucketsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListBuckets fail")
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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/bucket/", r.URL.String())
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
			assert.Equal(t, "/bucket/?encoding-type=url", r.URL.String())
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
			RequestPayer: Ptr("requester"),
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
  <Marker>test1.txt</Marker>
  <MaxKeys>1</MaxKeys>
  <Delimiter></Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>true</IsTruncated>
  <NextMarker>test100.txt</NextMarker>
  <CommonPrefixes/>
  <Contents>
        <Key></Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, *o.Marker, "test1.txt")
			assert.Equal(t, *o.Delimiter, "")
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, len(o.CommonPrefixes), 1)
			assert.Equal(t, o.MaxKeys, int32(1))
			assert.Equal(t, len(o.Contents), 1)
			assert.Equal(t, *o.Contents[0].Key, "")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Empty(t, o.Contents[0].RestoreInfo)
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
			assert.Equal(t, "/bucket/?encoding-type=url", r.URL.String())
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
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
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
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
			assert.Contains(t, err.Error(), "execute ListObjects fail")
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
  <Code>AccessDenied</Code>
  <Message>Access denied for requester pay bucket</Message>
  <RequestId>5C3D8D2A0ACA54D87B43****</RequestId>
  <HostId>test.oss-cn-hangzhou.aliyuncs.com</HostId>
  <BucketName>test</BucketName>
  <EC>0003-00000703</EC>
</Error>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=URL&marker&max-keys=3&prefix", strUrl)
		},
		&ListObjectsRequest{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			Marker:       Ptr(""),
			MaxKeys:      int32(3),
			Prefix:       Ptr(""),
			EncodingType: Ptr("URL"),
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			var serr *ServiceError
			errors.As(err, &serr)
			assert.NotNil(t, serr)
			assert.Equal(t, int(403), serr.StatusCode)
			assert.Equal(t, "AccessDenied", serr.Code)
			assert.Equal(t, "Access denied for requester pay bucket", serr.Message)
			assert.Equal(t, "0003-00000703", serr.EC)
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
	Request        *ListObjectsV2Request
	CheckOutputFn  func(t *testing.T, o *ListObjectsV2Result, err error)
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
			assert.Equal(t, "/bucket/?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsV2Request{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsV2Request{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsV2Request{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
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
  <StartAfter>test1.txt</StartAfter>
  <MaxKeys>1</MaxKeys>
  <Delimiter></Delimiter>
  <EncodingType>url</EncodingType>
  <IsTruncated>true</IsTruncated>
  <NextMarker>test100.txt</NextMarker>
  <CommonPrefixes/>
  <Contents>
        <Key></Key>
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
</ListBucketResult>`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
			assert.Equal(t, r.Header.Get("x-oss-request-payer"), "requester")
		},
		&ListObjectsV2Request{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
			RequestPayer: Ptr("requester"),
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "application/xml", o.Headers.Get("Content-Type"))
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

			assert.Equal(t, *o.Name, "examplebucket")
			assert.Equal(t, *o.Prefix, "")
			assert.Equal(t, *o.StartAfter, "test1.txt")
			assert.Equal(t, *o.Delimiter, "")
			assert.Equal(t, o.IsTruncated, true)
			assert.Equal(t, len(o.CommonPrefixes), 1)
			assert.Equal(t, o.MaxKeys, int32(1))
			assert.Equal(t, len(o.Contents), 1)
			assert.Equal(t, *o.Contents[0].Key, "")
			assert.Equal(t, *o.Contents[0].Type, "Normal")
			assert.Equal(t, *o.Contents[0].Owner.ID, "0022012****")
			assert.Equal(t, *o.Contents[0].Owner.DisplayName, "user-example")
			assert.Empty(t, o.Contents[0].RestoreInfo)
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
	Request        *ListObjectsV2Request
	CheckOutputFn  func(t *testing.T, o *ListObjectsV2Result, err error)
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
			assert.Equal(t, "/bucket/?encoding-type=url&list-type=2", strUrl)
		},
		&ListObjectsV2Request{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
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
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsV2Request{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
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
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?delimiter=%2F&encoding-type=url&fetch-owner=true&list-type=2&max-keys=3&prefix=a%2F&start-after=b", strUrl)
		},
		&ListObjectsV2Request{
			Bucket:       Ptr("bucket"),
			Delimiter:    Ptr("/"),
			StartAfter:   Ptr("b"),
			MaxKeys:      int32(3),
			Prefix:       Ptr("a/"),
			EncodingType: Ptr("url"),
			FetchOwner:   true,
		},
		func(t *testing.T, o *ListObjectsV2Result, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListObjectsV2 fail")
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
	<DataRedundancyType>LRS</DataRedundancyType>
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
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
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
			assert.Equal(t, *o.BucketInfo.DataRedundancyType, "LRS")
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
			assert.Nil(t, o.BucketInfo.BlockPublicAccess)
			assert.Nil(t, o.BucketInfo.Comment)
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
    <BlockPublicAccess>true</BlockPublicAccess>
    <Comment>test</Comment>
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
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
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

			assert.True(t, *o.BucketInfo.BlockPublicAccess)
			assert.Equal(t, *o.BucketInfo.Comment, "test")
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
    <BlockPublicAccess>false</BlockPublicAccess>
    <Comment>test</Comment>
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
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
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

			assert.False(t, *o.BucketInfo.BlockPublicAccess)
			assert.Equal(t, *o.BucketInfo.Comment, "test")
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
			assert.Equal(t, "/bucket/?bucketInfo", r.URL.String())
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
			assert.Equal(t, "/bucket/?bucketInfo", strUrl)
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
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?bucketInfo", strUrl)
		},
		&GetBucketInfoRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketInfoResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketInfo fail")
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
			assert.Equal(t, "/bucket/?location", r.URL.String())
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
			assert.Equal(t, "/bucket/?location", r.URL.String())
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
			assert.Equal(t, "/bucket/?location", r.URL.String())
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
			assert.Equal(t, "/bucket/?location", strUrl)
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
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?location", strUrl)
		},
		&GetBucketLocationRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketLocationResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketLocation fail")
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
  <DeepColdArchiveStorage>2340340840</DeepColdArchiveStorage>
  <DeepColdArchiveRealStorage>2340340840</DeepColdArchiveRealStorage>
  <DeepColdArchiveObjectCount>2</DeepColdArchiveObjectCount>
  <MultipartPartCount>1</MultipartPartCount>
  <DeleteMarkerCount>6276</DeleteMarkerCount>
</BucketStat>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "/bucket/?stat", r.URL.String())
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
			assert.Equal(t, int64(2340340840), o.DeepColdArchiveStorage)
			assert.Equal(t, int64(2340340840), o.DeepColdArchiveRealStorage)
			assert.Equal(t, int64(2), o.DeepColdArchiveObjectCount)
			assert.Equal(t, int64(1), o.MultipartPartCount)
			assert.Equal(t, int64(6276), o.DeleteMarkerCount)
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
			assert.Equal(t, "/bucket/?stat", r.URL.String())
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
			assert.Equal(t, "/bucket/?stat", strUrl)
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
	{
		200,
		map[string]string{
			"Content-Type":     "application/text",
			"x-oss-request-id": "5C3D8D2A0ACA54D87B43****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`StrField1>StrField1</StrField1><StrField2>StrField2<`),
		func(t *testing.T, r *http.Request) {
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?stat", strUrl)
		},
		&GetBucketStatRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *GetBucketStatResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketStat fail")
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


