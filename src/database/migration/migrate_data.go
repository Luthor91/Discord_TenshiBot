package migration

import (
	"log"

	"github.com/Luthor91/Tenshi/models"

	"gorm.io/gorm"
)

// MigrateAllPostgresql migre toutes les tables de la base de données
func MigrateAllPostgresql(db *gorm.DB) {
	// Migrate the schema
	err := db.AutoMigrate(
		&models.User{}, &models.Log{}, &models.ShopItem{},
		&models.Item{}, &models.UserShopCooldown{}, &models.GoodWord{},
		&models.BadWord{}, &models.Warn{}, &models.Investment{},
	)
	if err != nil {
		log.Fatalf("Error migrating the database schema: %v", err)
	}
	SeedShopItems(db)

}
func SeedShopItems(db *gorm.DB) {
	items := []models.ShopItem{
		{Name: "smolpack", Price: 100, Cooldown: 3600, Emoji: "1️⃣"}, // 1 heure en secondes
		{Name: "pack", Price: 1050, Cooldown: 3600 * 5, Emoji: "2️⃣"},
		{Name: "bigpack", Price: 1100, Cooldown: 3600 * 24, Emoji: "3️⃣"},
		{Name: "timeout", Price: 5000, Cooldown: 3600 * 48, Emoji: "4️⃣"}, // 48 heures en secondes
	}

	for _, item := range items {
		// Utilise FirstOrCreate pour vérifier si l'item existe déjà ou le créer
		if err := db.Where(models.ShopItem{Name: item.Name}).FirstOrCreate(&item).Error; err != nil {
			log.Println("Erreur lors de la création ou de la récupération de l'article dans la base de données:", err)
		}
	}
}
