package disgord

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/websocket"
)

// NewDispatch construct a Dispatch object for reacting to web socket events
// from discord
func NewDispatch(ws websocket.DiscordWebsocket) *Dispatch {
	dispatcher := &Dispatch{
		allChan:                      make(chan interface{}),
		readyChan:                    make(chan *Ready),
		resumedChan:                  make(chan *Resumed),
		channelCreateChan:            make(chan *ChannelCreate),
		channelUpdateChan:            make(chan *ChannelUpdate),
		channelDeleteChan:            make(chan *ChannelDelete),
		channelPinsUpdateChan:        make(chan *ChannelPinsUpdate),
		guildCreateChan:              make(chan *GuildCreate),
		guildUpdateChan:              make(chan *GuildUpdate),
		guildDeleteChan:              make(chan *GuildDelete),
		guildBanAddChan:              make(chan *GuildBanAdd),
		guildBanRemoveChan:           make(chan *GuildBanRemove),
		guildEmojisUpdateChan:        make(chan *GuildEmojisUpdate),
		guildIntegrationsUpdateChan:  make(chan *GuildIntegrationsUpdate),
		guildMemberAddChan:           make(chan *GuildMemberAdd),
		guildMemberRemoveChan:        make(chan *GuildMemberRemove),
		guildMemberUpdateChan:        make(chan *GuildMemberUpdate),
		guildMembersChunkChan:        make(chan *GuildMembersChunk),
		guildRoleUpdateChan:          make(chan *GuildRoleUpdate),
		guildRoleCreateChan:          make(chan *GuildRoleCreate),
		guildRoleDeleteChan:          make(chan *GuildRoleDelete),
		messageCreateChan:            make(chan *MessageCreate),
		messageUpdateChan:            make(chan *MessageUpdate),
		messageDeleteChan:            make(chan *MessageDelete),
		messageDeleteBulkChan:        make(chan *MessageDeleteBulk),
		messageReactionAddChan:       make(chan *MessageReactionAdd),
		messageReactionRemoveChan:    make(chan *MessageReactionRemove),
		messageReactionRemoveAllChan: make(chan *MessageReactionRemoveAll),
		presenceUpdateChan:           make(chan *PresenceUpdate),
		presencesReplaceChan:         make(chan *PresencesReplace),
		typingStartChan:              make(chan *TypingStart),
		userUpdateChan:               make(chan *UserUpdate),
		voiceStateUpdateChan:         make(chan *VoiceStateUpdate),
		voiceServerUpdateChan:        make(chan *VoiceServerUpdate),
		webhooksUpdateChan:           make(chan *WebhooksUpdate),

		ws: ws,

		listeners:      make(map[string][]interface{}),
		listenOnceOnly: make(map[string][]int),

		shutdown: make(chan struct{}),
	}

	return dispatcher
}

// Dispatch holds all the channels and internal state for all registered
// observers
type Dispatch struct {
	allChan                      chan interface{} // any event
	readyChan                    chan *Ready
	resumedChan                  chan *Resumed
	channelCreateChan            chan *ChannelCreate
	channelUpdateChan            chan *ChannelUpdate
	channelDeleteChan            chan *ChannelDelete
	channelPinsUpdateChan        chan *ChannelPinsUpdate
	guildCreateChan              chan *GuildCreate
	guildUpdateChan              chan *GuildUpdate
	guildDeleteChan              chan *GuildDelete
	guildBanAddChan              chan *GuildBanAdd
	guildBanRemoveChan           chan *GuildBanRemove
	guildEmojisUpdateChan        chan *GuildEmojisUpdate
	guildIntegrationsUpdateChan  chan *GuildIntegrationsUpdate
	guildMemberAddChan           chan *GuildMemberAdd
	guildMemberRemoveChan        chan *GuildMemberRemove
	guildMemberUpdateChan        chan *GuildMemberUpdate
	guildMembersChunkChan        chan *GuildMembersChunk
	guildRoleUpdateChan          chan *GuildRoleUpdate
	guildRoleCreateChan          chan *GuildRoleCreate
	guildRoleDeleteChan          chan *GuildRoleDelete
	messageCreateChan            chan *MessageCreate
	messageUpdateChan            chan *MessageUpdate
	messageDeleteChan            chan *MessageDelete
	messageDeleteBulkChan        chan *MessageDeleteBulk
	messageReactionAddChan       chan *MessageReactionAdd
	messageReactionRemoveChan    chan *MessageReactionRemove
	messageReactionRemoveAllChan chan *MessageReactionRemoveAll
	presenceUpdateChan           chan *PresenceUpdate
	presencesReplaceChan         chan *PresencesReplace
	typingStartChan              chan *TypingStart
	userUpdateChan               chan *UserUpdate
	voiceStateUpdateChan         chan *VoiceStateUpdate
	voiceServerUpdateChan        chan *VoiceServerUpdate
	webhooksUpdateChan           chan *WebhooksUpdate

	ws websocket.DiscordWebsocket

	listeners      map[string][]interface{}
	listenOnceOnly map[string][]int

	shutdown chan struct{}

	listenersLock sync.RWMutex
}

