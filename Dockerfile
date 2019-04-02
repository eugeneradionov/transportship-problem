FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN apk add --no-cache git

RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /go/bin/app .
COPY --from=builder /go/src/app/public ./public

USER appuser

EXPOSE $PORT
CMD [ "./app" ]
