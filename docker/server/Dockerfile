FROM golang:latest as builder

WORKDIR /

COPY ./go.mod /
RUN go mod download

COPY . /
RUN make server

FROM ubuntu:20.04

RUN apt update
RUN apt install -y ffmpeg

COPY --from=builder bin/server /bin/server

CMD /bin/server

