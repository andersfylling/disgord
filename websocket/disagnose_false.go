// +build !disgord_diagnosews

package websocket

const SaveIncomingPackets = false

func saveOutgoingPacket(c *baseClient, packet *clientPacket)                {}
func saveIncomingPacker(c *baseClient, event *discordPacket, packet []byte) {}
