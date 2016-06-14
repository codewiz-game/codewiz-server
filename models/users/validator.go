package users

import (
	"github.com/crob1140/codewiz/models"
	"net/mail"
)

type Validator struct {
	Dao *Dao
}

func NewValidator(dao *Dao) *Validator {
	return &Validator{Dao : dao}
}

func (validator *Validator) Validate(user *User) (models.ValidationErrors, error) {

	errs := make(models.ValidationErrors)

	if user.Username == "" {
		errs.Add("Username", "Username must be provided.")
	}

	if user.HashedPassword == "" {
		errs.Add("Password", "Password must be provided.")
	}

	if user.Email == "" {
		errs.Add("Email", "Email must be provided.")
	} else if _, err := mail.ParseAddress(user.Email); err != nil {
		errs.Add("Email", "Invalid email address.")
	}

	savedUser, err := validator.Dao.GetByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	if savedUser != nil && savedUser.ID != user.ID {
		errs.Add("Username", "A user with this username already exists.")
	}

	return errs, nil

}