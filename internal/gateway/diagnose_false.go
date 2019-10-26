// +build !disgord_diagnosews

package gateway

const SaveIncomingPackets = false

func saveOutgoingPacket(c *client, packet *clientPacket)                {}
func saveIncomingPacker(c *client, event *DiscordPacket, packet []byte) {}
