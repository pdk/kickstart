package migrate

// stupid/simple database migration

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// constants for managing migrations
const (
	Success               = "SUCCESS"
	Error                 = "ERROR"
	CreateTableMigrations = `
		create table if not exists schema_migrations (
			script_name varchar not null,
			executed_at timestamp with time zone not null,
			status varchar not null,
			error text
		);`
)

// Database executes scripts from fileList that have not successfully run, yet.
func Database(db *sql.DB, fileList []string) {

	_, err := db.Exec(CreateTableMigrations)
	if err != nil {
		log.Fatalf("failed to create table schema_migrations: %v", err)
	}

	skipCount, execCount := 0, 0

	for _, scriptPath := range fileList {

		scriptContent, err := ioutil.ReadFile(scriptPath)
		if err != nil {
			log.Fatalf("cannot read script %s: %v", scriptPath, err)
		}

		scriptName := strings.TrimSuffix(filepath.Base(scriptPath), ".sql")
		successRowCount := 0
		err = db.QueryRow(
			`select count(*) as success_count from schema_migrations where script_name = ? and status = ?`,
			scriptName, Success).Scan(&successRowCount)
		if err != nil {
			log.Fatalf("failed to query schema_migrations table: %v", err)
		}
		if successRowCount > 0 {
			skipCount++
			continue
		}

		_, err = db.Exec(string(scriptContent))
		if err != nil {
			log.Printf("failed to execute script %s: %v", scriptPath, err)

			_, err = db.Exec(
				`insert into schema_migrations (script_name, executed_at, status, error) values (?, ?, ?, ?)`,
				scriptName, time.Now(), Error, err.Error())

			if err != nil {
				log.Printf("failed to update schema_migrations with error: %v", err)
			}

			log.Fatalf("aborting migration")
		}

		_, err = db.Exec(`insert into schema_migrations (script_name, executed_at, status) values (?, ?, ?)`,
			scriptName, time.Now(), Success)
		if err != nil {
			log.Fatalf("script executed successfully, but failed to update schema_migrations for script %s: %v", scriptName, err)
		}

		log.Printf("executed %s", scriptPath)
		execCount++
	}

	log.Printf("skipped %d scripts", skipCount)
	log.Printf("executed %d scripts", execCount)
}
