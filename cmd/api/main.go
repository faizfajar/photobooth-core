// Package main adalah entry point utama untuk menjalankan layanan Photobooth Core API.
// Proyek ini menggunakan arsitektur Clean Architecture untuk skalabilitas SaaS.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// Import internal platform
	"photobooth-core/internal/platform/config"
	"photobooth-core/internal/platform/postgres"

	// Import module Tenant
	tHandler "photobooth-core/internal/tenant/handler"
	tRepo "photobooth-core/internal/tenant/repository"
	tUcase "photobooth-core/internal/tenant/usecase"

	// Import module Users
	uHandler "photobooth-core/internal/users/handler"
	uRepo "photobooth-core/internal/users/repository"
	uUcase "photobooth-core/internal/users/usecase"

	// Import domain untuk migrasi
	"photobooth-core/internal/domain"
)

func main() {
	// 1. INITIALIZATION: Memuat konfigurasi dari .env
	// Gunakan os.Getenv jika di-deploy menggunakan Docker/Kubernetes.
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variable sistem.")
	}
	cfg := config.LoadConfig()

	// 2. INFRASTRUCTURE: Inisialisasi koneksi Database PostgreSQL
	db, err := postgres.NewConnection(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Kritikal: Gagal terhubung ke database: %v", err)
	}

	// Menjalankan Auto Migration untuk memastikan skema tabel selalu sinkron
	db.AutoMigrate(&domain.Tenant{}, &domain.User{})

	// ---------------------------------------------------------
	// 3. WIRING (DEPENDENCY INJECTION): Menyambungkan modul-modul
	// ---------------------------------------------------------

	// --- Modul Users ---
	userRepository := uRepo.NewUserTenantRepository(db)
	userUsecase := uUcase.NewUserUsecase(userRepository)
	userHandler := uHandler.NewUserHandler(userUsecase)

	// --- Modul Tenant ---
	// TenantUsecase membutuhkan userRepository dan instance db untuk transaksi registrasi
	tenantRepository := tRepo.NewTenantRepository(db)
	tenantUsecase := tUcase.NewTenantUsecase(tenantRepository, userRepository, db)
	tenantHandler := tHandler.NewTenantHandler(tenantUsecase)

	// ---------------------------------------------------------
	// 4. ROUTER SETUP: Konfigurasi Endpoint API menggunakan Gin
	// ---------------------------------------------------------

	// Gunakan gin.ReleaseMode jika di production untuk performa maksimal.
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Route Welcome (Pusat Kendali Info)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":     "Photobooth Core API",
			"status":  "Active",
			"message": "Sistem SaaS Photobooth Berhasil Dijalankan",
		})
	})

	// Route Health Check untuk Monitoring (K8s/Docker Probe)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	// Group API Versi 1
	v1 := r.Group("/api/v1")
	{
		// Registrasi Tenant & Admin Account pertama kali
		v1.POST("/tenants", tenantHandler.Register)

		// Autentikasi User (Login) untuk mendapatkan JWT
		v1.POST("/login", userHandler.Login)
	}

	// 5. SERVER STARTUP
	log.Printf("Server Photobooth berjalan di port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
