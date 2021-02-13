package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var task = &Task{
	Id:        1,
	Status:    StatusOpen,
	Title:     "Manual test of ui",
	AsigneeId: "U123",
	ChannelId: "C123",
}

func TestPersistTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO TASK \\(STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID\\) VALUES \\(\\?,\\?,\\?,\\?\\)").ExpectExec().WithArgs(task.Status, task.Title, task.AsigneeId, task.ChannelId).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.PersistTask(task)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetTaskById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"}).
		AddRow(task.Id, task.Status, task.Title, task.AsigneeId, task.ChannelId)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = \\?").ExpectQuery().WithArgs(task.Id).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetTaskById(task.Id)
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.EqualValues(t, res, task)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetTaskByIdError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"})
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = \\?").ExpectQuery().WithArgs(task.Id).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetTaskById(task.Id)
	expectedError := sql.ErrNoRows
	if assert.Error(t, err) {
		assert.Nil(t, res)
		assert.Equal(t, expectedError, err)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetAllInChannel(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"}).
		AddRow(task.Id, task.Status, task.Title, task.AsigneeId, task.ChannelId).
		AddRow(2, task.Status, task.Title, task.AsigneeId, task.ChannelId)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = \\?").ExpectQuery().WithArgs(task.ChannelId).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetAllInChannel(task.ChannelId)
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetAllInChannelAssignedTo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"}).
		AddRow(task.Id, task.Status, task.Title, task.AsigneeId, task.ChannelId)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = \\? AND ASIGNEE_ID = \\?").ExpectQuery().WithArgs(task.ChannelId, task.AsigneeId).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetAllInChannelAssignedTo(task.ChannelId, task.AsigneeId)
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.Equal(t, 1, len(res))
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestAssignTaskTo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE TASK SET ASIGNEE_ID = \\? WHERE ID = \\?").ExpectExec().WithArgs(task.AsigneeId, task.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.AssignTaskTo(task.Id, task.AsigneeId)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSetStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE TASK SET STATUS = \\? WHERE ID = \\?").ExpectExec().WithArgs(StatusInProgress, task.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.SetStatus(task.Id, StatusInProgress)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}
