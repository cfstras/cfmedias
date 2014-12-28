package coreimpl

import (
	"os"
	"os/signal"
	"sync"
	"time"

	as "github.com/cfstras/cfmedias/audioscrobbler"
	"github.com/cfstras/cfmedias/config"
	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/db"
	"github.com/cfstras/cfmedias/ipod"
	log "github.com/cfstras/cfmedias/logger"
	filesync "github.com/cfstras/cfmedias/sync"
	"github.com/cfstras/cfmedias/web"
	"github.com/peterh/liner"
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

	jobMutex sync.Mutex
	jobs     []chan core.JobSignal

	db             *db.DB
	audioscrobbler *as.AS
	sync           *filesync.Sync
	ipod           *ipod.IPod
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
	go net.Start(c, c.db)

	// start audioscrobbler
	c.audioscrobbler = new(as.AS)
	c.audioscrobbler.Start(c, c.db)

	// start syncer
	c.sync = new(filesync.Sync)
	c.sync.Start(c, c.db)

	c.ipod = new(ipod.IPod)
	c.ipod.Start(c, c.db, c.sync)

	// update local files
	go c.db.Update()

	return nil
}

func (c *impl) Shutdown() error {
	if !c.shutUp {
		return nil
	}

	log.Log.Println("Shutting down.")

	log.Log.Println("Stopping jobs...")
	wait := false
	for _, j := range c.jobs {
		wait = true
		go func(j chan<- core.JobSignal) {
			j <- core.SignalTerminate
		}(j)
	}
	if wait {
		time.Sleep(time.Second * 3)
	}
	log.Log.Println("Killing jobs...")
	for _, j := range c.jobs {
		j <- core.SignalKill
		close(j)
	}

	log.Log.Println("Closing database...")
	if err := c.db.Close(); err != nil {
		log.Log.Println("Error closing database:", err)
	}

	log.Log.Println("Saving config...")
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

func (c *impl) RegisterJob() <-chan core.JobSignal {
	c.jobMutex.Lock()
	defer c.jobMutex.Unlock()
	ch := make(chan core.JobSignal)
	c.jobs = append(c.jobs, ch)
	return ch
}

func (c *impl) UnregisterJob(job <-chan core.JobSignal) {
	go func() {
		// drain the channel
		for _ = range job {
		}
	}()
	c.jobMutex.Lock()
	defer c.jobMutex.Unlock()
	for i, v := range c.jobs {
		if v == job {
			// delete from slice
			c.jobs[i], c.jobs[len(c.jobs)-1], c.jobs = c.jobs[len(c.jobs)-1],
				nil, c.jobs[:len(c.jobs)-1]
			break
		}
	}
}
