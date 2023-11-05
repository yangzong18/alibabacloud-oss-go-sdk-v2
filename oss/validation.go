package oss

import (
	"fmt"
	"net/url"
	"strings"
)

type InvalidParamError interface {
	error
	Field() string
	SetContext(string)
}

type invalidParamError struct {
	context string
	field   string
	reason  string
}

func (e invalidParamError) Error() string {
	return fmt.Sprintf("%s, %s.", e.reason, e.Field())
}

func (e invalidParamError) Field() string {
	sb := &strings.Builder{}
	sb.WriteString(e.context)
	if sb.Len() > 0 {
		sb.WriteRune('.')
	}
	sb.WriteString(e.field)
	return sb.String()
}

func (e *invalidParamError) SetContext(ctx string) {
	e.context = ctx
}

type ParamRequiredError struct {
	invalidParamError
}

func NewErrParamRequired(field string) *ParamRequiredError {
	return &ParamRequiredError{
		invalidParamError{
			field:  field,
			reason: fmt.Sprintf("missing required field"),
		},
	}
}

func NewErrParamInvalid(field string) *ParamRequiredError {
	return &ParamRequiredError{
		invalidParamError{
			field:  field,
			reason: fmt.Sprintf("invalid field"),
		},
	}
}

func NewErrParamNull(field string) *ParamRequiredError {
	return &ParamRequiredError{
		invalidParamError{
			field:  field,
			reason: fmt.Sprintf("null field"),
		},
	}
}

func isValidEndpoint(endpoint *url.URL) bool {
	return (endpoint != nil)
}

func isValidBucketName(bucketName *string) bool {
	if bucketName == nil {
		return false
	}

	nameLen := len(*bucketName)
	if nameLen < 3 || nameLen > 63 {
		return false
	}

	if (*bucketName)[0] == '-' || (*bucketName)[nameLen-1] == '-' {
		return false
	}

	for _, v := range *bucketName {
		if !(('a' <= v && v <= 'z') || ('0' <= v && v <= '9') || v == '-') {
			return false
		}
	}
	return true
}

func isValidObjectName(objectName *string) bool {
	if objectName == nil || len(*objectName) == 0 {
		return false
	}
	return true
}

var supportedMethod = map[string]interface{}{
	"GET":    nil,
	"PUT":    nil,
	"POST":   nil,
	"DELETE": nil,
	"OPTION": nil,
}

func isValidMethod(method string) bool {
	if _, ok := supportedMethod[method]; ok {
		return true
	}
	return false
}
