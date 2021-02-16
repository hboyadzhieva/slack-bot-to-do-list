package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var task = &Task{
	ID:        1,
	Status:    StatusOpen,
	Title:     "Manual test of ui",
	AsigneeID: "U123",
	ChannelID: "C123",
}

func TestPersistTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO TASK \\(STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID\\) VALUES \\(\\?,\\?,\\?,\\?\\)").ExpectExec().WithArgs(task.Status, task.Title, task.AsigneeID, task.ChannelID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.PersistTask(task)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetTaskByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"}).
		AddRow(task.ID, task.Status, task.Title, task.AsigneeID, task.ChannelID)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = \\?").ExpectQuery().WithArgs(task.ID).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetTaskByID(task.ID)
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.EqualValues(t, res, task)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestGetTaskByIDError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "ASIGNEE_ID", "CHANNEL_ID"})
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE ID = \\?").ExpectQuery().WithArgs(task.ID).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetTaskByID(task.ID)
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
		AddRow(task.ID, task.Status, task.Title, task.AsigneeID, task.ChannelID).
		AddRow(2, task.Status, task.Title, task.AsigneeID, task.ChannelID)
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("SELECT ID, STATUS, TITLE, ASIGNEE_ID, CHANNEL_ID FROM TASK WHERE CHANNEL_ID = \\?").ExpectQuery().WithArgs(task.ChannelID).WillReturnRows(rows)
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	res, err := mockService.GetAllInChannel(task.ChannelID)
	if assert.NoError(t, err) {
		assert.NotNil(t, res)
		assert.Equal(t, 2, len(res))
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
	mock.ExpectPrepare("UPDATE TASK SET ASIGNEE_ID = \\? WHERE ID = \\?").ExpectExec().WithArgs(task.AsigneeID, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.AssignTaskTo(task.ID, task.AsigneeID)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestAssignTaskErrNoRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE TASK SET ASIGNEE_ID = \\? WHERE ID = \\?").ExpectExec().WithArgs(task.AsigneeID, 57).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.AssignTaskTo(57, task.AsigneeID)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoRowOrMoreThanOne)
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
	mock.ExpectPrepare("UPDATE TASK SET STATUS = \\? WHERE ID = \\?").ExpectExec().WithArgs(StatusInProgress, task.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.SetStatus(task.ID, StatusInProgress)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestSetStatusErrNoRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("UPDATE TASK SET STATUS = \\? WHERE ID = \\?").ExpectExec().WithArgs(StatusInProgress, task.ID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.SetStatus(task.ID, StatusInProgress)
	assert.Error(t, err)
	assert.Equal(t, err, ErrNoRowOrMoreThanOne)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}
