package db

type QueueEvent struct {
	EventID   int
	Data      string
	Resource  string
	Processed bool
}

type QueueReader interface {
	Read(resource string) ([]*QueueEvent, error)
	MarkProcessed(eventID int) error
	Close() error
}

func OpenQueue(dsn string) (QueueReader, error) {
	return NewPsqlConnector(dsn)
}

type TargetWriter interface {
	Write(eventId int, data []byte) error
	Close() error
}

func OpenTarget(dsn string) (TargetWriter, error) {
	return NewPsqlConnector(dsn)
}

type LogWriter interface {
	WriteLog(data []byte) error
	Close() error
}

func OpenLog(dsn string) (LogWriter, error) {
	return NewPsqlConnector(dsn)
}
