package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ChetnaBhati343/RestAPI/database"
	_ "github.com/ChetnaBhati343/RestAPI/docs"
	"github.com/ChetnaBhati343/RestAPI/models"
	"github.com/ChetnaBhati343/RestAPI/repository"
	"github.com/jackc/pgx/v5/pgconn"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"os"
)

// @Summary Get all students
// @Description Returns a list of all students
// @Tags Students
// @Produce json
// @Success 200 {array} models.Student
// @Failure 500 {string} string "Failed to retrieve students"
// @Router /students [get]
func getStudents(w http.ResponseWriter, r *http.Request, studentRepo *repository.StudentRepository) {
	allStudents, err := studentRepo.GetAllStudents(r.Context())
	if err != nil {
		log.Printf("Failed to retrieve students: %v", err)
		http.Error(w, "Failed to retrieve students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allStudents)
}

// @Summary Get student by ID
// @Description Returns a student with the specified ID
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Invalid student ID"
// @Failure 404 {string} string "Student not found"
// @Failure 500 {string} string "Failed to retrieve student"
// @Router /students/{id} [get]
func getStudent(w http.ResponseWriter, r *http.Request, studentRepo *repository.StudentRepository) {
	w.Header().Set("Content-Type", "application/json")

	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)

	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, err := studentRepo.GetStudentByID(r.Context(), intId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		log.Printf("Failed to retrieve student: %v", err)
		http.Error(w, "Failed to retrieve student", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(student)
}

// @Summary Create a new student
// @Description Creates a new student record
// @Tags Students
// @Accept json
// @Produce json
// @Param student body models.CreateStudentRequest true "Student details"
// @Success 201 {object} models.Student
// @Failure 400 {string} string "Invalid JSON or invalid student details"
// @Failure 409 {string} string "Email already exists"
// @Failure 500 {string} string "Failed to create student"
// @Router /students [post]
func postStudent(w http.ResponseWriter, r *http.Request, studentRepo *repository.StudentRepository) {
	var request models.CreateStudentRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	if request.Age <= 0 {
		http.Error(w, "Age must be greater than zero", http.StatusBadRequest)
		return
	}

	if request.Email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	newStudent := models.Student{
		Name:  request.Name,
		Age:   request.Age,
		Email: request.Email,
	}

	createdStudent, err := studentRepo.CreateStudent(r.Context(), newStudent)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		log.Printf("Failed to create student: %v", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStudent)
}

// @Summary Update an existing student
// @Description Updates the details of an existing student
// @Tags Students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body models.UpdateStudentRequest true "Updated student details"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Invalid student ID, JSON or student details"
// @Failure 404 {string} string "Student not found"
// @Failure 409 {string} string "Email already exists"
// @Failure 500 {string} string "Failed to update student"
// @Router /students/{id} [put]
func updateStudent(w http.ResponseWriter, r *http.Request, studentRepo *repository.StudentRepository) {
	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)

	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var request models.UpdateStudentRequest

	err = json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	if request.Age <= 0 {
		http.Error(w, "Age must be greater than zero", http.StatusBadRequest)
		return
	}

	if request.Email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	updatedStudent := models.Student{
		Name:  request.Name,
		Age:   request.Age,
		Email: request.Email,
	}

	updatedStudent, err = studentRepo.UpdateStudent(r.Context(), intId, updatedStudent)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}

		log.Printf("Failed to update student: %v", err)
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

// @Summary Delete student by ID
// @Description Deletes a student with the specified ID
// @Tags Students
// @Produce json
// @Param id path int true "Student ID"
// @Success 204 "Student deleted successfully"
// @Failure 400 {string} string "Invalid student ID"
// @Failure 404 {string} string "Student not found"
// @Failure 500 {string} string "Failed to delete student"
// @Router /students/{id} [delete]
func deleteStudent(w http.ResponseWriter, r *http.Request, studentRepo *repository.StudentRepository) {
	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)

	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	err = studentRepo.DeleteStudent(r.Context(), intId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		log.Printf("Failed to delete student: %v", err)
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @title Student REST API
// @version 1.0
// @description REST API for managing student records.
func main() {
	db, err := database.Connect()

	if err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	defer db.Close()

	log.Println("Connected to PostgreSQL successfully")

	studentRepo := repository.NewStudentRepository(db)

	// Register student API routes.
	http.HandleFunc("GET /students", func(w http.ResponseWriter, r *http.Request) {
		getStudents(w, r, studentRepo)
	})

	http.HandleFunc("POST /students", func(w http.ResponseWriter, r *http.Request) {
		postStudent(w, r, studentRepo)
	})

	http.HandleFunc("GET /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		getStudent(w, r, studentRepo)
	})

	http.HandleFunc("PUT /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateStudent(w, r, studentRepo)
	})

	http.HandleFunc("DELETE /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteStudent(w, r, studentRepo)
	})

	// Register Swagger UI route.
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Read the IP address and port provided by the hosting platform.
ip := os.Getenv("IP")
port := os.Getenv("PORT")

// Use local development defaults when environment variables are unavailable.
if ip == "" {
	ip = "localhost"
}

if port == "" {
	port = "8081"
}

// JoinHostPort correctly handles IPv6 addresses.
address := net.JoinHostPort(ip, port)

log.Printf("Server starting on %s", address)
log.Fatal(http.ListenAndServe(address, nil))
}
