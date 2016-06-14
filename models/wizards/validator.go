package wizards

import (
	"github.com/crob1140/codewiz/models"
)

type Validator struct {
	Dao *Dao
}

func NewValidator(dao *Dao) *Validator {
	return &Validator{Dao : dao}
}

func (validator *Validator) Validate(wizard *Wizard) (models.ValidationErrors, error) {

	errs := make(models.ValidationErrors)

	if !(wizard.Sex == male || wizard.Sex == female) {
		errs.Add("Sex", "Must be either male or female.")
	}

	if wizard.Name == "" {
		errs.Add("Name", "This field cannot be empty.")
	}

	savedWizard, err := validator.Dao.GetByNameAndOwnerID(wizard.Name, wizard.OwnerID)
	if err != nil {
		return nil, err
	}

	if savedWizard != nil && savedWizard.ID != wizard.ID {
		errs.Add("Name", "A wizard with this name already exists.")
	}

	return errs, nil
}