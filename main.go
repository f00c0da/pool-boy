package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"pool-boy/app"
	"pool-boy/core"
	"pool-boy/ui"
	"time"
)

const (
	WatchCmd   = "watch"
	UsePoolCmd = "pool"
)

func main() {
	cliApp, appData := createCliApp()

	cliApp.Action = func(context *cli.Context) {

		if context.Bool(UsePoolCmd) {
			core.GetPool()

		} else if context.Bool(WatchCmd) {
			if context.NArg() > 0 {
				appData.TextToSpeech = context.Args().Get(0)
			}
			appData.SpeechActive = true
			startWatching(&appData)
		}
	}

	_ = cliApp.Run(os.Args)
}

func startWatching(appData *app.AppData) {
	config := core.GetConfig()
	events := core.GetEvents()

	if config.IsUrlNotValid() || len(events) == 0 {
		fmt.Printf("please use the option --pool\ncan not load config file or event file\n")
		return
	}

	uiApplication, uiElements := ui.InitializeUI(&events)
	app.StartEventPolling(appData, config, uiApplication, uiElements, &events)

	if err := uiApplication.Run(); err != nil {
		panic(err)
	}
}

func createCliApp() (*cli.App, app.AppData) {
	appData := app.AppData{
		StartTime:    time.Now(),
		RequestCount: 0,
		LastMessage:  "use <<enter>> to watch a event",
		SpeechActive: false,
		TextToSpeech: "Junge, Ticket koofen",
	}

	cliApp := cli.NewApp()
	cliApp.Name = "pool-boy"
	cliApp.Version = "1.3.1"
	cliApp.Usage = "Fetch pool events from the Berliner Bäder Betriebe"
	cliApp.Authors = []cli.Author{{Name: "__                         _                     \n  / _|   _   _      ___      | |__   ___  _   _ ____\n | |_   | | | |    / __|     | '_ \\ / _ \\| | | |_  /\n |  _|  | |_| |_  | (__ _    | |_) | (_) | |_| |/ / \n |_|(_)  \\__,_(_)  \\___(_)   |_.__/ \\___/ \\__, /___|\n                                         |___/     \n"}}
	cliApp.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  UsePoolCmd,
			Usage: "select a swimming pool",
		},
		cli.BoolFlag{
			Name:  WatchCmd,
			Usage: "view and interact with events of a pool",
		},
	}
	return cliApp, appData
}
