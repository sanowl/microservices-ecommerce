FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o user-service .

EXPOSE 8081

CMD ["./user-service"]
