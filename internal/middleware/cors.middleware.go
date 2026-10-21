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
		"http://localhost:8000", // Ditambahkan untuk Swagger UI
		"http://127.0.0.1:8000",
	}

	currentOrigin := ctx.GetHeader("Origin")

	if slices.Contains(allowedOrigins, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
	} else if currentOrigin == "" {
		// Jika dipanggil langsung lewat browser address bar
		ctx.Header("Access-Control-Allow-Origin", "*")
	}

	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, XXX-Header")
	// Perbaikan: Pakai 'Methods' (jamak) & tambahkan POST, PUT, DELETE
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

	// Preflight request
	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
