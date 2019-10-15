package main

import (
	"time"

	"github.com/andersfylling/disgord"
)

const guildID = 486833041486905345
const channelID = 486833041486905347
const msgID = 492360894244716554

func prepareRequests(client *disgord.Client) (rs []func() error) {
	// RESTAuditLogs
	rs = append(rs, func() error {
		_, err := client.GetGuildAuditLogs(guildID, disgord.IgnoreCache).Execute()
		return err
	})

	// RESTMessage
	rs = append(rs, func() error {
		p := &disgord.GetMessagesParams{
			Limit: 1,
		}
		_, err := client.GetMessages(channelID, p, disgord.IgnoreCache)
		return err
	})
	rs = append(rs, func() error {
		_, err := client.GetMessage(channelID, msgID, disgord.IgnoreCache)
		return err
	})
	rs = append(rs, func() error {
		<-time.After(50 * time.Millisecond)
		p := &disgord.CreateMessageParams{
			Content: "test",
		}
		msg, err := client.CreateMessage(channelID, p, disgord.IgnoreCache)
		if err != nil {
			return err
		}

		_, err = client.UpdateMessage(channelID, msg.ID).SetContent("asdas").Execute()
		if err != nil {
			return err
		}

		err = client.DeleteMessage(channelID, msg.ID)
		return err
	})
	rs = append(rs, func() error {
		<-time.After(50 * time.Millisecond)
		p := &disgord.CreateMessageParams{
			Content: "test2",
		}
		msg, err := client.CreateMessage(channelID, p, disgord.IgnoreCache)
		if err != nil {
			return err
		}

		p2 := &disgord.DeleteMessagesParams{
			Messages: []disgord.Snowflake{msg.ID},
		}
		err = client.DeleteMessages(channelID, p2)
		return err
	})

	// RESTReaction

	return rs
}
