package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// 处理网络请求,包含数据库处理，以及路由处理
type Server struct {
	config     util.Config
	store      db.Store
	tokenmaker token.Maker
	router     *gin.Engine
}

// 创建一个服务实例，并且完成所有路由设置
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("不能制造令牌")
	}
	server := &Server{store: store, config: config, tokenmaker: tokenMaker} //server是自己自定义的，里面包含方法存储器和路由
	//和gin.New的区别就是不用自己再加上gin.recover

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	} //获取gin当前使用的验证器引擎,binding是gin的一个子包
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRouters := router.Group("/", authMiddleware(server.tokenmaker))

	authRouters.POST("/accounts", server.createAccount)
	authRouters.GET("/accounts/:id", server.getAccount)
	authRouters.GET("/accounts/", server.listAccount)

	authRouters.POST("/transfers", server.createTransfer)
	server.router = router
}

// 开启监听地址服务
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// 返回错误状态吗对应的文本信息
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
