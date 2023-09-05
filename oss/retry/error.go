package retry

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type MaxAttemptsError struct {
	Attempt int
	Err     error
}

func (e *MaxAttemptsError) Error() string {
	return fmt.Sprintf("exceeded maximum number of attempts, %d, %v", e.Attempt, e.Err)
}

func (e *MaxAttemptsError) Unwrap() error {
	return e.Err
}

type walkFunc func(error) bool

type causer interface {
	Cause() error
}
type wrapper interface {
	Unwrap() error
}

func walk(err error, f walkFunc) {
	for prev := err; err != nil; prev = err {
		if f(err) {
			return
		}

		switch e := err.(type) {
		case causer:
			err = e.Cause()
		case wrapper:
			err = e.Unwrap()
		default:
			errType := reflect.TypeOf(err)
			errValue := reflect.ValueOf(err)
			if errValue.IsValid() && errType.Kind() == reflect.Ptr {
				errType = errType.Elem()
				errValue = errValue.Elem()
			}
			if errValue.IsValid() && errType.Kind() == reflect.Struct {
				if errField := errValue.FieldByName("Err"); errField.IsValid() {
					errFieldValue := errField.Interface()
					if newErr, ok := errFieldValue.(error); ok {
						err = newErr
					}
				}
			}
		}
		if reflect.DeepEqual(err, prev) {
			break
		}
	}
}

func cause(cause error) (retriable bool, err error) {
	walk(cause, func(c error) bool {
		// Check for net error Timeout()
		if x, ok := c.(interface {
			Timeout() bool
		}); ok && x.Timeout() {
			retriable = true
		}

		// Check for net error Temporary()
		if x, ok := c.(interface {
			Temporary() bool
		}); ok && x.Temporary() {
			retriable = true
		}
		err = c
		return false
	})
	return
}

var retriableErrorStrings = []string{
	"use of closed network connection",
	"unexpected EOF reading trailer",
	"transport connection broken",
	"http: ContentLength=",
	"server closed idle connection",
	"bad record MAC",
	"stream error:",
	"tls: use of closed connection",
}

var retriableErrors = []error{
	io.EOF,
	io.ErrUnexpectedEOF,
}

func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	retriable, err := cause(err)
	if retriable {
		return true
	}

	for _, retriableErr := range retriableErrors {
		if err == retriableErr {
			return true
		}
	}

	errString := err.Error()
	for _, phrase := range retriableErrorStrings {
		if strings.Contains(errString, phrase) {
			return true
		}
	}

	return false
}

var retryErrorCodes = []int{
	400, // Bad request
	401, // Unauthorized
	408, // Request Timeout
	429, // Rate exceeded.
}

func ShouldRetryHTTP(statusCode int) bool {
	for _, e := range retryErrorCodes {
		if statusCode == e {
			return true
		}
	}
	if statusCode >= 500 {
		return true
	}
	return false
}
