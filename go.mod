module github.com/hboyadzhieva/slack-bot-to-do-list

go 1.15

replace github.com/hboyadzhieva/slack-bot-to-do-list/mysql => ./mysql

replace github.com/hboyadzhieva/slack-bot-to-do-list/tododo => ./tododo

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Skellington-Closet/slack-mock v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/joho/godotenv v1.3.0
	github.com/nlopes/slack v0.6.0
	github.com/stretchr/testify v1.6.1
)
