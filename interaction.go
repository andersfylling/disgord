package disgord

import (
	"bytes"
	"mime/multipart"
	"strings"

	"github.com/andersfylling/disgord/internal/httd"
	"github.com/andersfylling/disgord/json"
)

type InteractionType = int

const (
	_ InteractionType = iota
	InteractionPing
	InteractionApplicationCommand
	InteractionMessageComponent
)

type OptionType = int

const (
	_ OptionType = iota
	OptionTypeSubCommand
	OptionTypeSubCommandGroup
	OptionTypeString
	OptionTypeInteger
	OptionTypeBoolean
	OptionTypeUser
	OptionTypeChannel
	OptionTypeRole
	OptionTypeMentionable
	OptionTypeNumber
)

type InteractionCallbackType = int

const (
	_ InteractionCallbackType = iota
	InteractionCallbackPong
	_
	_
	InteractionCallbackChannelMessageWithSource
	InteractionCallbackDeferredChannelMessageWithSource
	InteractionCallbackDeferredUpdateMessage
	InteractionCallbackUpdateMessage
)

// ApplicationCommandInteractionDataResolved
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ApplicationCommandInteractionDataResolved struct {
	Users    map[Snowflake]*User    `json:"users"`
	Members  map[Snowflake]*Member  `json:"members"`
	Roles    map[Snowflake]*Role    `json:"roles"`
	Channels map[Snowflake]*Channel `json:"channels"`
	Messages map[Snowflake]*Message `json:"messages"`
}

type ApplicationCommandInteractionData struct {
	ID            Snowflake                                  `json:"id"`
	Name          string                                     `json:"name"`
	Resolved      *ApplicationCommandInteractionDataResolved `json:"resolved"`
	Options       []*ApplicationCommandDataOption            `json:"options"`
	CustomID      string                                     `json:"custom_id"`
	Type          ApplicationCommandType                     `json:"type"`
	Values        []string                                   `json:"values"`
	ComponentType MessageComponentType                       `json:"component_type"`
	TargetID      Snowflake                                  `json:"target_id"`
}

type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User *User           `json:"user"`
}

type CreateInteractionResponseData struct {
	Content         string              `json:"content"`
	Tts             bool                `json:"tts,omitempty"`
	Embeds          []*Embed            `json:"embeds,omitempty"`
	Components      []*MessageComponent `json:"components"`
	Attachments     []*Attachment       `json:"attachments"`
	AllowedMentions *AllowedMentions    `json:"allowed_mentions,omitempty"`
	Flags           MessageFlag         `json:"flags,omitempty"` // Only SUPPRESS_EMBEDS and EPHEMERAL flags allowed.

	Files []CreateMessageFile `json:"-"`

	SpoilerTagContent        bool `json:"-"`
	SpoilerTagAllAttachments bool `json:"-"`
}

type CreateInteractionResponse struct {
	Type InteractionCallbackType        `json:"type"`
	Data *CreateInteractionResponseData `json:"data"`
}

func (res *CreateInteractionResponse) prepare() (postBody interface{}, contentType string, err error) {
	if res.Data == nil {
		return res, httd.ContentTypeJSON, nil
	}

	p := res.Data
	// spoiler tag
	if p.SpoilerTagContent && len(p.Content) > 0 {
		p.Content = "|| " + p.Content + " ||"
	}

	if len(p.Files) == 0 {
		postBody = res
		contentType = httd.ContentTypeJSON
		return
	}

	if p.SpoilerTagAllAttachments {
		for i := range p.Files {
			p.Files[i].SpoilerTag = true
		}
	}

	// check for spoilers
	for _, embed := range p.Embeds {
		for i := range p.Files {
			if p.Files[i].SpoilerTag && strings.Contains(embed.Image.URL, p.Files[i].FileName) {
				s := strings.Split(embed.Image.URL, p.Files[i].FileName)
				if len(s) > 0 {
					s[0] += AttachmentSpoilerPrefix + p.Files[i].FileName
					embed.Image.URL = strings.Join(s, "")
				}
			}
		}
	}

	// Set up a new multipart writer, as we'll be using this for the POST body instead
	buf := new(bytes.Buffer)
	mp := multipart.NewWriter(buf)

	// Write the existing JSON payload
	var payload []byte
	payload, err = json.Marshal(res)
	if err != nil {
		return
	}
	if err = mp.WriteField("payload_json", string(payload)); err != nil {
		return
	}

	// Iterate through all the files and write them to the multipart blob
	for i, file := range p.Files {
		if err = file.write(i, mp); err != nil {
			return
		}
	}

	mp.Close()

	postBody = buf
	contentType = mp.FormDataContentType()

	return
}