func (c *Dispatch) On(event string, handlers ...interface{}) {
	c.ws.RegisterEvent(event)

	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()
	for _, handler := range handlers {
		c.listeners[event] = append(c.listeners[event], handler)
	}
}

func (c *Dispatch) Once(event string, handlers ...interface{}) {
	c.ws.RegisterEvent(event) // TODO: remove event after firing. unless there are more handlers

	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()
	for _, handler := range handlers {
		index := len(c.listeners[event])
		c.listeners[event] = append(c.listeners[event], handler)
		c.listenOnceOnly[event] = append(c.listenOnceOnly[event], index)
	}
}

func (c *Dispatch) EventChan(event string) (channel interface{}, err error) {
	switch event {
	case EventReady:
		channel = c.Ready()
	case EventResumed:
		channel = c.Resumed()
	case EventChannelCreate:
		channel = c.ChannelCreate()
	case EventChannelUpdate:
		channel = c.ChannelUpdate()
	case EventChannelDelete:
		channel = c.ChannelDelete()
	case EventChannelPinsUpdate:
		channel = c.ChannelPinsUpdate()
	case EventGuildCreate:
		channel = c.GuildCreate()
	case EventGuildUpdate:
		channel = c.GuildUpdate()
	case EventGuildDelete:
		channel = c.GuildDelete()
	case EventGuildBanAdd:
		channel = c.GuildBanAdd()
	case EventGuildBanRemove:
		channel = c.GuildBanRemove()
	case EventGuildEmojisUpdate:
		channel = c.GuildEmojisUpdate()
	case EventGuildIntegrationsUpdate:
		channel = c.GuildIntegrationsUpdate()
	case EventGuildMemberAdd:
		channel = c.GuildMemberAdd()
	case EventGuildMemberRemove:
		channel = c.GuildMemberRemove()
	case EventGuildMemberUpdate:
		channel = c.GuildMemberUpdate()
	case EventGuildMembersChunk:
		channel = c.GuildMembersChunk()
	case EventGuildRoleCreate:
		channel = c.GuildRoleCreate()
	case EventGuildRoleUpdate:
		channel = c.GuildRoleUpdate()
	case EventGuildRoleDelete:
		channel = c.GuildRoleDelete()
	case EventMessageCreate:
		channel = c.MessageCreate()
	case EventMessageUpdate:
		channel = c.MessageUpdate()
	case EventMessageDelete:
		channel = c.MessageDelete()
	case EventMessageDeleteBulk:
		channel = c.MessageDeleteBulk()
	case EventMessageReactionAdd:
		channel = c.MessageReactionAdd()
	case EventMessageReactionRemove:
		channel = c.MessageReactionRemove()
	case EventMessageReactionRemoveAll:
		channel = c.MessageReactionRemoveAll()
	case EventPresenceUpdate:
		channel = c.PresenceUpdate()
	case EventTypingStart:
		channel = c.TypingStart()
	case EventUserUpdate:
		channel = c.UserUpdate()
	case EventVoiceStateUpdate:
		channel = c.VoiceStateUpdate()
	case EventVoiceServerUpdate:
		channel = c.VoiceServerUpdate()
	case EventWebhooksUpdate:
		channel = c.WebhooksUpdate()
	default:
		err = errors.New("no event channel exists for given event: " + event)
	}

	return
}

