package oss

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/v3/oss/readers"
	"github.com/stretchr/testify/assert"
)

type stubRequest struct {
	StrPrtField   *string `input:"query,str-field"`
	StrField      string  `input:"query,str-field"`
	IntPtrFiled   *int32  `input:"query,int32-field"`
	IntFiled      int32   `input:"query,int32-field"`
	BoolPtrFiled  *bool   `input:"query,bool-field"`
	HStrPrtField  *string `input:"header,x-oss-str-field"`
	HStrField     string  `input:"header,x-oss-str-field"`
	HIntPtrFiled  *int32  `input:"header,x-oss-int32-field"`
	HIntFiled     int32   `input:"header,x-oss-int32-field"`
	HBoolPtrFiled *bool   `input:"header,x-oss-bool-field"`
}

func TestMarshalInput(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *stubRequest
	var err error

	// nil request
	input = &OperationInput{}
	request = nil

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 0, len(input.Parameters))

	// emtpy request
	input = &OperationInput{}
	request = &stubRequest{}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 0, len(input.Parameters))

	// query ptr
	input = &OperationInput{}

	request = &stubRequest{
		StrPrtField:  Ptr("str1"),
		IntPtrFiled:  Ptr(int32(123)),
		BoolPtrFiled: Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str1", input.Parameters["str-field"])
	assert.Equal(t, "123", input.Parameters["int32-field"])
	assert.Equal(t, "true", input.Parameters["bool-field"])

	// query value
	input = &OperationInput{}

	request = &stubRequest{
		StrField:     "str2",
		IntFiled:     int32(223),
		BoolPtrFiled: Ptr(false),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Headers))
	assert.Equal(t, 3, len(input.Parameters))
	assert.Equal(t, "str2", input.Parameters["str-field"])
	assert.Equal(t, "223", input.Parameters["int32-field"])
	assert.Equal(t, "false", input.Parameters["bool-field"])

	// header ptr
	input = &OperationInput{}

	request = &stubRequest{
		HStrPrtField:  Ptr("str1"),
		HIntPtrFiled:  Ptr(int32(123)),
		HBoolPtrFiled: Ptr(true),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str1", input.Headers["x-oss-str-field"])
	assert.Equal(t, "123", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "true", input.Headers["x-oss-bool-field"])

	// header value
	input = &OperationInput{}

	request = &stubRequest{
		HStrField:     "str2",
		HIntFiled:     int32(223),
		HBoolPtrFiled: Ptr(false),
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(input.Parameters))
	assert.Equal(t, 3, len(input.Headers))
	assert.Equal(t, "str2", input.Headers["x-oss-str-field"])
	assert.Equal(t, "223", input.Headers["x-oss-int32-field"])
	assert.Equal(t, "false", input.Headers["x-oss-bool-field"])
}

type xmlbodyRequest struct {
	StrHostPrtField    *string        `input:"host,bucket,required"`
	StrQueryPrtField   *string        `input:"query,str-field"`
	StrHeaderPrtField  *string        `input:"header,x-oss-str-field"`
	StructBodyPrtField *xmlBodyConfig `input:"xmlbody,BodyConfiguration"`
}

type xmlBodyConfig struct {
	StrField1 *string `xml:"StrField1"`
	StrField2 string  `xml:"StrField2"`
}

func TestMarshalInput_xmlbody(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *xmlbodyRequest
	var err error

	input = &OperationInput{}
	request = &xmlbodyRequest{
		StrHostPrtField:   Ptr("bucket"),
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &xmlBodyConfig{
			StrField1: Ptr("StrField1"),
			StrField2: "StrField2",
		},
	}

	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(input.Parameters))
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, 1, len(input.Headers))
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Body)

	body, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>", string(body))
}

type commonStubRequest struct {
	StrHostPrtField    *string        `input:"host,bucket,required"`
	StrQueryPrtField   *string        `input:"query,str-field"`
	StrHeaderPrtField  *string        `input:"header,x-oss-str-field"`
	StructBodyPrtField *xmlBodyConfig `input:"xmlbody,BodyConfiguration"`
	RequestCommon
}

func TestMarshalInput_CommonFields(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *commonStubRequest
	var err error

	//default
	request = &commonStubRequest{}
	input = &OperationInput{}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Body)
	assert.Nil(t, input.Headers)
	assert.Nil(t, input.Parameters)

	//set by request
	request = &commonStubRequest{}
	request.Headers = map[string]string{
		"key": "value",
	}
	request.Parameters = map[string]string{
		"p": "",
	}
	request.Body = bytes.NewReader([]byte("hello"))
	input = &OperationInput{}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "value", input.Headers["key"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "", input.Parameters["p"])
	assert.NotNil(t, input.Body)
	data, err := io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(data))

	// priority
	// request commmn > input default
	input = &OperationInput{
		Headers: map[string]string{
			"key": "value1",
		},
		Parameters: map[string]string{
			"p": "value1",
		},
	}
	request = &commonStubRequest{}
	request.Headers = map[string]string{
		"key": "value2",
	}
	request.Parameters = map[string]string{
		"p": "value3",
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "value2", input.Headers["key"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "value3", input.Parameters["p"])
	assert.Nil(t, input.Body)

	// reuqest filed parametr > request commmn
	input = &OperationInput{}
	request = &commonStubRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		StructBodyPrtField: &xmlBodyConfig{
			StrField1: Ptr("StrField1"),
			StrField2: "StrField2",
		},
	}
	request.Headers = map[string]string{
		"x-oss-str-field": "value2",
	}
	request.Parameters = map[string]string{
		"str-field": "value3",
	}
	request.Body = bytes.NewReader([]byte("hello"))
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 1)
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 1)
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.NotNil(t, input.Body)
	data, err = io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>", string(data))

	// merge, replace
	//reuqest filed parametr > request commmn > input
	input = &OperationInput{
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &commonStubRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
	}
	request.Headers = map[string]string{
		"x-oss-str-field":  "value2",
		"x-oss-str-field1": "value2-1",
	}
	request.Parameters = map[string]string{
		"str-field1": "value3",
	}
	request.Body = bytes.NewReader([]byte("hello"))
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 3)
	assert.Equal(t, "value1", input.Headers["input-key"])
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.Equal(t, "value2-1", input.Headers["x-oss-str-field1"])
	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 4)
	assert.Equal(t, "value2", input.Parameters["input-param"])
	assert.Equal(t, "value2-1", input.Parameters["input-param1"])
	assert.Equal(t, "query", input.Parameters["str-field"])
	assert.Equal(t, "value3", input.Parameters["str-field1"])
	assert.NotNil(t, input.Body)
	data, err = io.ReadAll(input.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(data))
}

