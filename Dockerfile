#Build stage
FROM golang:1.19.1-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go


#RUN state
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY config.env .

EXPOSE 8080
CMD [ "/app/main" ]