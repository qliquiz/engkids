FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app .

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/app .
COPY scripts/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

ENV TZ=UTC
EXPOSE 3000

ENTRYPOINT ["/wait-for-it.sh", "engkids_db:5432", "--", "/wait-for-it.sh", "logstash:5000", "--", "/wait-for-it.sh", "elasticsearch:9200", "--", "./app"]
