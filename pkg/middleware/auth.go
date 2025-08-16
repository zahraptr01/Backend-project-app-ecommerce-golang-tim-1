package middleware

import (
	"net/http"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/pkg/response"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Repo   repository.Repository
	Logger *zap.Logger
}

func NewAuthMiddleware(repo repository.Repository, logger *zap.Logger) AuthMiddleware {
	return AuthMiddleware{
		Repo:   repo,
		Logger: logger,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authz := ctx.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			response.ResponseBadRequest(ctx, http.StatusUnauthorized, "missing or invalid Authorization header")
			ctx.Abort()
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))

		claims, err := utils.VerifyJWT(token)
		if err != nil {
			response.ResponseBadRequest(ctx, http.StatusUnauthorized, "invalid token")
			ctx.Abort()
			return
		}

		stored, err := m.Repo.RedisRepo.GetToken(ctx.Request.Context(), claims.UserID, claims.Role)
		if err != nil || stored != token {
			response.ResponseBadRequest(ctx, http.StatusUnauthorized, "session expired")
			ctx.Abort()
			return
		}

		ctx.Set("userID", claims.UserID)
		ctx.Set("userRole", claims.Role)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("userRole")
		roleStr, _ := role.(string)
		switch roleStr {
		case "superadmin", "admin", "staff":
			ctx.Next()
			return
		default:
			response.ResponseBadRequest(ctx, http.StatusForbidden, "admin only")
			ctx.Abort()
			return
		}
	}
}

func CustomerOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("userRole")
		if roleStr, _ := role.(string); roleStr == "customer" {
			ctx.Next()
			return
		}
		response.ResponseBadRequest(ctx, http.StatusForbidden, "customer only")
	}
}
