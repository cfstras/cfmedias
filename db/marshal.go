package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

func CustomMarshal(item interface{}) ([]byte, error) {
	res := make(map[string]interface{})
	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		switch field.Type.Name() {
		case "NullString":
			res[field.Name] = Str(v.Field(i).Interface().(sql.NullString))
		case "NullInt64":
			res[field.Name] = I64(v.Field(i).Interface().(sql.NullInt64))
		default:
			res[field.Name] = v.Field(i).Interface()
		}
	}
	return json.Marshal(&res)
}

func (item Item) MarshalJSON() ([]byte, error) {
	return CustomMarshal(item)
}

func (folder Folder) MarshalJSON() ([]byte, error) {
	return CustomMarshal(folder)
}

func NullStr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *v}
}

func NullI64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Valid: true, Int64: *v}
}

func Str(v sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}

func I64(v sql.NullInt64) *int64 {
	if v.Valid {
		return &v.Int64
	}
	return nil
}

func StrPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func StrStr(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func SqlStr(v string) sql.NullString {
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{Valid: true, String: v}
}
