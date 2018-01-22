package discordws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Connect establishes a websocket connection to the discord API
func (dws *Client) Connect() (err error) {
	dws.Lock()
	defer dws.Unlock()

	// check if web socket connection is already open
	if !dws.Dead() {
		return errors.New("websocket connection already established, cannot open a new one")
	}

	// discord API sends a web socket url which should be used.
	// It's required to be cached, and only re-requested whenever disgord is unable
	// to reconnect to the API..
	if !dws.Routed() {
		var resp *http.Response
		resp, err = dws.HTTPClient.Get(dws.URLAPIVersion + "/gateway")
		if err != nil {
			return
		}
		defer resp.Body.Close()

		var body []byte
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		gatewayResponse := &GetGatewayResponse{}
		err = json.Unmarshal([]byte(body), gatewayResponse)
		if err != nil {
			return
		}

		dws.URL = gatewayResponse.URL + "?v=" + strconv.Itoa(dws.DiscordAPIVersion) + "&encoding=" + dws.DiscordAPIEncoding
	}

	// establish ws connection
	dws.Conn, _, err = websocket.DefaultDialer.Dial(dws.URL, nil)
	if err != nil {
		return
	}

	// NOTE. I completely stole this defer func from discordgo. It was a too nice not to take advantage of.
	//  Will this return any potential errors?
	defer func() {
		if err != nil {
			dws.Close()
			dws.Conn = nil
		}
	}()

	// Once connected, the client should immediately receive an Opcode 10 Hello payload, with information on the
	// connection's heartbeat interval:
	// {
	//   "heartbeat_interval": 45000,
	//   "_trace": ["discord-gateway-prd-1-99"]
	// }
	messageType, packet, err := dws.ReadMessage()
	fmt.Printf("%+v\n", string(packet))
	if err != nil {
		return
	}
	if messageType != 1 {
		logrus.Fatal("encrypted hello package, missing support")
	}

	// handle the incoming hello data and store it in our ws configuration
	event := &GatewayPayload{
		Data: dws,
	}
	err = json.Unmarshal(packet, &event)
	if err != nil {
		return
	}
	if event.Op != 10 {
		err = errors.New("while opening a discord ws connection, discord responded with op " + string(event.Op) + " when 10 was expected")
		return
	}

	fmt.Println(dws.HeartbeatInterval)

	// create a new session, otherwise respond with a resume packet
	if dws.SessionID == "" && dws.sequenceNumber == 0 {
		err = dws.sendIdentity()
		if err != nil {
			return err
		}
	} else {
		resume := GatewayPayload{
			Op: 6,
			Data: struct {
				Token      string `json:"token"`
				SessionID  string `json:"session_id"`
				SequenceNr uint   `json:"seq"`
			}{dws.token, dws.SessionID, dws.sequenceNumber},
		}

		dws.wsMutex.Lock()
		err = dws.WriteJSON(resume)
		dws.wsMutex.Unlock()
		if err != nil {
			return
		}
	}

	// Expecting a resumed, invalid session or ready event
	messageType, packet, err = dws.ReadMessage()
	if err != nil {
		return
	}
	event, err = dws.onEvent(messageType, packet)
	if err != nil {
		return
	}

	// first incoming data recieved. We are now connected to discord.
	dws.connected = make(chan struct{})

	// start pulsating keep-alive packets
	go dws.pulsate(dws.Conn, dws.connected)
	go dws.listenForDiscordEvents(dws.Conn, dws.connected)

	return nil
}

// Disconnect closes the discord websocket connection
func (dws *Client) Disconnect() (err error) {
	dws.Lock()
	defer dws.Unlock()

	if dws.Conn == nil {
		err = errors.New("No websocket connection exist")
		return
	}

	defer dws.Close()
	done := make(chan struct{})

	go func() {
		defer dws.Close()
		defer close(done)

		if dws.connected != nil {
			close(dws.connected)
			dws.connected = nil
		}
		for {
			_, message, err1 := dws.ReadMessage()
			if err1 != nil {
				log.Println("read:", err1)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	dws.wsMutex.Lock()
	err = dws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	dws.wsMutex.Unlock()
	if err != nil {
		logrus.Warningln(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second * 2):
	}
	err = dws.Close()
	dws.Conn = nil
	return
}

// Reconnect to discord endpoint
func (dws *Client) Reconnect() error {
	dws.Lock()
	defer dws.Unlock()

	for try := 0; try < MaxReconnectTries; try++ {

	}

	return nil
}

func (dws *Client) pulsate(ws *websocket.Conn, disconnected <-chan struct{}) {
	// previous := time.Now().UTC()

	ticker := time.NewTicker(time.Millisecond * dws.HeartbeatInterval)
	defer ticker.Stop()

	for {

		dws.RLock()
		snr := dws.sequenceNumber
		dws.RUnlock()

		dws.wsMutex.Lock()
		ws.WriteJSON(struct {
			OP uint  `json:"op"`
			d  *uint `json:"d"`
		}{1, &snr})
		dws.wsMutex.Unlock()

		select {
		case <-ticker.C:
			continue
		case <-disconnected:
			return
		}
	}
}

func (dws *Client) sendIdentity() (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       struct {
			Since  *uint       `json:"since"`
			Game   interface{} `json:"game"`
			Status string      `json:"status"`
			AFK    bool        `json:"afk"`
		} `json:"presence"`
	}{
		Token: dws.token,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, dws.String(), dws.String()},
		//Shard: dws
	}

	if dws.ShardCount > 1 {
		identityPayload.Shard = &[2]uint{uint(dws.ShardID), dws.ShardCount}
	}

	identity := GatewayPayload{
		Op:   2,
		Data: identityPayload,
	}

	fmt.Printf("Sending: %+v\n", identity)

	dws.wsMutex.Lock()
	err = dws.WriteJSON(identity)
	dws.wsMutex.Unlock()
	if err != nil {
		return
	}

	return nil
}
