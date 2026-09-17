package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	fmt.Println((r.PathValue("id")))
}
func main() {
	http.HandleFunc("/students", getStudents)
	fmt.Println(students)
	http.HandleFunc("/students/{id}", getStudent)
	http.ListenAndServe(":8080", nil)
}
