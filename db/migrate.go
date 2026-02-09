package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer m.Close()

	versionBefore, _, errBefore := m.Version()
	if errBefore != nil && errBefore != migrate.ErrNilVersion {
		versionBefore = 0
	}

	migrationErr := m.Up()
	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", migrationErr)
	}

	versionAfter, _, _ := m.Version()

	if migrationErr == migrate.ErrNoChange || (versionBefore > 0 && versionBefore == versionAfter) {
		log.Println("Migrations: already up to date")
	} else {
		log.Println("Migrations: applied successfully")
	}

	return nil
}
