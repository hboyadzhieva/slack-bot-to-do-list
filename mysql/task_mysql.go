// Package mysql provides interaction with database tables for slack-bot tasks
package mysql

import (
	//"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	//"time"
)

const (
	statusOpen       = "Open"
	statusInProgress = "In Progress"
	statusDone       = "Done"
)

type Task struct {
	Id          int    `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AsigneeId   string `json:"asignee"`
	ChannelId   string `json:"channel"`
}

func NewTask(title string, channelId string) *Task {
	task := Task{}
	task.Status = statusOpen
	task.Title = title
	task.Description = ""
	task.AsigneeId = ""
	task.ChannelId = channelId
	return &task
}

type TaskRepositoryInterface interface {
	PersistTask(t *Task) error
	GetTaskById(id int) (*Task, error)
	GetAllInChannel(channelId string) ([]*Task, error)
	GetAllInChannelWithStatus(channelId string, status ...string)
	GetAllInChannelAssignedTo(channelId string, assigneeId string)
	AssignTaskTo(taskId int, assigneeId string)
	SetStatus(taskId int, status string)
}

type TaskRepository struct {
	DB *sql.DB
}

func (repo *TaskRepository) PersistTask(t *Task) error {
	query := "INSERT INTO TASK (STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID) VALUES (?,?,?,?,?)"

	txn, err := repo.DB.Begin()
	if err != nil {
		txn.Rollback()
		return err
	}
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	defer txn.Commit()

	_, err = txn.Exec(query, t.Status, t.Title, t.Description, t.AsigneeId, t.ChannelId)
	return err
}

/*func (ts *TaskRepository) GetTaskById(id int) (*Task, error) {
	var task Task

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := ts.DB.QueryRowContext(ctx, "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE ID=?", id)
	err := row.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.ListId, &task.AsigneeId)
	if err != nil {
		return nil, err
	}
	return &task, nil
}*/
/*
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
*/
