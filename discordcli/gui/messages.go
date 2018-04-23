package gui

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
	"gopkg.in/kyokomi/emoji.v1"
)

/*
MessagesView shows messages in current selected channel
*/
type MessagesView struct {
	gui   *gocui.Gui
	State *core.State
}

func (mv *MessagesView) render() error {
	maxX, maxY := mv.gui.Size()
	if v, err := mv.gui.SetView(messagesView, 21, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
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
		for i := len(messages) - 1; i >= 0; i-- {
			fmt.Fprintln(v, formatMessage(mv.State.Session, messages[i]))
		}
		return nil
	})
	return nil
}

func formatMessage(session *discordgo.Session, message *discordgo.Message) string {
	return formatAuthorNickname(session, message) + ": " + emoji.Sprintf(message.ContentWithMentionsReplaced())
}

func formatAuthorNickname(session *discordgo.Session, message *discordgo.Message) string {
	return core.GetColoredNick(message.Author.Username, session.State.UserColor(message.Author.ID, message.ChannelID))
}

func (mv *MessagesView) bindKeys() error {
	return nil
}
