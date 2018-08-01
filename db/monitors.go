package db

import (
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// Monitor is a known monitor
type Monitor struct {
	ID     string
	Name   string
	Type   string
	Status string
}

var monitorBucket = []byte("monitors")
var db *bolt.DB

// Init sets up boltdb instance
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(monitorBucket)
		return err
	})
}

// UpdateMonitor upserts a monitor to DB
func UpdateMonitor(m *Monitor) (string, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(monitorBucket)
		// Convert ID to []byte for storage
		key := []byte(m.ID)
		// Convert m to []byte for storage
		buf, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return b.Put(key, buf)
	})
	if err != nil {
		return "", err
	}
	return m.ID, nil
}

// DeleteMonitor removes monitor from DB
func DeleteMonitor(id string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(monitorBucket)
		return b.Delete([]byte(id))
	})
	return err
}
