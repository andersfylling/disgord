package disgord

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andersfylling/disgord/internal/endpoint"
	"github.com/andersfylling/disgord/internal/httd"
)

type GuildScheduledEvent struct {
	ID                 Snowflake   `json:"id"`
	GuildID            Snowflake   `json:"guild_id"`
	ChannelID          Snowflake   `json:"channel_id"`
	CreatorID          Snowflake   `json:"creator_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	ScheduledStartTime Time        `json:"scheduled_start_time"`
	ScheduledEndTime   Time        `json:"scheduled_end_time"`
	PrivacyLevel       int         `json:"privacy_level"`
	EventStatus        int         `json:"event_status"`
	EntityType         int         `json:"entity_type"`
	EntityMetadata     interface{} `json:"entity_metadata"`
	Creator            *User       `json:"creator"`
	UserCount          int         `json:"user_count"`
}

type GuildScheduledEventQueryBuilder interface {
	WithContext(ctx context.Context) GuildScheduledEventQueryBuilder
	WithFlags(flags ...Flag) GuildScheduledEventQueryBuilder

	Get(params *GetScheduledEvent) (*GuildScheduledEvent, error)
	Update(params *UpdateScheduledEvent) (*GuildScheduledEvent, error)
	Delete() error

	GetMembers(params *GetScheduledEventMembers) ([]*GuildScheduledEventUsers, error)
}

type guildScheduledEventQueryBuilder struct {
	ctx     context.Context
	flags   Flag
	client  *Client
	gid     Snowflake
	eventID Snowflake
}

func (gse guildScheduledEventQueryBuilder) WithContext(ctx context.Context) GuildScheduledEventQueryBuilder {
	gse.ctx = ctx
	return &gse
}

func (gse guildScheduledEventQueryBuilder) WithFlags(flags ...Flag) GuildScheduledEventQueryBuilder {
	gse.flags = mergeFlags(flags)
	return &gse
}

type GetScheduledEvents struct {
	WithUserCount bool `urlparam:"with_user_count,omitempty"`
}

var _ URLQueryStringer = (*GetScheduledEvents)(nil)

func (gse *GetScheduledEvents) FindErrors() error {
	return nil
}

func (gse guildScheduledEventQueryBuilder) Get(params *GetScheduledEvent) (*GuildScheduledEvent, error) {
	// TODO: add cache implementation
	if params == nil {
		params = &GetScheduledEvent{
			WithUserCount: false,
		}
	}

	r := gse.client.newRESTRequest(&httd.Request{
		Endpoint:    endpoint.ScheduledEvent(gse.gid, gse.eventID) + params.URLQueryString(),
		Ctx:         gse.ctx,
		ContentType: httd.ContentType,
	}, gse.flags)
	r.factory = func() interface{} {
		return &GuildScheduledEvent{}
	}

	return getScheduledEvent(r.Execute)
}

type GetScheduledEvent struct {
	WithUserCount bool `urlparam:"with_user_count,omitempty"`
}

var _ URLQueryStringer = (*GetScheduledEvent)(nil)

func (gse *GetScheduledEvent) FindErrors() error {
	return nil
}

func (gse guildScheduledEventQueryBuilder) Delete() error {
	r := gse.client.newRESTRequest(&httd.Request{
		Method:   http.MethodDelete,
		Endpoint: endpoint.ScheduledEvent(gse.gid, gse.eventID),
		Ctx:      gse.ctx,
	}, gse.flags)

	_, err := r.Execute()
	return err
}

// ScheduledEventEntityMetadata ...
// https://discord.com/developers/docs/resources/guild-scheduled-event#guild-scheduled-event-object-guild-scheduled-event-entity-metadata
type ScheduledEventEntityMetadata struct {
	Location string `json:"location,omitempty"` // required if EntityType is EXTERNAL
}

// CreateScheduledEvent ...
// https://discord.com/developers/docs/resources/guild-scheduled-event#create-guild-scheduled-event-json-params
type CreateScheduledEvent struct {
	ChannelID          Snowflake                       `json:"channel_id"` // optional if EntityType is EXTERNAL
	EntityMetadata     ScheduledEventEntityMetadata    `json:"entity_metadata"`
	Name               string                          `json:"name,omitempty"`
	PrivacyLevel       GuildScheduledEventPrivacyLevel `json:"privacy_level"`
	ScheduledStartTime Time                            `json:"scheduled_start_time"`
	ScheduledEndTime   Time                            `json:"scheduled_end_time,omitempty"`
	Description        string                          `json:"description,omitempty"`
	EntityType         GuildScheduledEventEntityTypes  `json:"entity_type"`

	AuditLogReason string `json:"-"`
}

func (cse *CreateScheduledEvent) validate() error {
	if cse.EntityType == 0 {
		return ErrMissingScheduledEventEntityType
	}
	if cse.EntityType != GuildScheduledEventEntityTypesExternal {
		if cse.ChannelID == 0 {
			return ErrMissingChannelID
		}
	}
	if cse.EntityType == GuildScheduledEventEntityTypesExternal {
		if cse.EntityMetadata.Location == "" {
			return ErrMissingScheduledEventLocation
		}
	}

	if cse.Name == "" {
		return ErrMissingScheduledEventName
	}
	if l := len(cse.Name); !(2 <= l && l <= 100) {
		return fmt.Errorf("scheduled event name must be 2 or more characters and no more than 100 characters: %w", ErrIllegalValue)
	}

	if cse.PrivacyLevel != GuildScheduledEventPrivacyLevelGuildOnly {
		return ErrIllegalScheduledEventPrivacyLevelValue
	}

	if cse.ScheduledStartTime.IsZero() {
		return ErrMissingScheduledEventStartTime
	}

	return nil
}

type UpdateScheduledEvent struct {
	ChannelID          *Snowflake                       `json:"channel_id"` // optional if EntityType is EXTERNAL
	EntityMetadata     *ScheduledEventEntityMetadata    `json:"entity_metadata,omitempty"`
	Name               *string                          `json:"name,omitempty"`
	PrivacyLevel       *GuildScheduledEventPrivacyLevel `json:"privacy_level,omitempty"`
	ScheduledStartTime *Time                            `json:"scheduled_start_time,omitempty"`
	ScheduledEndTime   *Time                            `json:"scheduled_end_time,omitempty"`
	Description        *string                          `json:"description,omitempty"`
	EntityType         *GuildScheduledEventEntityTypes  `json:"entity_type,omitempty"`
	Status             *GuildScheduledEventStatus       `json:"status,omitempty"`

	AuditLogReason string `json:"-"`
}

func (cse UpdateScheduledEvent) validate() error {
	if cse.EntityType != nil && *cse.EntityType == GuildScheduledEventEntityTypesExternal {
		if cse.EntityMetadata != nil && cse.EntityMetadata.Location == "" {
			return ErrMissingScheduledEventLocation
		}

		if cse.ScheduledEndTime == nil {
			return ErrMissingScheduledEventEndTime
		}

		cse.ChannelID = nil
	}

	if cse.Name != nil {
		if l := len(*cse.Name); !(2 <= l && l <= 100) {
			return fmt.Errorf("scheduled event name must be 2 or more characters and no more than 100 characters: %w", ErrIllegalValue)
		}
	}

	if cse.PrivacyLevel != nil && *cse.PrivacyLevel != GuildScheduledEventPrivacyLevelGuildOnly {
		return ErrIllegalScheduledEventPrivacyLevelValue
	}

	return nil
}

func (gse guildScheduledEventQueryBuilder) Update(params *UpdateScheduledEvent) (*GuildScheduledEvent, error) {
	if params == nil {
		return nil, ErrMissingRESTParams
	}

	if err := params.validate(); err != nil {
		return nil, err
	}

	r := gse.client.newRESTRequest(&httd.Request{
		Method:      http.MethodPatch,
		Ctx:         gse.ctx,
		Endpoint:    endpoint.ScheduledEvent(gse.gid, gse.eventID),
		Body:        params,
		ContentType: httd.ContentTypeJSON,
		Reason:      params.AuditLogReason,
	}, gse.flags)
	r.factory = func() interface{} {
		return &GuildScheduledEvent{}
	}

	return getScheduledEvent(r.Execute)
}

type GetScheduledEventMembers struct {
	Limit      uint32    `urlparam:"limit,omitempty"`
	WithMember bool      `urlparam:"with_member,omitempty"`
	Before     Snowflake `urlparam:"before,omitempty"`
	After      Snowflake `urlparam:"after,omitempty"`
}

var _ URLQueryStringer = (*GetScheduledEventMembers)(nil)

func (gse *GetScheduledEventMembers) FindErrors() error {
	return nil
}

type GuildScheduledEventUsers struct {
	GuildScheduledEventID Snowflake `json:"guild_scheduled_event_id"`
	User                  `json:"user"`
	Member                `json:"member"`
}

func (gse guildScheduledEventQueryBuilder) GetMembers(params *GetScheduledEventMembers) ([]*GuildScheduledEventUsers, error) {
	const QueryLimit uint32 = 100

	if params == nil {
		params = &GetScheduledEventMembers{
			Limit: QueryLimit,
		}
	}

	if params.Limit == 0 || params.Limit > QueryLimit {
		params.Limit = QueryLimit
	}

	if params.Before != 0 && params.After != 0 {
		params.After = 0
	}

	r := gse.client.newRESTRequest(&httd.Request{
		Endpoint:    endpoint.ScheduledEventUsers(gse.gid, gse.eventID) + params.URLQueryString(),
		Ctx:         gse.ctx,
		ContentType: httd.ContentType,
	}, gse.flags)
	r.factory = func() interface{} {
		gseusr := make([]*GuildScheduledEventUsers, 0)
		return &gseusr
	}

	return getScheduledEventUsers(r.Execute)
}
