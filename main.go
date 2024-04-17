package main

import (
	"database/sql"
	"fmt"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".") //记载配置 。表示就在当前目录下
	if err != nil {
		log.Fatal("不能载入配置")
	}
	fmt.Println(config.DBDriver)
	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("无法连接", err)
	}

	store := db.NewStore(connection) //初始化
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("不能开启服务", err)
	}
}
