# Build the application from source
FROM golang:1.20 as builder
RUN apt-get update && apt-get install -y ca-certificates openssl
ARG cert_location=/usr/local/share/ca-certificates

RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
RUN update-ca-certificates

WORKDIR /app
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download && go mod verify

COPY ./cmd/gophermart ./cmd/gophermart
COPY ./configs ./configs
COPY ./deployments/db/migrations ./deployments/db/migrations
COPY ./internal ./internal

RUN mkdir -p -m 0600 ./bin
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/gophermart ./cmd/gophermart/main.go


# Run the application from builder
FROM alpine:latest
RUN apk update && apk add bash
RUN apk add --no-cache bash
WORKDIR /app
COPY --from=builder /app/bin/gophermart .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/deployments/db/migrations ./deployments/db/migrations
EXPOSE 8080
CMD ./gophermart