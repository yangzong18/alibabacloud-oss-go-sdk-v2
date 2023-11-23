package oss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

func ReadSeekNopCloser(r io.Reader) ReadSeekerNopClose {
	return ReadSeekerNopClose{r}
}

type ReadSeekerNopClose struct {
	r io.Reader
}

func (r ReadSeekerNopClose) Read(p []byte) (int, error) {
	switch t := r.r.(type) {
	case io.Reader:
		return t.Read(p)
	}
	return 0, nil
}

func (r ReadSeekerNopClose) Seek(offset int64, whence int) (int64, error) {
	switch t := r.r.(type) {
	case io.Seeker:
		return t.Seek(offset, whence)
	}
	return int64(0), nil
}

func (r ReadSeekerNopClose) Close() error {
	return nil
}

func (r ReadSeekerNopClose) IsSeeker() bool {
	_, ok := r.r.(io.Seeker)
	return ok
}

func (r ReadSeekerNopClose) HasLen() (int, bool) {
	type lenner interface {
		Len() int
	}

	if lr, ok := r.r.(lenner); ok {
		return lr.Len(), true
	}

	return 0, false
}

func (r ReadSeekerNopClose) GetLen() (int64, error) {
	if l, ok := r.HasLen(); ok {
		return int64(l), nil
	}

	if s, ok := r.r.(io.Seeker); ok {
		return seekerLen(s)
	}

	return -1, nil
}

func seekerLen(s io.Seeker) (int64, error) {
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}

func isReaderSeekable(r io.Reader) bool {
	switch v := r.(type) {
	case ReadSeekerNopClose:
		return v.IsSeeker()
	case *ReadSeekerNopClose:
		return v.IsSeeker()
	case io.ReadSeeker:
		return true
	default:
		return false
	}
}

type buffer struct {
	buf    []byte
	err    error
	offset int
}

func (b *buffer) isEmpty() bool {
	if b == nil {
		return true
	}
	if len(b.buf)-b.offset <= 0 {
		return true
	}
	return false
}

func (b *buffer) read(rd io.Reader) error {
	var n int
	n, b.err = readFill(rd, b.buf)
	b.buf = b.buf[0:n]
	b.offset = 0
	return b.err
}

func (b *buffer) buffer() []byte {
	return b.buf[b.offset:]
}

func (b *buffer) increment(n int) {
	b.offset += n
}

const (
	AsyncReadeBufferSize = 1024 * 1024
	softStartInitial     = 512 * 1024
)

type AsyncReaderRangeGet func(context.Context, HTTPRange) (r io.ReadCloser, offset int64, etag string, err error)

type AsyncReader struct {
	in      io.ReadCloser // Input reader
	ready   chan *buffer  // Buffers ready to be handed to the reader
	token   chan struct{} // Tokens which allow a buffer to be taken
	exit    chan struct{} // Closes when finished
	buffers int           // Number of buffers
	err     error         // If an error has occurred it is here
	cur     *buffer       // Current buffer being served
	exited  chan struct{} // Channel is closed been the async reader shuts down
	size    int           // size of buffer to use
	closed  bool          // whether we have closed the underlying stream
	mu      sync.Mutex    // lock for Read/WriteTo/Abandon/Close

	//Range Getter
	rangeGet  AsyncReaderRangeGet
	httpRange HTTPRange

	// For reader
	offset  int64
	gotsize int64

	oriHttpRange HTTPRange

	context context.Context
	cancel  context.CancelFunc

	// Origin file pattern
	etag    string
	modTime string
}

// NewAsyncReader returns a reader that will asynchronously read from
// the Reader returued by getter from the given offset into a number of buffers each of size AsyncReadeBufferSize
// The input can be read from the returned reader.
// When done use Close to release the buffers and close the supplied input.
// The etag is used to identify the content of the object. If not set, the first ETag returned value will be used instead.
func NewAsyncReader(ctx context.Context,
	rangeGet AsyncReaderRangeGet, httpRange *HTTPRange, etag string, buffers int) (*AsyncReader, error) {

	if buffers <= 0 {
		return nil, errors.New("number of buffers too small")
	}
	if rangeGet == nil {
		return nil, errors.New("nil reader supplied")
	}

	context, cancel := context.WithCancel(ctx)

	range_ := HTTPRange{}
	if httpRange != nil {
		range_ = *httpRange
	}

	a := &AsyncReader{
		rangeGet:     rangeGet,
		context:      context,
		cancel:       cancel,
		httpRange:    range_,
		oriHttpRange: range_,
		offset:       range_.Offset,
		gotsize:      0,
		etag:         etag,
		buffers:      buffers,
	}

	//fmt.Printf("NewAsyncReader, range: %s, etag:%s, buffer count:%v\n", ToString(a.httpRange.FormatHTTPRange()), a.etag, a.buffers)

	a.init(buffers)
	return a, nil
}

