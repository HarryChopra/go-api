package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t, nil)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    -50,
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.WithinDuration(t, account.CreatedAt, entry.CreatedAt, time.Second)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t, nil)
	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	})
	require.NoError(t, err)

	entryFound, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryFound)
	require.Equal(t, entry.ID, entryFound.ID)
	require.Equal(t, entry.AccountID, entryFound.AccountID)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t, nil)
	arg1 := CreateEntryParams{
		AccountID: account.ID,
		Amount:    0,
	}
	for i := 0; i < 15; i++ {
		arg1.Amount++
		_, _ = testQueries.CreateEntry(context.Background(), arg1)
	}
	arg2 := ListEntriesParams{
		AccountID: account.ID,
		Limit:     10,
		Offset:    5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Len(t, entries, 10)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg2.AccountID, entry.AccountID)
	}
	require.Equal(t, int64(10), entries[4].Amount)
}

func TestDeleteEntry(t *testing.T) {
	account := createRandomAccount(t, nil)
	entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	})
	require.NoError(t, err)
	err = testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	deletedEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.Empty(t, deletedEntry)
	require.ErrorIs(t, sql.ErrNoRows, err)
}
