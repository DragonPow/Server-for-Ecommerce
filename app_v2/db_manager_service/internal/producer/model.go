package producer

type ProducerEvent struct {
	Key   string
	Value any
}

type InsertDatabaseEventValue int64 // id of new model

type UpdateDatabaseEventValue struct {
	Id       int64  `json:"id"`
	Variants []byte `json:"variants"`
}
