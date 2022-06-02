FROM golang:latest as builder

WORKDIR /

COPY ./go.mod /
RUN go mod download

COPY . /
RUN make server

FROM ubuntu:20.04

RUN apt update
RUN apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*
RUN add-apt-repository ppa:savoury1/ffmpeg5
RUN add-apt-repository ppa:savoury1/ffmpeg4
RUN apt update
RUN apt full-upgrade -y
RUN apt install -y ffmpeg

COPY --from=builder /bin/server /bin/server

CMD /bin/server -host=$HOST -port=$PORT -clients=$CLIENTS -ffmpeg=/usr/bin/ffmpeg

