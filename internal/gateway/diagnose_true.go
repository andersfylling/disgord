// +build disgord_diagnosews

package gateway

import (
	"bytes"
	"fmt"
	"github.com/andersfylling/disgord/json"
	"io/ioutil"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
)

const SaveIncomingPackets = true

const DiagnosePath = "diagnose-report"
const DiagnosePath_packets = "diagnose-report/packets"

var outgoingPacketSequence uint64
var dirExists bool

func formatFilename(incoming bool, clientType ClientType, shardID uint, opCode opcode.OpCode, sequencenr uint64, suffix string) (filename string) {

	unix := strconv.FormatInt(time.Now().UnixNano(), 10)
	shard := strconv.FormatUint(uint64(shardID), 10)
	op := strconv.FormatUint(uint64(opCode), 10)
	seq := strconv.FormatUint(uint64(sequencenr), 10)

	var direction string
	if incoming {
		direction = "IN"
	} else {
		direction = "OUT"
	}

	var t string
	if clientType == clientTypeVoice {
		t = "V"
	} else if clientType == clientTypeEvent {
		t = "E"
	} else {
		t = "-"
	}

	return unix + "_" + t + "_" + direction + "_id" + shard + "_op" + op + "_s" + seq + suffix + ".json"
}

func ensureDir() {
	if dirExists {
		return
	}

	if _, err := os.Stat(DiagnosePath_packets); err == nil {
		return
	}

	if err := os.MkdirAll(DiagnosePath_packets, os.ModePerm); err != nil {
		fmt.Println("unable to create directory " + DiagnosePath_packets + ", please create it and restart")
		panic(err)
	}

	dirExists = true
}

// saveOutgoingPacket saves raw json content to disk
// format: I_<seq>_<op>_<shard_id>_<unix>.json
// unix is the unix timestamp on save
// seq is the sequence number: outgoingPacketSequence
// op is the operation code
func saveOutgoingPacket(c *client, packet *clientPacket) {
	ensureDir()
	data, err := json.MarshalIndent(packet, "", "\t")
	if err != nil {
		c.log.Debug(c.getLogPrefix(), err)
	}

	nr := atomic.AddUint64(&outgoingPacketSequence, 1)
	filename := formatFilename(false, c.clientType, c.ShardID, packet.Op, nr, "")

	path := DiagnosePath_packets + "/" + filename
	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		c.log.Debug(c.getLogPrefix(), err)
	}

	c.log.Info(c.getLogPrefix(), "saved "+filename)
}

// saveIncomingPacker saves raw json content to disk
// format: O_<unix>_<seq>_<op>_<shard_id>[_<evt_name>].json
// unix is the unix timestamp on save. This is needed as the sequence number can be reset.
// seq is the sequence number
// op is the operation code
// evt_name is the event name (optional)
func saveIncomingPacker(c *client, evt *DiscordPacket, packet []byte) {
	ensureDir()
	evtStr := "_" + evt.EventName
	if evtStr == "_" {
		evtStr = "_EMPTY"
	}

	filename := formatFilename(true, c.clientType, c.ShardID, evt.Op, uint64(evt.SequenceNumber), evtStr)
	path := DiagnosePath_packets + "/" + filename

	// pretty
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, packet, "", "\t"); err != nil {
		c.log.Debug(c.getLogPrefix(), err)
	}

	if err := ioutil.WriteFile(path, prettyJSON.Bytes(), 0644); err != nil {
		c.log.Debug(c.getLogPrefix(), err)
	}

	c.log.Info(c.getLogPrefix(), "saved "+filename)
}