func (d *Dispatch) start() {
	// make sure every channel has a receiver to avoid deadlock
	// TODO: review, this feels hacky
	d.alwaysListenToChans()
}

func (d *Dispatch) stop() {
	close(d.shutdown)
}

// On places listeners into their respected stacks
// func (d *Dispatcher) OnEvent(evtName string, listener EventCallback) {
// 	d.listeners[evtName] = append(d.listeners[evtName], listener)
// }

// alwaysListenToChans makes sure no deadlocks occure
func (d *Dispatch) alwaysListenToChans() {
	go func() {
		stop := false
		for {
			select {
			case <-d.allChan:
			case <-d.readyChan:
			case <-d.resumedChan:
			case <-d.channelCreateChan:
			case <-d.channelDeleteChan:
			case <-d.channelPinsUpdateChan:
			case <-d.channelUpdateChan:
			case <-d.guildBanAddChan:
			case <-d.guildBanRemoveChan:
			case <-d.guildCreateChan:
			case <-d.guildDeleteChan:
			case <-d.guildEmojisUpdateChan:
			case <-d.guildIntegrationsUpdateChan:
			case <-d.guildMemberAddChan:
			case <-d.guildMemberRemoveChan:
			case <-d.guildMemberUpdateChan:
			case <-d.guildMembersChunkChan:
			case <-d.guildRoleCreateChan:
			case <-d.guildRoleDeleteChan:
			case <-d.guildRoleUpdateChan:
			case <-d.guildUpdateChan:
			case <-d.messageCreateChan:
			case <-d.messageDeleteBulkChan:
			case <-d.messageDeleteChan:
			case <-d.messageReactionAddChan:
			case <-d.messageReactionRemoveAllChan:
			case <-d.messageReactionRemoveChan:
			case <-d.messageUpdateChan:
			case <-d.presenceUpdateChan:
			case <-d.presencesReplaceChan:
			case <-d.typingStartChan:
			case <-d.userUpdateChan:
			case <-d.voiceStateUpdateChan:
			case <-d.voiceServerUpdateChan:
			case <-d.webhooksUpdateChan:
			case <-d.shutdown:
				stop = true
			}

			if stop {
				break
			}
		}
	}()
}

