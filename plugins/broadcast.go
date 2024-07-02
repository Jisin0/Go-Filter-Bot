// (c) Jisin0

package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson"
)

func Broadcast(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Function to handle /broadcast command
	if !utils.IsAdmin(ctx.EffectiveUser.Id) {
		ctx.Message.Reply(bot, "<b>Only bot admins can use this command !</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyParameters: &gotgbot.ReplyParameters{MessageId: ctx.Message.MessageId}})
		return nil
	}

	var (
		update = ctx.Message
		filter = bson.D{}
	)

	if strings.HasPrefix(update.Text, "/concast") {
		filter = bson.D{{Key: "connected", Value: bson.D{{Key: "$exists", Value: true}}}}
	}

	cursor, err := DB.Ucol.Find(context.TODO(), filter)
	if err != nil {
		fmt.Printf("broadcast.find: %v\n", err)
		return nil
	}

	if update.ReplyToMessage == nil {
		update.Reply(bot, "Please reply this command to the message you would like to broadcast !", &gotgbot.SendMessageOpts{})
		return nil
	}

	msg := update.ReplyToMessage

	var (
		isText  = false
		isMedia = false
		caption string
		markup  = &gotgbot.InlineKeyboardMarkup{}
		total   int
		sent    int
		failed  int
		id      int64
	)

	if msg.Text != "" {
		isText = true
	} else {
		isMedia = true
		caption = msg.OriginalCaptionHTML()
	}

	if msg.ReplyMarkup != nil {
		markup = &gotgbot.InlineKeyboardMarkup{InlineKeyboard: msg.ReplyMarkup.InlineKeyboard}
	}

	stat, err := update.Reply(bot, "<code>Starting broadcast ...</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		fmt.Printf("broadcast.startmsg: %v\n", err)
		return nil
	}

	// TODO: create User type and remove this bson.M trash.
	for cursor.Next(context.TODO()) {
		var doc bson.M

		cursor.Decode(&doc)

		// Just in case, to prevent unwanted crashes
		rawID, ok := doc["_id"]
		if !ok {
			continue
		}

		if _, ok := rawID.(int32); ok {
			id = int64(rawID.(int32))
		} else {
			id = rawID.(int64)
		}

		if isText {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{ReplyMarkup: markup},
			)
			if err != nil {
				fmt.Printf("broadcast.sendmsg: %v\n", err)

				failed += 1
			} else {
				sent += 1
			}
		} else if isMedia {
			_, err := msg.Copy(
				bot,
				id,
				&gotgbot.CopyMessageOpts{Caption: &caption, ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: markup},
			)
			if err != nil {
				failed += 1
			} else {
				sent += 1
			}
		}

		total += 1

		// Update stat message every 20 requests to prevent floodwaits
		if total%20 == 0 {
			_, _, er := stat.EditText(
				bot,
				fmt.Sprintf(`
<u>Live Broadcast Stats :</u>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
				&gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML},
			)
			if er != nil {
				fmt.Printf("broadcast.progress.edit: %v\n", err)
			}
		}
	}

	_, _, err = stat.EditText(
		bot,
		fmt.Sprintf(`
<b><u>Broadcast Completed :</u></b>

Success : %v
Failed  : %v
Total   : %v
`, sent, failed, total),
		&gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML},
	)
	if err != nil {
		fmt.Printf("broadcast.progress.complete: %v\n", err)
	}

	return nil
}
