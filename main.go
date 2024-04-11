package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:123456@localhost/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8081"
)

func main() {
	connection, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalln("无法连接", err)
	}

	store := db.NewStore(connection) //初始化
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("不能开启服务", err)
	}
}
