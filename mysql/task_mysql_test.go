package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var task = &models.Task{
	Id:          1,
	Status:      "Open",
	Title:       "Manual test of ui",
	Description: "Test manually home and about page",
	ListId:      4,
	AsigneeId:   3,
}

func TestFindTaskById(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "DESCRIPTION", "LIST", "ASIGNEE"}).
		AddRow(task.Id, task.Status, task.Title, task.Description, task.ListId, task.AsigneeId)
	mock.ExpectQuery(query).WithArgs(task.Id).WillReturnRows(rows)
	ta, err := mockService.FindTaskById(task.Id)
	assert.NotNil(t, ta)
	assert.NoError(t, err)
}

func TestFindTaskByIdError(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "DESCRIPTION", "LIST", "ASIGNEE"})
	mock.ExpectQuery(query).WithArgs(task.Id).WillReturnRows(rows)
	ta, err := mockService.FindTaskById(task.Id)
	assert.Empty(t, ta)
	assert.Error(t, err)
}

func TestAllTasksInList(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE FROM list WHERE LIST=\\?"
	rows := sqlmock.NewRows([]string{"ID", "STATUS", "TITLE", "DESCRIPTION", "LIST", "ASIGNEE"}).
		AddRow(task.Id, task.Status, task.Title, task.Description, task.ListId, task.AsigneeId)
	mock.ExpectQuery(query).WithArgs(task.ListId).WillReturnRows(rows)
	tasks, err := mockService.AllTasksInList(&models.List{Id: task.ListId})
	assert.NotEmpty(t, tasks)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)
}

func TestCreateTask(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "INSERT INTO list \\(STATUS, TITLE, DESCRIPTION, LIST, ASIGNEE\\) VALUES\\(\\?, \\?, \\?, \\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(task.Status, task.Title, task.Description, task.ListId, task.AsigneeId).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.CreateTask(task)
	assert.NoError(t, err)
}

func TestAssignTask(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()
	p := &models.Person{Id: 1, SlackId: "a1"}
	query := "UPDATE list SET ASIGNEE = \\? WHERE ID = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(task.Id, p.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.AssignTask(task, p)
	assert.NoError(t, err)
}

func TestSetStatus(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()

	status := "Done"
	query := "UPDATE list SET STATUS = \\? WHERE ID = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(status, task.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.SetStatus(task, status)
	assert.NoError(t, err)
}

func TestAddTaskToList(t *testing.T) {
	db, mock := NewMock()
	mockService := &TaskService{db}
	defer func() {
		mockService.DB.Close()
	}()
	l := &models.List{Id: 1}
	query := "UPDATE list SET LIST = \\? WHERE ID = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(task.Id, l.Id).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.AddTaskToList(task, l)
	assert.NoError(t, err)
}
