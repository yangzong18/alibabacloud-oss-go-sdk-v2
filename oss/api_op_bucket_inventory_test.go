package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_PutBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *PutBucketInventoryRequest
	var input *OperationInput
	var err error

	request = &PutBucketInventoryRequest{}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &PutBucketInventoryRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, InventoryId.")

	request = &PutBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
		InventoryId: Ptr("report1"),
	}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, InventoryConfiguration.")

	request = &PutBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
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
					InventoryOptionalFieldTransitionTime,
					InventoryOptionalFieldObjectAcl,
					InventoryOptionalFieldTaggingCount,
					InventoryOptionalFieldObjectType,
					InventoryOptionalFieldCRC64,
				},
			},
		},
	}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket><Prefix>prefix1</Prefix><Encryption><SSE-KMS><KeyId>keyId</KeyId></SSE-KMS></Encryption></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions><OptionalFields><Field>Size</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>StorageClass</Field><Field>IsMultipartUploaded</Field><Field>EncryptionStatus</Field><Field>TransitionTime</Field><Field>ObjectAcl</Field><Field>TaggingCount</Field><Field>ObjectType</Field><Field>Crc64</Field></OptionalFields></InventoryConfiguration>")

	request = &PutBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
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
	}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket></OSSBucketDestination></Destination><Schedule><Frequency>Daily</Frequency></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions></InventoryConfiguration>")

	request = &PutBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
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
				Frequency:  InventoryFrequencyMonthly,
				DayOfMonth: Ptr(15),
			},
			IncludedObjectVersions: Ptr("All"),
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
	}
	input = &OperationInput{
		OpName: "PutBucketInventory",
		Method: "PUT",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<InventoryConfiguration><Id>report1</Id><IsEnabled>true</IsEnabled><Destination><OSSBucketDestination><Format>CSV</Format><AccountId>1000000000000000</AccountId><RoleArn>acs:ram::1000000000000000:role/AliyunOSSRole</RoleArn><Bucket>acs:oss:::destination-bucket</Bucket></OSSBucketDestination></Destination><Schedule><Frequency>Monthly</Frequency><DayOfMonth>15</DayOfMonth></Schedule><Filter><LastModifyBeginTimeStamp>1637883649</LastModifyBeginTimeStamp><LastModifyEndTimeStamp>1638347592</LastModifyEndTimeStamp><LowerSizeBound>1024</LowerSizeBound><UpperSizeBound>1048576</UpperSizeBound><StorageClass>Standard,IA</StorageClass><Prefix>filterPrefix</Prefix></Filter><IncludedObjectVersions>All</IncludedObjectVersions><IncrementalInventory><IsEnabled>true</IsEnabled><Schedule><Frequency>600</Frequency></Schedule><OptionalFields><Field>SequenceNumber</Field><Field>RecordType</Field><Field>RecordTimestamp</Field><Field>Requester</Field><Field>SourceIp</Field><Field>RequestId</Field><Field>Size</Field><Field>StorageClass</Field><Field>LastModifiedDate</Field><Field>ETag</Field><Field>IsMultipartUploaded</Field><Field>ObjectType</Field><Field>ObjectAcl</Field><Field>Crc64</Field><Field>EncryptionStatus</Field></OptionalFields></IncrementalInventory></InventoryConfiguration>")
}

