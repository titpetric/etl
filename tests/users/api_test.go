package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/titpetric/etl/server"
)

func TestCQRSAPIDirect(t *testing.T) {
	// Initialize database
	if err := initializeDatabase(); err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	// Create HTTP server handler directly
	handler, err := server.NewHandler()
	if err != nil {
		t.Fatalf("failed to create server handler: %v", err)
	}

	// Create test server
	ts := httptest.NewServer(handler)
	defer ts.Close()

	// ====================
	// QUERY SIDE TESTS (READ)
	// ====================

	t.Run("Query/JSON/ListUsers", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		var users []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}

		// Check first user
		if users[0]["name"] != "Alice Johnson" {
			t.Errorf("expected name 'Alice Johnson', got %v", users[0]["name"])
		}
		if users[0]["email"] != "alice@example.com" {
			t.Errorf("expected email 'alice@example.com', got %v", users[0]["email"])
		}
	})

	t.Run("Query/JSON/GetUser", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users/1")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		var user map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user["id"] != "1" {
			t.Errorf("expected id '1', got %v", user["id"])
		}
		if user["name"] != "Alice Johnson" {
			t.Errorf("expected name 'Alice Johnson', got %v", user["name"])
		}
	})

	t.Run("Query/JSON/GetUserNotFound", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/users/999")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		// When no rows found, server typically returns 200 with null
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response status: %d, body: %s", resp.StatusCode, string(body))
	})

	t.Run("Query/HTML/ListUsers", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		if ct := resp.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
			t.Errorf("expected Content-Type 'text/html; charset=utf-8', got %s", ct)
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		if len(bodyStr) == 0 {
			t.Error("response body is empty")
			return
		}

		// Check for HTML structure
		if !bytes.Contains(body, []byte("<table>")) {
			t.Errorf("response should contain HTML table, got: %s", bodyStr[:min(200, len(bodyStr))])
		}
		if !bytes.Contains(body, []byte("Alice Johnson")) {
			t.Errorf("response should contain user data, got: %s", bodyStr[:min(200, len(bodyStr))])
		}
	})

	t.Run("Query/HTML/GetUser", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users/1")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		if !bytes.Contains(body, []byte("User Details")) {
			t.Errorf("response should contain user details title, got: %s", bodyStr[:min(200, len(bodyStr))])
		}
		if !bytes.Contains(body, []byte("alice@example.com")) {
			t.Errorf("response should contain user email, got: %s", bodyStr[:min(200, len(bodyStr))])
		}
	})

	// ====================
	// COMMAND SIDE TESTS (WRITE)
	// ====================

	t.Run("Command/CreateUser", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"name":"David Brown","email":"david@example.com"}`)
		req, err := http.NewRequest("POST", ts.URL+"/api/users", payload)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)

		var user map[string]interface{}
		if err := json.Unmarshal(body, &user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user["name"] != "David Brown" {
			t.Errorf("expected name 'David Brown', got %v", user["name"])
		}
		if user["email"] != "david@example.com" {
			t.Errorf("expected email 'david@example.com', got %v", user["email"])
		}
	})

	t.Run("Command/UpdateUser", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"name":"Alice Updated","email":"alice.updated@example.com"}`)
		req, err := http.NewRequest("PUT", ts.URL+"/api/users/1", payload)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)

		var user map[string]interface{}
		if err := json.Unmarshal(body, &user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user["name"] != "Alice Updated" {
			t.Errorf("expected updated name, got %v", user["name"])
		}
	})

	t.Run("Command/DeleteUser", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", ts.URL+"/api/users/3", nil)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
		}

		// Verify deletion
		verifyResp, err := http.Get(ts.URL + "/api/users/3")
		if err != nil {
			t.Fatalf("verification request failed: %v", err)
		}
		defer verifyResp.Body.Close()
		t.Logf("After delete, status: %d", verifyResp.StatusCode)
	})

	// ====================
	// ORDERS TESTS
	// ====================

	t.Run("Query/GetOrders", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/orders")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		var orders []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(orders) < 1 {
			t.Errorf("expected at least 1 order, got %d", len(orders))
		}
	})

	t.Run("Query/GetUserOrders", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/users/1/orders")
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
			return
		}

		var orders []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(orders) != 2 {
			t.Errorf("expected 2 orders for user 1, got %d", len(orders))
		}
	})

	t.Run("Command/CreateOrder", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"user_id":"1","total_amount":"50.00","status":"pending"}`)
		req, err := http.NewRequest("POST", ts.URL+"/orders", payload)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)

		var order map[string]interface{}
		if err := json.Unmarshal(body, &order); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if order["status"] != "pending" {
			t.Errorf("expected status 'pending', got %v", order["status"])
		}
	})

	// ====================
	// LOGIN/SESSION TESTS
	// ====================

	t.Run("Command/Login", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"user_id":"1","token":"test_token_12345"}`)
		req, err := http.NewRequest("POST", ts.URL+"/login", payload)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
		}

		body, _ := io.ReadAll(resp.Body)

		var session map[string]interface{}
		if err := json.Unmarshal(body, &session); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if session["token"] != "test_token_12345" {
			t.Errorf("expected token to match, got %v", session["token"])
		}
	})
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// initializeDatabase sets up the test database with schema and seed data
func initializeDatabase() error {
	// Remove existing database
	os.Remove("test.db")

	// Open database
	db, err := sqlx.Open("sqlite", "file:test.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Read and execute schema
	schema, err := os.ReadFile("users.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return err
	}

	return nil
}
