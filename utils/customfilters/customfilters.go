// (c) Jisin0

package customfilters

import (
	"fmt"

	"github.com/Jisin0/Go-Filter-Bot/database"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var DB *database.Database = database.NewDatabase()
var CachedAdmins map[int64][]int64 = map[int64][]int64{}

func Listen(m *gotgbot.Message) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return m.From.Id == msg.From.Id && m.MessageId < msg.MessageId && m.Chat.Id == msg.Chat.Id && msg.Text != "  "
	}
}

func PrivateOrGroup(msg *gotgbot.Message) bool {
	// A Function To Filter Group & SuperGroup Message
	return msg.Chat.Type == gotgbot.ChatTypeSupergroup || msg.Chat.Type == gotgbot.ChatTypeGroup || msg.Chat.Type == gotgbot.ChatTypePrivate
}

func Chats(chatID []int64) filters.Message {
	return func(msg *gotgbot.Message) bool {
		for _, c := range chatID {
			if c == msg.Chat.Id {
				return true
			}
		}

		return false
	}
}

//nolint:errcheck // too lazy
func Verify(bot *gotgbot.Bot, ctx *ext.Context) (int64, bool) {
	var (
		userID int64
		msg    gotgbot.MaybeInaccessibleMessage
	)

	if ctx.CallbackQuery != nil {
		msg = ctx.CallbackQuery.Message
		userID = ctx.CallbackQuery.From.Id
	} else {
		msg = ctx.Message
		userID = msg.(*gotgbot.Message).From.Id
	}

	chatType := msg.GetChat().Type
	chatID := msg.GetChat().Id

	var c int64

	switch chatType {
	case gotgbot.ChatTypeSupergroup, gotgbot.ChatTypeGroup:
		if userID == 0 {
			bot.SendMessage(
				chatID,
				"Sorry It Looks Like You Are Anonymous Please Connect From Pm And Use Me Or Turn Off Anonymous :(",
				&gotgbot.SendMessageOpts{
					ReplyParameters: &gotgbot.ReplyParameters{
						MessageId: msg.GetMessageId(),
					},
				},
			)

			return c, false
		}

		cachedAdmins, ok := CachedAdmins[chatID]
		if !ok {
			admins, e := msg.GetChat().GetAdministrators(bot, &gotgbot.GetChatAdministratorsOpts{})
			if e != nil {
				return c, false
			}

			var newAdmins []int64

			for _, admin := range admins {
				newAdmins = append(newAdmins, admin.GetUser().Id)
			}

			CachedAdmins[chatID] = newAdmins

			for _, admin := range newAdmins {
				if userID == admin {
					return c, true
				}
			}

			if ctx.CallbackQuery == nil {
				bot.SendMessage(
					chatID,
					"Who dis non-admin telling me what to do !",
					&gotgbot.SendMessageOpts{
						ReplyParameters: &gotgbot.ReplyParameters{
							MessageId: msg.GetMessageId(),
						},
					},
				)
			} else {
				ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Who dis non-admin telling me what to do !", ShowAlert: true})
			}

			return c, false
		} else {
			for _, admin := range cachedAdmins {
				if userID == admin {
					return c, true
				}
			}

			if ctx.CallbackQuery == nil {
				bot.SendMessage(
					chatID,
					"Hey You're Not An Admin, If You Are A New Admin Use The /updateadmins Command To Update Current List !",
					&gotgbot.SendMessageOpts{
						ReplyParameters: &gotgbot.ReplyParameters{
							MessageId: msg.GetMessageId(),
						},
					},
				)
			} else {
				ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Who dis non-admin telling me what to do !", ShowAlert: true})
			}

			return c, false
		}
	case gotgbot.ChatTypePrivate:
		c, ok := DB.GetConnection(userID)
		if !ok {
			bot.SendMessage(
				chatID,
				"Please connect to a chat first to use this operation !",
				&gotgbot.SendMessageOpts{
					ReplyParameters: &gotgbot.ReplyParameters{
						MessageId: msg.GetMessageId(),
					},
				},
			)

			return c, false
		}

		return c, true
	default:
		fmt.Println("Unknown ChatType ", chatType)
		return c, false
	}
}
