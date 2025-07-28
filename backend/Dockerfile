FROM golang:1.24.5-bookworm AS build
WORKDIR /build
COPY . .
RUN go mod tidy && go build -o app .

FROM debian:12.11-slim AS production
WORKDIR /production
COPY --from=build /build/app /production/app
EXPOSE 80
# CMD ["/production/app", "-path-to-cert='/env/certfile.crt'", "-path-to-key='/env/keyfile.key'"]
CMD ["/production/app"]
