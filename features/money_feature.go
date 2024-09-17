package features

import (
	"sync"
	"time"

	"github.com/Luthor91/Tenshi/models"
)

var (
	usersMapMu sync.Mutex
)

// AddMoney ajoute de la monnaie à un utilisateur donné
func AddMoney(userID string, amount int) {
	user, exists := usersMap[userID]
	if !exists {
		user = models.User{
			UserID:     userID,
			Money:      0,
			Affinity:   0,
			Experience: 0,
		}
	}
	user.Money += amount
	usersMap[userID] = user
	SaveUsers()
}

// GetUserMoney renvoie la quantité de monnaie d'un utilisateur
func GetUserMoney(userID string) int {
	user, exists := usersMap[userID]
	if !exists {
		return 0
	}
	return user.Money
}

// CanReceiveDailyReward vérifie si l'utilisateur peut recevoir une récompense quotidienne
func CanReceiveDailyReward(userID string) (bool, time.Duration) {
	user, exists := usersMap[userID]
	if !exists {
		return true, 0
	}

	// Convertir la chaîne en time.Time
	lastRewardTime, err := time.Parse(time.RFC3339, user.LastDailyReward)
	if err != nil {
		// En cas d'erreur de parsing, assumer que la récompense peut être reçue
		return true, 0
	}

	now := time.Now()
	if now.Sub(lastRewardTime).Hours() >= 24 {
		return true, 0
	}
	return false, time.Until(lastRewardTime.Add(24 * time.Hour))
}

// GiveDailyMoney accorde la récompense quotidienne et met à jour la dernière date de réception
func GiveDailyMoney(userID string, amount int) {
	usersMapMu.Lock()         // Verrouiller le mutex
	defer usersMapMu.Unlock() // Déverrouiller le mutex lorsque la fonction retourne

	user, exists := usersMap[userID]
	if !exists {
		user = models.User{
			UserID:     userID,
			Money:      0,
			Affinity:   0,
			Experience: 0,
		}
	}

	user.Money += amount
	user.LastDailyReward = time.Now().Format(time.RFC3339) // Convertir en string
	usersMap[userID] = user
	SaveUsers()
}
