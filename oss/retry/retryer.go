package retry

import (
	"fmt"
	"time"
)

type RetryMode string

const (
	RetryModeStandard RetryMode = "standard"
)

func ParseRetryMode(v string) (mode RetryMode, err error) {
	switch v {
	case "standard":
		return RetryModeStandard, nil
	default:
		return mode, fmt.Errorf("unknown RetryMode, %v", v)
	}
}

func (m RetryMode) String() string { return string(m) }

type Retryer interface {
	IsErrorRetryable(error) bool
	MaxAttempts() int
	RetryDelay(attempt int, opErr error) (time.Duration, error)
}

type NopRetryer struct{}

func (NopRetryer) IsErrorRetryable(error) bool { return false }

func (NopRetryer) MaxAttempts() int { return 1 }

func (NopRetryer) RetryDelay(int, error) (time.Duration, error) {
	return 0, fmt.Errorf("not retrying any attempt errors")
}

type RetryOptions struct {
	MaxAttempts     int
	MaxBackoff      time.Duration
	BaseDelay       time.Duration
	Backoff         BackoffDelayer
	ErrorRetryables []ErrorRetryable
}

type BackoffDelayer interface {
	BackoffDelay(attempt int, err error) (time.Duration, error)
}

type ErrorRetryable interface {
	IsErrorRetryable(error) bool
}
