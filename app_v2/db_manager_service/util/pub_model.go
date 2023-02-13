package util

import "time"

type InsertDatabaseEventValue int64 // id of new model

type UpdateDatabaseEventValue struct {
	Id         int64     `json:"id"`
	TimeUpdate time.Time `json:"time_update"`
	Variants   []byte    `json:"variants"`
}

func (u UpdateDatabaseEventValue) GetVersion() string {
	return u.TimeUpdate.Format(time.RFC3339)
}
