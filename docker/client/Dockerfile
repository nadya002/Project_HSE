FROM golang:latest as builder

WORKDIR /

RUN apt update
RUN apt install -y make
RUN apt install -y autoconf
RUN apt install -y automake
RUN apt install -y libtool

COPY ./third-party/rnnoise /third-party/rnnoise

WORKDIR /third-party/rnnoise
RUN ./autogen.sh
RUN ./configure
RUN make

WORKDIR /

COPY ./go.mod /
RUN go mod download

COPY . /
RUN make client

FROM ubuntu:20.04

RUN apt update
RUN apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*
RUN add-apt-repository ppa:savoury1/ffmpeg5
RUN add-apt-repository ppa:savoury1/ffmpeg4
RUN apt update
RUN apt full-upgrade -y
RUN apt install -y ffmpeg

COPY --from=builder bin/client /bin/client
COPY --from=builder third-party/rnnoise/ third-party/rnnoise/

CMD sleep 1 && /bin/client -host=$HOST -port=$PORT -ffmpeg=/usr/bin/ffmpeg -denoiser=./third-party/rnnoise/examples/rnnoise_demo -input=$INPUT -output=$OUTPUT

