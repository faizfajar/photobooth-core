package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response adalah struktur standar untuk API photobooth
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse adalah struktur standar untuk kegagalan API
type ErrorResponse struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Errors interface{} `json:"errors,omitempty"` // Untuk detail validasi field
}

// JSON memberikan respons sukses yang konsisten
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// Error memberikan respons gagal yang konsisten
func Error(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, ErrorResponse{
		Status: false,
		Error:  message,
		Errors: details,
	})
}

// ValidationErrorFormat mendefinisikan struktur error per field
type ValidationErrorFormat struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validation mengembalikan daftar error input yang user-friendly
func Validation(c *gin.Context, err error) {
	var errors []ValidationErrorFormat

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			msg := ""
			switch fe.Tag() {
			case "required":
				msg = "Field ini wajib diisi"
			case "email":
				msg = "Format email tidak valid"
			case "min":
				msg = "Karakter terlalu pendek"
			default:
				msg = "Input tidak valid"
			}
			errors = append(errors, ValidationErrorFormat{
				Field:   fe.Field(),
				Message: msg,
			})
		}
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Status: false,
		Error:  "Validasi gagal",
		Errors: errors,
	})

}

// Abort memberikan respons gagal dan langsung menghentikan rantai eksekusi (untuk Middleware)
func Abort(c *gin.Context, code int, message string, details interface{}) {
	Error(c, code, message, details)
	c.Abort()
}
