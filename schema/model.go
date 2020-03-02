package schema

type Contact struct {
	Id        int
	FirstName string
	LastName  string
	Phones []string
}

type Phone struct {
	ContactId          int
	PhoneNumber string
}

type Participant struct {
	Id        int
	FirstName string
	LastName  string
	Command string
	Date	string
	UsertokenId int
}

type UsersToken struct {
	Id int
	Login string
	FirstName string
	LastName  string
	email string
}

