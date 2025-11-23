FROM golang:1.25.1-bookworm AS deps
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download

FROM golang:1.25.1-bookworm AS builder
WORKDIR /app
COPY --from=deps /go /go
COPY . .
RUN GOOS=linux go build -o main .

FROM debian:12.11-slim
COPY --from=builder /app/main /app/main
EXPOSE 37190
CMD ["/app/main"]
