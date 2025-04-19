FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app .
CMD ["./app"]

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/app .
ENV TZ=UTC
EXPOSE 8080
CMD ["./app"]
