# Build stage
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/goloyaltyplatform ./cmd/server

# Run stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /usr/local/bin/goloyaltyplatform /usr/local/bin/goloyaltyplatform
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/goloyaltyplatform"]
