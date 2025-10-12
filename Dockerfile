FROM golang:1.25.2-alpine3.22 AS builder

RUN apk --update add ca-certificates make git
RUN echo 'easee:*:65532:' > /tmp/group && \
    echo 'easee:*:65532:65532:easee:/:/easee-exporter' > /tmp/passwd

WORKDIR /workspace
COPY . /workspace

RUN make build

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
