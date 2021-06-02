// Code generated by generate/interfaces; DO NOT EDIT.

package disgord

func (m *Member) copyOverTo(other interface{}) error {
	var dest *Member
	var valid bool
	if dest, valid = other.(*Member); !valid {
		return newErrorUnsupportedType("argument given is not a *Member type")
	}
	dest.GuildID = m.GuildID
	dest.User = m.User
	dest.Nick = m.Nick
	dest.Roles = make([]Snowflake, len(m.Roles))
	copy(dest.Roles, m.Roles)
	dest.JoinedAt = m.JoinedAt
	dest.PremiumSince = m.PremiumSince
	dest.Deaf = m.Deaf
	dest.Mute = m.Mute
	dest.Pending = m.Pending
	dest.UserID = m.UserID

	return nil
}

func (g *Guild) copyOverTo(other interface{}) error {
	var dest *Guild
	var valid bool
	if dest, valid = other.(*Guild); !valid {
		return newErrorUnsupportedType("argument given is not a *Guild type")
	}
	dest.ID = g.ID
	dest.ApplicationID = g.ApplicationID
	dest.Name = g.Name
	dest.Icon = g.Icon
	dest.Splash = g.Splash
	dest.Owner = g.Owner
	dest.OwnerID = g.OwnerID
	dest.Permissions = g.Permissions
	dest.Region = g.Region
	dest.AfkChannelID = g.AfkChannelID
	dest.AfkTimeout = g.AfkTimeout
	dest.VerificationLevel = g.VerificationLevel
	dest.DefaultMessageNotifications = g.DefaultMessageNotifications
	dest.ExplicitContentFilter = g.ExplicitContentFilter
	dest.Roles = make([]*Role, len(g.Roles))
	for i := 0; i < len(g.Roles); i++ {
		dest.Roles[i] = DeepCopy(g.Roles[i]).(*Role)
	}
	dest.Emojis = make([]*Emoji, len(g.Emojis))
	for i := 0; i < len(g.Emojis); i++ {
		dest.Emojis[i] = DeepCopy(g.Emojis[i]).(*Emoji)
	}
	dest.Features = make([]string, len(g.Features))
	copy(dest.Features, g.Features)
	dest.MFALevel = g.MFALevel
	dest.WidgetEnabled = g.WidgetEnabled
	dest.WidgetChannelID = g.WidgetChannelID
	dest.SystemChannelID = g.SystemChannelID
	dest.JoinedAt = g.JoinedAt
	dest.Large = g.Large
	dest.Unavailable = g.Unavailable
	dest.MemberCount = g.MemberCount
	dest.VoiceStates = make([]*VoiceState, len(g.VoiceStates))
	for i := 0; i < len(g.VoiceStates); i++ {
		dest.VoiceStates[i] = DeepCopy(g.VoiceStates[i]).(*VoiceState)
	}
	dest.Members = make([]*Member, len(g.Members))
	for i := 0; i < len(g.Members); i++ {
		dest.Members[i] = DeepCopy(g.Members[i]).(*Member)
	}
	dest.Channels = make([]*Channel, len(g.Channels))
	for i := 0; i < len(g.Channels); i++ {
		dest.Channels[i] = DeepCopy(g.Channels[i]).(*Channel)
	}
	dest.Presences = make([]*UserPresence, len(g.Presences))
	for i := 0; i < len(g.Presences); i++ {
		dest.Presences[i] = DeepCopy(g.Presences[i]).(*UserPresence)
	}

	return nil
}

func (r *Reaction) copyOverTo(other interface{}) error {
	var dest *Reaction
	var valid bool
	if dest, valid = other.(*Reaction); !valid {
		return newErrorUnsupportedType("argument given is not a *Reaction type")
	}
	dest.Count = r.Count
	dest.Me = r.Me
	dest.Emoji = r.Emoji

	return nil
}

