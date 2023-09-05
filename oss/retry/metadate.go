package retry

import (
	"time"
)

type AttemptResults struct {
	Results []AttemptResult
}

type AttemptResult struct {
	Err       error
	Retryable bool
	Retried   bool
}

type RetryMetadata struct {
	AttemptNum       int
	AttemptTime      time.Time
	MaxAttempts      int
	AttemptClockSkew time.Duration
}
