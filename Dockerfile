FROM golang:1.23-alpine
RUN apk add build-base

WORKDIR /app

COPY . .

RUN go env -w CGO_ENABLED=1 && go build -o main cmd/main.go

EXPOSE 5515

CMD ["./main"]