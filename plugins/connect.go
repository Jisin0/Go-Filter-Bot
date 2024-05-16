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
	//Check for any existing connections first
	_, k := DB.GetConnection(update.Message.From.Id)
	if k {
		update.Message.Reply(
			bot,
			"Looks Like You Are Already Connected To A Chat, Please /disconnect From It To Connect To Another One :)",
			&gotgbot.SendMessageOpts{},
		)
		return nil
	} else {
		if update.Message.Chat.Type == "private" {
			//Require a chat id if the command was used in a private chat
			args := strings.Split(update.Message.Text, " ")
			var chat_raw string
			if len(args) < 2 {
				m := utils.Ask(bot, "Please send the id of the chat you would like to connect to : ", update.EffectiveChat, update.EffectiveUser)
				if m == nil {
					return nil
				}

				chat_raw = m.Text
			} else {
				chat_raw = args[1]
			}
			chat_id, e := strconv.ParseInt(chat_raw, 0, 64)
			if e != nil {
				//If converion of raw chat_id to an int64 fails i.e it isnt a number
				update.Message.Reply(
					bot,
					"That Doesnt Seem Like A Valid ChatId A ChatId Looks Something Like -100xxxxxxxxxx :(",
					&gotgbot.SendMessageOpts{},
				)
				return nil
			} else {
				//Verify and connect
				admins, err := bot.GetChatAdministrators(chat_id, &gotgbot.GetChatAdministratorsOpts{})
				if err != nil {
					update.Message.Reply(
						bot,
						fmt.Sprintf("Sorry Looks Like I Couldnt Find That Chat With Id <code>%v</code>. Make Sure I'm Admin There With Full Permissions :(", chat_id),
						&gotgbot.SendMessageOpts{ParseMode: "HTML"},
					)
				} else {
					for _, admin := range admins {
						if update.Message.From.Id == admin.GetUser().Id {
							DB.ConnectUser(update.Message.From.Id, chat_id)
							update.Message.Reply(
								bot,
								"Awesome I've Succesfully Connected You To Your Group !",
								&gotgbot.SendMessageOpts{},
							)
							return nil
						}
					}

					update.Message.Reply(
						bot,
						"You Cant Connect To A Chat Where You're Not Admin :)",
						&gotgbot.SendMessageOpts{},
					)
					return nil
				}
			}
		} else if update.Message.Chat.Type == "supergroup" || update.Message.Chat.Type == "group" {
			//For groups or supergroups just connect
			if update.Message.From.Id == 0 {
				//Connect using button in case user is anonymous
				update.Message.Reply(
					bot,
					"It Looks Like You Are Anonymous Click The Button Below To Connect :(",
					&gotgbot.SendMessageOpts{
						ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "Connect Me", CallbackData: "cbconnect(con)"}}}},
					},
				)
				return nil
			} else {
				//Verification stuff
				admins, _ := bot.GetChatAdministrators(update.Message.Chat.Id, &gotgbot.GetChatAdministratorsOpts{})

				for _, admin := range admins {
					if update.Message.From.Id == admin.GetUser().Id {
						DB.ConnectUser(update.Message.From.Id, update.Message.Chat.Id)
						update.Message.Reply(
							bot,
							"Awesome I've Succesfully Connected You To This Group !",
							&gotgbot.SendMessageOpts{},
						)
						return nil
					}
				}

				update.Message.Reply(
					bot,
					"Ok Mr. Non-Admin :)",
					&gotgbot.SendMessageOpts{},
				)
				return nil

			}
		}
	}

	return nil
}

func CbConnect(bot *gotgbot.Bot, update *ext.Context) error {
	//Function to handle callback from connect button
	action := cbConnectRegex.FindStringSubmatch(update.CallbackQuery.Data)[1]
	if action == "con" {
		admins, _ := bot.GetChatAdministrators(update.CallbackQuery.Message.GetChat().Id, &gotgbot.GetChatAdministratorsOpts{})

		for _, admin := range admins {
			if update.CallbackQuery.From.Id == admin.GetUser().Id {
				DB.ConnectUser(update.CallbackQuery.From.Id, update.CallbackQuery.Message.GetChat().Id)
				update.CallbackQuery.Answer(
					bot,
					&gotgbot.AnswerCallbackQueryOpts{Text: "Awesome I've Succesfully Connected You To This Group !", ShowAlert: true},
				)
				update.CallbackQuery.Message.Delete(bot, &gotgbot.DeleteMessageOpts{})
				return nil
			}
		}

		update.CallbackQuery.Answer(
			bot,
			&gotgbot.AnswerCallbackQueryOpts{Text: "Ok Mr. Non-Admin :)", ShowAlert: true},
		)
		return nil
	} else if action == "dis" {
		DB.DeleteConnection(update.CallbackQuery.From.Id)
		update.CallbackQuery.Answer(
			bot,
			&gotgbot.AnswerCallbackQueryOpts{ShowAlert: true, Text: "All Of Your Connections Were Cleared :)"},
		)
	}

	return nil
}

func Disconnect(bot *gotgbot.Bot, update *ext.Context) error {
	//Function to handle /diconnect command
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
