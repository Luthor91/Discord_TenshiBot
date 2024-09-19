package migration

import (
	"log"

	"github.com/Luthor91/Tenshi/models"

	"gorm.io/gorm"
)

// MigrateAllPostgresql migre toutes les tables de la base de données
func MigrateAllPostgresql(db *gorm.DB) {
	// Migrate the schema
	err := db.AutoMigrate(&models.User{}, &models.Log{}, &models.ShopItem{}, &models.UserShopCooldown{}, &models.Item{}, &models.GoodWord{}, &models.BadWord{})
	if err != nil {
		log.Fatalf("Error migrating the database schema: %v", err)
	}
	SeedShopItems(db)

}

// SeedShopItems insère des articles dans la table des items du shop
func SeedShopItems(db *gorm.DB) {
	items := []models.ShopItem{
		{Name: "50 XP", Price: 100, Cooldown: 3600, Emoji: "1️⃣"}, // 1 heure en secondes
		{Name: "500 XP", Price: 1000, Cooldown: 3600, Emoji: "2️⃣"},
		{Name: "XP", Price: 100, Cooldown: 3600, Emoji: "3️⃣"},            // Prix sera calculé
		{Name: "Timeout", Price: 5000, Cooldown: 3600 * 10, Emoji: "4️⃣"}, // 5 minutes en secondes
	}

	for _, item := range items {
		// Vérifie si l'item existe déjà pour éviter les doublons
		var existingItem models.ShopItem
		if err := db.Where("name = ?", item.Name).First(&existingItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Si l'item n'existe pas, l'insérer
				if err := db.Create(&item).Error; err != nil {
					log.Println("Erreur lors de l'insertion de l'article dans la base de données:", err)
				}
			} else {
				log.Println("Erreur lors de la vérification de l'article:", err)
			}
		}
	}
}
