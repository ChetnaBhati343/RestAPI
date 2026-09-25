package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var students = []Student{
	{ID: 1, Name: "Chetna Bhati", Age: 23, Email: "chetnabhati343@gmail.com"},
	{ID: 2, Name: "test01", Age: 82, Email: "test01@gmail.com"},
	{ID: 3, Name: "test02", Age: 41, Email: "test00@gmail.com"},
}

var nextID = 4

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println((r.PathValue("id")))
	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	found := false
	for _, student := range students {
		if student.ID == intId {
			found = true
			json.NewEncoder(w).Encode(student)

		}
	}
	if !found {
		http.Error(w, "Student not found", http.StatusNotFound)
	}
}

func postStudent(w http.ResponseWriter, r *http.Request) {
	var newStudent Student
	err := json.NewDecoder(r.Body).Decode(&newStudent)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newStudent.Name == "" {
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}

	if newStudent.Age <= 0 {
		http.Error(w, "Age must be greater than zero", http.StatusBadRequest)
		return
	}

	if newStudent.Email == "" {
		http.Error(w, "Email cannot be empty", http.StatusBadRequest)
		return
	}

	for _, student := range students {
		if newStudent.Email == student.Email {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
	}

	newStudent.ID = nextID
	nextID++
	students = append(students, newStudent)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newStudent)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)

	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	found := false
	var updatedStudent Student

	for i := range students {
		if students[i].ID == intId {
			found = true
			err := json.NewDecoder(r.Body).Decode(&updatedStudent)

			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			updatedStudent.ID = intId
			if updatedStudent.Age <= 0 {
				http.Error(w, "Age must be greater than zero", http.StatusBadRequest)
				return
			}
			if updatedStudent.Name == "" {
				http.Error(w, "Name cannot be empty", http.StatusBadRequest)
				return
			}
			if updatedStudent.Email == "" {
				http.Error(w, "Email cannot be empty", http.StatusBadRequest)
				return
			}
			for _, student := range students {
				if updatedStudent.Email == student.Email && updatedStudent.ID != student.ID {
					http.Error(w, "Email already exists", http.StatusConflict)
					return
				}
			}
			students[i] = updatedStudent
			break
		}
	}
	if !found {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	getId := r.PathValue("id")
	intId, err := strconv.Atoi(getId)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}
	found := false
	for i := range students {
		if students[i].ID == intId {
			found = true
			students = append(students[:i], students[i+1:]...)
			break
		}
	}
	if !found {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	http.Error(w, "Student deleted successfully.", http.StatusOK)

}

func main() {
	http.HandleFunc("GET /students", getStudents)
	http.HandleFunc("POST /students", postStudent)
	http.HandleFunc("PUT /students/{id}", updateStudent)
	http.HandleFunc("/students/{id}", getStudent)
	http.HandleFunc("DELETE /students/{id}", deleteStudent)
	http.ListenAndServe(":8081", nil)
}
