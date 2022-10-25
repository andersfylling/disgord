//go:build !integration
// +build !integration

package disgord

import (
	"context"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/andersfylling/disgord/json"

	"github.com/andersfylling/disgord/internal/gateway"
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
	events[EvtReady] = &Ready{
		User: &User{},
	}
	events[EvtChannelCreate] = &ChannelCreate{
		Channel: &Channel{},
	}
	events[EvtChannelDelete] = &ChannelDelete{
		Channel: &Channel{},
	}
	events[EvtGuildCreate] = &GuildCreate{
		Guild: &Guild{},
	}
	events[EvtGuildDelete] = &GuildDelete{
		UnavailableGuild: &GuildUnavailable{},
	}
	events[EvtGuildBanRemove] = &GuildBanRemove{
		User: &User{},
	}
	events[EvtGuildIntegrationsUpdate] = &GuildIntegrationsUpdate{}
	events[EvtGuildMemberRemove] = &GuildMemberRemove{
		User: &User{},
	}
	events[EvtGuildMembersChunk] = &GuildMembersChunk{}
	events[EvtGuildRoleUpdate] = &GuildRoleUpdate{
		Role: &Role{},
	}
	events[EvtMessageCreate] = &MessageCreate{
		Message: &Message{},
	}
	events[EvtMessageDelete] = &MessageDelete{}
	events[EvtMessageReactionAdd] = &MessageReactionAdd{
		PartialEmoji: &Emoji{},
	}
	events[EvtMessageReactionRemoveAll] = &MessageReactionRemoveAll{}
	events[EvtTypingStart] = &TypingStart{}
	events[EvtVoiceStateUpdate] = &VoiceStateUpdate{
		VoiceState: &VoiceState{},
	}
	events[EvtWebhooksUpdate] = &WebhooksUpdate{}
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

func (g *mockerWSReceiveOnly) Open(ctx context.Context, endpoint string, requestHeader http.Header) (err error) {
	return
}

func (g *mockerWSReceiveOnly) WriteJSON(v interface{}) (err error) {
	return
}

func (g *mockerWSReceiveOnly) Close() (err error) {
	return
}

func (g *mockerWSReceiveOnly) Read(ctx context.Context) (packet []byte, err error) {
	packet = <-g.reading
	return
}

func (g *mockerWSReceiveOnly) Disconnected() bool {
	return true
}

var _ gateway.Conn = (*mockerWSReceiveOnly)(nil)

var sink1 int = 1

// TODO
// BenchmarkDiscordEventToHandler from the time Disgord gets the raw byte event data, to the event handler is triggered
//func Benchmark1000DiscordEventToHandler_cacheDisabled(b *testing.B) {
//	mocker := mockerWSReceiveOnly{
//		reading: make(chan []byte),
//	}
//	// starts receiver and operation handler
//	wsClient, wsShutdownChan := websocket.NewTestClient(nil, 0, &mocker)
//
//	d := Client{
//		shutdownChan: make(chan interface{}),
//		config: &EvtConfig{
//			IgnoreCache: true,
//		},
//		httpClient: &http.Client{
//			Timeout: time.Second * 10,
//		},
//		evtDispatch: NewDispatch(false, 20),
//	}
//	d.shardManager.shards[0].ws = wsClient
//	go d.eventHandler()
//
//	seq := uint(1)
//	wg := &sync.WaitGroup{}
//	d.On(event.Ready, func(s Session, evt *Ready) {
//		sink1++
//		wg.Done()
//	})
//
//	f := func(mocker *mockerWSReceiveOnly, wg *sync.WaitGroup, seq *uint) {
//		loops := 1000
//		wg.Add(loops)
//		for i := 0; i < loops; i++ {
//			//evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{}}`)
//			evt := []byte(`{"t":"READY","s":` + strconv.Itoa(int(*seq)) + `,"op":0,"d":{"v":6,"user_settings":{},"user":{"verified":true,"username":"Disgord tester","mfa_enabled":false,"id":"486832262592069632","email":null,"discriminator":"9338","bot":true,"avatar":null},"session_id":"d3954ff063fa8d387ec395fe65723624","relationships":[],"private_channels":[],"presences":[],"guilds":[{"unavailable":true,"id":"486833041486905345"},{"unavailable":true,"id":"486833611564253184"}],"_trace":["gateway-prd-main-kg6w","discord-sessions-prd-1-27"]}}`)
//			mocker.reading <- evt
//			*seq++
//		}
//
//		wg.Wait()
//	}
//
//	for i := 0; i < b.N; i++ {
//		f(&mocker, wg, &seq)
//	}
//
//	close(d.shutdownChan)
//	close(wsShutdownChan)
//}

func TestCtrl(t *testing.T) {
	var ctrl *Ctrl
	newCtrl := func(c *Ctrl) *Ctrl {
		c.OnInsert(nil)
		return c
	}

	ctrl = newCtrl(&Ctrl{})
	if ctrl.IsDead() {
		t.Error("ctrl is marked dead even though no conditions for it's death was defined")
	}

	ctrl = newCtrl(&Ctrl{Runs: -1})
	if ctrl.IsDead() {
		t.Error("ctrl is marked dead even though no conditions for it's death was defined")
	}

	t.Run("counter", func(t *testing.T) {
		ctrl = newCtrl(&Ctrl{Runs: 5})
		if ctrl.IsDead() {
			t.Error("ctrl is marked dead too early")
		}
		for i := 0; i < 4; i++ {
			ctrl.Update()
			if ctrl.IsDead() {
				t.Errorf("ctrl is marked dead too early. Counter is %d", ctrl.Runs)
			}
		}
		ctrl.Update()
		if !ctrl.IsDead() {
			t.Errorf("ctrl was not marked dead. But was expected to be. Counter is %d", ctrl.Runs)
		}
	})

	t.Run("until", func(t *testing.T) {
		ctrl = newCtrl(&Ctrl{Until: time.Now().Add(1 * time.Millisecond)})
		if ctrl.IsDead() {
			t.Error("ctrl is marked dead too early")
		}
		<-time.After(3 * time.Millisecond)
		if !ctrl.IsDead() {
			t.Error("ctrl is dead. But was not marked dead even thought the condition have been met")
		}
	})

	t.Run("duration", func(t *testing.T) {
		ctrl = newCtrl(&Ctrl{Duration: 1 * time.Millisecond})
		if ctrl.IsDead() {
			t.Error("ctrl is marked dead too early")
		}
		<-time.After(3 * time.Millisecond)
		if !ctrl.IsDead() {
			t.Error("ctrl is dead. But was not marked dead even thought the condition have been met")
		}
	})

}

func check(err error, t *testing.T) {
	// Hide function from stacktrace, PR#3
	t.Helper()

	if err != nil {
		t.Error(err)
	}
}

func TestError_InterfaceImplementations(t *testing.T) {
	var u interface{} = &ErrorUnsupportedType{}

	t.Run("error", func(t *testing.T) {
		if _, ok := u.(error); !ok {
			t.Error("ErrorUnsupportedType does not implement error")
		}
	})
}

func TestTime(t *testing.T) {
	t.Run("omitempty", func(t *testing.T) {
		b := struct {
			T Time `json:"time,omitempty"`
		}{}

		bBytes, err := json.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}

		if string(bBytes) != `{"time":""}` {
			t.Errorf("did not get an 'omitted' field. Got %s", string(bBytes))
		}
	})
}

