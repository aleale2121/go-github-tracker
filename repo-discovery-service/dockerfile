# Build stage
FROM golang:1.21.12-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN  CGO_ENABLED=0 go build -o app ./cmd/api/ 

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8081
CMD [ "/app/app" ]