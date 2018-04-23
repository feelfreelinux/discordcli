package core

import "github.com/bwmarrin/discordgo"

/*
State holds reference of current application state (selected channel, etc)
*/
type State struct {
	CurrentGuild   *discordgo.Guild
	CurrentChannel *discordgo.Channel
	Session        *discordgo.Session
}
