FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o app .

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add bash ca-certificates tzdata

COPY --from=builder /app/app .
COPY --chmod=+x scripts/wait-for-it.sh /wait-for-it.sh

ENV TZ=UTC
EXPOSE 3000

ENTRYPOINT ["/wait-for-it.sh", "engkids_db:5432", "-t", "90", "--", "./app"]