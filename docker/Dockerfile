FROM golang:1.16-buster AS builder
ENV PROJECT_DIR /go/src/github.com/chez-shanpu/traffic-generator
COPY ./ ${PROJECT_DIR}/
WORKDIR ${PROJECT_DIR}/
RUN go build -o bin/tg


FROM ubuntu:focal
COPY --from=builder /go/src/github.com/chez-shanpu/traffic-generator/bin/tg /
RUN apt update && apt install -y \
    iperf3 \
    && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["/tg"]
CMD ["--help"]