func (d *Dispatch) triggerChan(ctx context.Context, evtName string, session Session, box interface{}) {
	switch evtName {
	case EventReady:
		d.readyChan <- box.(*Ready)
	case EventResumed:
		d.resumedChan <- box.(*Resumed)
	case EventChannelCreate:
		d.channelCreateChan <- box.(*ChannelCreate)
	case EventChannelUpdate:
		d.channelUpdateChan <- box.(*ChannelUpdate)
	case EventChannelDelete:
		d.channelDeleteChan <- box.(*ChannelDelete)
	case EventChannelPinsUpdate:
		d.channelPinsUpdateChan <- box.(*ChannelPinsUpdate)
	case EventGuildCreate:
		for _, role := range (box.(*GuildCreate)).Guild.Roles {
			role.guildID = (box.(*GuildCreate)).Guild.ID
		}
		d.guildCreateChan <- box.(*GuildCreate)
	case EventGuildUpdate:
		for _, role := range (box.(*GuildCreate)).Guild.Roles {
			role.guildID = (box.(*GuildCreate)).Guild.ID
		}
		d.guildUpdateChan <- box.(*GuildUpdate)
	case EventGuildDelete:
		d.guildDeleteChan <- box.(*GuildDelete)
	case EventGuildBanAdd:
		d.guildBanAddChan <- box.(*GuildBanAdd)
	case EventGuildBanRemove:
		d.guildBanRemoveChan <- box.(*GuildBanRemove)
	case EventGuildEmojisUpdate:
		d.guildEmojisUpdateChan <- box.(*GuildEmojisUpdate)
	case EventGuildIntegrationsUpdate:
		d.guildIntegrationsUpdateChan <- box.(*GuildIntegrationsUpdate)
	case EventGuildMemberAdd:
		// Member.Roles is just a snowflake list
		d.guildMemberAddChan <- box.(*GuildMemberAdd)
	case EventGuildMemberRemove:
		d.guildMemberRemoveChan <- box.(*GuildMemberRemove)
	case EventGuildMemberUpdate:
		for _, role := range (box.(*GuildMemberUpdate)).Roles {
			role.guildID = (box.(*GuildMemberUpdate)).GuildID
		}
		d.guildMemberUpdateChan <- box.(*GuildMemberUpdate)
	case EventGuildMembersChunk:
		d.guildMembersChunkChan <- box.(*GuildMembersChunk)
	case EventGuildRoleCreate:
		(box.(*GuildRoleCreate)).Role.guildID = (box.(*GuildCreate)).Guild.ID
		d.guildRoleCreateChan <- box.(*GuildRoleCreate)
	case EventGuildRoleUpdate:
		(box.(*GuildRoleUpdate)).Role.guildID = (box.(*GuildRoleUpdate)).GuildID
		d.guildRoleUpdateChan <- box.(*GuildRoleUpdate)
	case EventGuildRoleDelete:
		d.guildRoleDeleteChan <- box.(*GuildRoleDelete)
	case EventMessageCreate:
		d.messageCreateChan <- box.(*MessageCreate)
	case EventMessageUpdate:
		d.messageUpdateChan <- box.(*MessageUpdate)
	case EventMessageDelete:
		d.messageDeleteChan <- box.(*MessageDelete)
	case EventMessageDeleteBulk:
		d.messageDeleteBulkChan <- box.(*MessageDeleteBulk)
	case EventMessageReactionAdd:
		d.messageReactionAddChan <- box.(*MessageReactionAdd)
	case EventMessageReactionRemove:
		d.messageReactionRemoveChan <- box.(*MessageReactionRemove)
	case EventMessageReactionRemoveAll:
		d.messageReactionRemoveAllChan <- box.(*MessageReactionRemoveAll)
	case EventPresenceUpdate:
		d.presenceUpdateChan <- box.(*PresenceUpdate)
	case EventPresencesReplace:
		d.presencesReplaceChan <- box.(*PresencesReplace)
	case EventTypingStart:
		d.typingStartChan <- box.(*TypingStart)
	case EventUserUpdate:
		d.userUpdateChan <- box.(*UserUpdate)
	case EventVoiceStateUpdate:
		d.voiceStateUpdateChan <- box.(*VoiceStateUpdate)
	case EventVoiceServerUpdate:
		d.voiceServerUpdateChan <- box.(*VoiceServerUpdate)
	case EventWebhooksUpdate:
		d.webhooksUpdateChan <- box.(*WebhooksUpdate)
	default:
		fmt.Printf("------\nTODO\nImplement channel for `%s`\n------\n\n", evtName)
	}
}

