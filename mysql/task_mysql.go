package mysql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"time"
)

type TaskService struct {
	DB *sql.DB
}

func NewTaskService(dialect, dsn string, idleConn, maxConn int) (*ListService, error) {
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

func (ts *TaskService) FindTaskById(id int) (*models.Task, error) {
	var task models.Task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ts.DB.QueryRowContext(ctx, "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE ID=?", id)
	err := row.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.ListId, &task.AsigneeId)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (ts *TaskService) AllTasksInList(list *models.List) ([]*models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := ts.DB.QueryContext(ctx, "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE LIST=?", list.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	allTasks := make([]*models.Task, 0)

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.ListId, &task.AsigneeId)
		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, &task)
	}

	return allTasks, nil
}

func (ts *TaskService) CreateTask(t *models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO list (STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE) VALUES(?, ?, ?, ?, ?)"
	stmt, err := ts.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, t.Status, t.Title, t.Description, t.ListId, t.AsigneeId)
	return err
}

func (ts *TaskService) AssignTask(t *models.Task, p *models.Person) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE list SET ASIGNEE = ? WHERE ID = ?"
	stmt, err := ts.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, p.Id, t.Id)
	return err
}

func (ts *TaskService) SetStatus(t *models.Task, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE list SET STATUS = ? WHERE ID = ?"
	stmt, err := ts.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, status, t.Id)
	return err
}

func (ts *TaskService) AddTaskToList(t *models.Task, l *models.List) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "UPDATE list SET LIST = ? WHERE ID = ?"
	stmt, err := ts.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, l.Id, t.Id)
	return err
}
