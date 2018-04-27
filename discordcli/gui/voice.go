package gui

import (
	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/core"
	"github.com/jroimartin/gocui"
)

/*
VoiceView shows messages in current selected channel
*/
type VoiceView struct {
	State           *core.State
	VoiceConnection *core.VoiceConnection
}

func (vv *VoiceView) render() error {
	_, maxY := vv.State.Gui.Size()
	if v, err := vv.State.Gui.SetView(voiceView, 0, maxY-25, 25, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil
}

func (vv *VoiceView) joinVoice(channel *discordgo.Channel) error {
	voice := core.CreateVoiceConnection(vv.State.Session, channel)
	vv.VoiceConnection = voice
	voice.Start()
	return nil
}
