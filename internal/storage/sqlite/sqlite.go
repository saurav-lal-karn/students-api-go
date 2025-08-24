package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saurav-lal-karn/students-api-go/internal/config"
	"github.com/saurav-lal-karn/students-api-go/internal/types"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
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

func (s Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// Prepare query to prevent from SQL injection
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	// Close the statement once the function completes
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	// Exec returns the last insert id and number of rows effected
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students where id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query Error: %w", err)
	}

	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var students []types.Student
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student types.Student
		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) (types.Student, error) {
	_, err := s.GetStudentById(id)
	if err != nil {
		return types.Student{}, err
	}

	// Prepare query to prevent from SQL injection
	stmt, err := s.Db.Prepare("UPDATE students SET name = ?,email = ?, age= ? WHERE id=?")
	if err != nil {
		return types.Student{}, err
	}
	// Close the statement once the function completes
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age, id)
	if err != nil {
		return types.Student{}, err
	}

	// Exec returns the last insert id and number of rows effected
	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return types.Student{}, err
	}

	if rowsUpdated > 0 {
		updatedStudent, err := s.GetStudentById(id)
		if err != nil {
			return types.Student{}, err
		}
		return updatedStudent, nil
	} else {
		return types.Student{}, fmt.Errorf("no rows found to be updated")
	}
}

func (s *Sqlite) DeleteStudent(id int64) error {
	_, err := s.GetStudentById(id)
	if err != nil {
		return err
	}
	// Prepare query to prevent from SQL injection
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id=?")
	if err != nil {
		return err
	}
	// Close the statement once the function completes
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
