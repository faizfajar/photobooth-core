// @title Photobooth Core API
// @version 1.0
// @description API Dokumentasi untuk Sistem SaaS Photobooth.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Import docs yang dihasilkan oleh swag init
	_ "photobooth-core/docs"

	"photobooth-core/internal/domain"
	"photobooth-core/internal/middleware"
	"photobooth-core/internal/platform/config"
	"photobooth-core/internal/platform/postgres"
	"photobooth-core/internal/platform/response"

	tHandler "photobooth-core/internal/tenant/handler"
	tRepo "photobooth-core/internal/tenant/repository"
	tUcase "photobooth-core/internal/tenant/usecase"
	uHandler "photobooth-core/internal/users/handler"
	uRepo "photobooth-core/internal/users/repository"
	uUcase "photobooth-core/internal/users/usecase"
)

func main() {
	// INITIALIZATION: Setup Logger (JSON format untuk produksi) dan Config
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan variable sistem.")
	}
	cfg := config.LoadConfig()

	// INFRASTRUCTURE: Database & Seeder
	db, err := postgres.NewConnection(cfg.DBDSN)
	if err != nil {
		slog.Error("Kritikal: Gagal terhubung ke database", "error", err)
		os.Exit(1)
	}

	db.AutoMigrate(&domain.Tenant{}, &domain.User{})
	postgres.SeedAdmin(db) // Membuat admin default jika belum ada

	// WIRING (DEPENDENCY INJECTION)
	userRepository := uRepo.NewUserTenantRepository(db)
	userUsecase := uUcase.NewUserUsecase(userRepository)
	userHandler := uHandler.NewUserHandler(userUsecase)

	tenantRepository := tRepo.NewTenantRepository(db)
	tenantUsecase := tUcase.NewTenantUsecase(tenantRepository, userRepository, db)
	tenantHandler := tHandler.NewTenantHandler(tenantUsecase)

	// ROUTER SETUP
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Menggunakan gin.New() agar kita bisa mengontrol middleware secara penuh
	r := gin.New()
	r.Use(gin.Logger())                // Logging request
	r.Use(middleware.GlobalRecovery()) // Menangkap panic
	r.Use(middleware.CORS())           // Izin akses domain (CORS)

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health Check & Welcome
	r.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Sistem SaaS Photobooth Berhasil Dijalankan", gin.H{
			"app":    "Photobooth Core API",
			"status": "Active",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "System is UP", nil)
	})

	// Group API Versi 1
	v1 := r.Group("/api/v1")
	{
		// Registrasi Tenant & Admin Account pertama kali
		v1.POST("/tenants", tenantHandler.Register)

		// Autentikasi User (Login) untuk mendapatkan JWT
		v1.POST("/login", userHandler.Login)
	}

	// SERVER STARTUP
	slog.Info("Server Photobooth berjalan", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
