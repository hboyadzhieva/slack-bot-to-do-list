// Package mysql provides interaction with database tables for slack-bot tasks
package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	StatusOpen       = "Open"
	StatusInProgress = "In Progress"
	StatusDone       = "Done"
)

type Task struct {
	Id          int
	Status      string
	Title       string
	Description string
	AsigneeId   string
	ChannelId   string
}

func NewTask(title string, channelId string) *Task {
	task := Task{}
	task.Status = StatusOpen
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
	GetAllInChannelWithStatus(channelId string, status ...string) ([]*Task, error)
	GetAllInChannelAssignedTo(channelId string, assigneeId string) ([]*Task, error)
	AssignTaskTo(taskId int, assigneeId string) error
	SetStatus(taskId int, status string) error
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

	_, err = stmt.Exec(t.Status, t.Title, t.Description, t.AsigneeId, t.ChannelId)
	return err
}

func (repo *TaskRepository) GetTaskById(id int) (*Task, error) {
	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = ?"
	txn, err := repo.DB.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	defer txn.Commit()
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	var task Task
	err = stmt.QueryRow(id).Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.AsigneeId, &task.ChannelId)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (repo *TaskRepository) GetAllInChannel(channelId string) ([]*Task, error) {
	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = ?"
	txn, err := repo.DB.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	defer txn.Commit()
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(channelId)
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.AsigneeId, &task.ChannelId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (repo *TaskRepository) GetAllInChannelWithStatus(channelId string, status ...string) ([]*Task, error) {
	statusCount := len(status)
	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = ?"
	if statusCount == 1 {
		query = query + " AND STATUS = ?"
	} else if statusCount > 1 {
		query = query + " AND ( "
		for i := 0; i < statusCount-1; i++ {
			query = query + "STATUS = ? OR "
		}
		query = query + "STATUS = ? )"
	}
	txn, err := repo.DB.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	defer txn.Commit()
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	args := make([]interface{}, 0)
	args = append(args, channelId)
	for _, stat := range status {
		args = append(args, stat)
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.AsigneeId, &task.ChannelId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (repo *TaskRepository) GetAllInChannelAssignedTo(channelId string, assigneeId string) ([]*Task, error) {
	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = ? AND ASIGNEE_ID = ?"
	txn, err := repo.DB.Begin()
	if err != nil {
		txn.Rollback()
		return nil, err
	}
	defer txn.Commit()
	stmt, err := repo.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(channelId, assigneeId)
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.Id, &task.Status, &task.Title, &task.Description, &task.AsigneeId, &task.ChannelId)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (repo *TaskRepository) AssignTaskTo(taskId int, assigneeId string) error {
	query := "UPDATE TASK SET ASIGNEE_ID = ? WHERE ID = ?"

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

	_, err = stmt.Exec(assigneeId, taskId)
	return err
}

func (repo *TaskRepository) SetStatus(taskId int, status string) error {
	query := "UPDATE TASK SET STATUS = ? WHERE ID = ?"

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

	_, err = stmt.Exec(status, taskId)
	return err
}
