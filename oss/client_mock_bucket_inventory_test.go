package oss

import (
	"testing"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/stretchr/testify/assert"
)

var testMockPutBucketInventorySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketInventoryResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions></InventoryConfiguration>")
		},
		&PutBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
			InventoryConfiguration: &InventoryConfiguration{
				Id:        Ptr("report1"),
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
						AccountId: Ptr("1000000000000000"),
						RoleArn:   Ptr("acs:ram::1000000000000000:role/AliyunOSSRole"),
						Bucket:    Ptr("acs:oss:::destination-bucket"),
					},
				},
				Schedule: &InventorySchedule{
					Frequency: InventoryFrequencyDaily,
				},
				IncludedObjectVersions: Ptr("All"),
			},
		},
		func(t *testing.T, o *PutBucketInventoryResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket><Prefix>prefix1</Prefix><Encryption><SSE-KMS><KeyId>keyId</KeyId></SSE-KMS></Encryption></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field><Field>EncryptionStatus</Field></OptionalFields></InventoryConfiguration>")
		},
		&PutBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
			InventoryConfiguration: &InventoryConfiguration{
				Id:        Ptr("report1"),
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
						AccountId: Ptr("1000000000000000"),
						RoleArn:   Ptr("acs:ram::1000000000000000:role/AliyunOSSRole"),
						Bucket:    Ptr("acs:oss:::destination-bucket"),
						Prefix:    Ptr("prefix1"),
						Encryption: &InventoryEncryption{
							SseKms: &SSEKMS{
								Ptr("keyId"),
							},
						},
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
		},
		func(t *testing.T, o *PutBucketInventoryResult, err error) {
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket><Prefix>prefix1</Prefix><Encryption><SSE-KMS><KeyId>keyId</KeyId></SSE-KMS></Encryption></OSSBucketDestination></Destination><Schedule><Frequency>Monthly</Frequency><DayOfMonth>15</DayOfMonth></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field><Field>EncryptionStatus</Field><Field>ObjectAcl</Field><Field>TaggingCount</Field><Field>ObjectType</Field><Field>Crc64</Field></OptionalFields><IncrementalInventory><IsEnabled>true</IsEnabled><Schedule><Frequency>600</Frequency></Schedule><OptionalFields><Field>SequenceNumber</Field><Field>RecordType</Field><Field>RecordTimestamp</Field><Field>Requester</Field><Field>SourceIp</Field><Field>RequestId</Field><Field>Size</Field><Field>StorageClass</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>IsMultipartUploaded</Field><Field>ObjectType</Field><Field>ObjectAcl</Field><Field>Crc64</Field><Field>EncryptionStatus</Field></OptionalFields></IncrementalInventory></InventoryConfiguration>")
		},
		&PutBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
			InventoryConfiguration: &InventoryConfiguration{
				Id:        Ptr("report1"),
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
						AccountId: Ptr("1000000000000000"),
						RoleArn:   Ptr("acs:ram::1000000000000000:role/AliyunOSSRole"),
						Bucket:    Ptr("acs:oss:::destination-bucket"),
						Prefix:    Ptr("prefix1"),
						Encryption: &InventoryEncryption{
							SseKms: &SSEKMS{
								Ptr("keyId"),
							},
						},
					},
				},
				Schedule: &InventorySchedule{
					Frequency:  InventoryFrequencyMonthly,
					DayOfMonth: Ptr(int(15)),
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
						InventoryOptionalFieldObjectAcl,
						InventoryOptionalFieldTaggingCount,
						InventoryOptionalFieldObjectType,
						InventoryOptionalFieldCRC64,
					},
				},
				IncrementalInventory: &IncrementalInventory{
					IsEnabled: Ptr(true),
					Schedule: &IncrementInventorySchedule{
						Frequency: Ptr(int64(600)),
					},
					OptionalFields: &IncrementalInventoryOptionalFields{
						Fields: []IncrementalInventoryOptionalFieldType{
							IncrementalInventoryOptionalFieldSequenceNumber,
							IncrementalInventoryOptionalFieldRecordType,
							IncrementalInventoryOptionalFieldRecordTimestamp,
							IncrementalInventoryOptionalFieldRequester,
							IncrementalInventoryOptionalFieldSourceIp,
							IncrementalInventoryOptionalFieldRequestId,
							IncrementalInventoryOptionalFieldSize,
							IncrementalInventoryOptionalFieldStorageClass,
							IncrementalInventoryOptionalFieldLastModifiedDate,
							IncrementalInventoryOptionalFieldETag,
							IncrementalInventoryOptionalFieldIsMultipartUploaded,
							IncrementalInventoryOptionalFieldObjectType,
							IncrementalInventoryOptionalFieldObjectAcl,
							IncrementalInventoryOptionalFieldCRC64,
							IncrementalInventoryOptionalFieldEncryptionStatus,
						},
					},
				},
			},
		},
		func(t *testing.T, o *PutBucketInventoryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
		},
	},
}

