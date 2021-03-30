package store

import "github.com/poncheska/iot-mousetrap/pkg/models"

type Store struct {
	Mousetrap
	Organisation
}

type Mousetrap interface {
	GetAll(OrgId int64) ([]models.Mousetrap, error)
	Create(mt models.Mousetrap) (int64, error)
	Update(mt models.Mousetrap) error
	GetByName(name, orgName string) (models.Mousetrap, error)
}

type Organisation interface {
	GetByCredentials(name, password string) (models.Organisation, error)
	Create(mt models.Organisation) (int64, error)
}
