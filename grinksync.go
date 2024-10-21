package grinksync

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

type Grafanadbsubscriber struct {
	DbLocation string
	DbVersion  string
}

func NewGrafanaDbSubscriber(dbLocation string, driverName string, printDest io.Writer) (Grafanadbsubscriber, error) {

	location, err := filepath.Abs(dbLocation)
	if err != nil {
		return Grafanadbsubscriber{}, err
	}

	g := Grafanadbsubscriber{DbLocation: location}
	sqlite3conn := []*sqlite3.SQLiteConn{}

	sql.Register(driverName,
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				sqlite3conn = append(sqlite3conn, conn)
				conn.RegisterUpdateHook(func(op int, db string, table string, rowid int64) {
					// switch op {
					// case sqlite3.SQLITE_CREATE_TABLE:
					fmt.Fprint(printDest, "Notified of operation", op, " on db", db, "table", table, "rowid", rowid)
					log.Println(printDest, "Notified of operation", op, " on db", db, "table", table, "rowid", rowid)
					// }
				})
				return nil
			},
		})

	err = os.Remove(g.DbLocation)
	if err != nil {
		return g, err
	}

	d, err := sql.Open(driverName, g.DbLocation)
	if err != nil {
		return g, err
	}

	defer d.Close()

	_, err = d.Exec("CREATE TABLE foo (id integer not null primary key, content text);")
	if err != nil {
		return g, err
	}

	_, err = d.Exec("INSERT INTO foo (content) VALUES ('spamalot');")
	if err != nil {
		return g, err
	}

	row := d.QueryRow("SELECT sqlite_version();")
	// Retrieve the version from the query result
	err = row.Scan(&g.DbVersion)
	if err != nil {
		return g, err
	}

	return g, nil
}

func Main(args []string, printDest io.Writer) int {
	dbLocation := "./sqlite.db"
	driverName := "grinksync"

	_, err := NewGrafanaDbSubscriber(dbLocation, driverName, os.Stdout)

	if err != nil {
		panic(err)
	}

	// hang out forever
	// cf. https://groups.google.com/g/golang-nuts/c/H6DwLS4mUeM/m/kRbWGj0jd8wJ
	<-make(chan int)

	return 0
}
