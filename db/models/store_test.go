package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	account1 := createRandomAccount(t, nil)
	account2 := createRandomAccount(t, nil)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	// Concurrent txs will attempt to access EXCLUSIVE LOCK ordered by row ids
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// Transfer record
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// Entry (from) record
		fromEntry := result.FromEntry
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// Entry (to) record
		toEntry := result.ToEntry
		require.NotEmpty(t, transfer)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// Account Balances
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.Equal(t, account1.Balance-(amount*(int64(i)+1)), fromAccount.Balance)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		require.Equal(t, account2.Balance+(amount*(int64(i)+1)), toAccount.Balance)
	}
}

func TestTransferTxDeadlock(t *testing.T) {
	accountA := createRandomAccount(t, nil)
	accountB := createRandomAccount(t, nil)

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountA.ID
		toAccountID := accountB.ID

		// Concurrent transactions will attempt to get "EXCLUSIVE LOCK" on each others' LOCKED rows
		if i%2 == 1 {
			fromAccountID = accountB.ID
			toAccountID = accountA.ID
		}
		go func() {
			_, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccountA, err := testStore.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)

	updatedAccountB, err := testStore.GetAccount(context.Background(), accountB.ID)
	require.NoError(t, err)

	// Equal number or credits and debits
	require.Equal(t, accountA.Balance, updatedAccountA.Balance)
	require.Equal(t, accountB.Balance, updatedAccountB.Balance)
}
