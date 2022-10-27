// (c) Jisin0

package customfilters

import (
	"fmt"

	"github.com/Jisin0/Go-Filter-Bot/database"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var DB database.Database = database.NewDatabase()
var CachedAdmins map[int64][]int64 = map[int64][]int64{}

func Listen(m *gotgbot.Message) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return m.From.Id == msg.From.Id && m.MessageId < msg.MessageId && m.Chat.Id == msg.Chat.Id && msg.Text != "  "
	}
}

func PrivateOrGroup(msg *gotgbot.Message) bool {
	//A Function To Filter Group & SuperGroup Message

	return msg.Chat.Type == "supergroup" || msg.Chat.Type == "group" || msg.Chat.Type == "private"
}

func Chats(chatId []int64) filters.Message {
	return func(msg *gotgbot.Message) bool {
		for _, c := range chatId {
			if c == msg.Chat.Id {
				return true
			}
		}
		return false
	}
}

func Verify(bot *gotgbot.Bot, ctx *ext.Context) (int64, bool) {
	var user_id int64
	msg := ctx.Message
	if ctx.CallbackQuery != nil {
		msg = ctx.CallbackQuery.Message
		user_id = ctx.CallbackQuery.From.Id
	} else {
		user_id = msg.From.Id
	}
	chatType := msg.Chat.Type
	var c int64
	if chatType == "supergroup" || chatType == "group" {
		if user_id == 0 {
			msg.Reply(
				bot,
				"Sorry It Looks Like You Are Anonymous Please Connect From Pm And Use Me Or Turn Off Anonymous :(",
				&gotgbot.SendMessageOpts{},
			)
			return c, false
		}
		cachedAdmins, ok := CachedAdmins[msg.Chat.Id]
		if !ok {
			admins, e := msg.Chat.GetAdministrators(bot)
			var newAdmins []int64
			if e != nil {
				return c, false
			}
			for _, admin := range admins {
				newAdmins = append(newAdmins, admin.GetUser().Id)
			}

			CachedAdmins[msg.Chat.Id] = newAdmins

			for _, admin := range newAdmins {
				if user_id == admin {
					return c, true
				}
			}
			if ctx.CallbackQuery == nil {
				msg.Reply(
					bot,
					"Who dis non-admin telling me what to do !",
					&gotgbot.SendMessageOpts{},
				)
			} else {
				ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Who dis non-admin telling me what to do !", ShowAlert: true})
			}
			return c, false
		} else {
			for _, admin := range cachedAdmins {
				if user_id == admin {
					return c, true
				}
			}

			if ctx.CallbackQuery == nil {
				msg.Reply(
					bot,
					"Hey You're Not An Admin, If You Are A New Admin Use The /updateadmins Command To Update Current List !",
					&gotgbot.SendMessageOpts{ReplyToMessageId: msg.MessageId},
				)
			} else {
				ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Who dis non-admin telling me what to do !", ShowAlert: true})
			}

			return c, false
		}
	} else if chatType == "private" {
		c, ok := DB.GetConnection(user_id)
		if !ok {
			msg.Reply(
				bot,
				"Sorry You Have To Connect To A Chat To Use This Command Here :(",
				&gotgbot.SendMessageOpts{},
			)
			return c, false
		}
		fmt.Println(c)
		return c, true
	} else {
		fmt.Println("Unknown ChatType ", chatType)
		return c, false
	}
}
