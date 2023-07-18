package db

import (
	"context"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountParams struct {
	Owner    string `json:"owner"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	createAccount := `
		INSERT INTO accounts (
			owner,
			balance,
			currency
		) VALUES (
			$1, $2, $3
		) RETURNING id, owner, balance, currency, created_at;
	`

	row := q.db.QueryRow(ctx, createAccount, arg.Owner, arg.Balance, arg.Currency)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	addAccountBalance := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING id, owner, balance, currency, created_at;
	`

	row := q.db.QueryRow(ctx, addAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	deleteAccount := `
		DELETE FROM accounts
		WHERE id = $1;
	`

	_, err := q.db.Exec(ctx, deleteAccount)

	return err
}

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	getAccount := `
		SELECT id, owner, balance, currency, created_ad
		FROM accounts
		WHERE id = $1 LIMMIT 1;	
	`
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	getAccountForUpdate := `
		SELECT id, owner, balance, currency, created_at
		FROM accounts
		WHERE id = $1 LIMIT 1
		FOR NO KEY UPDATE;
	`
	row := q.db.QueryRow(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}

type ListAccountsParams struct {
	Owner  string `json:"owner"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	listAccounts := `
		SELECT id, owner, balance, currency, created_at
		FROM accounts
		WHERE owner = $1
		ORDER BY id
		LIMIT $2
		OFFSET $3;
	`

	rows, err := q.db.Query(ctx, listAccounts, arg.Owner, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Account{}
	for rows.Next() {
		var i Account
		err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type UpdateAccountParams struct {
	ID      int64 `json:"id"`
	Balance int64 `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	updateAccount := `
		UPDATE accounts
		SET balance = $2
		WHERE id = $1
		RETURNING id, owner, balance, currency, created_at;	
	`
	row := q.db.QueryRow(ctx, updateAccount, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
	)

	return i, err
}
