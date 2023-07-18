package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
}

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	createUser := `
		INSERT INTO users (
			username,
			hashed_password,
			full_name,
			email
		) VALUES (
			$1, $2, $3, $4
		) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at, is_email_verified;
	`

	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
	)

	return i, err
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	getUser := `
		SELECT username, hashed_password, full_name, email, password_changed_at, created_at, is_email_verified
		FROM users
		WHERE username = $1 LIMIT 1;
	`

	row := q.db.QueryRow(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
	)

	return i, err
}

type UpdateUserParams struct {
	HashedPassword    pgtype.Text        `json:"hashed_password"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	FullName          pgtype.Text        `json:"full_name"`
	Email             pgtype.Text        `json:"email"`
	IsEmailVerified   pgtype.Bool        `json:"is_email_verified"`
	Username          pgtype.Text        `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	updateUser := `
		UPDATE users
		SET
			hashed_password = COALESCE($1, hashed_password),
			password_changed_at = COALESCE($2, password_changed_at),
			full_name = COALESCE($3, full_name),
			email = COALESCE($4, email),
			is_email_verified = COALESCE($5, is_email_verified),
		WHERE
			username = $6
		RETURNING username, hashed_password, full_name, email, password_changed_at, created_at, is_email_verified;
	`

	row := q.db.QueryRow(ctx, updateUser, arg.HashedPassword, arg.PasswordChangedAt, arg.FullName, arg.Email, arg.IsEmailVerified, arg.Username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailVerified,
	)

	return i, err
}
