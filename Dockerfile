FROM golang:1.19-alpine as builder
RUN apk add --no-cache git make
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY . /src
WORKDIR /src
RUN rm -f go.sum
RUN go get
RUN go test ./...
RUN make release

FROM alpine:3.16.2
WORKDIR /app
COPY --from=builder /src/vault-up /app/vault-up
CMD ["/app/vault-up"]
