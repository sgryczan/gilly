FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/pkg/app/
COPY . .

WORKDIR $GOPATH/src/pkg/app
RUN go get -d -v

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/bin/app

FROM alpine

RUN apk update && apk add --no-cache \
    ca-certificates \
    docker-cli
EXPOSE 8080

WORKDIR /go/bin
COPY --from=builder /go/bin/app /go/bin/app
COPY ssl ssl

ENTRYPOINT ["./app"]