package mysql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"time"
)

//Implement ListService interface
// https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e
type PersonService struct {
	DB *sql.DB
}

func NewPersonService(dialect, dsn string, idleConn, maxConn int) (*PersonService, error) {
	// dialect - "mysql", "dsn" "user:password@/dbname", idle max 10
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return &PersonService{db}, nil
}

func (ps *PersonService) FindPersonById(id int) (*models.Person, error) {
	var p models.Person

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ps.DB.QueryRowContext(ctx, "SELECT ID, SLACK_ID FROM person WHERE ID=?", id)
	err := row.Scan(&p.Id, &p.SlackId)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (ps *PersonService) CreatePerson(p *models.Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO person (SLACK_ID) VALUES(?)"
	stmt, err := ps.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, p.SlackId)
	return err
}
