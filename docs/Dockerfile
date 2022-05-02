# Build executable binary
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o docs

# Build a small image
FROM alpine
WORKDIR /app
COPY --from=builder app/docs /app
CMD ["./docs"]