func (e *EmbedVideo) copyOverTo(other interface{}) error {
	var dest *EmbedVideo
	var valid bool
	if dest, valid = other.(*EmbedVideo); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedVideo type")
	}
	dest.URL = e.URL
	dest.Height = e.Height
	dest.Width = e.Width

	return nil
}

func (a *AuditLogOption) copyOverTo(other interface{}) error {
	var dest *AuditLogOption
	var valid bool
	if dest, valid = other.(*AuditLogOption); !valid {
		return newErrorUnsupportedType("argument given is not a *AuditLogOption type")
	}
	dest.DeleteMemberDays = a.DeleteMemberDays
	dest.MembersRemoved = a.MembersRemoved
	dest.ChannelID = a.ChannelID
	dest.Count = a.Count
	dest.ID = a.ID
	dest.Type = a.Type
	dest.RoleName = a.RoleName

	return nil
}

func (u *UserPresence) copyOverTo(other interface{}) error {
	var dest *UserPresence
	var valid bool
	if dest, valid = other.(*UserPresence); !valid {
		return newErrorUnsupportedType("argument given is not a *UserPresence type")
	}
	dest.User = u.User
	dest.Roles = make([]Snowflake, len(u.Roles))
	copy(dest.Roles, u.Roles)
	dest.Game = u.Game
	dest.GuildID = u.GuildID
	dest.Nick = u.Nick
	dest.Status = u.Status

	return nil
}

func (i *Integration) copyOverTo(other interface{}) error {
	var dest *Integration
	var valid bool
	if dest, valid = other.(*Integration); !valid {
		return newErrorUnsupportedType("argument given is not a *Integration type")
	}
	dest.ID = i.ID
	dest.Name = i.Name
	dest.Type = i.Type
	dest.Enabled = i.Enabled
	dest.Syncing = i.Syncing
	dest.RoleID = i.RoleID
	dest.ExpireBehavior = i.ExpireBehavior
	dest.ExpireGracePeriod = i.ExpireGracePeriod
	dest.User = i.User
	dest.Account = i.Account

	return nil
}

func (u *User) copyOverTo(other interface{}) error {
	var dest *User
	var valid bool
	if dest, valid = other.(*User); !valid {
		return newErrorUnsupportedType("argument given is not a *User type")
	}
	dest.ID = u.ID
	dest.Username = u.Username
	dest.Discriminator = u.Discriminator
	dest.Avatar = u.Avatar
	dest.Bot = u.Bot
	dest.System = u.System
	dest.MFAEnabled = u.MFAEnabled
	dest.Locale = u.Locale
	dest.Verified = u.Verified
	dest.Email = u.Email
	dest.Flags = u.Flags
	dest.PremiumType = u.PremiumType
	dest.PublicFlags = u.PublicFlags
	dest.PartialMember = u.PartialMember

	return nil
}

func (a *Attachment) copyOverTo(other interface{}) error {
	var dest *Attachment
	var valid bool
	if dest, valid = other.(*Attachment); !valid {
		return newErrorUnsupportedType("argument given is not a *Attachment type")
	}
	dest.ID = a.ID
	dest.Filename = a.Filename
	dest.Size = a.Size
	dest.URL = a.URL
	dest.ProxyURL = a.ProxyURL
	dest.Height = a.Height
	dest.Width = a.Width
	dest.SpoilerTag = a.SpoilerTag

	return nil
}

func (a *AuditLogChanges) copyOverTo(other interface{}) error {
	var dest *AuditLogChanges
	var valid bool
	if dest, valid = other.(*AuditLogChanges); !valid {
		return newErrorUnsupportedType("argument given is not a *AuditLogChanges type")
	}
	dest.NewValue = a.NewValue
	dest.OldValue = a.OldValue
	dest.Key = a.Key

	return nil
}

