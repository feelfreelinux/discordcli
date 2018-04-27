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
	State                  *core.State
	channelChangedCallback func(channel *discordgo.Channel)
	GuildMapPositionIDS    []string
}

const (
	categoryMark     = "-"
	hashtagMark      = "#"
	treeSignUpMiddle = "├"
	treeSignUpEnding = "└"
)

func (cv *ChannelsView) render() error {
	_, maxY := cv.State.Gui.Size()
	if v, err := cv.State.Gui.SetView(channelsView, 0, 0, 20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if _, err := cv.State.Gui.SetCurrentView(channelsView); err != nil {
			return err
		}
	}
	return nil
}

func (cv *ChannelsView) drawGuilds(guilds []*discordgo.Guild) error {
	cv.State.Gui.Update(func(g *gocui.Gui) error {
		v, err := g.View(channelsView)
		if err != nil {
			return err
		}
		cv.GuildMapPositionIDS = []string{}
		for _, guild := range guilds {
			cv.State.GuildMap[guild.ID] =
				&core.GuildMapItem{
					Type:    core.GuildMapTypeGuild,
					Guild:   guild,
					Members: make(map[string]*discordgo.Member),
				}
			cv.GuildMapPositionIDS = append(cv.GuildMapPositionIDS, guild.ID)
			channels := drawGuild(v, guild)
			for _, channel := range channels {
				cv.State.GuildMap[channel.ID] = &core.GuildMapItem{
					Type:    core.GuildMapTypeChannel,
					Channel: channel,
				}
				cv.GuildMapPositionIDS = append(cv.GuildMapPositionIDS, channel.ID)
			}
		}
		return nil
	})
	return nil
}

func (cv *ChannelsView) joinChannel(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, xy := v.Origin()
	pos := cy + xy
	if pos < len(cv.GuildMapPositionIDS) {
		if cv.State.GuildMap[cv.GuildMapPositionIDS[pos]].Type == core.GuildMapTypeChannel {
			switch cv.State.GuildMap[cv.GuildMapPositionIDS[pos]].Channel.Type {
			case discordgo.ChannelTypeGuildText:
				cv.channelChangedCallback(cv.State.GuildMap[cv.GuildMapPositionIDS[pos]].Channel)

			case discordgo.ChannelTypeGuildVoice:
				cv.channelChangedCallback(cv.State.GuildMap[cv.GuildMapPositionIDS[pos]].Channel)
			}
		}
	}
	return nil
}

func drawGuild(w io.Writer, guild *discordgo.Guild) []*discordgo.Channel {
	fmt.Fprintln(w, guild.Name)

	return drawChannels(w, guild.Channels)
}

func drawChannels(w io.Writer, channels []*discordgo.Channel) []*discordgo.Channel {
	guildMapChannels := []*discordgo.Channel{}
	for _, category := range core.SortChannels(channels) {
		guildMapChannels = append(guildMapChannels, category.Channel)
		fmt.Fprintln(w, formatChannel(category.Channel), color.BlueString(category.Channel.Name))
		for _, channel := range category.Channels {
			guildMapChannels = append(guildMapChannels, channel)
			fmt.Fprintln(w, " ", formatChannel(channel), color.BlueString(channel.Name))
		}

	}
	return guildMapChannels
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
	if err := cv.State.Gui.SetKeybinding(channelsView, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := cv.State.Gui.SetKeybinding(channelsView, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := cv.State.Gui.SetKeybinding(channelsView, gocui.KeyEnter, gocui.ModNone, cv.joinChannel); err != nil {
		return err
	}
	return nil
}
