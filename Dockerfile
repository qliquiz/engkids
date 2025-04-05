FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app ./cmd/server
CMD ["./app"]
