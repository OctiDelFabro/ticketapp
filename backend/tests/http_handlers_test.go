package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
)

func performJSONRequest(t *testing.T, handler http.Handler, method string, path string, body any, headers ...http.Header) *httptest.ResponseRecorder {
	t.Helper()

	var payload bytes.Buffer
	if body != nil {
		mustNoError(t, json.NewEncoder(&payload).Encode(body))
	}

	request := httptest.NewRequest(method, path, &payload)
	request.Header.Set("Content-Type", "application/json")
	for _, header := range headers {
		for key, values := range header {
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}
	}

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func TestHealthEndpointReturnsOK(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)

	response := performJSONRequest(t, router, http.MethodGet, "/api/health", nil)

	mustEqual(t, http.StatusOK, response.Code)
}

func TestRegisterEndpointReturnsCreatedWithBody(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)

	response := performJSONRequest(t, router, http.MethodPost, "/api/auth/register", map[string]string{
		"name":     "HTTP Client",
		"email":    "http-client@example.com",
		"password": "secret123",
	})

	mustEqual(t, http.StatusCreated, response.Code)
	body := decodeJSONResponse(t, response)
	mustNotEmpty(t, body["token"])
	mustMap(t, body["user"])
	user := body["user"].(map[string]any)
	mustEqual(t, "http-client@example.com", user["email"])
	mustNotContains(t, user, "password_hash")
}

func TestLoginEndpointReturnsOKWithValidCredentials(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	createTestUser(t, db, func(user *domain.User) { user.Email = "login-http@example.com" })

	response := performJSONRequest(t, router, http.MethodPost, "/api/auth/login", map[string]string{
		"email":    "login-http@example.com",
		"password": "secret123",
	})

	mustEqual(t, http.StatusOK, response.Code)
	body := decodeJSONResponse(t, response)
	mustNotEmpty(t, body["token"])
}

func TestListEventsEndpointReturnsOK(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	createTestEvent(t, db)

	response := performJSONRequest(t, router, http.MethodGet, "/api/events", nil)

	mustEqual(t, http.StatusOK, response.Code)
}

func TestGetEventByIDEndpointReturnsOKForExistingEvent(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	event := createTestEvent(t, db)

	response := performJSONRequest(t, router, http.MethodGet, fmt.Sprintf("/api/events/%d", event.ID), nil)

	mustEqual(t, http.StatusOK, response.Code)
	body := decodeJSONResponse(t, response)
	mustEqual(t, float64(event.ID), body["id"])
}

func TestPurchaseTicketEndpointWithoutTokenReturnsUnauthorized(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	event := createTestEvent(t, db)

	response := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]uint{"event_id": event.ID})

	mustEqual(t, http.StatusUnauthorized, response.Code)
}

func TestTicketAuthenticatedHTTPFlow(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	user := createTestUser(t, db, func(user *domain.User) { user.Email = "ticket-owner@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "ticket-target@example.com" })
	firstEvent := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 3 })
	secondEvent := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	header := authHeader(t, user)

	purchase := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]any{"event_id": firstEvent.ID, "quantity": 2}, header)
	mustEqual(t, http.StatusCreated, purchase.Code)
	purchaseBody := decodeJSONResponse(t, purchase)
	mustEqual(t, float64(2), purchaseBody["quantity"])
	tickets := purchaseBody["tickets"].([]any)
	mustLen(t, tickets, 2)
	firstTicket := tickets[0].(map[string]any)
	mustEqual(t, float64(firstEvent.ID), firstTicket["event_id"])
	mustEqual(t, "ACTIVE", firstTicket["status"])

	myTickets := performJSONRequest(t, router, http.MethodGet, "/api/tickets/me", nil, header)
	mustEqual(t, http.StatusOK, myTickets.Code)

	cancel := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/cancel", uint(firstTicket["id"].(float64))), nil, header)
	mustEqual(t, http.StatusOK, cancel.Code)
	cancelBody := decodeJSONResponse(t, cancel)
	mustEqual(t, "CANCELLED", cancelBody["status"])

	transferPurchase := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]uint{"event_id": secondEvent.ID}, header)
	mustEqual(t, http.StatusCreated, transferPurchase.Code)
	transferPurchaseBody := decodeJSONResponse(t, transferPurchase)
	transferTickets := transferPurchaseBody["tickets"].([]any)
	transferTicket := transferTickets[0].(map[string]any)

	transfer := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", uint(transferTicket["id"].(float64))), map[string]string{"target_email": target.Email}, header)
	mustEqual(t, http.StatusOK, transfer.Code)
	transferBody := decodeJSONResponse(t, transfer)
	mustEqual(t, float64(target.ID), transferBody["user_id"])
	mustEqual(t, "ACTIVE", transferBody["status"])
}

func TestAdminCannotUseTicketEndpoints(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	admin := createTestUser(t, db, func(user *domain.User) {
		user.Email = "admin-ticket-blocked@example.com"
		user.Role = "ADMIN"
	})
	event := createTestEvent(t, db)

	response := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]any{"event_id": event.ID, "quantity": 1}, authHeader(t, admin))

	mustEqual(t, http.StatusForbidden, response.Code)
}

func TestClientPassesTicketRoleMiddleware(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	client := createTestUser(t, db, func(user *domain.User) { user.Email = "client-ticket-allowed@example.com" })
	event := createTestEvent(t, db)

	response := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]any{"event_id": event.ID, "quantity": 1}, authHeader(t, client))

	mustEqual(t, http.StatusCreated, response.Code)
}

func TestAdminCanStillUseAdminEndpoints(t *testing.T) {
	db := newTestDB(t)
	router := newTestRouter(t, db)
	admin := createTestUser(t, db, func(user *domain.User) {
		user.Email = "admin-still-allowed@example.com"
		user.Role = "ADMIN"
	})

	response := performJSONRequest(t, router, http.MethodGet, "/api/admin/stats/summary", nil, authHeader(t, admin))

	mustEqual(t, http.StatusOK, response.Code)
}