func (d *Dispatch) triggerCallbacks(ctx context.Context, evtName string, session Session, box interface{}) {
	switch evtName {
	case EventReady:
		for _, listener := range d.listeners[EventReady] {
			(listener.(ReadyCallback))(session, box.(*Ready))
		}
	case EventResumed:
		for _, listener := range d.listeners[EventResumed] {
			(listener.(ResumedCallback))(session, box.(*Resumed))
		}
	case EventChannelCreate:
		for _, listener := range d.listeners[EventChannelCreate] {
			(listener.(ChannelCreateCallback))(session, box.(*ChannelCreate))
		}
	case EventChannelUpdate:
		for _, listener := range d.listeners[EventChannelUpdate] {
			(listener.(ChannelUpdateCallback))(session, box.(*ChannelUpdate))
		}
	case EventChannelDelete:
		for _, listener := range d.listeners[EventChannelDelete] {
			(listener.(ChannelDeleteCallback))(session, box.(*ChannelDelete))
		}
	case EventChannelPinsUpdate:
		for _, listener := range d.listeners[EventChannelPinsUpdate] {
			(listener.(ChannelPinsUpdateCallback))(session, box.(*ChannelPinsUpdate))
		}
	case EventGuildCreate:
		for _, listener := range d.listeners[EventGuildCreate] {
			(listener.(GuildCreateCallback))(session, box.(*GuildCreate))
		}
	case EventGuildUpdate:
		for _, listener := range d.listeners[EventGuildUpdate] {
			(listener.(GuildUpdateCallback))(session, box.(*GuildUpdate))
		}
	case EventGuildDelete:
		for _, listener := range d.listeners[EventGuildDelete] {
			(listener.(GuildDeleteCallback))(session, box.(*GuildDelete))
		}
	case EventGuildBanAdd:
		for _, listener := range d.listeners[EventGuildBanAdd] {
			(listener.(GuildBanAddCallback))(session, box.(*GuildBanAdd))
		}
	case EventGuildBanRemove:
		for _, listener := range d.listeners[EventGuildBanRemove] {
			(listener.(GuildBanRemoveCallback))(session, box.(*GuildBanRemove))
		}
	case EventGuildEmojisUpdate:
		for _, listener := range d.listeners[EventGuildEmojisUpdate] {
			(listener.(GuildEmojisUpdateCallback))(session, box.(*GuildEmojisUpdate))
		}
	case EventGuildIntegrationsUpdate:
		for _, listener := range d.listeners[EventGuildIntegrationsUpdate] {
			(listener.(GuildIntegrationsUpdateCallback))(session, box.(*GuildIntegrationsUpdate))
		}
	case EventGuildMemberAdd:
		for _, listener := range d.listeners[EventGuildMemberAdd] {
			(listener.(GuildMemberAddCallback))(session, box.(*GuildMemberAdd))
		}
	case EventGuildMemberRemove:
		for _, listener := range d.listeners[EventGuildMemberRemove] {
			(listener.(GuildMemberRemoveCallback))(session, box.(*GuildMemberRemove))
		}
	case EventGuildMemberUpdate:
		for _, listener := range d.listeners[EventGuildMemberUpdate] {
			(listener.(GuildMemberUpdateCallback))(session, box.(*GuildMemberUpdate))
		}
	case EventGuildMembersChunk:
		for _, listener := range d.listeners[EventGuildMembersChunk] {
			(listener.(GuildMembersChunkCallback))(session, box.(*GuildMembersChunk))
		}
	case EventGuildRoleCreate:
		for _, listener := range d.listeners[EventGuildRoleCreate] {
			(listener.(GuildRoleCreateCallback))(session, box.(*GuildRoleCreate))
		}
	case EventGuildRoleUpdate:
		for _, listener := range d.listeners[EventGuildRoleUpdate] {
			(listener.(GuildRoleUpdateCallback))(session, box.(*GuildRoleUpdate))
		}
	case EventGuildRoleDelete:
		for _, listener := range d.listeners[EventGuildRoleDelete] {
			(listener.(GuildRoleDeleteCallback))(session, box.(*GuildRoleDelete))
		}
	case EventMessageCreate:
		for _, listener := range d.listeners[EventMessageCreate] {
			(listener.(MessageCreateCallback))(session, box.(*MessageCreate))
		}
	case EventMessageUpdate:
		for _, listener := range d.listeners[EventMessageUpdate] {
			(listener.(MessageUpdateCallback))(session, box.(*MessageUpdate))
		}
	case EventMessageDelete:
		for _, listener := range d.listeners[EventMessageDelete] {
			(listener.(MessageDeleteCallback))(session, box.(*MessageDelete))
		}
	case EventMessageDeleteBulk:
		for _, listener := range d.listeners[EventMessageDeleteBulk] {
			(listener.(MessageDeleteBulkCallback))(session, box.(*MessageDeleteBulk))
		}
	case EventMessageReactionAdd:
		for _, listener := range d.listeners[EventMessageReactionAdd] {
			(listener.(MessageReactionAddCallback))(session, box.(*MessageReactionAdd))
		}
	case EventMessageReactionRemove:
		for _, listener := range d.listeners[EventMessageReactionRemove] {
			(listener.(MessageReactionRemoveCallback))(session, box.(*MessageReactionRemove))
		}
	case EventMessageReactionRemoveAll:
		for _, listener := range d.listeners[EventMessageReactionRemoveAll] {
			(listener.(MessageReactionRemoveAllCallback))(session, box.(*MessageReactionRemoveAll))
		}
	case EventPresenceUpdate:
		for _, listener := range d.listeners[EventPresenceUpdate] {
			(listener.(PresenceUpdateCallback))(session, box.(*PresenceUpdate))
		}
	case EventPresencesReplace:
		for _, listener := range d.listeners[EventPresencesReplace] {
			(listener.(PresencesReplaceCallback))(session, box.(*PresencesReplace))
		}
	case EventTypingStart:
		for _, listener := range d.listeners[EventTypingStart] {
			(listener.(TypingStartCallback))(session, box.(*TypingStart))
		}
	case EventUserUpdate:
		for _, listener := range d.listeners[EventUserUpdate] {
			(listener.(UserUpdateCallback))(session, box.(*UserUpdate))
		}
	case EventVoiceStateUpdate:
		for _, listener := range d.listeners[EventVoiceStateUpdate] {
			(listener.(VoiceStateUpdateCallback))(session, box.(*VoiceStateUpdate))
		}
	case EventVoiceServerUpdate:
		for _, listener := range d.listeners[EventVoiceServerUpdate] {
			(listener.(VoiceServerUpdateCallback))(session, box.(*VoiceServerUpdate))
		}
	case EventWebhooksUpdate:
		for _, listener := range d.listeners[EventWebhooksUpdate] {
			(listener.(WebhooksUpdateCallback))(session, box.(*WebhooksUpdate))
		}
	default:
		fmt.Printf("------\nTODO\nImplement callback for `%s`\n------\n\n", evtName)
	}

	// remove the run only once listeners
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	for _, index := range d.listenOnceOnly[evtName] {
		// https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
		d.listeners[evtName][index] = d.listeners[evtName][len(d.listeners[evtName])-1]
		d.listeners[evtName][len(d.listeners[evtName])-1] = nil
		d.listeners[evtName] = d.listeners[evtName][:len(d.listeners[evtName])-1]

		if len(d.listeners[evtName]) == 0 {
			// TODO: call removeEvent from socket pkg
		}
	}

	// remove the once only register
	if _, exists := d.listenOnceOnly[evtName]; exists {
		delete(d.listenOnceOnly, evtName)
	}
}

