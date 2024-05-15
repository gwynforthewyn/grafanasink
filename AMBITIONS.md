1. verify go-sqlite lets you register c functions. 
1. Integration tests as a sidecar mounted inside the grafana container, unit tests using go-sqlite3
1. In your container, you'll need to install sqlite-devel or whatever the correct package is.


THere's a built in way of subscribing an update hook:
https://pkg.go.dev/github.com/mattn/go-sqlite3#SQLiteConn.RegisterUpdateHook

How about creating a hook that prints the row number of the row that's edited?

```go
package main

/*
#cgo darwin CFLAGS: -I/opt/homebrew/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -lsqlite3
#cgo linux CFLAGS: -I/usr/include
#cgo linux LDFLAGS: -L/usr/lib -lsqlite3

// Include the necessary SQLite headers
#include <sqlite3.h>
#include <stdlib.h>

// Function prototypes in C for Go to call
extern void goUpdateHook(int operation, const char* dbName, const char* tableName, long long rowid);

// Wrapper to register the update hook
static void update_hook_wrapper(void *arg, int operation, const char *database, const char *table, sqlite3_int64 rowid) {
    goUpdateHook(operation, database, table, rowid);
}

// Set the update hook
void set_update_hook(sqlite3 *db) {
    sqlite3_update_hook(db, update_hook_wrapper, NULL);
}
*/
import "C"
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SQLite operations
const (
	SQLITE_INSERT = C.SQLITE_INSERT
	SQLITE_UPDATE = C.SQLITE_UPDATE
	SQLITE_DELETE = C.SQLITE_DELETE
)

//export goUpdateHook
func goUpdateHook(operation C.int, dbName, tableName *C.char, rowid C.longlong) {
	var opType string
	switch operation {
	case SQLITE_INSERT:
		opType = "INSERT"
	case SQLITE_UPDATE:
		opType = "UPDATE"
	case SQLITE_DELETE:
		opType = "DELETE"
	default:
		opType = "UNKNOWN"
	}
	fmt.Printf("Operation: %s, Table: %s, RowID: %d\n", opType, C.GoString(tableName), rowid)
}

func main() {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set the update hook
	C.set_update_hook((*C.sqlite3)(db.Driver().(*sqlite3.SQLiteDriver).Conn(db).Conn))

	// Perform some database operations to trigger the update hook
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO test (name) VALUES ('Alice'), ('Bob');")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE test SET name = 'Charlie' WHERE id = 1;")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("DELETE FROM test WHERE id = 2;")
	if err != nil {
		log.Fatal(err)
	}

	// Wait to ensure the update hooks are called
	fmt.Println("Listening for database updates. Press Ctrl+C to exit.")
	select {}
}
```