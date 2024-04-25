package token

import (
	"simplebank/util"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payLoad, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payLoad)
	require.NotZero(t, payLoad.ID)
	require.Equal(t, username, payLoad.Username)
	require.WithinDuration(t, issuedAt, payLoad.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payLoad.ExpiredAt, time.Second)
}

func TestExJWTMake(t *testing.T) {
	// maker, err := NewJWTMaker(util.RandomString(32))
	payLoad, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payLoad)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payLoad, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payLoad)
}
