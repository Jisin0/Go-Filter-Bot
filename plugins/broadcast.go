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
		ctx.Message.Reply(bot, "<b>Sorry Thats An Admin Only Command :(</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyToMessageId: ctx.Message.MessageId})
		return nil
	}

	cursor, _ := DB.Ucol.Find(context.TODO(), bson.D{{}})
	update := ctx.Message

	var isText bool = false
	var isMedia bool = false
	var caption string
	var markup *gotgbot.InlineKeyboardMarkup
	var total int
	var sent int
	var failed int
	var id int64

	if update.Text != "" {
		isText = true
	} else {
		isMedia = true
		caption = update.OriginalCaptionHTML()
	}

	markup = update.ReplyMarkup

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
			_, err := update.Copy(
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
			_, err := update.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{Caption: caption, ParseMode: "HTML", ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		}

		total += 1

		_, _, er := update.EditText(
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

	update.EditText(
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
		ctx.Message.Reply(bot, "<b>Sorry Thats An Owner Only Command :(</b>", &gotgbot.SendMessageOpts{ParseMode: "HTML", ReplyToMessageId: ctx.Message.MessageId})
		return nil
	}

	cursor, _ := DB.Ucol.Find(context.TODO(), bson.D{{Key: "connected", Value: bson.D{{Key: "$exists", Value: true}}}})
	update := ctx.Message

	var isText bool = false
	var isMedia bool = false
	var caption string
	var markup *gotgbot.InlineKeyboardMarkup
	var total int
	var sent int
	var failed int
	var id int64

	if update.Text != "" {
		isText = true
	} else {
		isMedia = true
		caption = update.OriginalCaptionHTML()
	}

	markup = update.ReplyMarkup

	for cursor.Next(context.TODO()) {
		var doc bson.M

		cursor.Decode(&doc)

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
			_, err := update.Copy(
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
			_, err := update.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{Caption: caption, ParseMode: "HTML", ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		}

		total += 1

		_, _, er := update.EditText(
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

	update.EditText(
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
