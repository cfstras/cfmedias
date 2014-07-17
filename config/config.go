package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/cfstras/cfmedias/errrs"
)

type Configuration struct {
	DbFile            string
	MediaPath         string
	WebPort           uint
	CacheWebTemplates bool

	ListenedUpperThreshold float32
	ListenedLowerThreshold float32

	Plugins map[string]interface{}
}

var Current *Configuration
var Default *Configuration

// Initializes the config with standard parameters
func init() {
	log.Println("Initializing new configuation.")
	Default = &Configuration{
		DbFile:                 "db.sqlite",
		MediaPath:              "~/Music",
		WebPort:                38888,
		CacheWebTemplates:      true,
		ListenedUpperThreshold: 0.7,
		ListenedLowerThreshold: 0.3,

		Plugins: map[string]interface{}{},
	}
	Current = Default
}

func RegisterPlugin(pluginName string, defaults interface{}) {
	Default.Plugins[pluginName] = defaults
	if _, ok := Current.Plugins[pluginName]; !ok {
		Current.Plugins[pluginName] = defaults
	}
}

func Load(configFile string) error {
	if Current != nil && Current != Default {
		return errrs.New("Config already loaded.")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// save defaults
		Current = Default
		return Save(configFile)
	}

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	var read Configuration
	err = json.Unmarshal(b, &read)
	if err != nil {
		return err
	}

	Current = &read

	//TODO ensure minimal config is done
	return nil
}

func Save(configFile string) error {
	if Current == nil {
		return errrs.New("Config has not been initialized, cannot save!")
	}

	str, err := json.MarshalIndent(Current, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configFile, str, 0644)

	if err != nil {
		return err
	}

	return nil
}
