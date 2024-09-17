package features

import (
	"encoding/json"
	"log"
	"os"
)

// LoadExperience charge les informations sur l'expérience depuis users.json
func LoadExperience() {
	data, err := os.ReadFile("resources/users.json")
	if err != nil {
		log.Printf("Erreur lors du chargement des données d'expérience, création d'un nouveau fichier : %v", err)
		// Si le fichier n'existe pas, créer un fichier vide
		return
	}
	err = json.Unmarshal(data, &usersMap)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des données d'expérience: %v", err)
	}
}

// AddExperience ajoute de l'expérience à un utilisateur
func AddExperience(userID string, amount int) {
	if user, exists := usersMap[userID]; exists {
		user.Experience += amount
		usersMap[userID] = user
		SaveUsers() // Utiliser SaveUsers pour sauvegarder les modifications
	} else {
		log.Printf("Utilisateur %s non trouvé lors de l'ajout d'expérience", userID)
	}
}

// GetExperience renvoie l'expérience d'un utilisateur
func GetExperience(userID string) (int, bool) {
	if user, exists := usersMap[userID]; exists {
		return user.Experience, true
	}
	return 0, false
}
