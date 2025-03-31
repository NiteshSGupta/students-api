package storage

type Storage interface {
	// we created method belowe for storage means for database  (currntly we used sqlite)
	CreateStudent(name string, email string, age int) (int64, error)
}
