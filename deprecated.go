package disgord

// Deprecated: use UpdateBuilder
func (m messageQueryBuilder) Update(flags ...Flag) (builder *updateMessageBuilder) {
	return m.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (g guildMemberQueryBuilder) Update(flags ...Flag) UpdateGuildMemberBuilder {
	return g.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (g guildQueryBuilder) Update(flags ...Flag) UpdateGuildBuilder { return g.UpdateBuilder(flags...) }

// Deprecated: use UpdateBuilder
func (g guildEmojiQueryBuilder) Update(flags ...Flag) UpdateGuildEmojiBuilder {
	return g.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (c channelQueryBuilder) Update(flags ...Flag) (builder *updateChannelBuilder) {
	return c.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (guildQueryBuilderNop) Update(flags ...Flag) UpdateGuildBuilder {
	return nil
}

// Deprecated: use UpdateBuilder
func (currentUserQueryBuilderNop) Update(_ ...Flag) UpdateCurrentUserBuilder {
	return nil
}

// Deprecated: use UpdateBuilder
func (g guildRoleQueryBuilder) Update(flags ...Flag) UpdateGuildRoleBuilder {
	return g.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (c currentUserQueryBuilder) Update(flags ...Flag) UpdateCurrentUserBuilder {
	return c.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (w webhookQueryBuilder) Update(flags ...Flag) (builder *updateWebhookBuilder) {
	return w.UpdateBuilder(flags...)
}

// Deprecated: use UpdateBuilder
func (w webhookWithTokenQueryBuilder) Update(flags ...Flag) (builder *updateWebhookBuilder) {
	return w.UpdateBuilder(flags...)
}
