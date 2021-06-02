// Code generated by generate/interfaces; DO NOT EDIT.

package disgord

func (m *Member) reset() {
	m.GuildID = 0
	if m.User != nil {
		Reset(m.User)
	}
	m.Nick = ""
	m.Roles = nil
	m.JoinedAt = Time{}
	m.PremiumSince = Time{}
	m.Deaf = false
	m.Mute = false
	m.Pending = false
	m.UserID = 0
}

func (g *Guild) reset() {
	g.ID = 0
	g.ApplicationID = 0
	g.Name = ""
	g.Icon = ""
	g.Splash = ""
	g.Owner = false
	g.OwnerID = 0
	g.Permissions = 0
	g.Region = ""
	g.AfkChannelID = 0
	g.AfkTimeout = 0
	g.VerificationLevel = 0
	g.DefaultMessageNotifications = 0
	g.ExplicitContentFilter = 0
	g.Roles = nil
	g.Emojis = nil
	g.Features = nil
	g.MFALevel = 0
	g.WidgetEnabled = false
	g.WidgetChannelID = 0
	g.SystemChannelID = 0
	g.JoinedAt = nil
	g.Large = false
	g.Unavailable = false
	g.MemberCount = 0
	g.VoiceStates = nil
	g.Members = nil
	g.Channels = nil
	g.Presences = nil
}

func (r *Reaction) reset() {
	r.Count = 0
	r.Me = false
	if r.Emoji != nil {
		Reset(r.Emoji)
	}
}

func (m *MessageCreate) reset() {
	if m.Message != nil {
		Reset(m.Message)
	}
	m.ShardID = 0
}

func (u *User) reset() {
	u.ID = 0
	u.Username = ""
	u.Discriminator = 0
	u.Avatar = ""
	u.Bot = false
	u.System = false
	u.MFAEnabled = false
	u.Locale = ""
	u.Verified = false
	u.Email = ""
	u.Flags = 0
	u.PremiumType = 0
	u.PublicFlags = 0
	if u.PartialMember != nil {
		Reset(u.PartialMember)
	}
}

func (a *Activity) reset() {
	a.Name = ""
	a.Type = 0
	a.URL = ""
	a.CreatedAt = 0
	a.Timestamps = nil
	a.ApplicationID = 0
	a.Details = ""
	a.State = ""
	a.Emoji = nil
	a.Party = nil
	a.Assets = nil
	a.Secrets = nil
	a.Instance = false
	a.Flags = 0
}

func (v *VoiceRegion) reset() {
	v.ID = ""
	v.Name = ""
	v.SampleHostname = ""
	v.SamplePort = 0
	v.VIP = false
	v.Optimal = false
	v.Deprecated = false
	v.Custom = false
}

func (r *Role) reset() {
	r.ID = 0
	r.Name = ""
	r.Color = 0
	r.Hoist = false
	r.Position = 0
	r.Permissions = 0
	r.Managed = false
	r.Mentionable = false
	r.guildID = 0
}

func (v *VoiceState) reset() {
	v.GuildID = 0
	v.ChannelID = 0
	v.UserID = 0
	if v.Member != nil {
		Reset(v.Member)
	}
	v.SessionID = ""
	v.Deaf = false
	v.Mute = false
	v.SelfDeaf = false
	v.SelfMute = false
	v.Suppress = false
}

func (c *Channel) reset() {
	c.ID = 0
	c.Type = 0
	c.GuildID = 0
	c.Position = 0
	c.PermissionOverwrites = nil
	c.Name = ""
	c.Topic = ""
	c.NSFW = false
	c.LastMessageID = 0
	c.Bitrate = 0
	c.UserLimit = 0
	c.RateLimitPerUser = 0
	c.Recipients = nil
	c.Icon = ""
	c.OwnerID = 0
	c.ApplicationID = 0
	c.ParentID = 0
	c.LastPinTimestamp = Time{}
}

func (m *Message) reset() {
	m.ID = 0
	m.ChannelID = 0
	m.GuildID = 0
	if m.Author != nil {
		Reset(m.Author)
	}
	if m.Member != nil {
		Reset(m.Member)
	}
	m.Content = ""
	m.Timestamp = Time{}
	m.EditedTimestamp = Time{}
	m.Tts = false
	m.MentionEveryone = false
	m.Mentions = nil
	m.MentionRoles = nil
	m.MentionChannels = nil
	m.Attachments = nil
	m.Embeds = nil
	m.Reactions = nil
	m.Nonce = nil
	m.Pinned = false
	m.WebhookID = 0
	m.Type = 0
	m.Activity = MessageActivity{}
	m.Application = MessageApplication{}
	m.MessageReference = nil
	if m.ReferencedMessage != nil {
		Reset(m.ReferencedMessage)
	}
	m.Flags = 0
	m.Stickers = nil
	m.Components = nil
	m.SpoilerTagContent = false
	m.SpoilerTagAllAttachments = false
	m.HasSpoilerImage = false
}

func (e *Emoji) reset() {
	e.ID = 0
	e.Name = ""
	e.Roles = nil
	if e.User != nil {
		Reset(e.User)
	}
	e.RequireColons = false
	e.Managed = false
	e.Animated = false
}
