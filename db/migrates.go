package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"         // PostgreSQL драйвер
	"github.com/pressly/goose/v3" // Goose пакет
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:zxc@localhost:5432/wbtech?sslmode=disable"
)

func DatabaseDown() {
	// Подключение к базе данных
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()
	if err := goose.Down(db, "./migrates"); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}
	log.Println("Таблицам пиздец!")
}

func DatabaseUp() {
	// Подключение к базе данных
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Применение миграций
	if err := goose.Up(db, "./migrates"); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
}
