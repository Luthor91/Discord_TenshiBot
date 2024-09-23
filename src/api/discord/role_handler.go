package discord

import (
	"github.com/bwmarrin/discordgo"
)

// UserHasAdminRole vérifie si l'utilisateur possède le rôle d'administrateur ou est le propriétaire du serveur
func UserHasAdminRole(s *discordgo.Session, guildID, userID string) (bool, error) {
	// Vérifie si l'utilisateur est le propriétaire du serveur
	guild, err := s.Guild(guildID)
	if err != nil {
		return false, err
	}
	if guild.OwnerID == userID {
		// L'utilisateur est le propriétaire du serveur
		return true, nil
	}

	// Récupère les rôles de l'utilisateur
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	// Récupère les rôles de la guilde directement depuis Discord
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return false, err
	}

	// Parcourt les rôles de l'utilisateur et vérifie si un rôle a la permission d'administrateur
	for _, userRoleID := range member.Roles {
		for _, role := range roles {
			if role.ID == userRoleID && role.Permissions&discordgo.PermissionAdministrator != 0 {
				// L'utilisateur a un rôle avec la permission d'administrateur
				return true, nil
			}
		}
	}

	return false, nil
}

// UserHasModeratorRole vérifie si l'utilisateur possède le rôle de modérateur ou est le propriétaire du serveur
func UserHasModeratorRole(s *discordgo.Session, guildID, userID string) (bool, error) {
	// Vérifie si l'utilisateur est le propriétaire du serveur
	guild, err := s.Guild(guildID)
	if err != nil {
		return false, err
	}
	if guild.OwnerID == userID {
		// L'utilisateur est le propriétaire du serveur
		return true, nil
	}

	// Récupère les rôles de l'utilisateur
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	// Récupère les rôles de la guilde
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err == nil && role.Permissions&discordgo.PermissionManageMessages != 0 {
			// L'utilisateur a un rôle avec la permission de gérer les messages (modérateur)
			return true, nil
		}
	}

	return false, nil
}
