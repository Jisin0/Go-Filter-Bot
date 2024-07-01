// (c) Jisin0

package plugins

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jisin0/Go-Filter-Bot/database"
	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/Jisin0/Go-Filter-Bot/utils/customfilters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.mongodb.org/mongo-driver/bson"
)

// Regex expressions for callbacks and filtering
var buttonRegex *regexp.Regexp = regexp.MustCompile(`\[([^\[]+?)\]\((buttonurl|url|alert):(?:/{0,2})(.+?)\)`)
var parseRegex *regexp.Regexp = regexp.MustCompile(`^"([^"]+)"`)
var cbstopRegex *regexp.Regexp = regexp.MustCompile(`stopf\((.+)\)`)
var cbalertRegex *regexp.Regexp = regexp.MustCompile(`alert\((.+)\)`)

// Number used as chat id for global filters. you could change it to anything you like but you will lose any existing gfilters
var globalNumber int64 = 101

var DB *database.Database = database.NewDatabase()

// Manual filter function
func MFilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	var (
		chatID    int64
		update    = ctx.Message
		messageID = update.MessageId
		chatType  = update.Chat.Type
		message   = update.Text
	)

	if update.Caption != "" {
		message = update.Caption
	}

	switch chatType {
	case gotgbot.ChatTypePrivate:
		var ok bool
		chatID, ok = DB.GetConnection(update.Chat.Id)
		if !ok {
			return nil
		}
	case gotgbot.ChatTypeSupergroup, gotgbot.ChatTypeGroup:
		chatID = update.Chat.Id
	default:
		return nil
	}

	res, e := DB.GetMfilters(chatID)
	if e != nil {
		fmt.Printf("mfilter.getmfilters: %v\n", e)
		return nil
	} else {
		// Trying to find a regex match
		for res.Next(context.TODO()) {
			var f database.Filter

			res.Decode(&f)

			text := `(?i)( |^|[^\w])` + f.Text + `( |$|[^\w])`

			pattern := regexp.MustCompile(text)

			m := pattern.FindStringSubmatch(message)
			if len(m) > 0 {
				sendFilter(f, bot, update, chatID, messageID)
			}
		}
	}

	return nil
}

