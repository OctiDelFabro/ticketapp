package services

import (
	"errors"
	"strings"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/utils"
	"gorm.io/gorm"
)

const clientRole = "CLIENT"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func Register(db *gorm.DB, req RegisterRequest) (*AuthResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if req.Name == "" || req.Email == "" || req.Password == "" || len(req.Password) < 6 {
		return nil, ErrInvalidRequest
	}

	if _, err := dao.FindUserByEmail(db, req.Email); err == nil {
		return nil, ErrEmailAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         clientRole,
	}

	if err := dao.CreateUser(db, &user); err != nil {
		return nil, err
	}

	return buildAuthResponse(user)
}

func Login(db *gorm.DB, req LoginRequest) (*AuthResponse, error) {
	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidRequest
	}

	user, err := dao.FindUserByEmail(db, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return buildAuthResponse(*user)
}

func buildAuthResponse(user domain.User) (*AuthResponse, error) {
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
