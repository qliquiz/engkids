FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add bash ca-certificates tzdata

COPY --from=builder /app/app .
COPY scripts/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

ENV TZ=UTC
EXPOSE 3000

ENTRYPOINT ["/wait-for-it.sh", "engkids_db:5432", "-t", "90", "--", "./app"]