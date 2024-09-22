FROM --platform=linux/amd64 alpine/curl as downloader

WORKDIR /app

RUN curl -X GET -k https://github.com/gitleaks/gitleaks/releases/download/v8.19.1/gitleaks_8.19.1_linux_x64.tar.gz -L --output gitleaks.tar.gz \
    && mkdir "gitleaks" \
    && tar -xvzf gitleaks.tar.gz -C gitleaks

FROM --platform=linux/amd64 golang:1.23.0-alpine as builder

WORKDIR /app

COPY . .

RUN go get ./...

RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g ./cmd/main/main.go

RUN go build -tags=jsoniter -o app cmd/main/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=downloader /app/gitleaks gitleaks

COPY --from=builder /app/app .

CMD ["./app"]