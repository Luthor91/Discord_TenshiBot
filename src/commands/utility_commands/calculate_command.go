package utility_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
	"github.com/expr-lang/expr"
)

// CalculateCommand évalue une expression mathématique donnée
func CalculateCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%scalculate", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Extraire l'expression après la commande
		expression := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if expression == "" {
			s.ChannelMessageSend(m.ChannelID, "Merci de spécifier une expression mathématique à évaluer.")
			return
		}

		// Évaluer l'expression
		result, err := evaluateExpression(expression)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'évaluation de l'expression.")
			log.Println("Erreur lors de l'évaluation de l'expression :", err)
			return
		}

		// Envoyer le résultat de l'évaluation
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Résultat de l'expression `%s` : %v", expression, result))
	}
}

// evaluateExpression évalue une expression mathématique simple
func evaluateExpression(expression string) (interface{}, error) {
	program, err := expr.Compile(expression, expr.Env(map[string]interface{}{}))
	if err != nil {
		return nil, err
	}
	result, err := expr.Run(program, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
