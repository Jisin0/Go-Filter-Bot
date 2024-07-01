package plugins

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var cbConnectRegex *regexp.Regexp = regexp.MustCompile(`cbconnect\((.+)\)`)

func Connect(bot *gotgbot.Bot, update *ext.Context) error {
	// Check for any existing connections first
	_, k := DB.GetConnection(update.Message.From.Id)
	if k {
		_, err := update.Message.Reply(
			bot,
			"Looks Like You Are Already Connected To A Chat, Please /disconnect From It To Connect To Another One :)",
			&gotgbot.SendMessageOpts{},
		)
		if err != nil {
			fmt.Println(err)
		}

		return nil
	} else {
		if update.Message.Chat.Type == gotgbot.ChatTypePrivate {
			// Require a chat id if the command was used in a private chat
			var (
				args    = strings.Split(update.Message.Text, " ")
				chatRaw string
			)

			if len(args) < 2 {
				m := utils.Ask(bot, "Please send the id of the chat you would like to connect to : ", update.EffectiveChat, update.EffectiveUser)
				if m == nil {
					return nil
				}

				chatRaw = m.Text
			} else {
				chatRaw = args[1]
			}

			chatID, e := strconv.ParseInt(chatRaw, 0, 64)
			if e != nil {
				// If converion of raw chat_id to an int64 fails i.e it isnt a number
				_, err := update.Message.Reply(
					bot,
					"That Doesnt Seem Like A Valid ChatId A ChatId Looks Something Like -100xxxxxxxxxx :(",
					&gotgbot.SendMessageOpts{},
				)
				if err != nil {
					fmt.Println(err)
				}

				return nil
			} else {
				// Verify and connect
				admins, err := bot.GetChatAdministrators(chatID, &gotgbot.GetChatAdministratorsOpts{})
				if err != nil {
					_, err := update.Message.Reply(
						bot,
						fmt.Sprintf("Sorry Looks Like I Couldnt Find That Chat With Id <code>%v</code>. Make Sure I'm Admin There With Full Permissions :(", chatID),
						&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
					)
					if err != nil {
						fmt.Println(err)
					}

					return nil
				}

				for _, admin := range admins {
					if update.Message.From.Id == admin.GetUser().Id {
						DB.ConnectUser(update.Message.From.Id, chatID)
						_, err := update.Message.Reply(
							bot,
							"Awesome I've Successfully Connected You To Your Group !",
							&gotgbot.SendMessageOpts{},
						)

						if err != nil {
							fmt.Println(err)
						}

						return nil
					}
				}

				_, err = update.Message.Reply(
					bot,
					"You Cant Connect To A Chat Where You're Not Admin :)",
					&gotgbot.SendMessageOpts{},
				)
				if err != nil {
					fmt.Println(err)
				}

				return nil
			}
		} else if update.Message.Chat.Type == gotgbot.ChatTypeSupergroup || update.Message.Chat.Type == gotgbot.ChatTypeGroup {
			// For groups or supergroups just connect
			if update.Message.From.Id == 0 {
				// Connect using button in case user is anonymous
				_, err := update.Message.Reply(
					bot,
					"It Looks Like You Are Anonymous Click The Button Below To Connect :(",
					&gotgbot.SendMessageOpts{
						ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "Connect Me", CallbackData: "cbconnect(con)"}}}},
					},
				)
				if err != nil {
					fmt.Println(err)
				}

				return nil
			} else {
				// Verification stuff
				admins, err := bot.GetChatAdministrators(update.Message.Chat.Id, &gotgbot.GetChatAdministratorsOpts{})
				if err != nil {
					_, err := update.Message.Reply(
						bot,
						"Sorry I Couldn't access the admins list of this chat!\nPlease make sure I'm an admin here.",
						&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
					)
					if err != nil {
						fmt.Println(err)
					}

					return nil
				}

				for _, admin := range admins {
					if update.Message.From.Id == admin.GetUser().Id {
						DB.ConnectUser(update.Message.From.Id, update.Message.Chat.Id)
						_, err := update.Message.Reply(
							bot,
							"Awesome I've Successfully Connected You To This Group !",
							&gotgbot.SendMessageOpts{},
						)

						if err != nil {
							fmt.Println(err)
						}

						return nil
					}
				}

				_, err = update.Message.Reply(
					bot,
					"Ok Mr. Non-Admin :)",
					&gotgbot.SendMessageOpts{},
				)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return nil
}

// CbConnecthandles callback from connect button.
//
//nolint:errcheck // Ignore alert errors
func CbConnect(bot *gotgbot.Bot, ctx *ext.Context) error {
	update := ctx.CallbackQuery

	rMatch := cbConnectRegex.FindStringSubmatch(update.Data)
	if len(rMatch) < 2 {
		fmt.Printf("cbconnect: bad callback data %s", update.Data)
		return nil
	}

	action := rMatch[1]
	if action == "con" {
		admins, _ := bot.GetChatAdministrators(update.Message.GetChat().Id, &gotgbot.GetChatAdministratorsOpts{})

		for _, admin := range admins {
			if update.From.Id == admin.GetUser().Id {
				DB.ConnectUser(update.From.Id, update.Message.GetChat().Id)

				update.Answer(
					bot,
					&gotgbot.AnswerCallbackQueryOpts{Text: "Awesome I've Successfully Connected You To This Group !", ShowAlert: true},
				)

				update.Message.Delete(bot, &gotgbot.DeleteMessageOpts{})

				return nil
			}
		}

		update.Answer(
			bot,
			&gotgbot.AnswerCallbackQueryOpts{Text: "Ok Mr. Non-Admin :)", ShowAlert: true},
		)
	} else if action == "dis" {
		DB.DeleteConnection(update.From.Id)

		update.Answer(
			bot,
			&gotgbot.AnswerCallbackQueryOpts{ShowAlert: true, Text: "All Of Your Connections Were Cleared :)"},
		)
	}

	return nil
}

// Function to handle /diconnect command
//
//nolint:errcheck // no need
func Disconnect(bot *gotgbot.Bot, update *ext.Context) error {
	if update.Message.From.Id == 0 {
		update.Message.Reply(
			bot,
			"Sorry Looks Like You Are Anonymous Use The Button Below To Prove Your Identity :)",
			&gotgbot.SendMessageOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "Disconnect Me", CallbackData: "cbconnect(dis)"}}}}},
		)
	} else {
		DB.DeleteConnection(update.Message.From.Id)

		update.Message.Reply(
			bot,
			"Any Existing Connections Were Cleared Successfully :)",
			&gotgbot.SendMessageOpts{},
		)
	}

	return nil
}
