# HSE Coursework: Background noise removing system

## Requirements

Platform:

- Linux (Ubuntu, Debian, Fedora etc)
- Os X

Packets:

- Golang >1.16
- FFmpeg
- autoconf
- automake
- make

## Description

This project was built to use in video/audio conferencing systems to remove background noise. It contains 2 commands:

- *server* which accepts audio streams with TCP sockets, merges them and returns to clients
- *client*, this is example of client, which send audio file to server (for usage in audio/video conferences file should be replaced with audio stream from e.g WebRTC channel). After sending to server it receives combined stream with other clients, subtracts from it audio from current client and than uses RNNoise to reduce noise. RNNoise may also use GPU card to improve performance. For testing purposes it saves result to file

## Project structure

## Examples

In *testdata* folder you could check some examples of processed audio with this program

- *testdata/input* contains input files
- *testdata/output* contains result files with removed noise

## Build

To build everyting as binary files, run:

```shell
make
```

To build everything as docker containers, run:

```shell
make docker.build
```

## Run

### Server

```shell
./bin/server -host=$HOST -port=$PORT -clients=$NUMBER_OF_CLIENTS_TO_WAIT
```

### Client

```shell
./bin/client
```