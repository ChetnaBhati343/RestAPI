package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChetnaBhati343/RestAPI/database"
	"github.com/ChetnaBhati343/RestAPI/repository"
)

func TestGetStudentsWithCancelledRequestContext(t *testing.T) {
	// Connect to PostgreSQL.
	db, err := database.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize the student repository.
	studentRepo := repository.NewStudentRepository(db)

	// Create a context that can be cancelled manually.
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context before sending the request to the handler.
	cancel()

	// Create a simulated HTTP GET request.
	req := httptest.NewRequest(
		http.MethodGet,
		"/students",
		nil,
	).WithContext(ctx)

	// Record the HTTP response.
	recorder := httptest.NewRecorder()

	// Execute the actual getStudents handler.
	getStudents(recorder, req, studentRepo)

	// Verify that the handler received the database error.
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"Expected status %d, got %d",
			http.StatusInternalServerError,
			recorder.Code,
		)
	}
}