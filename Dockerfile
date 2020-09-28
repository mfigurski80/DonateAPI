FROM golang:latest AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest AS production

COPY --from=builder /app/main /
CMD ["/main"]

