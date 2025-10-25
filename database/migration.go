package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) error {
	sqlDb, _ := db.DB()
	driver, _ := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://database/migrations/sqlite",
		"sqlite3", driver)
	return m.Down()
}
