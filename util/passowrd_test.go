package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	str := RandomString(6)
	hashPassword, err := HashPassword(str)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	err = CheckPassword(str, hashPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
