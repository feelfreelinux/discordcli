package gui

import (
	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
)

const (
	channelsView = "channels"
	inputView    = "input"
	messagesView = "messages"
)

/*
MainView holds reference for all views and renders them
*/
type MainView struct {
	channels *ChannelsView
	messages *MessagesView
	input    *InputView
	State    *core.State
}

func (mv *MainView) layout(g *gocui.Gui) error {
	if err := mv.channels.render(); err != nil {
		return err
	}

	if err := mv.input.render(); err != nil {
		return err
	}

	if err := mv.messages.render(); err != nil {
		return err
	}
	return nil
}

/*
CreateMainView creates MainView and all of its child views
*/
func CreateMainView(dgsession *discordgo.Session, ui *gocui.Gui) error {
	ui.Cursor = true
	var state = &core.State{
		Gui:            ui,
		Session:        dgsession,
		CommandManager: &core.CommandManager{},
		GuildMap:       make(map[string]*core.GuildMapItem),
	}
	var mainView = &MainView{
		State: state,
		channels: &ChannelsView{
			State: state,
		},
		messages: &MessagesView{
			State: state,
		},
		input: &InputView{
			State: state,
		},
	}
	ui.SetManagerFunc(mainView.layout)
	err := mainView.bindKeys()

	mainView.setHandlers()
	return err
}

func (mv *MainView) setHandlers() error {
	mv.State.Session.AddHandler(mv.readyHandler)
	mv.State.Session.AddHandler(mv.messageHandler)
	return nil
}

func (mv *MainView) changeChannel(channel *discordgo.Channel) {
	mv.State.CurrentChannel = channel
	guild, _ := mv.State.Session.State.Guild(channel.GuildID)
	mv.State.CurrentGuild = guild
	messages, _ := mv.State.Session.ChannelMessages(channel.ID, 50, "", "", "")
	mv.messages.clearBuffer()
	mv.messages.showMessages(messages)
}

// Messy function just for testing, raw index will be replaced with correct implementation
func (mv *MainView) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	mv.channels.channelChangedCallback = mv.changeChannel
	mv.channels.drawGuilds(s.State.Guilds)
}

func (mv *MainView) messageHandler(s *discordgo.Session, event *discordgo.MessageCreate) {
	if mv.State.CurrentChannel != nil {
		if event.ChannelID == mv.State.CurrentChannel.ID {
			mv.messages.showMessages([]*discordgo.Message{event.Message})
		}
	}
}

func (mv *MainView) bindKeys() error {
	if err := mv.State.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := mv.State.Gui.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, changeScreenFocus); err != nil {
		return err
	}

	if err := mv.channels.bindKeys(); err != nil {
		return err
	}

	if err := mv.messages.bindKeys(); err != nil {
		return err
	}

	if err := mv.input.bindKeys(); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func changeScreenFocus(g *gocui.Gui, v *gocui.View) error {
	switch g.CurrentView().Name() {
	case messagesView:
		g.SetCurrentView(inputView)
	case channelsView:
		g.SetCurrentView(messagesView)
	case inputView:
		g.SetCurrentView(channelsView)
	}
	return nil
}
