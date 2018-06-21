package dbx

import (
	"database/sql"
	stdErrors "errors"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

// ErrNoChanges means that the migration process didn't change execute any file
var ErrNoChanges = stdErrors.New("nothing to migrate.")

// Migrate the database to latest version
func (db *Database) Migrate(path string) error {
	db.logger.Info("Running migrations...")
	dir, err := os.Open(env.Path(path))
	if err != nil {
		return errors.Wrap(err, "failed to open dir '%s'", path)
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return errors.Wrap(err, "failed to read files from dir '%s'", path)
	}

	versions := make([]int, len(files))
	versionFiles := make(map[int]string, len(files))
	for i, file := range files {
		fileName := file.Name()
		parts := strings.Split(fileName, "_")
		versions[i], err = strconv.Atoi(parts[0])
		versionFiles[versions[i]] = fileName
		if err != nil {
			return errors.Wrap(err, "failed to convert '%s' to number", parts[0])
		}
	}
	sort.Ints(versions)

	db.logger.Infof("Found total of @{Total} migration files.", log.Props{
		"Total": len(versions),
	})

	lastVersion, err := db.getLastMigration()
	if err != nil {
		return errors.Wrap(err, "failed to get last migration record")
	}

	db.logger.Infof("Current version is @{Version}", log.Props{
		"Version": lastVersion,
	})

	// Apply all migrations
	for _, version := range versions {
		if version > lastVersion {
			fileName := versionFiles[version]
			db.logger.Infof("Running Version: @{Version} (@{FileName})", log.Props{
				"Version":  version,
				"FileName": fileName,
			})
			err := db.runMigration(version, path, fileName)
			if err != nil {
				return errors.Wrap(err, "failed to run migration '%s'", fileName)
			}
		}
	}

	db.logger.Info("Migrations finished with success.")
	return nil
}

func (db Database) runMigration(version int, path, fileName string) error {
	filePath := env.Path(path + "/" + fileName)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "failed to read file '%s'", filePath)
	}

	trx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = trx.tx.Exec(string(content))
	if err != nil {
		return err
	}

	_, err = trx.Execute("INSERT INTO migrations_history (version, filename) VALUES ($1, $2)", version, fileName)
	if err != nil {
		return err
	}

	return trx.Commit()
}

func (db Database) getLastMigration() (int, error) {
	_, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS migrations_history (
		version     BIGINT PRIMARY KEY,
		filename    VARCHAR(100) null,
		date	 			TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return 0, err
	}

	var lastVersion sql.NullInt64
	row := db.conn.QueryRow("SELECT MAX(version) FROM migrations_history LIMIT 1")
	err = row.Scan(&lastVersion)
	if err != nil {
		return 0, err
	}

	if !lastVersion.Valid {
		// If it's the first run, maybe we have records on old migrations table, so try to get from it.
		// This SHOULD be removed in the far future.
		row := db.conn.QueryRow("SELECT version FROM schema_migrations LIMIT 1")
		row.Scan(&lastVersion)
	}

	return int(lastVersion.Int64), nil
}