func (g *GuildEmbed) copyOverTo(other interface{}) error {
	var dest *GuildEmbed
	var valid bool
	if dest, valid = other.(*GuildEmbed); !valid {
		return newErrorUnsupportedType("argument given is not a *GuildEmbed type")
	}
	dest.Enabled = g.Enabled
	dest.ChannelID = g.ChannelID

	return nil
}

func (a *Activity) copyOverTo(other interface{}) error {
	var dest *Activity
	var valid bool
	if dest, valid = other.(*Activity); !valid {
		return newErrorUnsupportedType("argument given is not a *Activity type")
	}
	dest.Name = a.Name
	dest.Type = a.Type
	dest.URL = a.URL
	dest.CreatedAt = a.CreatedAt
	dest.Timestamps = a.Timestamps
	dest.ApplicationID = a.ApplicationID
	dest.Details = a.Details
	dest.State = a.State
	dest.Emoji = a.Emoji
	dest.Party = a.Party
	dest.Assets = a.Assets
	dest.Secrets = a.Secrets
	dest.Instance = a.Instance
	dest.Flags = a.Flags

	return nil
}

func (w *Webhook) copyOverTo(other interface{}) error {
	var dest *Webhook
	var valid bool
	if dest, valid = other.(*Webhook); !valid {
		return newErrorUnsupportedType("argument given is not a *Webhook type")
	}
	dest.ID = w.ID
	dest.GuildID = w.GuildID
	dest.ChannelID = w.ChannelID
	dest.User = w.User
	dest.Name = w.Name
	dest.Avatar = w.Avatar
	dest.Token = w.Token

	return nil
}

func (e *EmbedField) copyOverTo(other interface{}) error {
	var dest *EmbedField
	var valid bool
	if dest, valid = other.(*EmbedField); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedField type")
	}
	dest.Name = e.Name
	dest.Value = e.Value
	dest.Inline = e.Inline

	return nil
}

func (v *VoiceRegion) copyOverTo(other interface{}) error {
	var dest *VoiceRegion
	var valid bool
	if dest, valid = other.(*VoiceRegion); !valid {
		return newErrorUnsupportedType("argument given is not a *VoiceRegion type")
	}
	dest.ID = v.ID
	dest.Name = v.Name
	dest.SampleHostname = v.SampleHostname
	dest.SamplePort = v.SamplePort
	dest.VIP = v.VIP
	dest.Optimal = v.Optimal
	dest.Deprecated = v.Deprecated
	dest.Custom = v.Custom

	return nil
}

func (e *EmbedFooter) copyOverTo(other interface{}) error {
	var dest *EmbedFooter
	var valid bool
	if dest, valid = other.(*EmbedFooter); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedFooter type")
	}
	dest.Text = e.Text
	dest.IconURL = e.IconURL
	dest.ProxyIconURL = e.ProxyIconURL

	return nil
}

func (b *Ban) copyOverTo(other interface{}) error {
	var dest *Ban
	var valid bool
	if dest, valid = other.(*Ban); !valid {
		return newErrorUnsupportedType("argument given is not a *Ban type")
	}
	dest.Reason = b.Reason
	dest.User = b.User

	return nil
}

func (a *ActivityParty) copyOverTo(other interface{}) error {
	var dest *ActivityParty
	var valid bool
	if dest, valid = other.(*ActivityParty); !valid {
		return newErrorUnsupportedType("argument given is not a *ActivityParty type")
	}
	dest.ID = a.ID
	dest.Size = make([]int, len(a.Size))
	copy(dest.Size, a.Size)

	return nil
}

func (m *MessageComponent) copyOverTo(other interface{}) error {
	var dest *MessageComponent
	var valid bool
	if dest, valid = other.(*MessageComponent); !valid {
		return newErrorUnsupportedType("argument given is not a *MessageComponent type")
	}
	dest.Type = m.Type
	dest.Style = m.Style
	dest.Label = m.Label
	dest.Emoji = m.Emoji
	dest.CustomID = m.CustomID
	dest.Url = m.Url
	dest.Disabled = m.Disabled
	dest.Components = make([]*MessageComponent, len(m.Components))
	for i := 0; i < len(m.Components); i++ {
		dest.Components[i] = DeepCopy(m.Components[i]).(*MessageComponent)
	}

	return nil
}

