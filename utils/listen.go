// (c) Jisin0

package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var Listening []Listeners

const (
	fiveMinutes = time.Minute * 5
)

type Listeners struct {
	// Filter which checks if message is a match
	Filter filters.Message
	// Function to be called if filter matches
	Callback handlers.Response
	// Unix timestamp at which the handler was added (Timeout after 300s)
	AddedTime int64
	// Channel to which message is passed
	Channel chan *gotgbot.Message
}

// Checks if update matches any saved listen filters
func RunListening(bot *gotgbot.Bot, update *ext.Context) error {
	for i, u := range Listening {
		if !u.Filter(update.Message) {
			continue
		}

		// Delete handler from slice
		Listening[i] = Listening[len(Listening)-1] // Copy last element to index i.
		Listening[len(Listening)-1] = Listeners{}  // Erase last element (write zero value).
		Listening = Listening[:len(Listening)-1]   // Completely Remove That Last Value

		// Send message to channel
		u.Channel <- update.Message
	}

	if update.Message.ForwardOrigin != nil && update.Message.ForwardOrigin.MergeMessageOrigin().Chat != nil && update.Message.Chat.Type == gotgbot.ChatTypePrivate {
		text := fmt.Sprintf("This Message Was Forwarded From : <code>%v</code>", update.Message.ForwardOrigin.MergeMessageOrigin().Chat.Id)
		update.Message.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML}) //nolint:errcheck // not a core feature
	}

	return nil
}

var ListenTimeout error

// Listens for a message that matches the given filter.
//
// ctx - Context with ideally a deadline set
// filter - The filter that the message must match (use customfilters package for prebuilts)
//
// Returns either the message if it was received or a ListenTimeout error
func ListenMessage(ctx context.Context, filter filters.Message) (*gotgbot.Message, error) {
	// Make a channel
	c := make(chan *gotgbot.Message, 1)

	// Add details to listening slice
	Listening = append(Listening, Listeners{Filter: filter, Channel: c})

	// Listen for either a message or timeout
	select {
	case <-ctx.Done():
		return nil, ListenTimeout
	case m := <-c:
		return m, nil
	}
}

// Custom filter used when requesting a user or chat.
//
// chat - chat the message was sent to.
// user - user who is expected to reply.
// msgID - id of the last message in the chat
func listenRequestFilter(chat *gotgbot.Chat, user *gotgbot.User, msgID int64) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return user.Id == msg.From.Id && msgID < msg.MessageId && chat.Id == msg.Chat.Id
	}
}

// Sends a text message to the given chat and listens for a reply.
//
// bot - The bot that handles the conversation.
// text - Content of the request message to be sent.
// chat - The chat in which the conversion takes place
// user - The user expected to answer
func Ask(bot *gotgbot.Bot, text string, chat *gotgbot.Chat, user *gotgbot.User) *gotgbot.Message {
	// initial msg which's id is later used as the pinpoint of the converation's start
	firstM, err := bot.SendMessage(chat.Id, text, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		fmt.Println("error while asking ", err)
		return nil
	}

	// create context with 5 min timeout
	ctx, cancel := context.WithTimeout(context.Background(), fiveMinutes)
	defer cancel()

	// listen for matching messages
	msg, err := ListenMessage(ctx, listenRequestFilter(chat, user, firstM.MessageId))
	if err != nil {
		bot.SendMessage(chat.Id, "<i>Request timed out ❗</i>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return nil
	}

	return msg
}
