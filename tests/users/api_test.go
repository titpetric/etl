package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/titpetric/platform"

	"github.com/titpetric/etl/server"
)

func TestCQRSAPIDirect(t *testing.T) {
	// Initialize database
	err := initializeDatabase()
	require.NoError(t, err, "failed to initialize database")

	// Load config
	conf, err := server.NewConfig()
	require.NoError(t, err, "failed to load config")

	// Create platform with the ETL module
	opts := platform.NewTestOptions()
	opts.Quiet = true
	svc := platform.New(opts)
	svc.Register(server.NewModule(conf))

	// Start the service (blocks until setup is complete, then returns)
	ctx, cancel := context.WithCancel(context.Background())
	err = svc.Start(ctx)
	require.NoError(t, err, "failed to start platform")

	// Cleanup: stop the service and cancel context
	t.Cleanup(func() {
		svc.Stop()
		cancel()
	})

	// Get the base URL from the platform
	baseURL := svc.URL()
	require.NotEmpty(t, baseURL, "platform URL should not be empty")

	// ====================
	// QUERY SIDE TESTS (READ)
	// ====================

	t.Run("Query/JSON/ListUsers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/users")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode, "expected 200, got %d", resp.StatusCode)

		var users []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&users)
		require.NoError(t, err)

		require.Len(t, users, 3, "expected 3 users, got %d", len(users))

		// Check first user
		require.Equal(t, "Alice Johnson", users[0]["name"])
		require.Equal(t, "alice@example.com", users[0]["email"])
	})

	t.Run("Query/JSON/GetUser", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/users/1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		var user map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&user)
		require.NoError(t, err)

		require.Equal(t, "1", user["id"])
		require.Equal(t, "Alice Johnson", user["name"])
	})

	t.Run("Query/JSON/GetUserNotFound", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/api/users/999")
		require.NoError(t, err)
		defer resp.Body.Close()

		// When no rows found, server typically returns 200 with null
		body, _ := io.ReadAll(resp.Body)
		t.Logf("Response status: %d, body: %s", resp.StatusCode, string(body))
	})

	t.Run("Query/HTML/ListUsers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		ct := resp.Header.Get("Content-Type")
		require.Equal(t, "text/html; charset=utf-8", ct)

		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)

		require.NotEmpty(t, bodyStr, "response body is empty")

		// Check for HTML structure
		require.Contains(t, bodyStr, "<table", "response should contain HTML table")
		require.Contains(t, bodyStr, "Alice Johnson", "response should contain user data")
	})

	t.Run("Query/HTML/GetUser", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users/1")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		require.Contains(t, string(body), "User Details", "response should contain user details title")
		require.Contains(t, string(body), "alice@example.com", "response should contain user email")
	})

	// ====================
	// COMMAND SIDE TESTS (WRITE)
	// ====================

	t.Run("Command/CreateUser", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"name":"David Brown","email":"david@example.com"}`)
		req, err := http.NewRequest("POST", baseURL+"/api/users", payload)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var user map[string]interface{}
		err = json.Unmarshal(body, &user)
		require.NoError(t, err)

		require.Equal(t, "David Brown", user["name"])
		require.Equal(t, "david@example.com", user["email"])
	})

	t.Run("Command/UpdateUser", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"name":"Alice Updated","email":"alice.updated@example.com"}`)
		req, err := http.NewRequest("PUT", baseURL+"/api/users/1", payload)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var user map[string]interface{}
		err = json.Unmarshal(body, &user)
		require.NoError(t, err)

		require.Equal(t, "Alice Updated", user["name"])
	})

	t.Run("Command/DeleteUser", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("DELETE", baseURL+"/api/users/3", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		// Verify deletion
		verifyResp, err := http.Get(baseURL + "/api/users/3")
		require.NoError(t, err)
		defer verifyResp.Body.Close()
		t.Logf("After delete, status: %d", verifyResp.StatusCode)
	})

	// ====================
	// ORDERS TESTS
	// ====================

	t.Run("Query/GetOrders", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/orders")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		var orders []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orders)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(orders), 1, "expected at least 1 order")
	})

	t.Run("Query/GetUserOrders", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/users/1/orders")
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		var orders []map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&orders)
		require.NoError(t, err)

		require.Len(t, orders, 2, "expected 2 orders for user 1")
	})

	t.Run("Command/CreateOrder", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"user_id":"1","total_amount":"50.00","status":"pending"}`)
		req, err := http.NewRequest("POST", baseURL+"/orders", payload)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var order map[string]interface{}
		err = json.Unmarshal(body, &order)
		require.NoError(t, err)

		require.Equal(t, "pending", order["status"])
	})

	// ====================
	// LOGIN/SESSION TESTS
	// ====================

	t.Run("Command/Login", func(t *testing.T) {
		client := &http.Client{}
		payload := bytes.NewBufferString(`{"user_id":"1","token":"test_token_12345"}`)
		req, err := http.NewRequest("POST", baseURL+"/login", payload)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)

		var session map[string]interface{}
		err = json.Unmarshal(body, &session)
		require.NoError(t, err)

		require.Equal(t, "test_token_12345", session["token"])
	})
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
