FROM golang:1.21 as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /build/tasks ./cmd/tasks/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

ENV DB_URL="postgresql://postgres:very_secure_password!....for_real@db:5432/tasks?sslmode=disable"
ENV TOKEN_SECRET="topSecretToke__..!!jsfdjq0324234234kk!!"
ENV ENV=local
ENV GRPC_PORT=9800


COPY --from=builder /build/tasks .
COPY --from=builder /build/migrations ./migrations

CMD ["./tasks"]
