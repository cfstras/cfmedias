package db

import (
	"core"
	"fmt"
	"io"
	log "logger"
)

//TODO cache these

func (d *DB) initStats(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"stats"},
		"Prints some statistics about the database",
		core.AuthUser,
		func(_ core.ArgMap, w io.Writer) error {
			fmt.Fprintf(w, " %7s %7s %7s\n", "Titles", "Folders", "Titles/Folder")
			fmt.Fprintf(w, " %7d %7d %7f\n", d.TitlesTotal(), d.FoldersTotal(),
				d.AvgFilesPerFolder())
			return nil
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
