FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /groupie-tracker

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /groupie-tracker .
COPY templates ./templates

LABEL maintainer="Juroll"
LABEL version="1.0"
LABEL description="Artists Informations Web App"

EXPOSE 8000

CMD ["./groupie-tracker"]
