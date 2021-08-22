package disgord

import "fmt"

type ApplicationCommandType = int

const (
	_ ApplicationCommandType = iota
	ApplicationCommandChatInput
	ApplicationCommandUser
	ApplicationCommandMessage
)

type ApplicationCommandOptionChoice struct {
	Name        string
	StringValue string
	NumberValue int
	DoubleValue float64
}

func (w *ApplicationCommandOptionChoice) MarshalJSON() ([]byte, error) {
	var data string
	if w.StringValue != "" {
		data = fmt.Sprintf(`{"name": "%s", "value": "%s"}`, w.Name, w.StringValue)
	}
	if w.NumberValue != 0 {
		data = fmt.Sprintf(`{"name": "%s", "value": "%d"}`, w.Name, w.NumberValue)
	}
	if w.DoubleValue != 0 {
		data = fmt.Sprintf(`{"name": "%s", "value": "%f"}`, w.Name, w.DoubleValue)
	}
	return []byte(data), nil
}

type ApplicationCommandOption struct {
	Type        OptionType                        `json:"type"`
	Name        string                            `json:"name"`
	Description string                            `json:"description"`
	Required    bool                              `json:"required"`
	Choices     []*ApplicationCommandOptionChoice `json:"choices"`
	Options     []*ApplicationCommandOption       `json:"options"`
}

type ApplicationCommandDataOption struct {
	Name    string                          `json:"name"`
	Type    OptionType                      `json:"type"`
	Value   interface{}                     `json:"value"`
	Options []*ApplicationCommandDataOption `json:"options"`
}

type ApplicationCommandPermissionType = int

const (
	_ ApplicationCommandPermissionType = iota
	ApplicationCommandPermissionRole
	ApplicationCommandPermissionUser
)

type ApplicationCommandPermissions struct {
	ID         Snowflake                        `json:"id"`
	Type       ApplicationCommandPermissionType `json:"type"`
	Permission bool                             `json:"permission"`
}
type GuildApplicationCommandPermissions struct {
	ID            Snowflake                        `json:"id"`
	ApplicationID Snowflake                        `json:"application_id"`
	GuildID       Snowflake                        `json:"guild_id"`
	Permissions   []*ApplicationCommandPermissions `json:"permissions"`
}

type ApplicationCommand struct {
	ID                Snowflake                   `json:"id"`
	Type              ApplicationCommandType      `json:"type"`
	ApplicationID     Snowflake                   `json:"application_id"`
	GuildID           Snowflake                   `json:"guild_id"`
	Name              string                      `json:"name"`
	Description       string                      `json:"description"`
	Options           []*ApplicationCommandOption `json:"options"`
	DefaultPermission bool                        `json:"default_permission,omitempty"`
}
