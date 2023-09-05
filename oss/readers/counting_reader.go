package readers

import "io"

func NewCountingReader(in io.Reader) *CountingReader {
	return &CountingReader{in: in}
}

type CountingReader struct {
	in   io.Reader
	read uint64
}

func (cr *CountingReader) Read(b []byte) (int, error) {
	n, err := cr.in.Read(b)
	cr.read += uint64(n)
	return n, err
}

func (cr *CountingReader) BytesRead() uint64 {
	return cr.read
}
