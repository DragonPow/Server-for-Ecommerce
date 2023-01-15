package producer

type ProducerEvent struct {
	Key   string
	Value any
}

type InsertDatabaseEventValue int64 // id of new model
