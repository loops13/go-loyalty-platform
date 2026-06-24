# GoLoyaltyPlatform monorepo

This repository is split into:

- `backend/` — Go API, Swagger docs, and backend build files
- `frontend/` — Angular 20 client

## Run locally

Backend:

```sh
cd backend
make run
```

Frontend:

```sh
cd frontend
npm install
npm start
```

The frontend uses Angular environments:

- `src/environments/environment.ts` for local development (`http://localhost:8080`)
- `src/environments/environment.prod.ts` for production (`/api`, suitable for a reverse proxy or ALB path rule)

## Docker

```sh
docker compose up --build
```

## Backend endpoints

- `GET /health`
- `GET /swagger/*`
- `GET /clients`
- `POST /clients`
- `GET /clients/{id}`
- `POST /clients/{id}/awards`
- `GET /clients/{id}/awards`
- `GET /rewards`
- `POST /clients/{id}/redeem`
