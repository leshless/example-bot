package text

import (
	"encoding/json"
	"io"
	"os"
)

// Global messages map
var Messages map[string]string

// Function that loads static ui data from /static/json/ui.json
func Load() error{
	path := "./static/json/text.json"
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

	err = json.Unmarshal(bytes, &Messages)
	if err != nil{
		return err
	}

	return nil
}
