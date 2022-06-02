package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"hse/pkg/pipe"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func getConnection(addr string, inputFile *os.File) net.Conn {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := io.Copy(conn, inputFile); err != nil {
		log.Fatalln(err)
	}
	return conn
}

func generateReadPipe(conn net.Conn) uintptr {
	r, w, err := pipe.Pipe()
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		if _, err := io.Copy(os.NewFile(uintptr(w), "pipe"), conn); err != nil {
			log.Fatalln(err)
		}
	}()
	return uintptr(r)
}


func main() {
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.Int("port", 7777, "listen port")
	input := flag.String("input", "", "input file (WAV)")
	output := flag.String("output", "", "output file (WAV)")
	denoiserLocation := flag.String("denoiser", "/bin/denoiser", "denoiser_location")
	ffmpeg := flag.String("ffmpeg", "/usr/local/bin/ffmpeg", "ffmpeg_location")
	flag.Parse()

	inputFile, err := os.Open(*input)
	defer inputFile.Close()
	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("%v:%v", *host, *port)
	conn := getConnection(addr, inputFile)

	readPipe := generateReadPipe(conn)

	r, w, err := pipe.Pipe()
	if err != nil {
		log.Fatalln(err)
	}

	args := []string{*ffmpeg}
	args = append(args, "-f", "wav", "-i", fmt.Sprintf("pipe:%v", readPipe))
	args = append(args, "-f", "s16le", "-acodec", "pcm_s16le", "-ar", "48k", fmt.Sprintf("pipe:%v", w))

	if _, err = syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: append([]uintptr{0, 1, 2, readPipe},),
	}); err != nil {
		log.Fatalln(errors.Wrap(err, "failed to start ffmpeg"))
	}

	args = []string{*denoiserLocation}

	r1, w1, err := pipe.Pipe()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: append([]uintptr{uintptr(r), uintptr(w1), 2},),
	}); err != nil {
		log.Fatalln(errors.Wrap(err, "failed to start denoiser"))
	}

	*denoiserLocation = ""

	args = []string{*ffmpeg}
	args = append(args, "-f", "s16le", "-ar", "48k", "-ac", "2", "-i", fmt.Sprintf("pipe:%v", r1))
	args = append(args, "-y", "-f", "wav", *output)

	if _, err = syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: append([]uintptr{0, 1, 2, uintptr(r1)},),
	}); err != nil {
		log.Fatalln(errors.Wrap(err, "failed to start ffmpeg"))
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	<- stopSignal

}
