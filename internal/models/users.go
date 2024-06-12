package models

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int32
	Name            string
	Email           string
	Hashed_password []byte
	Created         time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name string, email string, password string) error {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP());`
	row, err := m.DB.Exec(query, name, email, hashed_password)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	query := `SELECT id, name, email, hashed_password from users
	WHERE email=?`

	row := m.DB.QueryRow(query, email, password)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Hashed_password)
	if err != nil {
		return -1, nil
	}

	err = bcrypt.CompareHashAndPassword(user.Hashed_password, []byte(password))
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	query := `SELECT id, name, email, created from users
	WHERE id=?`

	row := m.DB.QueryRow(query, id)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Created)

	if err != nil {
		return false, err
	}

	return true, nil
}
