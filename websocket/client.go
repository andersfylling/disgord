package websocket

import (
	"errors"
	"github.com/andersfylling/disgord/websocket/event"
	"github.com/andersfylling/disgord/websocket/opcode"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// Client is responsible for the external logic, and is decoupled from as much managing as possible. By managing
// I refer to the possibility to handle invalid connection, respond to events, reconnect, etc.
//
// This client only holds the required functionality to allow interacting with Discord, while the manager dictates the
// behaviour of the client. This decoupling allows for proper unit testing.
type Client interface {
	Connect() error
	Disconnect() error
	Emit(command string, data interface{}) error
	Receive() <-chan *discordPacket
}

// WebsocketErr is used internally when the websocket package returns an error. It does not represent a Discord error(!)
type WebsocketErr struct {
	ID      uint
	message string
}

func (e *WebsocketErr) Error() string {
	return e.message
}

// DefaultClientConfig is the configuration struct used for initializing a DefaultClient.
type DefaultClientConfig struct {
	// Token Discord bot token
	Token string

	// HTTPClient custom http client to support the use of proxy
	HTTPClient *http.Client

	// ChannelBuffer is used to set the event channel buffer
	ChannelBuffer uint

	// Endpoint for establishing socket connection. Either endpoints, `Gateway` or `Gateway Bot`, is used to retrieve
	// a valid socket endpoint from Discord
	Endpoint string

	// Encoding make sure we support the correct encoding
	Encoding string

	// Version make sure we support the correct Discord version
	Version int
}

func NewDefaultClient(config *DefaultClientConfig) (*DefaultClient, error) {
	return &DefaultClient{
		conf:        config,
		receiveChan: make(chan *discordPacket),
		emitChan:    make(chan *clientPacket),
		connection:  make(chan int8),
	}, nil
}

// DefaultClient is the default implementation for handling external communication with the Discord API. The client is
// only aware of connecting, disconnecting, sending and receiving data. It does not understand that there exist
// different Discord operation codes, nor that they each expect reaction or way to be handling. That is the role of the
// manager. See Manager.
type DefaultClient struct {
	sync.Mutex

	conf *DefaultClientConfig

	receiveChan chan *discordPacket
	emitChan    chan *clientPacket
	conn        *websocket.Conn
	connection  chan int8
}

// Connect establishes a socket connection with the Discord API
func (c *DefaultClient) Connect() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.conn != nil {
		err = errors.New("cannot connect while a connection already exist")
		return
	}

	if c.conf.Endpoint == "" {
		c.conf.Endpoint, err = getGatewayRoute(c.conf.HTTPClient, c.conf.Version)
		if err != nil {
			return
		}
	}

	// by default we use gorilla's websocket dialer here, but if the passed http client uses a custom transport
	// we make sure we open the websocket over the same transport/proxy, in case the user uses this
	dialer := websocket.DefaultDialer
	if t, ok := c.conf.HTTPClient.Transport.(*http.Transport); ok {
		dialer = &websocket.Dialer{
			HandshakeTimeout: dialer.HandshakeTimeout,
			Proxy:            t.Proxy,
			NetDialContext:   t.DialContext,
			NetDial:          t.Dial, // even though Dial is deprecated in http.Transport, it isn't in websocket
		}
	}

	// ready the error handler
	defer func(err error) error {
		if err != nil {
			if c.conn != nil {
				c.conn.Close()
				c.conn = nil
				if c.connection != nil {
					close(c.connection)
					c.connection = nil
				}
			}
			return err
		}
		return nil
	}(err)

	// prepare the receiver
	c.connection = make(chan int8)

	// establish ws connection
	c.conn, _, err = dialer.Dial(c.conf.Endpoint, nil)
	if err != nil {
		return
	}

	// we can now interact with Discord
	go c.receiver()
	go c.emitter()
	return
}

// Disconnect disconnects the socket connection
func (c *DefaultClient) Disconnect() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		err = errors.New("already disconnected")
		return
	}

	// use the emitter to dispatch the close message
	c.Emit(event.Shutdown, nil)

	// close connection
	<-time.After(time.Second * 1)
	close(c.connection)

	// wait for other processes to finish
	<-time.After(time.Second * 1)
	c.connection = nil
	err = c.conn.Close()
	c.conn = nil
	return
}

// Emit emits a command, if supported, and its data to the Discord Socket API
func (c *DefaultClient) Emit(command string, data interface{}) (err error) {
	var op uint
	switch command {
	case event.Shutdown:
		op = opcode.Shutdown
	case event.Heartbeat:
		op = opcode.Heartbeat
	case event.Identify:
		op = opcode.Identify
	case event.Resume:
		op = opcode.Resume
	case event.RequestGuildMembers:
		op = opcode.RequestGuildMembers
	case event.VoiceStateUpdate:
		op = opcode.VoiceStateUpdate
	case event.StatusUpdate:
		op = opcode.StatusUpdate
	default:
		err = errors.New("unsupported command: " + command)
		return
	}

	c.emitChan <- &clientPacket{
		Op:   op,
		Data: data,
	}
	return
}

// Receive returns the channel for receiving Discord packets
func (c *DefaultClient) Receive() <-chan *discordPacket {
	return c.receiveChan
}

// emitter holds the actually dispatching logic for the Emit method. See DefaultClient#Emit.
func (c *DefaultClient) emitter() {
	for {
		var msg *clientPacket
		var open bool

		select {
		case <-c.connection:
			// c.connection got closed
		case msg, open = <-c.emitChan:
		}
		if !open || (msg.Data == nil && msg.Op == opcode.Shutdown) {
			c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		}

		err := c.conn.WriteJSON(&msg)
		if err != nil {
			// TODO-logging
		}
	}
}

func (c *DefaultClient) receiver() {
	for {
		messageType, packet, err := c.conn.ReadMessage()
		if err != nil && (websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || c.connection == nil) {
			logrus.Debug("closing readPump")
			return
		}

		logrus.Debugf("<-: %+v\n", string(packet))

		if messageType == websocket.BinaryMessage {
			packet, err = decompressBytes(packet)
			if err != nil {
				logrus.Panic(err)
				continue
			}
		}

		// parse to gateway payload object
		evt := &discordPacket{}
		err = evt.UnmarshalJSON(packet)
		if err != nil {
			logrus.Error(err)
			continue
		}

		// notify listeners
		c.receiveChan <- evt
	}
}
