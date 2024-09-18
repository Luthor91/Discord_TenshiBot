package lol_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/api/riot_games"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// SummonerProfileCommand récupère les informations d'un profil d'invocateur
func SummonerProfileCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si la commande correspond
	command := fmt.Sprintf("%slolprofile", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		// Extraire le pseudo avec le tag de la commande
		args := strings.Split(m.Content, " ")
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Please provide the summoner name and tag.")
			return
		}

		// Le pseudo avec tag est le dernier mot
		nameTag := args[len(args)-1]
		parts := strings.SplitN(nameTag, "#", 2)
		if len(parts) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Invalid format. Use name#tag.")
			return
		}
		summoner_name := parts[0]
		tag := parts[1]

		// Appeler GetSummonerProfile pour obtenir les données du profil
		profileDataInterface, err := riot_games.GetSummonerProfile(summoner_name, tag)
		if err != nil {
			log.Println("Error fetching summoner profile:", err)
			s.ChannelMessageSend(m.ChannelID, "Error fetching summoner profile.")
			return
		}

		// Convertir profileDataInterface en ProfileData
		profileData, ok := profileDataInterface.(map[string]interface{})
		if !ok {
			log.Println("Error: failed to assert type of profileData")
			s.ChannelMessageSend(m.ChannelID, "Error processing summoner profile data.")
			return
		}

		// Extraire les informations du profil
		id, _ := profileData["id"].(string)
		accountId, _ := profileData["accountId"].(string)
		puuid, _ := profileData["puuid"].(string)
		name, _ := profileData["name"].(string)
		tagLine, _ := profileData["tagLine"].(string)
		profileIcon, _ := profileData["profileIconId"].(float64)

		// Format the response with better structure
		response := fmt.Sprintf(
			"**Summoner Profile Information**:\n"+
				"**ID:** %s\n"+
				"**Account ID:** %s\n"+
				"**PUUID:** %s\n"+
				"**Name:** %s\n"+
				"**Tag Line:** %s\n"+
				"**Profile Icon ID:** %d",
			id, accountId, puuid, name, tagLine, int(profileIcon),
		)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
