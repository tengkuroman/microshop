# Build executable binary
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o shopping-service

# Build a small image
FROM alpine
WORKDIR /app
COPY --from=builder app/shopping-service /app
CMD ["./shopping-service"]