func TestDiscriminator(t *testing.T) {
	t.Run("String()", func(t *testing.T) {
		var d Discriminator
		d = Discriminator(0)
		if d.String() != "" {
			t.Errorf("got %s, wants \"\"", d.String())
		}

		d = Discriminator(1)
		if d.String() != "0001" {
			t.Errorf("got %s, wants \"0001\"", d.String())
		}

		d = Discriminator(4)
		if d.String() != "0004" {
			t.Errorf("got %s, wants \"0004\"", d.String())
		}

		d = Discriminator(12)
		if d.String() != "0012" {
			t.Errorf("got %s, wants \"0012\"", d.String())
		}

		d = Discriminator(120)
		if d.String() != "0120" {
			t.Errorf("got %s, wants \"0120\"", d.String())
		}

		d = Discriminator(1201)
		if d.String() != "1201" {
			t.Errorf("got %s, wants \"1201\"", d.String())
		}
	})
	t.Run("UnmarshalJSON(...)", func(t *testing.T) {
		var data []byte
		var d Discriminator
		var err error

		data = []byte{
			'"', '0', '0', '0', '1', '"',
		}
		if err = json.Unmarshal(data, &d); err != nil {
			t.Error(err)
		}
		executeInternalUpdater(d)

		if d.String() != "0001" {
			t.Errorf("got %s, wants \"0001\"", d.String())
		}

		data = []byte{
			'"', '0', '2', '0', '1', '"',
		}
		err = json.Unmarshal(data, &d)
		if err != nil {
			t.Error(err)
		}
		if d.String() != "0201" {
			t.Errorf("got %s, wants \"0201\"", d.String())
		}

		data = []byte{
			'"', '"',
		}
		if err = json.Unmarshal(data, &d); err != nil {
			t.Error(err)
		}
		executeInternalUpdater(d)

		if d.String() != "" {
			t.Errorf("got %s, wants \"\"", d.String())
		}
		if !d.NotSet() {
			t.Error("expected Discriminator to be NotSet")
		}

	})
	t.Run("MarshalJSON()", func(t *testing.T) {
		d := Discriminator(34)
		data, err := json.Marshal(&d)
		if err != nil {
			t.Error(err)
		}
		if string(data) != string([]byte("\"0034\"")) {
			t.Errorf("wrong. got %s, expects \"0034\"", string(data))
		}

		d = Discriminator(0)
		data, err = json.Marshal(&d)
		if err != nil {
			t.Error(err)
		}
		if string(data) != string([]byte("\"\"")) {
			t.Errorf("wrong. got %s, expects \"\"", string(data))
		}
	})
	t.Run("NotSet()", func(t *testing.T) {
		d := Discriminator(34)
		if d.NotSet() {
			t.Error("expected Discriminator.NotSet to be false, got true")
		}

		d = Discriminator(0)
		if !d.NotSet() {
			t.Error("expected Discriminator.NotSet to be true, got false")
		}
	})
}

