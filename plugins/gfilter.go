// (c) Jisin0

package plugins

import (
	"fmt"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/database"
	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/Jisin0/Go-Filter-Bot/utils/config"
	"github.com/Jisin0/Go-Filter-Bot/utils/customfilters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func GFilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		chatID    int64
		update    = ctx.Message
		messageID = update.MessageId
	)

	switch chatType := update.Chat.Type; chatType {
	case gotgbot.ChatTypePrivate:
		var ok bool

		chatID, ok = DB.GetConnection(update.From.Id)
		if !ok {
			return nil
		}
	case gotgbot.ChatTypeSupergroup, gotgbot.ChatTypeGroup:
		chatID = update.Chat.Id
	default:
		return nil
	}

	var message string
	switch {
	case update.Text != "":
		message = update.Text
	case update.Caption != "":
		message = update.Caption
	default:
		return nil
	}

	stopped := DB.GetCachedSetting(chatID).Stopped

	var results []*database.Filter

	fields := strings.Fields(message)
	if len(fields) <= 15 { // uses new method only if input has <=15 substrings
		results = DB.SearchMfilterNew(globalNumber, fields, config.MultiFilter)
	} else {
		results = DB.SearchMfilterClassic(globalNumber, message)
	}

	for _, f := range results {
		if utils.Contains(stopped, f.Text) {
			continue
		}

		sendFilter(f, bot, update, chatID, messageID)
	}

	return nil
}

// Function to handle the startglobal command
//
//nolint:errcheck // hmm
func StartGlobal(bot *gotgbot.Bot, ctx *ext.Context) error {
	update := ctx.Message

	c, v := customfilters.Verify(bot, ctx)
	if !v {
		return nil
	}

	if c == 0 {
		c = ctx.Message.Chat.Id
	}

	split := strings.SplitN(update.Text, " ", 2)
	if split[1] == "" {
		update.Reply(bot, "Bad Usage No Keyword Provided :(", &gotgbot.SendMessageOpts{})
	} else {
		key := split[1]

		_, ok := DB.GetMfilter(globalNumber, key)
		if !ok {
			update.Reply(bot, fmt.Sprintf("No Global Filter For %v Was Found To Restart !", key), &gotgbot.SendMessageOpts{})
		} else {
			for _, k := range DB.GetCachedSetting(c).Stopped {
				if k == key {
					DB.StartGfilter(c, key)
					update.Reply(bot, fmt.Sprintf("Restarted Global Filter For <i>%v</i> Successfully !", key), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})

					return nil
				}
			}

			update.Reply(bot, fmt.Sprintf("You Havent Stopped Any Global Filter For %v :(", key), &gotgbot.SendMessageOpts{})
		}
	}

	return nil
}

func Gfilters(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Function to handle /gfilters function
	text := DB.StringMfilter(globalNumber)

	_, err := ctx.Message.Reply(bot, "Aʟʟ ғɪʟᴛᴇʀs sᴀᴠᴇᴅ ғᴏʀ ɢʟᴏʙᴀʟ ᴜsᴀɢᴇ :\n"+text, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		fmt.Printf("gfilters: %v\n", err)
	}

	return nil
}

// Function to handle the gstop command
//
//nolint:errcheck // too many
func StopGfilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	update := ctx.Message

	if !utils.IsAdmin(ctx.EffectiveUser.Id) {
		update.Reply(bot, "Only bot admins can use this command !", &gotgbot.SendMessageOpts{})
		return nil
	}

	var (
		split = strings.SplitN(update.Text, " ", 2)
		key   string
	)

	switch {
	case len(split) < 2:
		m := utils.Ask(bot, "Ok Now Send Me The Name OF The Filter You Would Like To Stop ...", ctx.EffectiveChat, ctx.EffectiveUser)
		if m == nil {
			return nil
		}

		key = m.Text
	default:
		key = split[1]
	}

	// Checking if theres a global filter for the key
	_, ok := DB.GetMfilter(globalNumber, key)

	// If there isnt local or global
	if !ok {
		update.Reply(bot, fmt.Sprintf("I Couldnt Find Any Global Filter For <code>%v</code> To Stop !", key), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return nil
	}

	DB.DeleteMfilter(globalNumber, key)
	update.Reply(bot, fmt.Sprintf("Global Filter For <i>%v</i> Was Stopped Successfully !", key), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})

	return nil
}
