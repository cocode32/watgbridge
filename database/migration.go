//package database
//
//import (
//	"fmt"
//	"time"
//
//	"gorm.io/gorm"
//)
//
//type Migration struct {
//	Name string
//	Up   func(*gorm.DB) error
//}
//
//var migrations []Migration
//
//func RegisterMigration(name string, up func(*gorm.DB) error) {
//	migrations = append(migrations, Migration{Name: name, Up: up})
//}
//
//func RunMigrations(db *gorm.DB) error {
//	db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
//		name TEXT PRIMARY KEY,
//		applied_at TIMESTAMP DEFAULT NOW()
//	)`)
//
//	for _, m := range migrations {
//		var count int64
//		db.Raw("SELECT COUNT(*) FROM schema_migrations WHERE name = ?", m.Name).Scan(&count)
//		if count > 0 {
//			continue
//		}
//
//		fmt.Printf("Running migration: %s\n", m.Name)
//		if err := m.Up(db); err != nil {
//			return fmt.Errorf("migration %s failed: %w", m.Name, err)
//		}
//		db.Exec("INSERT INTO schema_migrations (name, applied_at) VALUES (?, ?)", m.Name, time.Now())
//	}
//
//	fmt.Println("All migrations complete.")
//	return nil
//}

package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	//time.Sleep(15 * time.Second)
	sqlDb, _ := db.DB()
	driver, _ := sqlite3.WithInstance(sqlDb, &sqlite3.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://database/migrations/sqlite",
		"sqlite3", driver)
	m.Up()
}
