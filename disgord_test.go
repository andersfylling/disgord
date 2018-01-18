package disgord

import "testing"

const DoSocketTest = false

func TestConnect(t *testing.T) {
	if !DoSocketTest {
		return
	}
	d, err := NewDisgord()
	if err != nil {
		t.Error(err.Error())
	}
	d.Connect()
	d.Disconnect()
}