// Function to handle filter and gfilter commands
//
//nolint:errcheck
func NewFilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	var c int64
	update := ctx.Message

	// I didnt wanna create a whole new function for gfilter so ...
	if strings.HasPrefix(update.Text, "/gfilter") {
		for _, admin := range Admins {
			if update.From.Id == admin {
				c = globalNumber
			}
		}

		if c != globalNumber {
			update.Reply(bot, "Thats an Admin-only Command :(", &gotgbot.SendMessageOpts{})
			return nil
		}
	} else {
		// Verifying and getting connections for private chats
		var v bool
		c, v = customfilters.Verify(bot, ctx)
		if !v {
			return nil
		} else if c == 0 {
			c = ctx.Message.Chat.Id
		}
	}

	args := strings.SplitN(update.Text, " ", 3)
	if len(args) < 2 && (update.ReplyToMessage == nil && len(args) < 3) {
		update.Reply(
			bot,
			"Not Enough Parameters :(\n\nExample:- <code>/filter hi hello</code>",
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
		)

		return nil
	}

	parse := parseQuotes(strings.SplitN(update.Text, " ", 2)[1])

	e := DB.Mcol.FindOne(context.TODO(), bson.D{{Key: "group_id", Value: c}, {Key: "text", Value: parse[0]}})
	if e.Err() == nil {
		update.Reply(
			bot,
			fmt.Sprintf("It Looks Like Another Filter For <code>%v</code> Has Already Been Saved In This Chat, If You Want To Stop It First Use The Button Below", parse[0]),
			&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{
				{{Text: "Stop Filter", CallbackData: fmt.Sprintf("stopf(%v|%v|y)", parse[0], "local")}},
			}}},
		)
		return nil
	}

	var (
		text   = parse[1]
		button [][]map[string]string
	)

	if update.ReplyToMessage != nil {
		text += update.ReplyToMessage.Text
		text += update.ReplyToMessage.Caption

		if update.ReplyToMessage.ReplyMarkup != nil {
			// Buttons are converted to maps to save to the database
			button = buttonToMap(ctx.Message.ReplyToMessage.ReplyMarkup.InlineKeyboard)
		}
	}

	// Creating a random string with length 15
	uniqueID := utils.RandString(15)

	text, button, alert := parseButtons(text, uniqueID, button)

	// Finding media if any
	var (
		fileID    string
		mediaType string
	)

	if msg := update.ReplyToMessage; msg != nil {
		switch {
		case msg.Document != nil:
			fileID = msg.Document.FileId
			mediaType = "document"
		case msg.Video != nil:
			fileID = msg.Video.FileId
			mediaType = "video"
		case msg.Audio != nil:
			fileID = msg.Audio.FileId
			mediaType = "audio"
		case msg.Sticker != nil:
			fileID = msg.Sticker.FileId
			mediaType = "sticker"
		case msg.Animation != nil:
			fileID = msg.Animation.FileId
			mediaType = "animation"
		case msg.Photo != nil:
			fileID = msg.Photo[len(msg.Photo)-1].FileId
			mediaType = "photo"
		}
	}

	f := &database.Filter{
		Id:        uniqueID,
		ChatId:    c,
		Text:      parse[0],
		Content:   text,
		FileID:    fileID,
		Markup:    button,
		Alerts:    alert,
		Length:    len(parse[0]),
		MediaType: mediaType,
	}

	DB.SaveMfilter(f)

	_, err := update.Reply(
		bot,
		fmt.Sprintf("<i>Successfully Saved A Manual Filter For <code>%v</code> !</i>", parse[0]),
		&gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML},
	)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func CbStop(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Function to handle callbacks from confirm buttons when stopping a filter
	update := ctx.CallbackQuery

	c, v := customfilters.Verify(bot, ctx)
	if !v {
		return nil
	} else if c == 0 {
		c = update.Message.GetChat().Id
	}

	// Making sure the callback data is valid
	args := strings.Split(cbstopRegex.FindStringSubmatch(update.Data)[1], "|")
	if len(args) < 3 {
		update.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Bad Button :(", ShowAlert: true})
		return nil
	}

	var (
		key   = args[0]
		ftype = args[1]
		erase = args[2]
	)

	if ftype == "local" {
		if erase == "y" {
			DB.DeleteMfilter(c, key)
			update.Message.EditText(bot, fmt.Sprintf("Manual Filter For <code>%v</code> Was Deleted Successfully !", key), &gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML})
		} else if erase == "n" {
			// Unused till now, might use it later
			update.Message.EditText(bot, fmt.Sprintf(`Are You Sure You Want To Permanently Delete The Manual Filter For %v ?\nClick The "Yes I'm Sure" Button To Confirm `, key), &gotgbot.EditMessageTextOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "Ignore", CallbackData: "close"}, {Text: "Yes I'm Sure", CallbackData: fmt.Sprintf("stopf(%v|local|y)", key)}}}}, ParseMode: gotgbot.ParseModeHTML})
		}
	} else if ftype == "global" {
		for _, admin := range Admins {
			if update.From.Id == admin {
				DB.DeleteMfilter(globalNumber, key)
				update.Message.EditText(bot, fmt.Sprintf("Global Filter For <code>%v</code> Was Deleted Successfully !", key), &gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML})
				return nil
			}
		}

		DB.StopGfilter(c, key)

		update.Message.EditText(bot, fmt.Sprintf("Global Filter For <code>%v</code> Has Been Stopped Successfully !", key), &gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML})
	}

	return nil
}

