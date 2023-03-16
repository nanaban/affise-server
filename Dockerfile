FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o server ./cmd/server/*.go

# Deployment container
FROM scratch

EXPOSE 8080

COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /app/server /server
CMD [ "/server", "-addr", ":8080" ]