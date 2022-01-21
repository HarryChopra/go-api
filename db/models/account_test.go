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
	user := createRandomUser(t, nil)
	if arg == nil {
		arg = &CreateAccountParams{
			Owner:    user.Username,
			Balance:  util.RandomMoney(),
			Currency: util.RandomCurrency(),
		}
	}
	if arg.Owner == "" {
		arg.Owner = user.Username
	}
	account, err := testQueries.CreateAccount(context.Background(), *arg)
	require.NoError(t, err)
	return account
}

func TestCreateAccount(t *testing.T) {
	arg := &CreateAccountParams{
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
	arg1 := &CreateUserParams{
		Username:       util.RandomName(),
		HashedPassword: "secret",
		FullName:       util.RandomName(),
		Email:          util.RandomEmail(),
	}
	user := createRandomUser(t, arg1)
	accountParams := []CreateAccountParams{
		{
			Owner:    user.Username,
			Balance:  1,
			Currency: "USD",
		},
		{
			Owner:    user.Username,
			Balance:  2,
			Currency: "CAD",
		},
		{
			Owner:    user.Username,
			Balance:  3,
			Currency: "GBP",
		},
		{
			Owner:    user.Username,
			Balance:  4,
			Currency: "EUR",
		},
		{
			Owner:    user.Username,
			Balance:  5,
			Currency: "AUD",
		},
	}
	for _, param := range accountParams {
		createRandomAccount(t, &param)
	}
	arg2 := ListAccountsParams{
		Owner:  user.Username,
		Limit:  2,
		Offset: 2,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 2)
	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, arg2.Owner, account.Owner)
	}
	require.Equal(t, int64(4), accounts[1].Balance)
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
