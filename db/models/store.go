package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// Begin the transaction
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// For each query within the transaction
	q := New(tx)

	// If the transaction fails, rollback the transaction
	if err = fn(q); err != nil {
		if rbErr := tx.Rollback(); err != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("tx error: %v", err)
	}
	return tx.Commit()
}

// Input for transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult struct contains result of each operation in the transaction
type TransferTxResult struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTX performs all money transfer operations within the transfer transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// a. A transfer record
		if result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg)); err != nil {
			return err
		}

		// b. Entry (from) record
		if result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}); err != nil {
			return err
		}

		// c. Entry (to) record
		if result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}); err != nil {
			return err
		}

		// d. & e. Account Balance update
		// Before concurrent row updates, Order each row operation by ID to prevent deadlock
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = updateBalance(
				ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = updateBalance(
				ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})
	return result, err
}

// updateBalance performs the adjustment of balance amount for two accounts
func updateBalance(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (
	account1, account2 Account, err error) {
	if account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount1,
		ID:     accountID1,
	}); err != nil {
		return
	}
	if account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		Amount: amount2,
		ID:     accountID2,
	}); err != nil {
		return
	}
	return
}
