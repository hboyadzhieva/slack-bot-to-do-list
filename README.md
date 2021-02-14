# Slack bot ToDo list
FMI Golang Course Project

## Overview
This is a simple Slack bot which uses Slack slash commands written in Golang. The commands can be used in a Slack channel by users to add tasks in the ToDo list of the channel, assign them to specific users and update their progress.
The task data is kept in MySQL database.

### Commands
- */tododo-help* - show all available commands
- */tododo-add [task]* - add a task to the list
- */tododo-show* - show all tasks in the list, the assignees and progress
- */tododo-assing [task id] [@user]* - assign a task to a user in the channel
- */tododo-start [task id]* - start progress on a task
- */tododo-done [task id]* - finish a task

## Local build and install

1. Get packages and install dependencies

     `go get -t github.com/hboyadzhieva/slack-bot-to-do-list`
2. Install 
    
    `go install github.com/hboyadzhieva/slack-bot-to-do-list`

3. Start docker container for MySQL
     
     Go to $GOPATH/src/github.com/hboyadzhieva/slack-bot-to-do-list and execute:
    
    `docker-compose up -d`
5. Prepare Slack bot and Slack Slash commands

    - For test on local machine install [ngrok](https://ngrok.com/)
    - Run ngrok http [port] and copy the url that forwards to localhost:[port]
    - Go to [https://api.slack.com/apps/](https://api.slack.com/apps/) and create a new app
    - Open your new app and go to Feature -> Slash commands
    - Create slash commands and in the field of Request URL paste the url from ngrok and append /tododo in the end for every command
    - Need to create commands */tododo-help*, */tododo-show*, */tododo-add*, */tododo-assign*, */tododo-start*, */tododo-done*
    - Install the app to a workspace of your choice
    <br/>
    <img alt="commands image" src="https://github.com/hboyadzhieva/slack-bot-to-do-list/blob/main/img/commands.png" width="500" height="500">
    
6. Set environment variable
    
    - From [slack-api](https://api.slack.com/apps/) go to your app -> Basic Information -> App Credentials, copy Verification Token and set environment variable SLACK_VERIFICATION_TOKEN to the value
    - for Windows 
      `set SLACK_VERIFICATION_TOKEN=<your verification token>`
    - for Linux/Mac
      `EXPORT SLACK_VERIFICATION_TOKEN="<your verification token>"`
      
7. Run slack-bot-to-do-list from $GOPATH/bin and type commands in a Slack channel
      
