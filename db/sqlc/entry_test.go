package db

import (
	"context"
	"simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, account Account) {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	account1, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account1)

	require.Equal(t, arg.AccountID, account1.AccountID)
	require.Equal(t, arg.Amount, account1.Amount)

	require.NotZero(t, account1.ID)
	require.NotZero(t, account1.CreatedAt)
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestListEntre(t *testing.T) {
	//先生成插入，再将生成随机偏移量，然后再去查找
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	account2, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, account2, 5)

	for _, v := range account2 {
		require.NotEmpty(t, v)
	}
}
