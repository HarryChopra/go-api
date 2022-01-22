package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/harrychopra/go-api/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomName()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomName(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	payload, err := NewPayload(util.RandomName(), time.Minute)
	require.NoError(t, err)

	// Create a jwt token with a different signing method to ours
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)

	// Signing token with None besides testing purposes is forbidden.
	// Only allowed, If UnsafeAllowNoneSignatureType constant is used as the signed string
	badToken, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(badToken)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
