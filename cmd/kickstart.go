package main

import (
	"database/sql"
	"flag"
	"log"
	"math"
	"strings"
	"time"

	"github.com/pdk/kickstart/server"
	"github.com/pdk/kickstart/watch"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	watchFiles := flag.String("watch", "", "shutdown when files match given patterns (wrap with script to auto-restart)")
	databaseName := flag.String("db", "data.db", "name of sqlite3 database file")

	flag.Parse()

	if err := run(*watchFiles, *databaseName); err != nil {
		log.Fatalf("%v", err)
	}
}

func run(watchFiles, databaseName string) error {

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("%v", err)
	}

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
