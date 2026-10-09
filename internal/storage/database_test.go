package storage

import (
	"bytes"
	"path/filepath"
	"testing"

	"go.etcd.io/bbolt"
)

func TestDatabaseSurvivesCloseAndReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "completion.db")
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte("persisted-value")
	if err := database.Update(func(transaction *bbolt.Tx) error {
		bucket, err := transaction.CreateBucketIfNotExists([]byte("example"))
		if err != nil {
			return err
		}
		return bucket.Put([]byte("key"), want)
	}); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}

	database, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()

	if err := database.View(func(transaction *bbolt.Tx) error {
		got := transaction.Bucket([]byte("example")).Get([]byte("key"))
		if !bytes.Equal(got, want) {
			t.Fatalf("value = %q, want %q", got, want)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
