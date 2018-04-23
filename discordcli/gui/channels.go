package gui

import (
	"fmt"
	"io"

	"github.com/fatih/color"

	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
)

/*
ChannelsView shows list of channels in current server
*/
type ChannelsView struct {
	gui   *gocui.Gui
	State *core.State
}

const (
	categoryMark     = "-"
	hashtagMark      = "#"
	treeSignUpMiddle = "├"
	treeSignUpEnding = "└"
)

func (cv *ChannelsView) render() error {
	_, maxY := cv.gui.Size()
	if v, err := cv.gui.SetView(channelsView, 0, 0, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if _, err := cv.gui.SetCurrentView(channelsView); err != nil {
			return err
		}
	}
	return nil
}

func (cv *ChannelsView) drawGuilds(guilds []*discordgo.Guild) error {
	cv.gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(channelsView)
		if err != nil {
			return err
		}
		for _, guild := range guilds {
			drawGuild(v, guild)
		}
		return nil
	})
	return nil
}

func drawGuild(w io.Writer, guild *discordgo.Guild) error {
	fmt.Fprintln(w, guild.Name)
	drawChannels(w, guild.Channels)
	fmt.Fprintln(w)
	return nil
}

func drawChannels(w io.Writer, channels []*discordgo.Channel) error {
	for _, category := range core.SortChannels(channels) {
		fmt.Fprintln(w, formatChannel(category.Channel), color.BlueString(category.Channel.Name))
		for _, channel := range category.Channels {
			fmt.Fprintln(w, " ", formatChannel(channel), color.BlueString(channel.Name))
		}

	}

	return nil
}

func formatChannel(channel *discordgo.Channel) string {
	switch channel.Type {
	case discordgo.ChannelTypeGuildVoice:
		return color.RedString("v")
	case discordgo.ChannelTypeGuildCategory:
		return color.YellowString("-")
	case discordgo.ChannelTypeGuildText:
		return color.GreenString("#")
	}
	return ""
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cv *ChannelsView) bindKeys() error {
	if err := cv.gui.SetKeybinding(channelsView, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := cv.gui.SetKeybinding(channelsView, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	return nil
}
