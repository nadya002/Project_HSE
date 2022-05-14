package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func pipe() (int, int, error) {
	var result [2]int
	if err := syscall.Pipe(result[0:]); err != nil {
		return 0, 0, err
	}
	return result[0], result[1], nil
}

func main() {
	file, err := syscall.Open("bin/xx.mp3", syscall.O_RDONLY,0)
	if err != nil {
		log.Fatalln(err)
	}
	r1, w1, _ := pipe()
	args :=  []string{"/usr/local/bin/ffmpeg",
		"-f",
		"mp3",
		"-i",
		fmt.Sprintf("pipe:%v", file),
		"-filter_complex",
		"amerge=inputs=1",
		"-ac",
		"2",
		"-f",
		"mp3",
		fmt.Sprintf("pipe:%v", w1),
	}
	fmt.Println(args)
	_, err = syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: []uintptr{0, 1, 2, uintptr(file), uintptr(w1)},
	})
	fmt.Println("FF")
	if err != nil {
		log.Fatalln(err)
	}
	res, err := os.Create("res.mp3")
	if err != nil {
		log.Fatalln(err)
	}
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	go func() {
		<- stopSignal
		res.Close()
		os.Exit(0)
	}()
	for {
		buf := make([]byte, 4096)
		n, err := syscall.Read(r1, buf)
		res.Write(buf[:n])
		fmt.Printf("New bytes %d\n", n)
		if err != nil {
			break
		}
	}

}
