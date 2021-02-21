// +build !integration

package std

import (
	"context"
	"fmt"
	"testing"

	"github.com/andersfylling/disgord"
)

type clientRESTMock struct {
	id disgord.Snowflake
}

// var _ disgord.RESTMethods = (*clientRESTMock)(nil)

func (c *clientRESTMock) CurrentUser() disgord.CurrentUserQueryBuilder {
	return &clientRESTMock_currentUser{id: c.id}
}

type clientRESTMock_currentUser struct {
	disgord.CurrentUserQueryBuilder
	id disgord.Snowflake
}

func (c clientRESTMock_currentUser) WithContext(_ context.Context) disgord.CurrentUserQueryBuilder {
	return &c
}

func (c clientRESTMock_currentUser) Get(_ ...disgord.Flag) (*disgord.User, error) {
	return &disgord.User{ID: c.id}, nil
}

func TestNewMsgFilter(t *testing.T) {
	var botID disgord.Snowflake = 123
	filter, err := newMsgFilter(context.Background(), &clientRESTMock{id: botID})
	if err != nil {
		t.Fatal(err)
	}

	if filter.botID != botID {
		t.Errorf("expected filter to have the same id as client. Got %d, wants %d", filter.botID, botID)
	}
}

func TestMsgFilter_NotByBot(t *testing.T) {
	var botID disgord.Snowflake = 123

	messageFromBot := &disgord.Message{
		Author: &disgord.User{Bot: true},
	}

	messageNotFromBot := &disgord.Message{
		Author: &disgord.User{Bot: false},
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{"MessageCreate_FromBot", &disgord.MessageCreate{Message: messageFromBot}, false},
		{"MessageUpdate_FromBot", &disgord.MessageUpdate{Message: messageFromBot}, false},
		{"MessageCreate_NotBot", &disgord.MessageCreate{Message: messageNotFromBot}, true},
		{"MessageUpdate_NotBot", &disgord.MessageUpdate{Message: messageNotFromBot}, true},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.NotByBot(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_IsByBot(t *testing.T) {
	var botID disgord.Snowflake = 123

	messageFromBot := &disgord.Message{
		Author: &disgord.User{Bot: true},
	}

	messageNotFromBot := &disgord.Message{
		Author: &disgord.User{Bot: false},
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{"MessageCreate_FromBot", &disgord.MessageCreate{Message: messageFromBot}, true},
		{"MessageUpdate_FromBot", &disgord.MessageUpdate{Message: messageFromBot}, true},
		{"MessageCreate_NotBot", &disgord.MessageCreate{Message: messageNotFromBot}, false},
		{"MessageUpdate_NotBot", &disgord.MessageUpdate{Message: messageNotFromBot}, false},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.IsByBot(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_NotByWebhook(t *testing.T) {
	var webhookID disgord.Snowflake = 456
	var botID disgord.Snowflake = 123

	messageWithWebhookID := &disgord.Message{
		Author:    &disgord.User{},
		WebhookID: webhookID,
	}

	messageWithoutWebhookID := &disgord.Message{
		Author:    &disgord.User{},
		WebhookID: 0,
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{"MessageCreate_FromWebhook", &disgord.MessageCreate{Message: messageWithWebhookID}, false},
		{"MessageUpdate_FromWebhook", &disgord.MessageUpdate{Message: messageWithWebhookID}, false},
		{"MessageCreate_NotWebhook", &disgord.MessageCreate{Message: messageWithoutWebhookID}, true},
		{"MessageUpdate_NotWebhook", &disgord.MessageUpdate{Message: messageWithoutWebhookID}, true},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.NotByWebhook(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_IsByWebhook(t *testing.T) {
	var webhookID disgord.Snowflake = 456
	var botID disgord.Snowflake = 123

	messageWithWebhookID := &disgord.Message{
		Author:    &disgord.User{},
		WebhookID: webhookID,
	}

	messageWithoutWebhookID := &disgord.Message{
		Author:    &disgord.User{},
		WebhookID: 0,
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{"MessageCreate_FromWebhook", &disgord.MessageCreate{Message: messageWithWebhookID}, true},
		{"MessageUpdate_FromWebhook", &disgord.MessageUpdate{Message: messageWithWebhookID}, true},
		{"MessageCreate_NotWebhook", &disgord.MessageCreate{Message: messageWithoutWebhookID}, false},
		{"MessageUpdate_NotWebhook", &disgord.MessageUpdate{Message: messageWithoutWebhookID}, false},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.IsByWebhook(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_ContainsBotMention(t *testing.T) {
	var botID disgord.Snowflake = 123
	var wrongBotID disgord.Snowflake = 126

	messageCreate := func(content string) interface{} {
		return &disgord.MessageCreate{
			Message: &disgord.Message{Content: content},
		}
	}

	messageUpdate := func(content string) interface{} {
		return &disgord.MessageUpdate{
			Message: &disgord.Message{Content: content},
		}
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{
			"MessageCreate_BotHello",
			messageCreate(fmt.Sprintf("<@%s> hello", botID.String())),
			true,
		},
		{
			"MessageCreate_BotHello_nickname",
			messageCreate(fmt.Sprintf("<@!%s> hello", botID.String())),
			true,
		},
		{
			"MessageUpdate_BotHello",
			messageUpdate(fmt.Sprintf("<@%s> hello", botID.String())),
			true,
		},
		{
			"MessageCreate_WrongBotHello",
			messageCreate(fmt.Sprintf("<@%s> hello", wrongBotID.String())),
			false,
		},
		{
			"MessageUpdate_WrongBotHello",
			messageUpdate(fmt.Sprintf("<@%s> hello", wrongBotID.String())),
			false,
		},
		{
			"MessageCreate_BotHellWithPrefix",
			messageCreate(fmt.Sprintf("diff prefix <@%s> hello", botID.String())),
			true,
		},
		{
			"MessageUpdate_BotHelloWithPrefix",
			messageUpdate(fmt.Sprintf("diff prefix <@%s> hello", botID.String())),
			true,
		},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.ContainsBotMention(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_HasBotMentionPrefix(t *testing.T) {
	var botID disgord.Snowflake = 123
	var wrongBotID disgord.Snowflake = 126

	messageCreate := func(content string) interface{} {
		return &disgord.MessageCreate{
			Message: &disgord.Message{Content: content},
		}
	}

	messageUpdate := func(content string) interface{} {
		return &disgord.MessageUpdate{
			Message: &disgord.Message{Content: content},
		}
	}

	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{
			"MessageCreate_BotHello",
			messageCreate(fmt.Sprintf("<@%s> hello", botID.String())),
			true,
		},
		{
			"MessageCreate_BotHello_nickname",
			messageCreate(fmt.Sprintf("<@!%s> hello", botID.String())),
			true,
		},
		{
			"MessageUpdate_BotHello",
			messageUpdate(fmt.Sprintf("<@%s> hello", botID.String())),
			true,
		},
		{
			"MessageCreate_WrongBotHello",
			messageCreate(fmt.Sprintf("<@%s> hello", wrongBotID.String())),
			false,
		},
		{
			"MessageUpdate_WrongBotHello",
			messageUpdate(fmt.Sprintf("<@%s> hello", wrongBotID.String())),
			false,
		},
		{
			"MessageCreate_BotHellWithDiffPrefix",
			messageCreate(fmt.Sprintf("diff prefix <@%s> hello", botID.String())),
			false,
		},
		{
			"MessageUpdate_BotHelloWithDiffPrefix",
			messageUpdate(fmt.Sprintf("diff prefix <@%s> hello", botID.String())),
			false,
		},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{id: botID})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.HasBotMentionPrefix(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_SetPrefix(t *testing.T) {
	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{})
	if filter.prefix != "" {
		t.Fatal("expected prefix to be empty")
	}

	filter.SetPrefix("!")
	if filter.prefix != "!" {
		t.Errorf("wrong prefix. Got %s, wants %s", filter.prefix, "!")
	}
}

func TestMsgFilter_HasPrefix(t *testing.T) {
	messageCreate := func(content string) interface{} {
		return &disgord.MessageCreate{
			Message: &disgord.Message{Content: content},
		}
	}

	messageUpdate := func(content string) interface{} {
		return &disgord.MessageUpdate{
			Message: &disgord.Message{Content: content},
		}
	}

	prefix := "!!"
	testCases := []struct {
		name              string
		evt               interface{}
		shouldPassThrough bool
	}{
		{"MessageCreate_CorrectPrefix", messageCreate(prefix + "hello"), true},
		{"MessageUpdate_CorrectPrefix", messageUpdate(prefix + "hello"), true},
		{"MessageCreate_WrongPrefix", messageCreate("diff prefix " + prefix + "hello"), false},
		{"MessageUpdate_WrongPrefix", messageUpdate("diff prefix " + prefix + "hello"), false},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{})
	filter.SetPrefix(prefix)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.HasPrefix(tc.evt)
			if tc.shouldPassThrough && result == nil {
				t.Error("expected to passthrough")
			}
			if !tc.shouldPassThrough && result != nil {
				t.Error("expected a filter match")
			}
		})
	}
}

func TestMsgFilter_StripPrefix(t *testing.T) {
	messageCreate := func(content string) interface{} {
		return &disgord.MessageCreate{
			Message: &disgord.Message{Content: content},
		}
	}

	messageUpdate := func(content string) interface{} {
		return &disgord.MessageUpdate{
			Message: &disgord.Message{Content: content},
		}
	}

	prefix := "!!"
	testCases := []struct {
		name string
		evt  interface{}
	}{
		{"MessageCreate", messageCreate(prefix + "hello")},
		{"MessageUpdate", messageUpdate(prefix + "hello")},
		{"MessageUpdate", messageUpdate("   " + prefix + "hello")},
		{"MessageCreate", messageCreate("   " + prefix + "  hello")},
		{"MessageUpdate", messageUpdate(prefix + "  hello")},
		{"MessageUpdate", messageUpdate("  hello")},
	}

	filter, _ := newMsgFilter(context.Background(), &clientRESTMock{})
	filter.SetPrefix(prefix)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filter.StripPrefix(tc.evt)
			if filter.HasPrefix(result) != nil {
				t.Error("Did not strip prefix off message")
			}
		})
	}
}