// AllChan sends all event types
func (d *Dispatch) All() <-chan interface{} {
	d.ws.RegisterEvent("*")
	return d.allChan
}

// ReadyChan for READY events
func (d *Dispatch) Ready() <-chan *Ready {
	d.ws.RegisterEvent(event.Ready)
	return d.readyChan
}

// ResumedChan for RESUME events
func (d *Dispatch) Resumed() <-chan *Resumed {
	d.ws.RegisterEvent(event.Resumed)
	return d.resumedChan
}

// ChannelCreateChan for CHANNEL_CREATE, when a channel is created
func (d *Dispatch) ChannelCreate() <-chan *ChannelCreate {
	d.ws.RegisterEvent(event.ChannelCreate)
	return d.channelCreateChan
}

// ChannelUpdateChan for CHANNEL_UPDATE, when a channel is updated
func (d *Dispatch) ChannelUpdate() <-chan *ChannelUpdate {
	d.ws.RegisterEvent(event.ChannelUpdate)
	return d.channelUpdateChan
}

// ChannelDeleteChan for CHANNEL_DELETE, when a channel is deleted
func (d *Dispatch) ChannelDelete() <-chan *ChannelDelete {
	d.ws.RegisterEvent(event.ChannelDelete)
	return d.channelDeleteChan
}

// ChannelPinsUpdateChan for CHANNEL_PIN_UPDATE. Message was pinned or unpinned
func (d *Dispatch) ChannelPinsUpdate() <-chan *ChannelPinsUpdate {
	d.ws.RegisterEvent(event.ChannelPinsUpdate)
	return d.channelPinsUpdateChan
}

