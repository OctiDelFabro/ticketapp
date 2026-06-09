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

func decodeJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	mustNoError(t, json.Unmarshal(response.Body.Bytes(), &body))
	return body
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
	firstEvent := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	secondEvent := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	header := authHeader(t, user)

	purchase := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]uint{"event_id": firstEvent.ID}, header)
	mustEqual(t, http.StatusCreated, purchase.Code)
	purchaseBody := decodeJSONResponse(t, purchase)
	mustEqual(t, float64(firstEvent.ID), purchaseBody["event_id"])
	mustEqual(t, "ACTIVE", purchaseBody["status"])

	duplicate := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]uint{"event_id": firstEvent.ID}, header)
	mustEqual(t, http.StatusConflict, duplicate.Code)

	myTickets := performJSONRequest(t, router, http.MethodGet, "/api/tickets/me", nil, header)
	mustEqual(t, http.StatusOK, myTickets.Code)

	cancel := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/cancel", uint(purchaseBody["id"].(float64))), nil, header)
	mustEqual(t, http.StatusOK, cancel.Code)
	cancelBody := decodeJSONResponse(t, cancel)
	mustEqual(t, "CANCELLED", cancelBody["status"])

	transferPurchase := performJSONRequest(t, router, http.MethodPost, "/api/tickets/purchase", map[string]uint{"event_id": secondEvent.ID}, header)
	mustEqual(t, http.StatusCreated, transferPurchase.Code)
	transferPurchaseBody := decodeJSONResponse(t, transferPurchase)

	transfer := performJSONRequest(t, router, http.MethodPatch, fmt.Sprintf("/api/tickets/%d/transfer", uint(transferPurchaseBody["id"].(float64))), map[string]string{"target_email": target.Email}, header)
	mustEqual(t, http.StatusOK, transfer.Code)
	transferBody := decodeJSONResponse(t, transfer)
	mustEqual(t, float64(target.ID), transferBody["user_id"])
	mustEqual(t, "ACTIVE", transferBody["status"])
}
