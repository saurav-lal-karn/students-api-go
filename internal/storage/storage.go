package storage

import "github.com/saurav-lal-karn/students-api-go/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(id int64, name string, email string, age int) (types.Student, error)
	DeleteStudent(id int64) error
}
