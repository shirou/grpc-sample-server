FROM golang:1.15.2 as base

WORKDIR /go/src/app
COPY . .

RUN make build

### App
FROM scratch as app
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /go/src/app/grpc-sample-server /grpc-sample-server

CMD ["/grpc-sample-server"]