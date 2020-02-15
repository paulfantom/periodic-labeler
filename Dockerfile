FROM golang:alpine as builder

WORKDIR /go/src/github.com/paulfantom/periodic-labeler
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/paulfantom/periodic-labeler/periodic-labeler /
ENTRYPOINT ["/periodic-labeler"]
