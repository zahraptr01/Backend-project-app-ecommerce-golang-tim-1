package adaptor

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerUser struct {
	UserUC usecase.UserUsecase
	Logger *zap.Logger
}

func NewHandlerUser(userUC usecase.UserUsecase, logger *zap.Logger) HandlerUser {
	return HandlerUser{
		UserUC: userUC,
		Logger: logger,
	}
}

func (h *HandlerUser) TestHandler(ctx *gin.Context) {
	response.ResponseSuccess(ctx, http.StatusOK, "Ini adalah test handler", nil)
}

// CREATE ADMIN/STAFF
func (h *HandlerUser) CreateAdmin(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.UserUC.CreateAdmin(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "admin/staff berhasil dibuat",
		"data":    res,
	})
}

// LIST USERS with sorting & pagination
func (h *HandlerUser) ListUsers(c *gin.Context) {
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	res, err := h.UserUC.ListUsers(c.Request.Context(), sort, order, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GET USER BY ID
func (h *HandlerUser) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	res, err := h.UserUC.GetUserByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UPDATE USER
func (h *HandlerUser) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.UserUC.UpdateUser(c.Request.Context(), uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE USER (role = admin only for delete)
func (h *HandlerUser) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := h.UserUC.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response JSON success
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User berhasil dihapus",
		"user_id": id,
	})
}
