package fake

import (
	"fmt"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
)

type Repository struct {
	Mousetraps    []models.Mousetrap
	Organisations []models.Organisation
}

type MousetrapStore struct {
	repo   *Repository
	lastId int64
}

type OrganisationStore struct {
	repo   *Repository
	lastId int64
}

func NewFakeStore() store.Store {
	repo := &Repository{}
	return store.Store{
		Mousetrap:    &MousetrapStore{repo: repo},
		Organisation: &OrganisationStore{repo: repo},
	}
}

func (ms *MousetrapStore) GetAll(OrgId int64) ([]models.Mousetrap, error) {
	name := ""
	for _, v := range ms.repo.Organisations {
		if v.Id == OrgId {
			name = v.Name
			break
		}
	}
	if name == "" {
		return nil, fmt.Errorf("no organisation with id = %v", OrgId)
	}
	var res []models.Mousetrap
	for _, v := range ms.repo.Mousetraps {
		if v.OrgName == name {
			name = v.Name
			res = append(res, v)
		}
	}
	return res, nil
}

func (ms *MousetrapStore) Create(mt models.Mousetrap) error {
	for _, v := range ms.repo.Mousetraps {
		if mt.Id == v.Id {
			return fmt.Errorf("mousetrap with id = %v already exist", mt.Id)
		}
		if mt.OrgName == v.OrgName && mt.Name == v.Name {
			return fmt.Errorf("mousetrap with org_name = %v and name = %v already exist",
				mt.OrgName, mt.Name)
		}
	}
	mt.Id = ms.lastId
	ms.lastId ++
	ms.repo.Mousetraps = append(ms.repo.Mousetraps, mt)
	return nil
}

func (os *OrganisationStore) GetByCredentials(name, password string) (models.Organisation, error) {
	return models.Organisation{}, nil
}

func (os *OrganisationStore) Create(org models.Organisation) error {
	for _, v := range os.repo.Organisations {
		if org.Id == v.Id {
			return fmt.Errorf("organisation with id = %v already exist", org.Id)
		}
		if org.Name == v.Name {
			return fmt.Errorf("organisation with name = %v already exist", org.Name)
		}
	}
	org.Id = os.lastId
	os.lastId ++
	os.repo.Organisations = append(os.repo.Organisations, org)
	return nil
}
