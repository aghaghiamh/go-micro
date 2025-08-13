package repository

import "time"

type LogEntry struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Data      string    `bson:"data"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}
