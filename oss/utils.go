package oss

import (
	"bytes"
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	escQuot = []byte("&#34;") // shorter than "&quot;"
	escApos = []byte("&#39;") // shorter than "&apos;"
	escAmp  = []byte("&amp;")
	escLT   = []byte("&lt;")
	escGT   = []byte("&gt;")
	escTab  = []byte("&#x9;")
	escNL   = []byte("&#xA;")
	escCR   = []byte("&#xD;")
	escFFFD = []byte("\uFFFD") // Unicode replacement character
)

func init() {
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

// getNowSec returns Unix time, the number of seconds elapsed since January 1, 1970 UTC.
// gets the current time in Unix time, in seconds.
func getNowSec() int64 {
	return time.Now().Unix()
}

// getNowGMT gets the current time in GMT format.
func getNowGMT() string {
	return time.Now().UTC().Format(http.TimeFormat)
}

func escapePath(path string, encodeSep bool) string {
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

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Pointer:
		return v.IsNil()
	}
	return false
}

func setTimeReflectValue(dst reflect.Value, value time.Time) (err error) {
	dst0 := dst
	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}
	if dst.CanAddr() {
		pv := dst.Addr()
		if pv.CanInterface() {
			if val, ok := pv.Interface().(encoding.TextUnmarshaler); ok {
				return val.UnmarshalText([]byte(value.Format(time.RFC3339)))
			}
		}
	}
	return errors.New("cannot unmarshal into " + dst0.Type().String())
}

func setReflectValue(dst reflect.Value, data string) (err error) {
	dst0 := dst
	src := []byte(data)

	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	switch dst.Kind() {
	case reflect.Invalid:
	default:
		return errors.New("cannot unmarshal into " + dst0.Type().String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(src) == 0 {
			dst.SetInt(0)
			return nil
		}
		itmp, err := strconv.ParseInt(strings.TrimSpace(string(src)), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetInt(itmp)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if len(src) == 0 {
			dst.SetUint(0)
			return nil
		}
		utmp, err := strconv.ParseUint(strings.TrimSpace(string(src)), 10, dst.Type().Bits())
		if err != nil {
			return err
		}
		dst.SetUint(utmp)
	case reflect.Bool:
		if len(src) == 0 {
			dst.SetBool(false)
			return nil
		}
		value, err := strconv.ParseBool(strings.TrimSpace(string(src)))
		if err != nil {
			return err
		}
		dst.SetBool(value)
	case reflect.String:
		dst.SetString(string(src))
	}
	return nil
}

func setMapStringReflectValue(dst reflect.Value, key interface{}, data interface{}) (err error) {
	dst0 := dst

	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}
		dst = dst.Elem()
	}

	switch dst.Kind() {
	case reflect.Invalid:
	default:
		return errors.New("cannot unmarshal into " + dst0.Type().String())
	case reflect.Map:
		if dst.IsNil() {
			dst.Set(reflect.MakeMap(dst.Type()))
		}
		mapValue := reflect.ValueOf(data)
		mapKey := reflect.ValueOf(key)
		dst.SetMapIndex(mapKey, mapValue)
	}
	return nil
}

func defaultUserAgent() string {
	return fmt.Sprintf("aliyun-sdk-go/%s (%s/%s/%s;%s)", Version(), runtime.GOOS,
		"-", runtime.GOARCH, runtime.Version())
}

func isContextError(ctx context.Context, perr *error) bool {
	if ctxErr := ctx.Err(); ctxErr != nil {
		if *perr == nil {
			*perr = ctxErr
		}
		return true
	}
	return false
}

func copySeekableBody(dst io.Writer, src io.ReadSeeker) (int64, error) {
	curPos, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(dst, src)
	if err != nil {
		return n, err
	}

	_, err = src.Seek(curPos, io.SeekStart)
	if err != nil {
		return n, err
	}

	return n, nil
}

func parseOffsetAndSizeFromHeaders(headers http.Header) (offset, size int64) {
	size = -1
	var contentLength = headers.Get("Content-Length")
	if len(contentLength) != 0 {
		var err error
		if size, err = strconv.ParseInt(contentLength, 10, 64); err != nil {
			return 0, -1
		}
	}

	var contentRange = headers.Get("Content-Range")
	if len(contentRange) == 0 {
		return 0, size
	}

	if !strings.HasPrefix(contentRange, "bytes ") {
		return 0, -1
	}

	// start offset
	dash := strings.IndexRune(contentRange, '-')
	if dash < 0 {
		return 0, -1
	}
	ret, err := strconv.ParseInt(contentRange[6:dash], 10, 64)
	if err != nil {
		return 0, -1
	}
	offset = ret

	// total size
	slash := strings.IndexRune(contentRange, '/')
	if slash < 0 {
		return 0, -1
	}
	ret, err = strconv.ParseInt(contentRange[slash+1:], 10, 64)
	if err != nil {
		return 0, -1
	}
	size = ret

	return offset, size
}

// escapeXml EscapeString writes to p the properly escaped XML equivalent
// of the plain text data s.
func escapeXml(s string) string {
	var p strings.Builder
	var esc []byte
	hextable := "0123456789ABCDEF"
	escPattern := []byte("&#x00;")
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
		case '"':
			esc = escQuot
		case '\'':
			esc = escApos
		case '&':
			esc = escAmp
		case '<':
			esc = escLT
		case '>':
			esc = escGT
		case '\t':
			esc = escTab
		case '\n':
			esc = escNL
		case '\r':
			esc = escCR
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				if r >= 0x00 && r < 0x20 {
					escPattern[3] = hextable[r>>4]
					escPattern[4] = hextable[r&0x0f]
					esc = escPattern
				} else {
					esc = escFFFD
				}
				break
			}
			continue
		}
		p.WriteString(s[last : i-width])
		p.Write(esc)
		last = i
	}
	p.WriteString(s[last:])
	return p.String()
}

// Decide whether the given rune is in the XML Character Range, per
// the Char production of https://www.xml.com/axml/testaxml.htm,
// Section 2.2 Characters.
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xD7FF ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}
