package records

import (
	"alphabet_typer/models"
	"database/sql"
	"fmt"

	"github.com/labstack/gommon/log"
)

type Repository struct {
	dbClient *sql.DB
}

func NewRepository(dbClient *sql.DB) Repository {
	return Repository{dbClient: dbClient}
}

func (r *Repository) GetAll(limit int64, filter int64) ([]models.Record, error) {
	var notes []models.Record
	query := `SELECT * FROM records`
	filterString := ""
	if filter == 1 {
		filterString = ` WHERE createdDate >= current_date at time zone 'UTC' - interval '1 days'`
	} else if filter == 7 {
		filterString = ` WHERE createdDate >= current_date at time zone 'UTC' - interval '7 days'`
	} else if filter == 30 {
		filterString = ` WHERE createdDate >= current_date at time zone 'UTC' - interval '30 days'`
	} else if filter == 365 {
		filterString = ` WHERE createdDate >= current_date at time zone 'UTC' - interval '365 days'`
	}
	orderString := ` ORDER BY record ASC LIMIT $1`
	query = fmt.Sprintf("%s%s%s", query, filterString, orderString)
	notes, err := r.fetch(query, limit)
	return notes, err
}

func (r *Repository) Post(record *models.Record) (*models.Record, error) {
	statement := `INSERT INTO records (uuid, record, username, createdDate) VALUES ($1, $2, $3, CURRENT_DATE)`
	_, err := r.dbClient.Exec(statement, record.UUID, record.Record, record.Username)
	return record, err
}

func (r *Repository) fetch(query string, limit int64) ([]models.Record, error) {
	var rows *sql.Rows
	var err error
	rows, err = r.dbClient.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Errorf("Datenbankverbindung konnte nicht korrekt geschlossen werden: %v", err)
		}
	}()
	record := make([]models.Record, 0)
	for rows.Next() {
		RecordDB := models.RecordDB{}
		err := rows.Scan(&RecordDB.UUID, &RecordDB.Record, &RecordDB.Username, &RecordDB.CreatedDate)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			log.Infof("Fehler beim Lesen der Daten: %v", err)
			return record, err
		}
		record = append(record, RecordDB.GetRecord())
	}
	return record, nil
}
