package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/middlewares"
	"github.com/OctiDelFabro/ticketapp/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createAdminUser(t *testing.T, db *gorm.DB) domain.User {
	t.Helper()
	return createTestUser(t, db, func(user *domain.User) { user.Role = "ADMIN" })
}

func createClientUser(t *testing.T, db *gorm.DB) domain.User {
	t.Helper()
	return createTestUser(t, db, func(user *domain.User) { user.Role = "CLIENT" })
}

func adminHeader(t *testing.T, admin domain.User) http.Header {
	t.Helper()
	return authHeader(t, admin)
}

func clientHeader(t *testing.T, client domain.User) http.Header {
	t.Helper()
	return authHeader(t, client)
}

func validEventBody(overrides map[string]any) map[string]any {
	body := map[string]any{
		"title":            "Admin Created Event",
		"description":      "Created from admin HTTP tests",
		"image_url":        "https://example.com/admin.jpg",
		"category":         "Música",
		"location":         "Buenos Aires",
		"start_date":       time.Now().Add(72 * time.Hour).UTC().Format(time.RFC3339),
		"duration_minutes": 90,
		"capacity":         10,
		"price":            25.5,
	}
	for key, value := range overrides {
		body[key] = value
	}
	return body
}

func TestAdminMiddlewareAuthorization(t *testing.T) {
	cases := []struct {
		name       string
		role       any
		setRole    bool
		wantStatus int
	}{
		{name: "missing role", wantStatus: http.StatusUnauthorized},
		{name: "empty role", setRole: true, role: "", wantStatus: http.StatusUnauthorized},
		{name: "client role", setRole: true, role: "CLIENT", wantStatus: http.StatusForbidden},
		{name: "admin role", setRole: true, role: "ADMIN", wantStatus: http.StatusOK},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(func(c *gin.Context) {
				if tt.setRole {
					c.Set("userRole", tt.role)
				}
			})
			router.Use(middlewares.AdminMiddleware())
			router.GET("/admin", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
			response := performJSONRequest(t, router, http.MethodGet, "/admin", nil)
			mustEqual(t, tt.wantStatus, response.Code)
		})
	}
}

func TestAdminEventsHTTPAuthorizationAndCreate(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	admin := createAdminUser(t, db)
	client := createClientUser(t, db)
	createTestEvent(t, db)
	mustEqual(t, http.StatusUnauthorized, performJSONRequest(t, router, http.MethodGet, "/api/admin/events", nil).Code)
	mustEqual(t, http.StatusForbidden, performJSONRequest(t, router, http.MethodGet, "/api/admin/events", nil, clientHeader(t, client)).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/admin/events", nil, adminHeader(t, admin)).Code)
	mustEqual(t, http.StatusCreated, performJSONRequest(t, router, http.MethodPost, "/api/admin/events", validEventBody(nil), adminHeader(t, admin)).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/admin/events", validEventBody(map[string]any{"category": "Invalid"}), adminHeader(t, admin)).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/admin/events", "not-an-object", adminHeader(t, admin)).Code)
	mustEqual(t, http.StatusForbidden, performJSONRequest(t, router, http.MethodPost, "/api/admin/events", validEventBody(nil), clientHeader(t, client)).Code)
}

func TestAdminEventsHTTPUpdateAndDelete(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	admin := createAdminUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	owner := createClientUser(t, db)
	secondOwner := createClientUser(t, db)
	createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, secondOwner.ID, event.ID, "ACTIVE")
	header := adminHeader(t, admin)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/admin/events/%d", event.ID), map[string]any{"title": "Updated", "category": "Teatro", "price": 10, "capacity": 2}, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, "/api/admin/events/abc", map[string]any{"title": "x"}, header).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodPatch, "/api/admin/events/99999", map[string]any{"title": "x"}, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/admin/events/%d", event.ID), map[string]any{"category": "Invalid"}, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/admin/events/%d", event.ID), map[string]any{"capacity": 1}, header).Code)
	toDelete := createTestEvent(t, db)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodDelete, fmt.Sprintf("/api/admin/events/%d", toDelete.ID), nil, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodDelete, "/api/admin/events/abc", nil, header).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodDelete, "/api/admin/events/99999", nil, header).Code)
	var disabled domain.Event
	mustNoError(t, db.First(&disabled, toDelete.ID).Error)
	mustEqual(t, false, disabled.Active)
}

func TestAdminStatsHTTP(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	admin := createAdminUser(t, db)
	client := createClientUser(t, db)
	event := createTestEvent(t, db)
	header := adminHeader(t, admin)
	mustEqual(t, http.StatusUnauthorized, performJSONRequest(t, router, http.MethodGet, "/api/admin/stats/summary", nil).Code)
	mustEqual(t, http.StatusForbidden, performJSONRequest(t, router, http.MethodGet, "/api/admin/stats/summary", nil, clientHeader(t, client)).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/admin/stats/summary", nil, header).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/admin/stats/events", nil, header).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/api/admin/events/%d/report", event.ID), nil, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodGet, "/api/admin/events/abc/report", nil, header).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodGet, "/api/admin/events/99999/report", nil, header).Code)
}

