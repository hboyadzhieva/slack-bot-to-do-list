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
type ListService struct {
	DB *sql.DB
}

func NewListService(dialect, dsn string, idleConn, maxConn int) (*ListService, error) {
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

	return &ListService{db}, nil
}

func (ls *ListService) FindListById(id int) (*models.List, error) {
	var list models.List

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ls.DB.QueryRowContext(ctx, "SELECT ID, TITLE, CHANNEL FROM list WHERE ID=?", id)
	err := row.Scan(&list.Id, &list.Title, &list.ChannelId)
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (ls *ListService) AllLists() ([]*models.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := ls.DB.QueryContext(ctx, "SELECT ID, TITLE, CHANNEL FROM list")

	if err != nil {
		return nil, err
	}

	allLists := make([]*models.List, 0)
	defer rows.Close()

	for rows.Next() {
		var list models.List
		err := rows.Scan(&list.Id, &list.Title, &list.ChannelId)
		if err != nil {
			return nil, err
		}
		allLists = append(allLists, &list)
	}

	return allLists, nil
}
func (ls *ListService) AllListsInChannel(c *models.Channel) ([]*models.List, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := ls.DB.QueryContext(ctx, "SELECT ID, TITLE, CHANNEL FROM list WHERE CHANNEL=?", c.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	allLists := make([]*models.List, 0)

	for rows.Next() {
		var list models.List
		err := rows.Scan(&list.Id, &list.Title, &list.ChannelId)
		if err != nil {
			return nil, err
		}
		allLists = append(allLists, &list)
	}

	return allLists, nil
}
func (ls *ListService) CreateList(l *models.List) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO list (TITLE, CHANNEL) VALUES(?, ?)"
	stmt, err := ls.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, l.Title, l.ChannelId)
	return err
}
func (ls *ListService) DeleteList(id int) error {
	return nil
}
func (ls *ListService) AddTask(t *models.Task) error {
	return nil
}
