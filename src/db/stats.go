package db

import (
	"core"
	log "logger"
)

//TODO cache these

func (d *DB) initStats(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"stats"},
		"Prints some statistics about the database",
		map[string]string{},
		core.AuthUser,
		func(_ core.CommandContext) core.Result {
			res := map[string]interface{}{
				"items_total":      d.TitlesTotal(),
				"folders_total":    d.FoldersTotal(),
				"items_per_folder": d.AvgFilesPerFolder(),
			}
			//TODO more stats
			//TODO stats per user level
			return core.Result{Status: core.StatusOK, Result: res}
		}})
}

func (d *DB) TitlesTotal() int64 {
	return d.selectInt(`select count(*) from ` + ItemTable)
}

func (d *DB) FoldersTotal() int64 {
	return d.selectInt(`select count(*) from ` + FolderTable)
}

func (d *DB) AvgFilesPerFolder() float32 {
	return float32(d.TitlesTotal()) / float32(d.FoldersTotal())
}

func (d *DB) selectInt(q string) int64 {
	if sum, err := d.dbmap.SelectInt(q); err != nil {
		log.Log.Println("query error:", err)
		return 0
	} else {
		return sum
	}
}
