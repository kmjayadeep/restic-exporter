# Stage 1: Build the Go application
FROM golang:1.20.7-alpine3.18 AS build

WORKDIR /app

COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o resticmon

FROM alpine:latest

# Install Restic and other dependencies
RUN apk add --no-cache restic

WORKDIR /app

COPY --from=build /app/resticmon .

EXPOSE 18090

CMD ["./resticmon"]
