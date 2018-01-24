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

	"github.com/andersfylling/disgord/event"
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

		gatewayResponse := &getGatewayResponse{}
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
	logrus.Debug("Connecting to url: " + c.url)
	c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return
	}
	logrus.Debug("established web socket connection")
	c.disconnected = make(chan struct{})

	// setup operation listeners
	// These handle sepecific "events" related to the socket connection
	go c.operationHandlers()

	// setup read and write goroutines
	go c.readPump()
	go c.writePump(c.conn)

	return
}

func (c *Client) operationHandlers() {
	logrus.Debug("Ready to recieve operation codes...")
	for {
		select {
		case gp, ok := <-c.operationChan:
			if !ok {
				logrus.Debug("operationChan is dead..")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// always store the sequence number
			c.Lock()
			c.sequenceNumber = gp.SequenceNumber
			c.Unlock()

			switch gp.Op {
			case 0:
				// discord events
				// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

				// always store the session id if it exists
				if gp.EventName == event.Ready {
					c.Lock()
					err := json.Unmarshal(gp.Data.ByteArr(), c)
					c.Unlock()
					if err != nil {
						logrus.Error(err)
					}
				} else if gp.EventName == event.Resumed {
					// eh? debugging.
				}

			case 1:
				// ping
				c.Lock()
				snr := c.sequenceNumber
				c.Unlock()

				c.sendChan <- &gatewayPayload{Op: 1, Data: &snr}
			case 7:
				// reconnect
			case 9:
				// invalid session. Must respond with a identify packet
				err := c.sendIdentity()
				if err != nil {
					logrus.Error(err)
				}
			case 10:
				// hello
				c.Lock()
				err := json.Unmarshal(gp.Data.ByteArr(), c)
				c.Unlock()
				if err != nil {
					logrus.Debug(err)
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
					resume := &gatewayPayload{
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
				// heartbeat recieved
				c.Lock()
				c.heartbeatAcquired = time.Now()
				c.Unlock()
			default:
				// unknown
				logrus.Debugf("Unknown operation: %+v\n", gp)
			}

		case <-c.disconnected:
			logrus.Debug("exiting operation handler")
			return
		}
	}
}

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
func (c *Client) Reconnect() error {
	c.Lock()
	defer c.Unlock()

	for try := 0; try < MaxReconnectTries; try++ {

	}

	return nil
}

func (c *Client) pulsate(ws *websocket.Conn, disconnected chan struct{}) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()
	// previous := time.Now().UTC()

	ticker := time.NewTicker(time.Millisecond * time.Duration(c.HeartbeatInterval))
	defer ticker.Stop()

	<-ticker.C

	for {

		c.Lock()
		last := c.heartbeatAcquired
		interval := c.HeartbeatInterval
		snr := c.sequenceNumber
		c.Unlock()

		if interval == 0 {
			logrus.Debug("heartbeat interval was 0")
			close(c.disconnected)
		}

		c.sendChan <- &gatewayPayload{Op: 1, Data: snr}

		// verify the heartbeat ACK
		go func(client *Client, last time.Time, disconnect chan struct{}) {
			time.Sleep((3 * time.Second) % (time.Duration(interval) * time.Millisecond))
			var die bool
			c.Lock()
			die = c.heartbeatAcquired == last
			c.Unlock()

			if die {
				logrus.Debug("heartbeat ACK was not recieved")
				close(disconnect)
			}
		}(c, last, disconnected)

		select {
		case <-ticker.C:
			continue
		case <-disconnected:
			logrus.Debug("Stopping pulse")
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

	identity := &gatewayPayload{
		Op:   2,
		Data: identityPayload,
	}

	c.sendChan <- identity
	return nil
}
