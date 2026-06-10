package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/instaagrammeta/rentacar/backend-go/internal/apperr"
	"github.com/instaagrammeta/rentacar/backend-go/internal/auth"
	"github.com/instaagrammeta/rentacar/backend-go/internal/models"
)

// Authenticate validates credentials and returns tokens + the user payload.
func (s *Service) Authenticate(username, password string, ip string) (map[string]interface{}, error) {
	var user models.User
	err := s.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.Auth("Неверное имя пользователя или пароль")
	}
	if err != nil {
		return nil, err
	}
	if !auth.VerifyPassword(password, user.PasswordHash) {
		return nil, apperr.Auth("Неверное имя пользователя или пароль")
	}
	if !user.IsActive {
		return nil, apperr.Auth("Учётная запись отключена")
	}

	access, refresh, err := s.Tokens.GenerateTokens(&user)
	if err != nil {
		return nil, err
	}

	s.record(&Actor{ID: uintPtr(user.ID), Username: user.Username, IP: ip}, "login", "user", uintPtr(user.ID), "")

	return map[string]interface{}{
		"access_token":  access,
		"refresh_token": refresh,
		"user":          user.ToMap(),
	}, nil
}

// GetUserByID returns a user by primary key.
func (s *Service) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("Пользователь не найден")
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers returns all users ordered by username.
func (s *Service) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.DB.Order("username asc").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// CreateUserInput holds the fields required to create an employee account.
type CreateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// CreateUser creates a new employee account.
func (s *Service) CreateUser(in CreateUserInput, actor *Actor) (*models.User, error) {
	if in.Username == "" || in.Password == "" {
		return nil, apperr.Validation("Имя пользователя и пароль обязательны")
	}
	if len(in.Password) < 6 {
		return nil, apperr.Validation("Пароль должен содержать не менее 6 символов")
	}

	var count int64
	s.DB.Model(&models.User{}).Where("username = ?", in.Username).Count(&count)
	if count > 0 {
		return nil, apperr.Conflict("Пользователь с таким именем уже существует")
	}

	role := models.UserRole(in.Role)
	if in.Role == "" {
		role = models.RoleOperator
	}
	if !role.Valid() {
		return nil, apperr.Validation("Неизвестная роль: " + in.Role)
	}

	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     in.Username,
		FullName:     in.FullName,
		Role:         role,
		PasswordHash: hash,
		IsActive:     true,
	}
	if in.Email != "" {
		email := in.Email
		user.Email = &email
	}
	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	s.record(actor, "create_user", "user", uintPtr(user.ID), "Создан пользователь "+in.Username+" ("+in.Role+")")
	return user, nil
}

// ChangePassword updates a user's password.
func (s *Service) ChangePassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return apperr.Validation("Пароль должен содержать не менее 6 символов")
	}
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}
	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.DB.Save(user).Error
}

// Logout revokes the given token id until its natural expiry (best effort).
func (s *Service) Logout(jti string, expiresAt time.Time) {
	if s.Redis == nil || jti == "" {
		return
	}
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		ttl = time.Minute
	}
	_ = s.Redis.Set(context.Background(), "jwt:blacklist:"+jti, "1", ttl).Err()
}
