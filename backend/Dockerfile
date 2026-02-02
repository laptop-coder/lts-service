FROM golang:1.25.6-bookworm AS deps
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

FROM golang:1.25.6-bookworm AS builder
WORKDIR /app
COPY --from=deps /go /go
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-w -s -extldflags "-static"' -o ./bin/main ./cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-w -s -extldflags "-static"' -o ./bin/migrate ./cmd/migration/main.go

FROM debian:12.11-slim
COPY --from=builder /app/bin/main /usr/local/bin/app
COPY --from=builder /app/bin/migrate /usr/local/bin/migrate
COPY ./entrypoint.sh /app/entrypoint.sh
RUN adduser appuser && \
    chown appuser /app/entrypoint.sh /usr/local/bin/app /usr/local/bin/migrate && \
    chmod u+x /app/entrypoint.sh /usr/local/bin/app /usr/local/bin/migrate
USER appuser
WORKDIR /app
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:37190/health || exit 1
EXPOSE 37190
ENTRYPOINT ["./entrypoint.sh"]
