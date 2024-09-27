FROM golang:bookworm AS builder
WORKDIR /src
COPY . /src
RUN go clean -modcache
RUN go build

FROM alpine:latest
COPY --from=builder /src/ccdb /usr/local/bin/ccdb
EXPOSE 6379
ENTRYPOINT ["/usr/local/bin/ccdb"]
