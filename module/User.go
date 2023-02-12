package module

type User struct {
	Id      string
	Name    string
	Email   string
	Balance float64
}

func (user *User) setName(name string) {
	user.Name = name
}
