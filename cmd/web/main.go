package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	_"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //драйвер для MYSQL(импорт только так"_")
)

var DB *sql.DB

func main() {
	addr := flag.String("addr", ":8080", "Сетевой адрес HTTP")
	dsn := flag.String("dsn", "web:1@/proxyserver?parseTime=true", "Имя источника данных MySQL")
	flag.Parse() 

	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"}) //доверенный IP-адрес
	
	DB, _ = sql.Open("mysql", *dsn)//Инициализировать базу данных MYSQL
	
	InitializeRoutes()

	s := &http.Server{		
		Addr:           *addr,
		Handler:        InitializeRoutes(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Сервер слушает на 127.0.0.1%s", s.Addr)
	s.ListenAndServe()
}
