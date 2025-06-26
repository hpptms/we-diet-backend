FROM golang:1.24.2-alpine
WORKDIR /app
COPY . .

RUN go mod download

# 'my-gin-app'という名前でバイナリをビルド
RUN go build -o my-gin-app .

CMD ["/app/my-gin-app"]
