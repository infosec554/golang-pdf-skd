FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/pdf-sdk ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache poppler-utils ca-certificates

WORKDIR /app

COPY --from=builder /app/pdf-sdk .

EXPOSE 8080

CMD ["./pdf-sdk"]
