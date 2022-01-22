package db

import (
	"context"
	"testing"
	"time"

	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T, arg *CreateUserParams) User {
	if arg == nil {
		arg = &CreateUserParams{
			Username: util.RandomName(),
			FullName: util.RandomName(),
			Email:    util.RandomEmail(),
		}
	}
	var err error
	arg.HashedPassword, err = util.HashedPassword(util.RandomString(8))
	require.NotEmpty(t, arg.HashedPassword)
	require.NoError(t, err)
	user, err := testQueries.CreateUser(context.Background(), *arg)
	require.NoError(t, err)
	return user
}
func TestCreateUser(t *testing.T) {
	arg := &CreateUserParams{
		Username: util.RandomName(),
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}
	user := createRandomUser(t, arg)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
}

func TestGetUser(t *testing.T) {
	userA := createRandomUser(t, nil)
	userB, err := testQueries.GetUser(context.Background(), userA.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userB)
	require.Equal(t, userA.Username, userA.Username)
	require.Equal(t, userA.HashedPassword, userA.HashedPassword)
	require.Equal(t, userA.FullName, userA.FullName)
	require.Equal(t, userA.Email, userA.Email)
	require.WithinDuration(t, userA.CreatedAt, userB.CreatedAt, time.Second)
}
