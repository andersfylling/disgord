package discordws

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	ReadyKey   string = "READY"
	ResumedKey string = "RESUMED"
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

	defer func(err error) error {
		if err != nil {
			if c.conn != nil {
				c.conn.Close()
				c.conn = nil
				close(c.disconnected)
			}
			logrus.Error(err)
			return err
		}
		return nil
	}(err)

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

			switch gp.Op {
			case 0:
				// discord events
				// events that directly correlates to the socket layer, will be dealt with here. But still dispatched.

				// increment the sequence number for each event to make sure everything is synced with discord
				c.Lock()
				c.sequenceNumber++ // = gp.SequenceNumber
				c.Unlock()

				// always store the session id
				if gp.EventName == ReadyKey {
					ready := &readyPacket{}
					err := json.Unmarshal(gp.Data.ByteArr(), ready)
					if err != nil {
						logrus.Error(err)
					}

					c.RLock()
					c.SessionID = ready.SessionID
					c.Trace = ready.Trace
					c.RUnlock()
				} else if gp.EventName == ResumedKey {
					// eh? debugging.
				}

				// dispatch events
				eventPkt := &Event{gp}
				c.iEventChan <- eventPkt
			case 1:
				// ping
				c.RLock()
				snr := c.sequenceNumber
				c.RUnlock()

				c.sendChan <- &gatewayPayload{Op: 1, Data: &snr}
			case 7:
				// reconnect
				c.Disconnect()
				go c.reconnect()
			case 9:
				time.Sleep(time.Second * time.Duration(rand.Intn(4)+1)) // [1,5]
				// invalid session. Must respond with a identify packet
				err := c.sendIdentity()
				if err != nil {
					logrus.Error(err)
				}
			case 10:
				// hello
				helloPk := &helloPacket{}
				err := json.Unmarshal(gp.Data.ByteArr(), helloPk)
				if err != nil {
					logrus.Debug(err)
				}
				c.Lock()
				c.HeartbeatInterval = helloPk.HeartbeatInterval
				c.Unlock()

				// TODO, this might create several idle goroutines..
				go c.pulsate(c.conn, c.disconnected)

				// send identify or resume packet
				if c.SessionID == "" && c.sequenceNumber == 0 {
					err = c.sendIdentity()
					if err != nil {
						logrus.Error(err)
					}
				} else {
					c.RLock()
					token := c.token
					session := c.SessionID
					sequence := c.sequenceNumber
					c.RUnlock()

					resume := &gatewayPayload{
						Op: 6,
						Data: struct {
							Token      string `json:"token"`
							SessionID  string `json:"session_id"`
							SequenceNr *uint  `json:"seq"`
						}{token, session, &sequence},
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
	c.RLock()
	if c.conn == nil {
		err = errors.New("No websocket connection exist")
		return
	}
	c.RUnlock()

	c.wsMutex.Lock()
	err = c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.wsMutex.Unlock()
	if err != nil {
		logrus.Warningln(err)
	}
	select {
	case <-c.disconnected:
		// might get closed by another process
	case <-time.After(time.Second * 1):
		close(c.disconnected)
	}

	// give remainding processes some time to exit
	<-time.After(time.Second * 1)
	c.disconnected = nil

	// close connection
	err = c.conn.Close()
	c.conn = nil
	return
}

// Reconnect to discord endpoint
func (c *Client) reconnect() (err error) {
	for try := 0; try <= MaxReconnectTries; try++ {

		logrus.Debugf("Reconnect attempt #%d\n", try)
		err = c.Connect()
		if err == nil {
			logrus.Info("successfully reconnected")

			// send resume package

			break
			// TODO voice
		}
		if try == MaxReconnectTries-1 {
			err = errors.New("Too many reconnect attempts")
			return err
		}

		// wait 5 seconds
		logrus.Info("reconnect failed, trying again in 5 seconds")
		<-time.After(time.Duration(5) * time.Second)
	}

	return
}

func (c *Client) pulsate(ws *websocket.Conn, disconnected chan struct{}) {
	c.pulseMutex.Lock()
	defer c.pulseMutex.Unlock()
	// previous := time.Now().UTC()

	ticker := time.NewTicker(time.Millisecond * time.Duration(c.HeartbeatInterval))
	defer ticker.Stop()

	for {

		c.Lock()
		last := c.heartbeatAcquired
		interval := c.HeartbeatInterval
		snr := c.sequenceNumber
		c.Unlock()

		c.sendChan <- &gatewayPayload{Op: 1, Data: snr}

		// verify the heartbeat ACK
		go func(client *Client, last time.Time) {
			time.Sleep((3 * time.Second) % (time.Duration(interval) * time.Millisecond))
			var die bool
			c.Lock()
			die = c.heartbeatAcquired == last
			c.Unlock()

			if die {
				logrus.Debug("heartbeat ACK was not recieved")
				c.Disconnect()
				go c.reconnect()
			}
		}(c, last)

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
