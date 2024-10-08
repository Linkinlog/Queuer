package internal

type dbParams struct {
	driver string
	dsn    string
}

type QueueEvent struct {
	sequence int
	payload  interface{}
}

type QueueReader interface {
	Read() (*QueueEvent, error)
}


func OpenQueue(d *dbParams) (QueueReader, error) {
	return nil, nil
}

type TargetWriter interface {
	Write(interface{}) error
}

func OpenTarget(d *dbParams) (TargetWriter, error) {
	return nil, nil
}
