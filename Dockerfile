FROM golang:1.18.3-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o server server.go

FROM alpine:3.16
WORKDIR /app
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /bin/migrate
COPY --from=builder /app/server .
COPY app.env .
COPY db/migration ./db/migration
COPY scripts ./scripts
EXPOSE 8080
ENTRYPOINT ["/app/scripts/start.sh"]
CMD ["/app/server"]