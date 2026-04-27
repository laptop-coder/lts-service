FROM golang:1.25.6-bookworm AS deps
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

FROM golang:1.25.6-bookworm AS builder
WORKDIR /app
COPY --from=deps /go /go
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags='-w -s -extldflags "-static"' -o ./bin/main ./cmd/app

FROM debian:12.11-slim
COPY --from=builder /app/bin/main /usr/local/bin/app
RUN adduser appuser && \
    chown appuser /usr/local/bin/app && \
    chmod u+x /usr/local/bin/app
USER appuser
WORKDIR /app
EXPOSE 37190

CMD /usr/local/bin/app
