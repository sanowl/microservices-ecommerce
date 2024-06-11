FROM golang:1.18-alpine

WORKDIR /app

COPY . .

RUN go build -o product-service .

EXPOSE 8082

CMD ["./product-service"]
