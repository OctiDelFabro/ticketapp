package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/middlewares"
	"github.com/OctiDelFabro/ticketapp/backend/utils"
	"github.com/gin-gonic/gin"
)

func TestPasswordHashAndCheck(t *testing.T) {
	hash, err := utils.HashPassword("secret123")
	mustNoError(t, err)
	mustNotEmpty(t, hash)
	mustTrue(t, utils.CheckPasswordHash("secret123", hash))
	mustEqual(t, false, utils.CheckPasswordHash("wrong", hash))
}

func TestJWTGenerateValidateAndRejectInvalid(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	token, err := utils.GenerateToken(42, "jwt@example.com", "ADMIN")
	mustNoError(t, err)
	claims, err := utils.ValidateToken(token)
	mustNoError(t, err)
	mustEqual(t, uint(42), claims.UserID)
	mustEqual(t, "jwt@example.com", claims.Email)
	mustEqual(t, "ADMIN", claims.Role)

	_, err = utils.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestAuthMiddlewareRejectsMissingTokenAndAcceptsValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", middlewares.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.GetUint("userID"), "email": c.GetString("userEmail")})
	})

	missing := httptest.NewRecorder()
	router.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/protected", nil))
	mustEqual(t, http.StatusUnauthorized, missing.Code)

	token, err := utils.GenerateToken(7, "middleware@example.com", "CLIENT")
	mustNoError(t, err)
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	valid := httptest.NewRecorder()
	router.ServeHTTP(valid, request)
	mustEqual(t, http.StatusOK, valid.Code)
	body := decodeJSONResponse(t, valid)
	mustEqual(t, float64(7), body["user_id"])
	mustEqual(t, "middleware@example.com", body["email"])
}
