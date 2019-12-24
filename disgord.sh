#!/usr/bin/env bash

VER="v2.0.0"

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

curl -o Dockerfile https://github.com/andersfylling/disgord/cmd/script/Dockerfile
curl -o main.go https://github.com/andersfylling/disgord/cmd/script/main.go
sed -i -e "s/REPLACE_ME/${BOT_PREFIX}/g" main.go

echo "# ${PROJECT_NAME}

## Congratulations!
You have successfully created a basic DisGord bot.

In order for your bot to start you must supply a environment variable with the name DISGORD_TOKEN that holds
the bot token you created in a Discord application or got from a friend.
See tutorial here to find/create the token: https://github.com/andersfylling/disgord/wiki/Get-bot-token-and-add-it-to-a-server

eg. \"export DISGORD_TOKEN=si7fisdgfsfushgsjdf.sdfksgjyefs.dfgysyefs\"

A dockerfile has also been created to build a proper production image. Note that you must supply the environment variable DISGORD_TOKEN when running the container.

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
    # TODO: ask for script repository url
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