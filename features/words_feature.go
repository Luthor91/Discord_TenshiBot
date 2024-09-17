package features

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Déclaration des variables pour stocker les mots
var goodwords []string
var badwords []string

// LoadWords charge les badwords et goodwords depuis leurs fichiers JSON
func LoadWords() {
	// Charger les mots interdits
	data, err := os.ReadFile("resources/badwords.json")
	if err != nil {
		log.Fatalf("Erreur lors du chargement des badwords: %v", err)
	}
	var result map[string][]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("Erreur lors du parsing des badwords: %v", err)
	}
	badwords = result["badwords"]

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

// AddGoodWord ajoute un mot à la liste des "goodwords" et met à jour le fichier JSON
func AddGoodWord(word string) {
	goodwords = append(goodwords, word)
	fmt.Println(word)
	// Sauvegarder les mots dans le fichier
	err := saveWords("resources/goodwords.json", goodwords, "goodwords")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde du bon mot '%s': %v", word, err)
	} else {
		log.Printf("Le bon mot '%s' a été ajouté avec succès.", word)
	}
}

// AddBadWord ajoute un mot à la liste des "badwords" et met à jour le fichier JSON
func AddBadWord(word string) {
	badwords = append(badwords, word)

	// Sauvegarder les mots dans le fichier
	err := saveWords("resources/badwords.json", badwords, "badwords")
	if err != nil {
		log.Printf("Erreur lors de la sauvegarde du mot interdit '%s': %v", word, err)
	} else {
		log.Printf("Le mot interdit '%s' a été ajouté avec succès.", word)
	}
}

// saveWords est une fonction utilitaire pour sauvegarder les mots dans un fichier JSON
func saveWords(filename string, words []string, key string) error {
	data := map[string][]string{
		key: words,
	}
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors du marshal des données: %w", err)
	}
	err = os.WriteFile(filename, fileData, 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier %s: %w", filename, err)
	}
	return nil
}
