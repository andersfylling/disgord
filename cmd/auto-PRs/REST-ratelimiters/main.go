package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func main() {
	client, err := disgord.NewClient(&disgord.Config{
		BotToken:     "NDg2ODMyMjYyNTkyMDY5NjMy.XVbNbQ.sZMxoQw1RxA9GnqU8PizhPMLIKE",
		DisableCache: true,
	})
	if err != nil {
		panic(err)
	}

	requests := prepareRequests(client)
	for _, f := range requests {
		for i := 0; i < 10; i++ {
			if err := f(); err != nil {
				fmt.Println("ERROR", err)
			}
		}
	}

	for k, v := range client.RESTBucketsRelations() {
		fmt.Println("\"" + k + "\":\"" + v + "\"")
	}
}
