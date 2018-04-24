package core

import "github.com/bwmarrin/discordgo"

const (
	GuildMapTypeGuild     = 1
	GuildMapTypeChannel   = 2
	GuildMapTypeSeparator = 3
)

type GuildMapItem struct {
	Type    int
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
}
