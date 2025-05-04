from golang:1.23-alpine as go-url-shortener

workdir /app

copy go.mod go.sum ./
run go mod download

copy . .

run go build -o main ./cmd/url-shortener/

expose 8000

cmd ["./main"]