// Function to handle the stop command
//
//nolint:errcheck // too many
func StopMfilter(bot *gotgbot.Bot, ctx *ext.Context) error {
	c, v := customfilters.Verify(bot, ctx)
	if !v {

		return nil
	} else if c == 0 {
		c = ctx.Message.Chat.Id
	}

	var (
		update = ctx.Message
		split  = strings.SplitN(update.Text, " ", 2)
		key    string
	)

	switch {
	case len(split) < 2 && update.Chat.Type == gotgbot.ChatTypePrivate:
		m := utils.Ask(bot, "Ok Now Send Me The Name OF The Filter You Would Like To Stop ...", ctx.EffectiveChat, ctx.EffectiveUser)
		if m == nil {
			return nil
		}

		key = m.Text
	case len(split) < 2:
		update.Reply(bot, "Whoops looks like you forgot to mention a filter to stop !", &gotgbot.SendMessageOpts{})
	default:
		key = split[1]
	}

	// Checking if theres a local/global filter for the key
	_, k := DB.GetMfilter(c, key)
	_, ok := DB.GetMfilter(globalNumber, key)

	// If there isnt local or global
	if !k && !ok {
		update.Reply(bot, fmt.Sprintf("I Couldnt Find Any Filter For <code>%v</code> To Stop :(", key), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return nil
	}

	if len(key) < 20 {
		// Both local and global
		if k && ok {
			markup := gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "Local", CallbackData: fmt.Sprintf("stopf(%v|local|y)", key)}, {Text: "Global", CallbackData: fmt.Sprintf("stopf(%v|global|y)", key)}}}}
			_, err := update.Reply(bot, "Please Select If You Would Like To Stop The Manual Filter (which you saved) Or Global Filter (saved by owners) For <code>"+key+"</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: markup})
			if err != nil {
				fmt.Println(err)
			}
		} else if k {
			// Only local
			markup := gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "CONFIRM", CallbackData: fmt.Sprintf("stopf(%v|local|y)", key)}}}}
			_, err := update.Reply(bot, "Please Press The Button Below To Confirm Deletion Of Manual Filter For <code>"+key+"</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: markup})
			if err != nil {
				fmt.Println(err)
			}
		} else if ok {
			// Only global
			markup := gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "CONFIRM", CallbackData: fmt.Sprintf("stopf(%v|global|y)", key)}}}}
			_, err := update.Reply(bot, "Please Press The Button Below To Stop Global Filter For <code>"+key+"</code>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML, ReplyMarkup: markup})
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		DB.DeleteMfilter(c, key)
		update.Reply(bot, fmt.Sprintf("Manual Filter For <i>%v</i> Was Stopped Successfully !", key), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	}

	return nil
}

func AllMfilters(bot *gotgbot.Bot, ctx *ext.Context) error {
	// Function to handle the /filters command
	update := ctx.Message
	c := update.Chat.Id

	if update.Chat.Type == gotgbot.ChatTypePrivate {
		if i, k := DB.GetConnection(update.From.Id); k {
			c = i
		}
	}

	text := DB.StringMfilter(c)

	_, err := update.Reply(bot, "Lɪsᴛ ᴏғ ғɪʟᴛᴇʀs ғᴏʀ ᴛʜɪs ᴄʜᴀᴛ :\n"+text, &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		fmt.Printf("allmfilters: %v\n", err)
	}

	return nil
}

// Function to handle callbacks from alert button in saved filters
//
//nolint:errcheck // unnecessary
func CbAlert(bot *gotgbot.Bot, ctx *ext.Context) error {
	update := ctx.CallbackQuery

	args := strings.Split(cbalertRegex.FindStringSubmatch(update.Data)[1], "|")
	if len(args) < 2 {
		update.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Bad Request :(", ShowAlert: true})
	} else {
		uniqueID := args[0]

		index, err := strconv.ParseInt(args[1], 0, 64)
		if err != nil {
			update.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Bad Request :(", ShowAlert: true})
		} else {
			text := DB.GetAlert(uniqueID, int(index))
			update.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: text, ShowAlert: true})
		}
	}

	return nil
}

func buttonToMap(btn [][]gotgbot.InlineKeyboardButton) [][]map[string]string {
	// Convert a button into a map
	var (
		totalButtons [][]map[string]string
		rowButtons   []map[string]string
	)

	for _, i := range btn {
		rowButtons = []map[string]string{}

		for _, j := range i {
			b := map[string]string{"Text": j.Text}

			if j.CallbackData != "" {
				b["CallbackData"] = j.CallbackData
			} else if j.Url != "" {
				b["Url"] = j.Url
			} else {
				continue
			}
			rowButtons = append(rowButtons, b)
		}
		totalButtons = append(totalButtons, rowButtons)
	}

	return totalButtons
}

