// (c) Jisin0

package plugins

import (
	"context"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson"
)

func Broadcast(bot *gotgbot.Bot, ctx *ext.Context) error {
	//Function to handle /broadcast command
	if !IsAdmin(ctx.Message.From.Id) {
		ctx.Message.Reply(bot, "<b>Sorry Thats An Admin Only Command :(</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyParameters: &gotgbot.ReplyParameters{MessageId: ctx.Message.MessageId}})
		return nil
	}

	cursor, _ := DB.Ucol.Find(context.TODO(), bson.D{{}})
	update := ctx.Message

	if update.ReplyToMessage == nil {
		update.Reply(bot, "Please reply this command to the message you would like to broadcast !", &gotgbot.SendMessageOpts{})
		return nil
	}

	msg := update.ReplyToMessage

	var isText bool = false
	var isMedia bool = false
	var caption string
	var markup *gotgbot.InlineKeyboardMarkup = &gotgbot.InlineKeyboardMarkup{}
	var total int
	var sent int
	var failed int
	var id int64

	if msg.Text != "" {
		isText = true
	} else {
		isMedia = true
		caption = msg.OriginalCaptionHTML()
	}

	if msg.ReplyMarkup != nil {
		markup = &gotgbot.InlineKeyboardMarkup{InlineKeyboard: msg.ReplyMarkup.InlineKeyboard}
	}

	stat, _ := update.Reply(bot, "<code>Starting broadcast ...</code>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	for cursor.Next(context.TODO()) {
		var doc bson.M

		cursor.Decode(&doc)

		//Just in case, to prevent unwanted crashes
		rawId, ok := doc["_id"]
		if !ok {
			continue
		}

		if _, ok := rawId.(int32); ok {
			id = int64(rawId.(int32))
		} else {
			id = rawId.(int64)
		}

		if isText {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{ReplyMarkup: markup},
			)
			if err != nil {
				fmt.Println(err)
				failed += 1
			} else {
				sent += 1
			}
		} else if isMedia {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{Caption: &caption, ParseMode: "HTML", ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		}

		total += 1

		_, _, er := stat.EditText(
			bot,
			fmt.Sprintf(`
<u>Live Broadcast Stats :</u>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
			&gotgbot.EditMessageTextOpts{ParseMode: "HTML"},
		)

		if er != nil {
			fmt.Println(er)
		}
	}

	stat.EditText(
		bot,
		fmt.Sprintf(`
<b><u>Broadcast Completed :</u></b>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
		&gotgbot.EditMessageTextOpts{ParseMode: "HTML"},
	)

	return nil
}

func ConCast(bot *gotgbot.Bot, ctx *ext.Context) error {
	//Function to handle the /concast command
	if !IsAdmin(ctx.Message.From.Id) {
		ctx.Message.Reply(bot, "<b>Sorry Thats An Owner Only Command :(</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyParameters: &gotgbot.ReplyParameters{MessageId: ctx.Message.MessageId}})
		return nil
	}

	cursor, _ := DB.Ucol.Find(context.TODO(), bson.D{{Key: "connected", Value: bson.D{{Key: "$exists", Value: true}}}})
	update := ctx.Message

	if update.ReplyToMessage == nil {
		update.Reply(bot, "Please reply this command to the message you would like to broadcast !", &gotgbot.SendMessageOpts{})
		return nil
	}

	msg := update.ReplyToMessage

	var isText bool = false
	var isMedia bool = false
	var caption string
	var markup *gotgbot.InlineKeyboardMarkup = &gotgbot.InlineKeyboardMarkup{}
	var total int
	var sent int
	var failed int
	var id int64

	if msg.Text != "" {
		isText = true
	} else {
		isMedia = true
		caption = msg.OriginalCaptionHTML()
	}

	if msg.ReplyMarkup != nil {
		markup = &gotgbot.InlineKeyboardMarkup{InlineKeyboard: msg.ReplyMarkup.InlineKeyboard}
	}

	stat, _ := update.Reply(bot, "<code>Starting concast ...</code>", &gotgbot.SendMessageOpts{ParseMode: "HTML"})

	for cursor.Next(context.TODO()) {
		var doc bson.M

		cursor.Decode(&doc)

		//Just in case, to prevent unwanted crashes
		rawId, ok := doc["_id"]
		if !ok {
			continue
		}

		if _, ok := rawId.(int32); ok {
			id = int64(rawId.(int32))
		} else {
			id = rawId.(int64)
		}

		if isText {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		} else if isMedia {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{Caption: &caption, ParseMode: "HTML", ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		}

		total += 1

		_, _, er := stat.EditText(
			bot,
			fmt.Sprintf(`
<u>Live Concast Stats :</u>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
			&gotgbot.EditMessageTextOpts{ParseMode: "HTML"},
		)

		if er != nil {
			fmt.Println(er)
		}
	}

	stat.EditText(
		bot,
		fmt.Sprintf(`
<b><u>Concast Completed :</u></b>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
		&gotgbot.EditMessageTextOpts{ParseMode: "HTML"},
	)
	return nil
}

func IsAdmin(user int64) bool {
	for _, admin := range Admins {
		if user == admin {
			return true
		}
	}
	return false
}
