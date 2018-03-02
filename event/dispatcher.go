package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/disgordctx"
	"github.com/andersfylling/disgord/guild"
)

func NewDispatch() *Dispatch {
	dispatcher := &Dispatch{
		allChan:                      make(chan interface{}),
		readyChan:                    make(chan *ReadyBox),
		resumedChan:                  make(chan *ResumedBox),
		channelCreateChan:            make(chan *ChannelCreateBox),
		channelUpdateChan:            make(chan *ChannelUpdateBox),
		channelDeleteChan:            make(chan *ChannelDeleteBox),
		channelPinsUpdateChan:        make(chan *ChannelPinsUpdateBox),
		guildCreateChan:              make(chan *GuildCreateBox),
		guildUpdateChan:              make(chan *GuildUpdateBox),
		guildDeleteChan:              make(chan *GuildDeleteBox),
		guildBanAddChan:              make(chan *GuildBanAddBox),
		guildBanRemoveChan:           make(chan *GuildBanRemoveBox),
		guildEmojisUpdateChan:        make(chan *GuildEmojisUpdateBox),
		guildIntegrationsUpdateChan:  make(chan *GuildIntegrationsUpdateBox),
		guildMemberAddChan:           make(chan *GuildMemberAddBox),
		guildMemberRemoveChan:        make(chan *GuildMemberRemoveBox),
		guildMemberUpdateChan:        make(chan *GuildMemberUpdateBox),
		guildMembersChunkChan:        make(chan *GuildMembersChunkBox),
		guildRoleUpdateChan:          make(chan *GuildRoleUpdateBox),
		guildRoleCreateChan:          make(chan *GuildRoleCreateBox),
		guildRoleDeleteChan:          make(chan *GuildRoleDeleteBox),
		messageCreateChan:            make(chan *MessageCreateBox),
		messageUpdateChan:            make(chan *MessageUpdateBox),
		messageDeleteChan:            make(chan *MessageDeleteBox),
		messageDeleteBulkChan:        make(chan *MessageDeleteBulkBox),
		messageReactionAddChan:       make(chan *MessageReactionAddBox),
		messageReactionRemoveChan:    make(chan *MessageReactionRemoveBox),
		messageReactionRemoveAllChan: make(chan *MessageReactionRemoveAllBox),
		presenceUpdateChan:           make(chan *PresenceUpdateBox),
		typingStartChan:              make(chan *TypingStartBox),
		userUpdateChan:               make(chan *UserUpdateBox),
		voiceStateUpdateChan:         make(chan *VoiceStateUpdateBox),
		voiceServerUpdateChan:        make(chan *VoiceServerUpdateBox),
		webhooksUpdateChan:           make(chan *WebhooksUpdateBox),

		listeners: make(map[string][]interface{}),
	}

	// make sure every channel has a reciever to avoid deadlock
	// hack...
	dispatcher.alwaysListenToChans()

	return dispatcher
}

