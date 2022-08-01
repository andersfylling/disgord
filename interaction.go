package disgord

import (
	"bytes"
	"context"
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
	InteractionApplicationCommandAutocomplete
	InteractionModalSubmit
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
	InteractionCallbackApplicationCommandAutocompleteResult
	InteractionCallbackModal
)

type Interactable interface {
	GetID() Snowflake
	GetToken() string
	Edit(context.Context, Session, *UpdateMessage) error
	Reply(context.Context, Session, *CreateInteractionResponse) error
}

// ApplicationCommandInteractionDataResolved
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-resolved-data-structure
type ApplicationCommandInteractionDataResolved struct {
	Users    map[Snowflake]*User    `json:"users"`
	Members  map[Snowflake]*Member  `json:"members"`
	Roles    map[Snowflake]*Role    `json:"roles"`
	Channels map[Snowflake]*Channel `json:"channels"`
	Messages map[Snowflake]*Message `json:"messages"`
}

// ApplicationCommandInteractionData
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-application-command-data-structure
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
	Components    []*MessageComponent                        `json:"components"`
}

// MessageComponentInteractionData
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-message-component-data-structure
type MessageComponentInteractionData struct {
	CustomID string               `json:"custom_id"`
	Type     MessageComponentType `json:"type"`
	Values   []*SelectMenuOption  `json:"values,omitempty"`
}

// ModalSubmitInteractionData
// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-object-modal-submit-data-structure
type ModalSubmitInteractionData struct{}

type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User *User           `json:"user"`
}

type InteractionResponseData interface {
	prepareData(*CreateInteractionResponse) (interface{}, string, error)
}

type ModalInteractionResponseData struct {
	CustomID   string              `json:"custom_id"`
	Title      string              `json:"title"`
	Components []*MessageComponent `json:"components"`
}

func (d *ModalInteractionResponseData) prepareData(res *CreateInteractionResponse) (interface{}, string, error) {
	return res, httd.ContentTypeJSON, nil
}

type MessageInteractionResponseData struct {
	Content         string              `json:"content"`
	Title           string              `json:"title"`
	CustomID        string              `json:"custom_id"`
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

func (d *MessageInteractionResponseData) prepareData(res *CreateInteractionResponse) (postBody interface{}, contentType string, err error) {
	// spoiler tag
	if d.SpoilerTagContent && len(d.Content) > 0 {
		d.Content = "|| " + d.Content + " ||"
	}

	if len(d.Files) == 0 {
		postBody = res
		contentType = httd.ContentTypeJSON
		return
	}

	if d.SpoilerTagAllAttachments {
		for i := range d.Files {
			d.Files[i].SpoilerTag = true
		}
	}

	// check for spoilers
	for _, embed := range d.Embeds {
		for i := range d.Files {
			if d.Files[i].SpoilerTag && strings.Contains(embed.Image.URL, d.Files[i].FileName) {
				s := strings.Split(embed.Image.URL, d.Files[i].FileName)
				if len(s) > 0 {
					s[0] += AttachmentSpoilerPrefix + d.Files[i].FileName
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
	for i, file := range d.Files {
		if err = file.write(i, mp); err != nil {
			return
		}
	}

	mp.Close()

	postBody = buf
	contentType = mp.FormDataContentType()
	return
}

type CreateInteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
	Data InteractionResponseData `json:"data"`
}

func (res *CreateInteractionResponse) prepare() (postBody interface{}, contentType string, err error) {
	if res.Data == nil {
		return res, httd.ContentTypeJSON, nil
	}

	return res.Data.prepareData(res)
}
