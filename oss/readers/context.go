package readers

import (
	"context"
	"io"
)

func NewContextReader(ctx context.Context, r io.Reader) io.Reader {
	return &contextReader{
		ctx: ctx,
		r:   r,
	}
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr *contextReader) Read(p []byte) (n int, err error) {
	err = cr.ctx.Err()
	if err != nil {
		return 0, err
	}
	return cr.r.Read(p)
}