func (a *ActivityEmoji) copyOverTo(other interface{}) error {
	var dest *ActivityEmoji
	var valid bool
	if dest, valid = other.(*ActivityEmoji); !valid {
		return newErrorUnsupportedType("argument given is not a *ActivityEmoji type")
	}
	dest.Name = a.Name
	dest.ID = a.ID
	dest.Animated = a.Animated

	return nil
}

func (u *UserConnection) copyOverTo(other interface{}) error {
	var dest *UserConnection
	var valid bool
	if dest, valid = other.(*UserConnection); !valid {
		return newErrorUnsupportedType("argument given is not a *UserConnection type")
	}
	dest.ID = u.ID
	dest.Name = u.Name
	dest.Type = u.Type
	dest.Revoked = u.Revoked
	dest.Integrations = make([]*IntegrationAccount, len(u.Integrations))
	for i := 0; i < len(u.Integrations); i++ {
		dest.Integrations[i] = DeepCopy(u.Integrations[i]).(*IntegrationAccount)
	}

	return nil
}

func (r *Role) copyOverTo(other interface{}) error {
	var dest *Role
	var valid bool
	if dest, valid = other.(*Role); !valid {
		return newErrorUnsupportedType("argument given is not a *Role type")
	}
	dest.ID = r.ID
	dest.Name = r.Name
	dest.Color = r.Color
	dest.Hoist = r.Hoist
	dest.Position = r.Position
	dest.Permissions = r.Permissions
	dest.Managed = r.Managed
	dest.Mentionable = r.Mentionable
	dest.guildID = r.guildID

	return nil
}

func (v *VoiceState) copyOverTo(other interface{}) error {
	var dest *VoiceState
	var valid bool
	if dest, valid = other.(*VoiceState); !valid {
		return newErrorUnsupportedType("argument given is not a *VoiceState type")
	}
	dest.GuildID = v.GuildID
	dest.ChannelID = v.ChannelID
	dest.UserID = v.UserID
	dest.Member = v.Member
	dest.SessionID = v.SessionID
	dest.Deaf = v.Deaf
	dest.Mute = v.Mute
	dest.SelfDeaf = v.SelfDeaf
	dest.SelfMute = v.SelfMute
	dest.Suppress = v.Suppress

	return nil
}

func (a *AuditLog) copyOverTo(other interface{}) error {
	var dest *AuditLog
	var valid bool
	if dest, valid = other.(*AuditLog); !valid {
		return newErrorUnsupportedType("argument given is not a *AuditLog type")
	}
	dest.Webhooks = make([]*Webhook, len(a.Webhooks))
	for i := 0; i < len(a.Webhooks); i++ {
		dest.Webhooks[i] = DeepCopy(a.Webhooks[i]).(*Webhook)
	}
	dest.Users = make([]*User, len(a.Users))
	for i := 0; i < len(a.Users); i++ {
		dest.Users[i] = DeepCopy(a.Users[i]).(*User)
	}
	dest.AuditLogEntries = make([]*AuditLogEntry, len(a.AuditLogEntries))
	for i := 0; i < len(a.AuditLogEntries); i++ {
		dest.AuditLogEntries[i] = DeepCopy(a.AuditLogEntries[i]).(*AuditLogEntry)
	}

	return nil
}

func (e *EmbedAuthor) copyOverTo(other interface{}) error {
	var dest *EmbedAuthor
	var valid bool
	if dest, valid = other.(*EmbedAuthor); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedAuthor type")
	}
	dest.Name = e.Name
	dest.URL = e.URL
	dest.IconURL = e.IconURL
	dest.ProxyIconURL = e.ProxyIconURL

	return nil
}

