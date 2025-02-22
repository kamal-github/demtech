FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main cmd/main.go

FROM alpine:latest
RUN apk add --no-cache redis
COPY --from=builder /app/main /bin/main
RUN chmod +x /bin/main
EXPOSE 8080
CMD ["/bin/main"]
