package models

import (
	"alphabet_typer/utility"
	"database/sql"
)

type Record struct {
	UUID        string  `json:"uuid"`
	Record      float64 `json:"record"`
	Username    string  `json:"username"`
	CreatedDate string  `json:"created_date"`
}

type RecordDB struct {
	UUID        string
	Record      float64
	Username    string
	CreatedDate sql.NullString
}

func (dbV *RecordDB) GetRecord() (r Record) {
	r.UUID = dbV.UUID
	r.Record = dbV.Record
	r.Username = dbV.Username
	r.CreatedDate = utility.GetStringValue(dbV.CreatedDate)
	return r
}
