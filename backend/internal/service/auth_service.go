package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aimerneige/auto-you-koma/internal/auth"
	"github.com/aimerneige/auto-you-koma/internal/config"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

type AuthService struct {
	repo repository.UserRepository
	cfg  config.AuthConfig
}

func NewAuthService(repo repository.UserRepository, cfg config.AuthConfig) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*model.User, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hashed),
		QuotaLimit:   100,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if user.TOTPEnabled {
		return "", user, errors.New("2fa_required")
	}

	token, err := auth.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) Setup2FA(ctx context.Context, userID string) (string, string, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	key, err := auth.GenerateTOTPSecret(user.Email)
	if err != nil {
		return "", "", err
	}

	user.TOTPSecret = key.Secret()
	if err := s.repo.Update(ctx, user); err != nil {
		return "", "", err
	}

	return key.URL(), key.Secret(), nil
}

func (s *AuthService) VerifyAndEnable2FA(ctx context.Context, userID, passcode string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !auth.ValidateTOTP(passcode, user.TOTPSecret) {
		return errors.New("invalid passcode")
	}

	user.TOTPEnabled = true
	return s.repo.Update(ctx, user)
}

func (s *AuthService) LoginWith2FA(ctx context.Context, email, password, passcode string) (string, *model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.TOTPEnabled {
		return "", nil, errors.New("2fa not enabled")
	}

	if !auth.ValidateTOTP(passcode, user.TOTPSecret) {
		return "", nil, errors.New("invalid passcode")
	}

	token, err := auth.GenerateToken(user.ID, user.Email, s.cfg)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
