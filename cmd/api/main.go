// @title Photobooth Core API
// @version 1.0
// @description API Dokumentasi untuk Sistem SaaS Photobooth.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

	// MODULE: Booth (Scaffolded)
	bHandler "photobooth-core/internal/booth/handler"
	bRepo "photobooth-core/internal/booth/repository"
	bUcase "photobooth-core/internal/booth/usecase"

	// MODULE: Tenant & Users
	tHandler "photobooth-core/internal/tenant/handler"
	tRepo "photobooth-core/internal/tenant/repository"
	tUcase "photobooth-core/internal/tenant/usecase"
	uHandler "photobooth-core/internal/users/handler"
	uRepo "photobooth-core/internal/users/repository"
	uUcase "photobooth-core/internal/users/usecase"

	trHandler "photobooth-core/internal/transaction/handler"
	trRepo "photobooth-core/internal/transaction/repository"
	trUcase "photobooth-core/internal/transaction/usecase"
)

func main() {
	// INITIALIZATION: Setup JSON Logger & Env
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan variable sistem.")
	}
	cfg := config.LoadConfig()

	// INFRASTRUCTURE: Database & Migration
	db, err := postgres.NewConnection(cfg.DBDSN)
	if err != nil {
		slog.Error("Kritikal: Gagal terhubung ke database", "error", err)
		os.Exit(1)
	}

	// migration
	db.AutoMigrate(&domain.Tenant{}, &domain.User{}, &domain.Booth{}, &domain.Transaction{})
	postgres.SeedAdmin(db)

	// WIRING: Dependency Injection (User & Tenant)
	userRepository := uRepo.NewUserTenantRepository(db)
	userUsecase := uUcase.NewUserUsecase(userRepository)
	userHandler := uHandler.NewUserHandler(userUsecase)

	tenantRepository := tRepo.NewTenantRepository(db)
	tenantUsecase := tUcase.NewTenantUsecase(tenantRepository, userRepository, db)
	tenantHandler := tHandler.NewTenantHandler(tenantUsecase)

	// WIRING: Dependency Injection (Booth Module)
	boothRepository := bRepo.NewBoothRepository(db)
	boothUsecase := bUcase.NewBoothUsecase(boothRepository)
	boothHandler := bHandler.NewBoothHandler(boothUsecase)

	// transaction
	trxRepo := trRepo.NewTransactionRepository(db)
	trxUcase := trUcase.NewTransactionUsecase(trxRepo)
	trxHandler := trHandler.NewTransactionHandler(trxUcase)

	// ROUTER SETUP
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.GlobalRecovery())
	r.Use(middleware.CORS())

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// BASE ROUTES: Health & Welcome
	r.GET("/", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Sistem SaaS Photobooth Berhasil Dijalankan", gin.H{
			"app":    "Photobooth Core API",
			"status": "Active",
		})
	})
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "System is UP", nil)
	})

	// API VERSION 1
	v1 := r.Group("/api/v1")
	{
		// AUTHENTICATION: Public Routes
		v1.POST("/tenants", tenantHandler.Register)
		v1.POST("/login", userHandler.Login)

		// BOOTH HANDSHAKE: Public endpoint agar mesin bisa "Pairing" tanpa JWT user
		v1.POST("/booths/pair", boothHandler.Pair)

		// AUTHORIZED ROUTES: Perlu Bearer Token
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			// BOOTH MANAGEMENT (Untuk Dashboard Owner)
			authorized.POST("/booths", boothHandler.Register)   // Buat mesin baru
			authorized.GET("/booths", boothHandler.GetAllBooth) // List mesin milik tenant

			authorized.POST("/transactions/session", middleware.DeviceOnly(), trxHandler.StartSession)
		}
	}

	v1.POST("/print", func(c *gin.Context) {
		var input struct {
			Image string `json:"image"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 1. Decode Base64 ke File
		b64data := input.Image[strings.IndexByte(input.Image, ',')+1:]
		data, _ := base64.StdEncoding.DecodeString(b64data)

		tmpFile := "temp_print.jpg"
		if err := os.WriteFile(tmpFile, data, 0644); err != nil {
			c.JSON(500, gin.H{"error": "Gagal menulis file sementara"})
			return
		}

		// 2. PANGGIL FUNGSI executePrint DI SINI (Inilah kuncinya!)
		err := executePrint(tmpFile) // Memanggil fungsi yang berisi Debug Print & SumatraPDF

		if err != nil {
			// Jika error, kirimkan detailnya agar kita tahu masalahnya di mana
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Gagal mencetak",
				"detail":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"status": "success", "message": "Printing started"})
	})

	for _, route := range r.Routes() {
		fmt.Printf("Method: %s, Path: %s, Name: %s\n", route.Method, route.Path, route.Handler)
	}

	// SERVER STARTUP
	slog.Info("Server Photobooth berjalan", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func executePrint(filePath string) error {
	if runtime.GOOS == "windows" {
		printerName := "Brother HL-L5100DN series"

		// 1. Pastikan file gambar BENAR-BENAR ada sebelum diprint
		absImagePath, _ := filepath.Abs(filePath)
		if _, err := os.Stat(absImagePath); os.IsNotExist(err) {
			return fmt.Errorf("File gambar tidak ditemukan di: %s", absImagePath)
		}

		// 2. Gunakan alamat absolut untuk SumatraPDF
		sumatraPath, _ := filepath.Abs("./SumatraPDF.exe")

		// 3. LOGGING: Print perintah ke terminal Go agar bisa kamu copy-paste buat tes
		fmt.Printf("\n--- DEBUG PRINT ---\n")
		fmt.Printf("Command: %s -print-to \"%s\" -print-settings \"fit,paper=4r\" \"%s\"\n",
			sumatraPath, printerName, absImagePath)
		fmt.Printf("-------------------\n\n")

		cmd := exec.Command(sumatraPath,
			"-print-to", printerName,
			"-print-settings", "fit,paper=A4",
			absImagePath,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("SumatraPDF Error: %s, Output: %s", err, string(output))
		}

		return nil
	}
	return nil
}
