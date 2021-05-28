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
	SUB_COMMAND
	SUB_COMMAND_GROUP
	STRING
	INTEGER
	BOOLEAN
	USER
	CHANNEL
	ROLE
	MENTIONABLE
)

//TODO ApplicationCommandInteractionDataResolved https://discord.com/developers/docs/interactions/slash-commands#interaction-applicationcommandinteractiondataresolved
type ApplicationCommandInteractionDataResolved struct {
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                     `json:"name"`
	Type    OptionType                                 `json:"type"`
	Options []*ApplicationCommandInteractionDataOption `json:"options"`
}

type ApplicationCommandInteractionData struct {
	ID       Snowflake                                    `json:"id"`
	Name     string                                       `json:"name"`
	Resolved []*ApplicationCommandInteractionDataResolved `json:"resolved"`
	Options  []*ApplicationCommandInteractionDataOption   `json:"options"`
	CustomID string                                       `json:"custom_id"`
	Type     MessageComponentType                         `json:"component_type"`
}
