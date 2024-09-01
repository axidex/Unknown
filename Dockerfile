FROM --platform=linux/amd64 golang:1.23.0-alpine as builder

WORKDIR /app

COPY . .

RUN go get ./...

RUN go build -tags=jsoniter -o app cmd/main/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]