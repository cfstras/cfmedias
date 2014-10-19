package db

import (
	"strings"

	"github.com/cfstras/cfmedias/core"
	"github.com/cfstras/cfmedias/errrs"
	"github.com/cfstras/cfmedias/util"
)

func (db *DB) initList(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"list", "ls"},
		"Searches and lists items from the database",
		map[string]string{
			"text":  "(Optional) Search string for full text search",
			"query": "(Optional) Input for SQL-ish query",
			"tags":  "(Optional) Comma-separated list of tags to include",
		},
		core.AuthGuest,
		db.list})
}

func (db *DB) list(ctx core.CommandContext) core.Result {
	//TODO
	args := ctx.Args
	var err error
	text, err := util.GetArg(args, "text", false, err)
	query, err := util.GetArg(args, "query", false, err)
	tags, err := util.GetArg(args, "tags", false, err)
	if text != nil {
		return core.ResultByError(core.ErrorNotImplemented)
	}
	if query != nil {
		res, err := db.ListQuery(*query)
		if err != nil {
			return core.ResultByError(errrs.New(err.Error()))
		}
		return core.Result{core.StatusOK, res, err, false}
	}
	if tags != nil {
		return core.ResultByError(core.ErrorNotImplemented)
	}

	res, err := db.ListAll()
	if err != nil {
		return core.ResultByError(errrs.New(err.Error()))
	}
	return core.Result{core.StatusOK, res, err, false}
}

func (db *DB) ListQuery(query string) ([]interface{}, error) {
	if strings.ContainsAny(query, `";`) ||
		strings.Contains(strings.ToLower(query), "union") {
		return nil, core.ErrorInvalidQuery
	}
	res := make([]Item, 0)
	err := db.db.Where(query).Find(&res).Error
	for i := range res { //TODO do this with join
		db.db.Find(&res[i].Folder, res[i].FolderId)
	}

	return ItemToInterfaceSlice(res), err
}

func (db *DB) ListAll() ([]interface{}, error) {
	res := make([]Item, 0)
	err := db.db.Find(&res).Error
	for i := range res { //TODO do this with join
		db.db.Find(&res[i].Folder, res[i].FolderId)
	}

	return ItemToInterfaceSlice(res), err
}
