FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest
COPY --from=builder /go/bin/app /app
EXPOSE $PORT

CMD [ "./app" ]
