FROM golang:1.24-alpine AS builder

WORKDIR /order-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk add --no-cache git

RUN git clone https://github.com/pressly/goose /tmp/goose \
    && cd /tmp/goose \
    && go mod tidy \
    && go build -tags='no_sqlite3 no_clickhouse no_mssql no_mysql no_ydb no_libsql no_vertica' -o /order-service/goose ./cmd/goose


RUN go build -o order-service ./cmd/order-service/main.go


FROM alpine:3.18

WORKDIR /order-service

COPY --from=builder /order-service/order-service .
COPY --from=builder /order-service/goose .
COPY db/migrations ./migrations

COPY .env .

EXPOSE 8080

CMD ["./order-service"]
