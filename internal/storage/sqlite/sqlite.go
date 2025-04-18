package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/NiteshSGupta/students-api/internal/config"
	"github.com/NiteshSGupta/students-api/internal/types"

	//in sqlite we put underscroe in front , because we not using directly , it's working on background
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

// golang always return first value , second is error
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storagepath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	//this will save us form sql injection
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id=? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	//from here we have to pass the data from db, to stuct doing serialize
	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}

// func GetStudents(s *Sqlite) ([]types.Student, error) {
// 	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students")
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer stmt.Close()

// 	rows, err := stmt.Query()
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	//belowe slice going to store the data which come from database
// 	var students []types.Student

// 	for rows.Next() {
// 		var student types.Student
// 		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
// 		if err != nil {
// 			//do you know , why it returning the nil,err because out function returning the two things : []types.Student, error
// 			return nil, err
// 		}

// 		// so our loop will will get the data from rows Query , and pass to the : var students []types.Student
// 		students = append(students, student)
// 	}

// 	return students, nil

// }

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}