func (a *AsyncReader) init(buffers int) {
	a.ready = make(chan *buffer, buffers)
	a.token = make(chan struct{}, buffers)
	a.exit = make(chan struct{})
	a.exited = make(chan struct{})
	a.buffers = buffers
	a.cur = nil
	a.size = softStartInitial

	// Create tokens
	for i := 0; i < buffers; i++ {
		a.token <- struct{}{}
	}

	// Start async reader
	go func() {
		// Ensure that when we exit this is signalled.
		defer close(a.exited)
		defer close(a.ready)
		for {
			select {
			case <-a.token:
				b := a.getBuffer()
				if a.size < AsyncReadeBufferSize {
					b.buf = b.buf[:a.size]
					a.size <<= 1
				}

				if a.httpRange.Count > 0 && a.gotsize > a.httpRange.Count {
					b.buf = b.buf[0:0]
					b.err = io.EOF
					//fmt.Printf("a.gotsize > a.httpRange.Count, err:%v\n", b.err)
					a.ready <- b
					return
				}

				if a.in == nil {
					body, off, etag, err := a.rangeGet(a.context, a.httpRange)
					if a.etag == "" {
						a.etag = etag
					}
					if err == nil {
						if off != a.httpRange.Offset {
							err = fmt.Errorf("Range get fail, expect offset:%v, got offset:%v", a.httpRange.Offset, off)
						}
						if etag != a.etag {
							err = fmt.Errorf("Source file is changed, expect etag:%s ,got offset:%s", a.etag, etag)
						}
					}
					if err != nil {
						b.buf = b.buf[0:0]
						b.err = err
						//fmt.Printf("call getFunc fail, err:%v\n", err)
						a.ready <- b
						return
					}
					a.in = body
					//fmt.Printf("call getFunc done, range:%s\n", ToString(a.httpRange.FormatHTTPRange()))
				}

				// ignore err from read
				err := b.read(a.in)
				a.httpRange.Offset += int64(len(b.buf))
				a.gotsize += int64(len(b.buf))
				if err != io.EOF {
					b.err = nil
				}
				//fmt.Printf("read into buffer, size:%v, next begin:%v, err:%v\n", len(b.buf), a.httpRange.Offset, err)
				a.ready <- b
				if err != nil {
					a.in.Close()
					a.in = nil
					if err == io.EOF {
						return
					}
				}
			case <-a.exit:
				return
			}
		}
	}()
}

func (a *AsyncReader) fill() (err error) {
	if a.cur.isEmpty() {
		if a.cur != nil {
			a.putBuffer(a.cur)
			a.token <- struct{}{}
			a.cur = nil
		}
		b, ok := <-a.ready
		if !ok {
			// Return an error to show fill failed
			if a.err == nil {
				return errors.New("stream abandoned")
			}
			return a.err
		}
		a.cur = b
	}
	return nil
}

// Read will return the next available data.
func (a *AsyncReader) Read(p []byte) (n int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	defer func() {
		a.offset += int64(n)
	}()

	// Swap buffer and maybe return error
	err = a.fill()
	if err != nil {
		return 0, err
	}

	// Copy what we can
	n = copy(p, a.cur.buffer())
	a.cur.increment(n)

	// If at end of buffer, return any error, if present
	if a.cur.isEmpty() {
		a.err = a.cur.err
		return n, a.err
	}
	return n, nil
}

func (a *AsyncReader) Close() (err error) {
	a.abandon()
	if a.closed {
		return nil
	}
	a.closed = true

	if a.in != nil {
		err = a.in.Close()
	}
	return
}

func (a *AsyncReader) abandon() {
	a.stop()
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cur != nil {
		a.putBuffer(a.cur)
		a.cur = nil
	}
	for b := range a.ready {
		a.putBuffer(b)
	}
}

func (a *AsyncReader) stop() {
	select {
	case <-a.exit:
		return
	default:
	}
	a.cancel()
	close(a.exit)
	<-a.exited
}

// bufferPool is a global pool of buffers
var bufferPool *sync.Pool
var bufferPoolOnce sync.Once

// TODO use pool
func (a *AsyncReader) putBuffer(b *buffer) {
	b.buf = b.buf[0:cap(b.buf)]
	bufferPool.Put(b.buf)
}

func (a *AsyncReader) getBuffer() *buffer {
	bufferPoolOnce.Do(func() {
		// Initialise the buffer pool when used
		bufferPool = &sync.Pool{
			New: func() any {
				//fmt.Printf("make([]byte, BufferSize)\n")
				return make([]byte, AsyncReadeBufferSize)
			},
		}
	})
	return &buffer{
		buf: bufferPool.Get().([]byte),
	}
}

func readFill(r io.Reader, buf []byte) (n int, err error) {
	var nn int
	for n < len(buf) && err == nil {
		nn, err = r.Read(buf[n:])
		n += nn
	}
	return n, err
}
