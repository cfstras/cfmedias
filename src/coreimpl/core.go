package coreimpl

import (
	"config"
	"core"
	"db"
	"github.com/peterh/liner"
	log "logger"
	"os"
	"os/signal"
	"web"
)

type impl struct {
	shutUp         bool
	signals        chan os.Signal
	currentVersion version

	// stores the loaded commands, sorted by verb.
	// multiple verbs may point to the same command.
	commandMap map[string]core.Command

	// stores the loaded commands, sorted by first verb
	commandSet map[string]core.Command

	// the REPL state
	repl       *liner.State
	replActive bool
	reading    bool

	db *db.DB
}

func New() core.Core {
	return new(impl)
}

func (c *impl) Start() error {
	c.initVersion()

	// load config
	err := config.Load("config.json")
	if err != nil {
		return err
	}

	// set up signal handler
	c.signals = make(chan os.Signal, 2)
	signal.Notify(c.signals, os.Interrupt, os.Kill)
	go func() {
		for _ = range c.signals {
			// interrupted!
			c.Shutdown()
		}
	}()

	c.initCmd()

	c.db = new(db.DB)
	// connect to db
	if err = c.db.Open(c); err != nil {
		return err
	}
	c.shutUp = true

	//TODO call plugin loads

	// start web
	net := new(web.NetCmdLine)
	go net.Start(c)

	// update local files
	go c.db.Update()

	return nil
}

func (c *impl) Shutdown() error {
	if !c.shutUp {
		return nil
	}

	log.Log.Println("shutting down.")

	// disconnect from db
	if err := c.db.Close(); err != nil {
		log.Log.Println("Error closing database:", err)
	}

	// save config
	err := config.Save("config.json")
	if err != nil {
		//TODO don't catch if this is an init error
		log.Log.Println("Error while saving config:", err.Error())
	}
	c.shutUp = false

	if err := c.exitCmd(); err != nil {
		log.Log.Println("cmd exit error", err)
	}

	return err
}
