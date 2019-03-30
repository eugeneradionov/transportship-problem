FROM golang:alpine AS builder

RUN adduser -D -g '' appuser

WORKDIR /go/src/app

COPY . .

RUN apk add --no-cache git

RUN go get -d -v ./...
RUN go install -v ./...

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/app /app

USER appuser

EXPOSE $PORT
CMD [ "./app" ]
