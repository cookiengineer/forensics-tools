package api

type User struct {
	Identifier   uint64   `json:"id"`
	Name         string   `json:"login"`
	Aliases      []string `json:"aliases"`
	Emails       []string `json:"emails"`
	Repositories []string `json:"repositories"`
	Type         string   `json:"type"`
}

func NewUser(name string) User {

	var user User

	user.Name = name
	user.Aliases = make([]string, 0)
	user.Emails = make([]string, 0)
	user.Repositories = make([]string, 0)

	return user

}

func (user *User) AddAlias(value string) {

	var found bool = false

	for a := 0; a < len(user.Aliases); a++ {

		if user.Aliases[a] == value {
			found = true
			break
		}

	}

	if found == false {
		user.Aliases = append(user.Aliases, value)
	}

}

func (user *User) AddEmail(value string) {

	var found bool = false

	for e := 0; e < len(user.Emails); e++ {

		if user.Emails[e] == value {
			found = true
			break
		}

	}

	if found == false {
		user.Emails = append(user.Emails, value)
	}

}

func (user *User) AddRepository(value string) {

	var found bool = false

	for r := 0; r < len(user.Repositories); r++ {

		if user.Repositories[r] == value {
			found = true
			break
		}

	}

	if found == false {
		user.Repositories = append(user.Repositories, value)
	}

}

func (user *User) HasRepository(value string) bool {

	var found bool = false

	for r := 0; r < len(user.Repositories); r++ {

		if user.Repositories[r] == value {
			found = true
			break
		}

	}

	return found

}

