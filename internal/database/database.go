package database

import (
	"database/sql"
	"event-registration/internal/config"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Init(cfg config.DBConfig) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Ошибка подключения к БД: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ Ошибка пинга БД: %v", err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL!")

	runMigrations(db)

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Ошибка инициализации драйвера миграций: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("❌ Ошибка создания экземпляра миграций: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Ошибка применения миграций: %v", err)
	}
	fmt.Println("✅ Миграции успешно применены!")
}