package disgordutil

import (
	"fmt"
	"github.com/andersfylling/disgord"
	"reflect"
	"sort"
	"strings"
)

func Sort(v interface{}, sortByField SortFieldType, sortOrder SortOrderType) {
	if v == nil {
		return
	}
	if sortByField == SortByID {
		sortByID(v, sortOrder)
	} else if sortByField == SortByGuildID {
		sortByGuildID(v, sortOrder)
	} else if sortByField == SortByChannelID {
		sortByChannelID(v, sortOrder)
	} else if sortByField == SortByName {
		sortByName(v, sortOrder)
	} else if sortByField == SortByHoist {
		sortByHoist(v, sortOrder)
	} else if list, ok := v.(sort.Interface); ok {
		if sortOrder == OrderDescending {
			sort.Sort(sort.Reverse(list))
		} else {
			sort.Sort(list)
		}
	} else {
		panic("type is missing sort.Interface implementation")
	}
}

func derefSliceP(v interface{}) (s interface{}) {
	switch t := v.(type) {
	case *[]*disgord.AuditLog:
		s = *t
	case *[]*disgord.AuditLogChanges:
		s = *t
	case *[]*disgord.AuditLogEntry:
		s = *t
	case *[]*disgord.AuditLogOption:
		s = *t
	case *[]*disgord.BasicCache:
		s = *t
	case *[]*disgord.AllowedMentions:
		s = *t
	case *[]*disgord.Attachment:
		s = *t
	case *[]*disgord.Channel:
		s = *t
	case *[]*disgord.CreateMessageFileParams:
		s = *t
	case *[]*disgord.CreateMessageParams:
		s = *t
	case *[]*disgord.CreateWebhookParams:
		s = *t
	case *[]*disgord.DeleteMessagesParams:
		s = *t
	case *[]*disgord.GetMessagesParams:
		s = *t
	case *[]*disgord.GroupDMParticipant:
		s = *t
	case *[]*disgord.PartialChannel:
		s = *t
	case *[]*disgord.PermissionOverwrite:
		s = *t
	case *[]*disgord.UpdateChannelPermissionsParams:
		s = *t
	case *[]*disgord.Client:
		s = *t
	case *[]*disgord.Config:
		s = *t
	case *[]*disgord.ErrorEmptyValue:
		s = *t
	case *[]*disgord.ErrorMissingSnowflake:
		s = *t
	case *[]*disgord.Embed:
		s = *t
	case *[]*disgord.EmbedAuthor:
		s = *t
	case *[]*disgord.EmbedField:
		s = *t
	case *[]*disgord.EmbedFooter:
		s = *t
	case *[]*disgord.EmbedImage:
		s = *t
	case *[]*disgord.EmbedProvider:
		s = *t
	case *[]*disgord.EmbedThumbnail:
		s = *t
	case *[]*disgord.EmbedVideo:
		s = *t
	case *[]*disgord.Emoji:
		s = *t
	case *[]*disgord.ChannelCreate:
		s = *t
	case *[]*disgord.ChannelDelete:
		s = *t
	case *[]*disgord.ChannelPinsUpdate:
		s = *t
	case *[]*disgord.ChannelUpdate:
		s = *t
	case *[]*disgord.GuildBanAdd:
		s = *t
	case *[]*disgord.GuildBanRemove:
		s = *t
	case *[]*disgord.GuildCreate:
		s = *t
	case *[]*disgord.GuildDelete:
		s = *t
	case *[]*disgord.GuildEmojisUpdate:
		s = *t
	case *[]*disgord.GuildIntegrationsUpdate:
		s = *t
	case *[]*disgord.GuildMemberAdd:
		s = *t
	case *[]*disgord.GuildMemberRemove:
		s = *t
	case *[]*disgord.GuildMemberUpdate:
		s = *t
	case *[]*disgord.GuildMembersChunk:
		s = *t
	case *[]*disgord.GuildRoleCreate:
		s = *t
	case *[]*disgord.GuildRoleDelete:
		s = *t
	case *[]*disgord.GuildRoleUpdate:
		s = *t
	case *[]*disgord.GuildUpdate:
		s = *t
	case *[]*disgord.InteractionCreate:
		s = *t
	case *[]*disgord.InviteCreate:
		s = *t
	case *[]*disgord.InviteDelete:
		s = *t
	case *[]*disgord.MessageCreate:
		s = *t
	case *[]*disgord.MessageDelete:
		s = *t
	case *[]*disgord.MessageDeleteBulk:
		s = *t
	case *[]*disgord.MessageReactionAdd:
		s = *t
	case *[]*disgord.MessageReactionRemove:
		s = *t
	case *[]*disgord.MessageReactionRemoveAll:
		s = *t
	case *[]*disgord.MessageReactionRemoveEmoji:
		s = *t
	case *[]*disgord.MessageUpdate:
		s = *t
	case *[]*disgord.PresenceUpdate:
		s = *t
	case *[]*disgord.Ready:
		s = *t
	case *[]*disgord.Resumed:
		s = *t
	case *[]*disgord.TypingStart:
		s = *t
	case *[]*disgord.UserUpdate:
		s = *t
	case *[]*disgord.VoiceServerUpdate:
		s = *t
	case *[]*disgord.VoiceStateUpdate:
		s = *t
	case *[]*disgord.WebhooksUpdate:
		s = *t
	case *[]*disgord.AddGuildMemberParams:
		s = *t
	case *[]*disgord.Ban:
		s = *t
	case *[]*disgord.BanMemberParams:
		s = *t
	case *[]*disgord.CreateGuildChannelParams:
		s = *t
	case *[]*disgord.CreateGuildEmojiParams:
		s = *t
	case *[]*disgord.CreateGuildIntegrationParams:
		s = *t
	case *[]*disgord.CreateGuildParams:
		s = *t
	case *[]*disgord.CreateGuildRoleParams:
		s = *t
	case *[]*disgord.GetMembersParams:
		s = *t
	case *[]*disgord.Guild:
		s = *t
	case *[]*disgord.GuildEmbed:
		s = *t
	case *[]*disgord.GuildUnavailable:
		s = *t
	case *[]*disgord.Integration:
		s = *t
	case *[]*disgord.IntegrationAccount:
		s = *t
	case *[]*disgord.Member:
		s = *t
	case *[]*disgord.PartialBan:
		s = *t
	case *[]*disgord.UpdateGuildChannelPositionsParams:
		s = *t
	case *[]*disgord.UpdateGuildIntegrationParams:
		s = *t
	case *[]*disgord.UpdateGuildRolePositionsParams:
		s = *t
	case *[]*disgord.ApplicationCommandInteractionData:
		s = *t
	case *[]*disgord.ApplicationCommandInteractionDataOption:
		s = *t
	case *[]*disgord.ApplicationCommandInteractionDataResolved:
		s = *t
	case *[]*disgord.InteractionApplicationCommandCallbackData:
		s = *t
	case *[]*disgord.InteractionResponse:
		s = *t
	case *[]*disgord.MessageInteraction:
		s = *t
	case *[]*disgord.Invite:
		s = *t
	case *[]*disgord.InviteMetadata:
		s = *t
	case *[]*disgord.MentionChannel:
		s = *t
	case *[]*disgord.Message:
		s = *t
	case *[]*disgord.MessageActivity:
		s = *t
	case *[]*disgord.MessageApplication:
		s = *t
	case *[]*disgord.MessageComponent:
		s = *t
	case *[]*disgord.MessageReference:
		s = *t
	case *[]*disgord.MessageSticker:
		s = *t
	case *[]*disgord.StickerItem:
		s = *t
	case *[]*disgord.GetReactionURLParams:
		s = *t
	case *[]*disgord.Reaction:
		s = *t
	case *[]*disgord.Ctrl:
		s = *t
	case *[]*disgord.RESTBuilder:
		s = *t
	case *[]*disgord.Role:
		s = *t
	case *[]*disgord.ErrorUnsupportedType:
		s = *t
	case *[]*disgord.Time:
		s = *t
	case *[]*disgord.Activity:
		s = *t
	case *[]*disgord.ActivityAssets:
		s = *t
	case *[]*disgord.ActivityEmoji:
		s = *t
	case *[]*disgord.ActivityParty:
		s = *t
	case *[]*disgord.ActivitySecrets:
		s = *t
	case *[]*disgord.ActivityTimestamp:
		s = *t
	case *[]*disgord.ClientStatus:
		s = *t
	case *[]*disgord.CreateGroupDMParams:
		s = *t
	case *[]*disgord.GetCurrentUserGuildsParams:
		s = *t
	case *[]*disgord.User:
		s = *t
	case *[]*disgord.UserConnection:
		s = *t
	case *[]*disgord.UserPresence:
		s = *t
	case *[]*disgord.VoiceRegion:
		s = *t
	case *[]*disgord.VoiceState:
		s = *t
	case *[]*disgord.ExecuteWebhookParams:
		s = *t
	case *[]*disgord.Webhook:
		s = *t
	default:
		s = t
	}

	return s
}
func sortByID(v interface{}, sortOrder SortOrderType) {
	v = derefSliceP(v)
	if !reflectIsSlice(v) {
		return
	}

	var less func(i, j int) bool
	switch s := v.(type) {
	case []*disgord.AuditLogEntry:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.AuditLogOption:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Attachment:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Channel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.PartialChannel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.PermissionOverwrite:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Emoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.InteractionCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.CreateGuildIntegrationParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Guild:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.GuildUnavailable:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Integration:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.IntegrationAccount:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.UpdateGuildChannelPositionsParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.UpdateGuildRolePositionsParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.ApplicationCommandInteractionData:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.MessageInteraction:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.MentionChannel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Message:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.MessageApplication:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.MessageSticker:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.StickerItem:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Role:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.ActivityEmoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.ActivityParty:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.User:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.UserConnection:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.VoiceRegion:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	case []*disgord.Webhook:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ID > s[j].ID }
		} else {
			less = func(i, j int) bool { return s[i].ID < s[j].ID }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", s))
	}

	sort.Slice(v, less)
}
func sortByGuildID(v interface{}, sortOrder SortOrderType) {
	v = derefSliceP(v)
	if !reflectIsSlice(v) {
		return
	}

	var less func(i, j int) bool
	switch s := v.(type) {
	case []*disgord.Channel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.ChannelPinsUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildBanAdd:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildBanRemove:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildEmojisUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildIntegrationsUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildMemberRemove:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildMembersChunk:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildRoleCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildRoleDelete:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.GuildRoleUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.InteractionCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.InviteCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.InviteDelete:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.MessageDelete:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.MessageReactionRemoveEmoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.PresenceUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.TypingStart:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.VoiceServerUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.WebhooksUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.Member:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.MentionChannel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.Message:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.MessageReference:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.UserPresence:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.VoiceState:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	case []*disgord.Webhook:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].GuildID > s[j].GuildID }
		} else {
			less = func(i, j int) bool { return s[i].GuildID < s[j].GuildID }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", s))
	}

	sort.Slice(v, less)
}
func sortByChannelID(v interface{}, sortOrder SortOrderType) {
	v = derefSliceP(v)
	if !reflectIsSlice(v) {
		return
	}

	var less func(i, j int) bool
	switch s := v.(type) {
	case []*disgord.AuditLogOption:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.ChannelPinsUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.InteractionCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.InviteCreate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.InviteDelete:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageDelete:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageDeleteBulk:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageReactionAdd:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageReactionRemove:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageReactionRemoveAll:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageReactionRemoveEmoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.TypingStart:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.WebhooksUpdate:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.GuildEmbed:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.Message:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.MessageReference:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.VoiceState:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	case []*disgord.Webhook:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].ChannelID > s[j].ChannelID }
		} else {
			less = func(i, j int) bool { return s[i].ChannelID < s[j].ChannelID }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", s))
	}

	sort.Slice(v, less)
}
func sortByName(v interface{}, sortOrder SortOrderType) {
	v = derefSliceP(v)
	if !reflectIsSlice(v) {
		return
	}

	var less func(i, j int) bool
	switch s := v.(type) {
	case []*disgord.Channel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.CreateWebhookParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.PartialChannel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.EmbedAuthor:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.EmbedField:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.EmbedProvider:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Emoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.CreateGuildChannelParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.CreateGuildEmojiParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.CreateGuildParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.CreateGuildRoleParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Guild:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Integration:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.IntegrationAccount:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.ApplicationCommandInteractionData:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.ApplicationCommandInteractionDataOption:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.MessageInteraction:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.MentionChannel:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.MessageApplication:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.MessageSticker:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.StickerItem:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Role:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Activity:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.ActivityEmoji:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.UserConnection:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.VoiceRegion:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	case []*disgord.Webhook:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) > strings.ToLower(s[j].Name) }
		} else {
			less = func(i, j int) bool { return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name) }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", s))
	}

	sort.Slice(v, less)
}
func sortByHoist(v interface{}, sortOrder SortOrderType) {
	v = derefSliceP(v)
	if !reflectIsSlice(v) {
		return
	}

	var less func(i, j int) bool
	switch s := v.(type) {
	case []*disgord.CreateGuildRoleParams:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].Hoist && !s[j].Hoist }
		} else {
			less = func(i, j int) bool { return !s[i].Hoist && s[j].Hoist }
		}
	case []*disgord.Role:
		if sortOrder == OrderDescending {
			less = func(i, j int) bool { return s[i].Hoist && !s[j].Hoist }
		} else {
			less = func(i, j int) bool { return !s[i].Hoist && s[j].Hoist }
		}
	default:
		panic(fmt.Sprintf("type %+v does not support sorting", s))
	}

	sort.Slice(v, less)
}

func reflectIsSlice(v interface{}) bool {
	ValueIface := reflect.ValueOf(v)
	kind := ValueIface.Type().Kind()
	return kind == reflect.Slice
}

// Reflect if an interface is either a struct or a pointer to a struct
// and has the defined member field, if error is nil, the given
// FieldName exists and is accessible with reflect.
func reflectStructField(Iface interface{}, FieldName string) error {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return fmt.Errorf("Interface `%s` does not have the field `%s`", ValueIface.Type(), FieldName)
	}
	return nil
}
