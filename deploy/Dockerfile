FROM golang:1.20 AS builder

WORKDIR /cmd

COPY . .

RUN go mod download
RUN go build -o rag ./cmd/rag

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /cmd/rag ./

CMD ["./rag"]