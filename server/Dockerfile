FROM golang:1.24 as builder

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server /server

CMD ["/server"]