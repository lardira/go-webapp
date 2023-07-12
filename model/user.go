package model

import (
	"database/sql"
	"errors"
	"time"
)

const (
	MIN_LOGIN_LENGTH    = 2
	MIN_PASSWORD_LENGTH = 5
)

type User struct {
	Id         int64     `db:"id"`
	Login      string    `db:"login"`
	Password   string    `db:"password"`
	LastAuth   time.Time `db:"last_auth"`
	LastUnauth time.Time `db:"last_unauth"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IsAuth     bool      `db:"is_auth"`
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *UserRequest) Validate() (err error) {
	if len(u.Login) < MIN_LOGIN_LENGTH {
		return errors.New("invalid user login data")
	}
	if len(u.Password) < MIN_PASSWORD_LENGTH {
		return errors.New("invalid user password data")
	}
	return nil
}

func CreateUser(db *sql.DB, login, password string) (User, error) {
	query := `
		INSERT INTO users
		(login, password, last_auth, last_unauth, created_at, updated_at, is_auth)
		VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	timeNow := time.Now()

	var id int64

	err := db.QueryRow(
		query,
		login,
		password,
		nil,
		nil,
		timeNow,
		timeNow,
		false,
	).Scan(&id)

	if err != nil {
		return User{}, err
	}

	return User{
			Id:        id,
			Login:     login,
			Password:  password,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			IsAuth:    false,
		},
		nil
}

func LogOutUser(db *sql.DB, login, password string) error {
	_, err := GetUserByLoginAndPasword(db, login, password)
	if err != nil {
		return err
	}

	query := `
		UPDATE users 
		SET is_auth=$1, last_unauth=$2, updated_at=$3
		WHERE login=$4 AND password=$5
	`
	timeNow := time.Now()

	_, err = db.Exec(
		query,
		false,
		timeNow,
		timeNow,
		login,
		password,
	)

	return err
}

func GetUserByLoginAndPasword(db *sql.DB, login, password string) (User, error) {
	query := `
		SELECT id, is_auth FROM users
		WHERE login=$1 AND password=$2
	`
	var user User = User{}

	err := db.QueryRow(query, login, password).Scan(
		&user.Id,
		&user.IsAuth,
	)

	return user, err
}

func Authorize(db *sql.DB, login, password string) error {
	user, err := GetUserByLoginAndPasword(db, login, password)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET is_auth=$1, last_auth=$2, updated_at=$3
		WHERE id=$4
	`
	timeNow := time.Now()

	result, err := db.Exec(
		query,
		true,
		timeNow,
		timeNow,
		user.Id,
	)

	if rows, _ := result.RowsAffected(); rows != 1 {
		return errors.New("user not found")
	}

	return err
}
