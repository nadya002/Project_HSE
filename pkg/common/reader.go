package main

import (
	"context"
	"io"
)

type CancelableReader struct {
	reader io.Reader
	ctx    context.Context
	err    error
}

func NewCancelableReader(reader io.Reader, ctx context.Context) io.Reader {
	result := &CancelableReader{
		reader: reader,
		ctx:    ctx,
	}
	return result
}

func (reader *CancelableReader) Read(buf []byte) (int, error) {
	dataCh := make(chan int)
	go func() {
		n, err := reader.reader.Read(buf)
		reader.err = err
		dataCh <- n
		close(dataCh)
	}()
	select {
	case n := <-dataCh:
		return n, reader.err
	case <-reader.ctx.Done():
		return 0, reader.ctx.Err()
	}
}
