// (c) Jisin0

package plugins

import (
	"fmt"

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

// Create updater and Dispatcher.
var Dispatcher = ext.NewDispatcher(&ext.DispatcherOpts{
	// If an error is returned by a handler, log it and continue going.
	Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
		fmt.Println("an error occurred while handling update:", err.Error())
		return ext.DispatcherActionNoop
	},
	MaxRoutines: ext.DefaultMaxRoutines,
})

func init() {
	Dispatcher.AddHandlerToGroup(handlers.NewMessage(customfilters.PrivateOrGroup, FilterHandler), filterHandlerGroup)

	Dispatcher.AddHandlerToGroup(handlers.NewCommand("start", Start), basicCommandsGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("about", About), basicCommandsGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("help", Help), basicCommandsGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("stats", Stats), basicCommandsGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("id", GetID), basicCommandsGroup)

	Dispatcher.AddHandlerToGroup(handlers.NewCommand("filter", NewFilter), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilter", NewFilter), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("filters", AllMfilters), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("stop", StopMfilter), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("gstop", StopGfilter), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("connect", Connect), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("disconnect", Disconnect), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("startglobal", StartGlobal), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("gfilters", Gfilters), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("broadcast", Broadcast), commandHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCommand("concast", Broadcast), commandHandlerGroup)

	Dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("cbconnect("), CbConnect), callbackHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("stopf("), CbStop), callbackHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("edit("), CbEdit), callbackHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("alert("), CbAlert), callbackHandlerGroup)
	Dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Equal("stats"), CbStats), callbackHandlerGroup)

	Dispatcher.AddHandlerToGroup(handlers.NewMessage(message.All, utils.RunListening), miscHandlerGroup)
}
