FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=America/La_Paz

COPY --from=builder /api /api
COPY migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/api"]
