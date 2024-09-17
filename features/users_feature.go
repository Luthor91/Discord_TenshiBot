package features

import (
	"encoding/json"
	"log"
	"os"
	"sort"

	"github.com/Luthor91/Tenshi/models"
)

var usersMap map[string]models.User

func initializeUsersMap() {
	usersMap = make(map[string]models.User)
}

// LoadUsers charge les informations sur les utilisateurs depuis le fichier JSON
func LoadUsers() {
	data, err := os.ReadFile("resources/users.json")
	if err != nil {
		log.Printf("Erreur lors du chargement de users.json: %v", err)
		initializeUsersMap()
		return
	}

	if len(data) == 0 {
		log.Println("users.json est vide, création d'un fichier valide.")
		initializeUsersMap()
		return
	}

	usersMapMu.Lock()         // Verrouiller le mutex
	defer usersMapMu.Unlock() // Déverrouiller le mutex lorsque la fonction retourne

	if err := json.Unmarshal(data, &usersMap); err != nil {
		log.Fatalf("Erreur lors du parsing de users.json: %v", err)
	}
}

// SaveUsers sauvegarde les informations sur les utilisateurs dans le fichier JSON
func SaveUsers() {
	data, err := json.MarshalIndent(usersMap, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde de users.json: %v", err)
		return
	}
	if err := os.WriteFile("resources/users.json", data, 0644); err != nil {
		log.Printf("Erreur lors de l'écriture de users.json: %v", err)
	}
}

// AddUserIfNotExists ajoute un utilisateur s'il n'existe pas déjà dans usersMap
func AddUserIfNotExists(userID, username string) {
	if _, exists := usersMap[userID]; !exists {
		usersMap[userID] = models.User{
			UserID:          userID,
			Username:        username,
			Affinity:        0,
			Money:           0,
			Experience:      0,
			LastDailyReward: "",
			RankMoney:       0,
			RankExperience:  0,
			RankAffinity:    0,
		}
		SaveUsers()
	}
}

// GetUserByID renvoie les données d'un utilisateur par son ID
func GetUserByID(userID string) (models.User, bool) {
	user, exists := usersMap[userID]
	return user, exists
}

// UpdateUser met à jour les informations d'un utilisateur
func UpdateUser(userID string, user models.User) bool {
	if _, exists := usersMap[userID]; exists {
		usersMap[userID] = user
		SaveUsers()
		return true
	}
	return false
}

// GetAllUsers renvoie tous les utilisateurs triés par score combiné décroissant
func GetAllUsers() []models.User {
	return sortUsersByScore(usersMap)
}

// UpdateUserRanks met à jour les rangs des utilisateurs pour différentes catégories
func UpdateUserRanks() {
	updateUserRanks("Money", func(u models.User) int { return u.Money })
	updateUserRanks("Experience", func(u models.User) int { return u.Experience })
	updateUserRanks("Affinity", func(u models.User) int { return u.Affinity })
	updateUserRanks("General", func(u models.User) int {
		return u.Money + u.Experience + u.Affinity
	})
}

// updateUserRanks trie les utilisateurs et met à jour leurs rangs pour une catégorie spécifique
func updateUserRanks(category string, scoreFunc func(models.User) int) {
	users := sortUsersByScore(usersMap, scoreFunc)
	for rank, user := range users {
		switch category {
		case "Money":
			user.RankMoney = rank + 1
		case "Experience":
			user.RankExperience = rank + 1
		case "Affinity":
			user.RankAffinity = rank + 1
		case "General":
			user.Rank = rank + 1
		}
		usersMap[user.UserID] = user
	}
	SaveUsers()
}

// sortUsersByScore trie les utilisateurs par score calculé à l'aide de scoreFunc
func sortUsersByScore(usersMap map[string]models.User, scoreFunc ...func(models.User) int) []models.User {
	var users []models.User
	for _, user := range usersMap {
		users = append(users, user)
	}
	if len(scoreFunc) > 0 {
		sort.Slice(users, func(i, j int) bool {
			return scoreFunc[0](users[i]) > scoreFunc[0](users[j])
		})
	}
	return users
}

// GetUserRankAndScore renvoie le rang et le score d'un utilisateur
func GetUserRankAndScore(userID string) (int, int) {
	users := GetAllUsers()
	for rank, user := range users {
		if user.UserID == userID {
			return rank + 1, user.Money + user.Experience + user.Affinity
		}
	}
	return 0, 0
}

// GetAllUsersByCategory renvoie tous les utilisateurs triés par une catégorie spécifique
func GetAllUsersByCategory(category string) []models.User {
	var scoreFunc func(models.User) int
	switch category {
	case "money":
		scoreFunc = func(u models.User) int { return u.Money }
	case "affinity":
		scoreFunc = func(u models.User) int { return u.Affinity }
	case "xp":
		scoreFunc = func(u models.User) int { return u.Experience }
	case "general":
		scoreFunc = func(u models.User) int {
			return u.Money + u.Experience + u.Affinity
		}
	default:
		return []models.User{}
	}
	return sortUsersByScore(usersMap, scoreFunc)
}

// GetUserRankAndScoreByCategory renvoie le rang et le score d'un utilisateur pour une catégorie spécifique
func GetUserRankAndScoreByCategory(userID, category string) (int, int, bool) {
	users := GetAllUsersByCategory(category)
	for rank, user := range users {
		if user.UserID == userID {
			var score int
			switch category {
			case "money":
				score = user.Money
			case "affinity":
				score = user.Affinity
			case "xp":
				score = user.Experience
			case "general":
				score = user.Money + user.Experience + user.Affinity
			}
			return rank + 1, score, true
		}
	}
	return 0, 0, false
}

// GetUserScore renvoie le score d'un utilisateur en fonction de la catégorie demandée
func GetUserScore(user models.User, category string) int {
	switch category {
	case "money":
		return user.Money
	case "affinity":
		return user.Affinity
	case "xp":
		return user.Experience
	case "general":
		return user.Money + user.Experience + user.Affinity
	default:
		return 0
	}
}

func GetUsersByMoney() []models.User {
	return GetAllUsersByCategory("money")
}

func GetUsersByAffinity() []models.User {
	return GetAllUsersByCategory("affinity")
}

func GetUsersByXP() []models.User {
	return GetAllUsersByCategory("xp")
}

func GetUsersByGeneral() []models.User {
	return GetAllUsersByCategory("general")
}