//
type Dispatcher interface {
	AllChan() <-chan interface{} // any event
	ReadyChan() <-chan *ReadyBox
	ResumedChan() <-chan *ResumedBox
	ChannelCreateChan() <-chan *ChannelCreateBox
	ChannelUpdateChan() <-chan *ChannelUpdateBox
	ChannelDeleteChan() <-chan *ChannelDeleteBox
	ChannelPinsUpdateChan() <-chan *ChannelPinsUpdateBox
	GuildCreateChan() <-chan *GuildCreateBox
	GuildUpdateChan() <-chan *GuildUpdateBox
	GuildDeleteChan() <-chan *GuildDeleteBox
	GuildBanAddChan() <-chan *GuildBanAddBox
	GuildBanRemoveChan() <-chan *GuildBanRemoveBox
	GuildEmojisUpdateChan() <-chan *GuildEmojisUpdateBox
	GuildIntegrationsUpdateChan() <-chan *GuildIntegrationsUpdateBox
	GuildMemberAddChan() <-chan *GuildMemberAddBox
	GuildMemberRemoveChan() <-chan *GuildMemberRemoveBox
	GuildMemberUpdateChan() <-chan *GuildMemberUpdateBox
	GuildMembersChunkChan() <-chan *GuildMembersChunkBox
	GuildRoleUpdateChan() <-chan *GuildRoleUpdateBox
	GuildRoleCreateChan() <-chan *GuildRoleCreateBox
	GuildRoleDeleteChan() <-chan *GuildRoleDeleteBox
	MessageCreateChan() <-chan *MessageCreateBox
	MessageUpdateChan() <-chan *MessageUpdateBox
	MessageDeleteChan() <-chan *MessageDeleteBox
	MessageDeleteBulkChan() <-chan *MessageDeleteBulkBox
	MessageReactionAddChan() <-chan *MessageReactionAddBox
	MessageReactionRemoveChan() <-chan *MessageReactionRemoveBox
	MessageReactionRemoveAllChan() <-chan *MessageReactionRemoveAllBox
	PresenceUpdateChan() <-chan *PresenceUpdateBox
	TypingStartChan() <-chan *TypingStartBox
	UserUpdateChan() <-chan *UserUpdateBox
	VoiceStateUpdateChan() <-chan *VoiceStateUpdateBox
	VoiceServerUpdateChan() <-chan *VoiceServerUpdateBox
	WebhooksUpdateChan() <-chan *WebhooksUpdateBox

	Trigger(evtName string, session disgordctx.Context, ctx context.Context, data []byte)
	AddHandler(evtName string, listener interface{})
}

type Dispatch struct {
	allChan                      chan interface{} // any event
	readyChan                    chan *ReadyBox
	resumedChan                  chan *ResumedBox
	channelCreateChan            chan *ChannelCreateBox
	channelUpdateChan            chan *ChannelUpdateBox
	channelDeleteChan            chan *ChannelDeleteBox
	channelPinsUpdateChan        chan *ChannelPinsUpdateBox
	guildCreateChan              chan *GuildCreateBox
	guildUpdateChan              chan *GuildUpdateBox
	guildDeleteChan              chan *GuildDeleteBox
	guildBanAddChan              chan *GuildBanAddBox
	guildBanRemoveChan           chan *GuildBanRemoveBox
	guildEmojisUpdateChan        chan *GuildEmojisUpdateBox
	guildIntegrationsUpdateChan  chan *GuildIntegrationsUpdateBox
	guildMemberAddChan           chan *GuildMemberAddBox
	guildMemberRemoveChan        chan *GuildMemberRemoveBox
	guildMemberUpdateChan        chan *GuildMemberUpdateBox
	guildMembersChunkChan        chan *GuildMembersChunkBox
	guildRoleUpdateChan          chan *GuildRoleUpdateBox
	guildRoleCreateChan          chan *GuildRoleCreateBox
	guildRoleDeleteChan          chan *GuildRoleDeleteBox
	messageCreateChan            chan *MessageCreateBox
	messageUpdateChan            chan *MessageUpdateBox
	messageDeleteChan            chan *MessageDeleteBox
	messageDeleteBulkChan        chan *MessageDeleteBulkBox
	messageReactionAddChan       chan *MessageReactionAddBox
	messageReactionRemoveChan    chan *MessageReactionRemoveBox
	messageReactionRemoveAllChan chan *MessageReactionRemoveAllBox
	presenceUpdateChan           chan *PresenceUpdateBox
	typingStartChan              chan *TypingStartBox
	userUpdateChan               chan *UserUpdateBox
	voiceStateUpdateChan         chan *VoiceStateUpdateBox
	voiceServerUpdateChan        chan *VoiceServerUpdateBox
	webhooksUpdateChan           chan *WebhooksUpdateBox

	listeners map[string][]interface{}
}

// On places listeners into their respected stacks
// func (d *Dispatcher) OnEvent(evtName string, listener EventCallback) {
// 	d.listeners[evtName] = append(d.listeners[evtName], listener)
// }

// alwaysListenToChans makes sure no deadlocks occure
func (d *Dispatch) alwaysListenToChans() {
	go func() {
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
			case <-d.typingStartChan:
			case <-d.userUpdateChan:
			case <-d.voiceStateUpdateChan:
			case <-d.voiceServerUpdateChan:
			case <-d.webhooksUpdateChan:
			}
		}
	}()
}

