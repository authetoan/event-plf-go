package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Human-readable error message
	Details string `json:"details"` // Detailed error for debugging (optional)
}

func NewAppError(code int, message string, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Details: %s", e.Code, e.Message, e.Details)
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors[0].Err
			var appErr *AppError
			if errors.As(err, &appErr) {
				// Handle custom application errors
				c.JSON(appErr.Code, gin.H{
					"error": gin.H{
						"message": appErr.Message,
						"details": appErr.Details,
					},
				})
			}
		}
	}
}

func HandleError(c *gin.Context, err error, code int, message string) {
	if err != nil {
		c.Error(NewAppError(code, message, err.Error())) // Attach error to Gin context
		c.Abort()                                        // Stop further processing
	}
}
