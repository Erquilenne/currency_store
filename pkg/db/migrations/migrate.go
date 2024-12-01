package migrations

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations выполняет миграции базы данных
func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %v", err)
	}

	// Получаем текущую версию
	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get current migration version: %v", err)
	}

	if dirty {
		log.Printf("WARNING: Database is in dirty state. Version: %d", version)
		// Опционально можно добавить автоматическое исправление dirty state
		// if err := m.Force(int(version)); err != nil {
		//     return fmt.Errorf("could not force migration version: %v", err)
		// }
	}

	// Если версия не установлена (новая база), выполняем миграции
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Println("No migrations applied yet. Running migrations...")
		if err := m.Up(); err != nil {
			return fmt.Errorf("could not run migrate up: %v", err)
		}
		log.Println("Migrations successfully executed")
		return nil
	}

	// Проверяем, есть ли новые миграции для применения
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("Database is up to date. Current version: %d", version)
			return nil
		}
		return fmt.Errorf("could not run migrate up: %v", err)
	}

	newVersion, _, err := m.Version()
	if err != nil {
		return fmt.Errorf("could not get new migration version: %v", err)
	}

	log.Printf("Migrations successfully executed. Version changed from %d to %d", version, newVersion)
	return nil
}
