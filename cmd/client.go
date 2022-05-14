package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.Int("port", 7777, "listen port")
	input := flag.String("input", "", "input file (WAV)")
	output := flag.String("output", "", "output file (WAV)")
	denoiserLocation := flag.String("denoiser", "/bin/denoiser", "denoiser_location")
	flag.Parse()

	inputFile, err := os.Open(*input)
	if err != nil {
		log.Fatalln(err)
	}
	tmp, err := os.CreateTemp(os.TempDir(), "*.wav")
	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("%v:%v", *host, *port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := io.Copy(conn, inputFile); err != nil {
		log.Fatalln(err)
	}

	if _, err := io.Copy(tmp, conn); err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(*denoiserLocation, *denoiserLocation, *input, *output)
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

}
