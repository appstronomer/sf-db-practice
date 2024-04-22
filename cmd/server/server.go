package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
	"os"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	var db storage.Interface
	var err error

	switch os.Getenv("GONEWS_DB") {
	case "postgres":
		// Реляционная БД PostgreSQL.
		db, err = postgres.New("postgres://user:password@postgres:5432/sf")
	case "mongo":
		// Документная БД MongoDB.
		db, err = mongo.New("mongodb://mongo:27017")
	default:
		// БД в памяти.
		db = memdb.New()
	}

	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":80", srv.api.Router())
}
