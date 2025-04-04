package storage

import "github.com/NiteshSGupta/students-api/internal/types"

type Storage interface {
	// we created method belowe for storage means for database  (currntly we used sqlite)
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}
