FROM golang:1.24-alpine

WORKDIR /app

Copy go.mod go.sum ./

RUN apk add --no-cache git
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