func BenchmarkDiscriminator(b *testing.B) {
	b.Run("comparison", func(b *testing.B) {
		b.Run("string", func(b *testing.B) {
			val := "0401"
			vals := []string{
				"0401", "0001", "0400", "0011", "0101", "5435", "0010",
			}
			var i int
			var match bool
			for n := 0; n < b.N; n++ {
				for i = range vals {
					match = val == vals[i]
				}
			}
			if match {
			}
		})
		b.Run("uint16", func(b *testing.B) {
			val := Discriminator(0401)
			vals := []Discriminator{
				0401, 0001, 0400, 0011, 0101, 5435, 0010,
			}
			var i int
			var match bool
			for n := 0; n < b.N; n++ {
				for i = range vals {
					match = val == vals[i]
				}
			}
			if match {
			}
		})
	})
	b.Run("Unmarshal", func(b *testing.B) {
		dataSets := [][]byte{
			[]byte("\"0001\""),
			[]byte("\"0452\""),
			[]byte("\"4342\""),
			[]byte("\"0100\""),
			[]byte("\"5100\""),
			[]byte("\"1000\""),
			[]byte("\"5129\""),
			[]byte("\"3020\""),
			[]byte("\"0010\""),
		}
		b.Run("string", func(b *testing.B) {
			var result string
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				result = string(dataSets[i])
				if i == length {
					i = 0
				}
			}
			if result == "" {
			}
		})
		b.Run("uint16-a", func(b *testing.B) {
			var result uint16
			var i int
			lengthi := len(dataSets)
			for n := 0; n < b.N; n++ {
				data := dataSets[i]
				result = 0
				length := len(data) - 1
				for j := 1; j < length; j++ {
					result = result*10 + uint16(data[j]-'0')
				}
				if i == lengthi {
					i = 0
				}
			}
			if result == 0 {
			}
		})
		b.Run("uint16-b", func(b *testing.B) {
			var result uint16
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				data := dataSets[i]
				var tmp uint64
				tmp, _ = strconv.ParseUint(string(data), 10, 16)
				result = uint16(tmp)
				if i == length {
					i = 0
				}
			}
			if result == 0 {
			}
		})
		type fooOld struct {
			Foo string `json:"discriminator,omitempty"`
		}
		type fooNew struct {
			Foo Discriminator `json:"discriminator,omitempty"`
		}
		b.Run("string-struct", func(b *testing.B) {
			foo := &fooOld{}
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				_ = json.Unmarshal(dataSets[i], foo)
				executeInternalUpdater(foo)
				if i == length {
					i = 0
				}
			}
		})
		b.Run("uint16-struct", func(b *testing.B) {
			foo := &fooNew{}
			var i int
			length := len(dataSets)
			for n := 0; n < b.N; n++ {
				_ = json.Unmarshal(dataSets[i], foo)
				executeInternalUpdater(foo)
				if i == length {
					i = 0
				}
			}
		})
	})
}
