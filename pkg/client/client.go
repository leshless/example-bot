package client

import (
	st "bot/pkg/storage"
	sm "bot/pkg/syncmap"
	"bot/pkg/text"
	"bot/pkg/ui"
	"fmt"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

var bot *tg.Bot

// Data type representing the callback query hidden data.
// First - the name of coresponding handler to call or next inline menu to open (or both).
// Second - some additional data for handler to process.
// If second field is empty, it will be ommited.
type Callback struct{
	First string
	Second int64
}

func GetCallback(data string) (Callback, error){
	var d Callback
	_, err := fmt.Sscanf(data, "%v %v", &d.First, &d.Second)
	if err != nil{
		_, err = fmt.Sscanf(data, "%v", &d.First)
		if err != nil{
			return d, err
		}
	}

	return d, nil
}

func (d Callback) String() string{
	return fmt.Sprintf("%v_%v", d.First, d.Second)
}

// The following types represent the database tables.
// Allowed field types - int64 or string
// The "id" field is always a primary key for all tables.
// Нужно еще позже добавить метатэги для того, чтобы указывать параметры полей... (Или вообще для их названий)
type User struct{
	Id int64
	State string
}

// Data structures
// During the runtime the state must remain similar to db tables. (We always update these maps and db at the same time.)
var Users sm.SyncMap[int64, User]

func AddUser(id int64){
	user := User{
		id,
		"main",
	}

	Users.Set(id, user)
	st.Insert[User](user)
}

var (
	MessageHandlers map[string] func(tg.Message)
	CommandHandlers map[string] func(tg.Message)
	QueryHandlers map[string] func(tg.CallbackQuery)
)

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

func EditMessageText(message tg.Message, text string){
	editparams := &tg.EditMessageTextParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
		Text: text,
		ParseMode: "HTML",
	}

	bot.EditMessageText(editparams)
}

func EditMessageMarkup(message tg.Message, menu *ui.Menu){
	editparams := &tg.EditMessageReplyMarkupParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
		ReplyMarkup: menu.Build(),
	}

	bot.EditMessageReplyMarkup(editparams)
}

func DeleteMessage(message tg.Message) {
	deleteparams := &tg.DeleteMessageParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
	}

	bot.DeleteMessage(deleteparams)
}


// Base handlers + middleware

func CommandHandler(bot *tg.Bot, message tg.Message){
	id := message.From.ID
	command := message.Text
	handler, ok := CommandHandlers[command]
	
	if ok{
		handler(message)
	}else{
		SendMessage(id, text.Messages["command_not_found"])
	}
}

func MessageHandler(bot *tg.Bot, message tg.Message){
	id := message.From.ID
	user, _ := Users.Get(id)
	handler, ok := MessageHandlers[user.State]
	
	if ok{
		handler(message)
	}else{
		menu, ok := ui.Menus[user.State]

		if ok{
			SendMessageWithMarkup(id, "", menu)
		}else{
			SendMessage(id, text.Messages["message_not_found"])
		}
	}
}

func QueryHandler(bot *tg.Bot, query tg.CallbackQuery){
	callback, _ := GetCallback(query.Data)

	handler, ok := QueryHandlers[callback.First]
	if ok{
		handler(query)
	}

	menu, ok := ui.Menus[callback.First]
	if ok && query.Message.IsAccessible(){
		message := *query.Message.(*tg.Message)
		EditMessageMarkup(message, menu)
		// change state?
	}
}