func (d *Dispatch) AddHandler(evtName string, listener interface{}) {
	d.listeners[evtName] = append(d.listeners[evtName], listener)
}

func (d *Dispatch) TriggerChan(evtName string, session disgordctx.Context, ctx context.Context, box interface{}) {
	switch evtName {
	case ReadyKey:
		d.readyChan <- box.(*ReadyBox)
	case ResumedKey:
		d.resumedChan <- box.(*ResumedBox)
	case ChannelCreateKey:
		d.channelCreateChan <- box.(*ChannelCreateBox)
	case ChannelUpdateKey:
		d.channelUpdateChan <- box.(*ChannelUpdateBox)
	case ChannelDeleteKey:
		d.channelDeleteChan <- box.(*ChannelDeleteBox)
	case ChannelPinsUpdateKey:
		d.channelPinsUpdateChan <- box.(*ChannelPinsUpdateBox)
	case GuildCreateKey:
		d.guildCreateChan <- box.(*GuildCreateBox)
	case GuildUpdateKey:
		d.guildUpdateChan <- box.(*GuildUpdateBox)
	case GuildDeleteKey:
		d.guildDeleteChan <- box.(*GuildDeleteBox)
	case GuildBanAddKey:
		d.guildBanAddChan <- box.(*GuildBanAddBox)
	case GuildBanRemoveKey:
		d.guildBanRemoveChan <- box.(*GuildBanRemoveBox)
	case GuildEmojisUpdateKey:
		d.guildEmojisUpdateChan <- box.(*GuildEmojisUpdateBox)
	case GuildIntegrationsUpdateKey:
		d.guildIntegrationsUpdateChan <- box.(*GuildIntegrationsUpdateBox)
	case GuildMemberAddKey:
		d.guildMemberAddChan <- box.(*GuildMemberAddBox)
	case GuildMemberRemoveKey:
		d.guildMemberRemoveChan <- box.(*GuildMemberRemoveBox)
	case GuildMemberUpdateKey:
		d.guildMemberUpdateChan <- box.(*GuildMemberUpdateBox)
	case GuildMembersChunkKey:
		d.guildMembersChunkChan <- box.(*GuildMembersChunkBox)
	case GuildRoleCreateKey:
		d.guildRoleCreateChan <- box.(*GuildRoleCreateBox)
	case GuildRoleUpdateKey:
		d.guildRoleUpdateChan <- box.(*GuildRoleUpdateBox)
	case GuildRoleDeleteKey:
		d.guildRoleDeleteChan <- box.(*GuildRoleDeleteBox)
	case MessageCreateKey:
		d.messageCreateChan <- box.(*MessageCreateBox)
	case MessageUpdateKey:
		d.messageUpdateChan <- box.(*MessageUpdateBox)
	case MessageDeleteKey:
		d.messageDeleteChan <- box.(*MessageDeleteBox)
	case MessageDeleteBulkKey:
		d.messageDeleteBulkChan <- box.(*MessageDeleteBulkBox)
	case MessageReactionAddKey:
		d.messageReactionAddChan <- box.(*MessageReactionAddBox)
	case MessageReactionRemoveKey:
		d.messageReactionRemoveChan <- box.(*MessageReactionRemoveBox)
	case MessageReactionRemoveAllKey:
		d.messageReactionRemoveAllChan <- box.(*MessageReactionRemoveAllBox)
	case PresenceUpdateKey:
		d.presenceUpdateChan <- box.(*PresenceUpdateBox)
	case TypingStartKey:
		d.typingStartChan <- box.(*TypingStartBox)
	case UserUpdateKey:
		d.userUpdateChan <- box.(*UserUpdateBox)
	case VoiceStateUpdateKey:
		d.voiceStateUpdateChan <- box.(*VoiceStateUpdateBox)
	case VoiceServerUpdateKey:
		d.voiceServerUpdateChan <- box.(*VoiceServerUpdateBox)
	case WebhooksUpdateKey:
		d.webhooksUpdateChan <- box.(*WebhooksUpdateBox)
	default:
		fmt.Printf("------\nTODO\nImplement channel for `%s`\n------\n\n", evtName)
	}
}

