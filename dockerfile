FROM golang:1.25-alpine3.22 AS builder
WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
COPY . .
RUN go build -o main main.go

FROM alpine:3.22 AS stage
WORKDIR /app
COPY --from=builder /app/main ./main
ENTRYPOINT [ "./main" ]