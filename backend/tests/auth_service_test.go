package tests

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/services"
)

func TestRegisterSuccessCreatesClientAndReturnsSafeResponse(t *testing.T) {
	db := newTestDB(t)
	t.Setenv("JWT_SECRET", "test-secret")

	response, err := services.Register(db, services.RegisterRequest{
		Name:     "Client User",
		Email:    "client@example.com",
		Password: "secret123",
	})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustNotEmpty(t, response.Token)
	mustEqual(t, "Client User", response.User.Name)
	mustEqual(t, "client@example.com", response.User.Email)
	mustEqual(t, "CLIENT", response.User.Role)

	var user domain.User
	mustNoError(t, db.Where("email = ?", "client@example.com").First(&user).Error)
	mustEqual(t, "CLIENT", user.Role)
	mustNotEmpty(t, user.PasswordHash)
	mustNotEqual(t, "secret123", user.PasswordHash)

	body, err := json.Marshal(response)
	mustNoError(t, err)
	mustNotContains(t, string(body), "password_hash")
	mustNotContains(t, string(body), user.PasswordHash)
}

func TestRegisterDuplicateEmailReturnsError(t *testing.T) {
	db := newTestDB(t)
	createTestUser(t, db, func(user *domain.User) { user.Email = "duplicate@example.com" })

	_, err := services.Register(db, services.RegisterRequest{
		Name:     "Other User",
		Email:    "duplicate@example.com",
		Password: "secret123",
	})

	mustErrorIs(t, err, services.ErrEmailAlreadyExists)
}

func TestRegisterShortPasswordReturnsError(t *testing.T) {
	db := newTestDB(t)

	_, err := services.Register(db, services.RegisterRequest{
		Name:     "Client User",
		Email:    "client@example.com",
		Password: "12345",
	})

	mustErrorIs(t, err, services.ErrInvalidRequest)
}

func TestLoginSuccessReturnsTokenAndUser(t *testing.T) {
	db := newTestDB(t)
	t.Setenv("JWT_SECRET", "test-secret")
	createTestUser(t, db, func(user *domain.User) {
		user.Name = "Login User"
		user.Email = "login@example.com"
	})

	response, err := services.Login(db, services.LoginRequest{
		Email:    "login@example.com",
		Password: "secret123",
	})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustNotEmpty(t, response.Token)
	mustEqual(t, "Login User", response.User.Name)
	mustEqual(t, "login@example.com", response.User.Email)
	mustEqual(t, "CLIENT", response.User.Role)
}

func TestLoginIncorrectPasswordReturnsError(t *testing.T) {
	db := newTestDB(t)
	createTestUser(t, db, func(user *domain.User) { user.Email = "login@example.com" })

	_, err := services.Login(db, services.LoginRequest{
		Email:    "login@example.com",
		Password: "wrong-password",
	})

	mustTrue(t, errors.Is(err, services.ErrInvalidCredentials))
}
