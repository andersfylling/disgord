package disgord

import "testing"

func TestConnect(t *testing.T) {
	d, err := NewDisgord()
	if err != nil {
		t.Error(err.Error())
	}
	d.Connect()
	d.Disconnect()
}
