# GoLoyaltyPlatform

A small Go microservice demonstrating clients, awards, and rewards.

Run locally:

- go run ./cmd/server
- or using Makefile: make run

Endpoints (JSON):
- POST /clients {name,email}
- GET /clients/{id}
- POST /clients/{id}/awards {type}
- GET /clients/{id}/awards
- GET /rewards
- POST /clients/{id}/redeem {rewardId}

Configuration:
- PORT environment variable (default 8080)

Docker:
- docker build -t GoLoyaltyPlatform:latest .
- docker run -p 8080:8080 GoLoyaltyPlatform:latest
