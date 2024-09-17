package features

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	cooldownFile     = "resources/shop_cooldown.json"
	cooldownDuration = time.Hour // Durée du cooldown (1 heure)
)

// Structure pour stocker le cooldown des utilisateurs pour chaque option
type UserCooldowns struct {
	Option1 time.Time `json:"option1"`
	Option2 time.Time `json:"option2"`
	Option3 time.Time `json:"option3"`
}

// Fonction pour définir le cooldown du shop pour un utilisateur et une option spécifique
func SetShopCooldown(userID string, option string, lastUsed time.Time) error {
	cooldowns, err := LoadCooldowns()
	if err != nil {
		return err
	}

	// Assurez-vous que l'utilisateur a une entrée dans les cooldowns
	if _, exists := cooldowns[userID]; !exists {
		cooldowns[userID] = UserCooldowns{}
	}

	// Récupérer les cooldowns existants pour cet utilisateur
	userCooldowns := cooldowns[userID]

	// Mettre à jour le cooldown approprié
	switch option {
	case "option1":
		userCooldowns.Option1 = lastUsed
	case "option2":
		userCooldowns.Option2 = lastUsed
	case "option3":
		userCooldowns.Option3 = lastUsed
	default:
		return fmt.Errorf("option inconnue : %s", option)
	}

	// Réassigner les cooldowns modifiés dans la carte
	cooldowns[userID] = userCooldowns

	// Sauvegarder les cooldowns modifiés dans le fichier
	return SaveCooldowns(cooldowns)
}

// Fonction pour vérifier si le cooldown est expiré pour une option spécifique
func IsCooldownExpired(userID string, option string) (bool, error) {
	cooldowns, err := LoadCooldowns()
	if err != nil {
		return false, err
	}

	// Vérifier si l'utilisateur a des cooldowns enregistrés
	userCooldowns, exists := cooldowns[userID]
	if !exists {
		return true, nil // Pas de données de cooldown, donc pas de cooldown
	}

	var lastUsed time.Time
	switch option {
	case "option1":
		lastUsed = userCooldowns.Option1
	case "option2":
		lastUsed = userCooldowns.Option2
	case "option3":
		lastUsed = userCooldowns.Option3
	default:
		return false, fmt.Errorf("option inconnue : %s", option)
	}

	return time.Since(lastUsed) > cooldownDuration, nil
}

// Fonction pour obtenir le cooldown d'une option spécifique pour un utilisateur
func GetShopCooldown(userID string, option string) (time.Time, error) {
	cooldowns, err := LoadCooldowns()
	if err != nil {
		return time.Time{}, err
	}

	// Vérifier si l'utilisateur a des cooldowns enregistrés
	userCooldowns, exists := cooldowns[userID]
	if !exists {
		return time.Time{}, fmt.Errorf("aucun cooldown trouvé pour l'utilisateur : %s", userID)
	}

	var cooldownTime time.Time
	switch option {
	case "option1":
		cooldownTime = userCooldowns.Option1
	case "option2":
		cooldownTime = userCooldowns.Option2
	case "option3":
		cooldownTime = userCooldowns.Option3
	default:
		return time.Time{}, fmt.Errorf("option inconnue : %s", option)
	}

	return cooldownTime, nil
}

// Fonction pour charger les cooldowns depuis le fichier JSON
func LoadCooldowns() (map[string]UserCooldowns, error) {
	// Vérifie si le fichier existe
	if _, err := os.Stat(cooldownFile); os.IsNotExist(err) {
		// Le fichier n'existe pas, donc on le crée
		file, err := os.Create(cooldownFile)
		if err != nil {
			return nil, fmt.Errorf("impossible de créer le fichier de cooldowns: %v", err)
		}
		defer file.Close()

		// Initialise un map vide
		emptyCooldowns := make(map[string]UserCooldowns)

		// Écrire le map vide dans le fichier
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(emptyCooldowns); err != nil {
			return nil, fmt.Errorf("erreur lors de l'écriture dans le fichier: %v", err)
		}

		return emptyCooldowns, nil
	}

	// Ouvre le fichier s'il existe
	file, err := os.Open(cooldownFile)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'ouverture du fichier: %v", err)
	}
	defer file.Close()

	// Décode les cooldowns à partir du fichier
	var cooldowns map[string]UserCooldowns
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cooldowns); err != nil {
		return nil, fmt.Errorf("erreur lors du décodage du fichier: %v", err)
	}

	return cooldowns, nil
}

// Fonction pour sauvegarder les cooldowns dans le fichier JSON
func SaveCooldowns(cooldowns map[string]UserCooldowns) error {
	file, err := os.Create(cooldownFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Formatage lisible
	if err := encoder.Encode(cooldowns); err != nil {
		return err
	}

	return nil
}
