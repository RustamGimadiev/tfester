FROM golang:1.12.5-alpine AS builder
LABEL maintainer="Rustam Gimadiev <gimadiev.kzn@gmail.com>"
WORKDIR /go/src/tfester
RUN apk add --no-cache git ca-certificates
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/tfester/
RUN dep ensure -vendor-only

COPY . /go/src/tfester
RUN go get -t && CGO_ENABLED=0 GOOS=linux go test -c

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/tfester/tfester.test /bin/
ENTRYPOINT ["/bin/tfester.test"]
