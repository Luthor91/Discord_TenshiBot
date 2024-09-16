package features

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/Luthor91/Tenshi/models"
)

var users map[string]models.User
var mu_users sync.Mutex

// LoadUsers charge les informations sur les utilisateurs depuis users.json
func LoadUsers() {
	mu_users.Lock()
	defer mu_users.Unlock()

	data, err := os.ReadFile("../resources/users.json")
	if err != nil {
		log.Printf("Erreur lors du chargement des utilisateurs, création d'un nouveau fichier : %v", err)
		users = make(map[string]models.User)
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des utilisateurs: %v", err)
	}
}

// SaveUsers sauvegarde les utilisateurs dans le fichier users.json
func SaveUsers() {
	mu_users.Lock()
	defer mu_users.Unlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde des utilisateurs: %v", err)
		return
	}
	err = os.WriteFile("../resources/users.json", data, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture des utilisateurs dans le fichier: %v", err)
	}
}

// AddUserIfNotExists ajoute un utilisateur s'il n'existe pas déjà dans users.json
func AddUserIfNotExists(userID string, username string) {
	mu_users.Lock()
	defer mu_users.Unlock()

	if _, exists := users[userID]; !exists {
		users[userID] = models.User{
			Username: username,
			Affinity: 0,
		}
		SaveUsers()
	}
}
