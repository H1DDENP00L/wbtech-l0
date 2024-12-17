package main

import (
	"database/sql"
	"l0/db"
	"log"
)

func main() {
	pg, err := sql.Open(db.DbDriver, db.DbSource)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	err = db.DatabaseDown(pg)
	if err != nil {
		log.Fatal("Ошибка запуска миграций", err)
	}
}