func TestTicketGiftHTTP(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	giver := createClientUser(t, db)
	target := createClientUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	body := map[string]any{"event_id": event.ID, "target_email": target.Email, "message": "enjoy"}
	mustEqual(t, http.StatusUnauthorized, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", body).Code)
	mustEqual(t, http.StatusCreated, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", body, authHeader(t, giver)).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", "bad", authHeader(t, giver)).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", map[string]any{"event_id": event.ID, "target_email": "missing@example.com"}, authHeader(t, giver)).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", map[string]any{"event_id": event.ID, "target_email": giver.Email}, authHeader(t, giver)).Code)
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", body, authHeader(t, giver)).Code)
	emptyEvent := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 1 })
	createTestTicket(t, db, giver.ID, emptyEvent.ID, "ACTIVE")
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPost, "/api/tickets/gift", map[string]any{"event_id": emptyEvent.ID, "target_email": createClientUser(t, db).Email}, authHeader(t, giver)).Code)
}

func TestTicketCancelAndTransferHTTPErrorPaths(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	owner := createClientUser(t, db)
	other := createClientUser(t, db)
	target := createClientUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 5 })
	ownerTicket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	otherTicket := createTestTicket(t, db, other.ID, event.ID, "ACTIVE")
	cancelled := createTestTicket(t, db, owner.ID, event.ID, "CANCELLED")
	header := authHeader(t, owner)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, "/api/tickets/abc/cancel", nil, header).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodPatch, "/api/tickets/99999/cancel", nil, header).Code)
	mustEqual(t, http.StatusForbidden, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/cancel", otherTicket.ID), nil, header).Code)
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/cancel", cancelled.ID), nil, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, "/api/tickets/abc/transfer", map[string]string{"target_email": target.Email}, header).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", ownerTicket.ID), "bad", header).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", ownerTicket.ID), map[string]string{"target_email": "missing@example.com"}, header).Code)
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", ownerTicket.ID), map[string]string{"target_email": owner.Email}, header).Code)
	mustEqual(t, http.StatusForbidden, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", otherTicket.ID), map[string]string{"target_email": target.Email}, header).Code)
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", cancelled.ID), map[string]string{"target_email": target.Email}, header).Code)
}

func TestAuthControllerHTTPErrorPaths(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	user := createTestUser(t, db, func(user *domain.User) { user.Email = "duplicate@example.com" })
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/auth/register", "bad").Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/auth/register", map[string]string{"name": "", "email": "", "password": ""}).Code)
	mustEqual(t, http.StatusConflict, performJSONRequest(t, router, http.MethodPost, "/api/auth/register", map[string]string{"name": "Duplicate", "email": user.Email, "password": "secret123"}).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodPost, "/api/auth/login", "bad").Code)
	mustEqual(t, http.StatusUnauthorized, performJSONRequest(t, router, http.MethodPost, "/api/auth/login", map[string]string{"email": "missing@example.com", "password": "secret123"}).Code)
	mustEqual(t, http.StatusUnauthorized, performJSONRequest(t, router, http.MethodPost, "/api/auth/login", map[string]string{"email": user.Email, "password": "wrongpass"}).Code)
}

func TestPublicEventsHTTPFiltersAndErrors(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	music := createTestEvent(t, db, func(event *domain.Event) {
		event.Title = "Searchable Music"
		event.Category = "Música"
		event.Capacity = 1
	})
	createTestEvent(t, db, func(event *domain.Event) { event.Category = "Teatro" })
	inactive := createTestEvent(t, db)
	mustNoError(t, services.DisableEvent(db, inactive.ID))
	user := createClientUser(t, db)
	createTestTicket(t, db, user.ID, music.ID, "ACTIVE")
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/events?search=Searchable", nil).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/events?category=Teatro", nil).Code)
	mustEqual(t, http.StatusOK, performJSONRequest(t, router, http.MethodGet, "/api/events?available_only=true", nil).Code)
	mustEqual(t, http.StatusBadRequest, performJSONRequest(t, router, http.MethodGet, "/api/events/abc", nil).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodGet, "/api/events/99999", nil).Code)
	mustEqual(t, http.StatusNotFound, performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/api/events/%d", inactive.ID), nil).Code)
}

func TestEventAndStatsServicesAdditionalCoverage(t *testing.T) {
	db := newTestDB(t)
	event := createTestEvent(t, db)
	mustNoError(t, services.DisableEvent(db, event.ID))
	var disabled domain.Event
	mustNoError(t, db.First(&disabled, event.ID).Error)
	mustEqual(t, false, disabled.Active)
	mustErrorIs(t, services.DisableEvent(db, 99999), services.ErrEventNotFound)
	badDuration := 0
	_, err := services.UpdateEvent(db, event.ID, services.UpdateEventRequest{DurationMinutes: &badDuration})
	mustErrorIs(t, err, services.ErrInvalidEventRequest)
	negativePrice := -1.0
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Price: &negativePrice})
	mustErrorIs(t, err, services.ErrInvalidEventRequest)
	blank := "   "
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Title: &blank})
	mustErrorIs(t, err, services.ErrInvalidEventRequest)
	event = createTestEvent(t, db)
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Description: &blank})
	mustErrorIs(t, err, services.ErrInvalidEventRequest)
	event = createTestEvent(t, db)
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Location: &blank})
	mustErrorIs(t, err, services.ErrInvalidEventRequest)
	event = createTestEvent(t, db)
	active := false
	updated, err := services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Active: &active})
	mustNoError(t, err)
	mustEqual(t, false, updated.Active)
	reportEvent := createTestEvent(t, db)
	report, err := services.GetEventReport(db, reportEvent.ID)
	mustNoError(t, err)
	mustEqual(t, reportEvent.ID, report.EventID)
	_, err = services.GetEventReport(db, 99999)
	mustErrorIs(t, err, services.ErrEventNotFound)
	summary, err := services.GetStatsSummary(newTestDB(t))
	mustNoError(t, err)
	mustEqual(t, float64(0), summary.OccupancyRatePercent)
}
