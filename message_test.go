// +build !integration

package disgord

import (
	"testing"
)

func TestMessage_updateInternals(t *testing.T) {
	m := &Message{}
	m.updateInternals()

	if m.SpoilerTagContent {
		t.Error("expects spoiler tag for content to be false. Got true")
	}
	if m.SpoilerTagAllAttachments {
		t.Error("expects spoiler tag for attachments to be false. Got true")
	}

	m.Content = "||||"
	m.updateInternals()
	if !m.SpoilerTagContent {
		t.Error("expects spoiler tag for content to be true. Got false")
	}

	m.Content = "|.||"
	m.updateInternals()
	if m.SpoilerTagContent {
		t.Error("expects spoiler tag for content to be false. Got true")
	}

	m.Content = "|| testing ||"
	m.updateInternals()
	if !m.SpoilerTagContent {
		t.Error("expects spoiler tag for content to be true. Got false")
	}

	m.Attachments = append(m.Attachments, &Attachment{
		Filename: AttachmentSpoilerPrefix,
	})
	m.updateInternals()
	if !m.SpoilerTagAllAttachments {
		t.Error("expects spoiler tag for attachments to be true. Got false")
	}

	m.Attachments = append(m.Attachments, &Attachment{
		Filename: "random",
	})
	m.updateInternals()
	if m.SpoilerTagAllAttachments {
		t.Error("expects spoiler tag for attachments to be false. Got true")
	}
}
