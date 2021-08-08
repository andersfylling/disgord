package disgord

import "github.com/andersfylling/disgord/internal/httd"

// Deprecated: use Update instead
func (m messageQueryBuilder) UpdateBuilder() UpdateMessageBuilder {
	builder := &updateMessageBuilder{}
	builder.r.itemFactory = func() interface{} {
		return &Message{}
	}
	builder.r.flags = m.flags
	builder.r.addPrereq(m.cid.IsZero(), "channelID must be set to get channel messages")
	builder.r.addPrereq(m.mid.IsZero(), "msgID must be set to edit the message")
	builder.r.setup(m.client.req, &httd.Request{
		Method:      httd.MethodPatch,
		Ctx:         m.ctx,
		Endpoint:    "/channels/" + m.cid.String() + "/messages/" + m.mid.String(),
		ContentType: httd.ContentTypeJSON,
	}, nil)

	return builder
}

//////////////////////////////////////////////////////
//
// REST Builders
//
//////////////////////////////////////////////////////

// updateMessageBuilder, params here
//  https://discord.com/developers/docs/resources/channel#edit-message-json-params
//generate-rest-params: content:string, embed:*Embed,
//generate-rest-basic-execute: message:*Message,
type updateMessageBuilder struct {
	r RESTBuilder
}

// SetAllowedMentions sets the allowed mentions for the updateMessageBuilder then returns the builder to allow chaining.
func (b *updateMessageBuilder) SetAllowedMentions(mentions *AllowedMentions) *updateMessageBuilder {
	b.r.param("allowed_mentions", mentions)
	return b
}

//////////////////////////////////////////////////////
//
// REST Wrappers
//
//////////////////////////////////////////////////////

// Deprecated: use Update instead
func (m messageQueryBuilder) SetContent(content string) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetContent(content).
		Execute()
}

// Deprecated: use Update instead
func (m messageQueryBuilder) SetEmbed(embed *Embed) (*Message, error) {
	builder := m.WithContext(m.ctx).UpdateBuilder()
	return builder.
		SetEmbed(embed).
		Execute()
}
