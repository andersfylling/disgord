package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/internal/constant"
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
	c := disgord.New(disgord.Config{
		BotToken: os.Getenv("DISGORD_TOKEN"),
		Logger:   disgord.DefaultLogger(true),
	})

	msgID := disgord.Snowflake(540519319814275089)
	chanID := disgord.Snowflake(540519296640614416)

	e, err := c.GetGuildEmojis(486833041486905345)
	if err != nil {
		panic(err)
	}

	_ = c.DeleteAllReactions(chanID, msgID)
	wg := sync.WaitGroup{}
	for i := range e {
		wg.Add(1)
		go func(index int) {
			start := time.Now()
			var msg string
			err = c.CreateReaction(chanID, msgID, e[index])
			if err != nil {
				msg = fmt.Sprint(index, ": ", err, " ### ")
			}

			fmt.Println(msg, time.Now().Sub(start).String())
			wg.Done()
		}(i)
	}

	wg.Wait()
}
