package producer

type ProducerEvent struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type InsertDatabaseEventValue int64 // id of new model

type UpdateDatabaseEventValue struct {
	Id       int64  `json:"id"`
	Variants string `json:"variants"`
}
