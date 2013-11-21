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

	//TODO ensure minimal config is done
	return nil
}

func Save(configFile string) error {
	if config == nil {
		return errors.New("Config has not been initialized, cannot save!")
	}

	file, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
