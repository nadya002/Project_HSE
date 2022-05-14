package server

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

func getAddr(host string, port int) (*unix.SockaddrInet4, error) {
	addrSlices := strings.Split(host, ".")
	var addr [4]byte
	if len(addrSlices) == 0 {
		addr[0] = 0
		addr[1] = 0
		addr[2] = 0
		addr[3] = 0
	} else if len(addrSlices) == 4 {
		for index, slice := range addrSlices {
			if value, err := strconv.Atoi(slice); err != nil {
				return nil, errors.Errorf("invalid host %s", host)
			} else {
				addr[index] = byte(value)
			}
		}
	} else {
		return nil, errors.Errorf("invalid host %s", host)
	}
	return &unix.SockaddrInet4{
		Port: port,
		Addr: addr,
	}, nil
}
