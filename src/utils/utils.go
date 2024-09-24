package utils

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func PrintMemoryUsage(intervalSeconds int) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Conversion des valeurs en MB pour plus de clarté
		fmt.Printf("Mémoire allouée: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("Mémoire totale allouée: %.2f MB\n", float64(memStats.TotalAlloc)/1024/1024)
		fmt.Printf("Mémoire système obtenue: %.2f MB\n", float64(memStats.Sys)/1024/1024)
		fmt.Printf("Nombre de Garbage Collection : %d\n", memStats.NumGC)
		fmt.Println("-----------------------------------")
	}
}

// parseDuration parse la durée au format '10s', '5m', '2h', '1d'
func ParseDuration(durationStr string) (time.Duration, error) {
	if len(durationStr) < 2 {
		return 0, fmt.Errorf("durée invalide")
	}

	unit := durationStr[len(durationStr)-1]
	valueStr := durationStr[:len(durationStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("valeur numérique invalide")
	}

	switch unit {
	case 's': // secondes
		return time.Duration(value) * time.Second, nil
	case 'm': // minutes
		return time.Duration(value) * time.Minute, nil
	case 'h': // heures
		return time.Duration(value) * time.Hour, nil
	case 'd': // jours
		return time.Duration(value) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("unité de temps invalide")
	}
}

func CheckErr(e error) {
	if e != nil {
		log.Fatalf("Erreur: %v", e)
	}
}

// Réponse d'erreur standard pour l'envoi de messages
func SendErrorMessage(s *discordgo.Session, channelID, errMessage string) {
	s.ChannelMessageSend(channelID, fmt.Sprintf("Erreur : %s", errMessage))
}
