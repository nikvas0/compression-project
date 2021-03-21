FROM golang:latest AS builder

WORKDIR /chat_service/
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o /go/bin/chat_service

FROM golang:latest
COPY --from=builder /go/bin/chat_service /go/bin/chat_service

EXPOSE 8083

ENTRYPOINT ["/go/bin/chat_service"]