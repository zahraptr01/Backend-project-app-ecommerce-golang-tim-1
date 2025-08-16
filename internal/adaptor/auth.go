package adaptor

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HandlerAuth struct {
	Auth   usecase.AuthService
	Logger *zap.Logger
}

func NewHandlerAuth(auth usecase.AuthService, logger *zap.Logger) HandlerAuth {
	return HandlerAuth{
		Auth:   auth,
		Logger: logger,
	}
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.Auth.Login(ctx.Request.Context(), req)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusUnauthorized, "login failed")
		return
	}

	response.ResponseSuccess(ctx, http.StatusOK, "login successful", res)
}

func (h *HandlerAuth) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Auth.ForgotPassword(ctx.Request.Context(), req)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, "failed to send otp")
		return
	}

	response.ResponseSuccess(ctx, http.StatusOK, "OTP sent successfully", nil)
}

func (h *HandlerAuth) ValidateOTP(ctx *gin.Context) {
	var req dto.ValidateOtpRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Auth.ValidateOtp(ctx.Request.Context(), req)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, "OTP validation failed")
		return
	}

	response.ResponseSuccess(ctx, http.StatusOK, "OTP validated successfully", nil)
}

func (h *HandlerAuth) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ResponseBadRequest(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.Auth.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, "failed to reset password")
		return
	}

	response.ResponseSuccess(ctx, http.StatusOK, "Password reset successfully", nil)
}

func (h *HandlerAuth) Logout(ctx *gin.Context) {
	userID, ok1 := ctx.Get("userID")
	role, ok2 := ctx.Get("userRole")
	if !ok1 || !ok2 {
		response.ResponseBadRequest(ctx, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.Auth.Logout(ctx.Request.Context(), userID.(uint), role.(string))
	if err != nil {
		response.ResponseBadRequest(ctx, http.StatusInternalServerError, "logout failed")
		return
	}

	response.ResponseSuccess(ctx, http.StatusOK, "logout successful", nil)
}
