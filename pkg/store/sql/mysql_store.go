package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
)

type MousetrapStore struct {
	db *sqlx.DB
}

type OrganisationStore struct {
	db *sqlx.DB
}

func NewFakeStore(db *sqlx.DB) store.Store {
	return store.Store{
		Mousetrap:    &MousetrapStore{db: db},
		Organisation: &OrganisationStore{db: db},
	}
}

func (ms *MousetrapStore) GetAll(OrgId int64) ([]models.Mousetrap, error) {
	return []models.Mousetrap{}, nil
}

func (ms *MousetrapStore) Create(mt models.Mousetrap) (int64, error) {
	return 0, nil
}

func (ms *MousetrapStore) Update(mt models.Mousetrap) error {
	return nil
}

func (ms *MousetrapStore) GetByName(name, orgName string) (models.Mousetrap, error) {
	return models.Mousetrap{}, nil
}

func (os *OrganisationStore) GetByCredentials(name, password string) (models.Organisation, error) {
	return models.Organisation{}, nil
}

func (os *OrganisationStore) Create(org models.Organisation) (int64, error) {
	return 0, nil
}
