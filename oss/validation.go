package oss

import (
	"bytes"
	"fmt"
	"strings"
)

type InvalidParamsError struct {
	Context string
	errs    []InvalidParamError
}

func (e *InvalidParamsError) Add(err InvalidParamError) {
	err.SetContext(e.Context)
	e.errs = append(e.errs, err)
}

func (e *InvalidParamsError) Len() int {
	return len(e.errs)
}

func (e InvalidParamsError) Error() string {
	w := &bytes.Buffer{}
	fmt.Fprintf(w, "%d validation error(s) found.\n", len(e.errs))

	for _, err := range e.errs {
		fmt.Fprintf(w, "- %s\n", err.Error())
	}

	return w.String()
}

func (e InvalidParamsError) Errs() []error {
	errs := make([]error, len(e.errs))
	for i := 0; i < len(errs); i++ {
		errs[i] = e.errs[i]
	}

	return errs
}

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
