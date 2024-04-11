package api

import (
	db "simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)

// 处理网络请求,包含数据库处理，以及路由处理
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// 创建一个服务实例，并且完成所有路由设置
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts/", server.listAccount)
	server.router = router
	return server
}

// 开启监听地址服务
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// 返回错误状态吗对应的文本信息
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
