package gui

import "github.com/bwmarrin/discordgo"

/*
State holds reference of current application state (selected channel, etc)
*/
type State struct {
	currentGuild   *discordgo.Guild
	currentChannel *discordgo.Channel
}
