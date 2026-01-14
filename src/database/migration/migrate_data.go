package migration

import (
	"log"

	"github.com/Luthor91/DiscordBot/models"

	"gorm.io/gorm"
)

// MigrateAllPostgresql migre toutes les tables de la base de données
func MigrateAllPostgresql(db *gorm.DB) {
	// Migrate the schema
	err := db.AutoMigrate(
		&models.User{},
		&models.Log{},
		&models.BadWord{},
		&models.Warn{},
	)
	if err != nil {
		log.Fatalf("Error migrating the database schema: %v", err)
	}

}

