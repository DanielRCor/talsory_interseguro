FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-qr-api ./cmd/server

FROM alpine:3.22

WORKDIR /app

COPY --from=build /go-qr-api /usr/local/bin/go-qr-api

EXPOSE 8080

CMD ["go-qr-api"]
