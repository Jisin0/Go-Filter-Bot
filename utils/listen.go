// (c) Jisin0

package utils

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var Listening []Listeners

type Listeners struct {
	Filter   filters.Message
	Callback handlers.Response
}

func Listen(filter filters.Message, callback_func handlers.Response) error {
	Listening = append(Listening, Listeners{Filter: filter, Callback: callback_func})
	return nil
}

func IsListening(bot *gotgbot.Bot, update *ext.Context) error {

	for i, u := range Listening {
		if u.Filter(update.Message) {
			u.Callback(bot, update)

			Listening[i] = Listening[len(Listening)-1] // Copy last element to index i.
			Listening[len(Listening)-1] = Listeners{}  // Erase last element (write zero value).
			Listening = Listening[:len(Listening)-1]   //Completely Remove That Last Value
			return nil

		}
	}
	if update.Message.ForwardFromChat != nil && update.Message.ForwardFromChat.Id != 0 && update.Message.Chat.Type == "private" {
		text := fmt.Sprintf("This Message Was Forwarded From : <code>%v</code>", update.Message.ForwardFromChat.Id)
		update.Message.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	}
	return nil
}
