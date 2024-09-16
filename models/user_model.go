package models

// User représente les informations sur la monnaie d'un utilisateur
type User struct {
	Username string `json:"username"`
	Affinity int    `json:"affinity"`
}
