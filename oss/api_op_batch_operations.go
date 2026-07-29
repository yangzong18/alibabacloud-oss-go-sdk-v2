package oss

import (
	"context"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/signer"
)

type CreateJobRequest struct {
	// The name of the bucket to create.
	Bucket *string `input:"host,bucket,required"`

	CreateJob *CreateJobConfig `input:"body,CreateJobRequest,xml"`

	RequestCommon
}

type CreateJobConfig struct {
	ConfirmationRequired *bool `xml:"ConfirmationRequired"`

	Operation *Operation `xml:"Operation"`

	Report *Report `xml:"Report"`

	ClientRequestToken *string `xml:"ClientRequestToken"`

	Manifest *Manifest `xml:"Manifest"`

	Description *string `xml:"Description"`

	Priority *int32 `xml:"Priority"`

	RoleArn *string `xml:"RoleArn"`
}

type Operation struct {
	PutObjectTagging *Tagging `xml:"PutObjectTagging"`

	DeleteObjectTagging *string `xml:"DeleteObjectTagging"`

	AddObjectTagging *Tagging `xml:"AddObjectTagging"`

	PutObjectAcl *PutObjectAcl `xml:"PutObjectAcl"`

	RestoreObject *RestoreObject `xml:"RestoreObject"`
}

type PutObjectAcl struct {
	ObjectAcl ObjectACLType `xml:"ObjectAcl"`
}

type RestoreObject struct {
	Days int32 `xml:"Days"`

	Tier *string `xml:"Tier"`
}

type Report struct {
	Bucket *string `xml:"Bucket"`

	Enabled *bool `xml:"Enabled"`

	Prefix *string `xml:"Prefix"`

	ReportScope *string `xml:"ReportScope"`
}

type Manifest struct {
	Location *Location `xml:"Location"`

	Spec *Spec `xml:"Spec"`
}

type Location struct {
	ETag *string `xml:"ETag"`

	Bucket *string `xml:"Bucket"`

	Object *string `xml:"Object"`

	VersionId *string `xml:"VersionId"`
}

type Spec struct {
	Fields *string `xml:"Fields"`
	Format *string `xml:"Format"`
}

type CreateJobResult struct {
	JobId *string `xml:"JobId"`

	ResultCommon
}

// CreateJob Creates a batch job that performs a specified operation on multiple objects.
func (c *Client) CreateJob(ctx context.Context, request *CreateJobRequest, optFns ...func(*Options)) (*CreateJobResult, error) {
	var err error
	if request == nil {
		request = &CreateJobRequest{}
	}
	input := &OperationInput{
		OpName: "CreateJob",
		Method: "POST",
		Headers: map[string]string{
			HTTPHeaderContentType: contentTypeXML,
		},
		Parameters: map[string]string{
			"batchJob": "",
		},
		Bucket: request.Bucket,
	}
	input.OpMetadata.Set(signer.SubResource, []string{"batchJob"})
	if err = c.marshalInput(request, input, updateContentMd5); err != nil {
		return nil, err
	}

	output, err := c.invokeOperation(ctx, input, optFns)
	if err != nil {
		return nil, err
	}

	result := &CreateJobResult{}

	if err = c.unmarshalOutput(result, output, discardBody); err != nil {
		return nil, c.toClientError(err, "UnmarshalOutputFail", output)
	}

	return result, err
}
