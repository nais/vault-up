FROM golang:1.13-alpine as builder
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

FROM alpine:3
MAINTAINER Johnny Horvi <johnny.horvi@gmail.com>
WORKDIR /app
COPY --from=builder /src/up /app/up
CMD ["/app/up"]
