package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/cfstras/cfmedias/logger"

	"github.com/cfstras/cfmedias/errrs"
	"github.com/mitchellh/mapstructure"
)

type Configuration struct {
	DbFile            string
	MediaPath         string
	WebPort           uint
	CacheWebTemplates bool

	ListenedUpperThreshold float32
	ListenedLowerThreshold float32

	Plugins map[string]interface{}

	pluginTypes map[string]interface{}
}

var (
	Current *Configuration
	Default *Configuration
	loaded  bool = false
)

// Initializes the config with standard parameters
func init() {
	log.Log.Println("config init")
	Default = &Configuration{
		DbFile:                 "db.sqlite",
		MediaPath:              "~/Music",
		WebPort:                38888,
		CacheWebTemplates:      true,
		ListenedUpperThreshold: 0.7,
		ListenedLowerThreshold: 0.3,

		Plugins:     map[string]interface{}{},
		pluginTypes: map[string]interface{}{},
	}
	Current = copyDefault()
}

func copyDefault() *Configuration {
	str, err := json.MarshalIndent(Default, "", "  ")
	if err != nil {
		log.Log.Fatalln(err)
	}
	var c2 Configuration
	err = json.Unmarshal(str, &c2)
	if err != nil {
		log.Log.Fatalln(err)
	}
	return &c2
}

func RegisterPlugin(pluginName string, defaults interface{}, emptyConfigPtr interface{}) {
	log.Log.Println("Plugin", pluginName, "loaded")
	Default.Plugins[pluginName] = defaults
	Default.pluginTypes[pluginName] = emptyConfigPtr

	if _, ok := Current.Plugins[pluginName]; !ok {
		log.Log.Println("Loading defaults for", pluginName)
		Current.Plugins[pluginName] = defaults
	} else {
		log.Log.Println("config for", pluginName, ":", Current.Plugins[pluginName])
	}
}

func Load(configFile string) error {
	if loaded {
		return errrs.New("Config already loaded.")
	}
	loaded = true

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Log.Println("Initializing new configuation.")
		return nil
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

	// convert loaded maps to their config structs
	for k := range read.Plugins {
		err := mapstructure.Decode(read.Plugins[k], Default.pluginTypes[k])
		if err != nil {
			log.Log.Fatalln("Error loading config for plugin", k, "-", err)
		}
		read.Plugins[k] = Default.pluginTypes[k]
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