func (i *IntegrationAccount) copyOverTo(other interface{}) error {
	var dest *IntegrationAccount
	var valid bool
	if dest, valid = other.(*IntegrationAccount); !valid {
		return newErrorUnsupportedType("argument given is not a *IntegrationAccount type")
	}
	dest.ID = i.ID
	dest.Name = i.Name

	return nil
}

func (a *ActivitySecrets) copyOverTo(other interface{}) error {
	var dest *ActivitySecrets
	var valid bool
	if dest, valid = other.(*ActivitySecrets); !valid {
		return newErrorUnsupportedType("argument given is not a *ActivitySecrets type")
	}
	dest.Join = a.Join
	dest.Spectate = a.Spectate
	dest.Match = a.Match

	return nil
}

func (e *EmbedThumbnail) copyOverTo(other interface{}) error {
	var dest *EmbedThumbnail
	var valid bool
	if dest, valid = other.(*EmbedThumbnail); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedThumbnail type")
	}
	dest.URL = e.URL
	dest.ProxyURL = e.ProxyURL
	dest.Height = e.Height
	dest.Width = e.Width

	return nil
}

func (a *AuditLogEntry) copyOverTo(other interface{}) error {
	var dest *AuditLogEntry
	var valid bool
	if dest, valid = other.(*AuditLogEntry); !valid {
		return newErrorUnsupportedType("argument given is not a *AuditLogEntry type")
	}
	dest.TargetID = a.TargetID
	dest.Changes = make([]*AuditLogChanges, len(a.Changes))
	for i := 0; i < len(a.Changes); i++ {
		dest.Changes[i] = DeepCopy(a.Changes[i]).(*AuditLogChanges)
	}
	dest.UserID = a.UserID
	dest.ID = a.ID
	dest.Event = a.Event
	dest.Options = a.Options
	dest.Reason = a.Reason

	return nil
}

func (m *MentionChannel) copyOverTo(other interface{}) error {
	var dest *MentionChannel
	var valid bool
	if dest, valid = other.(*MentionChannel); !valid {
		return newErrorUnsupportedType("argument given is not a *MentionChannel type")
	}
	dest.ID = m.ID
	dest.GuildID = m.GuildID
	dest.Type = m.Type
	dest.Name = m.Name

	return nil
}

func (a *ActivityAssets) copyOverTo(other interface{}) error {
	var dest *ActivityAssets
	var valid bool
	if dest, valid = other.(*ActivityAssets); !valid {
		return newErrorUnsupportedType("argument given is not a *ActivityAssets type")
	}
	dest.LargeImage = a.LargeImage
	dest.LargeText = a.LargeText
	dest.SmallImage = a.SmallImage
	dest.SmallText = a.SmallText

	return nil
}

func (m *MessageSticker) copyOverTo(other interface{}) error {
	var dest *MessageSticker
	var valid bool
	if dest, valid = other.(*MessageSticker); !valid {
		return newErrorUnsupportedType("argument given is not a *MessageSticker type")
	}
	dest.ID = m.ID
	dest.PackID = m.PackID
	dest.Name = m.Name
	dest.Description = m.Description
	dest.Tags = m.Tags
	dest.Asset = m.Asset
	dest.PreviewAsset = m.PreviewAsset
	dest.FormatType = m.FormatType

	return nil
}

