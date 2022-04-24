package records

import (
	"alphabet_typer/models"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	recordRepository Repository
}

func NewService(recordRepository Repository) Service {
	return Service{recordRepository: recordRepository}
}

func (s *Service) GetAll(limit int64, filter int64, offset int64) ([]models.Record, error) {
	notes, err := s.recordRepository.GetAll(limit, filter, offset)
	if err != nil {
		log.Warnf("RecordsService.GetAll(): Could not get Records: %s", err)
		return notes, err
	}
	return notes, err
}

func (s *Service) Post(record *models.Record) (*models.Record, error) {
	id, err := uuid.NewV4()
	record.UUID = id.String()
	record, err = s.recordRepository.Post(record)
	if err != nil {
		log.Warnf("RecordService.Post(): Could not post record: %s", err)
		return record, err
	}
	return record, err
}
