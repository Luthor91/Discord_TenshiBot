package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features" // Importer le package pour gérer l'affinité
	"github.com/bwmarrin/discordgo"
)

func checkErr(e error) {
	if e != nil {
		log.Fatalf("Erreur: %v", e)
	}
}

func Run() {
	// Créer une nouvelle session Discord
	discord, err := discordgo.New("Bot " + config.AppConfig.BotToken)
	checkErr(err)

	discord.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions

	//here
	RegisterHandlers(discord)
	err = discord.Open()
	//checkErr(err)

	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("Erreur lors de la fermeture de la connexion: %v", err)
		}
	}()

	// Charger les mots, la monnaie, et les utilisateurs
	features.LoadWords()
	features.LoadUsers()

	// Garder le bot en fonctionnement jusqu'à une interruption système (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
