package db

import (
	"database/sql"
	"fmt"
	"github.com/heroku/changemomentum/logger"
	"github.com/heroku/changemomentum/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx"
)

type Postgresrepo struct {
	Db *sql.DB
}

func NewPostgresrepo(dsn *string) (*Postgresrepo, error) {
	db, err := sqlx.Open("postgres", *dsn)
	if err != nil {
		return nil, err
	}
	return &Postgresrepo{
		db,
	}, nil
}

type dbError struct {
	method string
	Err    error
}

func (re *dbError) Error() string {
	return fmt.Sprintf(
		"DB error:  %s, err: %v",
		re.method,
		re.Err,
	)
}

func (db Postgresrepo) Close() {
	db.Db.Close()
}

func (db Postgresrepo) AddContact(firstname string, lastname string) error {
	result, err := db.Db.Exec(
		"INSERT INTO Contacts (firstname, lastname) VALUES ($1, $2)",
		firstname,
		lastname,
	)
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}

	lastID, err := result.LastInsertId()

	logger.Info("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return nil
}

func (db Postgresrepo) AddPhone(idContact int, phone string) error {
	result, err := db.Db.Exec(
		"INSERT INTO Phonenumber (contact_id,phonenumber) VALUES ($1,$2)",
		idContact,
		phone,
	)
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}

	lastID, err := result.LastInsertId()

	logger.Info("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return nil
}

func (db Postgresrepo) List() ([]schema.Participant, error) {
	//participants := make(map[int]schema.Participant)
	var participants []schema.Participant
	sqlStr := "select id, firstname, lastname, command, data_registration, usertokenid from participants"
	rows, err := db.Db.Query(sqlStr)
	for rows.Next() {
		participant := &schema.Participant{}
		err = rows.Scan(&participant.Id, &participant.FirstName, &participant.LastName, &participant.Command, &participant.Date, &participant.UsertokenId)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil,  err
		}
		//db.selectItemUsertoken(participant)
		participants[participant.Id] = *participant
	}

	rows.Close()
	return participants, nil
}

func (db Postgresrepo) Delete(id int) error {
	result, err := db.Db.Exec(
		"DELETE FROM Contacts WHERE id = $1",
		id,
	)
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Can't add rows", err)
		return err
	}

	logger.Info("Delete - RowsAffected", affected)

	return nil

}



func (db Postgresrepo) Update(contact schema.Contact, phones []string) error {

	tx, err := db.Db.Begin()

	sqlstr := "UPDATE Contacts SET firstname =$1 , lastname=$2 where id=$3"

	if _, err = tx.Exec(sqlstr, contact.FirstName, contact.LastName, contact.Id); err != nil {
		tx.Rollback()
		logger.Error("Can't insert rows", err)
		return err
	}

	sqlstr1 := "delete from Phonenumber where contact_id = $1 "

	if _, err = tx.Exec(sqlstr1, contact.Id); err != nil {
		tx.Rollback()
		logger.Error("Can't insert rows", err)
		return err
	}

	for _, number := range phones {
		if number != "" {
			sqlstr2 := "insert into Phonenumber (contact_id, Phonenumber) VALUES ($1,$2)"
			if _, err = tx.Exec(sqlstr2, contact.Id, number); err != nil {
				tx.Rollback()
				logger.Error("Can't insert rows", err)
				return err
			}
		}
	}
	err = tx.Commit()

	return err
}

func (db Postgresrepo) selectItemUsertoken(participant *schema.Participant) error {
	rows, err := db.Db.Query("select login from UsersToken where id = $1", participant.UsertokenId)
	usertoken := &schema.UsersToken{}
	err = rows.Scan(&usertoken.Login)

		if err != nil {
			logger.Error("Can't select rows", err)
			return err
		}
	rows.Close()
	return nil
}

func (db Postgresrepo) SelectItem(id int) (schema.Contact, error) {
	sqlStr := "select id, firstname, lastname from Contacts where id=$1"
	rowscontact := db.Db.QueryRow(sqlStr, id)
	contact := &schema.Contact{}
	err := rowscontact.Scan(&contact.Id, &contact.FirstName, &contact.LastName)
	if err != nil {
		logger.Error("Can't select rows", err)
		return schema.Contact{}, err
	}
	//db.selectItemPhones(contact)

	return *contact, nil
}


func (db Postgresrepo) Search(search string) (map[int]schema.Contact, error) {
	contacts := make(map[int]schema.Contact)
	sqlStr := "select id, firstname, lastname from Contacts where upper(firstname) like upper('%'||$1||'%') or upper(lastname)  like upper('%'||$2||'%')"
	rows, err := db.Db.Query(sqlStr, search, search)
	if err != nil {
		logger.Error("Can't select rows", err)
		return nil, err
	}
	for rows.Next() {
		contact := &schema.Contact{}
		err = rows.Scan(&contact.Id, &contact.FirstName, &contact.LastName)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil, err
		}
		//db.selectItemPhones(contact)
		contacts[contact.Id] = *contact
	}
	rows.Close()

	rows, err = db.Db.Query("select contact_id from Phonenumber where upper(Phonenumber) like upper('%'||$1||'%')", search)
	for rows.Next() {
		contactId := new(int)
		err = rows.Scan(&contactId)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil, err
		}
		contact, err := db.SelectItem(*contactId)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil, err
		}

		contacts[contact.Id] = contact
	}

	rows.Close()
	return contacts, nil
}

