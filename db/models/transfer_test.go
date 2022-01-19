package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t, nil)
	account2 := createRandomAccount(t, nil)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.WithinDuration(t, account1.CreatedAt, transfer.CreatedAt, time.Second)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t, nil)
	account2 := createRandomAccount(t, nil)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	transferFound, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transferFound)
	require.Equal(t, arg.FromAccountID, transferFound.FromAccountID)
	require.Equal(t, arg.ToAccountID, transferFound.ToAccountID)
	require.Equal(t, arg.Amount, transferFound.Amount)
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t, nil)
	account2 := createRandomAccount(t, nil)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        0,
	}
	for i := 0; i < 15; i++ {
		arg.Amount++
		_, _ = testQueries.CreateTransfer(context.Background(), arg)
	}

	arg2 := ListTransfersParams{
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Limit:         10,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 10)
	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, arg2.FromAccountID, transfer.FromAccountID)
		require.Equal(t, arg2.ToAccountID, transfer.ToAccountID)
	}
	require.Equal(t, int64(10), transfers[4].Amount)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createRandomAccount(t, nil)
	account2 := createRandomAccount(t, nil)

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)

	err = testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	deletedTransfer, err := testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Empty(t, deletedTransfer)
	require.ErrorIs(t, sql.ErrNoRows, err)
}
