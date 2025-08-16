package usecase

import (
	"context"
	"errors"
	"project-app-ecommerce-golang-tim-1/internal/data/entity"
	"project-app-ecommerce-golang-tim-1/internal/data/repository"
	"project-app-ecommerce-golang-tim-1/internal/dto"
	"project-app-ecommerce-golang-tim-1/pkg/utils"
	"time"

	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.ResponseUser, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ValidateOtp(ctx context.Context, req dto.ValidateOtpRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	Logout(ctx context.Context, userID uint, role string) error
}

type authService struct {
	Repo   repository.Repository
	Logger *zap.Logger
	Config utils.Configuration
}

func NewAuthService(repo repository.Repository, logger *zap.Logger, config utils.Configuration) AuthService {
	return &authService{
		Repo:   repo,
		Logger: logger,
		Config: config,
	}
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.ResponseUser, error) {
	user, err := s.Repo.AuthRepo.FindCustomerByEmailOrPhone(ctx, req.EmailOrPhone)
	if err != nil {
		s.Logger.Error("login failed: user not found", zap.String("error", err.Error()))
		return dto.ResponseUser{}, errors.New("invalid username or password")
	}

	isValid := utils.CheckPassword(req.Password, user.Password)
	if !isValid {
		s.Logger.Error("login failed: wrong password", zap.String("email/phone", req.EmailOrPhone))
		return dto.ResponseUser{}, errors.New("invalid username or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return dto.ResponseUser{}, err
	}

	if err := s.Repo.RedisRepo.SetToken(ctx, user.ID, user.Role, token, 24*time.Hour); err != nil {
		return dto.ResponseUser{}, err
	}

	return dto.ResponseUser{
		Name:  user.Fullname,
		Email: utils.Deref(user.Email),
		Token: token,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error {
	user, err := s.Repo.AuthRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("email not registered")
	}

	otpCode := utils.GenerateOTP(4)
	otp := &entity.AuthOTP{
		UserID:    user.ID,
		OtpCode:   otpCode,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := s.Repo.AuthRepo.SaveOTP(ctx, otp); err != nil {
		return err
	}

	return utils.SendOTPEmail(s.Config, req.Email, otpCode)
}

func (s *authService) ValidateOtp(ctx context.Context, req dto.ValidateOtpRequest) error {
	otp, err := s.Repo.AuthRepo.FindOTP(ctx, req.Email)
	if err != nil {
		return errors.New("otp not found")
	}

	if time.Now().After(otp.ExpiresAt) {
		return errors.New("otp expired")
	}

	if otp.OtpCode != req.OTP {
		return errors.New("invalid otp")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	hashed := utils.HashPassword((req.NewPassword))
	return s.Repo.AuthRepo.UpdatePasswordByEmail(ctx, req.Email, hashed)
}

func (s *authService) Logout(ctx context.Context, userID uint, role string) error {
	return s.Repo.RedisRepo.DeleteToken(ctx, userID, role)
}
