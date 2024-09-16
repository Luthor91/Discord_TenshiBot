package models

// UserRank représente le rang d'un utilisateur dans le classement
type UserRank struct {
	UserID string `json:"user_id"`
	Rank   int    `json:"rank"`
	Money  int    `json:"money"`
}
