// +build !integration

package gateway

import (
	"errors"
	"testing"

	"github.com/andersfylling/disgord/internal/gateway/opcode"
)

func TestClientPktQueue_Add(t *testing.T) {
	limit := 5
	q := newClientPktQueue(limit)
	if !q.IsEmpty() {
		t.Error("should be empty")
	}

	if err := q.Add(&clientPacket{}); err != nil {
		t.Error(err)
	}

	if q.IsEmpty() {
		t.Error("should not be empty")
	}

	for i := 1; i < limit; i++ {
		_ = q.Add(&clientPacket{})
	}
	if len(q.messages) != q.limit {
		t.Fatal("number of entries in queue was less than the limit")
	}

	if err := q.Add(&clientPacket{}); err == nil {
		t.Error("expected an error when trying to add to a full queue")
	}

	q.Steal()
	_ = q.Add(&clientPacket{Op: opcode.EventStatusUpdate})
	if len(q.messages) != 1 {
		t.Fatal("number of entries in queue should be 1")
	}
	_ = q.Add(&clientPacket{Op: opcode.EventStatusUpdate})
	if len(q.messages) != 1 {
		t.Fatal("number of entries in queue should be 1")
	}
	_ = q.Add(&clientPacket{})
	if len(q.messages) != 2 {
		t.Fatal("number of entries in queue should be 2")
	}
}

func TestClientPktQueue_AddByOverwrite(t *testing.T) {
	q := newClientPktQueue(10)
	if !q.IsEmpty() {
		t.Error("should be empty")
	}

	if err := q.AddByOverwrite(&clientPacket{}); err == nil {
		t.Error("should complain that no similar entry to be overwritten was found")
	}

	if !q.IsEmpty() {
		t.Error("should be empty")
	}

	pkt := &clientPacket{}
	if err := q.Add(pkt); err != nil {
		t.Error(err)
	}

	if q.IsEmpty() {
		t.Error("should not be empty")
	}

	if err := q.AddByOverwrite(pkt); err != nil {
		t.Error(err)
	}

	if q.IsEmpty() {
		t.Error("should not be empty")
	}

	if len(q.messages) != 1 {
		t.Error("there should only be one item in the queue")
	}
}

func TestClientPktQueue_Steal(t *testing.T) {
	q := newClientPktQueue(10)
	_ = q.Add(&clientPacket{})
	if len(q.messages) != 1 {
		t.Error("there should only be one item in the queue")
	}

	q.Steal()
	if !q.IsEmpty() {
		t.Error("should be empty")
	}
}

func TestClientPktQueue_Try(t *testing.T) {
	q := newClientPktQueue(10)

	_ = q.Add(&clientPacket{})
	if err := q.Try(func(msg *clientPacket) error { return errors.New("") }); err == nil {
		t.Error("Try should fail when cb returns an error")
	}
	if len(q.messages) != 1 {
		t.Error("the number of entries in the queue should not reduce after a failed Try execution")
	}

	if err := q.Try(func(msg *clientPacket) error { return nil }); err != nil {
		t.Error("Try should not have failed", err)
	}
	if !q.IsEmpty() {
		t.Error("queue should be empty")
	}

	_ = q.Add(&clientPacket{})
	_ = q.Add(&clientPacket{})
	_ = q.Add(&clientPacket{})
	if err := q.Try(func(msg *clientPacket) error { return nil }); err != nil {
		t.Error("Try should not have failed", err)
	}
	if len(q.messages) != 2 {
		t.Error("the number of entries in the queue should reduce after Try execution")
	}
}
