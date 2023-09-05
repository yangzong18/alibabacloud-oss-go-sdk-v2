package util

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

func init() {
	NowTime = time.Now
	Sleep = time.Sleep
	SleepWithContext = sleepWithContext

	for i := 0; i < len(noEscape); i++ {
		noEscape[i] = (i >= 'A' && i <= 'Z') ||
			(i >= 'a' && i <= 'z') ||
			(i >= '0' && i <= '9') ||
			i == '-' ||
			i == '.' ||
			i == '_' ||
			i == '~'
	}
}

var NowTime func() time.Time

var Sleep func(time.Duration)

var SleepWithContext func(context.Context, time.Duration) error

var noEscape [256]bool

func sleepWithContext(ctx context.Context, dur time.Duration) error {
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// GetNowSec returns Unix time, the number of seconds elapsed since January 1, 1970 UTC.
// gets the current time in Unix time, in seconds.
func GetNowSec() int64 {
	return time.Now().Unix()
}

// GetNowNanoSec returns t as a Unix time, the number of nanoseconds elapsed
// since January 1, 1970 UTC. The result is undefined if the Unix time
// in nanoseconds cannot be represented by an int64. Note that this
// means the result of calling UnixNano on the zero Time is undefined.
// gets the current time in Unix time, in nanoseconds.
func GetNowNanoSec() int64 {
	return time.Now().UnixNano()
}

// GetNowGMT gets the current time in GMT format.
func GetNowGMT() string {
	return time.Now().UTC().Format(http.TimeFormat)
}

func EscapePath(path string, encodeSep bool) string {
	var buf bytes.Buffer
	for i := 0; i < len(path); i++ {
		c := path[i]
		if noEscape[c] || (c == '/' && !encodeSep) {
			buf.WriteByte(c)
		} else {
			fmt.Fprintf(&buf, "%%%02X", c)
		}
	}
	return buf.String()
}
