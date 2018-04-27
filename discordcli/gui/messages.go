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
	State *core.State
}

func (mv *MessagesView) render() error {
	maxX, maxY := mv.State.Gui.Size()
	if v, err := mv.State.Gui.SetView(messagesView, 21, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func (mv *MessagesView) showMessages(messages []*discordgo.Message) error {
	mv.State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(messagesView)
		if err != nil {
			return err
		}
		for i := len(messages) - 1; i >= 0; i-- {
			fmt.Fprintln(v, mv.formatMessage(messages[i]))
		}
		return nil
	})
	return nil
}

func (mv *MessagesView) clearBuffer() error {
	mv.State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(messagesView)
		if err != nil {
			return err
		}
		v.Clear()
		return nil
	})
	return nil
}

func (mv *MessagesView) formatMessage(message *discordgo.Message) string {
	return message.Author.Mention() + " " + mv.formatAuthorNickname(message) + ": " + emoji.Sprintf(message.ContentWithMentionsReplaced())
}

func (mv *MessagesView) formatAuthorNickname(message *discordgo.Message) string {
	author, err := mv.State.Session.State.Member(mv.State.CurrentGuild.ID, message.Author.ID)

	if err == nil {
		memberNick := author.Nick
		if memberNick == "" {
			memberNick = message.Author.Username
		}
		mv.State.GuildMap[mv.State.CurrentGuild.ID].Members[memberNick] = author
		return core.GetColoredNick(memberNick, mv.State.Session.State.UserColor(message.Author.ID, message.ChannelID))
	}
	return core.GetColoredNick(message.Author.Username, mv.State.Session.State.UserColor(message.Author.ID, message.ChannelID))

}

func (mv *MessagesView) bindKeys() error {
	return nil
}
