package gateway

import (
	"time"

	"github.com/andersfylling/disgord/internal/gateway/cmd"
	"github.com/andersfylling/disgord/internal/gateway/event"
	"github.com/andersfylling/disgord/internal/gateway/opcode"
)

func CmdNameToOpCode(command string, t ClientType) (op opcode.OpCode) {
	op = opcode.None
	// TODO: refactor command and event name to avoid conversion (?)
	if t == clientTypeVoice {
		switch command {
		case cmd.VoiceSpeaking:
			op = opcode.VoiceSpeaking
		case cmd.VoiceIdentify:
			op = opcode.VoiceIdentify
		case cmd.VoiceSelectProtocol:
			op = opcode.VoiceSelectProtocol
		case cmd.VoiceHeartbeat:
			op = opcode.VoiceHeartbeat
		case cmd.VoiceResume:
			op = opcode.VoiceResume
		}
	} else if t == clientTypeEvent {
		switch command {
		case event.Heartbeat:
			op = opcode.EventHeartbeat
		case event.Identify:
			op = opcode.EventIdentify
		case event.Resume:
			op = opcode.EventResume
		case cmd.RequestGuildMembers:
			op = opcode.EventRequestGuildMembers
		case cmd.UpdateVoiceState:
			op = opcode.EventVoiceStateUpdate
		case cmd.UpdateStatus:
			op = opcode.EventStatusUpdate
		}
	}

	return op
}

// rotateByTime every timestamp after "limit" is deleted. Assumes the oldest entries are first.
func rotateByTime(times []time.Time, limit time.Time) []time.Time {
	var delim int
	for i := range times {
		if times[i].After(limit) {
			delim++
		} else {
			break
		}
	}

	// shift
	for i := 0; i+delim < len(times); i++ {
		times[i] = times[delim+i]
	}

	return times[:len(times)-delim]
}
