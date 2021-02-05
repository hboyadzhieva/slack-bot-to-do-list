package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var task = Task{
	Id:          1,
	Status:      "Open",
	Title:       "Manual test of ui",
	Description: "Test manually home and about page",
	AsigneeId:   "U123",
	ChannelId:   "C123",
}

func TestPersistTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	defer db.Close()
	mock.MatchExpectationsInOrder(true)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO TASK \\(STATUS, TITLE, DESCRIPTION, ASIGNEE_ID, CHANNEL_ID\\) VALUES \\(\\?,\\?,\\?,\\?,\\?\\)").ExpectExec().WithArgs(task.Status, task.Title, task.Description, task.AsigneeId, task.ChannelId).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mockService := &TaskRepository{db}
	err = mockService.PersistTask(&task)
	assert.NoError(t, err)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

/*
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
*/
