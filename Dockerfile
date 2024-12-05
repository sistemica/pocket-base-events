FROM golang:1.23.4-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -a -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/main .
ENV TZ=UTC
ENV REDIS_URL=localhost:6379
EXPOSE 8090
ENTRYPOINT ["./main"]
CMD ["serve", "--http=0.0.0.0:8090"]