package lol_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/api/riot_games"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// RotationData représente la structure des données de rotation des champions
type RotationData struct {
	FreeChampionIds              []float64 `json:"freeChampionIds"`
	FreeChampionIdsForNewPlayers []float64 `json:"freeChampionIdsForNewPlayers"`
}

// ChampionRotationCommand récupère la rotation des champions
func ChampionRotationCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the command matches
	command := fmt.Sprintf("%slolrotation", config.AppConfig.BotPrefix)
	if m.Content == command {
		// Appeler GetChampionRotations pour obtenir les données de rotation des champions
		rotationDataInterface, err := riot_games.GetChampionRotations()
		if err != nil {
			log.Println("Error fetching free champion rotation:", err)
			s.ChannelMessageSend(m.ChannelID, "Error fetching free champion rotation.")
			return
		}

		// Convertir rotationDataInterface en RotationData
		rotationData, ok := rotationDataInterface.(map[string]interface{})
		if !ok {
			log.Println("Error: failed to assert type of rotationData")
			s.ChannelMessageSend(m.ChannelID, "Error processing champion rotation data.")
			return
		}

		// Extraire les IDs des champions gratuits
		var freeChampionIds []float64
		if ids, ok := rotationData["freeChampionIds"].([]interface{}); ok {
			for _, id := range ids {
				if idFloat, ok := id.(float64); ok {
					freeChampionIds = append(freeChampionIds, idFloat)
				}
			}
		}

		// Extraire les IDs des champions gratuits pour les nouveaux joueurs
		var freeChampionIdsNewPlayers []float64
		if ids, ok := rotationData["freeChampionIdsForNewPlayers"].([]interface{}); ok {
			for _, id := range ids {
				if idFloat, ok := id.(float64); ok {
					freeChampionIdsNewPlayers = append(freeChampionIdsNewPlayers, idFloat)
				}
			}
		}

		// Convertir les IDs des champions en noms
		freeChampionNames, err := riot_games.GetChampionsNameByIds(freeChampionIds)
		if err != nil {
			log.Println("Error getting champion names:", err)
			s.ChannelMessageSend(m.ChannelID, "Error getting champion names.")
			return
		}

		// Convertir les IDs des champions gratuits pour les nouveaux joueurs en noms
		freeChampionNamesNewPlayers, err := riot_games.GetChampionsNameByIds(freeChampionIdsNewPlayers)
		if err != nil {
			log.Println("Error getting new player champion names:", err)
			s.ChannelMessageSend(m.ChannelID, "Error getting new player champion names.")
			return
		}

		// Format the response with better structure and line breaks
		response := fmt.Sprintf(
			"**Current Free Champion Rotations**:\n- %s\n\n**Free Champions for New Players**:\n- %s",
			formatChampionList(freeChampionNames),
			formatChampionList(freeChampionNamesNewPlayers),
		)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// Helper function to format champion list with line breaks, ignoring empty names
func formatChampionList(champions []string) string {
	var filteredChampions []string
	for _, champion := range champions {
		if champion != "" { // Ignore empty champion names
			filteredChampions = append(filteredChampions, champion)
		}
	}
	return strings.Join(filteredChampions, ", ")
}
