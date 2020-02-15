FROM golang:alpine as builder

WORKDIR /go/src/github.com/paulfantom/periodic-labeler
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM gcr.io/distroless/static
COPY --from=builder /go/src/github.com/paulfantom/periodic-labeler/periodic-labeler /
ENTRYPOINT ["/periodic-labeler"]
