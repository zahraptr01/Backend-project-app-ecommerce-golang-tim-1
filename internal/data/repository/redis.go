package repository

import (
	"context"
	"fmt"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RedisRepository interface {
	SetToken(ctx context.Context, userID uint, role string, token string, duration time.Duration) error
	GetToken(ctx context.Context, userID uint, role string) (string, error)
	DeleteToken(ctx context.Context, userID uint, role string) error
}

type redisRepositoryImpl struct {
	DB  *gorm.DB
	Log *zap.Logger
}

func NewRedisRepository(db *gorm.DB, log *zap.Logger) RedisRepository {
	return &redisRepositoryImpl{
		DB:  db,
		Log: log,
	}
}

func (r *redisRepositoryImpl) SetToken(ctx context.Context, userID uint, role string, token string, duration time.Duration) error {
	key := buildKey(userID, role)
	return utils.RDB.Set(ctx, key, token, duration).Err()
}

func (r *redisRepositoryImpl) GetToken(ctx context.Context, userID uint, role string) (string, error) {
	key := buildKey(userID, role)
	return utils.RDB.Get(ctx, key).Result()
}

// DeleteToken menghapus token dari Redis
func (r *redisRepositoryImpl) DeleteToken(ctx context.Context, userID uint, role string) error {
	key := buildKey(userID, role)
	return utils.RDB.Del(ctx, key).Err()
}

func buildKey(userID uint, role string) string {
	return fmt.Sprintf("auth:token:%s:%d", role, userID)
}
