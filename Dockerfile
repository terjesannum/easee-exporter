FROM golang:1.20.3-alpine3.17 as builder

RUN apk --update add ca-certificates
RUN echo 'easee:*:65532:' > /tmp/group && \
    echo 'easee:*:65532:65532:easee:/:/easee-exporter' > /tmp/passwd

WORKDIR /workspace
COPY go.* ./
RUN go mod download

COPY . /workspace

RUN CGO_ENABLED=0 go build -a -o easee-exporter .

FROM scratch

LABEL org.opencontainers.image.title="easee-exporter" \
      org.opencontainers.image.description="Prometheus exporter for Easee EV chargers" \
      org.opencontainers.image.authors="Terje Sannum <terje@offpiste.org>" \
      org.opencontainers.image.url="https://github.com/terjesannum/easee-exporter" \
      org.opencontainers.image.source="https://github.com/terjesannum/easee-exporter"

WORKDIR /
EXPOSE 8080

COPY --from=builder /tmp/passwd /tmp/group /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /workspace/easee-exporter .

USER 65532:65532

ENTRYPOINT ["/easee-exporter"]
