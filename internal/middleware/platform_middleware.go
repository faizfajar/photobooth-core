package middleware

import (
	"log/slog"
	"net/http"
	"photobooth-core/internal/platform/response"

	// "photobooth-core/internal/platform/response"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS setup untuk mengizinkan akses dari frontend (Web/Electron)
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Sesuaikan dengan domain di production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

// GlobalRecovery menangkap panic agar server tidak mati total (crash)
func GlobalRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		slog.Error("RECOVERED_FROM_PANIC", "error", recovered)
		response.Abort(c, http.StatusInternalServerError, "Terjadi kesalahan internal pada server", recovered)
	})
}
