package features

import (
	"sort"

	"github.com/Luthor91/Tenshi/models"
)

// GetUserRankAndMoney renvoie le rang et la monnaie d'un utilisateur
func GetUserRankAndMoney(userID string) (int, int, bool) {
	mu_money.Lock()
	defer mu_money.Unlock()

	users := make([]models.UserMoney, 0, len(userMoneyMap))
	for _, user := range userMoneyMap {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Money > users[j].Money
	})

	for rank, user := range users {
		if user.UserID == userID {
			return rank + 1, user.Money, true
		}
	}

	return 0, 0, false
}

// GetAllUsersMoney renvoie la liste de tous les utilisateurs avec leur monnaie
func GetAllUsersMoney() []models.UserMoney {
	mu_money.Lock()
	defer mu_money.Unlock()

	users := make([]models.UserMoney, 0, len(userMoneyMap))
	for _, user := range userMoneyMap {
		users = append(users, user)
	}
	return users
}
