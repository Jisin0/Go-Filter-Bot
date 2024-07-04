package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/plugins"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var (
	allowedTokens    = strings.Split(os.Getenv("BOT_TOKEN"), " ")
	lenAllowedTokens = len(allowedTokens)
)

const (
	statusCodeSuccess = 200
)

// Handles all incoming traffic from webhooks.
func Bot(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path

	_, botToken := path.Split(url)

	bot, _ := gotgbot.NewBot(botToken, &gotgbot.BotOpts{DisableTokenCheck: true})

	// Delete the webhook incase token is unauthorized.
	if lenAllowedTokens > 0 && allowedTokens[0] != "" && !slices.Contains(allowedTokens, botToken) {
		bot.DeleteWebhook(&gotgbot.DeleteWebhookOpts{}) //nolint:errcheck // It doesn't matter if it errors
		w.WriteHeader(statusCodeSuccess)

		return
	}

	var update gotgbot.Update

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading request body: %v", err)
		w.WriteHeader(statusCodeSuccess)

		return
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		fmt.Println("failed to unmarshal body ", err)
		w.WriteHeader(statusCodeSuccess)

		return
	}

	ctx := ext.NewContext(&update, map[string]interface{}{})

	if msg := ctx.Message; msg != nil {
		if len(msg.Entities) > 0 {
			if msg.Entities[0].Type == "bot_command" {
				split := strings.Split(strings.ToLower(strings.Fields(msg.Text)[0]), "@")
				cmd := split[0][1:]

				if cmd == "start" {
					err = plugins.Start(bot, ctx)
				}
			}
		}
	} else if ctx.InlineQuery != nil {
		err = plugins.InlineQueryHandler(bot, ctx)
	} else if ctx.ChosenInlineResult != nil {
		err = plugins.InlineResultHandler(bot, ctx)
	}

	if err != nil {
		fmt.Printf("error while handling update %v\n", err)
	}

	w.WriteHeader(statusCodeSuccess)
}
