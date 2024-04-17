package db

import (
	"context"
	"fmt"
	"log"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	User, err := testQueries.CreateUser(context.Background(), arg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(User)
	require.NoError(t, err)
	require.NotEmpty(t, User)

	require.Equal(t, arg.Username, User.Username)
	require.Equal(t, arg.HashedPassword, User.HashedPassword)
	require.Equal(t, arg.FullName, User.FullName)
	require.Equal(t, arg.Email, User.Email)

	require.True(t, User.PasswordChangedAt.IsZero())
	require.NotZero(t, User.CreatedAt)
	return User
}

func TestCrteateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
