package gui

import (
	"regexp"
	"strings"

	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
)

const (
	usernameRegex = "@([^\\s]+)"
)

/*
InputView shows message input
*/
type InputView struct {
	gui   *gocui.Gui
	State *core.State
}

func (iv *InputView) render() error {
	maxX, maxY := iv.gui.Size()
	if v, err := iv.gui.SetView(inputView, 21, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = false
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = " message "
		v.Wrap = true
		if _, err := iv.gui.SetCurrentView(inputView); err != nil {
			return err
		}
	}
	return nil
}

func (iv *InputView) sendMessage(g *gocui.Gui, v *gocui.View) error {
	v, err := g.View(inputView)
	if err != nil {
		return err
	}
	msg := v.Buffer()
	v.Clear()
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	if iv.State.CommandManager.IsCommand(msg) {
		err := iv.State.CommandManager.ExecuteCommand(msg)
		//iv.gui.
		return err
	}
	_, err = iv.State.Session.ChannelMessageSend(iv.State.CurrentChannel.ID, iv.formatNewMessage(msg))

	return err
}

func (iv *InputView) formatNewMessage(message string) string {
	members := iv.State.CurrentGuild.Members
	mentions := regexp.MustCompile(usernameRegex).FindAllString(message, -1)

	replaceMessage := message
	for _, mention := range mentions {
		for _, member := range members {
			memberNick := member.Nick
			if memberNick == "" {
				memberNick = member.User.Username
			}

			if strings.TrimPrefix(mention, "@") == memberNick {
				replaceMessage = strings.Replace(replaceMessage, mention, member.User.Mention(), -1)
				break
			}
		}
	}
	return replaceMessage
}

func (iv *InputView) bindKeys() error {
	if err := iv.gui.SetKeybinding(inputView, gocui.KeyEnter, gocui.ModNone, iv.sendMessage); err != nil {
		return err
	}
	return nil
}
