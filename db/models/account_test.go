package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T, arg *CreateAccountParams) Account {
	if arg == nil {
		arg = &CreateAccountParams{
			Owner:    util.RandomName(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}
	}
	account, err := testQueries.CreateAccount(context.Background(), *arg)
	require.NoError(t, err)
	return account
}

func TestCreateAccount(t *testing.T) {
	arg := &CreateAccountParams{
		Owner:    util.RandomName(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account := createRandomAccount(t, arg)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	accountA := createRandomAccount(t, nil)
	accountB, err := testQueries.GetAccount(context.Background(), accountA.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountB)
	require.Equal(t, accountA.ID, accountB.ID)
	require.Equal(t, accountA.Owner, accountB.Owner)
	require.Equal(t, accountA.Balance, accountB.Balance)
	require.Equal(t, accountA.Currency, accountB.Currency)
	require.WithinDuration(t, accountA.CreatedAt, accountB.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	arg1 := CreateAccountParams{
		Owner:    "TestTest",
		Balance:  0,
		Currency: "GBP",
	}
	for i := 0; i < 15; i++ {
		arg1.Balance++
		createRandomAccount(t, &arg1)
	}
	arg2 := ListAccountsParams{
		Owner:  arg1.Owner,
		Limit:  10,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 10)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, arg2.Owner, account.Owner)
	}
	require.Equal(t, int64(10), accounts[4].Balance)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t, nil)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)
	require.Equal(t, arg.ID, updatedAccount.ID)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t, nil)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	deletedAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	require.ErrorIs(t, sql.ErrNoRows, err)
	require.Empty(t, deletedAccount)
}
