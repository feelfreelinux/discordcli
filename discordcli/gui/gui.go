package gui

import (
	"github.com/bwmarrin/discordgo"
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
	gui      *gocui.Gui
	session  *discordgo.Session
	channels *ChannelsView
	messages *MessagesView
	input    *InputView
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

func CreateMainView(dgsession *discordgo.Session, ui *gocui.Gui) error {
	ui.Cursor = true
	var mainView = &MainView{
		gui:     ui,
		session: dgsession,
		channels: &ChannelsView{
			session: dgsession,
			gui:     ui,
		},
		messages: &MessagesView{
			session: dgsession,
			gui:     ui,
		},
		input: &InputView{
			session: dgsession,
			gui:     ui,
		},
	}
	ui.SetManagerFunc(mainView.layout)
	err := mainView.bindKeys()

	mainView.setHandlers()
	return err
}

func (mv *MainView) setHandlers() error {
	mv.session.AddHandler(mv.readyHandler)
	return nil
}

// Messy function just for testing, raw index will be replaced with correct implementation
func (mv *MainView) readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	mv.channels.showChannelsForGuild(s.State.Guilds[1])
	messages, err := s.ChannelMessages(s.State.Guilds[1].Channels[28].ID, 50, "", "", "")
	if err != nil {
		panic(err)
	}
	mv.messages.showMessages(messages)
}

func (mv *MainView) bindKeys() error {
	if err := mv.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
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
