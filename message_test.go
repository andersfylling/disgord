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

func TestMessage_DeepCopy(t *testing.T) {
	// no nice way of comparing structs ..
	// original := &Message{}
	// wrapper := &struct {
	// 	D *Message `json:"d"`
	// }{original}
	// if data, err := ioutil.ReadFile("./testdata/phases/startup-smooth-1/13_0_MESSAGE_CREATE.json"); err != nil {
	// 	t.Fatal(err)
	// } else {
	// 	if err = httd.Unmarshal(data, wrapper); err != nil {
	// 		t.Fatal(err)
	// 	}
	// }
	// c := original.DeepCopy().(*Message)
	//
	// prettyPrint := func(i interface{}) string {
	// 	s, _ := json.MarshalIndent(i, "", "\t")
	// 	return string(s)
	// }

	// fmt.Println(reflect.DeepEqual([]*Message{}, nil))
	//
	// if !reflect.DeepEqual(original, c) {
	// 	t.Errorf("expect messages to be equal after deep copy.\n Got \n%s,\n\n wants \n%s", prettyPrint(c), prettyPrint(original))
	// }
}
