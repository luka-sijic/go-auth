FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o go-auth-app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/go-auth-app .

EXPOSE 8082

CMD ["./go-auth-app"]