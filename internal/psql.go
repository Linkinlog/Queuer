package internal

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type psqlConnector struct {
	conn *pgx.Conn
}

func NewPsqlConnector(dsn string) (*psqlConnector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	return &psqlConnector{conn: conn}, nil
}

const readQuery = `SELECT
	events.id as id, data, resources.name as name, processed
	FROM events
	JOIN resources ON events.resource_id = resources.id
	WHERE processed = false
	AND resources.name = $1
	ORDER BY events.id ASC;`

func (q *psqlConnector) Read(resource string) ([]*QueueEvent, error) {
	rows, err := q.conn.Query(context.Background(), readQuery, resource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*QueueEvent, 0)
	for rows.Next() {
		var e QueueEvent
		err := rows.Scan(&e.EventID, &e.Data, &e.Resource, &e.Processed)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}

const markProcessedQuery = `UPDATE events SET processed = true WHERE id = $1;`

func (q *psqlConnector) MarkProcessed(eventID int) error {
	_, err := q.conn.Exec(context.Background(), markProcessedQuery, eventID)
	return err
}
