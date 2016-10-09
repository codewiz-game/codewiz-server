package users

const (
	Visitor = "visitor"
	Standard = "standard"
	Admin = "admin"
)

func (user *User) Roles() []string {
	if user == nil {
		return []string{Visitor}
	}

	return []string{Standard}
}