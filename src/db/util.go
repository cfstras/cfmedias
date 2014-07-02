package db

import (
	"database/sql"
)

func NullStr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *v}
}

func Str(v string) sql.NullString {
	return sql.NullString{Valid: true, String: v}
}

func NullI64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: *v}
}
