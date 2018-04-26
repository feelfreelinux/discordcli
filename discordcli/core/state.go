package core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jroimartin/gocui"
)

/*
State holds reference of current application state (selected channel, etc)
*/
type State struct {
	CurrentGuild   *discordgo.Guild
	CurrentChannel *discordgo.Channel
	Gui            *gocui.Gui
	Session        *discordgo.Session
	GuildMap       map[string]*GuildMapItem
	CommandManager *CommandManager
}
