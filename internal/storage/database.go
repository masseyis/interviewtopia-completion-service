// Package storage opens the embedded database supplied with the starter. It
// deliberately does not define buckets, records or a domain repository: those
// persistence decisions are part of the exercise.
package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

// Open creates or opens a bbolt database at path. bbolt is embedded, so no
// database server or external installation is required.
func Open(path string) (*bbolt.DB, error) {
	if path == "" {
		return nil, errors.New("database path is empty")
	}

	directory := filepath.Dir(path)
	if directory != "." {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}

	database, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	return database, nil
}
