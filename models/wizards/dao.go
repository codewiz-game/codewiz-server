package wizards

import (
	"github.com/crob1140/codewiz-server/datastore"
)

type Dao struct {
	DB *datastore.DB
}

func NewDao(db *datastore.DB) *Dao {
	db.AddTableWithName(Wizard{}, "Wizards")
	return &Dao{DB : db}
}

func (dao *Dao) GetByID(id uint64) (*Wizard, error) {
	wizard, err := dao.DB.Get(Wizard{}, "SELECT * FROM Wizards WHERE ID = ?", id)
	if err != nil || wizard == nil {
		return nil, err
	}
	return wizard.(*Wizard), err
}

func (dao *Dao) GetByOwnerID(ownerID uint64) ([]*Wizard, error) {
	var wizards []*Wizard
	_, err := dao.DB.Select(&wizards, "SELECT * FROM Wizards WHERE OwnerID = ?", ownerID)
	return wizards, err
}

func (dao *Dao) GetByNameAndOwnerID(name string, ownerID uint64) (*Wizard, error) {
	wizard, err := dao.DB.Get(Wizard{}, "SELECT * FROM Wizards WHERE Name = ? AND OwnerID = ?", name, ownerID)
	if err != nil || wizard == nil {
		return nil, err
	}
	return wizard.(*Wizard), err
}

func (dao *Dao) Update(wizard *Wizard) error {
	return dao.DB.Update(wizard)
}

func (dao *Dao) Delete(wizard *Wizard) error {
	return dao.DB.Delete(wizard)
}

func (dao *Dao) Insert(wizard *Wizard) error {
	return dao.DB.Insert(wizard)
}