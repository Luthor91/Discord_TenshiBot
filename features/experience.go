package features

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Experience struct {
	Username   string `json:"username"`
	Experience int    `json:"experience"`
}

var experienceData map[string]Experience
var mu_exp sync.Mutex

// LoadExperience charge les informations sur l'expérience depuis experience.json
func LoadExperience() {
	mu_exp.Lock()
	defer mu_exp.Unlock()

	data, err := os.ReadFile("../resources/experience.json")
	if err != nil {
		log.Printf("Erreur lors du chargement des données d'expérience, création d'un nouveau fichier : %v", err)
		experienceData = make(map[string]Experience)
		return
	}
	err = json.Unmarshal(data, &experienceData)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des données d'expérience: %v", err)
	}
}

// SaveExperience sauvegarde les données d'expérience dans le fichier experience.json
func SaveExperience() {
	mu_exp.Lock()
	defer mu_exp.Unlock()

	data, err := json.MarshalIndent(experienceData, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde des données d'expérience: %v", err)
		return
	}
	err = os.WriteFile("resources/experience.json", data, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture des données d'expérience dans le fichier: %v", err)
	}
}

// AddExperience ajoute de l'expérience à un utilisateur
func AddExperience(userID string, username string, amount int) {
	mu_exp.Lock()
	defer mu_exp.Unlock()

	if user, exists := experienceData[userID]; exists {
		user.Experience += amount
		experienceData[userID] = user
	} else {
		experienceData[userID] = Experience{
			Username:   username,
			Experience: amount,
		}
	}
	SaveExperience()
}

// GetExperience renvoie l'expérience d'un utilisateur
func GetExperience(userID string) (int, bool) {
	mu_exp.Lock()
	defer mu_exp.Unlock()

	if user, exists := experienceData[userID]; exists {
		return user.Experience, true
	}
	return 0, false
}
