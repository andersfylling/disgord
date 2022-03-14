//go:build integration
// +build integration

package disgord

import (
	"strings"
	"testing"
	"time"
)

func TestGetScheduledEvents(t *testing.T) {
	client := New(Config{BotToken: token})
	evts, err := client.Guild(guildAdmin.ID).GetScheduledEvents(&GetScheduledEvents{
		WithUserCount: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(evts) == 0 {
		t.Errorf("Expected at least 1 scheduled event, got %d", len(evts))
	}
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

	evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(cEvt)
	if err != nil {
		t.Fatal(err)
	}

	gEvt, err := client.Guild(guildAdmin.ID).ScheduledEvent(evt.ID).Get(nil)
	if err != nil {
		t.Fatal(err)
	}

	if gEvt == nil {
		t.Fatal("Expected scheduled event, got nil")
	}
	if gEvt.Name != cEvt.Name {
		t.Errorf("Expected name %s, got %s", cEvt.Name, gEvt.Name)
	}
	if gEvt.Description != cEvt.Description {
		t.Errorf("Expected description %s, got %s", cEvt.Description, gEvt.Description)
	}
	if GuildScheduledEventPrivacyLevel(gEvt.PrivacyLevel) != cEvt.PrivacyLevel {
		t.Errorf("Expected privacy level %d, got %d", cEvt.PrivacyLevel, GuildScheduledEventPrivacyLevel(gEvt.PrivacyLevel))
	}
}

func TestGetScheduledEventUsers(t *testing.T) {
	client := New(Config{BotToken: token})
	params := &GetScheduledEventMembers{
		Limit:      2,
		WithMember: false,
	}

	gEvtUsr, err := client.Guild(guildAdmin.ID).ScheduledEvent(935710181805936730).GetMembers(params)
	if err != nil {
		t.Fatal(err)
	}

	if len(gEvtUsr) == 0 {
		t.Errorf("Expected at least 1 scheduled event user, got %d", len(gEvtUsr))
	}
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

	evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(cEvt)
	if err != nil {
		t.Fatal(err)
	}

	if evt == nil {
		t.Fatal("Expected scheduled event, got nil")
	}
	if evt.Name != cEvt.Name {
		t.Errorf("Expected name %s, got %s", cEvt.Name, evt.Name)
	}
	if evt.Description != cEvt.Description {
		t.Errorf("Expected description %s, got %s", cEvt.Description, evt.Description)
	}
	if GuildScheduledEventPrivacyLevel(evt.PrivacyLevel) != cEvt.PrivacyLevel {
		t.Errorf("Expected privacy level %d, got %d", cEvt.PrivacyLevel, GuildScheduledEventPrivacyLevel(evt.PrivacyLevel))
	}

	err = client.Guild(guildAdmin.ID).ScheduledEvent(evt.ID).Delete()
	if err != nil {
		t.Fatal(err)
	}
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
			evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(v.evt)

			if v.wantErr != nil && v.wantErr != err {
				t.Error(err)
			}

			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if evt != nil {
				t.Fatal("Expected nil, got event")
			}
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

		evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(cEvt)
		if err != nil {
			t.Fatal(err)
		}

		if evt == nil {
			t.Fatal("Expected event, got nil")
		}
		if evt.Name != cEvt.Name {
			t.Errorf("Expected event name %s, got %s", cEvt.Name, evt.Name)
		}
		if evt.Description != cEvt.Description {
			t.Errorf("Expected event description %s, got %s", cEvt.Description, evt.Description)
		}

		err = client.Guild(guildAdmin.ID).ScheduledEvent(evt.ID).Delete()
		if err != nil {
			t.Fatal(err)
		}
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
			evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(cEvt)
			if err != nil {
				t.Fatal(err)
			}

			gEvt, err := client.Guild(guildAdmin.ID).ScheduledEvent(evt.ID).Update(v.evt)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if gEvt != nil {
				t.Fatal("Expected nil, got event")
			}
		})
	}

	t.Run("Update event with valid data", func(t *testing.T) {
		evt, err := client.Guild(guildAdmin.ID).CreateScheduledEvent(cEvt)
		if err != nil {
			t.Fatal(err)
		}

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

		gEvt, err := client.Guild(guildAdmin.ID).ScheduledEvent(evt.ID).Update(uEvt)
		if err != nil {
			t.Fatal(err)
		}
		if gEvt == nil {
			t.Fatal("Expected event, got nil")
		}
	})
}
