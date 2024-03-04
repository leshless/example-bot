package client

import (
	st "bot/pkg/storage"
	ctx "context"
	"log"

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

	err = st.Insert(User{100, "Gleb", "Gleb"})
	if err != nil{
		return err
	}

	return nil
}

func MessageHandler(context ctx.Context, bot *tg.Bot, message tg.Message){
	
}

func QueryHandler(context ctx.Context, bot *tg.Bot, query tg.InlineQuery){

}	

func DeleteMessage(bot *tg.Bot, message *tg.Message) {
	deleteparams := &tg.DeleteMessageParams{
		ChatID:    tu.ID(message.Chat.ID),
		MessageID: message.MessageID,
	}

	bot.DeleteMessage(deleteparams)
}

func UnlockInlineButtons(bot *tg.Bot, message *tg.Message) {
	if message.Caption != "" {
		text := message.Caption
		editparams := &tg.EditMessageCaptionParams{
			ChatID:      tu.ID(message.Chat.ID),
			MessageID:   message.MessageID,
			ReplyMarkup: message.ReplyMarkup,
		}
		bot.EditMessageCaption(editparams.WithCaption(text + "　"))
		bot.EditMessageCaption(editparams.WithCaption(text))
	} else {
		text := message.Text
		editparams := &tg.EditMessageTextParams{
			ChatID:      tu.ID(message.Chat.ID),
			MessageID:   message.MessageID,
			ReplyMarkup: message.ReplyMarkup,
		}
		bot.EditMessageText(editparams.WithText(text + "　"))
		bot.EditMessageText(editparams.WithText(text))
	}
}

func RetrieveMessagePhoto(bot *tg.Bot, message tg.Message) ([]byte, bool) {
	if message.Photo == nil {
		return nil, false
	}
	file, err := bot.GetFile(&tg.GetFileParams{
		FileID: message.Photo[len(message.Photo)-1].FileID,
	})
	if err != nil {
		log.Print(err)
		return nil, false
	}

	photo, err := tu.DownloadFile(bot.FileDownloadURL(file.FilePath))
	if err != nil {
		log.Print(err)
		return nil, false
	}

	return photo, true
}