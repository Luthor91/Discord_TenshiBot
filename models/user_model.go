package models

// User représente un utilisateur avec sa monnaie, expérience et affinité
type User struct {
	UserID          string `json:"userid"`
	Username        string `json:"username"`
	Affinity        int    `json:"affinity"`
	Money           int    `json:"money"`
	Experience      int    `json:"experience"`
	LastDailyReward string `json:"lastdailyreward"`
	Rank            int    `json:"rank"`
	RankMoney       int    `json:"rank_money"`
	RankExperience  int    `json:"rank_experience"`
	RankAffinity    int    `json:"rank_affinity"`
}

var usersMap map[string]User
