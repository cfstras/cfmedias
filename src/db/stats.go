package db

import (
	log "logger"
)

//TODO cache these

func TitlesTotal() int64 {
	return selectInt(`select count(*) from ` + ItemTable)
}

func FoldersTotal() int64 {
	return selectInt(`select count(*) from ` + FolderTable)
}

func AvgFilesPerFolder() float32 {
	return float32(TitlesTotal()) / float32(FoldersTotal())
}

func selectInt(q string) int64 {
	if sum, err := dbmap.SelectInt(q); err != nil {
		log.Log.Println("query error:", err)
		return 0
	} else {
		return sum
	}
}
