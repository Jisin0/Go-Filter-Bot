// (c) Jisin0

package config

import "github.com/PaulSonOfLars/gotgbot/v2"

var TEXT map[string]string = map[string]string{

	"START": `
<b>Hᴇʏ %v ɪᴍ %v ᴀɴ Aᴡᴇsᴏᴍᴇ Filter bot with global filter support</b>

<i>I can save a custom reply for a word in any chat. Check my help menu for more details.</i>
	`,

	"ABOUT": `
<b>○ 𝖫𝖺𝗇𝗀𝗎𝖺𝗀𝖾 :</b> : <a href='https://go.dev'>GO</a>
<b>○ 𝖫𝗂𝖻𝗋𝖺𝗋𝗒</b> : <a href='github.com/PaulSonOfLars/gotgbot'>gotgbot</a>
<b>○ 𝖣𝖺𝗍𝖺𝖻𝖺𝗌𝖾</b> : <a href='mongodb.org'>mongoDB</a>
<b>○ 𝖲𝗎𝗉𝗉𝗈𝗋𝗍</b> : <a href='t.me/Jisin0'>Here</a>
	`,

	"MF": `
Mᴀɴᴜᴀʟ ғɪʟᴛᴇʀs ᴀʟʟᴏᴡ ʏᴏᴜ ᴛᴏ sᴀᴠᴇ ᴄᴜsᴛᴏᴍ ғɪʟᴛᴇʀs ᴏᴛʜᴇʀ ᴛʜᴀɴ ᴛʜᴇ ᴀᴜᴛᴏᴍᴀᴛɪᴄ ᴏɴᴇs. Fɪʟᴛᴇʀs ᴄᴀɴ ʙᴇ ᴏғ ᴛᴇxᴛ/ᴘʜᴏᴛᴏ/ᴅᴏᴄᴜᴍᴇɴᴛ/ᴀᴜᴅɪᴏ/ᴀɴɪᴍᴀᴛɪᴏɴ/ᴠɪᴅᴇᴏ .

<b><u>Nᴇᴡ ғɪʟᴛᴇʀ :</u></b>

<b>Fᴏʀᴍᴀᴛ</b>
  <code>/filter "keyword" text</code> or
  Rᴇᴘʟʏ ᴛᴏ ᴀ ᴍᴇssᴀɢᴇ -><code>/filter "keyword"</code>
<b>Usᴀɢᴇ</b>
  <code>/filter "hi" hello</code>
  [<code>hello</code>] -> Reply -> <code>/filter hi</code>

<b><u>Sᴛᴏᴘ ғɪʟᴛᴇʀ :</u></b>

<b>Fᴏʀᴍᴀᴛ</b>
  <code>/stop "keyword"</code>
<b>Usᴀɢᴇ</b>
  <code>/stop "hi"</code>

<b><u>Vɪᴇᴡ ғɪʟᴛᴇʀs :</u></b>
  <code>/filters</code>
`,

	"GF": `
<b>Gʟᴏʙᴀʟ ғɪʟᴛᴇʀs ᴀʀᴇ ᴍᴀɴᴜᴀʟ ғɪʟᴛᴇʀs sᴀᴠᴇᴅ ʙʏ ʙᴏᴛ ᴀᴅᴍɪɴs ᴛʜᴀᴛ ᴡᴏʀᴋ ɪɴ ᴀʟʟ ᴄʜᴀᴛs. Tʜᴇʏ ᴘʀᴏᴠɪᴅᴇ ʟᴀᴛᴇsᴛ ᴍᴏᴠɪᴇs ɪɴ ᴀ ᴇᴀsʏ ᴛᴏ ᴜsᴇ ғᴏʀᴍᴀᴛ.</b>

<b><u>Sᴛᴏᴘ ғɪʟᴛᴇʀ :</u></b>

<u>Fᴏʀᴍᴀᴛ</u>
<code>/stop "keyword"</code>
<u>Usᴀɢᴇ</u>
<code>/stop "et"</code>

<b><u>Vɪᴇᴡ ғɪʟᴛᴇʀs :</u></b>
/gfilters
`,
	"CONNECT": `
<b>Cᴏɴɴᴇᴄᴛɪᴏɴs ᴀʟʟᴏᴡ ʏᴏᴜ ᴛᴏ ᴍᴀɴᴀɢᴇ ʏᴏᴜʀ ɢʀᴏᴜᴘ ʜᴇʀᴇ ɪɴ ᴘᴍ ɪɴsᴛᴇᴀᴅ ᴏғ sᴇɴᴅɪɴɢ ᴛʜᴏsᴇ ᴄᴏᴍᴍᴀɴᴅs ᴘᴜʙʟɪᴄʟʏ ɪɴ ᴛʜᴇ ɢʀᴏᴜᴘ ⠘⁾</b>

<b><u>Cᴏɴɴᴇᴄᴛ :</u></b>
-> Fɪʀsᴛ ɢᴇᴛ ʏᴏᴜʀ ɢʀᴏᴜᴘ's ɪᴅ ʙʏ sᴇɴᴅɪɴɢ /id ɪɴ ʏᴏᴜʀ ɢʀᴏᴜᴘ
-> <code>/connect [group_id]</code>

<b><u>Dɪsᴄᴏɴɴᴇᴄᴛ :</u></b>
<code>/disconnect</code>
`,

	"BROADCAST": `
<b>The broadcast feature allows bot admins to broadcast a message to all of the bot's users.</b>

<I>Broadcasts are of two types :</i>
 ~ Full Broadcast - Broadcast to all of the bot users <code>/broadcast</code>
 ~ Concast - Broadcast to only users who are connected to a chat <code>/concast</code>
`,

	"HELP": `
<b>To know how to use my modules use the buttons below to get all my commands with usage examples !</b>
`,

	"BTN": `
Here's a format of how you can add buttons to filters.
Buttons are split into different rows by using a new line.

<b>URL Button</b>
  <code>[Button Text](url:https://example.com)</code>

<b>Alert Button</b>
  <code>[Button Text](alert:This is an alert)</code>
`,
}

