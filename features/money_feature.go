package features

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Luthor91/Tenshi/models"
)

var userMoneyMap map[string]models.UserMoney
var mu_money sync.Mutex

// LoadMoney charge les informations de monnaie depuis le fichier JSON
func LoadMoney() {
	data, err := os.ReadFile("resources/money.json")
	if err != nil {
		log.Printf("Erreur lors du chargement de money.json: %v", err)
		userMoneyMap = make(map[string]models.UserMoney) // Initialiser la map si fichier inexistant
		return
	}
	err = json.Unmarshal(data, &userMoneyMap)
	if err != nil {
		log.Fatalf("Erreur lors du parsing de money.json: %v", err)
	}
}

// SaveMoney sauvegarde la monnaie des utilisateurs dans money.json
func SaveMoney() {
	mu_money.Lock()
	defer mu_money.Unlock()

	data, err := json.MarshalIndent(userMoneyMap, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde des monnaies: %v", err)
		return
	}
	err = os.WriteFile("../resources/money.json", data, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture du fichier money.json: %v", err)
	}
}

// AddMoney ajoute de la monnaie à un utilisateur donné
func AddMoney(userID string, amount int) {
	mu_money.Lock()
	defer mu_money.Unlock()

	user, exists := userMoneyMap[userID]
	if !exists {
		user = models.UserMoney{
			UserID: userID,
			Money:  0,
		}
	}
	user.Money += amount
	userMoneyMap[userID] = user

	SaveMoney()
}

// GetUserMoney renvoie la quantité de monnaie d'un utilisateur
func GetUserMoney(userID string) int {
	mu_money.Lock()
	defer mu_money.Unlock()

	user, exists := userMoneyMap[userID]
	if !exists {
		return 0
	}
	return user.Money
}

// CanReceiveDailyReward vérifie si l'utilisateur peut recevoir une récompense quotidienne
func CanReceiveDailyReward(userID string) (bool, time.Duration) {
	mu_money.Lock()
	defer mu_money.Unlock()

	user, exists := userMoneyMap[userID]
	if !exists {
		return true, 0
	}

	now := time.Now()
	if now.Sub(user.LastDailyReward).Hours() >= 24 {
		return true, 0
	}
	return false, time.Until(user.LastDailyReward.Add(24 * time.Hour))
}

// GiveDailyMoney accorde la récompense quotidienne et met à jour la dernière date de réception
func GiveDailyMoney(userID string, amount int) {
	mu_money.Lock()
	defer mu_money.Unlock()

	user, exists := userMoneyMap[userID]
	if !exists {
		user = models.UserMoney{
			UserID: userID,
			Money:  0,
		}
	}

	user.Money += amount
	user.LastDailyReward = time.Now()
	userMoneyMap[userID] = user

	SaveMoney()
}