func (d *Dispatch) TriggerCallbacks(evtName string, session disgordctx.Context, ctx context.Context, box interface{}) {
	switch evtName {
	case ReadyKey:
		for _, listener := range d.listeners[ReadyKey] {
			go (listener.(ReadyCallback))(session, box.(*ReadyBox))
		}
	case ResumedKey:
		for _, listener := range d.listeners[ResumedKey] {
			go (listener.(ResumedCallback))(session, box.(*ResumedBox))
		}
	case ChannelCreateKey:
		for _, listener := range d.listeners[ChannelCreateKey] {
			go (listener.(ChannelCreateCallback))(session, box.(*ChannelCreateBox))
		}
	case ChannelUpdateKey:
		for _, listener := range d.listeners[ChannelUpdateKey] {
			go (listener.(ChannelUpdateCallback))(session, box.(*ChannelUpdateBox))
		}
	case ChannelDeleteKey:
		for _, listener := range d.listeners[ChannelDeleteKey] {
			go (listener.(ChannelDeleteCallback))(session, box.(*ChannelDeleteBox))
		}
	case ChannelPinsUpdateKey:
		for _, listener := range d.listeners[ChannelPinsUpdateKey] {
			go (listener.(ChannelPinsUpdateCallback))(session, box.(*ChannelPinsUpdateBox))
		}
	case GuildCreateKey:
		for _, listener := range d.listeners[GuildCreateKey] {
			go (listener.(GuildCreateCallback))(session, box.(*GuildCreateBox))
		}
	case GuildUpdateKey:
		for _, listener := range d.listeners[GuildUpdateKey] {
			go (listener.(GuildUpdateCallback))(session, box.(*GuildUpdateBox))
		}
	case GuildDeleteKey:
		for _, listener := range d.listeners[GuildDeleteKey] {
			go (listener.(GuildDeleteCallback))(session, box.(*GuildDeleteBox))
		}
	case GuildBanAddKey:
		for _, listener := range d.listeners[GuildBanAddKey] {
			go (listener.(GuildBanAddCallback))(session, box.(*GuildBanAddBox))
		}
	case GuildBanRemoveKey:
		for _, listener := range d.listeners[GuildBanRemoveKey] {
			go (listener.(GuildBanRemoveCallback))(session, box.(*GuildBanRemoveBox))
		}
	case GuildEmojisUpdateKey:
		for _, listener := range d.listeners[GuildEmojisUpdateKey] {
			go (listener.(GuildEmojisUpdateCallback))(session, box.(*GuildEmojisUpdateBox))
		}
	case GuildIntegrationsUpdateKey:
		for _, listener := range d.listeners[GuildIntegrationsUpdateKey] {
			go (listener.(GuildIntegrationsUpdateCallback))(session, box.(*GuildIntegrationsUpdateBox))
		}
	case GuildMemberAddKey:
		for _, listener := range d.listeners[GuildMemberAddKey] {
			go (listener.(GuildMemberAddCallback))(session, box.(*GuildMemberAddBox))
		}
	case GuildMemberRemoveKey:
		for _, listener := range d.listeners[GuildMemberRemoveKey] {
			go (listener.(GuildMemberRemoveCallback))(session, box.(*GuildMemberRemoveBox))
		}
	case GuildMemberUpdateKey:
		for _, listener := range d.listeners[GuildMemberUpdateKey] {
			go (listener.(GuildMemberUpdateCallback))(session, box.(*GuildMemberUpdateBox))
		}
	case GuildMembersChunkKey:
		for _, listener := range d.listeners[GuildMembersChunkKey] {
			go (listener.(GuildMembersChunkCallback))(session, box.(*GuildMembersChunkBox))
		}
	case GuildRoleCreateKey:
		for _, listener := range d.listeners[GuildRoleCreateKey] {
			go (listener.(GuildRoleCreateCallback))(session, box.(*GuildRoleCreateBox))
		}
	case GuildRoleUpdateKey:
		for _, listener := range d.listeners[GuildRoleUpdateKey] {
			go (listener.(GuildRoleUpdateCallback))(session, box.(*GuildRoleUpdateBox))
		}
	case GuildRoleDeleteKey:
		for _, listener := range d.listeners[GuildRoleDeleteKey] {
			go (listener.(GuildRoleDeleteCallback))(session, box.(*GuildRoleDeleteBox))
		}
	case MessageCreateKey:
		for _, listener := range d.listeners[MessageCreateKey] {
			go (listener.(MessageCreateCallback))(session, box.(*MessageCreateBox))
		}
	case MessageUpdateKey:
		for _, listener := range d.listeners[MessageUpdateKey] {
			go (listener.(MessageUpdateCallback))(session, box.(*MessageUpdateBox))
		}
	case MessageDeleteKey:
		for _, listener := range d.listeners[MessageDeleteKey] {
			go (listener.(MessageDeleteCallback))(session, box.(*MessageDeleteBox))
		}
	case MessageDeleteBulkKey:
		for _, listener := range d.listeners[MessageDeleteBulkKey] {
			go (listener.(MessageDeleteBulkCallback))(session, box.(*MessageDeleteBulkBox))
		}
	case MessageReactionAddKey:
		for _, listener := range d.listeners[MessageReactionAddKey] {
			go (listener.(MessageReactionAddCallback))(session, box.(*MessageReactionAddBox))
		}
	case MessageReactionRemoveKey:
		for _, listener := range d.listeners[MessageReactionRemoveKey] {
			go (listener.(MessageReactionRemoveCallback))(session, box.(*MessageReactionRemoveBox))
		}
	case MessageReactionRemoveAllKey:
		for _, listener := range d.listeners[MessageReactionRemoveAllKey] {
			go (listener.(MessageReactionRemoveAllCallback))(session, box.(*MessageReactionRemoveAllBox))
		}
	case PresenceUpdateKey:
		for _, listener := range d.listeners[PresenceUpdateKey] {
			go (listener.(PresenceUpdateCallback))(session, box.(*PresenceUpdateBox))
		}
	case TypingStartKey:
		for _, listener := range d.listeners[TypingStartKey] {
			go (listener.(TypingStartCallback))(session, box.(*TypingStartBox))
		}
	case UserUpdateKey:
		for _, listener := range d.listeners[UserUpdateKey] {
			go (listener.(UserUpdateCallback))(session, box.(*UserUpdateBox))
		}
	case VoiceStateUpdateKey:
		for _, listener := range d.listeners[VoiceStateUpdateKey] {
			go (listener.(VoiceStateUpdateCallback))(session, box.(*VoiceStateUpdateBox))
		}
	case VoiceServerUpdateKey:
		for _, listener := range d.listeners[VoiceServerUpdateKey] {
			go (listener.(VoiceServerUpdateCallback))(session, box.(*VoiceServerUpdateBox))
		}
	case WebhooksUpdateKey:
		for _, listener := range d.listeners[WebhooksUpdateKey] {
			go (listener.(WebhooksUpdateCallback))(session, box.(*WebhooksUpdateBox))
		}
	default:
		fmt.Printf("------\nTODO\nImplement callback for `%s`\n------\n\n", evtName)
	}
}

