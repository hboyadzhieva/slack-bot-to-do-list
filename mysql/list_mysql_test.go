package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var list = &models.List{
	Id:        4,
	Title:     "QA Tasks",
	ChannelId: 2,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Error: '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestFindListById(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, TITLE, CHANNEL FROM list WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "TITLE", "CHANNEL"}).
		AddRow(list.Id, list.Title, list.ChannelId)
	mock.ExpectQuery(query).WithArgs(list.Id).WillReturnRows(rows)
	l, err := mockService.FindListById(list.Id)
	assert.NotNil(t, l)
	assert.NoError(t, err)
}

func TestFindListByIdError(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, TITLE, CHANNEL FROM list WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "TITLE", "CHANNEL"})
	mock.ExpectQuery(query).WithArgs(list.Id).WillReturnRows(rows)
	l, err := mockService.FindListById(list.Id)
	assert.Empty(t, l)
	assert.Error(t, err)
}

func TestAllLists(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, TITLE, CHANNEL FROM list"
	rows := sqlmock.NewRows([]string{"ID", "TITLE", "CHANNEL"}).AddRow(list.Id, list.Title, list.ChannelId)
	mock.ExpectQuery(query).WillReturnRows(rows)
	lists, err := mockService.AllLists()
	assert.NotEmpty(t, lists)
	assert.NoError(t, err)
	assert.Len(t, lists, 1)
}

func TestAllListsEmpty(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, TITLE, CHANNEL FROM list"
	rows := sqlmock.NewRows([]string{"ID", "TITLE", "CHANNEL"})
	mock.ExpectQuery(query).WillReturnRows(rows)
	lists, err := mockService.AllLists()
	assert.Empty(t, lists)
	assert.NoError(t, err)
}

func TestAllListsInChannel(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, TITLE, CHANNEL FROM list WHERE CHANNEL=\\?"
	rows := sqlmock.NewRows([]string{"ID", "TITLE", "CHANNEL"}).AddRow(1, "Title1", 2).AddRow(2, "Title2", 2)
	mock.ExpectQuery(query).WithArgs(2).WillReturnRows(rows)
	lists, err := mockService.AllListsInChannel(&models.Channel{Id: 2, SlackId: "slackId"})
	assert.NotEmpty(t, lists)
	assert.NoError(t, err)
	assert.Len(t, lists, 2)
}

func TestCreateList(t *testing.T) {
	db, mock := NewMock()
	mockService := &ListService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "INSERT INTO list \\(TITLE, CHANNEL\\) VALUES\\(\\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(list.Title, list.ChannelId).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.CreateList(list)
	assert.NoError(t, err)
}
