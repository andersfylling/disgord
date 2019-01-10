package disgord

import (
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/websocket"
)

func injectRandomEvents(t *testing.T, callback func(name string, evt interface{}) error) {
	events := map[string]interface{}{}

	iterate := func(t *testing.T, events map[string]interface{}) {
		var err error
		for k, v := range events {
			err = callback(k, v)
			if err != nil {
				t.Error("event{" + k + "}: " + err.Error())
				err = nil
			}
		}
	}

	// first wave, just empty content
	// looks for incorrect type casting
	events[EventReady] = &Ready{
		User: NewUser(),
	}
	events[EventChannelCreate] = &ChannelCreate{
		Channel: NewChannel(),
	}
	events[EventChannelDelete] = &ChannelDelete{
		Channel: NewChannel(),
	}
	events[EventGuildCreate] = &GuildCreate{
		Guild: NewGuild(),
	}
	events[EventGuildDelete] = &GuildDelete{
		UnavailableGuild: &GuildUnavailable{},
	}
	events[EventGuildBanRemove] = &GuildBanRemove{
		User: NewUser(),
	}
	events[EventGuildIntegrationsUpdate] = &GuildIntegrationsUpdate{}
	events[EventGuildMemberRemove] = &GuildMemberRemove{
		User: NewUser(),
	}
	events[EventGuildMembersChunk] = &GuildMembersChunk{}
	events[EventGuildRoleUpdate] = &GuildRoleUpdate{
		Role: NewRole(),
	}
	events[EventMessageCreate] = &MessageCreate{
		Message: NewMessage(),
	}
	events[EventMessageDelete] = &MessageDelete{}
	events[EventMessageReactionAdd] = &MessageReactionAdd{
		PartialEmoji: &Emoji{},
	}
	events[EventMessageReactionRemoveAll] = &MessageReactionRemoveAll{}
	events[EventTypingStart] = &TypingStart{}
	events[EventVoiceStateUpdate] = &VoiceStateUpdate{
		VoiceState: &VoiceState{},
	}
	events[EventWebhooksUpdate] = &WebhooksUpdate{}
	iterate(t, events)

}

func TestValidateUsername(t *testing.T) {
	var err error

	if err = ValidateUsername(""); err == nil {
		t.Error("expected empty error")
	}

	if err = ValidateUsername("a"); err == nil {
		t.Error("expected username to be too short")
	}

	if err = ValidateUsername("gk523526hdfgdfjdghlkjdhfglksjhdfg"); err == nil {
		t.Error("expected username to be too long")
	}

	if err = ValidateUsername("  anders"); err == nil {
		t.Error("expected username to have whitespace prefix error")
	}

	if err = ValidateUsername("anders  "); err == nil {
		t.Error("expected username to have whitespace suffix error")
	}

	if err = ValidateUsername("and  ers"); err == nil {
		t.Error("expected username to have excessive whitespaces error")
	}

	if err = ValidateUsername("@anders"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("#anders"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("and:ers"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("and```ers"); err == nil {
		t.Error("expected illegal char error")
	}

	if err = ValidateUsername("discordtag"); err == nil {
		t.Error("expected illegal username error")
	}

	if err = ValidateUsername("everyone"); err == nil {
		t.Error("expected illegal username error")
	}

	if err = ValidateUsername("here"); err == nil {
		t.Error("expected illegal username error")
	}
}

type mockerWSReceiveOnly struct {
	reading chan []byte
}

func (g *mockerWSReceiveOnly) Open(endpoint string, requestHeader http.Header) (err error) {
	return
}

func (g *mockerWSReceiveOnly) WriteJSON(v interface{}) (err error) {
	return
}

func (g *mockerWSReceiveOnly) Close() (err error) {
	return
}

func (g *mockerWSReceiveOnly) Read() (packet []byte, err error) {
	packet = <-g.reading
	return
}

func (g *mockerWSReceiveOnly) Disconnected() bool {
	return true
}

var _ websocket.Conn = (*mockerWSReceiveOnly)(nil)

var sink1 int = 1

// BenchmarkDiscordEventToHandler from the time Disgord gets the raw byte event data, to the event handler is triggered
func Benchmark1000DiscordEventToHandler_cacheDisabled(b *testing.B) {
	mocker := mockerWSReceiveOnly{
		reading: make(chan []byte),
	}
	// starts receiver and operation handler
	wsClient, wsShutdownChan := websocket.NewTestClient(nil, 0, &mocker)

	d := Client{
		shutdownChan: make(chan interface{}),
		config: &Config{
			DisableCache: true,
		},
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		evtDispatch: NewDispatch(false, 20),
	}
	d.shardManager.shards[0].ws = wsClient
	go d.eventHandler()

	seq := uint(1)
	wg := &sync.WaitGroup{}
	d.On(event.Ready, func(s Session, evt *Ready) {
		sink1++
		wg.Done()
	})

	f := func(mocker *mockerWSReceiveOnly, wg *sync.WaitGroup, seq *uint) {
		loops := 1000
		wg.Add(loops)
		for i := 0; i < loops; i++ {
			//evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`)
			evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{"v":6,"user_settings":{},"user":{"verified":true,"username":"Disgord tester","mfa_enabled":false,"id":"486832262592069632","email":null,"discriminator":"9338","bot":true,"avatar":null},"session_id":"d3954ff063fa8d387ec395fe65723624","relationships":[],"private_channels":[],"presences":[],"guilds":[{"unavailable":true,"id":"486833041486905345"},{"unavailable":true,"id":"486833611564253184"}],"_trace":["gateway-prd-main-kg6w","discord-sessions-prd-1-27"]}}`)
			mocker.reading <- evt
			*seq++
		}

		wg.Wait()
	}

	for i := 0; i < b.N; i++ {
		f(&mocker, wg, &seq)
	}

	close(d.shutdownChan)
	close(wsShutdownChan)
}
