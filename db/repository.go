package db

import (
	"github.com/heroku/changemomentum/schema"
)

type Repository interface {
	Close()
	AddParticipant(string, string, string, int) error
	List() ([]schema.Participant, error)
	AddPhone (int , string) error
	SelectItem (int) (schema.Participant, error)
	Delete (int) error
	Update (schema.Participant)  error
	Search(string) ([]schema.Participant,  error)
	Registration (int) error
	UnRegistration (int) error

}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}


func List() ([]schema.Participant, error)  {
	return impl.List()
}

func AddContact(firstname string, lastname string, command string,id int) error  {
	return impl.AddParticipant(firstname,lastname, command, id)
}

func AddPhone(id int, number string) error  {
	return impl.AddPhone(id,number)
}

func SelectItem( id int) (schema.Participant,  error)  {
	return impl.SelectItem(id)
}

func Delete(id int) error {
	return impl.Delete(id)
}

func Update(participant schema.Participant) error  {
	return impl.Update(participant)
}

func Search(search string) ([]schema.Participant, error)  {
	return impl.Search(search)
}

func Registration(id int) (error)  {
	return impl.Registration(id)
}

func UnRegistration(id int) (error)  {
	return impl.UnRegistration(id)
}
