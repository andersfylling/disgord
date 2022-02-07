//go:build integration
// +build integration

package disgord

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetScheduledEvents(t *testing.T) {
	client := New(Config{BotToken: token})
	evts, err := client.GuildScheduledEvent(guildAdmin.ID).Gets(&GetScheduledEvents{
		WithUserCount: true,
	})

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(evts), 0)
}

func TestGetScheduledEvent(t *testing.T) {
	client := New(Config{BotToken: token})
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

	evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(cEvt)
	assert.Nil(t, err)

	gEvt, err := client.GuildScheduledEvent(guildAdmin.ID).Get(evt.ID, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, gEvt)
	assert.Equal(t, gEvt.Name, cEvt.Name)
	assert.Equal(t, gEvt.Description, cEvt.Description)
	assert.Equal(t, GuildScheduledEventPrivacyLevel(gEvt.PrivacyLevel), cEvt.PrivacyLevel)
}

func TestGetScheduledEventUsers(t *testing.T) {
	client := New(Config{BotToken: token})
	params := &GetScheduledEventMembers{
		Limit:      2,
		WithMember: false,
	}

	gEvtUsr, err := client.GuildScheduledEvent(guildAdmin.ID).GetMembers(935710181805936730, params)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(gEvtUsr), 0)
}

func TestDeleteScheduledEvent(t *testing.T) {
	client := New(Config{BotToken: token})
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

	evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(cEvt)
	assert.Nil(t, err)
	assert.NotEmpty(t, evt)
	assert.Equal(t, evt.Name, cEvt.Name)
	assert.Equal(t, evt.Description, cEvt.Description)
	assert.Equal(t, GuildScheduledEventPrivacyLevel(evt.PrivacyLevel), cEvt.PrivacyLevel)

	err = client.GuildScheduledEvent(guildAdmin.ID).Delete(evt.ID)
	assert.Nil(t, err)
}

func TestCreate(t *testing.T) {
	client := New(Config{BotToken: token})
	cTableTest := []struct {
		name    string
		evt     *CreateScheduledEvent
		wantErr error
	}{
		{
			name: "Create event with empty entity type",
			evt: &CreateScheduledEvent{
				Name: "Test event",
			},
			wantErr: ErrMissingScheduledEventEntityType,
		},
		{
			name: "Create event with empty channel ID and entity type is stage instance",
			evt: &CreateScheduledEvent{
				Name:       "Test event",
				EntityType: GuildScheduledEventEntityTypesStageInstance,
			},
			wantErr: ErrMissingChannelID,
		},
		{
			name: "Create event with empty channel ID and entity type is voice",
			evt: &CreateScheduledEvent{
				Name:       "Test event",
				EntityType: GuildScheduledEventEntityTypesVoice,
			},
			wantErr: ErrMissingChannelID,
		},
		{
			name: "Create event with empty location and entity type is external",
			evt: &CreateScheduledEvent{
				Name:       "Test event",
				EntityType: GuildScheduledEventEntityTypesExternal,
			},
			wantErr: ErrMissingScheduledEventLocation,
		},
		{
			name: "Create event with empty event name",
			evt: &CreateScheduledEvent{
				EntityType: GuildScheduledEventEntityTypesExternal,
				EntityMetadata: ScheduledEventEntityMetadata{
					Location: "Malang, Indonesia",
				},
			},
			wantErr: ErrMissingScheduledEventName,
		},
		{
			name: "Create event with less than minimum length of name",
			evt: &CreateScheduledEvent{
				Name:       "M",
				EntityType: GuildScheduledEventEntityTypesExternal,
				EntityMetadata: ScheduledEventEntityMetadata{
					Location: "Malang, Indonesia",
				},
			},
			wantErr: nil,
		},
		{
			name: "Create event with greater than max length of name",
			evt: &CreateScheduledEvent{
				Name:       strings.Repeat("AAA", 1000),
				EntityType: GuildScheduledEventEntityTypesExternal,
				EntityMetadata: ScheduledEventEntityMetadata{
					Location: "Malang, Indonesia",
				},
			},
			wantErr: nil,
		},
		{
			name: "Create event with privacy level is not guild",
			evt: &CreateScheduledEvent{
				Name:       "Name",
				EntityType: GuildScheduledEventEntityTypesExternal,
				EntityMetadata: ScheduledEventEntityMetadata{
					Location: "Malang, Indonesia",
				},
				PrivacyLevel: 0,
			},
			wantErr: ErrIllegalScheduledEventPrivacyLevelValue,
		},
		{
			name: "Create event with empty start time",
			evt: &CreateScheduledEvent{
				Name:       "Name",
				EntityType: GuildScheduledEventEntityTypesExternal,
				EntityMetadata: ScheduledEventEntityMetadata{
					Location: "Malang, Indonesia",
				},
				PrivacyLevel: GuildScheduledEventPrivacyLevelGuildOnly,
			},
			wantErr: ErrMissingScheduledEventStartTime,
		},
	}

	for _, v := range cTableTest {
		t.Run(v.name, func(t *testing.T) {
			evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(v.evt)

			if v.wantErr != nil {
				assert.Equal(t, v.wantErr, err)
			}

			assert.NotNil(t, err)
			assert.Nil(t, evt)
		})
	}

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

		evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(cEvt)
		assert.Nil(t, err)
		assert.NotNil(t, evt)
		assert.Equal(t, cEvt.Name, evt.Name)
		assert.Equal(t, cEvt.Description, evt.Description)

		err = client.GuildScheduledEvent(guildAdmin.ID).Delete(evt.ID)
		assert.Nil(t, err)
	})
}

func TestUpdateScheduledEvent(t *testing.T) {
	client := New(Config{BotToken: token})
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

	etVoice := GuildScheduledEventEntityTypesVoice
	etExt := GuildScheduledEventEntityTypesExternal

	plZero := GuildScheduledEventPrivacyLevel(0)

	minName := "M"
	maxName := strings.Repeat("AAA", 1000)

	uTableTest := []struct {
		name string
		evt  *UpdateScheduledEvent
	}{
		{
			name: "Update event with empty struct",
			evt:  nil,
		},
		{
			name: "Update event with minimum length of event name",
			evt: &UpdateScheduledEvent{
				EntityType: &etVoice,
				Name:       &minName,
			},
		},
		{
			name: "Update event with maximum length of event name",
			evt: &UpdateScheduledEvent{
				EntityType: &etVoice,
				Name:       &maxName,
			},
		},
		{
			name: "Update event with non allowed permission guild",
			evt: &UpdateScheduledEvent{
				EntityType:   &etExt,
				Name:         &maxName,
				PrivacyLevel: &plZero,
			},
		},
	}

	for _, v := range uTableTest {
		t.Run(v.name, func(t *testing.T) {
			evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(cEvt)
			assert.Nil(t, err)

			gEvt, err := client.GuildScheduledEvent(guildAdmin.ID).Update(evt.ID, v.evt)
			assert.NotNil(t, err)
			assert.Nil(t, gEvt)
		})
	}

	t.Run("Update event with valid data", func(t *testing.T) {
		evt, err := client.GuildScheduledEvent(guildAdmin.ID).Create(cEvt)
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

		gEvt, err := client.GuildScheduledEvent(guildAdmin.ID).Update(evt.ID, uEvt)
		assert.Nil(t, err)
		assert.NotNil(t, gEvt)
	})
}

// Run this test only if you want to delete all scheduled events test
func TestCleanUp(t *testing.T) {
	client := New(Config{BotToken: token})
	evts, err := client.GuildScheduledEvent(guildAdmin.ID).Gets(&GetScheduledEvents{
		WithUserCount: true,
	})

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(evts), 0)

	for i, v := range evts {
		err = client.GuildScheduledEvent(guildAdmin.ID).Delete(v.ID)
		assert.Nil(t, err)

		if (i%5 == 0) && (i != 0) {
			time.Sleep(time.Second * 2)
		}
	}
}