func (c *Channel) copyOverTo(other interface{}) error {
	var dest *Channel
	var valid bool
	if dest, valid = other.(*Channel); !valid {
		return newErrorUnsupportedType("argument given is not a *Channel type")
	}
	dest.ID = c.ID
	dest.Type = c.Type
	dest.GuildID = c.GuildID
	dest.Position = c.Position
	dest.PermissionOverwrites = make([]PermissionOverwrite, len(c.PermissionOverwrites))
	copy(dest.PermissionOverwrites, c.PermissionOverwrites)
	dest.Name = c.Name
	dest.Topic = c.Topic
	dest.NSFW = c.NSFW
	dest.LastMessageID = c.LastMessageID
	dest.Bitrate = c.Bitrate
	dest.UserLimit = c.UserLimit
	dest.RateLimitPerUser = c.RateLimitPerUser
	dest.Recipients = make([]*User, len(c.Recipients))
	for i := 0; i < len(c.Recipients); i++ {
		dest.Recipients[i] = DeepCopy(c.Recipients[i]).(*User)
	}
	dest.Icon = c.Icon
	dest.OwnerID = c.OwnerID
	dest.ApplicationID = c.ApplicationID
	dest.ParentID = c.ParentID
	dest.LastPinTimestamp = c.LastPinTimestamp

	return nil
}

func (m *Message) copyOverTo(other interface{}) error {
	var dest *Message
	var valid bool
	if dest, valid = other.(*Message); !valid {
		return newErrorUnsupportedType("argument given is not a *Message type")
	}
	dest.ID = m.ID
	dest.ChannelID = m.ChannelID
	dest.GuildID = m.GuildID
	dest.Author = m.Author
	dest.Member = m.Member
	dest.Content = m.Content
	dest.Timestamp = m.Timestamp
	dest.EditedTimestamp = m.EditedTimestamp
	dest.Tts = m.Tts
	dest.MentionEveryone = m.MentionEveryone
	dest.Mentions = make([]*User, len(m.Mentions))
	for i := 0; i < len(m.Mentions); i++ {
		dest.Mentions[i] = DeepCopy(m.Mentions[i]).(*User)
	}
	dest.MentionRoles = make([]Snowflake, len(m.MentionRoles))
	copy(dest.MentionRoles, m.MentionRoles)
	dest.MentionChannels = make([]*MentionChannel, len(m.MentionChannels))
	for i := 0; i < len(m.MentionChannels); i++ {
		dest.MentionChannels[i] = DeepCopy(m.MentionChannels[i]).(*MentionChannel)
	}
	dest.Attachments = make([]*Attachment, len(m.Attachments))
	for i := 0; i < len(m.Attachments); i++ {
		dest.Attachments[i] = DeepCopy(m.Attachments[i]).(*Attachment)
	}
	dest.Embeds = make([]*Embed, len(m.Embeds))
	for i := 0; i < len(m.Embeds); i++ {
		dest.Embeds[i] = DeepCopy(m.Embeds[i]).(*Embed)
	}
	dest.Reactions = make([]*Reaction, len(m.Reactions))
	for i := 0; i < len(m.Reactions); i++ {
		dest.Reactions[i] = DeepCopy(m.Reactions[i]).(*Reaction)
	}
	dest.Nonce = m.Nonce
	dest.Pinned = m.Pinned
	dest.WebhookID = m.WebhookID
	dest.Type = m.Type
	dest.Activity = m.Activity
	dest.Application = m.Application
	dest.MessageReference = m.MessageReference
	dest.ReferencedMessage = m.ReferencedMessage
	dest.Flags = m.Flags
	dest.Stickers = make([]*MessageSticker, len(m.Stickers))
	for i := 0; i < len(m.Stickers); i++ {
		dest.Stickers[i] = DeepCopy(m.Stickers[i]).(*MessageSticker)
	}
	dest.Components = make([]*MessageComponent, len(m.Components))
	for i := 0; i < len(m.Components); i++ {
		dest.Components[i] = DeepCopy(m.Components[i]).(*MessageComponent)
	}
	dest.SpoilerTagContent = m.SpoilerTagContent
	dest.SpoilerTagAllAttachments = m.SpoilerTagAllAttachments
	dest.HasSpoilerImage = m.HasSpoilerImage

	return nil
}

