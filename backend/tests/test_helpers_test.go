package tests

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/routes"
	"github.com/OctiDelFabro/ticketapp/backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	mustNoError(t, err)

	sqlDB, err := db.DB()
	mustNoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() {
		mustNoError(t, sqlDB.Close())
	})

	mustNoError(t, db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}))
	return db
}

func createTestEvent(t *testing.T, db *gorm.DB, attrs ...func(*domain.Event)) domain.Event {
	t.Helper()

	event := domain.Event{
		Title:           fmt.Sprintf("Test Event %d", time.Now().UnixNano()),
		Description:     "A test event description",
		ImageURL:        "https://example.com/event.jpg",
		Category:        "Music",
		Location:        "Buenos Aires",
		StartDate:       time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second),
		DurationMinutes: 120,
		Capacity:        10,
		Active:          true,
	}
	for _, attr := range attrs {
		attr(&event)
	}

	mustNoError(t, db.Create(&event).Error)
	return event
}

func createTestUser(t *testing.T, db *gorm.DB, attrs ...func(*domain.User)) domain.User {
	t.Helper()

	passwordHash, err := utils.HashPassword("secret123")
	mustNoError(t, err)

	user := domain.User{
		Name:         "Test User",
		Email:        fmt.Sprintf("user-%d@example.com", time.Now().UnixNano()),
		PasswordHash: passwordHash,
		Role:         "CLIENT",
	}
	for _, attr := range attrs {
		attr(&user)
	}

	mustNoError(t, db.Create(&user).Error)
	return user
}

func createTestTicket(t *testing.T, db *gorm.DB, userID uint, eventID uint, status string) domain.Ticket {
	t.Helper()

	ticket := domain.Ticket{
		UserID:       userID,
		EventID:      eventID,
		Status:       status,
		PurchaseDate: time.Now().UTC().Truncate(time.Second),
	}
	mustNoError(t, db.Create(&ticket).Error)
	mustNoError(t, db.Preload("User").Preload("Event").First(&ticket, ticket.ID).Error)
	return ticket
}

func generateTestToken(t *testing.T, user domain.User) string {
	t.Helper()

	t.Setenv("JWT_SECRET", "test-secret")
	token, err := utils.GenerateToken(user.ID, user.Email, user.Role)
	mustNoError(t, err)
	return token
}

func newTestRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	t.Helper()

	t.Setenv("JWT_SECRET", "test-secret")
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router, db)
	return router
}

func authHeader(t *testing.T, user domain.User) http.Header {
	t.Helper()

	header := http.Header{}
	header.Set("Authorization", "Bearer "+generateTestToken(t, user))
	return header
}

func mustNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func mustErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}

func mustNotNil(t *testing.T, value any) {
	t.Helper()
	if value == nil || reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
		t.Fatal("expected non-nil value")
	}
}

func mustEqual(t *testing.T, expected any, actual any) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v (%T), got %v (%T)", expected, expected, actual, actual)
	}
}

func mustNotEqual(t *testing.T, unexpected any, actual any) {
	t.Helper()
	if reflect.DeepEqual(unexpected, actual) {
		t.Fatalf("expected value different from %v", unexpected)
	}
}

func mustTrue(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Fatal("expected true")
	}
}

func mustLen(t *testing.T, value any, expected int) {
	t.Helper()
	actual := reflect.ValueOf(value).Len()
	if actual != expected {
		t.Fatalf("expected length %d, got %d", expected, actual)
	}
}

func mustNotEmpty(t *testing.T, value any) {
	t.Helper()
	if value == nil {
		t.Fatal("expected non-empty value")
	}
	if reflect.ValueOf(value).Len() == 0 {
		t.Fatal("expected non-empty value")
	}
}

func mustNotContains(t *testing.T, container any, item any) {
	t.Helper()
	switch typed := container.(type) {
	case string:
		needle, ok := item.(string)
		if !ok {
			t.Fatalf("expected string item for string container, got %T", item)
		}
		if strings.Contains(typed, needle) {
			t.Fatalf("expected %q not to contain %q", typed, needle)
		}
	case map[string]any:
		key, ok := item.(string)
		if !ok {
			t.Fatalf("expected string key for map container, got %T", item)
		}
		if _, exists := typed[key]; exists {
			t.Fatalf("expected map not to contain key %q", key)
		}
	default:
		t.Fatalf("unsupported container type %T", container)
	}
}

func mustMap(t *testing.T, value any) {
	t.Helper()
	if _, ok := value.(map[string]any); !ok {
		t.Fatalf("expected map[string]any, got %T", value)
	}
}
