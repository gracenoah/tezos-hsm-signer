FROM golang:1.11 as builder
RUN mkdir -p /build 
WORKDIR /build
COPY . . 
RUN go get -d -v ./...
RUN go mod verify
RUN go install -v ./...

FROM ubuntu:18.04
RUN mkdir /opt/signer
WORKDIR /opt/signer
COPY --from=builder /go/bin/tezos-hsm-signer .
COPY --from=builder /build/keys.yaml .
CMD ["./tezos-hsm-signer"] 
