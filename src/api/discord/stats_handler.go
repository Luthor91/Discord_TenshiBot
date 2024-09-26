package discord

import (
	"fmt"
	"runtime"

	"github.com/bwmarrin/discordgo"
)

// PrintUserStats affiche des statistiques sur un utilisateur
func PrintUserStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des statistiques de l'utilisateur.")
		return
	}

	response := fmt.Sprintf("**Statistiques de l'utilisateur :**\n- Nom: %s\n- ID: %s\n- Rôles: %d",
		member.User.Username, member.User.ID, len(member.Roles))
	s.ChannelMessageSend(m.ChannelID, response)
}

// PrintServerStats affiche des statistiques sur le serveur
func PrintServerStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des statistiques du serveur.")
		return
	}

	response := fmt.Sprintf(
		"**Statistiques du serveur :**\n- Nom: %s\n- ID: %s\n- Membres: %d\n- Rôles: %d\n- Canaux: %d",
		guild.Name, guild.ID, guild.MemberCount, len(guild.Roles), len(guild.Channels),
	)
	s.ChannelMessageSend(m.ChannelID, response)
}

// PrintBotStats affiche des statistiques sur le bot
func PrintBotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Compter le nombre total de membres sur tous les serveurs
	totalUsers := 0
	for _, guild := range s.State.Guilds {
		guild, err := s.Guild(guild.ID) // Récupère la guilde pour assurer la synchronisation
		if err == nil {
			totalUsers += guild.MemberCount
		}
	}

	response := fmt.Sprintf("**Statistiques du bot :**\n- Serveurs: %d\n- Utilisateurs: %d\n- Latence: %d ms",
		len(s.State.Guilds), totalUsers, s.HeartbeatLatency().Milliseconds())
	s.ChannelMessageSend(m.ChannelID, response)
}

// PrintChannelStats affiche des statistiques sur le canal
func PrintChannelStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des statistiques du canal.")
		return
	}

	response := fmt.Sprintf("**Statistiques du canal :**\n- Nom: %s\n- ID: %s\n- Type: %d\n- Position: %d",
		channel.Name, channel.ID, channel.Type, channel.Position)
	s.ChannelMessageSend(m.ChannelID, response)
}

// BotPerfsCommand affiche les performances du bot sur Discord
func PrintBotPerfs(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Récupérer les statistiques de la mémoire
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memoryAlloc := float64(memStats.Alloc) / 1024 / 1024           // Mémoire allouée en MB
	totalMemoryAlloc := float64(memStats.TotalAlloc) / 1024 / 1024 // Mémoire totale allouée
	systemMemory := float64(memStats.Sys) / 1024 / 1024            // Mémoire système
	garbageCollections := memStats.NumGC                           // Nombre de garbage collections

	// Récupérer des informations sur les utilisateurs et les serveurs
	totalServers := len(s.State.Guilds)
	totalUsers := 0
	for _, guild := range s.State.Guilds {
		totalUsers += guild.MemberCount
	}

	// Latence
	latency := s.HeartbeatLatency().Milliseconds()

	// Créer un message avec les statistiques
	response := fmt.Sprintf("**Statistiques du bot :**\n"+
		"- Mémoire allouée : %.2f MB\n"+
		"- Mémoire totale allouée : %.2f MB\n"+
		"- Mémoire système obtenue : %.2f MB\n"+
		"- Nombre de Garbage Collection : %d\n"+
		"- Serveurs : %d\n"+
		"- Utilisateurs : %d\n"+
		"- Latence : %d ms",
		memoryAlloc, totalMemoryAlloc, systemMemory, garbageCollections,
		totalServers, totalUsers, latency)

	// Envoyer les statistiques dans le canal
	s.ChannelMessageSend(m.ChannelID, response)
}
