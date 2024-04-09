package client

import (
	st "bot/pkg/storage"
	sm "bot/pkg/syncmap"
	"bot/pkg/ui"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// The following types represent the database tables.
// Allowed field types - int64 or string
// The "id" field is always a primary key for all tables.
// Нужно еще позже добавить метатэги для того, чтобы указывать параметры полей... (Или вообще для их названий)
type User struct{
	Id int64
	Name string
	Data string
}

// Data structures
// During the runtime the state must remain similar to db tables. (We always update these maps and db at the same time.)
var Users sm.SyncMap[int64, User]

var bot *tg.Bot

// Initialization
func Load() error{
	err := st.TouchTable(User{})
	if err != nil{
		return err
	}

	Users.Init()
	users, err := st.SelectAll[User]()
	if err != nil{
		return err
	}

	for _, user := range users{
		Users.Set(user.Id, user)
	}

	return nil
}

func SetBot(tgbot *tg.Bot){
	bot = tgbot
}

// Utility functions

func SendMessage(id int64, text string){
	sendparams := &tg.SendMessageParams{
		ChatID: tu.ID(id),
		Text: text,
		ParseMode: "HTML",
	}

	bot.SendMessage(sendparams)
}

func SendMessageWithMarkup(id int64, text string, menu *ui.Menu){
	sendparams := &tg.SendMessageParams{
		ChatID: tu.ID(id),
		Text: text,
		ReplyMarkup: menu.Build(),
		ParseMode: "HTML",
	}

	bot.SendMessage(sendparams)
}

func EditMessageText(message *tg.Message, text string){
	editparams := &tg.EditMessageTextParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
		Text: text,
		ParseMode: "HTML",
	}

	bot.EditMessageText(editparams)
}

func EditMessageMarkup(message *tg.Message, menu *ui.Menu){
	editparams := &tg.EditMessageReplyMarkupParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
		ReplyMarkup: menu.Build(),
	}

	bot.EditMessageReplyMarkup(editparams)
}

func DeleteMessage(message *tg.Message) {
	deleteparams := &tg.DeleteMessageParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
	}

	bot.DeleteMessage(deleteparams)
}


// Base handlers + middleware

func CommandHandler(bot *tg.Bot, message tg.Message){

}

func MessageHandler(bot *tg.Bot, message tg.Message){
	id := message.From.ID
	
	SendMessageWithMarkup(id, "Hello World", ui.Menus["main"])
}

func QueryHandler(bot *tg.Bot, query tg.CallbackQuery){

}	