func (e *Embed) copyOverTo(other interface{}) error {
	var dest *Embed
	var valid bool
	if dest, valid = other.(*Embed); !valid {
		return newErrorUnsupportedType("argument given is not a *Embed type")
	}
	dest.Title = e.Title
	dest.Type = e.Type
	dest.Description = e.Description
	dest.URL = e.URL
	dest.Timestamp = e.Timestamp
	dest.Color = e.Color
	dest.Footer = e.Footer
	dest.Image = e.Image
	dest.Thumbnail = e.Thumbnail
	dest.Video = e.Video
	dest.Provider = e.Provider
	dest.Author = e.Author
	dest.Fields = make([]*EmbedField, len(e.Fields))
	for i := 0; i < len(e.Fields); i++ {
		dest.Fields[i] = DeepCopy(e.Fields[i]).(*EmbedField)
	}

	return nil
}

func (e *EmbedProvider) copyOverTo(other interface{}) error {
	var dest *EmbedProvider
	var valid bool
	if dest, valid = other.(*EmbedProvider); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedProvider type")
	}
	dest.Name = e.Name
	dest.URL = e.URL

	return nil
}

func (i *InviteMetadata) copyOverTo(other interface{}) error {
	var dest *InviteMetadata
	var valid bool
	if dest, valid = other.(*InviteMetadata); !valid {
		return newErrorUnsupportedType("argument given is not a *InviteMetadata type")
	}
	dest.Inviter = i.Inviter
	dest.Uses = i.Uses
	dest.MaxUses = i.MaxUses
	dest.MaxAge = i.MaxAge
	dest.Temporary = i.Temporary
	dest.CreatedAt = i.CreatedAt
	dest.Revoked = i.Revoked

	return nil
}

func (e *Emoji) copyOverTo(other interface{}) error {
	var dest *Emoji
	var valid bool
	if dest, valid = other.(*Emoji); !valid {
		return newErrorUnsupportedType("argument given is not a *Emoji type")
	}
	dest.ID = e.ID
	dest.Name = e.Name
	dest.Roles = make([]Snowflake, len(e.Roles))
	copy(dest.Roles, e.Roles)
	dest.User = e.User
	dest.RequireColons = e.RequireColons
	dest.Managed = e.Managed
	dest.Animated = e.Animated

	return nil
}

func (i *Invite) copyOverTo(other interface{}) error {
	var dest *Invite
	var valid bool
	if dest, valid = other.(*Invite); !valid {
		return newErrorUnsupportedType("argument given is not a *Invite type")
	}
	dest.Code = i.Code
	dest.Guild = i.Guild
	dest.Channel = i.Channel
	dest.Inviter = i.Inviter
	dest.CreatedAt = i.CreatedAt
	dest.MaxAge = i.MaxAge
	dest.MaxUses = i.MaxUses
	dest.Temporary = i.Temporary
	dest.Uses = i.Uses
	dest.Revoked = i.Revoked
	dest.Unique = i.Unique
	dest.ApproximatePresenceCount = i.ApproximatePresenceCount
	dest.ApproximateMemberCount = i.ApproximateMemberCount

	return nil
}

func (e *EmbedImage) copyOverTo(other interface{}) error {
	var dest *EmbedImage
	var valid bool
	if dest, valid = other.(*EmbedImage); !valid {
		return newErrorUnsupportedType("argument given is not a *EmbedImage type")
	}
	dest.URL = e.URL
	dest.ProxyURL = e.ProxyURL
	dest.Height = e.Height
	dest.Width = e.Width

	return nil
}

func (a *ActivityTimestamp) copyOverTo(other interface{}) error {
	var dest *ActivityTimestamp
	var valid bool
	if dest, valid = other.(*ActivityTimestamp); !valid {
		return newErrorUnsupportedType("argument given is not a *ActivityTimestamp type")
	}
	dest.Start = a.Start
	dest.End = a.End

	return nil
}
