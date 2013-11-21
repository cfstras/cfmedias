package config

import (
	"encoding/json"
	"errors"
	"os"
)

type configuration struct {
	dbFile string
}

var config *configuration

func Load(configFile string) error {
	if config != nil {
		return errors.New("Config already loaded.")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var read configuration
	err = decoder.Decode(&read)
	if err != nil {
		return err
	}

	config = &read
	return nil
}
