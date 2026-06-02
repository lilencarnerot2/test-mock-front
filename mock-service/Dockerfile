FROM golang:1.25-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o mock-service
EXPOSE 6524
CMD ["./mock-service"]