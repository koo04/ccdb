FROM golang:bookworm AS builder
WORKDIR /src
COPY . /src
RUN go clean -modcache
RUN go build

FROM ubuntu:latest
COPY --from=builder /src/ccdb /usr/local/bin/ccdb
COPY --from=builder /src/config.example.toml /etc/ccdb/config.toml
EXPOSE 6969
ENTRYPOINT ["/usr/local/bin/ccdb"]
