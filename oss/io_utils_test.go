package oss

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsyncRangeReader(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return buf, 0, "", nil
	}
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, []byte(data), dst[:n])

	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	// Test read after error
	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	err = ar.Close()
	require.NoError(t, err)
	// Test double close
	err = ar.Close()
	require.NoError(t, err)

	// Test Close without reading everything
	buf = io.NopCloser(bytes.NewBuffer(make([]byte, 50000)))
	getFn = func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return buf, 0, "", nil
	}
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)
	err = ar.Close()
	require.NoError(t, err)
}

func TestAsyncRangeReaderErrors(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"

	// test nil reader
	_, err := NewAsyncRangeReader(ctx, nil, nil, "", 4)
	require.Error(t, err)

	// invalid buffer number
	buf := io.NopCloser(bytes.NewBufferString(data))
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return buf, 0, "", nil
	}
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", 0)
	require.Error(t, err)
	_, err = NewAsyncRangeReader(ctx, getFn, nil, "", -1)
	require.Error(t, err)
}

type readMaker struct {
	name string
	fn   func(io.Reader) io.Reader
}

var readMakers = []readMaker{
	{"full", func(r io.Reader) io.Reader { return r }},
	{"byte", iotest.OneByteReader},
	{"half", iotest.HalfReader},
	{"data+err", iotest.DataErrReader},
	{"timeout", iotest.TimeoutReader},
}

// Call Read to accumulate the text of a file
func reads(buf io.Reader, m int) string {
	var b [1000]byte
	nb := 0
	for {
		n, err := buf.Read(b[nb : nb+m])
		nb += n
		if err == io.EOF {
			break
		} else if err != nil && err != iotest.ErrTimeout {
			panic("Data: " + err.Error())
		} else if err != nil {
			break
		}
	}
	return string(b[0:nb])
}

type bufReader struct {
	name string
	fn   func(io.Reader) string
}

var bufreaders = []bufReader{
	{"1", func(b io.Reader) string { return reads(b, 1) }},
	{"2", func(b io.Reader) string { return reads(b, 2) }},
	{"3", func(b io.Reader) string { return reads(b, 3) }},
	{"4", func(b io.Reader) string { return reads(b, 4) }},
	{"5", func(b io.Reader) string { return reads(b, 5) }},
	{"7", func(b io.Reader) string { return reads(b, 7) }},
}

const minReadBufferSize = 16

var bufsizes = []int{
	0, minReadBufferSize, 23, 32, 46, 64, 93, 128, 1024, 4096,
}

// Test various  input buffer sizes, number of buffers and read sizes.
func TestAsyncRangeReaderSizes(t *testing.T) {
	ctx := context.Background()
	var texts [31]string
	str := ""
	all := ""
	for i := 0; i < len(texts)-1; i++ {
		texts[i] = str + "\n"
		all += texts[i]
		str += string(rune(i)%26 + 'a')
	}
	texts[len(texts)-1] = all

	for h := 0; h < len(texts); h++ {
		text := texts[h]
		for i := 0; i < len(readMakers); i++ {
			for j := 0; j < len(bufreaders); j++ {
				for k := 0; k < len(bufsizes); k++ {
					for l := 1; l < 10; l++ {
						readmaker := readMakers[i]
						bufreader := bufreaders[j]
						bufsize := bufsizes[k]
						read := readmaker.fn(strings.NewReader(text))
						buf := bufio.NewReaderSize(read, bufsize)
						getFn := func(_ context.Context, httpRange HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
							return io.NopCloser(buf), httpRange.Offset, "", nil
						}
						ar, _ := NewAsyncRangeReader(ctx, getFn, nil, "", 1)
						s := bufreader.fn(ar)
						// "timeout" expects the Reader to recover, AsyncRangeReader does not.
						if s != text && readmaker.name != "timeout" {
							t.Errorf("reader=%s fn=%s bufsize=%d want=%q got=%q",
								readmaker.name, bufreader.name, bufsize, text, s)
						}
						err := ar.Close()
						require.NoError(t, err)
					}
				}
			}
		}
	}
}

