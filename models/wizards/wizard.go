package wizards

import (
	"github.com/crob1140/codewiz/datastore"
)

const (
	male = "M"
	female = "F"
)

type Wizard struct {
	datastore.BaseRecord
	OwnerID uint64 `db:"OwnerID"`
	Sex string	`db:"Sex"`
	Name string 	`db:"Name"`
}

func NewWizard(name string, sex string, ownerID uint64) *Wizard {
	return &Wizard{Name : name, Sex : sex, OwnerID : ownerID}
}