FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o kasi .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/kasi .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
RUN mkdir data
EXPOSE 8080
CMD ["./kasi"]
