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
			res = append(res, v)
		}
	}
	return res, nil
}

func (ms *MousetrapStore) Create(mt models.Mousetrap) (int64, error) {
	org := true
	for _, v := range ms.repo.Organisations {
		if v.Name == mt.OrgName {
			org = false
			break
		}
	}
	if org {
		return 0, fmt.Errorf("invalid org_name")
	}
	for _, v := range ms.repo.Mousetraps {
		if mt.OrgName == v.OrgName && mt.Name == v.Name {
			return 0, fmt.Errorf("mousetrap with org_name = %v and name = %v already exist",
				mt.OrgName, mt.Name)
		}
	}
	mt.Id = ms.lastId
	ms.lastId++
	ms.repo.Mousetraps = append(ms.repo.Mousetraps, mt)
	return mt.Id, nil
}

func (ms *MousetrapStore) Update(mt models.Mousetrap) error {
	for i, v := range ms.repo.Mousetraps {
		if v.Id == mt.Id {
			if v.Name != mt.Name || v.OrgName != mt.OrgName {
				return fmt.Errorf("invalid mousetrap %v/%v", v.OrgName, v.Name)
			}
			ms.repo.Mousetraps[i] = mt
			return nil
		}
	}
	return nil
}

func (ms *MousetrapStore) GetByName(name, orgName string) (models.Mousetrap, error) {
	for _, v := range ms.repo.Mousetraps {
		if v.Name == name && v.OrgName == orgName {
			return v, nil
		}
	}
	return models.Mousetrap{}, fmt.Errorf("no mousetrap with name = %v and org_name = %v", name, orgName)
}

func (os *OrganisationStore) GetByCredentials(name, password string) (models.Organisation, error) {
	for _, v := range os.repo.Organisations {
		if v.Name == name {
			if v.Password == password {
				return v, nil
			} else {
				return models.Organisation{}, fmt.Errorf("invalid password")
			}
		}
	}
	return models.Organisation{}, fmt.Errorf("no such organisation")
}

func (os *OrganisationStore) Create(org models.Organisation) (int64, error) {
	for _, v := range os.repo.Organisations {
		if org.Name == v.Name {
			return 0, fmt.Errorf("organisation with name = %v already exist", org.Name)
		}
	}
	org.Id = os.lastId
	os.lastId++
	os.repo.Organisations = append(os.repo.Organisations, org)
	return org.Id, nil
}
