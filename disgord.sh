#!/usr/bin/env bash

VER="v1.0.1"

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
DEFAULT_K8S_SUPPORT="y"

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

read -p "kubernetes deployment script (y/n): " K8S_SUPPORT
if [[ ${K8S_SUPPORT} == "" ]]; then
    K8S_SUPPORT=${DEFAULT_K8S_SUPPORT}
fi


# Create the project
echo "Creating project"
mkdir -p ${PROJECT_PATH}
cd ${PROJECT_PATH}

echo 'package main

import (
	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"os"
)

// replyPongToPing is a handler that replies pong to ping messages
func replyPongToPing(s disgord.Session, data *disgord.MessageCreate) {
    msg := data.Message

    // whenever the message written is "ping", the bot replies "pong"
    if msg.Content == "ping" {
        msg.Reply(s, "pong")
    }
}

func main() {
    client := disgord.New(&disgord.Config{
        BotToken: os.Getenv("DISGORD_TOKEN"),
        Logger: disgord.DefaultLogger(false), // debug=false
    })
    defer client.StayConnectedUntilInterrupted()

	log, _ := std.NewLogFilter(client)
    filter, _ := std.NewMsgFilter(client)
	filter.SetPrefix("'${BOT_PREFIX}'")

    // create a handler and bind it to new message events
    // tip: read the documentation for std.CopyMsgEvt and understand why it is used here.
    client.On(disgord.EvtMessageCreate,
    	// middleware
    	filter.NotByBot,    // ignore bot messages
    	filter.HasPrefix,   // read original
    	log.LogMsg,         // log command message
    	std.CopyMsgEvt,     // read & copy original
    	filter.StripPrefix, // write copy
    	// handler
    	replyPongToPing) // handles copy
}
' >> main.go

echo "FROM andersfylling/disgord:v0.10 as builder
MAINTAINER https://github.com/andersfylling
WORKDIR /build
COPY . /build
RUN export GO111MODULE=on
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags \"-static\"' -o discordbot .

FROM gcr.io/distroless/base
WORKDIR /bot
COPY --from=builder /build/discordbot .
CMD [\"/bot/discordbot\"]
" >> Dockerfile

echo "# ${PROJECT_NAME}
Congratulations! You have successfully created a basic DisGord bot.

In order for your bot to start you must supply a environment variable with the name DISGORD_TOKEN that holds
the bot token you created in a Discord application or got from a friend.

eg. \"export DISGORD_TOKEN=si7fisdgfsfushgsjdf.sdfksgjyefs.dfgysyefs\"

A dockerfile has also been created for you if this is a preference. Note that you must supply the environment variable during run.

" >> README.md

if [[ ${K8S_SUPPORT} == "y" ]]; then
    echo "# kubernetes deployment file (GKE)
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: ${PROJECT_NAME}
spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: ${PROJECT_NAME}
    spec:
      containers:
        - name: ${PROJECT_NAME}
          image: *DOCKERHUBUSERNAME*/${PROJECT_NAME}-disgord:latest
          env:
            - name: DISGORD_TOKEN
              valueFrom:
                secretKeyRef: # needs to be manually created
                  name: discord-tokens
                  key: ${PROJECT_NAME}
" >> deployment.yaml
    echo "> remember to change *DOCKERHUBUSERNAME* in deployment.yaml with your actual username on hub.docker.com or other hosting site."
    echo "> note that deployment.yaml is just a suggestion. You will still need to manually edit it ot make it work."
    # TODO: ask for docker repository url
fi

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
    # TODO: ask for git service, then username, create repository if it does not exist and set upstream,
    #  commit changes and push.
fi