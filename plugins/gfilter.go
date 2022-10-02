// (c) Jisin0

package plugins

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/database"
	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/Jisin0/Go-Filter-Bot/utils/customfilters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson"
)

var Admins []int64 = utils.GetAdmins()

func GFilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	var chat_id int64
	update := ctx.Message
	message_id := update.MessageId
	chat_type := update.Chat.Type

	if chat_type == "private" {
		var ok bool
		chat_id, ok = DB.GetConnection(update.From.Id)
		if !ok {
			return nil
		}
	} else if chat_type == "supergroup" || chat_type == "group" {
		chat_id = update.Chat.Id
	} else {
		return nil
	}
	res, e := DB.GetMfilters(globalNumber)
	stopped := DB.GetCachedSetting(chat_id).Stopped
	if e != nil {
		fmt.Println(e)
		return nil
	} else {

		for res.Next(context.TODO()) {
			var f bson.M
			res.Decode(&f)
			key := f["text"].(string)
			text := `(?i)( |^|[^\w])` + key + `( |$|[^\w])`
			pattern := regexp.MustCompile(text)
			m := pattern.FindStringSubmatch(update.Text)
			if len(m) > 0 {
				var isStopped bool = false
				for _, k := range stopped {
					if key == k {
						isStopped = true
					}
				}
				if isStopped {
					continue
				}
				var filter database.Filter
				res.Decode(&filter)
				sendFilter(filter, bot, update, chat_id, message_id)
			}
		}
	}

	return nil
}

func StartGlobal(bot *gotgbot.Bot, ctx *ext.Context) error {
	//Function to handle the startglobal command
	update := ctx.Message
	var c int64
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
					update.Reply(bot, fmt.Sprintf("Restarted Global Filter For <i>%v</i> Successfully !", key), &gotgbot.SendMessageOpts{ParseMode: "HTML"})
					return nil
				}
			}

			update.Reply(bot, fmt.Sprintf("You Havent Stopped Any Global Filter For %v :(", key), &gotgbot.SendMessageOpts{})
		}
	}
	return nil
}

func Gfilters(bot *gotgbot.Bot, ctx *ext.Context) error {
	//Function to handle /gfilters function
	text := DB.StringMfilter(globalNumber)

	ctx.Message.Reply(bot, "Aʟʟ ғɪʟᴛᴇʀs sᴀᴠᴇᴅ ғᴏʀ ɢʟᴏʙᴀʟ ᴜsᴀɢᴇ :\n"+text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	return nil
}
