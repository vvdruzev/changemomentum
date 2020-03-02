package db

import (
	"database/sql"
	"fmt"
	"github.com/heroku/changemomentum/logger"
	"github.com/heroku/changemomentum/schema"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

type Postgresrepo struct {
	Db *sql.DB
}

func NewPostgresrepo(dsn *string) (*Postgresrepo, error) {
	db, err := sql.Open("postgres", *dsn)
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

func (db Postgresrepo) AddParticipant(firstname string, lastname string, command string, tokenId int) error {
	result, err := db.Db.Exec(
		"INSERT INTO participants (firstname, lastname, command, usertokenid) VALUES ($1, $2, $3, $4)",
		firstname,
		lastname,
		command,
		tokenId,
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
	var comm sql.NullString
	var date sql.NullTime
	for rows.Next() {
		participant := &schema.Participant{}
		err = rows.Scan(&participant.Id, &participant.FirstName, &participant.LastName, &comm, &date, &participant.UsertokenId)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil,  err
		}
		participant.Command = comm.String
		participant.Date = date.Time.Format("2006-01-02 15:04:05")
		//db.selectItemUsertoken(participant)
		//participants[participant.Id] = *participant
		participants = append(participants, *participant)
	}

	rows.Close()
	return participants, nil
}

func (db Postgresrepo) Delete(id int) error {
	result, err := db.Db.Exec(
		"DELETE FROM participants WHERE id = $1",
		id,
	)
	if err != nil {
		logger.Error("Can't delete rows", err)
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logger.Error("Can't delete rows", err)
		return err
	}

	logger.Info("Delete - RowsAffected", affected)

	return nil

}

func (db Postgresrepo) Registration(id int) error {
	result, err := db.Db.Exec(
		"UPDATE participants SET data_registration =$1 WHERE id = $2",
		time.Now(),
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

	lastID, err := result.LastInsertId()

	logger.Info("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return nil

}

func (db Postgresrepo) UnRegistration(id int) error {
	result, err := db.Db.Exec(
		"UPDATE participants SET data_registration = null WHERE id = $1",
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

	lastID, err := result.LastInsertId()

	logger.Info("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return nil

}






func (db Postgresrepo) Update(participant schema.Participant) error {

	tx, err := db.Db.Begin()

	sqlstr := "UPDATE participants SET firstname =$1 , lastname=$2, command=$3 where id=$4"

	if _, err = tx.Exec(sqlstr, participant.FirstName, participant.LastName,participant.Command, participant.Id); err != nil {
		tx.Rollback()
		logger.Error("Can't insert rows", err)
		return err
	}
tx.Commit()
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

func (db Postgresrepo) SelectItem(id int) (schema.Participant, error) {
	sqlStr := "select id, firstname, lastname, command from participants where id=$1"
	rowscontact := db.Db.QueryRow(sqlStr, id)
	participant := &schema.Participant{}
	var command sql.NullString
	err := rowscontact.Scan(&participant.Id, &participant.FirstName, &participant.LastName, &command)
	if err != nil {
		logger.Error("Can't select rows", err)
		return schema.Participant{}, err
	}
	//db.selectItemPhones(contact)
	participant.Command = command.String
	return *participant, nil
}


func (db Postgresrepo) Search(search string) ([]schema.Participant, error) {
	//contacts := make(map[int]schema.Contact)
	var participants []schema.Participant
	sqlStr := "select id, firstname, lastname, command, data_registration, usertokenid from participants where upper(firstname) like upper('%'||$1||'%') or upper(lastname)  like upper('%'||$2||'%')"
	rows, err := db.Db.Query(sqlStr, search,search)
	var comm sql.NullString
	var date sql.NullTime
	for rows.Next() {
		participant := &schema.Participant{}
		err = rows.Scan(&participant.Id, &participant.FirstName, &participant.LastName, &comm, &date, &participant.UsertokenId)
		if err != nil {
			logger.Error("Can't select rows", err)
			return nil,  err
		}
		participant.Command = comm.String
		participant.Date = date.Time.Format("2006-01-02 15:04:05")
		//db.selectItemUsertoken(participant)
		//participants[participant.Id] = *participant
		participants = append(participants, *participant)
	}

	rows.Close()
	return participants, nil

}

