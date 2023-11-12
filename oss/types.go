package oss

import (
	"io"
	"net/http"
)

type OperationMetadata struct {
	values map[interface{}]interface{}
}

func (m OperationMetadata) Get(key interface{}) interface{} {
	return m.values[key]
}

func (m OperationMetadata) Clone() OperationMetadata {
	vs := make(map[interface{}]interface{}, len(m.values))
	for k, v := range m.values {
		vs[k] = v
	}

	return OperationMetadata{
		values: vs,
	}
}

func (m *OperationMetadata) Set(key, value interface{}) {
	if m.values == nil {
		m.values = map[interface{}]interface{}{}
	}
	m.values[key] = value
}

func (m OperationMetadata) Has(key interface{}) bool {
	if m.values == nil {
		return false
	}
	_, ok := m.values[key]
	return ok
}

type RequestCommon struct {
	Headers    map[string]string
	Parameters map[string]string
	Body       io.Reader
}

type RequestCommonInterface interface {
	GetCommonFileds() (map[string]string, map[string]string, io.Reader)
}

func (r *RequestCommon) GetCommonFileds() (map[string]string, map[string]string, io.Reader) {
	return r.Headers, r.Parameters, r.Body
}

type ResultCommon struct {
	Status     string
	StatusCode int
	Headers    http.Header
	OpMetadata OperationMetadata
}

type ResultCommonInterface interface {
	CopyIn(status string, statusCode int, headers http.Header, meta OperationMetadata)
}

func (r *ResultCommon) CopyIn(status string, statusCode int, headers http.Header, meta OperationMetadata) {
	r.Status = status
	r.StatusCode = statusCode
	r.Headers = headers
	r.OpMetadata = meta
}

type OperationInput struct {
	OpName     string
	Method     string
	Headers    map[string]string
	Parameters map[string]string
	Body       io.Reader

	Bucket *string
	Key    *string

	OpMetadata OperationMetadata
}

type OperationOutput struct {
	Input *OperationInput

	Status     string
	StatusCode int
	Headers    http.Header
	Body       io.ReadCloser

	OpMetadata OperationMetadata

	httpRequest *http.Request
}
