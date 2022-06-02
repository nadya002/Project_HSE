package server

import (
	"fmt"
	"github.com/pkg/errors"
	"hse/pkg/pipe"
	"hse/pkg/runner"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type RunnerConfig struct {
	Host            string
	Port            int
	NumberOfClients int
	FFmpeg          string
}

type serverRunner struct {
	RunnerConfig
}

func (config *serverRunner) Run() error {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt)
	server, err := New()
	if err != nil {
		return err
	}
	defer func() {
		if err := server.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	log.Printf("Start listening on (%v:%v)...\n", config.Host, config.Port)
	if err := server.Listen(config.Host, config.Port); err != nil {
		log.Fatalln(err)
	}
	log.Println("Waiting for connections...")
	connections := make([]Conn, 0, config.NumberOfClients)
	for config.NumberOfClients > 0 {
		log.Printf("Rest to start %v\n", config.NumberOfClients)
		config.NumberOfClients--
		connCh, errCh := server.Accept()
		select {
		case err := <-errCh:
			return err
		case conn := <-connCh:
			connections = append(connections, conn)
		case <-stopSignal:
			log.Println("stop signal received")
			return nil
		}
	}
	log.Println("Connections were successfully accepted")

	args := []string{config.FFmpeg}
	fds := make([]uintptr, 0, len(connections))
	for _, conn := range connections {
		fds = append(fds, conn.Fd())
		args = append(args, "-f", "wav", "-i", fmt.Sprintf("pipe:%v", conn.Fd()))
	}
	args = append(args, "-filter_complex", fmt.Sprintf("amerge=inputs=%v", len(connections)))
	args = append(args, "-ac", "2")

	r, w, err := pipe.Pipe()
	if err != nil {
		return errors.Wrap(err, "failed to open pipe")
	}

	fds = append(fds, uintptr(w))
	args = append(args, "-f", "wav", fmt.Sprintf("pipe:%v", w))

	if _, err = syscall.ForkExec(args[0], args, &syscall.ProcAttr{
		Files: append([]uintptr{0, 1, 2}, fds...),
	}); err != nil {
		log.Fatalln(errors.Wrap(err, "failed to start ffmpeg"))
	}

	buf := make([]byte, 4096)
	for {
		n, err := syscall.Read(r, buf)
		if n > 0 {
			for _, conn := range connections {
				if err := conn.Send(buf[:n]); err != nil {
					return errors.Wrapf(err, "failed to send data")
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return errors.Wrapf(err, "failed to read ffmpeg output")
		}
	}
	return nil
}

func NewRunner(config RunnerConfig) runner.Runner {
	return &serverRunner{
		RunnerConfig: config,
	}
}
