package gui

import (
	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
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

	_, err = iv.State.Session.ChannelMessageSend(iv.State.CurrentChannel.ID, v.Buffer())
	v.Clear()
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	return err
}

func (iv *InputView) bindKeys() error {
	if err := iv.gui.SetKeybinding(inputView, gocui.KeyEnter, gocui.ModNone, iv.sendMessage); err != nil {
		return err
	}
	return nil
}
