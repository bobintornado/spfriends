// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"context"

	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"go.opencensus.io/trace"
)

type DB struct {
	connUrl string
	dbPool  bolt.ClosableDriverPool
}

// New returns a new DB neo4j pool.
func New(url string, max int) (*DB, error) {
	dbPool, err := bolt.NewClosableDriverPool(url, max)
	if err != nil {
		return nil, err
	}

	db := DB{
		connUrl: url,
		dbPool:  dbPool,
	}

	return &db, nil
}

// OpenPool returns a new pooled connection.
func (db *DB) OpenPool() (conn *bolt.Conn, err error) {
	c, err := db.dbPool.OpenPool()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *DB) Close() error {
	err := db.dbPool.Close()
	if err != nil {
		return err
	}

	return nil
}

// StatusCheck validates the DB status good.
func (db *DB) StatusCheck(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "platform.DB.StatusCheck")
	defer span.End()

	c, err := db.dbPool.OpenPool()
	if err != nil {
		return err
	}
	defer c.Close()

	return nil
}
