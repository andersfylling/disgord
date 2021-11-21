// Code generated by generate/interfaces; DO NOT EDIT.

package disgord

func (a *Activity) reset() {
	a.ApplicationID = 0
	a.Assets = nil
	a.CreatedAt = 0
	a.Details = ""
	a.Emoji = nil
	a.Flags = 0
	a.Instance = false
	a.Name = ""
	a.Party = nil
	a.Secrets = nil
	a.State = ""
	a.Timestamps = nil
	a.Type = 0
	a.URL = ""
}

func (c *Channel) reset() {
	c.ApplicationID = 0
	c.Bitrate = 0
	c.DefaultAutoArchiveDuration = 0
	c.GuildID = 0
	c.Icon = ""
	c.ID = 0
	c.LastMessageID = 0
	c.LastPinTimestamp = Time{}
	c.Member = ThreadMember{}
	c.MemberCount = 0
	c.MessageCount = 0
	c.Name = ""
	c.NSFW = false
	c.OwnerID = 0
	c.ParentID = 0
	c.PermissionOverwrites = nil
	c.Position = 0
	c.RateLimitPerUser = 0
	c.Recipients = nil
	c.ThreadMetadata = ThreadMetadata{}
	c.Topic = ""
	c.Type = 0
	c.UserLimit = 0
}

func (e *Emoji) reset() {
	e.Animated = false
	e.Available = false
	e.ID = 0
	e.Managed = false
	e.Name = ""
	e.RequireColons = false
	e.Roles = nil
	if e.User != nil {
		Reset(e.User)
	}
}

func (g *Guild) reset() {
	g.AfkChannelID = 0
	g.AfkTimeout = 0
	g.ApplicationID = 0
	g.Banner = ""
	g.Channels = nil
	g.DefaultMessageNotifications = 0
	g.Description = ""
	g.DiscoverySplash = ""
	g.Emojis = nil
	g.ExplicitContentFilter = 0
	g.Features = nil
	g.Icon = ""
	g.ID = 0
	g.JoinedAt = nil
	g.Large = false
	g.MemberCount = 0
	g.Members = nil
	g.MFALevel = 0
	g.Name = ""
	g.Owner = false
	g.OwnerID = 0
	g.Permissions = 0
	g.PremiumSubscriptionCount = 0
	g.PremiumTier = 0
	g.Presences = nil
	g.Region = ""
	g.Roles = nil
	g.Splash = ""
	g.SystemChannelID = 0
	g.Unavailable = false
	g.VanityUrl = ""
	g.VerificationLevel = 0
	g.VoiceStates = nil
	g.WidgetChannelID = 0
	g.WidgetEnabled = false
}

func (m *Member) reset() {
	m.Deaf = false
	m.GuildID = 0
	m.JoinedAt = Time{}
	m.Mute = false
	m.Nick = ""
	m.Pending = false
	m.PremiumSince = Time{}
	m.Roles = nil
	if m.User != nil {
		Reset(m.User)
	}
	m.UserID = 0
}

func (m *Message) reset() {
	m.Activity = MessageActivity{}
	m.Application = MessageApplication{}
	m.Attachments = nil
	if m.Author != nil {
		Reset(m.Author)
	}
	m.ChannelID = 0
	m.Components = nil
	m.Content = ""
	m.EditedTimestamp = Time{}
	m.Embeds = nil
	m.Flags = 0
	m.GuildID = 0
	m.HasSpoilerImage = false
	m.ID = 0
	m.Interaction = nil
	if m.Member != nil {
		Reset(m.Member)
	}
	m.MentionChannels = nil
	m.MentionEveryone = false
	m.MentionRoles = nil
	m.Mentions = nil
	m.MessageReference = nil
	m.Nonce = nil
	m.Pinned = false
	m.Reactions = nil
	if m.ReferencedMessage != nil {
		Reset(m.ReferencedMessage)
	}
	m.SpoilerTagAllAttachments = false
	m.SpoilerTagContent = false
	m.StickerItems = nil
	m.Timestamp = Time{}
	m.Tts = false
	m.Type = 0
	m.WebhookID = 0
}

func (m *MessageCreate) reset() {
	if m.Message != nil {
		Reset(m.Message)
	}
	m.ShardID = 0
}

func (r *Reaction) reset() {
	r.Count = 0
	if r.Emoji != nil {
		Reset(r.Emoji)
	}
	r.Me = false
}

func (r *Role) reset() {
	r.Color = 0
	r.guildID = 0
	r.Hoist = false
	r.ID = 0
	r.Managed = false
	r.Mentionable = false
	r.Name = ""
	r.Permissions = 0
	r.Position = 0
}

func (u *User) reset() {
	u.Avatar = ""
	u.Bot = false
	u.Discriminator = 0
	u.Email = ""
	u.Flags = 0
	u.ID = 0
	u.Locale = ""
	u.MFAEnabled = false
	if u.PartialMember != nil {
		Reset(u.PartialMember)
	}
	u.PremiumType = 0
	u.PublicFlags = 0
	u.System = false
	u.Username = ""
	u.Verified = false
}

func (v *VoiceRegion) reset() {
	v.Custom = false
	v.Deprecated = false
	v.ID = ""
	v.Name = ""
	v.Optimal = false
	v.SampleHostname = ""
	v.SamplePort = 0
	v.VIP = false
}

func (v *VoiceState) reset() {
	v.ChannelID = 0
	v.Deaf = false
	v.GuildID = 0
	if v.Member != nil {
		Reset(v.Member)
	}
	v.Mute = false
	v.SelfDeaf = false
	v.SelfMute = false
	v.SessionID = ""
	v.Suppress = false
	v.UserID = 0
}
