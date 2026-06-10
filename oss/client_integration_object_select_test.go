//go:build integration

package oss

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSelectObjectMeta(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := "name,school,company,age\r\nLora Francis,School A,Staples Inc,27\r\n" + "Eleanor Little,School B,\"Conectiv, Inc\",43\r\n" + "Rosie Hughes,School C,Western Gas Resources Inc,44\r\n" + "Lawrence Ross,School D,MetLife Inc.,24"
	objectNameCsv := objectNamePrefix + randLowStr(6) + ".csv"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	csvMeta := &CreateSelectObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		MetaRequest: &CsvMetaRequest{
			OverwriteIfExists: Ptr(true),
		},
	}
	result, err := client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.RowsCount, int64(5))

	body = "{\n" +
		"\t\"name\": \"Lora Francis\",\n" +
		"\t\"age\": 27,\n" +
		"\t\"company\": \"Staples Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Eleanor Little\",\n" +
		"\t\"age\": 43,\n" +
		"\t\"company\": \"Conectiv, Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Rosie Hughes\",\n" +
		"\t\"age\": 44,\n" +
		"\t\"company\": \"Western Gas Resources Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Lawrence Ross\",\n" +
		"\t\"age\": 24,\n" +
		"\t\"company\": \"MetLife Inc.\"\n" +
		"}"
	objectNameJson := objectNamePrefix + randLowStr(6) + ".json"
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		Body:   strings.NewReader(string(body)),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	csvMeta = &CreateSelectObjectMetaRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		MetaRequest: &JsonMetaRequest{
			InputSerialization: &InputSerialization{
				JSON: &InputSerializationJSON{
					JSONType: Ptr("LINES"),
				},
			},
		},
	}
	result, err = client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.Nil(t, err)
	assert.Equal(t, result.RowsCount, int64(4))

	_, err = client.CreateSelectObjectMeta(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	csvMeta = &CreateSelectObjectMetaRequest{
		Bucket:      Ptr(bucketNameNotExist),
		Key:         Ptr(objectNameCsv),
		MetaRequest: &CsvMetaRequest{},
	}
	result, err = client.CreateSelectObjectMeta(context.TODO(), csvMeta)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}

func TestSelectObject(t *testing.T) {
	after := before(t)
	defer after(t)

	bucketName := bucketNamePrefix + randLowStr(6)
	//TODO
	putRequest := &PutBucketRequest{
		Bucket: Ptr(bucketName),
	}

	client := getDefaultClient()
	_, err := client.PutBucket(context.TODO(), putRequest)
	assert.Nil(t, err)

	body := "name,school,company,age\r\nLora Francis,School A,Staples Inc,27\r\n" + "Eleanor Little,School B,\"Conectiv, Inc\",43\r\n" + "Rosie Hughes,School C,Western Gas Resources Inc,44\r\n" + "Lawrence Ross,School D,MetLife Inc.,24"
	objectNameCsv := objectNamePrefix + randLowStr(6) + ".csv"
	putObjRequest := &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		Body:   strings.NewReader(body),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)

	request := &SelectObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err := client.SelectObject(context.TODO(), request)
	assert.Nil(t, err)
	dataByte, err := io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte), "name\nLora Francis\nEleanor Little\nRosie Hughes\nLawrence Ross\n")

	body = "{\n" +
		"\t\"name\": \"Lora Francis\",\n" +
		"\t\"age\": 27,\n" +
		"\t\"company\": \"Staples Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Eleanor Little\",\n" +
		"\t\"age\": 43,\n" +
		"\t\"company\": \"Conectiv, Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Rosie Hughes\",\n" +
		"\t\"age\": 44,\n" +
		"\t\"company\": \"Western Gas Resources Inc\"\n" +
		"}\n" +
		"{\n" +
		"\t\"name\": \"Lawrence Ross\",\n" +
		"\t\"age\": 24,\n" +
		"\t\"company\": \"MetLife Inc.\"\n" +
		"}"
	objectNameJson := objectNamePrefix + randLowStr(6) + ".json"
	putObjRequest = &PutObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameJson),
		Body:   strings.NewReader(string(body)),
	}
	_, err = client.PutObject(context.TODO(), putObjRequest)
	assert.Nil(t, err)
	request = &SelectObjectRequest{
		Bucket: Ptr(bucketName),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err = client.SelectObject(context.TODO(), request)
	assert.Nil(t, err)
	dataByte, err = io.ReadAll(result.Body)
	assert.Equal(t, string(dataByte), "name\nLora Francis\nEleanor Little\nRosie Hughes\nLawrence Ross\n")

	_, err = client.SelectObject(context.TODO(), nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing required field")

	var serr *ServiceError
	bucketNameNotExist := bucketName + "-not-exist"
	request = &SelectObjectRequest{
		Bucket: Ptr(bucketNameNotExist),
		Key:    Ptr(objectNameCsv),
		SelectRequest: &SelectRequest{
			Expression: Ptr("select name from ossobject"),
			InputSerializationSelect: InputSerializationSelect{
				CsvBodyInput: &CSVSelectInput{
					FileHeaderInfo: Ptr("Use"),
				},
			},
			OutputSerializationSelect: OutputSerializationSelect{
				OutputHeader: Ptr(true),
			},
		},
	}
	result, err = client.SelectObject(context.TODO(), request)
	assert.NotNil(t, err)
	errors.As(err, &serr)
	assert.Equal(t, int(404), serr.StatusCode)
	assert.Equal(t, "NoSuchBucket", serr.Code)
	assert.Equal(t, "The specified bucket does not exist.", serr.Message)
	assert.Equal(t, "0015-00000101", serr.EC)
	assert.NotEmpty(t, serr.RequestID)
}
