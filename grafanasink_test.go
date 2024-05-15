package grafanasink_test

import (
	"database/sql"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/playtechnique/grafanasink"
)

func TestOpeningDatabaseByVerifyingVersion(t *testing.T) {
	d := t.TempDir() + "/db.sqlite"
	db, err := sql.Open("sqlite3", d)
	db.Close()

	if err != nil {
		t.Fatal(err)
	}

	expected := "3.45.1"
	g, err := grafanasink.NewGrafanaDbSubscriber(d)

	if err != nil {
		t.Fatal(err)
	}

	if g.DbVersion != expected {
		t.Fatalf(cmp.Diff(expected, g.DbVersion), expected, g.DbVersion)
	}
}

// func TestSubscribingToStreamOfEvents(t *testing.T) {
// 	g, err := grafanasink.NewGrafanaDbSubscriber()

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	responseFunction := func() string {
// 		fmt.Println()
// 	}

// 	g.Subscribe(responseFunction())

// }
