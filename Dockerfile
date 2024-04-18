FROM golang:alpine3.19 AS  build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /htccli

FROM scratch

COPY --from=build /htccli /htccli
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["/htccli"]
