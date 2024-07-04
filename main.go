package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
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
		port := os.Getenv("PORT")

		if port == "" {
			port = "8080"
		}

		http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "Waku Waku")
		})

		http.ListenAndServe(":"+port, nil)
	}()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("Exiting Because No BOT_TOKEN Provided :(")
	}

	// Create bot from environment value.
	b, err := gotgbot.NewBot(token, &gotgbot.BotOpts{
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

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
