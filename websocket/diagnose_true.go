// +build disgord_diagnosews

package websocket

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/andersfylling/disgord/httd"
)

const SaveIncomingPackets = true

var outgoingPacketSequence uint64 = 0

// saveOutgoingPacket saves raw json content to disk
// format: I_<seq>_<op>_<shard_id>_<unix>.json
// unix is the unix timestamp on save
// seq is the sequence number: outgoingPacketSequence
// op is the operation code
func saveOutgoingPacket(c *baseClient, packet *clientPacket) {
	data, err := httd.Marshal(packet)
	if err != nil {
		fmt.Println(err)
	}

	unix := strconv.FormatInt(time.Now().UnixNano(), 10)

	seq := strconv.FormatUint(uint64(outgoingPacketSequence), 10)
	outgoingPacketSequence++

	shardID := strconv.FormatUint(uint64(c.ShardID), 10)

	op := strconv.FormatUint(uint64(packet.Op), 10)

	filename := "O_" + seq + "_" + op + "_" + shardID + "_" + unix + ".json"
	err = ioutil.WriteFile(DiagnosePath_packets+"/"+filename, data, 0644)
	if err != nil {
		c.Error(err.Error())
	}
}

// saveIncomingPacker saves raw json content to disk
// format: O_<unix>_<seq>_<op>_<shard_id>[_<evt_name>].json
// unix is the unix timestamp on save. This is needed as the sequence number can be reset.
// seq is the sequence number
// op is the operation code
// evt_name is the event name (optional)
func saveIncomingPacker(c *baseClient, evt *discordPacket, packet []byte) {
	evtStr := "_" + evt.EventName
	if evtStr == "_" {
		evtStr = ""
	}
	unix := strconv.FormatInt(time.Now().UnixNano(), 10)
	seq := strconv.FormatUint(uint64(evt.SequenceNumber), 10)
	op := strconv.FormatUint(uint64(evt.Op), 10)
	shardID := strconv.FormatUint(uint64(c.ShardID), 10)

	filename := "I_" + unix + "_" + seq + "_" + op + "_" + shardID + evtStr + ".json"
	err := ioutil.WriteFile(DiagnosePath_packets+"/"+filename, packet, 0644)
	if err != nil {
		c.Error(err.Error())
	}
}
