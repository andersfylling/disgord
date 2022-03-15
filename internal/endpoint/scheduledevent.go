package endpoint

import "fmt"

// ScheduledEvents /guilds/{guild.id}/scheduled-events
func ScheduledEvents(id fmt.Stringer) string {
	return Guild(id) + scheduledEvents
}

// ScheduledEvent /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}
func ScheduledEvent(guildID, gseID fmt.Stringer) string {
	return Guild(guildID) + scheduledEvents + "/" + gseID.String()
}

// ScheduledEvent /guilds/{guild.id}/scheduled-events/{guild_scheduled_event.id}/users
func ScheduledEventUsers(guildID, gseID fmt.Stringer) string {
	return Guild(guildID) + scheduledEvents + "/" + gseID.String() + users
}
