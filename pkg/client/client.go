package client

import (
	st "bot/pkg/storage"
	ui "bot/pkg/ui"

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

func Load() error{
	err := st.TouchTable(User{})
	if err != nil{
		return err
	}

	// err = st.Insert(User{100, "Gleb", "Gleb"})
	// if err != nil{
	// 	return err
	// }

	return nil
}

func CommandHandler(bot *tg.Bot, message tg.Message){

}

func MessageHandler(bot *tg.Bot, message tg.Message){
	SendMessage(bot, "HELLO WORLD!", message.From.ID)

}

func QueryHandler(bot *tg.Bot, query tg.InlineQuery){

}	

func SendMessage(bot *tg.Bot, text string, id int64){
	markup := ui.Menus["main"].Build()
	bot.SendMessage(tu.Message(tu.ID(id), text).WithReplyMarkup(markup))
}

func DeleteMessage(bot *tg.Bot, message *tg.Message) {
	deleteparams := &tg.DeleteMessageParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
	}

	bot.DeleteMessage(deleteparams)
}
