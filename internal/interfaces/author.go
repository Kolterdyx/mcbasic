package interfaces

type Author struct {
	Name  string
	Email string
}

func (a Author) String() string {
	return a.Name + " <" + a.Email + ">"
}
