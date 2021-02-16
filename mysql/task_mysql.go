// Package mysql provides interaction with database tables for slack-bot tasks
package mysql

import (
	"database/sql"
	"errors"
	//SQL Driver
	_ "github.com/go-sql-driver/mysql"
)

// String constants to set to Status in task database table.
const (
	StatusOpen       = "Open"
	StatusInProgress = "In Progress"
	StatusDone       = "Done"
)

// Task entity to represent database records
type Task struct {
	ID        int
	Status    string
	Title     string
	AsigneeID string
	ChannelID string
}

// ErrNoRowOrMoreThanOne database error when exactly 1 result is expected.
var ErrNoRowOrMoreThanOne = errors.New("sql: Expected exactly one row to be affected")

// NewTask constructs a task object. Pass title and channel id.
// Default status: "Open", default asignee value: "Not assigned".
func NewTask(title string, channelID string) *Task {
	task := Task{}
	task.Status = StatusOpen
	task.Title = title
	task.AsigneeID = "Not assigned"
	task.ChannelID = channelID
	return &task
}

// TaskRepositoryInterface provides functions for database operation execution on table TASK
type TaskRepositoryInterface interface {
	PersistTask(t *Task) error
	GetTaskByID(ID int) (*Task, error)
	GetAllInChannel(channelID string) ([]*Task, error)
	AssignTaskTo(taskID int, assigneeID string) error
	SetStatus(taskID int, status string) error
}

// TaskRepository implements TaskRepositoryInterface
type TaskRepository struct {
	DB *sql.DB
}

// PersistTask saves task in database.
// Task id is automatically incremented.
func (repo *TaskRepository) PersistTask(t *Task) error {
	query := "INSERT INTO TASK (STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID) VALUES (?,?,?,?)"

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

	_, err = stmt.Exec(t.Status, t.Title, t.AsigneeID, t.ChannelID)
	return err
}

// GetTaskByID returns reference to a task with this id.
// Return error if there is no such task.
func (repo *TaskRepository) GetTaskByID(ID int) (*Task, error) {
	query := "SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = ?"
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
	err = stmt.QueryRow(ID).Scan(&task.ID, &task.Status, &task.Title, &task.AsigneeID, &task.ChannelID)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAllInChannel accepts channel ID and returns all tasks in the specified channel
func (repo *TaskRepository) GetAllInChannel(channelID string) ([]*Task, error) {
	query := "SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = ?"
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
	rows, err := stmt.Query(channelID)
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, 0)
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.ID, &task.Status, &task.Title, &task.AsigneeID, &task.ChannelID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

// AssignTaskTo sets the assigneeID to assigneeID of the task with ID taskID. Returns error if there is no task with ID taskID.
func (repo *TaskRepository) AssignTaskTo(taskID int, assigneeID string) error {
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

	result, err := stmt.Exec(assigneeID, taskID)
	rows, err := result.RowsAffected()
	if rows != 1 {
		return ErrNoRowOrMoreThanOne
	}
	return err
}

// SetStatus sets the status to status of the task with ID taskID. Returns error if there is no task with ID taskID.
func (repo *TaskRepository) SetStatus(taskID int, status string) error {
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

	result, err := stmt.Exec(status, taskID)
	rows, err := result.RowsAffected()
	if rows != 1 {
		return ErrNoRowOrMoreThanOne
	}
	return err
}