// Read an infinite number of zeros
type zeroReader struct {
	closed bool
}

func (z *zeroReader) Read(p []byte) (n int, err error) {
	if z.closed {
		return 0, io.EOF
	}
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

func (z *zeroReader) Close() error {
	if z.closed {
		panic("double close on zeroReader")
	}
	z.closed = true
	return nil
}

// Test closing and abandoning
func TestAsyncRangeReaderClose(t *testing.T) {
	ctx := context.Background()
	zr := &zeroReader{}
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return zr, 0, "", nil
	}
	a, err := NewAsyncRangeReader(ctx, getFn, nil, "", 16)
	require.NoError(t, err)
	var copyN int64
	var copyErr error
	var wg sync.WaitGroup
	started := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		close(started)
		{
			// exercise the Read path
			buf := make([]byte, 64*1024)
			for {
				var n int
				n, copyErr = a.Read(buf)
				copyN += int64(n)
				if copyErr != nil {
					break
				}
			}
		}
	}()
	// Do some copying
	<-started
	time.Sleep(100 * time.Millisecond)
	// abandon the copy
	a.abandon()
	wg.Wait()
	assert.Contains(t, copyErr.Error(), "stream abandoned")
	// t.Logf("Copied %d bytes, err %v", copyN, copyErr)
	assert.True(t, copyN > 0)
}

func TestAsyncRangeReaderEtagCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return io.NopCloser(bytes.NewBufferString(data)), 0, "etag", nil
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[0:n]))

	// set etag to "etag"
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 10, n)
	assert.Equal(t, data, string(dst[:n]))

	// set etag to "invalid-etag"
	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "invalid-etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Contains(t, err.Error(), "Source file is changed, expect etag:invalid-etag")
}

func TestAsyncRangeReaderOffsetCheck(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBufferString(data)))), 0, "etag", nil
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = ar.Read(dst)
	assert.Equal(t, 0, n)
	assert.Contains(t, err.Error(), "Range get fail, expect offset")

	//
	getFn = func(ctx context.Context, range_ HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			return io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))), 0, "etag", nil
		} else {
			return io.NopCloser(bytes.NewBuffer(b[range_.Offset:])), range_.Offset, "etag", nil
		}
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)

	dst = make([]byte, 100)
	n, err = ar.Read(dst)
	assert.Equal(t, 1, n)
	n, err = ar.Read(dst[n:])
	assert.Equal(t, 9, n)
	assert.Equal(t, data, string(dst[:10]))
}

func TestAsyncRangeReaderGetError(t *testing.T) {
	ctx := context.Background()
	data := "Testbuffer"
	getFn := func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		return nil, 0, "", errors.New("range get fail")
	}

	// don't set etag
	ar, err := NewAsyncRangeReader(ctx, getFn, nil, "", 4)
	require.NoError(t, err)

	var dst = make([]byte, 100)
	_, err = ar.Read(dst)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")

	getFn = func(ctx context.Context, range_ HTTPRange) (r io.ReadCloser, offset int64, etag string, err error) {
		b := []byte(data)
		if range_.Offset == 0 {
			return io.NopCloser(iotest.TimeoutReader(iotest.OneByteReader(bytes.NewBuffer(b[range_.Offset:])))), 0, "etag", nil
		} else {
			return nil, 0, "", errors.New("range get fail")
		}
	}

	ar, err = NewAsyncRangeReader(ctx, getFn, nil, "etag", 4)
	require.NoError(t, err)
	dst = make([]byte, 100)
	n, err := ar.Read(dst)
	assert.Equal(t, 1, n)
	assert.Equal(t, int64(1), ar.offset)
	n, err = ar.Read(dst[n:])
	require.Error(t, err)
	assert.Contains(t, err.Error(), "range get fail")
}
