package response

import (
	"github.com/gin-gonic/gin"
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
func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// ERROR memberikan respons gagal yang konsisten
func ERROR(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, ErrorResponse{
		Status: false,
		Error:  message,
		Errors: details,
	})
}
