package retry

import (
	"time"
)

const (
	DefaultMaxAttempts int = 3

	DefaultMaxBackoff time.Duration = 20 * time.Second
	DefaultBaseDelay  time.Duration = 200 * time.Millisecond
)

var DefaultErrorRetryables = []ErrorRetryable{}

type Standard struct {
	maxAttempts int
	retryables  []ErrorRetryable
	backoff     BackoffDelayer
}

func NewStandard(fnOpts ...func(*RetryOptions)) *Standard {
	o := RetryOptions{
		MaxAttempts:     DefaultMaxAttempts,
		MaxBackoff:      DefaultMaxBackoff,
		BaseDelay:       DefaultBaseDelay,
		ErrorRetryables: append([]ErrorRetryable{}, DefaultErrorRetryables...),
	}

	for _, fn := range fnOpts {
		fn(&o)
	}

	if o.MaxAttempts <= 0 {
		o.MaxAttempts = DefaultMaxAttempts
	}

	if o.BaseDelay <= 0 {
		o.BaseDelay = DefaultBaseDelay
	}

	if o.Backoff == nil {
		o.Backoff = NewFullJitterBackoff(o.BaseDelay, o.MaxBackoff)
	}

	return &Standard{
		maxAttempts: o.MaxAttempts,
		retryables:  o.ErrorRetryables,
		backoff:     o.Backoff,
	}
}

func (s *Standard) MaxAttempts() int {
	return s.maxAttempts
}

func (s *Standard) IsErrorRetryable(err error) bool {
	for _, re := range s.retryables {
		if v := re.IsErrorRetryable(err); v {
			return v
		}
	}
	return false
}

func (s *Standard) RetryDelay(attempt int, err error) (time.Duration, error) {
	return s.backoff.BackoffDelay(attempt, err)
}
