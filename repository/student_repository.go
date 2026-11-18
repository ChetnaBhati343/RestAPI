package repository

import (
	"context"
	"database/sql"

	"github.com/ChetnaBhati343/RestAPI/models"
)

type StudentRepository struct {
	db *sql.DB
}

// NewStudentRepository creates a new StudentRepository.
func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{
		db: db,
	}
}

// GetAllStudents retrieves all student records from PostgreSQL.
func (r *StudentRepository) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	query := `SELECT id, name, age, email, created_at
	          FROM students
	          ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.Age,
			&student.Email,
			&student.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByID retrieves a student using their ID.
func (r *StudentRepository) GetStudentByID(ctx context.Context, id int) (models.Student, error) {

	var student models.Student

	query := `SELECT id, name, age, email, created_at
	          FROM students
	          WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&student.ID,
		&student.Name,
		&student.Age,
		&student.Email,
		&student.CreatedAt,
	)

	if err != nil {
		return models.Student{}, err
	}

	return student, nil
}

// CreateStudent inserts a new student record into PostgreSQL.
func (r *StudentRepository) CreateStudent(ctx context.Context, student models.Student) (models.Student, error) {

	query := `INSERT INTO students (name, age, email)
	          VALUES ($1, $2, $3)
	          RETURNING id, name, age, email, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		student.Name,
		student.Age,
		student.Email,
	).Scan(
		&student.ID,
		&student.Name,
		&student.Age,
		&student.Email,
		&student.CreatedAt,
	)

	if err != nil {
		return models.Student{}, err
	}

	return student, nil
}

// UpdateStudent updates an existing student record in PostgreSQL.
func (r *StudentRepository) UpdateStudent(ctx context.Context, id int, student models.Student) (models.Student, error) {

	query := `UPDATE students
              SET name = $1, age = $2, email = $3
              WHERE id = $4
              RETURNING id, name, age, email, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		student.Name,
		student.Age,
		student.Email,
		id,
	).Scan(
		&student.ID,
		&student.Name,
		&student.Age,
		&student.Email,
		&student.CreatedAt,
	)

	if err != nil {
		return models.Student{}, err
	}

	return student, nil
}

// DeleteStudent deletes a student record from PostgreSQL.
func (r *StudentRepository) DeleteStudent(ctx context.Context, id int) error {
	query := `DELETE FROM students
			  WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
