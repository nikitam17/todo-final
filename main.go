package main

import (
	"fmt"
	"net/http"
	"os"

	"todo-final/pkg/api"
	"todo-final/pkg/db"

	"github.com/joho/godotenv"
)

// Открытие/создание БД и запуск веб-сервера
func main() {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	// Устанавливаем директорию для web
	webDir := "./web"
	// Устанавливаем имя БД
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	db, err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	// Устанавливаем имя Порт
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	fmt.Println("Запускаем сервер")
	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
