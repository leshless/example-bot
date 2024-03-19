package ui

import (
	"encoding/json"
	"io"
	"os"

	tg "github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Inline menu button
type Button struct{
	Label string `json:"label"` // Text to show on button
	Data string `json:"data"` // Callback data. Also defines handler for button action. May include additional information.
}

func ButtonNew(label string, data string) Button{
	return Button{
		Label: label,
		Data: data,
	}
}

// Inline menu row
type Row []Button

func RowNew() Row{
	return Row{}
}

func (row *Row) AddButton(button Button){
	*row = append(*row, button)
}

// Inline menu 
type Menu struct{
	Name string `json:"name"` // Name of menu. Should not differ from the string key in global map.
	Rows []Row `json:"rows"` // Menu rows of buttons.
}

func MenuNew(name string) Menu{
	return Menu{
		Name: name,
		Rows: []Row{},
	}
}

func (menu *Menu) AddRow(row Row){
	menu.Rows = append(menu.Rows, row)
}

// Build the abstract menu into telego.InlineKeyboardMarkup object.
func (menu Menu) Build() *tg.InlineKeyboardMarkup{
	inlineRows := [][]tg.InlineKeyboardButton{}
	
	for _, row := range menu.Rows{
		inlineRow := []tg.InlineKeyboardButton{}
		for _, button := range row{
			inlineButton := tu.InlineKeyboardButton(button.Label).WithCallbackData(button.Data)
			inlineRow = append(inlineRow, inlineButton)
		}
		inlineRows = append(inlineRows, inlineRow)
	}

	return tu.InlineKeyboard(inlineRows...)
}

// Global menu map
var Menus map[string]*Menu

// Function that loads static ui data from /static/json/ui.json
func Load() error{
	path := "./static/json/ui.json"
	var file *os.File

	// create .json file if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	}else{
		file, err = os.Open(path)
		if err != nil{
			return err
		}
	}

	bytes, err := io.ReadAll(file)
	if err != nil{
		return err
	}

	err = file.Close()
	if err != nil{
		return err
	}

	err = json.Unmarshal(bytes, &Menus)
	if err != nil{
		return err
	}

	return nil
}


// Nobody uses keyboard menus, so their implementation is very unlikely to appear there...