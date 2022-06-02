package pipe

import "syscall"

func Pipe() (int, int, error) {
	var result [2]int
	if err := syscall.Pipe(result[0:]); err != nil {
		return 0, 0, err
	}
	return result[0], result[1], nil
}