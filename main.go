package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/Jisin0/Go-Filter-Bot/utils/customfilters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

const (
	filterHandlerGroup   = 1 // handler group for filters
	basicCommandsGroup   = 2 // handler group for basic commands
	commandHandlerGroup  = 3 // handler group for other cammands
	callbackHandlerGroup = 4 // handler group for callbacks
	miscHandlerGroup     = 5 // handler group for everything else
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

	dispatcher.AddHandlerToGroup(handlers.NewMessage(customfilters.PrivateOrGroup, plugins.FilterHandler), filterHandlerGroup)

	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", plugins.Start), basicCommandsGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("about", plugins.About), basicCommandsGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("help", plugins.Help), basicCommandsGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("stats", plugins.Stats), basicCommandsGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("id", plugins.GetID), basicCommandsGroup)

	dispatcher.AddHandlerToGroup(handlers.NewCommand("filter", plugins.NewFilter), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilter", plugins.NewFilter), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("filters", plugins.AllMfilters), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("stop", plugins.StopMfilter), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("gstop", plugins.StopGfilter), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("connect", plugins.Connect), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("disconnect", plugins.Disconnect), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("startglobal", plugins.StartGlobal), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilters", plugins.Gfilters), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("broadcast", plugins.Broadcast), commandHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("concast", plugins.Broadcast), commandHandlerGroup)

	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("cbconnect("), plugins.CbConnect), callbackHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("stopf("), plugins.CbStop), callbackHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("edit("), plugins.CbEdit), callbackHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("alert("), plugins.CbAlert), callbackHandlerGroup)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Equal("stats"), plugins.CbStats), callbackHandlerGroup)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, utils.RunListening), miscHandlerGroup)

	// Start receiving updates.
	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	fmt.Printf("@%s Started !\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}
