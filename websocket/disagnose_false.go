// +build !disgord_diagnosews

package websocket

const SaveIncomingPackets = false

func saveOutgoingPacket(c *Client, packet *clientPacket)                {}
func saveIncomingPacker(c *Client, event *discordPacket, packet []byte) {}
