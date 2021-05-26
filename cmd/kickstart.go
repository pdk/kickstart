package main

import (
	"database/sql"
	"flag"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/pdk/kickstart/migrate"
	"github.com/pdk/kickstart/server"
	"github.com/pdk/kickstart/watch"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	migrateDB := flag.Bool("migrate", false, "execute the given migration scripts if not executed")
	watchFiles := flag.String("watch", "", "shutdown when files match given patterns (wrap with script to auto-restart)")
	databaseName := flag.String("db", "data.db", "name of sqlite3 database file")

	flag.Parse()

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *migrateDB {
		scripts := os.Args[2:] // grab all the args AFTER the -migrate
		migrate.Database(db, scripts)
		return
	}

	if err := run(db, *watchFiles, *databaseName); err != nil {
		log.Fatalf("%v", err)
	}
}

func run(db *sql.DB, watchFiles, databaseName string) error {

	srv := server.New(db)

	shutdown := srv.ListenAndServe()

	watchFilesOrWaitForever(strings.Fields(watchFiles))

	shutdown()

	return nil
}

func watchFilesOrWaitForever(watchFiles []string) {

	if len(watchFiles) > 0 {
		watch.Files(".", watchFiles)
		return
	}

	// otherwise just hang forever
	time.Sleep(time.Duration(math.MaxInt64))
}
