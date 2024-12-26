package plugins

import (
	"github.com/Jisin0/Go-Filter-Bot/utils/autodelete"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// MyChatMember handles the MyChatMember event and saves the chat id in db.
func MyChatMember(bot *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.MyChatMember == nil || ctx.MyChatMember.NewChatMember.GetUser().Id != bot.Id {
		// new chat member not present or isnt the bot
		return nil
	}

	if chatType := ctx.MyChatMember.Chat.Type; chatType != gotgbot.ChatTypeSupergroup && chatType != gotgbot.ChatTypeGroup {
		return nil
	}

	DB.SetDefaultSettings(ctx.MyChatMember.Chat.Id)

	text := `<b><i>👋ᴛʜᴀɴᴋ ʏᴏᴜ ғᴏʀ ᴀᴅᴅɪɴɢ ᴍᴇ ᴛᴏ ʏᴏᴜʀ ɢʀᴏᴜᴘ.
I ᴄᴀɴ'ᴛ ᴡᴀɪᴛ ᴛᴏ sᴛᴀʀᴛ ʜᴇʟᴘɪɴɢ ʏᴏᴜ ᴏᴜᴛ, ᴍᴀᴋᴇ sᴜʀᴇ ʏᴏᴜ'ᴠᴇ ᴍᴀᴅᴇ ᴍᴇ ᴀɴ ᴀᴅᴍɪɴ ғɪʀsᴛ !</i></b>`

	m, err := bot.SendMessage(ctx.MyChatMember.Chat.Id, text, &gotgbot.SendMessageOpts{ParseMode: "HTML"})
	if err == nil {
		autodelete.InsertAutodel(autodelete.AutodelData{ChatID: m.Chat.Id, MessageID: m.MessageId}, 240) // autodel in 2 mins
	}

	return nil
}
