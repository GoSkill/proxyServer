package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	//драйвер для MYSQL(импорт только так"_")
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	InitializeRoutes()

	s := &http.Server{
		Addr:           ":8080", //для второй реплики :8081
		Handler:        InitializeRoutes(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if len(os.Args[:]) > 1 {
		s.Addr = os.Args[1]
	}
	log.Printf("Сервер слушает на 127.0.0.1%s", s.Addr)
	s.ListenAndServe()

}

// Функция openDB() подключает MySQL из Go и возвращает пул соединений
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn) //sql.Open инициализирует пул для будущего использования
	if err != nil {
		log.Printf("Ошибка базы данных" + err.Error())
		return nil, err
	}
	//	defer db.Close() не включать !!!
	if err = db.Ping(); err != nil { //db.Ping метод для создания соединения и проверки на наличие ошибок
		log.Printf("Ошибка при подключении к DB в следующем: " + err.Error())
		return nil, err
	}
	return db, nil
}
