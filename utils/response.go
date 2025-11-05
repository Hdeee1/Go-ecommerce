package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success		bool		`json:"success"`
	Message		string		`json:"message,omitempty"`
	Data		interface{}	`json:"data,omitempty"`
	Error		string		`json:"error,omitempty"`
}

func SuccessResponse(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, Response{
		Success: true,
		Message: message,
		Data: data,
	})
}

func ErrorResponse(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Response{
		Success: false,
		Error: message,
	})
}