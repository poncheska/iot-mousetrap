package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/poncheska/iot-mousetrap/pkg/models"
	"github.com/poncheska/iot-mousetrap/pkg/store"
	"log"
	"time"
)

type MousetrapStore struct {
	db *sqlx.DB
}

type OrganisationStore struct {
	db *sqlx.DB
}

type DBMousetrap struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	OrgId       int64  `db:"org_id"`
	Status      bool   `db:"status"`
	LastTrigger int64  `db:"last_trig"`
}

type MousetrapResp struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	OrgId       int64  `db:"org_id"`
	Status      bool   `db:"status"`
	LastTrigger string  `db:"last_trig"`
}

func (dmt DBMousetrap) Parse() models.Mousetrap {
	return models.Mousetrap{
		Id:          dmt.Id,
		Name:        dmt.Name,
		OrgId:       dmt.OrgId,
		Status:      dmt.Status,
		LastTrigger: time.Unix(0, dmt.LastTrigger),
	}
}

func (dmt DBMousetrap) ParseResp() MousetrapResp {
	return MousetrapResp{
		Id:          dmt.Id,
		Name:        dmt.Name,
		OrgId:       dmt.OrgId,
		Status:      dmt.Status,
		LastTrigger: time.Unix(0, dmt.LastTrigger).Format(time.UnixDate),
	}
}

func NewMySQLStore(db *sqlx.DB) store.Store {
	return store.Store{
		Mousetrap:    &MousetrapStore{db: db},
		Organisation: &OrganisationStore{db: db},
	}
}

func (ms *MousetrapStore) GetAll(OrgId int64) ([]models.Mousetrap, error) {
	res := []models.Mousetrap{}
	dmt := DBMousetrap{}
	rows, err := ms.db.Queryx("SELECT * FROM Mousetrap WHERE org_id = ?", OrgId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(&dmt)
		if err != nil {
			log.Fatalln(err)
		}
		res = append(res, dmt.ParseResp())
	}
	return res, nil
}

func (ms *MousetrapStore) Create(mt models.Mousetrap) (int64, error) {
	res, err := ms.db.Exec("INSERT INTO Mousetrap (name, org_id, status, last_trig) VALUES (?,?,?,?)",
		mt.Name, mt.OrgId, mt.Status, mt.LastTrigger.UTC().UnixNano())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (ms *MousetrapStore) Update(mt models.Mousetrap) error {
	_, err := ms.db.Exec("UPDATE Mousetrap SET status = ?, last_trig = ? WHERE id = ?",
		mt.Status, mt.LastTrigger.UTC().UnixNano(), mt.Id)
	return err
}

func (ms *MousetrapStore) GetByName(name string, orgId int64) (models.Mousetrap, error) {
	dmt := DBMousetrap{}
	err := ms.db.Get(&dmt, "SELECT * FROM Mousetrap WHERE name = ? AND org_id = ?", name, orgId)
	if err != nil {
		return models.Mousetrap{}, err
	}
	return dmt.Parse(), err
}

func (os *OrganisationStore) GetByCredentials(name, password string) (models.Organisation, error) {
	org := models.Organisation{}
	err := os.db.Get(&org, "SELECT * FROM Organisation WHERE name = ? AND password = ?", name, password)
	if err != nil {
		return models.Organisation{}, err
	}
	return org, err
}

func (os *OrganisationStore) Create(org models.Organisation) (int64, error) {
	res, err := os.db.Exec("INSERT INTO Organisation (name, password) VALUES (?, ?)", org.Name, org.Password)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
