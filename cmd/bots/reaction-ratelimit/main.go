package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/constant"
	"github.com/andersfylling/snowflake/v3"
)

type keys struct {
	GuildAdmin   disgord.Snowflake
	GuildDefault disgord.Snowflake
}

func notARateLimitIssue(err error) bool {
	return !strings.Contains(err.Error(), "You are being rate limited.")
}

func setupKeys() *keys {
	keys := &keys{}

	str1 := os.Getenv(constant.DisgordTestGuildDefault)
	g1, err := disgord.GetSnowflake(str1)
	if err != nil {
		panic("missing default guild id")
	}
	keys.GuildDefault = g1

	str2 := os.Getenv(constant.DisgordTestGuildAdmin)
	g2, err := disgord.GetSnowflake(str2)
	if err != nil {
		panic("missing admin guild id")
	}
	keys.GuildAdmin = g2

	return keys
}

func main() {
	c, err := disgord.NewClient(&disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
	})
	if err != nil {
		panic(err)
	}

	msgID := snowflake.ID(540004388891262976)
	chanID := snowflake.ID(540004372231225344)

	e, err := c.GetGuildEmojis(486833611564253184).Execute()
	if err != nil {
		panic(err)
	}

	_ = c.DeleteAllReactions(chanID, msgID)
	for i := range e {
		start := time.Now()
		err = c.CreateReaction(chanID, msgID, e[i])
		if err != nil {
			fmt.Println(i, ": ", err)
			break
		}

		fmt.Println(time.Now().Sub(start).Seconds())
	}
}
