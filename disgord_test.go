package disgord

import "testing"

func TestConnect(t *testing.T) {
	d := NewDisgord()
	d.Connect()
	d.Disconnect()
}
