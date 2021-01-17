package mysql

// refer to https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e blog
import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var p = &models.Person{
	Id:      1,
	SlackId: "SlackId",
}

func TestFindPersonById(t *testing.T) {
	db, mock := NewMock()
	mockService := &PersonService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, SLACK_ID FROM person WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "SLACK_ID"}).
		AddRow(p.Id, p.SlackId)
	mock.ExpectQuery(query).WithArgs(p.Id).WillReturnRows(rows)
	ta, err := mockService.FindPersonById(p.Id)
	assert.NotNil(t, ta)
	assert.NoError(t, err)
}

func TestCreatePerson(t *testing.T) {
	db, mock := NewMock()
	mockService := &PersonService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "INSERT INTO person \\(SLACK_ID\\) VALUES\\(\\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(p.SlackId).WillReturnResult(sqlmock.NewResult(0, 1))
	err := mockService.CreatePerson(p)
	assert.NoError(t, err)
}
