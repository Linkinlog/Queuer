package internal

type dbParams struct {
	dsn string
}

type QueueEvent struct {
	EventID   int
	Data      string
	Resource  string
	Processed bool
}

type QueueReader interface {
	Read(resource string) ([]*QueueEvent, error)
	MarkProcessed(eventID int) error
}

func OpenQueue(d *dbParams) (QueueReader, error) {
	return NewPsqlConnector(d.dsn)
}

type TargetWriter interface {
	Write(interface{}) error
}

func OpenTarget(d *dbParams) (TargetWriter, error) {
	return nil, nil
}
