package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type configuration struct {
	DbFile string
}

var config *configuration

// Initializes the config with standard parameters
func Init() {
	log.Println("Initializing new configuation.")
	config = &configuration{
		DbFile: "db.sqlite"}
}

func Load(configFile string) error {
	if config != nil {
		return errors.New("Config already loaded.")
	}

	file, err := os.Open(configFile)
	if os.IsNotExist(err) {
		Init()
		return Save(configFile)
	}
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

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	str, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(str); err != nil {
		return err
	}

	return nil
}
