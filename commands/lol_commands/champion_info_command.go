package lol_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/api/riot_games"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// ChampionInfoCommand récupère les informations sur un champion spécifique à partir de son nom
func ChampionInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the command matches
	command := fmt.Sprintf("%slolchamp", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		// Extraire le nom du champion de la commande
		args := strings.Split(m.Content, " ")
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier le nom du champion.")
			return
		}

		championName := strings.Join(args[1:], " ")

		// Appeler GetChampionData pour obtenir les données du champion
		championData, err := riot_games.GetChampionData()
		if err != nil {
			log.Println("Erreur lors de la récupération des données des champions:", err)
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des données des champions.")
			return
		}

		// Trouver les informations du champion spécifié par son nom
		championInfo := findChampionInfoByName(championData, championName)
		if championInfo == nil {
			s.ChannelMessageSend(m.ChannelID, "Champion non trouvé.")
			return
		}

		// Format the response with champion details
		response := formatChampionInfo(championInfo)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// Fonction helper pour trouver les informations du champion par son nom
func findChampionInfoByName(championData map[string]riot_games.ChampionDataExtended, name string) *riot_games.ChampionDataExtended {
	for _, champion := range championData {
		if strings.EqualFold(champion.Name, name) {
			return &champion
		}
	}
	return nil
}

// Fonction helper pour formater les informations du champion pour la réponse
func formatChampionInfo(championInfo *riot_games.ChampionDataExtended) string {
	name := championInfo.Name
	title := championInfo.Title
	blurb := championInfo.Blurb

	// Ajouter d'autres détails selon ce que vous souhaitez afficher
	return fmt.Sprintf("**%s** (%s)\n\n%s", name, title, blurb)
}
