ARG BASE_IMAGE=alpine@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a
FROM $BASE_IMAGE as build

#RUN apk update && apk add \
#  wget \
#  g++ \
#  make \
#  libffi-dev \
#  openssl-dev \
#  tar \
#  xz \
#  zlib-dev
#
#RUN mkdir -p /opt/rescale /usr/local/src
#
#RUN cd /usr/local/src &&  \
#    wget http://www.musl-libc.org/releases/musl-1.1.10.tar.gz && \
#    tar -xzf musl-1.1.10.tar.gz
#
#RUN cd /usr/local/src/musl-1.1.10 && \
#    ./configure --prefix=/opt/rescale --with-libs='-lpthread'
#
#RUN cd /usr/local/src/musl-1.1.10 && \
#    make -j $(grep -cE 'processor\s+: [0-9]+' /proc/cpuinfo)
#
#RUN cd /usr/local/src/musl-1.1.10 && make install
#
#COPY main.go go.mod commands/  /usr/local/src/htc/
#
#RUN cd /usr/local/src/htc/ && \
#    CC=/opt/rescale/musl/bin/musl-gcc go build --ldflags '-linkmode external -extldflags "-static"' main.go -o htccli \
#
#RUN file /usr/local/src/htc/htccli
#
#CMD ["/usr/local/src/htc/htccli"]

FROM golang:alpine3.19 AS  build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -ldflags="-extldflags=-static" -o /htccli

FROM scratch

COPY VERSION /opt/rescale/VERSION
COPY --from=build /htccli /opt/rescale/bin/htccli
