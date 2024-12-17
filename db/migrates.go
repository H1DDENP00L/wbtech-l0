package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"         // PostgreSQL драйвер
	"github.com/pressly/goose/v3" // Goose пакет
)

const (
	DbDriver      = "postgres"
	DbSource      = "postgres://postgres:zxc@localhost:5432/wbtech?sslmode=disable"
	migrationsDir = "./migrates"
)

func DatabaseDown(db *sql.DB) error {
	if err := goose.DownTo(db, migrationsDir, 0); err != nil {
		return fmt.Errorf("Не удалось применить миграции: %w", err)
	}
	log.Println("Миграции применены успешно Все таблицы удалены!")
	return nil
}

func DatabaseUp(db *sql.DB) error {
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("Не удалось применить миграции: %w", err)
	}
	log.Println("Миграции применены успешно Все таблицы работают!")
	return nil
}
