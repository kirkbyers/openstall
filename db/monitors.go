package db

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// Monitor is a known monitor
type Monitor struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
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

// AllMonitors gets and returns all monitors from BoltDB
func AllMonitors() ([]Monitor, error) {
	var monitors []Monitor
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(monitorBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var m Monitor
			d := json.NewDecoder(bytes.NewBuffer(v))
			if err := d.Decode(&m); err != nil {
				return err
			}
			monitors = append(monitors, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return monitors, nil
}
