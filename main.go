package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
	"github.com/Jisin0/Go-Filter-Bot/utils/autodelete"
	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// Print reason for panic + stack for some sort of helpful log output
			fmt.Println(r)
			fmt.Println(string(debug.Stack()))
		}
	}()

	// Run a useless http server to get a healthy build on koyeb
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "Waku Waku")
		})

		http.ListenAndServe(":"+config.Port, nil)
	}()

	if config.BotToken == "" {
		panic("Exiting Because No BOT_TOKEN Provided :(")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(config.BotToken, &gotgbot.BotOpts{
		BotClient: &gotgbot.BaseBotClient{
			Client: http.Client{},
			DefaultRequestOpts: &gotgbot.RequestOpts{
				Timeout: gotgbot.DefaultTimeout, // Customise the default request timeout here
				APIURL:  gotgbot.DefaultAPIURL,  // As well as the Default API URL here (in case of using local bot API servers)
			},
		},
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// To make sure no other instance of the bot is running
	_, err = b.GetUpdates(&gotgbot.GetUpdatesOpts{})
	if err != nil {
		fmt.Println("waiting 10s because : " + err.Error())
		time.Sleep(10 * time.Second)
	}

	updater := ext.NewUpdater(plugins.Dispatcher, &ext.UpdaterOpts{})

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	fmt.Printf("@%s Started !\n", b.User.Username)

	if plugins.AutoDelete > 0 {
		go autodelete.RunAutodel(b)
	}

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