type usermetaRequest struct {
	StrQueryPrtField  *string           `input:"query,str-field"`
	StrHeaderPrtField *string           `input:"header,x-oss-str-field"`
	UserMetaField1    map[string]string `input:"usermeta,x-oss-meta-"`
	UserMetaField2    map[string]string `input:"usermeta,x-oss-meta1-"`
}

func TestMarshalInput_usermeta(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var input *OperationInput
	var request *usermetaRequest
	var err error

	input = &OperationInput{}
	request = &usermetaRequest{}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.Nil(t, input.Headers)

	input = &OperationInput{
		Headers: map[string]string{
			"input-key": "value1",
		},
		Parameters: map[string]string{
			"input-param":  "value2",
			"input-param1": "value2-1",
		}}
	request = &usermetaRequest{
		StrQueryPrtField:  Ptr("query"),
		StrHeaderPrtField: Ptr("header"),
		UserMetaField1: map[string]string{
			"user1": "value1",
			"user2": "value2",
		},
		UserMetaField2: map[string]string{
			"user3": "value3",
			"user4": "value4",
		},
	}
	err = c.marshalInput(request, input)
	assert.Nil(t, err)
	assert.NotNil(t, input.Headers)
	assert.Len(t, input.Headers, 6)
	assert.Equal(t, "value1", input.Headers["input-key"])
	assert.Equal(t, "header", input.Headers["x-oss-str-field"])
	assert.Equal(t, "value1", input.Headers["x-oss-meta-user1"])
	assert.Equal(t, "value2", input.Headers["x-oss-meta-user2"])
	assert.Equal(t, "value3", input.Headers["x-oss-meta1-user3"])
	assert.Equal(t, "value4", input.Headers["x-oss-meta1-user4"])

	assert.NotNil(t, input.Parameters)
	assert.Len(t, input.Parameters, 3)
	assert.Equal(t, "value2", input.Parameters["input-param"])
	assert.Equal(t, "value2-1", input.Parameters["input-param1"])
	assert.Equal(t, "query", input.Parameters["str-field"])
}

type stubResult struct {
	ResultCommon
}

type xmlBodyResult struct {
	ResultCommon
	StrField1 *string `xml:"StrField1"`
	StrField2 *string `xml:"StrField2"`
}

func TestUnmarshalOutput(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var result *stubResult
	var err error

	//empty
	output = &OperationOutput{}
	assert.Nil(t, output.Input)
	assert.Nil(t, output.Body)
	assert.Nil(t, output.Headers)
	assert.Empty(t, output.Status)
	assert.Empty(t, output.StatusCode)

	// with default values
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
	}

	result = &stubResult{}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Nil(t, result.Headers)

	// has header
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Headers: http.Header{
			"Expires":          {"-1"},
			"Content-Length":   {"0"},
			"Content-Encoding": {"gzip"},
		},
	}

	result = &stubResult{}
	err = c.unmarshalOutput(result, output)
	assert.Nil(t, err)
	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "OK", result.Status)
	assert.Equal(t, "-1", result.Headers.Get("Expires"))
	assert.Equal(t, "0", result.Headers.Get("Content-Length"))
	assert.Equal(t, "gzip", result.Headers.Get("Content-Encoding"))

	// extract body
	body := "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       readers.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	xmlresult := &xmlBodyResult{}
	err = c.unmarshalOutput(xmlresult, output, unmarshalBodyXml)
	assert.Nil(t, err)
	assert.Equal(t, 200, xmlresult.StatusCode)
	assert.Equal(t, "OK", xmlresult.Status)
	assert.Equal(t, "StrField1", *xmlresult.StrField1)
	assert.Equal(t, "StrField2", *xmlresult.StrField2)
}

func TestUnmarshalOutput_error(t *testing.T) {
	c := Client{}
	assert.NotNil(t, c)
	var output *OperationOutput
	var err error

	// unsupport content-type
	body := "<BodyConfiguration><StrField1>StrField1</StrField1><StrField2>StrField2</StrField2></BodyConfiguration>"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       readers.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/text"},
		},
	}
	xmlresult := &xmlBodyResult{}
	err = c.unmarshalOutput(xmlresult, output, unmarshalBodyDefault)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupport contentType:application/text")

	// xml decode fail
	body = "StrField1>StrField1</StrField1><StrField2>StrField2<"
	output = &OperationOutput{
		StatusCode: 200,
		Status:     "OK",
		Body:       readers.ReadSeekNopCloser(bytes.NewReader([]byte(body))),
		Headers: http.Header{
			"Content-Type": {"application/xml"},
		},
	}
	result := &stubResult{}
	err = c.unmarshalOutput(result, output, unmarshalBodyXml)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "XML syntax error on line 1")
}
