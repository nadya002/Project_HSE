package server

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"syscall"
)

const ListenBacklog = 5

type Server interface {
	Listen(host string, port int) error
	Accept() (chan Conn, chan error)
	Close() error
}

func New() (Server, error) {
	return (&serverImpl{
		closed: false,
	}).init()
}

type serverImpl struct {
	sockFd int
	closed bool
}

func (s *serverImpl) Listen(host string, port int) error {
	if addr, err := getAddr(host, port); err != nil {
		return err
	} else {
		if err := unix.Bind(s.sockFd, addr); err != nil {
			return errors.Wrapf(err, "failed to bind server to addr, (close socket result: %v)", s.Close())
		}
		if err := unix.Listen(s.sockFd, ListenBacklog); err != nil {
			return errors.Wrapf(err, "failed to listen, (close socket result: %v)", s.Close())
		}
		return nil
	}
}

func (s *serverImpl) Accept() (chan Conn, chan error) {
	resultCh := make(chan Conn)
	errCh := make(chan error)
	go func() {
		if fd, _, err := unix.Accept(s.sockFd); err != nil {
			errCh <- errors.Wrapf(err, "failed to accept new connection")
			close(resultCh)
		} else {
			resultCh <- &connImpl{fd: fd}
			close(errCh)
		}
	}()
	return resultCh, errCh
}

func (s *serverImpl) Close() error {
	if err := unix.Close(s.sockFd); err != nil {
		return errors.Wrapf(err, "failed to close server socket (%v)", s.sockFd)
	}
	return nil
}

func (s *serverImpl) init() (*serverImpl, error) {
	var err error
	syscall.ForkLock.Lock()
	defer syscall.ForkLock.Unlock()
	s.sockFd, err = unix.Socket(
		unix.AF_INET,
		unix.SOCK_STREAM,
		unix.IPPROTO_IP,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create server")
	}
	return s, nil
}
