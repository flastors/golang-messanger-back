package app

import (
	"errors"
	"log"

	"github.com/pressly/goose/v3"
)

// Функция автоматически накатывает миграции, для отката меняю goose.Up на goose.Down и перезапускаю приложение :)
// Либо можно использовать goose CLI инструмент, в .env прописано всё необходимое
func RunPostgresMigrations(dsn string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := goose.Up(db, "./database/postgresql/migrations"); err != nil {
		if errors.Is(err, goose.ErrAlreadyApplied) {
			log.Println("Migrations already up-to-date")
			return nil
		}
		return err
	}
	log.Println("Migrations applied successfully")
	return nil
}