func TestMockPutBucketInventory_Success(t *testing.T) {
	for _, c := range testMockPutBucketInventorySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockPutBucketInventoryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *PutBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *PutBucketInventoryResult, err error)
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
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions></InventoryConfiguration>")
		},
		&PutBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
			InventoryConfiguration: &InventoryConfiguration{
				Id:        Ptr("report1"),
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
						AccountId: Ptr("1000000000000000"),
						RoleArn:   Ptr("acs:ram::1000000000000000:role/AliyunOSSRole"),
						Bucket:    Ptr("acs:oss:::destination-bucket"),
					},
				},
				Schedule: &InventorySchedule{
					Frequency: InventoryFrequencyDaily,
				},
				IncludedObjectVersions: Ptr("All"),
			},
		},
		func(t *testing.T, o *PutBucketInventoryResult, err error) {
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
			assert.Equal(t, "PUT", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			data, _ := io.ReadAll(r.Body)
			assert.Equal(t, string(data), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions></InventoryConfiguration>")
		},
		&PutBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
			InventoryConfiguration: &InventoryConfiguration{
				Id:        Ptr("report1"),
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
						AccountId: Ptr("1000000000000000"),
						RoleArn:   Ptr("acs:ram::1000000000000000:role/AliyunOSSRole"),
						Bucket:    Ptr("acs:oss:::destination-bucket"),
					},
				},
				Schedule: &InventorySchedule{
					Frequency: InventoryFrequencyDaily,
				},
				IncludedObjectVersions: Ptr("All"),
			},
		},
		func(t *testing.T, o *PutBucketInventoryResult, err error) {
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

func TestMockPutBucketInventory_Error(t *testing.T) {
	for _, c := range testMockPutBucketInventoryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)
		output, err := client.PutBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInventorySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInventoryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<InventoryConfiguration>
     <Id>report1</Id>
     <IsEnabled>true</IsEnabled>
     <Destination>
        <OSSBucketDestination>
           <Format>CSV</Format>
           <AccountId>1000000000000000</AccountId>
           <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
           <Bucket>acs:oss:::bucket_0001</Bucket>
           <Prefix>prefix1</Prefix>
           <Encryption>
              <SSE-KMS>
                 <KeyId>keyId</KeyId>
              </SSE-KMS>
           </Encryption>
        </OSSBucketDestination>
     </Destination>
     <Schedule>
        <Frequency>Daily</Frequency>
     </Schedule>
     <Filter>
        <LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp>
        <LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp>
        <LowerSizeBound>1024</LowerSizeBound>
        <UpperSizeBound>1048576</UpperSizeBound>
        <StorageClass>Standard,IA</StorageClass>
       	<Prefix>myprefix/</Prefix>
     </Filter>
     <IncludedObjectVersions>All</IncludedObjectVersions>
     <OptionalFields>
        <Field>Size</Field>
        <Field>LastModifiedDate</Field>
        <Field>ETag</Field>
        <Field>StorageClass</Field>
        <Field>IsMultipartUploaded</Field>
        <Field>EncryptionStatus</Field>
     </OptionalFields>
  </InventoryConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.InventoryConfiguration.Id, "report1")
			assert.True(t, *o.InventoryConfiguration.IsEnabled)
			assert.Equal(t, o.InventoryConfiguration.Destination.OSSBucketDestination.Format, InventoryFormatCSV)
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.AccountId, "1000000000000000")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.Bucket, "acs:oss:::bucket_0001")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.Prefix, "prefix1")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.Encryption.SseKms.KeyId, "keyId")
			assert.Equal(t, o.InventoryConfiguration.Schedule.Frequency, InventoryFrequencyDaily)
			assert.Equal(t, *o.InventoryConfiguration.IncludedObjectVersions, "All")
			assert.Equal(t, len(o.InventoryConfiguration.OptionalFields.Fields), 6)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)

			assert.Equal(t, *o.InventoryConfiguration.Filter.Prefix, "myprefix/")
			assert.Equal(t, *o.InventoryConfiguration.Filter.LastModifyBeginTimeStamp, int64(1637883649))
			assert.Equal(t, *o.InventoryConfiguration.Filter.LastModifyEndTimeStamp, int64(1638347592))
			assert.Equal(t, *o.InventoryConfiguration.Filter.LowerSizeBound, int64(1024))
			assert.Equal(t, *o.InventoryConfiguration.Filter.UpperSizeBound, int64(1048576))
			assert.Equal(t, *o.InventoryConfiguration.Filter.StorageClass, "Standard,IA")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<InventoryConfiguration>
    <Id>report1</Id>
    <IsEnabled>true</IsEnabled>
    <Destination>
        <OSSBucketDestination>
            <Format>CSV</Format>
            <AccountId>1000000000000000</AccountId>
            <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
            <Bucket>acs:oss:::destination-bucket</Bucket>
        </OSSBucketDestination>
    </Destination>
    <Schedule>
        <Frequency>Weekly</Frequency>
    </Schedule>
    <IncludedObjectVersions>Current</IncludedObjectVersions>
</InventoryConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.InventoryConfiguration.Id, "report1")
			assert.True(t, *o.InventoryConfiguration.IsEnabled)
			assert.Equal(t, o.InventoryConfiguration.Destination.OSSBucketDestination.Format, InventoryFormatCSV)
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.AccountId, "1000000000000000")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
			assert.Equal(t, o.InventoryConfiguration.Schedule.Frequency, InventoryFrequencyWeekly)
			assert.Equal(t, *o.InventoryConfiguration.IncludedObjectVersions, "Current")
		},
	},
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<InventoryConfiguration>
    <Id>report1</Id>
    <IsEnabled>true</IsEnabled>
    <Destination>
        <OSSBucketDestination>
            <Format>CSV</Format>
            <AccountId>1000000000000000</AccountId>
            <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
            <Bucket>acs:oss:::destination-bucket</Bucket>
        </OSSBucketDestination>
    </Destination>
    <Schedule>
        <Frequency>Weekly</Frequency>
    </Schedule>
    <IncludedObjectVersions>Current</IncludedObjectVersions>
      <OptionalFields>
        <Field>Size</Field>
        <Field>LastModifiedDate</Field>
        <Field>ETag</Field>
        <Field>StorageClass</Field>
        <Field>IsMultipartUploaded</Field>
        <Field>EncryptionStatus</Field>
		<Field>TransitionTime</Field>
		<Field>ObjectAcl</Field>
		<Field>TaggingCount</Field>
		<Field>ObjectType</Field>
		<Field>Crc64</Field>
     </OptionalFields>
    <IncrementalInventory>
        <IsEnabled>true</IsEnabled>
		<Schedule>
        	<Frequency>600</Frequency>
		</Schedule>
		<OptionalFields>
         <Field>SequenceNumber</Field>
         <Field>RecordType</Field>
         <Field>RecordTimestamp</Field>
         <Field>Requester</Field>
         <Field>SourceIp</Field>
         <Field>RequestId</Field>
         <Field>Size</Field>
         <Field>StorageClass</Field>
         <Field>LastModifiedDate</Field>
         <Field>ETag</Field>
         <Field>IsMultipartUploaded</Field>
         <Field>ObjectType</Field>
         <Field>ObjectAcl</Field>
         <Field>Crc64</Field>
         <Field>EncryptionStatus</Field>
      </OptionalFields>
     </IncrementalInventory>
