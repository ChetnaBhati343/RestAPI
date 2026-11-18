
package models

// CreateStudentRequest represents the data required
// to create a new student.
type CreateStudentRequest struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}

// UpdateStudentRequest represents the data required
// to update an existing student.
type UpdateStudentRequest struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}
