package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/gui"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	sess, err := discordgo.New("xD", "xD")
	if err != nil {
		panic(err)
	}

	err = gui.CreateMainView(sess, g)
	if err != nil {
		panic(err)
	}
	sess.Open()
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
