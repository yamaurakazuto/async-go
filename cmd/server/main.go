//main.goの役割
//アプリケーションの起動
//外部リソースの初期化
//依存関係の組み立て

package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:Yamakazu510@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("mysqlへ接続できません")
	}

	log.Println("接続成功")
}
