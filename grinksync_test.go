package grinksync_test

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/playtechnique/grinksync"
)

func TestOpeningDatabaseByVerifyingVersion(t *testing.T) {
	d := t.TempDir() + "/db.sqlite"
	db, err := sql.Open("sqlite3", d)
	db.Close()

	if err != nil {
		t.Fatal(err)
	}

	expected := "3.45.1"
	received := new(bytes.Buffer)

	g, err := grinksync.NewGrafanaDbSubscriber(d, "TestOpeningDatabaseByVerifyingVersion", received)

	if err != nil {
		t.Fatal(err)
	}

	if g.DbVersion != expected {
		t.Fatalf(cmp.Diff(expected, g.DbVersion), expected, g.DbVersion)
	}
}

func TestSubscribingToDatabaseUpdates(t *testing.T) {
	d := t.TempDir() + "/db.sqlite"
	db, err := sql.Open("sqlite3", d)

	if err != nil {
		t.Fatal(err)
	}

	db.Close()

	received := new(bytes.Buffer)
	_, err = grinksync.NewGrafanaDbSubscriber(d, "TestSubscribingToDatabaseUpdates", received)

	if err != nil {
		t.Fatal(err)
	}

	db, err = sql.Open("sqlite3", d)

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// _, err = db.Exec("CREATE TABLE foo (id integer not null primary key, content text);")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	if received.String() != "roflcopter" {
		t.Fatal(cmp.Diff(received.String(), "roflcopter"))
	}

}
