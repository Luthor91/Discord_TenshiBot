package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Luthor91/Tenshi/features" // Importer le package pour gérer l'affinité
	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkErr(e error) {
	if e != nil {
		log.Fatalf("Erreur: %v", e)
	}
}

func Run() {
	// Créer une nouvelle session Discord
	discord, err := discordgo.New("Bot " + BotToken)
	checkErr(err)
	// Configurer les intents nécessaires
	discord.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessageReactions
	// Ouvrir la connexion
	//here
	err = discord.Open()
	checkErr(err)

	defer func() {
		if err := discord.Close(); err != nil {
			log.Printf("Erreur lors de la fermeture de la connexion: %v", err)
		}
	}()

	// Charger les mots, la monnaie, et les utilisateurs
	features.LoadWords()
	features.LoadMoney()
	features.LoadUsers()

	// Récupérer les guildes et les membres, et les ajouter dans users.json
	guilds, err := discord.UserGuilds(100, "", "", false)
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des guildes: %v", err)
	}

	for _, guild := range guilds {
		members, err := discord.GuildMembers(guild.ID, "", 1000)
		if err != nil {
			log.Printf("Erreur lors de la récupération des membres pour la guilde %s: %v", guild.ID, err)
			continue
		}
		for _, member := range members {
			features.AddUserIfNotExists(member.User.ID, member.User.Username)
		}
	}

	// Ajouter un gestionnaire d'événements pour les messages
	discord.AddHandler(newMessage)

	// Garder le bot en fonctionnement jusqu'à une interruption système (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Éviter que le bot réponde à ses propres messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Ajouter de la monnaie à l'utilisateur
	features.AddMoney(message.Author.ID, 1)
	features.AddExperience(message.Author.ID, message.Author.Username, 1)

	// Appeler la fonction pour ajuster l'affinité en fonction des mots utilisés
	features.AdjustAffinity(discord, message)
}
