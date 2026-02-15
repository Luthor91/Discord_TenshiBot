package services

import (
	"github.com/Luthor91/DiscordBot/controllers"
	"github.com/Luthor91/DiscordBot/models"
	"github.com/bwmarrin/discordgo"
)

type ReactionRoleService struct {
	controller     *controllers.ReactionRoleController
	discordSession *discordgo.Session
}

func NewReactionRoleService(session *discordgo.Session) *ReactionRoleService {
	return &ReactionRoleService{
		controller:     controllers.NewReactionRoleController(),
		discordSession: session,
	}
}

// Création des mappings
func (s *ReactionRoleService) CreateReactionRoles(guildID, channelID, messageID string, mappings map[string]string) error {

	for emoji, roleID := range mappings {

		rr := models.ReactionRole{
			GuildID:   guildID,
			ChannelID: channelID,
			MessageID: messageID,
			Emoji:     emoji,
			RoleID:    roleID,
		}

		if err := s.controller.Create(&rr); err != nil {
			return err
		}
	}

	return nil
}

// Gestion ajout réaction
func (s *ReactionRoleService) HandleReactionAdd(r *discordgo.MessageReactionAdd) {

	if r.UserID == s.discordSession.State.User.ID {
		return
	}

	dbRoles, err := s.controller.GetByMessageID(r.MessageID)
	if err != nil || len(dbRoles) == 0 {
		return
	}

	emoji := r.Emoji.Name
	if r.Emoji.ID != "" {
		emoji = r.Emoji.APIName()
	}

	for _, rr := range dbRoles {
		if rr.Emoji == emoji {
			_ = s.discordSession.GuildMemberRoleAdd(r.GuildID, r.UserID, rr.RoleID)
			return
		}
	}
}

// Gestion suppression réaction
func (s *ReactionRoleService) HandleReactionRemove(r *discordgo.MessageReactionRemove) {

	dbRoles, err := s.controller.GetByMessageID(r.MessageID)
	if err != nil || len(dbRoles) == 0 {
		return
	}

	emoji := r.Emoji.Name
	if r.Emoji.ID != "" {
		emoji = r.Emoji.APIName()
	}

	for _, rr := range dbRoles {
		if rr.Emoji == emoji {
			_ = s.discordSession.GuildMemberRoleRemove(r.GuildID, r.UserID, rr.RoleID)
			return
		}
	}
}
