package db

import (
	"core"
	"errrs"
	"strings"
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
		res, err := db.listQuery(*query)
		if err != nil {
			return core.ResultByError(errrs.New(err.Error()))
		}
		return core.Result{core.StatusOK, res, err, false}
	}
	if tags != nil {
		return core.ResultByError(core.ErrorNotImplemented)
	}

	res := make([]interface{}, 0)
	err = db.db.Table(ItemTable).Find(res).Error

	return core.Result{core.StatusOK, res, err, false}
}

func (db *DB) listQuery(query string) ([]interface{}, error) {
	if strings.ContainsAny(query, `";`) ||
		strings.Contains(strings.ToLower(query), "union") {
		return nil, core.ErrorInvalidQuery
	}
	res := make([]interface{}, 0)
	err := db.db.Table(ItemTable).Where(query).Find(res).Error

	return res, err
}
