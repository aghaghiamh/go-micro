FROM golang:1.24-alpine3.22 AS builder

RUN mkdir /app

COPY . /app
WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o brockerApp ./cmd/api

FROM alpine:3.22

RUN mkdir /app

COPY --from=builder /app/brockerApp /app

CMD [ "/app/brockerApp" ]