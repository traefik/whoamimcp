# syntax=docker/dockerfile:1.2
FROM golang:1-alpine AS builder

RUN apk --no-cache --no-progress add git ca-certificates tzdata make \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

# syntax=docker/dockerfile:1.2
# Create a minimal container to run a Golang static binary
FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY whoamimcp /

# nobody
USER 65534

ENTRYPOINT ["/whoamimcp"]
EXPOSE 80
