// (c) Jisin0

package plugins

import (
	"fmt"
	"regexp"

	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

var cbEditPattern *regexp.Regexp = regexp.MustCompile(`edit\((.+)\)`)

func Start(bot *gotgbot.Bot, update *ext.Context) error {
	go DB.AddUser(update.EffectiveMessage.From.Id)
	bot.SendMessage(
		update.Message.Chat.Id,
		fmt.Sprintf(utils.TEXT["START"], update.Message.From.FirstName, bot.FirstName),
		&gotgbot.SendMessageOpts{
			ParseMode: "HTML",
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: utils.BUTTONS["START"],
			},
			AllowSendingWithoutReply: true,
		})
	return nil
}

func Stats(bot *gotgbot.Bot, update *ext.Context) error {
	update.EffectiveMessage.Reply(bot, DB.Stats(), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	return nil
}

func GetId(bot *gotgbot.Bot, ctx *ext.Context) error {
	var text string
	update := ctx.Message

	if update.ReplyToMessage != nil {
		text += fmt.Sprintf("\nReplied to user : <code>%v</code>", update.ReplyToMessage.From.Id)
		if update.ReplyToMessage.ForwardDate != 0 {
			text += fmt.Sprintf("\nForwarded from : <code>%v</code>", update.ReplyToMessage.ForwardFromChat.Id)
		}
	}
	text += fmt.Sprintf("\nUser id : <code>%v</code>", update.From.Id)
	if update.Chat.Type != "private" {
		text += fmt.Sprintf("\nChat id : <code>%v</code>", update.From.Id)
	}

	update.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyToMessageId: update.MessageId})

	return nil
}

func CbStats(bot *gotgbot.Bot, update *ext.Context) error {
	update.CallbackQuery.Message.EditText(bot, DB.Stats(), &gotgbot.EditMessageTextOpts{
		ChatId:      update.CallbackQuery.Message.Chat.Id,
		MessageId:   update.CallbackQuery.Message.MessageId,
		ParseMode:   "HTML",
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: utils.BUTTONS["STATS"]},
	})
	return nil
}

func FilterHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	//What R U Lookking At Its Just a Pro Function ;)

	go MFilter(b, ctx)
	go GFilter(b, ctx)

	return nil
}

func CbEdit(bot *gotgbot.Bot, update *ext.Context) error {
	//Function to handle edit() callbacks from the Start, About and Help menus

	key := cbEditPattern.FindStringSubmatch(update.CallbackQuery.Data)[1]
	markup, ok := utils.BUTTONS[key]
	if !ok {
		markup = [][]gotgbot.InlineKeyboardButton{{{Text: "⤝ Bᴀᴄᴋ", CallbackData: "edit(MAP)"}}}
	}
	options := gotgbot.EditMessageTextOpts{
		ChatId:                update.CallbackQuery.Message.Chat.Id,
		MessageId:             update.CallbackQuery.Message.MessageId,
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
		ReplyMarkup:           gotgbot.InlineKeyboardMarkup{InlineKeyboard: markup},
	}
	var text string
	if key == "START" {
		text = fmt.Sprintf(utils.TEXT["START"], update.CallbackQuery.From.FirstName, bot.FirstName)
	} else {
		text = utils.TEXT[key]
	}
	_, _, err := update.CallbackQuery.Message.EditText(bot,
		text,
		&options,
	)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func About(b *gotgbot.Bot, update *ext.Context) error {
	update.Message.Reply(b, utils.TEXT["ABOUT"], &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	return nil
}

func Help(b *gotgbot.Bot, update *ext.Context) error {
	update.Message.Reply(b, utils.TEXT["HELP"], &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	return nil
}
