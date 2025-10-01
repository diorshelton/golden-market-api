package models

// import (
// 	"database/sql"

// 	"github.com/google/uuid"
// )

// type UserRepo interface {
// 	CreateUser(u *User) error
// 	GetUser(id int) (*User, error)
// }

// type SQLiteUserRepo struct {
// 	db *sql.DB
// }

// func (r *SQLiteUserRepo) CreateUser(u *User) error {
// 	user := &User{
// 		ID:    uuid.New(),
// 		Email: email,
// 	}
// 	query := `
// 	INSERT INTO users (username, first_name, last_name, email, password_hash, balance,created_at) VALUES (?,?,?,?,?,?,?)
// 	`
// 	result, err := r.db.Exec(query, u.Username, u.FirstName, u.LastName, u.Email, u.PasswordHash, u.Balance, u.CreatedAt)
// 	if err != nil {
// 		return err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}

// 	u.ID = id
// 	return nil
// }

// func (r *SQLiteUserRepo) GetUser(id int) (*User, error) {
// 	row := r.db.QueryRow("SELECT id, username, first_name, last_name, email, password_hash, balance, created_at FROM users WHERE id = ?", id)
// 	var u User
// 	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.Email, &u.PasswordHash, &u.Balance, &u.CreatedAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &u, nil
// }

// func (r *SQLiteUserRepo) GetAllUsers() ([]User, error) {
// 	rows, err := r.db.Query(`SELECT id, username, first_name, last_name, email, password_hash, balance, created_at FROM users`)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		var u User
// 		var createdAt sql.NullTime

// 		err = rows.Scan(
// 			&u.ID,
// 			&u.Username,
// 			&u.FirstName,
// 			&u.LastName,
// 			&u.Email,
// 			&u.PasswordHash,
// 			&u.Balance,
// 			&u.CreatedAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if createdAt.Valid {
// 			u.CreatedAt = createdAt.Time
// 		}

// 		users = append(users, u)

// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }
