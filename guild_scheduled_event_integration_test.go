package disgord

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var token string
var guildID Snowflake

var maxName = `Light it creeping a herb evening seas you're stars. 
	Lights let dry land let female is blessed also blessed
	they're life were rule that subdue, third may. Greater without, 
	given can't And bring she'd created fruitful third sea. 
	Good, dominion whose i blessed and second a appear replenish shall great, 
	void two sea which god. A place female abundantly, 
	seas fruitful moveth us heaven. Forth beginning to and image own seasons land had dry. 
	Given that they're the face without. Wherein he first also. 
	Fill hath. Sea. Have waters. Deep over earth grass fill had was it the a of.`

func init() {
	token = os.Getenv("DISCORD_BOT_TOKEN")
	gID, _ := strconv.Atoi(os.Getenv("DISCORD_GUILD_ID"))
	guildID = Snowflake(gID)
}

func TestGetScheduledEvents(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	evts, err := client.GuildScheduledEvent(guildID).Gets(&GetScheduledEvents{
		WithUserCount: true,
	})

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(evts), 0)
}

func TestGetScheduledEvent(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	cEvt := &CreateScheduledEvent{
		Name:       "Test Scheduled Event",
		EntityType: GuildScheduledEventEntityTypesExternal,
		EntityMetadata: ScheduledEventEntityMetadata{
			Location: "Malang, Indonesia",
		},
		PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
		ScheduledStartTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7),
		},
		ScheduledEndTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7 * 2),
		},
		Description:    "Test description",
		AuditLogReason: "integration test",
	}

	evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
	assert.Nil(t, err)

	gEvt, err := client.GuildScheduledEvent(guildID).Get(evt.ID, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, gEvt)
	assert.Equal(t, gEvt.Name, cEvt.Name)
	assert.Equal(t, gEvt.Description, cEvt.Description)
	assert.Equal(t, GuildScheduledEventPrivacyLevel(gEvt.PrivacyLevel), cEvt.PrivacyLevel)
}

func TestDeleteScheduledEvent(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	cEvt := &CreateScheduledEvent{
		Name:       "Test Scheduled Event",
		EntityType: GuildScheduledEventEntityTypesExternal,
		EntityMetadata: ScheduledEventEntityMetadata{
			Location: "Malang, Indonesia",
		},
		PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
		ScheduledStartTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7),
		},
		ScheduledEndTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7 * 2),
		},
		Description:    "Test description",
		AuditLogReason: "integration test",
	}

	evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
	assert.Nil(t, err)
	assert.NotEmpty(t, evt)
	assert.Equal(t, evt.Name, cEvt.Name)
	assert.Equal(t, evt.Description, cEvt.Description)
	assert.Equal(t, GuildScheduledEventPrivacyLevel(evt.PrivacyLevel), cEvt.PrivacyLevel)

	err = client.GuildScheduledEvent(guildID).Delete(evt.ID)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	t.Run("Create event with empty entity type", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name: "Test event",
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingScheduledEventEntityType, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with empty channel ID and entity type is stage instance", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Test event",
			EntityType: GuildScheduledEventEntityTypesStageInstance,
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingChannelID, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with empty channel ID and entity type is voice", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Test event",
			EntityType: GuildScheduledEventEntityTypesVoice,
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingChannelID, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with empty location and entity type is external", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Test event",
			EntityType: GuildScheduledEventEntityTypesExternal,
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingScheduledEventLocation, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with empty event name", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingScheduledEventName, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with less than minimum length of name", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "M",
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with greater than max length of name", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       maxName,
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with privacy level is not guild", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Name",
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
			PrivacyLevel: 0,
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrIllegalScheduledEventPrivacyLevelValue, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with empty start time", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Name",
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
			PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.NotNil(t, err)
		assert.Equal(t, ErrMissingScheduledEventStartTime, err)
		assert.Nil(t, evt)
	})

	t.Run("Create event with valid value", func(t *testing.T) {
		cEvt := &CreateScheduledEvent{
			Name:       "Test Scheduled Event",
			EntityType: GuildScheduledEventEntityTypesExternal,
			EntityMetadata: ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
			PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
			ScheduledStartTime: Time{
				Time: time.Now().Add(time.Hour * 24 * 7),
			},
			ScheduledEndTime: Time{
				Time: time.Now().Add(time.Hour * 24 * 7 * 2),
			},
			Description:    "Test description",
			AuditLogReason: "integration test",
		}

		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)
		assert.NotNil(t, evt)
		assert.Equal(t, cEvt.Name, evt.Name)
		assert.Equal(t, cEvt.Description, evt.Description)

		err = client.GuildScheduledEvent(guildID).Delete(evt.ID)
		assert.Nil(t, err)
	})
}

func TestUpdateScheduledEvent(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	cEvt := &CreateScheduledEvent{
		Name:       "Test Scheduled Update Event",
		EntityType: GuildScheduledEventEntityTypesExternal,
		EntityMetadata: ScheduledEventEntityMetadata{
			Location: "Malang, Indonesia",
		},
		PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
		ScheduledStartTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7),
		},
		ScheduledEndTime: Time{
			Time: time.Now().Add(time.Hour * 24 * 7 * 2),
		},
		Description:    "Test description",
		AuditLogReason: "integration test",
	}

	t.Run("Update event with empty struct", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)

		gEvt, err := client.GuildScheduledEvent(guildID).Update(evt.ID, nil)
		assert.NotNil(t, err)
		assert.Nil(t, gEvt)
	})

	t.Run("Update event with minimum length of event name", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)

		et := GuildScheduledEventEntityTypesVoice
		name := "M"
		uEvt := &UpdateScheduledEvent{
			EntityType: &et,
			Name:       &name,
		}

		gEvt, err := client.GuildScheduledEvent(guildID).Update(evt.ID, uEvt)
		assert.NotNil(t, err)
		assert.Nil(t, gEvt)
	})

	t.Run("Update event with maximum length of event name", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)

		et := GuildScheduledEventEntityTypesVoice
		name := maxName
		uEvt := &UpdateScheduledEvent{
			EntityType: &et,
			Name:       &name,
		}

		gEvt, err := client.GuildScheduledEvent(guildID).Update(evt.ID, uEvt)
		assert.NotNil(t, err)
		assert.Nil(t, gEvt)
	})

	t.Run("Update event with non allowed permission guild", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)

		et := GuildScheduledEventEntityTypesExternal
		name := maxName
		pl := GuildScheduledEventPrivacyLevel(0)
		uEvt := &UpdateScheduledEvent{
			EntityType:   &et,
			Name:         &name,
			PrivacyLevel: &pl,
		}

		gEvt, err := client.GuildScheduledEvent(guildID).Update(evt.ID, uEvt)
		assert.NotNil(t, err)
		assert.Nil(t, gEvt)
	})

	t.Run("Update event with valid data", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildID).Create(cEvt)
		assert.Nil(t, err)

		et := GuildScheduledEventEntityTypesExternal
		name := "Update Test Scheduled Event"
		pl := GuildScheduledEventPrivacyLevel(GuildScheduledEventPrivacyLevelGuildOnly)
		st := GuildScheduledEventStatusScheduled

		uEvt := &UpdateScheduledEvent{
			EntityType: &et,
			EntityMetadata: &ScheduledEventEntityMetadata{
				Location: "Malang, Indonesia",
			},
			ScheduledEndTime: &Time{
				Time: time.Now().Add(time.Hour * 24 * 7 * 3),
			},
			Name:         &name,
			PrivacyLevel: &pl,
			Status:       &st,
		}

		gEvt, err := client.GuildScheduledEvent(guildID).Update(evt.ID, uEvt)
		assert.Nil(t, err)
		assert.NotNil(t, gEvt)
	})

}

// Run this test only if you want to delete all scheduled events test
func TestCleanUp(t *testing.T) {
	client := New(Config{
		BotToken: token,
	})

	evts, err := client.GuildScheduledEvent(guildID).Gets(&GetScheduledEvents{
		WithUserCount: true,
	})

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(evts), 0)

	for i, v := range evts {
		err = client.GuildScheduledEvent(guildID).Delete(v.ID)
		assert.Nil(t, err)

		if (i%5 == 0) && (i != 0) {
			time.Sleep(time.Second * 2)
		}
	}
}