</InventoryConfiguration>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, *o.InventoryConfiguration.Id, "report1")
			assert.True(t, *o.InventoryConfiguration.IsEnabled)
			assert.Equal(t, o.InventoryConfiguration.Destination.OSSBucketDestination.Format, InventoryFormatCSV)
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.AccountId, "1000000000000000")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
			assert.Equal(t, *o.InventoryConfiguration.Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
			assert.Equal(t, o.InventoryConfiguration.Schedule.Frequency, InventoryFrequencyWeekly)
			assert.Equal(t, *o.InventoryConfiguration.IncludedObjectVersions, "Current")

			assert.Equal(t, len(o.InventoryConfiguration.OptionalFields.Fields), 11)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[0], InventoryOptionalFieldSize)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[1], InventoryOptionalFieldLastModifiedDate)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[2], InventoryOptionalFieldETag)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[4], InventoryOptionalFieldIsMultipartUploaded)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[5], InventoryOptionalFieldEncryptionStatus)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[6], InventoryOptionalFieldTransitionTime)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[7], InventoryOptionalFieldObjectAcl)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[8], InventoryOptionalFieldTaggingCount)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[9], InventoryOptionalFieldObjectType)
			assert.Equal(t, o.InventoryConfiguration.OptionalFields.Fields[10], InventoryOptionalFieldCRC64)

			assert.Equal(t, *o.InventoryConfiguration.IncrementalInventory.Schedule.Frequency, int64(600))
			assert.Equal(t, *o.InventoryConfiguration.IncrementalInventory.IsEnabled, true)
			assert.Equal(t, len(o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields), 15)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[0], IncrementalInventoryOptionalFieldSequenceNumber)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[1], IncrementalInventoryOptionalFieldRecordType)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[2], IncrementalInventoryOptionalFieldRecordTimestamp)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[3], IncrementalInventoryOptionalFieldRequester)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[4], IncrementalInventoryOptionalFieldSourceIp)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[5], IncrementalInventoryOptionalFieldRequestId)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[6], IncrementalInventoryOptionalFieldSize)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[7], IncrementalInventoryOptionalFieldStorageClass)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[8], IncrementalInventoryOptionalFieldLastModifiedDate)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[9], IncrementalInventoryOptionalFieldETag)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[10], IncrementalInventoryOptionalFieldIsMultipartUploaded)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[11], IncrementalInventoryOptionalFieldObjectType)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[12], IncrementalInventoryOptionalFieldObjectAcl)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[13], IncrementalInventoryOptionalFieldCRC64)
			assert.Equal(t, o.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[14], IncrementalInventoryOptionalFieldEncryptionStatus)
		},
	},
}

func TestMockGetBucketInventory_Success(t *testing.T) {
	for _, c := range testMockGetBucketInventorySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockGetBucketInventoryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *GetBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *GetBucketInventoryResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", urlStr)
			assert.Equal(t, "GET", r.Method)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", strUrl)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", strUrl)
		},
		&GetBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *GetBucketInventoryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute GetBucketInventory fail")
		},
	},
}