var BUTTONS map[string][][]gotgbot.InlineKeyboardButton = map[string][][]gotgbot.InlineKeyboardButton{
	"START": {
		{
			{Text: "About", CallbackData: "edit(ABOUT)"},
			{Text: "Help", CallbackData: "edit(HELP)"},
		},
	},
	"ABOUT": {
		{
			{Text: "Home", CallbackData: "edit(START)"},
			{Text: "Stats", CallbackData: "stats"},
		}, {
			{Text: "Source 🔗", Url: "https://github.com/Jisin0/Go-Filter-Bot"},
		},
	},
	"STATS": {
		{
			{Text: "⬅ Back", CallbackData: "edit(ABOUT)"},
			{Text: "Refresh 🔄", CallbackData: "stats"},
		},
	},
	"HELP": {
		{{Text: "Fɪʟᴛᴇʀ", CallbackData: "edit(MF)"},
			{Text: "Gʟᴏʙᴀʟ", CallbackData: "edit(GF)"},
		}, {
			{Text: "Cᴏɴɴᴇᴄᴛ", CallbackData: "edit(CONNECT)"}, {Text: "Broadcast", CallbackData: "edit(BROADCAST)"},
		},
		{{Text: "Bᴀᴄᴋ ➔", CallbackData: "edit(START)"}},
	},
	"MF": {{
		{Text: "ʙᴜᴛᴛᴏɴs", CallbackData: "edit(BTN)"},
		{Text: "bᴀᴄᴋ ➔", CallbackData: "edit(HELP)"},
	}},
	"BTN": {{{Text: "bᴀᴄᴋ ➔", CallbackData: "edit(MF)"}}},
}
