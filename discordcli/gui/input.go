package gui

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jroimartin/gocui"
)

/*
InputView shows message input
*/
type InputView struct {
	gui     *gocui.Gui
	session *discordgo.Session
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

func (iv *InputView) bindKeys() error {
	return nil
}
