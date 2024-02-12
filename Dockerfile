FROM golang:1.21 as builder

WORKDIR /build

COPY . .
RUN apt-get update && apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
RUN mkdir /proto/gen
RUN protoc -I proto/proto  --go-grpc_out=proto/gen proto/proto/api.proto  --go_out=proto/gen proto/proto/api.proto --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative

# if .env do not exist just create a blank file, by default it will use os vars
RUN touch .env

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /build/tasks ./cmd/tasks/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/tasks .

COPY --from=builder /build/.env .


COPY --from=builder /build/migrations ./migrations

CMD ["./tasks"]
