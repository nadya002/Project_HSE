package main

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"log"
	"os/exec"
	"sync"
)

type AudioMixer struct {
	finish map[string]context.CancelFunc
	ctx    context.Context
	lock   sync.Mutex
	update chan struct{}
	cancel context.CancelFunc
}

var AlreadyExistsErr = errors.New("stream already exists")
var StreamNotFoundErr = errors.New("stream not found")
var EmptyStreamIdErr = errors.New("stream already must not be empty")
var MixerIsClosed = errors.New("mixer is closed")

func (mixer *AudioMixer) runFFmpegProcess(ctx context.Context) error {
	// надо убить старый процесс (если он был) и запустить новый
	// используй exec.CommandContext и в качестве контекста передавай ctx из аргументов
	// не забудь взять мютекс
	cmd := exec.CommandContext(ctx, "ffmpeg")
	err := cmd.Run()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			// очистить ресурсы
			return nil
		}
		return err
	}

	return nil
}

func (mixer *AudioMixer) createStreamData(id string) error {
	// создать что-то для нового стрима

	return nil
}

func (mixer *AudioMixer) removeStreamData(id string) error {

	return nil
}

func (mixer *AudioMixer) processStreamBytes(id string, bytes []byte) error {
	// в стрим прилетели новые байты

	return nil
}

func (mixer *AudioMixer) ffmpegUpdateLoop() {
	var ffmpegCancel context.CancelFunc
	go func() {
		for {
			ctx, cancel := context.WithCancel(mixer.ctx)
			ffmpegCancel = cancel
			if err := mixer.runFFmpegProcess(ctx); err != nil {
				log.Printf("ffmpeg process failed (%v)\n", err)
				mixer.cancel()
				break
			}
			if mixer.ctx.Err() != nil {
				break
			}
		}
	}()
loop:
	for {
		select {
		case <-mixer.ctx.Done():
			if ffmpegCancel != nil {
				ffmpegCancel()
			}
			break loop
		case <-mixer.update:
			if ffmpegCancel != nil {
				ffmpegCancel()
			}
		}
	}
}

func (mixer *AudioMixer) processStream(reader io.Reader, id string) {
	mixer.update <- struct{}{}
	for {
		buf := make([]byte, 4096)
		n, err := reader.Read(buf)
		if errors.Is(err, context.Canceled) {
			log.Printf("stream (%s) closed\n", id)
			break
		} else if errors.Is(err, io.EOF) {
			log.Printf("stream (%s) finished\n", id)
			break
		} else if err != nil {
			log.Printf("ERROR: stream (%s) read failed: (%v)\n", id, err)
			break
		}
		if err := mixer.processStreamBytes(id, buf[:n]); err != nil {
			log.Printf("ERROR: stream (%s) bytes process failed: (%v)\n", id, err)
		}
	}
	mixer.lock.Lock()
	delete(mixer.finish, id)
	if err := mixer.removeStreamData(id); err != nil {
		log.Printf("ERROR: stream (%s) data remove failed: (%v)\n", id, err)
	}
	mixer.lock.Unlock()
	mixer.update <- struct{}{}
}

func (mixer *AudioMixer) AddStream(reader io.Reader, id string) error {
	if mixer.ctx.Err() != nil {
		return MixerIsClosed
	}
	if id == "" {
		return EmptyStreamIdErr
	}
	{
		mixer.lock.Lock()
		defer mixer.lock.Unlock()
		if _, ok := mixer.finish[id]; ok {
			return errors.Wrapf(AlreadyExistsErr, "stream_id (%s)", id)
		} else {
			if err := mixer.createStreamData(id); err != nil {
				return err
			}
			ctx, cancel := context.WithCancel(mixer.ctx)
			mixer.finish[id] = cancel
			cancelableReader := NewCancelableReader(reader, ctx)
			go mixer.processStream(cancelableReader, id)
		}
	}
	return nil
}

func (mixer *AudioMixer) RemoveStream(id string) error {
	if id == "" {
		return EmptyStreamIdErr
	}
	mixer.lock.Lock()
	defer mixer.lock.Unlock()
	if cancel, ok := mixer.finish[id]; !ok {
		return errors.Wrapf(StreamNotFoundErr, "stream_id (%s)", id)
	} else {
		cancel()
		return nil
	}
}

func NewAudioMixer(ctx context.Context) *AudioMixer {
	mixerCtx, cancel := context.WithCancel(ctx)
	mixer := &AudioMixer{
		finish: make(map[string]context.CancelFunc),
		ctx:    mixerCtx,
		cancel: cancel,
	}
	go mixer.ffmpegUpdateLoop()
	return mixer
}
