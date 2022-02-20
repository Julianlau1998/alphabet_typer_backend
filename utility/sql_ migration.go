package utility

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

func Migrate(conn *sql.DB) error {
	if !hasMigrationTable(conn) {
		if err := createMigrationTable(conn); err != nil {
			return err
		}
	}
	var migrated []string
	rows, err := conn.Query("SELECT filename FROM migration")
	if err != nil {
		log.Errorf("fehler beim Ausführen der migration: %v \n", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var migrateString string
		err := rows.Scan(&migrateString)
		if err != nil {
			log.Errorf("fehler beim Ausführen der migration: %v \n", err)
			return err
		}
		migrated = append(migrated, migrateString)
	}

	var files []os.FileInfo
	var migratePath = "migration/"
	err = filepath.Walk(migratePath, func(path string, info os.FileInfo, err error) error {
		files = append(files, info)
		return nil
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			if !isMigrated(migrated, file.Name()) {
				err = migrateFile(migratePath+file.Name(), file, conn)
				if err != nil {
					log.Errorf("fehler beim Ausführen der migration (%s): %v \n", file.Name(), err)
					return err
				}
				log.Debugf("%s erfolgreich migriert", file.Name())
			}
		}
	}
	return nil
}

func hasMigrationTable(conn *sql.DB) bool {
	var dummy int64
	err := conn.QueryRow("SELECT 1 FROM migration LIMIT 1").Scan(&dummy)
	if err != nil && err != sql.ErrNoRows {
		log.Debugf("die Tabelle 'migration' wurde nicht gefunden")
		return false
	}
	fmt.Println("Die Tabelle wurde gefunden")
	return true
}

func createMigrationTable(conn *sql.DB) error {
	if err := os.Setenv("IMPORT", "TRUE"); err != nil {
		log.Infof("IMPORT ENV variable konnte nicht gesetzt werden, testdaten werden nicht importiert!")
	}

	_, err := conn.Exec("CREATE TABLE migration (" +
		"id uuid NOT NULL PRIMARY KEY, " +
		"filename VARCHAR(255) NOT NULL" +
		");")
	if err != nil {

	}
	return err
}

func isMigrated(migrated []string, filename string) bool {
	for _, migrate := range migrated {
		if migrate == filename {
			return true
		}
	}
	return false
}

func migrateFile(filepath string, info os.FileInfo, conn *sql.DB) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	_, err = conn.Exec(string(data))
	if err != nil {
		return err
	}
	u, err := uuid.NewV4()
	_, err = conn.Exec("INSERT INTO migration(id, filename) VALUES ($1, $2);", u.String(), info.Name())
	return err
}
