FROM golang:latest as builder

WORKDIR /

COPY ./go.mod /
RUN go mod download

COPY . /
RUN make client

FROM ubuntu:20.04

RUN apt update
RUN apt install -y ffmpeg

COPY --from=builder bin/client /bin/client
COPY ./third-party/rnnoise/examples/rnnoise_demo bin/client /bin/denoiser

CMD /bin/client

