package lvl

// DefaultMessageNotification ...
// https://discordapp.com/developers/docs/resources/guild#guild-object-default-message-notification-level
type DefaultMessageNotification uint

func (dmnl *DefaultMessageNotification) AllMessages() bool {
	return *dmnl == 0
}
func (dmnl *DefaultMessageNotification) OnlyMentions() bool {
	return *dmnl == 1
}
