package discordws

import (
	"encoding/json"
	"errors"
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
func (c *Client) Connect() (err error) {
	c.Lock()
	defer c.Unlock()

	// check if web socket connection is already open
	if !c.Dead() {
		return errors.New("websocket connection already established, cannot open a new one")
	}

	// discord API sends a web socket url which should be used.
	// It's required to be cached, and only re-requested whenever disgord is unable
	// to reconnect to the API..
	if !c.Routed() {
		var resp *http.Response
		resp, err = c.httpClient.Get(c.urlAPIVersion + "/gateway")
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

		c.url = gatewayResponse.URL + "?v=" + strconv.Itoa(c.dAPIVersion) + "&encoding=" + c.dAPIEncoding
	}

	defer func() {
		if err != nil {
			c.conn.Close()
			c.conn = nil
			close(c.disconnected)
			logrus.Error(err)
		}
	}()

	// establish ws connection
	logrus.Info("Connecting to url: " + c.url)
	c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return
	}
	logrus.Info("established web socket connection")
	c.disconnected = make(chan struct{})

	// setup operation listeners
	// These handle sepecific "events" related to the socket connection
	go c.operationHandlers()

	// setup read and write goroutines
	go c.readPump()
	go c.writePump()

	return
}

func (c *Client) operationHandlers() {
	logrus.Info("Ready to recieve operation codes...")
	for {
		select {
		case gp, ok := <-c.operationChan:
			if !ok {
				logrus.Info("operationChan is dead..")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			switch gp.Op {
			case 0:
				// normal package
				logrus.Info("Got a normal package")
			case 1:
				// ping
				c.Lock()
				snr := c.sequenceNumber
				c.Unlock()

				c.sendChan <- struct {
					OP uint  `json:"op"`
					d  *uint `json:"d"`
				}{1, &snr}

			// case 2:
			// case 3:
			// case 4:
			// case 5:
			// case 6:
			case 7:
				// reconnect
			// case 8:
			case 9:
				// invalid session. Must respond with a identify packet
				err := c.sendIdentity()
				if err != nil {
					logrus.Error(err)
				}
			case 10:
				// hello
				b, err := json.Marshal(gp)
				if err != nil {
					logrus.Error(err)
				}

				gp2 := &GatewayPayload{
					Data: c,
				}
				err = json.Unmarshal(b, gp2)
				if err != nil {
					logrus.Error(err)
				}

				// TODO, this might create several idle goroutines..
				go c.pulsate(c.conn, c.disconnected)

				// send identify or resume packet
				if c.SessionID == "" && c.sequenceNumber == 0 {
					err = c.sendIdentity()
					if err != nil {
						logrus.Error(err)
					}
				} else {
					resume := GatewayPayload{
						Op: 6,
						Data: struct {
							Token      string `json:"token"`
							SessionID  string `json:"session_id"`
							SequenceNr uint   `json:"seq"`
						}{c.token, c.SessionID, c.sequenceNumber},
					}

					c.sendChan <- resume
				}
			case 11:
				// zombied or failed connection
			default:
				// unknown
				logrus.Warningf("Unknown operation: %+v\n", gp)
			}

		case <-c.disconnected:
			logrus.Info("exiting operation handler")
			return
		}
	}
}

//
// 	// NOTE. I completely stole this defer func from discordgo. It was a too nice not to take advantage of.
// 	//  Will this return any potential errors?
// 	defer func() {
// 		if err != nil {
// 			dws.Close()
// 			dws.Conn = nil
// 		}
// 	}()
//
// 	// Once connected, the client should immediately receive an Opcode 10 Hello payload, with information on the
// 	// connection's heartbeat interval:
// 	// {
// 	//   "heartbeat_interval": 45000,
// 	//   "_trace": ["discord-gateway-prd-1-99"]
// 	// }
// 	messageType, packet, err := dws.ReadMessage()
// 	fmt.Printf("%+v\n", string(packet))
// 	if err != nil {
// 		return
// 	}
// 	if messageType != 1 {
// 		logrus.Fatal("encrypted hello package, missing support")
// 	}
//
// 	// handle the incoming hello data and store it in our ws configuration
// 	event := &GatewayPayload{
// 		Data: dws,
// 	}
// 	err = json.Unmarshal(packet, &event)
// 	if err != nil {
// 		return
// 	}
// 	if event.Op != 10 {
// 		err = errors.New("while opening a discord ws connection, discord responded with op " + string(event.Op) + " when 10 was expected")
// 		return
// 	}
//
// 	fmt.Println(dws.HeartbeatInterval)
//
// 	// create a new session, otherwise respond with a resume packet
// 	if dws.SessionID == "" && dws.sequenceNumber == 0 {
// 		err = dws.sendIdentity()
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		resume := GatewayPayload{
// 			Op: 6,
// 			Data: struct {
// 				Token      string `json:"token"`
// 				SessionID  string `json:"session_id"`
// 				SequenceNr uint   `json:"seq"`
// 			}{dws.token, dws.SessionID, dws.sequenceNumber},
// 		}
//
// 		dws.wsMutex.Lock()
// 		err = dws.WriteJSON(resume)
// 		dws.wsMutex.Unlock()
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	// Expecting a resumed, invalid session or ready event
// 	messageType, packet, err = dws.ReadMessage()
// 	if err != nil {
// 		return
// 	}
// 	event, err = dws.onEvent(messageType, packet)
// 	if err != nil {
// 		return
// 	}
//
// 	// first incoming data recieved. We are now connected to discord.
// 	dws.connected = make(chan struct{})
//
// 	// start pulsating keep-alive packets
// 	go dws.pulsate(dws.Conn, dws.connected)
// 	go dws.listenForDiscordEvents(dws.Conn, dws.connected)
//
// 	return nil
// }

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		err = errors.New("No websocket connection exist")
		return
	}

	defer c.conn.Close()
	done := make(chan struct{})

	go func() {
		defer c.conn.Close()
		defer close(done)

		if c.disconnected != nil {
			close(c.disconnected)
		}
		for {
			_, message, err1 := c.conn.ReadMessage()
			if err1 != nil {
				log.Println("read:", err1)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	c.wsMutex.Lock()
	err = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.wsMutex.Unlock()
	if err != nil {
		logrus.Warningln(err)
	}
	select {
	case <-done:
	case <-time.After(time.Second * 2):
	}
	err = c.conn.Close()
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

func (c *Client) pulsate(ws *websocket.Conn, disconnected <-chan struct{}) {
	//c.pulseMutex.Lock()
	//defer c.pulseMutex.Unlock()
	// previous := time.Now().UTC()

	ticker := time.NewTicker(time.Millisecond * c.HeartbeatInterval)
	defer ticker.Stop()

	for {

		c.Lock()
		snr := c.sequenceNumber
		c.Unlock()

		data := struct {
			OP uint `json:"op"`
			d  uint `json:"d"`
		}{1, snr}

		c.sendChan <- data
		select {
		case <-ticker.C:
			continue
		case <-disconnected:
			return
		}
	}
}

func (c *Client) sendIdentity() (err error) {
	// https://discordapp.com/developers/docs/topics/gateway#identify
	identityPayload := struct {
		Token          string      `json:"token"`
		Properties     interface{} `json:"properties"`
		Compress       bool        `json:"compress"`
		LargeThreshold uint        `json:"large_threshold"`
		Shard          *[2]uint    `json:"shard,omitempty"`
		Presence       interface{} `json:"presence,omitempty"`
	}{
		Token: c.token,
		Properties: struct {
			OS      string `json:"$os"`
			Browser string `json:"$browser"`
			Device  string `json:"$device"`
		}{runtime.GOOS, c.String(), c.String()},
		LargeThreshold: 250,
		// Presence: struct {
		// 	Since  *uint       `json:"since"`
		// 	Game   interface{} `json:"game"`
		// 	Status string      `json:"status"`
		// 	AFK    bool        `json:"afk"`
		// }{Status: "online"},
	}

	if c.ShardCount > 1 {
		identityPayload.Shard = &[2]uint{uint(c.ShardID), c.ShardCount}
	}

	identity := GatewayPayload{
		Op:   2,
		Data: identityPayload,
	}

	c.sendChan <- identity
	return nil
}