func TestMockGetBucketInventory_Error(t *testing.T) {
	for _, c := range testMockGetBucketInventoryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.GetBucketInventory(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketInventorySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketInventoryResult, err error)
}{
	{
		200,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<ListInventoryConfigurationsResult>
     <InventoryConfiguration>
        <Id>report1</Id>
        <IsEnabled>true</IsEnabled>
        <Destination>
           <OSSBucketDestination>
              <Format>CSV</Format>
              <AccountId>1000000000000000</AccountId>
              <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
              <Bucket>acs:oss:::destination-bucket</Bucket>
              <Prefix>prefix1</Prefix>
           </OSSBucketDestination>
        </Destination>
        <Schedule>
           <Frequency>Daily</Frequency>
        </Schedule>
        <Filter>
           <Prefix>prefix/One</Prefix>
        </Filter>
        <IncludedObjectVersions>All</IncludedObjectVersions>
        <OptionalFields>
           <Field>Size</Field>
           <Field>LastModifiedDate</Field>
           <Field>ETag</Field>
           <Field>StorageClass</Field>
           <Field>IsMultipartUploaded</Field>
           <Field>EncryptionStatus</Field>
        </OptionalFields>
     </InventoryConfiguration>
     <InventoryConfiguration>
        <Id>report2</Id>
        <IsEnabled>true</IsEnabled>
        <Destination>
           <OSSBucketDestination>
              <Format>CSV</Format>
              <AccountId>1000000000000000</AccountId>
              <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
              <Bucket>acs:oss:::destination-bucket</Bucket>
              <Prefix>prefix2</Prefix>
           </OSSBucketDestination>
        </Destination>
        <Schedule>
           <Frequency>Daily</Frequency>
        </Schedule>
        <Filter>
           <Prefix>prefix/Two</Prefix>
        </Filter>
        <IncludedObjectVersions>All</IncludedObjectVersions>
        <OptionalFields>
           <Field>Size</Field>
           <Field>LastModifiedDate</Field>
           <Field>ETag</Field>
           <Field>StorageClass</Field>
           <Field>IsMultipartUploaded</Field>
           <Field>EncryptionStatus</Field>
        </OptionalFields>
     </InventoryConfiguration>
     <InventoryConfiguration>
        <Id>report3</Id>
        <IsEnabled>true</IsEnabled>
        <Destination>
           <OSSBucketDestination>
              <Format>CSV</Format>
              <AccountId>1000000000000000</AccountId>
              <RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn>
              <Bucket>acs:oss:::destination-bucket</Bucket>
              <Prefix>prefix3</Prefix>
           </OSSBucketDestination>
        </Destination>
        <Schedule>
           <Frequency>Daily</Frequency>
        </Schedule>
        <Filter>
           <Prefix>prefix/Three</Prefix>
        </Filter>
        <IncludedObjectVersions>All</IncludedObjectVersions>
        <OptionalFields>
           <Field>Size</Field>
           <Field>LastModifiedDate</Field>
           <Field>ETag</Field>
           <Field>StorageClass</Field>
           <Field>IsMultipartUploaded</Field>
           <Field>EncryptionStatus</Field>
        </OptionalFields>
     </InventoryConfiguration>
     <IsTruncated>false</IsTruncated>
  </ListInventoryConfigurationsResult>`),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory", urlStr)
		},
		&ListBucketInventoryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketInventoryResult, err error) {
			assert.Equal(t, 200, o.StatusCode)
			assert.Equal(t, "200 OK", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))
			assert.Equal(t, len(o.ListInventoryConfigurationsResult.InventoryConfigurations), 3)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Id, "report1")
			assert.True(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].IsEnabled)
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Format, InventoryFormatCSV)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.AccountId, "1000000000000000")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Prefix, "prefix1")
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Schedule.Frequency, InventoryFrequencyDaily)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].IncludedObjectVersions, "All")
			assert.Equal(t, len(o.ListInventoryConfigurationsResult.InventoryConfigurations[0].OptionalFields.Fields), 6)
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[0].OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)

			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[0].Filter.Prefix, "prefix/One")

			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Id, "report2")
			assert.True(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].IsEnabled)
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Format, InventoryFormatCSV)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.AccountId, "1000000000000000")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Prefix, "prefix2")
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Schedule.Frequency, InventoryFrequencyDaily)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].IncludedObjectVersions, "All")
			assert.Equal(t, len(o.ListInventoryConfigurationsResult.InventoryConfigurations[1].OptionalFields.Fields), 6)
			assert.Equal(t, o.ListInventoryConfigurationsResult.InventoryConfigurations[1].OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)
			assert.Equal(t, *o.ListInventoryConfigurationsResult.InventoryConfigurations[1].Filter.Prefix, "prefix/Two")
			assert.False(t, *o.ListInventoryConfigurationsResult.IsTruncated)
		},
	},
}

func TestMockListBucketInventory_Success(t *testing.T) {
	for _, c := range testMockListBucketInventorySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockListBucketInventoryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *ListBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *ListBucketInventoryResult, err error)
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
			urlStr := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory", urlStr)
			assert.Equal(t, "GET", r.Method)
		},
		&ListBucketInventoryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketInventoryResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory", strUrl)
		},
		&ListBucketInventoryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketInventoryResult, err error) {
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
			assert.Equal(t, "GET", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory", strUrl)
		},
		&ListBucketInventoryRequest{
			Bucket: Ptr("bucket"),
		},
		func(t *testing.T, o *ListBucketInventoryResult, err error) {
			assert.Nil(t, o)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), "execute ListBucketInventory fail")
		},
	},
}

func TestMockListBucketInventory_Error(t *testing.T) {
	for _, c := range testMockListBucketInventoryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.ListBucketInventory(context.TODO(), c.Request)

		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketInventorySuccessCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketInventoryResult, err error)
}{
	{
		204,
		map[string]string{
			"x-oss-request-id": "534B371674E88A4D8906****",
			"Date":             "Fri, 24 Feb 2017 03:15:40 GMT",
		},
		[]byte(``),
		func(t *testing.T, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", strUrl)
		},
		&DeleteBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *DeleteBucketInventoryResult, err error) {
			assert.Equal(t, 204, o.StatusCode)
			assert.Equal(t, "204 No Content", o.Status)
			assert.Equal(t, "534B371674E88A4D8906****", o.Headers.Get("x-oss-request-id"))
			assert.Equal(t, "Fri, 24 Feb 2017 03:15:40 GMT", o.Headers.Get("Date"))

		},
	},
}

func TestMockDeleteBucketInventory_Success(t *testing.T) {
	for _, c := range testMockDeleteBucketInventorySuccessCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}

var testMockDeleteBucketInventoryErrorCases = []struct {
	StatusCode     int
	Headers        map[string]string
	Body           []byte
	CheckRequestFn func(t *testing.T, r *http.Request)
	Request        *DeleteBucketInventoryRequest
	CheckOutputFn  func(t *testing.T, o *DeleteBucketInventoryResult, err error)
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
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", strUrl)
		},
		&DeleteBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *DeleteBucketInventoryResult, err error) {
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
			assert.Equal(t, "DELETE", r.Method)
			strUrl := sortQuery(r)
			assert.Equal(t, "/bucket/?inventory&inventoryId=report1", strUrl)
		},
		&DeleteBucketInventoryRequest{
			Bucket:      Ptr("bucket"),
			InventoryId: Ptr("report1"),
		},
		func(t *testing.T, o *DeleteBucketInventoryResult, err error) {
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

func TestMockDeleteBucketInventory_Error(t *testing.T) {
	for _, c := range testMockDeleteBucketInventoryErrorCases {
		server := testSetupMockServer(t, c.StatusCode, c.Headers, c.Body, c.CheckRequestFn)
		defer server.Close()
		assert.NotNil(t, server)

		cfg := LoadDefaultConfig().
			WithCredentialsProvider(credentials.NewAnonymousCredentialsProvider()).
			WithRegion("cn-hangzhou").
			WithEndpoint(server.URL)

		client := NewClient(cfg)
		assert.NotNil(t, c)

		output, err := client.DeleteBucketInventory(context.TODO(), c.Request)
		c.CheckOutputFn(t, output, err)
	}
}


