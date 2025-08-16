package main

import (
	"log"
	"project-app-ecommerce-golang-tim-1/cmd"
	"project-app-ecommerce-golang-tim-1/internal/data"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/wire"
	"project-app-ecommerce-golang-tim-1/pkg/database"
	"project-app-ecommerce-golang-tim-1/pkg/middleware"
	"project-app-ecommerce-golang-tim-1/pkg/utils"

	"go.uber.org/zap"
)

func main() {

	// read config
	config, err := utils.ReadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	// init logger
	logger, err := utils.InitLogger(config.PathLogger, config)
	if err != nil {
		log.Fatal("can't init logger %w", zap.Error(err))
	}

	//Init db
	db, err := database.InitDB(config)
	if err != nil {
		logger.Fatal("can't connect to database ", zap.Error(err))
	}

	// migration
	if err := data.AutoMigrate(db); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	// seeder
	if err := data.SeedAll(db); err != nil {
		logger.Fatal("failed to seed initial data", zap.Error(err))
	}

	// Inisialisasi Redis
	utils.InitRedis(config)
	defer func() {
		if err := utils.CloseRedis(); err != nil {
			logger.Fatal("error closing redis", zap.Error(err))
		}
	}()

	repo := repository.NewRepository(db, logger)
	mLogger := middleware.NewLoggerMiddleware(logger)
	mAuth := middleware.NewAuthMiddleware(repo, logger)
	emailSender := utils.NewEmailSender(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPEmail,
		config.SMTPPassword,
	)
	router := wire.Wiring(repo, mLogger, mAuth, logger, config, emailSender)

	cmd.ApiServer(config, logger, router)
}
