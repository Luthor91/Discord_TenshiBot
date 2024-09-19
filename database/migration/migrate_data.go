package migration

import (
	"log"

	"github.com/Luthor91/Tenshi/models"

	"gorm.io/gorm"
)

// MigrateAllPostgresql migre toutes les tables de la base de données
func MigrateAllPostgresql(db *gorm.DB) {
	// Migrate the schema
	err := db.AutoMigrate(&models.User{}, &models.Item{}, &models.Log{}, &models.ShopItem{}, &models.UserShopCooldown{}, &models.GoodWord{}, &models.BadWord{})
	if err != nil {
		log.Fatalf("Error migrating the database schema: %v", err)
	}
}
