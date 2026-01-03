// main.goの役割
// アプリケーションの起動
// 外部リソースの初期化
// 依存関係の組み立て

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"async-go/internal/handler"
	"async-go/internal/repository"
	"async-go/internal/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/banking?parseTime=true"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		log.Fatal("mysqlへ接続できません")
	}

	repo := repository.NewBankRepository(db)
	service := service.NewBankService(repo)
	handler := handler.NewBankHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("server started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