// Trigger listeners based on the event type
func (d *Dispatch) Trigger(evtName string, session disgordctx.Context, ctx context.Context, data []byte) {
	// TODO: send data to allChan
	switch evtName {
	case ReadyKey:
		box := &ReadyBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case ResumedKey:
		box := &ResumedBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case ChannelCreateKey, ChannelUpdateKey, ChannelDeleteKey:
		chanContent := &channel.Channel{}
		Unmarshal(data, chanContent)

		switch evtName { // internal switch statement for ChannelEvt
		case ChannelCreateKey:
			box := &ChannelCreateBox{Channel: chanContent, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case ChannelUpdateKey:
			box := &ChannelUpdateBox{Channel: chanContent, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case ChannelDeleteKey:
			box := &ChannelDeleteBox{Channel: chanContent, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		} // END internal switch statement for ChannelEvt
	case ChannelPinsUpdateKey:
		box := &ChannelPinsUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildCreateKey, GuildUpdateKey, GuildDeleteKey:
		g := &guild.Guild{}
		Unmarshal(data, g)

		switch evtName { // internal switch statement for guild events
		case GuildCreateKey:
			box := &GuildCreateBox{Guild: g, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case GuildUpdateKey:
			box := &GuildUpdateBox{Guild: g, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case GuildDeleteKey:
			unavailGuild := guild.NewGuildUnavailable(g.ID)
			box := &GuildDeleteBox{UnavailableGuild: unavailGuild, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		} // END internal switch statement for guild events
	case GuildBanAddKey:
		box := &GuildBanAddBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildBanRemoveKey:
		box := &GuildBanRemoveBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildEmojisUpdateKey:
		box := &GuildEmojisUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildIntegrationsUpdateKey:
		box := &GuildIntegrationsUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildMemberAddKey:
		box := &GuildMemberAddBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildMemberRemoveKey:
		box := &GuildMemberRemoveBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildMemberUpdateKey:
		box := &GuildMemberUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildMembersChunkKey:
		box := &GuildMembersChunkBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildRoleCreateKey:
		box := &GuildRoleCreateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildRoleUpdateKey:
		box := &GuildRoleUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case GuildRoleDeleteKey:
		box := &GuildRoleDeleteBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case MessageCreateKey, MessageUpdateKey, MessageDeleteKey:
		msg := channel.NewMessage()
		Unmarshal(data, msg)

		switch evtName { // internal switch statement for MessageEvt
		case MessageCreateKey:
			box := &MessageCreateBox{Message: msg, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case MessageUpdateKey:
			box := &MessageUpdateBox{Message: msg, Ctx: ctx}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		case MessageDeleteKey:
			box := &MessageDeleteBox{MessageID: msg.ID, ChannelID: msg.ChannelID}
			d.TriggerChan(evtName, session, ctx, box)
			d.TriggerCallbacks(evtName, session, ctx, box)
		} // END internal switch statement for MessageEvt
	case MessageDeleteBulkKey:
		box := &MessageDeleteBulkBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case MessageReactionAddKey:
		box := &MessageReactionAddBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case MessageReactionRemoveKey:
		box := &MessageReactionRemoveBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case MessageReactionRemoveAllKey:
		box := &MessageReactionRemoveAllBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case PresenceUpdateKey:
		box := &PresenceUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case TypingStartKey:
		box := &TypingStartBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case UserUpdateKey:
		box := &UserUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case VoiceStateUpdateKey:
		box := &VoiceStateUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case VoiceServerUpdateKey:
		box := &VoiceServerUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	case WebhooksUpdateKey:
		box := &WebhooksUpdateBox{}
		box.Ctx = ctx
		Unmarshal(data, box)

		d.TriggerChan(evtName, session, ctx, box)
		d.TriggerCallbacks(evtName, session, ctx, box)
	default:
		fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evtName, string(data))
	}

	// trigger callbacks
	// for _, listener := range d.listeners[evtName] {
	// 	go listener(ctx, box)
	// }
}

func (d *Dispatch) AllChan() <-chan interface{} {
	return d.allChan
}
func (d *Dispatch) ReadyChan() <-chan *ReadyBox {
	return d.readyChan
}
func (d *Dispatch) ResumedChan() <-chan *ResumedBox {
	return d.resumedChan
}
func (d *Dispatch) ChannelCreateChan() <-chan *ChannelCreateBox {
	return d.channelCreateChan
}
func (d *Dispatch) ChannelUpdateChan() <-chan *ChannelUpdateBox {
	return d.channelUpdateChan
}
func (d *Dispatch) ChannelDeleteChan() <-chan *ChannelDeleteBox {
	return d.channelDeleteChan
}
func (d *Dispatch) ChannelPinsUpdateChan() <-chan *ChannelPinsUpdateBox {
	return d.channelPinsUpdateChan
}
func (d *Dispatch) GuildCreateChan() <-chan *GuildCreateBox {
	return d.guildCreateChan
}
func (d *Dispatch) GuildUpdateChan() <-chan *GuildUpdateBox {
	return d.guildUpdateChan
}
func (d *Dispatch) GuildDeleteChan() <-chan *GuildDeleteBox {
	return d.guildDeleteChan
}
func (d *Dispatch) GuildBanAddChan() <-chan *GuildBanAddBox {
	return d.guildBanAddChan
}
func (d *Dispatch) GuildBanRemoveChan() <-chan *GuildBanRemoveBox {
	return d.guildBanRemoveChan
}
func (d *Dispatch) GuildEmojisUpdateChan() <-chan *GuildEmojisUpdateBox {
	return d.guildEmojisUpdateChan
}
func (d *Dispatch) GuildIntegrationsUpdateChan() <-chan *GuildIntegrationsUpdateBox {
	return d.guildIntegrationsUpdateChan
}
func (d *Dispatch) GuildMemberAddChan() <-chan *GuildMemberAddBox { return d.guildMemberAddChan }
func (d *Dispatch) GuildMemberRemoveChan() <-chan *GuildMemberRemoveBox {
	return d.guildMemberRemoveChan
}
func (d *Dispatch) GuildMemberUpdateChan() <-chan *GuildMemberUpdateBox {
	return d.guildMemberUpdateChan
}
func (d *Dispatch) GuildMembersChunkChan() <-chan *GuildMembersChunkBox {
	return d.guildMembersChunkChan
}
func (d *Dispatch) GuildRoleUpdateChan() <-chan *GuildRoleUpdateBox {
	return d.guildRoleUpdateChan
}
func (d *Dispatch) GuildRoleCreateChan() <-chan *GuildRoleCreateBox {
	return d.guildRoleCreateChan
}
func (d *Dispatch) GuildRoleDeleteChan() <-chan *GuildRoleDeleteBox {
	return d.guildRoleDeleteChan
}
func (d *Dispatch) MessageCreateChan() <-chan *MessageCreateBox {
	return d.messageCreateChan
}
func (d *Dispatch) MessageUpdateChan() <-chan *MessageUpdateBox {
	return d.messageUpdateChan
}
func (d *Dispatch) MessageDeleteChan() <-chan *MessageDeleteBox {
	return d.messageDeleteChan
}
func (d *Dispatch) MessageDeleteBulkChan() <-chan *MessageDeleteBulkBox {
	return d.messageDeleteBulkChan
}
func (d *Dispatch) MessageReactionAddChan() <-chan *MessageReactionAddBox {
	return d.messageReactionAddChan
}
func (d *Dispatch) MessageReactionRemoveChan() <-chan *MessageReactionRemoveBox {
	return d.messageReactionRemoveChan
}
func (d *Dispatch) MessageReactionRemoveAllChan() <-chan *MessageReactionRemoveAllBox {
	return d.messageReactionRemoveAllChan
}
func (d *Dispatch) PresenceUpdateChan() <-chan *PresenceUpdateBox {
	return d.presenceUpdateChan
}
func (d *Dispatch) TypingStartChan() <-chan *TypingStartBox {
	return d.typingStartChan
}
func (d *Dispatch) UserUpdateChan() <-chan *UserUpdateBox {
	return d.userUpdateChan
}
func (d *Dispatch) VoiceStateUpdateChan() <-chan *VoiceStateUpdateBox {
	return d.voiceStateUpdateChan
}
func (d *Dispatch) VoiceServerUpdateChan() <-chan *VoiceServerUpdateBox {
	return d.voiceServerUpdateChan
}
func (d *Dispatch) WebhooksUpdateChan() <-chan *WebhooksUpdateBox {
	return d.webhooksUpdateChan
}

// wtf is this
func Unmarshal(data []byte, box interface{}) {
	err := json.Unmarshal(data, box)
	if err != nil {
		panic(err)
	}
}
