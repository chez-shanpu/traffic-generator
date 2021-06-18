FROM golang:1.16-buster AS builder
ENV PROJECT_DIR /go/src/github.com/chez-shanpu/traffic-generator
COPY ./ ${PROJECT_DIR}/
WORKDIR ${PROJECT_DIR}/
RUN go build -o bin/tg


FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/src/github.com/chez-shanpu/traffic-generator/bin/tg /
ENTRYPOINT ["/tg"]
CMD ["--help"]
