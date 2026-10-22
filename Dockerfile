FROM golang:1.27.3 AS builder
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod main.go ./
RUN go build -trimpath -ldflags="-s -w" -o /vault-up .

FROM cgr.dev/chainguard/static:latest
WORKDIR /app
COPY --from=builder /vault-up /app/vault-up
ENTRYPOINT ["/app/vault-up"]
