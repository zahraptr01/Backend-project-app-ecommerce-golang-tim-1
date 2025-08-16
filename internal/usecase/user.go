package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
)

type UserUsecase interface {
	CreateAdmin(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error)
	ListUsers(ctx context.Context, sort, order string, page, limit int) ([]dto.UserResponse, error)
	GetUserByID(ctx context.Context, id uint) (dto.UserResponse, error)
	UpdateUser(ctx context.Context, id uint, req dto.UpdateUserRequest) (dto.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type userUsecase struct {
	repo        repository.UserRepository
	logger      *zap.Logger
	Config      utils.Configuration
	emailSender utils.EmailSender
}

// NewUserService sekarang menerima emailSender sebagai argumen
func NewUserService(
	repo repository.UserRepository,
	logger *zap.Logger,
	config utils.Configuration,
	emailSender utils.EmailSender,
) UserUsecase {
	return &userUsecase{
		repo:        repo,
		logger:      logger,
		Config:      config,
		emailSender: emailSender,
	}
}

func generateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (uc *userUsecase) CreateAdmin(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	if req.Role != "admin" && req.Role != "staff" {
		return dto.UserResponse{}, fmt.Errorf("role harus admin atau staff")
	}

	// generate password random
	rawPassword := generateRandomPassword()
	hashedPassword := utils.HashPassword(rawPassword)

	user := entity.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
		Model: entity.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	newUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// send email, cek uc.emailSender != nil
	if uc.emailSender != nil && req.Email != nil && *req.Email != "" {
		subject := "Akun Admin Baru"
		body := fmt.Sprintf("Halo %s,\n\nAkun admin Anda telah dibuat.\nEmail: %s\nPassword: %s\nSilakan login dan ubah password segera.", req.Fullname, *req.Email, rawPassword)
		if err := uc.emailSender.SendEmail(*req.Email, subject, body); err != nil {
			uc.logger.Error("failed to send email", zap.Error(err))
		}
	}

	return mapToUserResponse(newUser), nil
}

func (uc *userUsecase) ListUsers(ctx context.Context, sort, order string, page, limit int) ([]dto.UserResponse, error) {
	users, err := uc.repo.FindAll(ctx, sort, order, page, limit)
	if err != nil {
		return nil, err
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, dto.UserResponse{
			ID:        u.ID,
			Fullname:  u.Fullname,
			Email:     u.Email,
			Phone:     u.Phone,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
			UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		})
	}

	return res, nil
}

func (uc *userUsecase) GetUserByID(ctx context.Context, id uint) (dto.UserResponse, error) {
	u, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        u.ID,
		Fullname:  u.Fullname,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, id uint, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	user, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if req.Fullname != "" {
		user.Fullname = req.Fullname
	}
	if req.Email != nil {
		user.Email = req.Email
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	user.UpdatedAt = time.Now()

	updated, err := uc.repo.Update(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:        updated.ID,
		Fullname:  updated.Fullname,
		Email:     updated.Email,
		Phone:     updated.Phone,
		Role:      updated.Role,
		IsActive:  updated.IsActive,
		CreatedAt: updated.CreatedAt.Format(time.RFC3339),
		UpdatedAt: updated.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (uc *userUsecase) DeleteUser(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func mapToUserResponse(u entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Fullname:  u.Fullname,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
		UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
	}
}
