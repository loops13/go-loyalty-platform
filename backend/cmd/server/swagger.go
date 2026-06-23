package main

import "net/http"

// @title Awesome Project API
// @version 1.0
// @description A small rewards microservice for clients, awards, and redemptions.
// @BasePath /
// @schemes http

// Health handles GET /health.
// @Summary Health check
// @Tags system
// @Produce plain
// @Success 200 {string} string "ok"
// @Router /health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
