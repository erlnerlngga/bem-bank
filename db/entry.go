package db

import (
	"context"
	"time"
)

type Entry struct {
	ID        int64     `json:"id"`
	AccountID int64     `json:"account_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg CreateEntryParams) (Entry, error) {
	createEntry := `
		INSERT INTO entries(
			account_id,
			amount
		) VALUES (
			$1, $2
		) RETURNING id, account_id, amount, created_at;
	`

	row := q.db.QueryRow(ctx, createEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)

	return i, err
}

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	getEntry := `
		SELECT id, account_id, amount, created_at
		FROM entries
		WHERE id = $1 LIMIT 1;
	`

	row := q.db.QueryRow(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)

	return i, err
}

type ListEntryParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntryParams) ([]Entry, error) {
	listEntries := `
		SELECT id, account_id, amount, created_at
		FROM entries
		WHERE account_id = $1
		ORDER by id
		LIMIT $2
		OFFSET $3;
	`

	rows, err := q.db.Query(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Entry{}
	for rows.Next() {
		var i Entry
		err := rows.Scan(
			&i.ID,
			&i.AccountID,
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
