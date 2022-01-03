package disgord

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

// ApplicationCommandInteractionDataResolved ..
//TODO implement https://discord.com/developers/docs/interactions/slash-commands#interaction-applicationcommandinteractiondataresolved
type ApplicationCommandInteractionDataResolved struct {
}

type ApplicationCommandInteractionData struct {
	ID       Snowflake                                    `json:"id"`
	Name     string                                       `json:"name"`
	Resolved []*ApplicationCommandInteractionDataResolved `json:"resolved"`
	Options  []*ApplicationCommandDataOption              `json:"options"`
	CustomID string                                       `json:"custom_id"`
	Type     MessageComponentType                         `json:"component_type"`
	Values   []*SelectMenuOption                          `json:"values"`
}

type MessageInteraction struct {
	ID   Snowflake       `json:"id"`
	Type InteractionType `json:"type"`
	Name string          `json:"name"`
	User *User           `json:"user"`
}

type InteractionApplicationCommandCallbackData struct {
	Tts             bool                `json:"tts"`
	Content         string              `json:"content"`
	Embeds          []*Embed            `json:"embeds"`
	Flags           int                 `json:"flags"`
	AllowedMentions *AllowedMentions    `json:"allowed_mentions"`
	Components      []*MessageComponent `json:"components"`
}

type InteractionResponse struct {
	Type InteractionCallbackType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data"`
}
