package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(config utils.Configuration) (*gorm.DB, error) {
	// Format connection string
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s TimeZone=%s",
		config.DB.Username, config.DB.Password, config.DB.Name, config.DB.Host, config.DB.TimeZone)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Setup logger for GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// Set connection pool options
	conn.SetConnMaxIdleTime(time.Duration(config.DB.MaxIdleTime) * time.Minute)
	conn.SetConnMaxLifetime(time.Duration(config.DB.MaxLifeTime) * time.Hour)
	conn.SetMaxIdleConns(config.DB.MaxIdleConns)
	conn.SetMaxOpenConns(config.DB.MaxOpenConns)

	// Open a connection to the PostgreSQL databas
	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: conn}),
		&gorm.Config{Logger: newLogger, PrepareStmt: config.DB.Logging})

	if err != nil {
		return nil, err
	}

	return db, nil
}
