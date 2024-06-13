package models

import (
	"database/sql"
	"errors"
	"strings"
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
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, hashed_password, created)
	VALUES (?, ?, ?, UTC_TIMESTAMP());`

	_, err = m.DB.Exec(query, name, email, string(hashed_password))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	query := `SELECT id, name, email, hashed_password from users
	WHERE email=?`

	row := m.DB.QueryRow(query, email)

	user := &User{}

	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Hashed_password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}

	}

	err = bcrypt.CompareHashAndPassword(user.Hashed_password, []byte(password))

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return int(user.Id), nil
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
