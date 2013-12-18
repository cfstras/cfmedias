package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Configuration struct {
	DbFile    string
	MediaPath string
	WebPort   uint
}

var Current *Configuration
var Default *Configuration

// Initializes the config with standard parameters
func Init() {
	log.Println("Initializing new configuation.")
	Default = &Configuration{
		DbFile:    "db.sqlite",
		MediaPath: "~/Music",
		WebPort:   38888,
	}
	Current = Default
}

func Load(configFile string) error {
	if Current != nil && Current != Default {
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
	var read Configuration
	err = decoder.Decode(&read)
	if err != nil {
		return err
	}

	Current = &read

	//TODO ensure minimal config is done
	return nil
}

func Save(configFile string) error {
	if Current == nil {
		return errors.New("Config has not been initialized, cannot save!")
	}

	file, err := os.Create(configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	str, err := json.MarshalIndent(Current, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(str); err != nil {
		return err
	}

	return nil
}
