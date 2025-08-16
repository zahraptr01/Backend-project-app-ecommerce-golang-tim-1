package wire

import (
	"project-app-ecommerce-golang-tim-1/internal/adaptor"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/usecase"
	"project-app-ecommerce-golang-tim-1/pkg/middleware"
	"project-app-ecommerce-golang-tim-1/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Wiring(repo repository.Repository, mLogger middleware.LoggerMiddleware, middlwareAuth middleware.AuthMiddleware, logger *zap.Logger, config utils.Configuration, emailSender utils.EmailSender) *gin.Engine {
	router := gin.New()
	router.Use(mLogger.LoggingMiddleware())
	api := router.Group("/api/v1")
	wireUser(api, middlwareAuth, repo, logger, config, emailSender)
	wireAuth(api, middlwareAuth, repo, logger, config)
	wireCustomer(api, middlwareAuth, repo, logger, config)
	wireStock(api, middlwareAuth, repo, logger, config)
	wireCategory(api, middlwareAuth, repo, logger, config)
	wireBanner(api, middlwareAuth, repo, logger, config)
	return router
}

func wireUser(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration, emailSender utils.EmailSender) {
	usecaseUser := usecase.NewUserService(repo.UserRepo, logger, config, emailSender)
	adaptorUser := adaptor.NewHandlerUser(usecaseUser, logger)

	// Testing endpoint
	router.GET("/test-handler", adaptorUser.TestHandler)

	// Admin management
	router.POST("/admins", middlwareAuth.Auth(), middleware.AdminOnly(), adaptorUser.CreateAdmin)
	router.GET("/admins", adaptorUser.ListUsers)
	router.GET("/admins/:id", adaptorUser.GetUserByID)
	router.PUT("/admins/:id", adaptorUser.UpdateUser)
	router.DELETE("/admins/:id", middlwareAuth.Auth(), middleware.AdminOnly(), adaptorUser.DeleteUser)
}

func wireAuth(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration) {
	usecaseAuth := usecase.NewAuthService(repo, logger, config)
	adaptorAuth := adaptor.NewHandlerAuth(usecaseAuth, logger)
	router.POST("/auth/login", adaptorAuth.Login)
	router.POST("/auth/forgot-password", adaptorAuth.ForgotPassword)
	router.POST("/auth/verify-otp", adaptorAuth.ValidateOTP)
	router.POST("/auth/reset-password", adaptorAuth.ResetPassword)
	router.POST("/auth/logout", middlwareAuth.Auth(), adaptorAuth.Logout)
}

func wireCustomer(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration) {
	usecaseCustomer := usecase.NewCustomerService(repo, logger, config)
	adaptorCustomer := adaptor.NewHandlerCustomer(usecaseCustomer, logger)
	router.POST("/register", adaptorCustomer.RegisterCustomer)
	// Address routes for customers
	usecaseAddress := usecase.NewAddressService(repo, logger)
	adaptorAddress := adaptor.NewHandlerAddress(usecaseAddress, logger)
	customerGroup := router.Group("/customer")
	customerGroup.Use(middlwareAuth.Auth())
	customerGroup.POST("/address", adaptorAddress.Create)
	customerGroup.GET("/address", adaptorAddress.List)
	customerGroup.PUT("/address/:id", adaptorAddress.Update)
	customerGroup.DELETE("/address/:id", adaptorAddress.Delete)
	customerGroup.PATCH("/address/:id/default", adaptorAddress.SetDefault)
	// Order routes
	usecaseOrder := usecase.NewOrderService(repo, logger)
	adaptorOrder := adaptor.NewHandlerOrder(usecaseOrder, logger)
	customerGroup.POST("/order", middlwareAuth.Auth(), adaptorOrder.CreateOrder)
	customerGroup.GET("/cart", middlwareAuth.Auth(), adaptorOrder.Cart)
	customerGroup.GET("/order/:id", middlwareAuth.Auth(), adaptorOrder.GetOrderDetail)
	customerGroup.GET("/order/history", middlwareAuth.Auth(), adaptorOrder.ListOrderHistory)
}

func wireStock(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration) {
	usecaseStock := usecase.NewStockService(repo, logger, config)
	adaptorStock := adaptor.NewHandlerStock(usecaseStock, logger)
	router.GET("/admin/stock", adaptorStock.List)
	router.GET("/admin/stock/:variant_id", adaptorStock.Detail)
	router.POST("/admin/stock/add", adaptorStock.Add)
	router.PUT("/admin/stock/set", adaptorStock.Set)
	router.DELETE("/admin/stock", adaptorStock.Delete)
	router.GET("/admin/stock/variants", adaptorStock.VariantsDropdown)
}

func wireCategory(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration) {
	usecaseCategory := usecase.NewCategoryService(repo, logger, config)
	adaptorCategory := adaptor.NewHandlerCategory(usecaseCategory, logger)
	router.GET("/admin/categories", adaptorCategory.List)
	router.GET("/admin/categories/:id", adaptorCategory.Get)
	router.POST("/admin/categories", adaptorCategory.Create)
	router.PUT("/admin/categories", adaptorCategory.Update)
	router.PATCH("/admin/categories/publish", adaptorCategory.TogglePublished)
	router.DELETE("/admin/categories/:id", adaptorCategory.Delete)
}

func wireBanner(router *gin.RouterGroup, middlwareAuth middleware.AuthMiddleware, repo repository.Repository, logger *zap.Logger, config utils.Configuration) {
	usecaseBanner := usecase.NewBannerService(repo, logger, config)
	adaptorBanner := adaptor.NewHandlerBanner(usecaseBanner, logger)
	router.GET("/admin/banners", adaptorBanner.List)
	router.GET("/admin/banners/:id", adaptorBanner.Get)
	router.POST("/admin/banners", adaptorBanner.Create)
	router.PUT("/admin/banners", adaptorBanner.Update)
	router.PATCH("/admin/banners/publish", adaptorBanner.TogglePublished)
	router.DELETE("/admin/banners/:id", adaptorBanner.Delete)
}
