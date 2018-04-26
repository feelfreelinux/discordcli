package core

import "github.com/bwmarrin/discordgo"

const (
	GuildMapTypeGuild = iota
	GuildMapTypeChannel
	GuildMapTypeSeparator
)

type GuildMapItem struct {
	Type    int
	Channel *discordgo.Channel
	Guild   *discordgo.Guild
	Members map[string]*discordgo.Member
}
