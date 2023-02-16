# Build App
FROM golang:1.20.0-alpine3.17 AS builder

WORKDIR ${GOPATH}/src/github.com/mpostument/awstaghelper
COPY . ${GOPATH}/src/github.com/mpostument/awstaghelper

RUN go build -o /go/bin/awstaghelper .


# Create small image with binary
FROM alpine:3.17

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/awstaghelper /usr/bin/awstaghelper

ENTRYPOINT ["/usr/bin/awstaghelper"]