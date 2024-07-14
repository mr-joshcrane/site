package backend_test

import (
	"testing"

	"github.com/mr-joshcrane/site/backend"
	"github.com/mr-joshcrane/site/store"
)

func TestWorkers(t *testing.T) {
	t.Parallel()
	// path := t.TempDir() + "/test.db"
	// s, err := store.NewSQLiteStore(path)
	s, err := store.NewSQLiteStore("test.db")
	if err != nil {
		t.Fatal(err)
	}
	err = backend.Workers(s)
	if err != nil {
		t.Fatal(err)
	}
}
