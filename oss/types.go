package oss

import (
	"io"
	"net/http"
)

type OperationMetadata struct {
	values map[any][]any
}

func (m OperationMetadata) Get(key any) any {
	if m.values == nil {
		return nil
	}
	v := m.values[key]
	if len(v) == 0 {
		return nil
	}
	return v[0]
}

func (m OperationMetadata) Values(key any) []any {
	if m.values == nil {
		return nil
	}
	return m.values[key]
}

func (m *OperationMetadata) Add(key, value any) {
	if m.values == nil {
		m.values = map[any][]any{}
	}
	m.values[key] = append(m.values[key], value)
}

func (m *OperationMetadata) Set(key, value any) {
	if m.values == nil {
		m.values = map[any][]any{}
	}
	m.values[key] = []any{value}
}

func (m OperationMetadata) Has(key any) bool {
	if m.values == nil {
		return false
	}
	_, ok := m.values[key]
	return ok
}

type RequestCommon struct {
	Headers    map[string]string
	Parameters map[string]string
	Payload    io.Reader
}

type RequestCommonInterface interface {
	GetCommonFileds() (map[string]string, map[string]string, io.Reader)
}

func (r *RequestCommon) GetCommonFileds() (map[string]string, map[string]string, io.Reader) {
	return r.Headers, r.Parameters, r.Payload
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

type RequestBodyTracker interface {
	io.Writer
	Reset()
}
