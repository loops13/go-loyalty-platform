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
API_PROXY_TARGET=http://localhost:8080 npm start
```

PowerShell:

```powershell
cd frontend
npm install
$env:API_PROXY_TARGET="http://localhost:8080"; npm start
```

The app itself uses relative API paths (`/api`) in both dev and prod, so there are no hardcoded localhost URLs in application code.

## Docker

```sh
docker compose up --build
```

On EC2, expose only port **80** on the instance and run the same Compose stack.
Nginx in the frontend container serves Angular and proxies `/api/*`, `/swagger/*`, and `/health` to the Go backend over the internal Docker network.

## Backend endpoints

- `GET /health`
- `GET /swagger/*`
- `GET /clients`
- `POST /clients`
- `GET /clients/{id}`
- `DELETE /clients/{id}`
- `POST /clients/{id}/awards`
- `GET /clients/{id}/awards`
- `GET /rewards`
- `POST /clients/{id}/redeem`
