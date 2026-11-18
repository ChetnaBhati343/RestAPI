package database

import (
	"errors"
	"sync"
	"testing"
	"context"
	"time"
	"database/sql"

)

func TestConnectionPoolConfiguration(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	stats := db.Stats()

	if stats.MaxOpenConnections != 25 {
		t.Errorf(
			"Expected maximum open connections to be 25, got %d",
			stats.MaxOpenConnections,
		)
	}

	if stats.Idle < 0 {
		t.Errorf(
			"Expected idle connections to be non-negative, got %d",
			stats.Idle,
		)
	}
}

func TestConcurrentDatabaseConnections(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	const numberOfGoroutines = 10

	var wg sync.WaitGroup

	wg.Add(numberOfGoroutines)

	for i := 0; i < numberOfGoroutines; i++ {
		go func() {
			defer wg.Done()

			var result int

			err := db.QueryRow("SELECT 1").Scan(&result)
			if err != nil {
				t.Errorf("Database query failed: %v", err)
				return
			}

			if result != 1 {
				t.Errorf("Expected result 1, got %d", result)
			}
		}()
	}

	wg.Wait()
}


func TestConcurrentConnectionUsage(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	const numberOfConnections = 10

	ready := make(chan struct{}, numberOfConnections)
	release := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(numberOfConnections)

	for i := 0; i < numberOfConnections; i++ {
		go func() {
			defer wg.Done()

			conn, err := db.Conn(context.Background())
			if err != nil {
				t.Errorf("Failed to acquire database connection: %v", err)
				return
			}
			defer conn.Close()

			ready <- struct{}{}

			<-release
		}()
	}

	for i := 0; i < numberOfConnections; i++ {
		<-ready
	}

	stats := db.Stats()

	t.Logf("Maximum Open Connections: %d", stats.MaxOpenConnections)
	t.Logf("Open Connections: %d", stats.OpenConnections)
	t.Logf("In-Use Connections: %d", stats.InUse)
	t.Logf("Idle Connections: %d", stats.Idle)

	close(release)
	wg.Wait()

	stats = db.Stats()

	t.Logf("After releasing connections:")
	t.Logf("Open Connections: %d", stats.OpenConnections)
	t.Logf("In-Use Connections: %d", stats.InUse)
	t.Logf("Idle Connections: %d", stats.Idle)
}


func TestConnectionPoolWait(t *testing.T) {
	db, err := Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Temporarily reduce the connection limit for this test.
	db.SetMaxOpenConns(5)

	// Set a timeout to prevent the test from waiting indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Acquire all five available connections.
	connections := make([]*sql.Conn, 0, 5)

	for i := 0; i < 5; i++ {
		conn, err := db.Conn(ctx)
		if err != nil {
			t.Fatalf("Failed to acquire connection %d: %v", i+1, err)
		}

		connections = append(connections, conn)
	}

	// Start a sixth request in another goroutine.
	sixthResult := make(chan error, 1)
	started := make(chan struct{})

	go func() {
		close(started)

		conn, err := db.Conn(ctx)
		if err != nil {
			sixthResult <- err
			return
		}

		conn.Close()
		sixthResult <- nil
	}()

	<-started

	// Wait until the connection pool records a waiting request.
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.NewTimer(2 * time.Second)
	defer timeout.Stop()

	for db.Stats().WaitCount == 0 {
		select {
		case <-ticker.C:
		case <-timeout.C:
			t.Fatal("Sixth connection request did not enter the wait queue")
		}
	}

	// Inspect the pool while the sixth request is waiting.
	stats := db.Stats()

	t.Logf("Maximum Open Connections: %d", stats.MaxOpenConnections)
	t.Logf("Open Connections: %d", stats.OpenConnections)
	t.Logf("In-Use Connections: %d", stats.InUse)
	t.Logf("Idle Connections: %d", stats.Idle)
	t.Logf("Wait Count: %d", stats.WaitCount)

	// Release one connection so the sixth request can proceed.
	connections[0].Close()
	connections[0] = nil

	err = <-sixthResult
	if err != nil {
		t.Fatalf("Sixth connection request failed: %v", err)
	}

	// Release the remaining connections.
	for _, conn := range connections {
		if conn != nil {
			conn.Close()
		}
	}
}


func TestQueryContextCancellation(t *testing.T) {
    db, err := Connect()
    if err != nil {
        t.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Create a context with a 100-millisecond timeout.
    ctx, cancel := context.WithTimeout(
        context.Background(),
        100*time.Millisecond,
    )
    defer cancel()

    var result int

    // PostgreSQL normally takes 5 seconds to execute this query.
    err = db.QueryRowContext(
        ctx,
        "SELECT 1 FROM pg_sleep(5)",
    ).Scan(&result)

    // The query should be cancelled when the context deadline expires.
    if !errors.Is(err, context.DeadlineExceeded) {
        t.Fatalf("Expected context deadline exceeded, got: %v", err)
    }
}