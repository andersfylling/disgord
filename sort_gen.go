package disgord

import (
	"fmt"
	"sort"
)

//////////////////////////////////////////////////////
//
// demultiplexer
//
//////////////////////////////////////////////////////

func Sort(s interface{}, f ...Flag) {
	flags := mergeFlags(f)
	if (flags & SortByID) > 0 {
		sortByID(s, f...)
	}
	if (flags & SortByName) > 0 {
		sortByName(s, f...)
	} else if list, ok := s.(sort.Interface); ok {
		sort.Sort(list) // TODO: asc/desc
	} else {
		panic("type is missing sort.Interface implementation")
	}
}
func sortByID(s interface{}, flags ...Flag) {
	var descending bool
	if (mergeFlags(flags) & OrderDescending) > 0 {
		descending = true
	}

	var less func(i, j int) bool
	switch t := s.(type) {
	case []*AuditLogEntry:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*AuditLogOption:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Attachment:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Channel:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*PartialChannel:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*PermissionOverwrite:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Emoji:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*CreateGuildIntegrationParams:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Guild:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*GuildUnavailable:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Integration:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*IntegrationAccount:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*UpdateGuildChannelPositionsParams:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*UpdateGuildRolePositionsParams:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Message:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*MessageApplication:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*rest:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Role:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*ActivityParty:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*User:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*UserConnection:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*userJSON:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*VoiceRegion:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	case []*Webhook:
		if descending {
			less = func(i, j int) bool { return t[i].ID > t[j].ID }
		} else {
			less = func(i, j int) bool { return t[i].ID < t[j].ID }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", t))
	}

	sort.Slice(s, less)
}
func sortByName(s interface{}, flags ...Flag) {
	var descending bool
	if (mergeFlags(flags) & OrderDescending) > 0 {
		descending = true
	}

	var less func(i, j int) bool
	switch t := s.(type) {
	case []*Channel:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*PartialChannel:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*EmbedAuthor:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*EmbedField:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*EmbedProvider:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*CreateGuildEmojiParams:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Emoji:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*CreateGuildChannelParams:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*CreateGuildParams:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Guild:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Integration:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*IntegrationAccount:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*MessageApplication:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*CreateGuildRoleParams:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Role:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Activity:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*UserConnection:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*VoiceRegion:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*CreateWebhookParams:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	case []*Webhook:
		if descending {
			less = func(i, j int) bool { return t[i].Name > t[j].Name }
		} else {
			less = func(i, j int) bool { return t[i].Name < t[j].Name }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", t))
	}

	sort.Slice(s, less)
}
