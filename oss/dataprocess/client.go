package dataprocess

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
)

// Client is the client for accessing OSS Data Process API
type Client struct {
	client *oss.Client
}

// NewClient creates a new DataProcess client with the given configuration
func NewClient(cfg *oss.Config, optFns ...func(*oss.Options)) *Client {
	return &Client{
		client: oss.NewClient(cfg, optFns...),
	}
}

// Unwrap returns the underlying OSS client
func (c *Client) Unwrap() *oss.Client { return c.client }

// toClientError converts an error to a client error
func (c *Client) toClientError(err error, code string, output *oss.OperationOutput) error {
	if err == nil {
		return nil
	}

	return &oss.ClientError{
		Code: code,
		Message: fmt.Sprintf("execute %s fail, error code is %s, request id:%s",
			output.Input.OpName,
			code,
			output.Headers.Get(oss.HeaderOssRequestID),
		),
		Err: err,
	}
}

func unmarshalBodyXmlMix(result any, output *oss.OperationOutput) error {
	var err error
	var body []byte
	if output.Body != nil {
		defer output.Body.Close()
		if body, err = io.ReadAll(output.Body); err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return nil
	}

	val := reflect.ValueOf(result)
	switch val.Kind() {
	case reflect.Pointer, reflect.Interface:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct || output == nil {
		return nil
	}

	t := val.Type()
	idx := -1
	for k := 0; k < t.NumField(); k++ {
		if tag, ok := t.Field(k).Tag.Lookup("output"); ok {
			tokens := strings.Split(tag, ",")
			if len(tokens) < 2 {
				continue
			}
			// header|query|body,filed_name,[required,time,usermeta...]
			switch tokens[0] {
			case "body":
				idx = k
				break
			}
		}
	}

	if idx >= 0 {
		dst := val.Field(idx)
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		err = xml.Unmarshal(body, dst.Interface())
	} else {
		err = xml.Unmarshal(body, result)
	}

	if err != nil {
		err = &oss.DeserializationError{
			Err:      err,
			Snapshot: body,
		}
	}

	return err
}
