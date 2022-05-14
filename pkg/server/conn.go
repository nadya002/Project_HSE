package server

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type Conn interface {
	Fd() uintptr
	Send(buf []byte) error
	Recv(buf []byte) (int, error)
	Close() error
}

type connImpl struct {
	fd int
}

func (c *connImpl) Fd() uintptr {
	return uintptr(c.fd)
}

func (c *connImpl) Send(buf []byte) error {
	if err := unix.Send(c.fd, buf, unix.MSG_NOSIGNAL); err != nil {
		return errors.Wrapf(err, "failed to send to connection (%v)", c.fd)
	}
	return nil
}

func (c *connImpl) Recv(buf []byte) (int, error) {
	if n, err := unix.Read(c.fd, buf); err != nil {
		return n, errors.Wrapf(err, "failed to recv from connection (%v)", c.fd)
	} else {
		return n, nil
	}
}

func (c *connImpl) Close() error {
	if err := unix.Close(c.fd); err != nil {
		return errors.Wrapf(err, "failed to close connection (%v)", c.fd)
	}
	return nil
}
