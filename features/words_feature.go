package features

import (
	"encoding/json"
	"log"
	"os"
)

// Déclaration des variables pour stocker les mots
var goodwords []string
var banwords []string

// LoadWords charge les banwords et goodwords depuis leurs fichiers JSON
func LoadWords() {
	// Charger les mots interdits
	data, err := os.ReadFile("resources/badwords.json")
	if err != nil {
		log.Fatalf("Erreur lors du chargement des banwords: %v", err)
	}
	var result map[string][]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des banwords: %v", err)
	}
	banwords = result["banwords"]

	// Charger les bons mots
	data, err = os.ReadFile("resources/goodwords.json")
	if err != nil {
		log.Fatalf("Erreur lors du chargement des goodwords: %v", err)
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des goodwords: %v", err)
	}
	goodwords = result["goodwords"]
}