// GuildCreateChan for GUILD_CREATE. Lazy-load for unavailable guild, guild
// became available, or user joined a new guild
func (d *Dispatch) GuildCreate() <-chan *GuildCreate {
	d.ws.RegisterEvent(event.GuildCreate)
	return d.guildCreateChan
}

// GuildUpdateChan for GUILD_UPDATE. Guild was updated
func (d *Dispatch) GuildUpdate() <-chan *GuildUpdate {
	d.ws.RegisterEvent(event.GuildUpdate)
	return d.guildUpdateChan
}

// GuildDeleteChan for GUILD_DELETE, guild became unavailable, or user
// left/was removed from a guild
func (d *Dispatch) GuildDelete() <-chan *GuildDelete {
	d.ws.RegisterEvent(event.GuildDelete)
	return d.guildDeleteChan
}

// GuildBanAddChan for GUILD_BAN_ADD. A user was banned from a guild
func (d *Dispatch) GuildBanAdd() <-chan *GuildBanAdd {
	d.ws.RegisterEvent(event.GuildBanAdd)
	return d.guildBanAddChan
}

// GuildBanRemoveChan for GUILD_BAN_REMOVE. A user was unbanned from a guild
func (d *Dispatch) GuildBanRemove() <-chan *GuildBanRemove {
	d.ws.RegisterEvent(event.GuildBanRemove)
	return d.guildBanRemoveChan
}

// GuildEmojisUpdateChan for GUILD_EMOJI_UPDATE. Guild emojis were updated
func (d *Dispatch) GuildEmojisUpdate() <-chan *GuildEmojisUpdate {
	d.ws.RegisterEvent(event.GuildEmojisUpdate)
	return d.guildEmojisUpdateChan
}

// GuildIntegrationsUpdateChan for GUILD_INTEGRATIONS_UPDATE. Guild integration
// was updated
func (d *Dispatch) GuildIntegrationsUpdate() <-chan *GuildIntegrationsUpdate {
	d.ws.RegisterEvent(event.GuildIntegrationsUpdate)
	return d.guildIntegrationsUpdateChan
}

// GuildMemberAddChan for GUILD_MEMBER_ADD. New user joined a guild.
func (d *Dispatch) GuildMemberAdd() <-chan *GuildMemberAdd {
	d.ws.RegisterEvent(event.GuildMemberAdd)
	return d.guildMemberAddChan
}

// GuildMemberRemoveChan for GUILD_MEMBER_REMOVE. User was removed from guild.
func (d *Dispatch) GuildMemberRemove() <-chan *GuildMemberRemove {
	d.ws.RegisterEvent(event.GuildMemberRemove)
	return d.guildMemberRemoveChan
}

// GuildMemberUpdateChan for GUILD_MEMBER_UPDATE. Guild member was updated.
func (d *Dispatch) GuildMemberUpdate() <-chan *GuildMemberUpdate {
	d.ws.RegisterEvent(event.GuildMemberUpdate)
	return d.guildMemberUpdateChan
}

// GuildMembersChunkChan for GUILD_MEMBERS_CHUNK. Response to socket command
// 'Request Guild Members'
func (d *Dispatch) GuildMembersChunk() <-chan *GuildMembersChunk {
	d.ws.RegisterEvent(event.GuildMembersChunk)
	return d.guildMembersChunkChan
}

// GuildRoleCreateChan for GUILD_ROLE_CREATE. Guild role was created.
func (d *Dispatch) GuildRoleCreate() <-chan *GuildRoleCreate {
	d.ws.RegisterEvent(event.GuildRoleCreate)
	return d.guildRoleCreateChan
}

// GuildRoleUpdateChan for GUILD_ROLE_UPDATE. Guild role was updated.
func (d *Dispatch) GuildRoleUpdate() <-chan *GuildRoleUpdate {
	d.ws.RegisterEvent(event.GuildRoleUpdate)
	return d.guildRoleUpdateChan
}

// GuildRoleDeleteChan for GUILD_ROLE_DELETE. Guild role was deleted.
func (d *Dispatch) GuildRoleDelete() <-chan *GuildRoleDelete {
	d.ws.RegisterEvent(event.GuildRoleDelete)
	return d.guildRoleDeleteChan
}

