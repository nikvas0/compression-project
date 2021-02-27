# Based on https://ergoz.ru/create-the-smallest-and-secured-golang-docker-image-based-on-scratch/
FROM golang:latest AS builder

WORKDIR /video_service/
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o /go/bin/video_service

FROM ubuntu:latest
# golang:latest
COPY --from=builder /go/bin/video_service /go/bin/video_service
RUN apt-get update && apt-get install -y ffmpeg

EXPOSE 80
ENTRYPOINT ["/go/bin/video_service"]
