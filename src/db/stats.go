package db

import (
	"core"
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

func (d *DB) TitlesTotal() (num int64) {
	d.db.Table(ItemTable).Count(&num)
	return
}

func (d *DB) FoldersTotal() (num int64) {
	d.db.Table(FolderTable).Count(&num)
	return
}

func (d *DB) AvgFilesPerFolder() float32 {
	return float32(d.TitlesTotal()) / float32(d.FoldersTotal())
}
