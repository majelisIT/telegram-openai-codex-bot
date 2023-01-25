FROM golang:1.19-alpine as builder
WORKDIR /app


COPY . .

RUN go mod tidy

RUN go build -o app .


ENTRYPOINT ["/app/app"]

