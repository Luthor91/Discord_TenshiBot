package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Luthor91/DiscordBot/commands"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
)

func Run() {
	// Créer une nouvelle session Discord
	discord, err := discordgo.New("Bot " + config.AppConfig.BotToken)
	utils.CheckErr(err)

	discord.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions

	// Enregistrer les gestionnaires d'événements
	commands.RegisterHandlers(discord)

	// Ouvrir la connexion à Discord
	err = discord.Open()
	utils.CheckErr(err)

	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("Erreur lors de la fermeture de la connexion: %v", err)
		}
	}()

	// Garder le bot en fonctionnement jusqu'à une interruption système (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
