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
	"regexp"
	"runtime"
	"strings"
	"time"

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

	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 15 << 20)
		c.Next()
	})

	// STATIC SERVING: Agar file di storage bisa diakses via browser/QR Code
	r.Static("/storage", "./storage")

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// BASE ROUTES
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
		v1.POST("/tenants", tenantHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.POST("/booths/pair", boothHandler.Pair)

		v1.POST("/save-history", func(c *gin.Context) {
			// Batasi ukuran body (Misal: max 10MB) agar server tidak hang
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20) 

			var input struct {
				Image     string `json:"image"`
				FrameName string `json:"frameName"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				slog.Error("Gagal bind JSON history", "error", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Payload terlalu besar atau format salah"})
				return
			}

			// Setup Folder
			storagePath := "./storage/history"
			if err := os.MkdirAll(storagePath, os.ModePerm); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder storage"})
				return
			}

			// Decode Base64
			idx := strings.Index(input.Image, ",")
			if idx == -1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format gambar tidak valid"})
				return
			}
			
			rawB64 := input.Image[idx+1:]
			data, err := base64.StdEncoding.DecodeString(rawB64)
			if err != nil {
				slog.Error("Gagal decode base64", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses gambar"})
				return
			}

			// Penamaan File
			timestamp := time.Now().Unix()
			reg := regexp.MustCompile("[^a-zA-Z0-9]+")
			cleanFrameName := reg.ReplaceAllString(input.FrameName, "_")
			
			fileName := fmt.Sprintf("%d_%s.jpg", timestamp, strings.ToLower(cleanFrameName))
			filePath := filepath.Join(storagePath, fileName)

			// Simpan ke Disk
			if err := os.WriteFile(filePath, data, 0644); err != nil {
				slog.Error("Gagal menulis file", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
				return
			}

			slog.Info("History saved successfully", "file", fileName)
			c.JSON(http.StatusOK, gin.H{
				"status": "success",
				"file":   fileName,
				"path":   "/storage/history/" + fileName, // Path relative untuk QR Code
			})
		})

		// AUTHORIZED ROUTES
		authorized := v1.Group("/")
		authorized.Use(middleware.AuthMiddleware())
		{
			authorized.POST("/booths", boothHandler.Register)
			authorized.GET("/booths", boothHandler.GetAllBooth)
			authorized.POST("/transactions/session", middleware.DeviceOnly(), trxHandler.StartSession)
		}
	}

	// LEGACY PRINT ROUTE
	v1.POST("/print", func(c *gin.Context) {
		var input struct {
			Image string `json:"image"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		b64data := input.Image[strings.IndexByte(input.Image, ',')+1:]
		data, _ := base64.StdEncoding.DecodeString(b64data)

		tmpFile := "temp_print.jpg"
		_ = os.WriteFile(tmpFile, data, 0644)

		if err := executePrint(tmpFile); err != nil {
			c.JSON(500, gin.H{"status": "error", "detail": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "success", "message": "Printing started"})
	})

	// SERVER STARTUP
	slog.Info("Server Photobooth berjalan", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

func executePrint(filePath string) error {
	if runtime.GOOS == "windows" {
		printerName := "Brother HL-L5100DN series"
		absImagePath, _ := filepath.Abs(filePath)
		if _, err := os.Stat(absImagePath); os.IsNotExist(err) {
			return fmt.Errorf("File gambar tidak ditemukan")
		}

		sumatraPath, _ := filepath.Abs("./SumatraPDF.exe")
		cmd := exec.Command(sumatraPath,
			"-print-to", printerName,
			"-print-settings", "fit,paper=A4",
			absImagePath,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("SumatraPDF Error: %s, Output: %s", err, string(output))
		}
	}
	return nil
}