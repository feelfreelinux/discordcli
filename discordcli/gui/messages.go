package gui

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jroimartin/gocui"
)

/*
MessagesView shows messages in current selected channel
*/
type MessagesView struct {
	gui     *gocui.Gui
	session *discordgo.Session
}

func (mv *MessagesView) render() error {
	maxX, maxY := mv.gui.Size()
	if v, err := mv.gui.SetView(messagesView, 21, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true

	}
	return nil
}

func (mv *MessagesView) showMessages(messages []*discordgo.Message) error {
	mv.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(messagesView)
		if err != nil {
			return err
		}
		for _, message := range messages {
			fmt.Fprintln(v, formatMessage(message))
		}
		return nil
	})
	return nil
}

func formatMessage(message *discordgo.Message) string {
	return message.Author.Username + ": " + message.ContentWithMentionsReplaced()
}

func (mv *MessagesView) bindKeys() error {
	return nil
}
