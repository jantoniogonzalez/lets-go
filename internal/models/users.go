package models

import (
	"database/sql"
	"time"
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

	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP());`
	row, err := m.DB.Exec(query, name, email, password)
	if err != nil {

	}
	return nil
}

func (m *UserModel) Authenticate(email string, password string) (int, error) {
	query := `SELECT id, name, email from users
	WHERE email=? AND hashed_password=password`

	row := m.DB.QueryRow(query, email, password)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return -1, nil
	}
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
