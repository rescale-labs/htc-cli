FROM golang:alpine3.19 AS  build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY *.go ./
COPY commands commands

RUN go mod download

RUN go build -ldflags="-extldflags=-static" -o /htccli

FROM scratch

COPY VERSION /opt/rescale/VERSION
COPY --from=build /htccli /opt/rescale/bin/htccli
