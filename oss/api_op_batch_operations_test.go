package oss

import (
	"io"
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
	"github.com/stretchr/testify/assert"
)

func TestMarshalInput_CreateJob(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var request *CreateJobRequest
	var input *OperationInput
	var err error

	request = &CreateJobRequest{}
	input = &OperationInput{
		OpName: "CreateJob",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field, CreateJob.")

	request = &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				DeleteObjectTagging: Ptr(""),
			},
			KeyPrefixManifestGenerator: &KeyPrefixManifestGenerator{
				SourceBucket: Ptr("bucket"),
				Prefix:       Ptr("batch-manifests/"),
			},
			Report: &Report{
				Bucket:      Ptr("bucket"),
				Enabled:     Ptr(true),
				Prefix:      Ptr("batch-reports/"),
				ReportScope: Ptr("AllTasks"),
			},
			Priority: Ptr(int32(10)),
			RoleArn:  Ptr("acs:ram::1234567890:role/oss-sdk-batch-test"),
		},
	}
	input = &OperationInput{
		OpName: "CreateJob",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ := io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<CreateJobRequest><Operation><DeleteObjectTagging></DeleteObjectTagging></Operation><Report><Bucket>bucket</Bucket><Enabled>true</Enabled><Prefix>batch-reports/</Prefix><ReportScope>AllTasks</ReportScope></Report><KeyPrefixManifestGenerator><SourceBucket>bucket</SourceBucket><Prefix>batch-manifests/</Prefix></KeyPrefixManifestGenerator><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")

	request = &CreateJobRequest{
		CreateJob: &CreateJobConfig{
			Operation: &Operation{
				PutObjectTagging: &Tagging{
					TagSet: &TagSet{
						Tags: []Tag{
							{
								Key:   Ptr("key1"),
								Value: Ptr("value1"),
							},
						},
					},
				},
			},
			ClientRequestToken: Ptr("unique-token-123"),
			Manifest: &Manifest{
				Location: &Location{
					ETag:   Ptr("d41d8cd98f00b204e9800998ecf8427e"),
					Bucket: Ptr("manifest-bucket"),
					Object: Ptr("manifest.csv"),
				},
				Spec: &Spec{
					Fields: Ptr("Bucket,Key"),
					Format: Ptr("OSS_BatchOperations_CSV_20250611"),
				},
			},
			Description: Ptr("批量设置对象标签任务"),
			Priority:    Ptr(int32(10)),
			RoleArn:     Ptr("acs:ram::1234567890:role/oss-sdk-batch-test"),
		},
	}
	input = &OperationInput{
		OpName: "CreateJob",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	err = c.marshalInput(request, input, updateContentMd5)
	assert.Nil(t, err)
	body, _ = io.ReadAll(input.Body)
	assert.Equal(t, string(body), "<CreateJobRequest><Operation><PutObjectTagging><TagSet><Tag><Key>key1</Key><Value>value1</Value></Tag></TagSet></PutObjectTagging></Operation><ClientRequestToken>unique-token-123</ClientRequestToken><Manifest><Location><ETag>d41d8cd98f00b204e9800998ecf8427e</ETag><Bucket>manifest-bucket</Bucket><Object>manifest.csv</Object></Location><Spec><Fields>Bucket,Key</Fields><Format>OSS_BatchOperations_CSV_20250611</Format></Spec></Manifest><Description>批量设置对象标签任务</Description><Priority>10</Priority><RoleArn>acs:ram::1234567890:role/oss-sdk-batch-test</RoleArn></CreateJobRequest>")
}