#!/usr/bin/env bash

VER="v0.0.0"

echo "
# # # # # # # # # # # # # # # # #
#        DisGord utility        #
#            ${VER}             #
# # # # # # # # # # # # # # # # #
* Simple script to assist you   *
* with creating a basic DisGord *
* bot.                          *
* * * * * * * * * * * * * * * * *
"

DEFAULT_PROJECT_NAME="mybot"
DEFAULT_BOT_PREFIX="!"
DEFAULT_PROJECT_PATH="."
DEFAULT_GIT_SUPPORT="y"

read -p "project name (${DEFAULT_PROJECT_NAME}): " PROJECT_NAME
if [[ ${PROJECT_NAME} == "" ]]; then
    PROJECT_NAME=${DEFAULT_PROJECT_NAME}
fi

read -p "bot prefix (${DEFAULT_BOT_PREFIX}): " BOT_PREFIX
if [[ ${BOT_PREFIX} == "" ]]; then
    BOT_PREFIX=${DEFAULT_BOT_PREFIX}
fi

read -p "project path (${DEFAULT_PROJECT_PATH}): " PROJECT_PATH
if [[ ${PROJECT_PATH} == "" ]]; then
    PROJECT_PATH=${DEFAULT_PROJECT_PATH}
fi
if [[ ${PROJECT_PATH} == "." ]]; then
    PROJECT_PATH=""
fi
if [[ ${PROJECT_PATH} == /* ]]; then
    PROJECT_PATH="${PROJECT_PATH}"
else
    PROJECT_PATH="$(pwd)/${PROJECT_PATH}"
fi
PROJECT_PATH="${PROJECT_PATH}/${PROJECT_NAME}"

read -p "git support (y/n): " GIT_SUPPORT
if [[ ${GIT_SUPPORT} == "" ]]; then
    GIT_SUPPORT=${DEFAULT_GIT_SUPPORT}
fi



# Create the project
echo "Creating project"
mkdir -p ${PROJECT_PATH}
cd ${PROJECT_PATH}

echo 'package main

import (
	"github.com/andersfylling/disgord"
	"os"
)

func replyPongToPing(session disgord.Session, data *disgord.MessageCreate) {
    msg := data.Message

    // whenever the message written is "ping", the bot replies "pong"
    if msg.Content == "ping" {
        msg.RespondString(session, "pong")
    }
}

func main() {
	var err error

	botConfig := &disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
        Logger: disgord.DefaultLogger(false), // optional logging, debug=false
    }

    // create a Disgord session
    var client *disgord.Client
    if client, err = disgord.NewClient(botConfig); err != nil {
        panic(err)
    }

    // create a handler and bind it to new message events
    client.On(disgord.EventMessageCreate, replyPongToPing)

    // connect to the discord gateway to receive events
    if err = client.Connect(); err != nil {
        panic(err)
    }

    // Keep the socket connection alive, until you terminate the application (eg. Ctrl + C)
    if err = client.DisconnectOnInterrupt(); err != nil {
    	botConfig.Logger.Error(err) // reuse the logger from DisGord
    }
}
' >> main.go

echo "FROM golang:1.11.5 as builder

WORKDIR /build
COPY . /build

RUN export GO111MODULE=on

RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags \"-static\"' -o discordbot .


FROM scratch
WORKDIR /bot

COPY . /bot
COPY --from=builder /build/discordbot .
RUN chmod +x /bot/discordbot

CMD [\"./discordbot\"]
" >> Dockerfile

echo "# ${PROJECT_NAME}
Congratulations! You have successfully created a basic DisGord bot.

In order for your bot to start you must supply a environment variable with the name DISGORD_TOKEN that holds
the bot token you created in a Discord application or got from a friend.

eg. \"export DISGORD_TOKEN=si7fisdgfsfushgsjdf.sdfksgjyefs.dfgysyefs\"

A dockerfile has also been created for you if this is a preference. Note that you must supply the environment variable during run.

" >> README.md

if [[ -z ${GO111MODULE} ]] || [[ ${GO111MODULE} == "off" ]]; then
    export GO111MODULE="auto"
fi

# create go.mod
go mod init "${PROJECT_NAME}"
go build .
go mod tidy

if [[ ${GIT_SUPPORT} == "y" ]]; then
    git init
    git add .
    git status
fi