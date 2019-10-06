package websocket

import (
	cmd2 "github.com/andersfylling/disgord/internal/websocket/cmd"
	event2 "github.com/andersfylling/disgord/internal/websocket/event"
	opcode2 "github.com/andersfylling/disgord/internal/websocket/opcode"
)

func CmdNameToOpCode(command string, t ClientType) (op opcode2.OpCode) {
	op = opcode2.None
	// TODO: refactor command and event name to avoid conversion (?)
	if t == clientTypeVoice {
		switch command {
		case cmd2.VoiceSpeaking:
			op = opcode2.VoiceSpeaking
		case cmd2.VoiceIdentify:
			op = opcode2.VoiceIdentify
		case cmd2.VoiceSelectProtocol:
			op = opcode2.VoiceSelectProtocol
		case cmd2.VoiceHeartbeat:
			op = opcode2.VoiceHeartbeat
		case cmd2.VoiceResume:
			op = opcode2.VoiceResume
		}
	} else if t == clientTypeEvent {
		switch command {
		case event2.Heartbeat:
			op = opcode2.EventHeartbeat
		case event2.Identify:
			op = opcode2.EventIdentify
		case event2.Resume:
			op = opcode2.EventResume
		case cmd2.RequestGuildMembers:
			op = opcode2.EventRequestGuildMembers
		case cmd2.UpdateVoiceState:
			op = opcode2.EventVoiceStateUpdate
		case cmd2.UpdateStatus:
			op = opcode2.EventStatusUpdate
		}
	}

	return op
}
