package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andersfylling/disgord/discordws"
	"github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("DISGORD_TOKEN")
	if token == "" {
		token = "NDA0NzY4MzUxMjgyMzk3MTg1.DUapbA.9uF6sXXIiOs7NzWC-nYdBz6Oaos"
		// panic("Missing disgord token in env var: DISGORD_TOKEN")
	}
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	d := discordws.NewRequiredClient(&discordws.Config{
		Token:        token,
		DAPIVersion:  6,
		DAPIEncoding: discordws.EncodingJSON,
	})
	err := d.Connect()
	if err != nil {
		panic(err)
	}
	<-termSignal
	fmt.Println("Closing connection")
	err = d.Disconnect()
	if err != nil {
		logrus.Fatal(err)
	}
}
