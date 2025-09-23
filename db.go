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
	result, err := r.db.Exec("INSERT INTO users (username,first_name, last_name, date_of_birth, email, password, created_at) VALUES (?,?,?,?,?,?,?)", u.Username, u.FirstName, u.LastName, u.DateOfBirth, u.Email, u.Password, u.CreatedAt)
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
	row := r.db.QueryRow("SELECT id, username, first_name, last_name, date_of_birth, email, password, balance, created_at FROM users WHERE id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.DateOfBirth, &u.Email, &u.Password, &u.Balance, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
