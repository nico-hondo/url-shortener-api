package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:8000", // Swagger UI
		"http://127.0.0.1:8000",
	}

	currentOrigin := ctx.GetHeader("Origin")

	if slices.Contains(allowedOrigins, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
		ctx.Header("Access-Control-Allow-Credentials", "true") // <-- WAJIB untuk credentials: "include"
		ctx.Header("Vary", "Origin")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	}

	// Preflight: jawab langsung
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
