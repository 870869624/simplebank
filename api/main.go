package api

import (
	"os"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// 测试配置用的
func NewTestServer(t *testing.T, store db.Store) *Server {
	config := &util.Config{
		TokenKey:           util.RandomString(32),
		AcessTokenDuration: time.Minute,
	}
	Server, err := NewServer(*config, store)
	require.NoError(t, err)
	return Server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
