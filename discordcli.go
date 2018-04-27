package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/feelfreelinux/discordcli/discordcli/gui"
	"github.com/gordonklaus/portaudio"
	"github.com/jroimartin/gocui"
	"github.com/shibukawa/configdir"
)

type config struct {
	Token string `json:"token"`
}

func loadConfig() (cfg config, err error) {
	dir := configdir.New("feelfreelinux", "discordcli")
	cfgFile := dir.QueryFolderContainsFile("config.json")

	if cfgFile == nil { // file not found, ask user for data
		// TODO: add nice gui
		fmt.Println("--- Copy token from app and paste ---")
		fmt.Print("Token: ")
		fmt.Scanln(&cfg.Token)
		var jsonFile *os.File
		jsonFile, err = dir.QueryFolders(configdir.Global)[0].Create("config.json")
		if err != nil {
			return
		}
		encoder := json.NewEncoder(jsonFile)
		encoder.SetIndent("", "    ")
		err = encoder.Encode(&cfg)
		if err != nil {
			return
		}
	} else {
		var jsonData []byte
		jsonData, err = cfgFile.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsonData, &cfg)
		if err != nil {
			return
		}
	}
	return
}

func main() {
	cfg, cfgErr := loadConfig()
	if cfgErr != nil {
		log.Panicf("Failed to load config %v", cfgErr)
	}
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	portaudio.Initialize()
	defer portaudio.Terminate()
	sess, err := discordgo.New(cfg.Token)
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
