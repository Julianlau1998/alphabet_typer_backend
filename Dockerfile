FROM golang:latest as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app .

EXPOSE 8080

CMD ["/app/main"]