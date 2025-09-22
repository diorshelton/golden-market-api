package main

import "database/sql"

type UserRepo interface {
	CreateUser(u *User) error
	GetUser(id int) (*User, error)
}

type SQLiteUserRepo struct {
	db *sql.DB
}

func (r *SQLiteUserRepo) CreateUser(u *User) error {
	result, err := r.db.Exec("INSERT INTO users (username, date_of_birth, email, password) VALUES (?,?,?,?)", u.Username, u.DateOfBirth, u.Email, u.Password)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *SQLiteUserRepo) GetUser(id int) (*User, error) {
	row := r.db.QueryRow("SELECT id, username, date_of_birth, email, password, balance FROM users WHERE id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.DateOfBirth, &u.Email, &u.Password, &u.Balance)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
