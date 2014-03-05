package db

import (
	"core"
	"util"
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
		return core.ResultByError(core.ErrorNotImplemented)
	}
	if tags != nil {
		return core.ResultByError(core.ErrorNotImplemented)
	}
	res, err := db.dbmap.Select(Item{}, "select * from "+ItemTable)
	return core.Result{core.StatusOK, res, err}
}