func TestUnmarshalOutput_PutBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &PutBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
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
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &PutBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &PutBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_GetBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *GetBucketInventoryRequest
	var input *OperationInput
	var err error

	request = &GetBucketInventoryRequest{}
	input = &OperationInput{
		OpName: "GetBucketInventory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &GetBucketInventoryRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketInventory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, InventoryId.")

	request = &GetBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
		InventoryId: Ptr("report1"),
	}
	input = &OperationInput{
		OpName: "GetBucketInventory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_GetBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
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
  </InventoryConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &GetBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.InventoryConfiguration.Id, "report1")
	assert.True(t, *result.InventoryConfiguration.IsEnabled)
	assert.Equal(t, result.InventoryConfiguration.Destination.OSSBucketDestination.Format, InventoryFormatCSV)
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.AccountId, "1000000000000000")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.Bucket, "acs:oss:::bucket_0001")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.Prefix, "prefix1")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.Encryption.SseKms.KeyId, "keyId")
	assert.Equal(t, result.InventoryConfiguration.Schedule.Frequency, InventoryFrequencyDaily)
	assert.Equal(t, *result.InventoryConfiguration.IncludedObjectVersions, "All")
	assert.Equal(t, len(result.InventoryConfiguration.OptionalFields.Fields), 11)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[0], InventoryOptionalFieldSize)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[1], InventoryOptionalFieldLastModifiedDate)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[2], InventoryOptionalFieldETag)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[4], InventoryOptionalFieldIsMultipartUploaded)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[5], InventoryOptionalFieldEncryptionStatus)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[6], InventoryOptionalFieldTransitionTime)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[7], InventoryOptionalFieldObjectAcl)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[8], InventoryOptionalFieldTaggingCount)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[9], InventoryOptionalFieldObjectType)
	assert.Equal(t, result.InventoryConfiguration.OptionalFields.Fields[10], InventoryOptionalFieldCRC64)

	assert.Equal(t, *result.InventoryConfiguration.IncrementalInventory.Schedule.Frequency, int64(600))
	assert.Equal(t, *result.InventoryConfiguration.IncrementalInventory.IsEnabled, true)
	assert.Equal(t, len(result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields), 15)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[0], IncrementalInventoryOptionalFieldSequenceNumber)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[1], IncrementalInventoryOptionalFieldRecordType)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[2], IncrementalInventoryOptionalFieldRecordTimestamp)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[3], IncrementalInventoryOptionalFieldRequester)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[4], IncrementalInventoryOptionalFieldSourceIp)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[5], IncrementalInventoryOptionalFieldRequestId)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[6], IncrementalInventoryOptionalFieldSize)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[7], IncrementalInventoryOptionalFieldStorageClass)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[8], IncrementalInventoryOptionalFieldLastModifiedDate)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[9], IncrementalInventoryOptionalFieldETag)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[10], IncrementalInventoryOptionalFieldIsMultipartUploaded)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[11], IncrementalInventoryOptionalFieldObjectType)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[12], IncrementalInventoryOptionalFieldObjectAcl)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[13], IncrementalInventoryOptionalFieldCRC64)
	assert.Equal(t, result.InventoryConfiguration.IncrementalInventory.OptionalFields.Fields[14], IncrementalInventoryOptionalFieldEncryptionStatus)

	assert.Equal(t, *result.InventoryConfiguration.Filter.Prefix, "myprefix/")
	assert.Equal(t, *result.InventoryConfiguration.Filter.LastModifyBeginTimeStamp, int64(1637883649))
	assert.Equal(t, *result.InventoryConfiguration.Filter.LastModifyEndTimeStamp, int64(1638347592))
	assert.Equal(t, *result.InventoryConfiguration.Filter.LowerSizeBound, int64(1024))
	assert.Equal(t, *result.InventoryConfiguration.Filter.UpperSizeBound, int64(1048576))
	assert.Equal(t, *result.InventoryConfiguration.Filter.StorageClass, "Standard,IA")
	body = `<InventoryConfiguration>
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
</InventoryConfiguration>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result = &GetBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, *result.InventoryConfiguration.Id, "report1")
	assert.True(t, *result.InventoryConfiguration.IsEnabled)
	assert.Equal(t, result.InventoryConfiguration.Destination.OSSBucketDestination.Format, InventoryFormatCSV)
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.AccountId, "1000000000000000")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
	assert.Equal(t, *result.InventoryConfiguration.Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
	assert.Equal(t, result.InventoryConfiguration.Schedule.Frequency, InventoryFrequencyWeekly)
	assert.Equal(t, *result.InventoryConfiguration.IncludedObjectVersions, "Current")

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &GetBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &GetBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_ListBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *ListBucketInventoryRequest
	var input *OperationInput
	var err error

	request = &ListBucketInventoryRequest{}
	input = &OperationInput{
		OpName: "ListBucketInventory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &ListBucketInventoryRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "GetBucketInventory",
		Method: "GET",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_ListBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	body := `<?xml version="1.0" encoding="UTF-8"?>
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
           <Frequency>Weekly</Frequency>
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
           <Frequency>Monthly</Frequency>
           <DayOfMonth>1</DayOfMonth>
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
  </ListInventoryConfigurationsResult>`
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &ListBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 200)
	assert.Equal(t, result.Status, "OK")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, len(result.ListInventoryConfigurationsResult.InventoryConfigurations), 3)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Id, "report1")
	assert.True(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].IsEnabled)
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Format, InventoryFormatCSV)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.AccountId, "1000000000000000")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Destination.OSSBucketDestination.Prefix, "prefix1")
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Schedule.Frequency, InventoryFrequencyDaily)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].IncludedObjectVersions, "All")
	assert.Equal(t, len(result.ListInventoryConfigurationsResult.InventoryConfigurations[0].OptionalFields.Fields), 6)
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[0].OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)

	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[0].Filter.Prefix, "prefix/One")

	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Id, "report2")
	assert.True(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].IsEnabled)
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Format, InventoryFormatCSV)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.AccountId, "1000000000000000")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Destination.OSSBucketDestination.Prefix, "prefix2")
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Schedule.Frequency, InventoryFrequencyWeekly)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].IncludedObjectVersions, "All")
	assert.Equal(t, len(result.ListInventoryConfigurationsResult.InventoryConfigurations[1].OptionalFields.Fields), 6)
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[1].OptionalFields.Fields[3], InventoryOptionalFieldStorageClass)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[1].Filter.Prefix, "prefix/Two")
	assert.False(t, *result.ListInventoryConfigurationsResult.IsTruncated)

	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Id, "report3")
	assert.True(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].IsEnabled)
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Destination.OSSBucketDestination.Format, InventoryFormatCSV)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Destination.OSSBucketDestination.AccountId, "1000000000000000")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Destination.OSSBucketDestination.RoleArn, "acs:ram::1000000000000000:role/AliyunOSSRole")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Destination.OSSBucketDestination.Bucket, "acs:oss:::destination-bucket")
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Destination.OSSBucketDestination.Prefix, "prefix3")
	assert.Equal(t, result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Schedule.Frequency, InventoryFrequencyMonthly)
	assert.Equal(t, *result.ListInventoryConfigurationsResult.InventoryConfigurations[2].Schedule.DayOfMonth, 1)

	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &ListBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &ListBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}

func TestMarshalInput_DeleteBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *DeleteBucketInventoryRequest
	var input *OperationInput
	var err error

	request = &DeleteBucketInventoryRequest{}
	input = &OperationInput{
		OpName: "DeleteBucketInventory",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, Bucket.")

	request = &DeleteBucketInventoryRequest{
		Bucket: Ptr("oss-demo"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketInventory",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, InventoryId.")

	request = &DeleteBucketInventoryRequest{
		Bucket:      Ptr("oss-demo"),
		InventoryId: Ptr("report1"),
	}
	input = &OperationInput{
		OpName: "DeleteBucketInventory",
		Method: "DELETE",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"inventory": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"inventory", "inventoryId"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
}

func TestUnmarshalOutput_DeleteBucketInventory(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error
	output = &OperationOutput{
		StatusCode: 204,
		Status:     "No Content",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
		},
	}
	result := &DeleteBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 204)
	assert.Equal(t, result.Status, "No Content")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")

	body := `<?xml version="1.0" encoding="UTF-8"?>
		<Error>
		<Code>NoSuchBucket</Code>
		<Message>The specified bucket does not exist.</Message>
		<RequestId>66C2FF09FDF07830343C72EC</RequestId>
		<HostId>bucket.oss-cn-hangzhou.aliyuncs.com</HostId>
		<BucketName>bucket</BucketName>
		<EC>0015-00000101</EC>
	</Error>`
	output = &OperationOutput{
		StatusCode: 404,
		Status:     "NoSuchBucket",
		Body:       io.NopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 404)
	assert.Equal(t, result.Status, "NoSuchBucket")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")

	output = &OperationOutput{
		StatusCode: 400,
		Status:     "InvalidArgument",
		Headers: http.Header{
			"X-Oss-Request-Id": {"534B371674E88A4D8906****"},
			"Content-Type":     {"application/xml"},
		},
	}
	result = &DeleteBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 400)
	assert.Equal(t, result.Status, "InvalidArgument")
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
	result = &DeleteBucketInventoryResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXmlMix)
	assert.Nil(t, err)
	assert.Equal(t, result.StatusCode, 403)
	assert.Equal(t, result.Status, "AccessDenied")
	assert.Equal(t, result.Headers.Get("X-Oss-Request-Id"), "534B371674E88A4D8906****")
	assert.Equal(t, result.Headers.Get("Content-Type"), "application/xml")
}
