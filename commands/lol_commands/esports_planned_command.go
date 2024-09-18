package lol_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/api/riot_games"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// EsportsMatchesCommand récupère les informations des prochains matchs eSports
func EsportsMatchesPlannedCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si la commande correspond
	command := fmt.Sprintf("%sesportsmatches", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		// Appeler GetUpcomingEsportMatches pour obtenir les données des matchs eSports
		matchesData, err := riot_games.GetUpcomingEsportMatches()
		if err != nil {
			log.Println("Error fetching upcoming esports matches:", err)
			s.ChannelMessageSend(m.ChannelID, "Error fetching upcoming esports matches.")
			return
		}

		// Vérifier que nous avons des données
		if len(matchesData) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No upcoming esports matches found.")
			return
		}

		// Format the response with better structure
		response := "**Upcoming Esports Matches**:\n"
		for _, match := range matchesData {
			response += fmt.Sprintf(
				"**Match:** %s vs %s\n**Start Time:** %s\n\n",
				match.Team1, match.Team2, match.StartTime,
			)
		}
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