func mapToButton(data [][]map[string]string) [][]gotgbot.InlineKeyboardButton {
	// Convert a map back into button
	var totalButtons [][]gotgbot.InlineKeyboardButton

	for _, i := range data {
		var rowButtons []gotgbot.InlineKeyboardButton

		for _, b := range i {
			text := b["Text"]

			if u, k := b["Url"]; k {
				rowButtons = append(rowButtons, gotgbot.InlineKeyboardButton{Text: text, Url: u})
			} else if c, k := b["CallbackData"]; k {
				rowButtons = append(rowButtons, gotgbot.InlineKeyboardButton{Text: text, CallbackData: c})
			}
		}

		totalButtons = append(totalButtons, rowButtons)
	}

	return totalButtons
}

func parseQuotes(text string) []string {
	res := parseRegex.FindStringSubmatch(text)
	if len(res) > 0 {
		return []string{res[1], strings.Replace(text, res[1], "", 1)}
	} else {
		split := strings.SplitN(text, " ", 2)
		return []string{split[0], strings.Replace(text, split[0], "", 1)}
	}
}

func parseButtons(text string, uniqueID string, totalButtons [][]map[string]string) (string, [][]map[string]string, []string) {
	var (
		returnText = text
		rowButtons []map[string]string
		alert      []string
	)

	for _, rows := range strings.Split(text, "\n") {
		for _, m := range buttonRegex.FindAllStringSubmatch(rows, 5) {
			if len(m) < 4 {
				continue
			}

			if m[2] == "url" || m[2] == "buttonurl" {
				rowButtons = append(rowButtons, map[string]string{"Text": m[1], "Url": m[3]})
			} else if m[2] == "alert" || m[2] == "buttonalert" {
				alert = append(alert, m[3])
				rowButtons = append(rowButtons, map[string]string{"Text": m[1], "CallbackData": fmt.Sprintf("alert(%v|%v)", uniqueID, len(alert)-1)})
			}

			returnText = strings.Replace(returnText, m[0], "", 1)
		}

		if len(rowButtons) > 0 {
			totalButtons = append(totalButtons, rowButtons)
		}

		rowButtons = []map[string]string{}
	}

	return strings.Trim(returnText, " "), totalButtons, alert
}

func sendFilter(f database.Filter, bot *gotgbot.Bot, update *gotgbot.Message, chatID int64, messageID int64) {
	// A function to send a filter if the regex matches
	// I just made this func bcuz i'd have to duplicate it for mfilter and gfilter

	// Find buttons saved for the filter and convert it from map
	var (
		buttons = mapToButton(f.Markup)
		markup  = gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons}
		content = f.Content
		err     error
	)

	mediaType := f.MediaType
	if mediaType != "" {
		fileID := f.FileID

		switch mediaType {
		case "document":
			_, err = bot.SendDocument(chatID, fileID, &gotgbot.SendDocumentOpts{Caption: content, ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
		case "sticker":
			_, err = bot.SendSticker(chatID, fileID, &gotgbot.SendStickerOpts{ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup})
		case "video":
			_, err = bot.SendVideo(chatID, fileID, &gotgbot.SendVideoOpts{Caption: content, ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
		case "photo":
			_, err = bot.SendPhoto(chatID, fileID, &gotgbot.SendPhotoOpts{Caption: content, ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
		case "audio":
			_, err = bot.SendAudio(chatID, fileID, &gotgbot.SendAudioOpts{Caption: content, ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
		case "animation":
			_, err = bot.SendAnimation(chatID, fileID, &gotgbot.SendAnimationOpts{Caption: content, ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
		default:
			fmt.Println("Unknown media type " + mediaType)

		}
	} else {
		_, err = update.Reply(bot, content, &gotgbot.SendMessageOpts{ReplyParameters: &gotgbot.ReplyParameters{MessageId: messageID}, ReplyMarkup: markup, ParseMode: gotgbot.ParseModeHTML})
	}

	if err != nil {
		fmt.Println(err)
		fmt.Println(f)
	}
}
