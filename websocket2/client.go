package websocket2

//
//import (
//	"time"
//
//	"github.com/andersfylling/disgord/websocket/event"
//	"github.com/pkg/errors"
//
//	"github.com/andersfylling/disgord/constant"
//)
//
//type wsPermission chan interface{}
//
//// client works as a base client for events, voice, etc.
//// it's designed by choreographic programming to properly synchronise state
//// and actions that are state dependent. eg. Connect should not connect if
//// the current state is anything but disconnected.
//type client struct {
//	conState
//
//	reqPermChan chan wsPermission
//}
//
//type step struct {
//	name      string // of step
//	condition string // state
//	cmd       string // send-id, wc-con, ping, send-resume
//	onStart   string // state while running step
//	onSuccess string // new state
//	onFailure string // new state
//}
//
//type Process struct {
//	name  string
//	steps []step
//}
//
//var connectProcess Process = Process{
//	name: "connect",
//	steps: []step{
//		{
//			name:      "connect to gateway",
//			condition: "disconnected",
//			onStart:   "connecting",
//			onSuccess: "connected",
//		},
//		{
//			name:      "wait for Hello",
//			condition: "connected",
//			cmd:       "wait-for-op-10",
//			onFailure: "disconnected",
//		},
//		{
//			name:      "start ping routine",
//			condition: "connected",
//			cmd:       "start-heartbeat",
//			onFailure: "disconnected",
//		},
//		{
//			name:      "send identify",
//			condition: "connected",
//			cmd:       "send-identify",
//			onFailure: "disconnected",
//		},
//	},
//}
//
//func connect() {
//	var state int
//	for {
//		select {
//		case <-event:
//		}
//
//		switch event.Op {
//		case <-stop: // shard is stopping/shutting down
//			return
//		case <-disconnected: // https://discordapp.com/developers/docs/topics/gateway#disconnections
//			return
//		case 11:
//			// hello
//		case 0:
//			// check if ready
//		case 9:
//			// invalid session
//			// check if we sent a resume or identify
//
//			// ETC.
//		}
//	}
//}
//
//type heartbeatConfig struct {
//	log      constant.Logger
//	stop     <-chan interface{}
//	timeout  uint
//	helloAck <-chan interface{}
//}
//
//func heartbeat(conf *heartbeatConfig) {
//	// read helloAck or incoming evt after it's parsed
//	// update
//	// way to send hello pkt
//
//	for range discordRequestedHeartbeat {
//		// let's remove old requests if there are any
//	}
//
//	sequence := 0
//	for {
//		select {
//		case <-conf.stop:
//			return
//		case <-time.After(time.Duration(conf.timeout) * time.Second):
//		case <-discordRequestedHearthbeat:
//		}
//		sendHeartbeat(sequence)
//
//		start := time.Now()
//		select {
//		case <-conf.stop:
//			return
//		case <-time.After(5 * time.Second):
//			// change state to "disconnected"
//			return
//		case <-conf.helloAck:
//			// store latency somehow
//			delay := time.Now().Sub(start)
//
//		}
//	}
//}
//
//func sendHearthbeat() {}
//
//func sendHello() {
//	canResume := c.sessionID != "" && c.sequenceNumber > 0
//	if !canResume {
//		err := sendIdentityPacket(c)
//		if err != nil {
//			c.Error(err.Error())
//		}
//		return
//	}
//
//	c.RLock()
//	token := c.conf.BotToken
//	session := c.sessionID
//	sequence := c.sequenceNumber
//	c.RUnlock()
//
//	err := c.Emit(event.Resume, struct {
//		Token      string `json:"token"`
//		SessionID  string `json:"session_id"`
//		SequenceNr uint   `json:"seq"`
//	}{token, session, sequence})
//	if err != nil {
//		c.Error(err.Error())
//	}
//
//	err = releaseConnectPermission(c)
//	if err != nil {
//		err = errors.New("unable to release connection permission. Err: " + err.Error())
//		c.Error(err.Error())
//	}
//}