// MessageCreateChan for MESSAGE_CREATE. New message was created.
func (d *Dispatch) MessageCreate() <-chan *MessageCreate {
	d.ws.RegisterEvent(event.MessageCreate)
	return d.messageCreateChan
}

// MessageUpdateChan for MESSAGE_UPDATE. Message was updated.
func (d *Dispatch) MessageUpdate() <-chan *MessageUpdate {
	d.ws.RegisterEvent(event.MessageUpdate)
	return d.messageUpdateChan
}

// MessageDeleteChan for MESSAGE_DELETE. Message was deleted.
func (d *Dispatch) MessageDelete() <-chan *MessageDelete {
	d.ws.RegisterEvent(event.MessageDelete)
	return d.messageDeleteChan
}

// MessageDeleteBulkChan for MESSAGE_DELETE_BULK. Multiple messages were
// deleted at once.
func (d *Dispatch) MessageDeleteBulk() <-chan *MessageDeleteBulk {
	d.ws.RegisterEvent(event.MessageDeleteBulk)
	return d.messageDeleteBulkChan
}

// MessageReactionAddChan for MESSAGE_REACTION_ADD. A user reacted to a message.
func (d *Dispatch) MessageReactionAdd() <-chan *MessageReactionAdd {
	d.ws.RegisterEvent(event.MessageReactionAdd)
	return d.messageReactionAddChan
}

// MessageReactionRemoveChan for MESSAGE_REACTION_REMOVE. A user removed a
// a reaction to a message.
func (d *Dispatch) MessageReactionRemove() <-chan *MessageReactionRemove {
	d.ws.RegisterEvent(event.MessageReactionRemove)
	return d.messageReactionRemoveChan
}

// MessageReactionRemoveAllChan for MESSAGE_REACTION_REMOVE_ALL. All reactions
// were explicitly removed from a message
func (d *Dispatch) MessageReactionRemoveAll() <-chan *MessageReactionRemoveAll {
	d.ws.RegisterEvent(event.MessageReactionRemoveAll)
	return d.messageReactionRemoveAllChan
}

// PresenceUpdateChan for PRESENCE_UPDATE. A user's presence was updated in a
// guild.
func (d *Dispatch) PresenceUpdate() <-chan *PresenceUpdate {
	d.ws.RegisterEvent(event.PresenceUpdate)
	return d.presenceUpdateChan
}

// PresenceUpdateChan for PRESENCE_UPDATE. A user's presence was updated in a
// guild.
func (d *Dispatch) PresencesReplace() <-chan *PresencesReplace {
	d.ws.RegisterEvent(event.PresencesReplace)
	return d.presencesReplaceChan
}

// TypingStartChan for TYPING_START. A user started typing in a channel.
func (d *Dispatch) TypingStart() <-chan *TypingStart {
	d.ws.RegisterEvent(event.TypingStart)
	return d.typingStartChan
}

// UserUpdateChan for USER_UPDATE. Properties about a user changed
func (d *Dispatch) UserUpdate() <-chan *UserUpdate {
	d.ws.RegisterEvent(event.UserUpdate)
	return d.userUpdateChan
}

// VoiceStateUpdateChan for VOICE_STATE_UPDATE. Someone joined, left, or moved
// a voice channel
func (d *Dispatch) VoiceStateUpdate() <-chan *VoiceStateUpdate {
	d.ws.RegisterEvent(event.VoiceStateUpdate)
	return d.voiceStateUpdateChan
}

// VoiceServerUpdateChan for VOICE_SERVER_UPDATE. Guild's voice server was
// updated
func (d *Dispatch) VoiceServerUpdate() <-chan *VoiceServerUpdate {
	d.ws.RegisterEvent(event.VoiceServerUpdate)
	return d.voiceServerUpdateChan
}

// WebhooksUpdateChan for WEBHOOK_UPDATE. A guild channel webhook was created,
// update, or deleted
func (d *Dispatch) WebhooksUpdate() <-chan *WebhooksUpdate {
	d.ws.RegisterEvent(event.WebhooksUpdate)
	return d.webhooksUpdateChan
}
