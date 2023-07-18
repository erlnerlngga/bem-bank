package db

import (
	"context"
	"time"
)

type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateTransferParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	createTransfer := `
		INSERT INTO transfers(
			from_account_id,
			to_account_id,
			amount
		) VALUES (
			$1, $2, $3
		) RETURNING id, from_account_id, to_account_id, amount, created_at;
	`

	row := q.db.QueryRow(ctx, createTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)

	return i, err
}

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	getTransfer := `
		SELECT id, from_account_id, to_account_id, amount, created_at
		FROM transfers
		WHERE id = $1 LIMIT 1;
	`

	row := q.db.QueryRow(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)

	return i, err
}

type ListTransfersParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `jsoN:"offset"`
}

func (q *Queries) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]Transfer, error) {
	listTransfers := `
		SELECT id, from_account_id, to_account_id, amount, created_at
		FROM transfers
		WHERE from_account_id $1 OR to_account_id = $2
		ORDER BY id
		LIMIT $3
		OFFSET $4; 
	`

	rows, err := q.db.Query(ctx, listTransfers, arg.FromAccountID, arg.ToAccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Transfer{}
	for rows.Next() {
		var i Transfer
		err := rows.Scan(
			&i.ID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Amount,
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
