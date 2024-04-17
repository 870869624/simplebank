package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalln("无法加载配置", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalln("无法连接", err)
	}
	fmt.Println(testDB)
	testQueries = New(testDB)
	os.Exit(m.Run())
}
