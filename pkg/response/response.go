package response

import (
	"project-app-ecommerce-golang-tim-1/internal/dto"

	"github.com/gin-gonic/gin"
)

type Reponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Reponse2 struct {
	Status  bool `json:"status"`
	Message any  `json:"message"`
	Data    any  `json:"data,omitempty"`
}

func ResponseSuccess(ctx *gin.Context, code int, message string, data any) {
	response := Reponse{
		Status:  true,
		Message: message,
		Data:    data,
	}
	ctx.JSON(code, response)
}

func ResponseBadRequest(ctx *gin.Context, code int, message string) {
	response := Reponse{
		Status:  false,
		Message: message,
	}
	ctx.JSON(code, response)
}

func ResponseBadRequest2(ctx *gin.Context, code int, message any) {
	response := Reponse2{
		Status:  false,
		Message: message,
	}
	ctx.JSON(code, response)
}

func ResponsePagination(ctx *gin.Context, code int, message string, data any, pagination dto.Pagination) {
	response := map[string]interface{}{
		"status":     true,
		"message":    message,
		"data":       data,
		"pagination": pagination,
	}

	ctx.JSON(code, response)
}
