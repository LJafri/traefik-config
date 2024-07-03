FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go build -o proxy_service .

EXPOSE 3000

CMD ["./proxy_service"]
