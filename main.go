package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/Jisin0/Go-Filter-Bot/utils/customfilters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func main() {

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
		fmt.Println("Exiting because : " + err.Error())
		return
	}

	// Create updater and dispatcher.
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			fmt.Println("an error occurred while handling update:", err.Error())
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})
	updater := ext.NewUpdater(dispatcher, &ext.UpdaterOpts{})

	// Add update handlers
	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", plugins.Start), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Equal("stats"), plugins.CbStats), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("stats", plugins.Stats), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("filter", plugins.NewFilter), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilter", plugins.NewFilter), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("filters", plugins.AllMfilters), 2)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("stop", plugins.StopMfilter), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("connect", plugins.Connect), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("disconnect", plugins.Disconnect), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("cbconnect("), plugins.CbConnect), 2)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("startglobal", plugins.StartGlobal), 2)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilters", plugins.Gfilters), 2)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("broadcast", plugins.Broadcast), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("concast", plugins.ConCast), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("id", plugins.GetId), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("about", plugins.About), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("help", plugins.Help), 1)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("stopf("), plugins.CbStop), 3)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("edit("), plugins.CbEdit), 3)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(customfilters.PrivateOrGroup, plugins.FilterHandler), 1)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, utils.RunListening), 5)

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}
	fmt.Printf("@%s Started !\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
