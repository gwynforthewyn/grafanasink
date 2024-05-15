package grafanasink

import (
	"database/sql"
	"io"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Grafanadbsubscriber struct {
	DbLocation string
	DbVersion  string
}

func NewGrafanaDbSubscriber(dbLocation string) (Grafanadbsubscriber, error) {
	g := Grafanadbsubscriber{DbLocation: dbLocation}

	d, err := sql.Open("sqlite3", g.DbLocation)
	if err != nil {
		return g, err
	}

	defer d.Close()

	row := d.QueryRow("SELECT sqlite_version();")
	// Retrieve the version from the query result
	err = row.Scan(&g.DbVersion)
	if err != nil {
		return g, err
	}

	return g, nil
}

func Main(args []string, printDest io.Writer) int {
	return 0
}
