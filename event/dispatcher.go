package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/disgordctx"
	"github.com/andersfylling/disgord/guild"
)

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
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

		//listeners:     make(map[KeyType][]EventCallback),
	}

	dispatcher.alwaysListenToChans()

	return dispatcher
}

//
type DispatcherInterface interface {
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

	Trigger(evtName KeyType, ctx disgordctx.Context, data []byte)
	//OnEvent(evtName KeyType, listener EventCallback)
}

type Dispatcher struct {
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

	//listeners     map[KeyType][]EventCallback
}

// On places listeners into their respected stacks
// func (d *Dispatcher) OnEvent(evtName KeyType, listener EventCallback) {
// 	d.listeners[evtName] = append(d.listeners[evtName], listener)
// }

// alwaysListenToChans makes sure no deadlocks occure
func (d *Dispatcher) alwaysListenToChans() {
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

// Trigger listeners based on the event type
func (d *Dispatcher) Trigger(evtName KeyType, session disgordctx.Context, ctx context.Context, data []byte) {
	// TODO: send data to allChan
	switch evtName {
	case ReadyKey:
		r := &ReadyBox{}
		r.Ctx = ctx
		Unmarshal(data, r)

		d.readyChan <- r
	case ResumedKey:
		resumed := &ResumedBox{}
		resumed.Ctx = ctx
		Unmarshal(data, resumed)

		d.resumedChan <- resumed
	case ChannelCreateKey, ChannelUpdateKey, ChannelDeleteKey:
		chanContent := &channel.Channel{}
		Unmarshal(data, chanContent)

		switch evtName { // internal switch statement for ChannelEvt
		case ChannelCreateKey:
			d.channelCreateChan <- &ChannelCreateBox{Channel: chanContent, Ctx: ctx}
		case ChannelUpdateKey:
			d.channelUpdateChan <- &ChannelUpdateBox{Channel: chanContent, Ctx: ctx}
		case ChannelDeleteKey:
			d.channelDeleteChan <- &ChannelDeleteBox{Channel: chanContent, Ctx: ctx}
		} // END internal switch statement for ChannelEvt
	case ChannelPinsUpdateKey:
		cpu := &ChannelPinsUpdateBox{}
		cpu.Ctx = ctx
		Unmarshal(data, cpu)
		d.channelPinsUpdateChan <- cpu
	case GuildCreateKey, GuildUpdateKey, GuildDeleteKey:
		g := &guild.Guild{}
		Unmarshal(data, g)

		switch evtName { // internal switch statement for guild events
		case GuildCreateKey:
			d.guildCreateChan <- &GuildCreateBox{Guild: g, Ctx: ctx}
		case GuildUpdateKey:
			d.guildUpdateChan <- &GuildUpdateBox{Guild: g, Ctx: ctx}
		case GuildDeleteKey:
			unavailGuild := guild.NewGuildUnavailable(g.ID)
			d.guildDeleteChan <- &GuildDeleteBox{UnavailableGuild: unavailGuild, Ctx: ctx}
		} // END internal switch statement for guild events
	case GuildBanAddKey:
		gba := &GuildBanAddBox{}
		gba.Ctx = ctx
		Unmarshal(data, gba)

		d.guildBanAddChan <- gba
	case GuildBanRemoveKey:
		gbr := &GuildBanRemoveBox{}
		gbr.Ctx = ctx
		Unmarshal(data, gbr)

		d.guildBanRemoveChan <- gbr
	case GuildEmojisUpdateKey:
		geu := &GuildEmojisUpdateBox{}
		geu.Ctx = ctx
		Unmarshal(data, geu)

		d.guildEmojisUpdateChan <- geu
	case GuildIntegrationsUpdateKey:
		giu := &GuildIntegrationsUpdateBox{}
		giu.Ctx = ctx
		Unmarshal(data, giu)

		d.guildIntegrationsUpdateChan <- giu
	case GuildMemberAddKey:
		gma := &GuildMemberAddBox{}
		gma.Ctx = ctx
		Unmarshal(data, gma)

		d.guildMemberAddChan <- gma
	case GuildMemberRemoveKey:
		gmr := &GuildMemberRemoveBox{}
		gmr.Ctx = ctx
		Unmarshal(data, gmr)

		d.guildMemberRemoveChan <- gmr
	case GuildMemberUpdateKey:
		gmu := &GuildMemberUpdateBox{}
		gmu.Ctx = ctx
		Unmarshal(data, gmu)

		d.guildMemberUpdateChan <- gmu
	case GuildMembersChunkKey:
		gmc := &GuildMembersChunkBox{}
		gmc.Ctx = ctx
		Unmarshal(data, gmc)

		d.guildMembersChunkChan <- gmc
	case GuildRoleCreateKey:
		r := &GuildRoleCreateBox{}
		r.Ctx = ctx
		Unmarshal(data, r)

		d.guildRoleCreateChan <- r
	case GuildRoleUpdateKey:
		r := &GuildRoleUpdateBox{}
		r.Ctx = ctx
		Unmarshal(data, r)

		d.guildRoleUpdateChan <- r
	case GuildRoleDeleteKey:
		r := &GuildRoleDeleteBox{}
		r.Ctx = ctx
		Unmarshal(data, r)

		d.guildRoleDeleteChan <- r
	case MessageCreateKey, MessageUpdateKey, MessageDeleteKey:
		msg := channel.NewMessage()
		Unmarshal(data, msg)

		switch evtName { // internal switch statement for MessageEvt
		case MessageCreateKey:
			d.messageCreateChan <- &MessageCreateBox{Message: msg, Ctx: ctx}
		case MessageUpdateKey:
			d.messageUpdateChan <- &MessageUpdateBox{Message: msg, Ctx: ctx}
		case MessageDeleteKey:
			d.messageDeleteChan <- &MessageDeleteBox{MessageID: msg.ID, ChannelID: msg.ChannelID}
		} // END internal switch statement for MessageEvt
	case MessageDeleteBulkKey:
		mdb := &MessageDeleteBulkBox{}
		mdb.Ctx = ctx
		Unmarshal(data, mdb)

		d.messageDeleteBulkChan <- mdb
	case MessageReactionAddKey:
		mra := &MessageReactionAddBox{}
		mra.Ctx = ctx
		Unmarshal(data, mra)

		d.messageReactionAddChan <- mra
	case MessageReactionRemoveKey:
		mrr := &MessageReactionRemoveBox{}
		mrr.Ctx = ctx
		Unmarshal(data, mrr)

		d.messageReactionRemoveChan <- mrr
	case MessageReactionRemoveAllKey:
		mrra := &MessageReactionRemoveAllBox{}
		mrra.Ctx = ctx
		Unmarshal(data, mrra)

		d.messageReactionRemoveAllChan <- mrra
	case PresenceUpdateKey:
		pu := &PresenceUpdateBox{}
		pu.Ctx = ctx
		Unmarshal(data, pu)

		d.presenceUpdateChan <- pu
	case TypingStartKey:
		ts := &TypingStartBox{}
		ts.Ctx = ctx
		Unmarshal(data, ts)

		d.typingStartChan <- ts
	case UserUpdateKey:
		u := &UserUpdateBox{}
		u.Ctx = ctx
		Unmarshal(data, u)

		// dispatch event
		d.userUpdateChan <- u
	case VoiceStateUpdateKey:
		vsu := &VoiceStateUpdateBox{}
		vsu.Ctx = ctx
		Unmarshal(data, vsu)

		d.voiceStateUpdateChan <- vsu
	case VoiceServerUpdateKey:
		vsu := &VoiceServerUpdateBox{}
		vsu.Ctx = ctx
		Unmarshal(data, vsu)

		d.voiceServerUpdateChan <- vsu
	case WebhooksUpdateKey:
		wsu := &WebhooksUpdateBox{}
		wsu.Ctx = ctx
		Unmarshal(data, wsu)

		d.webhooksUpdateChan <- wsu
	default:
		fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evtName, string(data))
	}

	// trigger callbacks
	// for _, listener := range d.listeners[evtName] {
	// 	go listener(ctx, box)
	// }
}

func (d *Dispatcher) AllChan() <-chan interface{} {
	return d.allChan
}
func (d *Dispatcher) ReadyChan() <-chan *ReadyBox {
	return d.readyChan
}
func (d *Dispatcher) ResumedChan() <-chan *ResumedBox {
	return d.resumedChan
}
func (d *Dispatcher) ChannelCreateChan() <-chan *ChannelCreateBox {
	return d.channelCreateChan
}
func (d *Dispatcher) ChannelUpdateChan() <-chan *ChannelUpdateBox {
	return d.channelUpdateChan
}
func (d *Dispatcher) ChannelDeleteChan() <-chan *ChannelDeleteBox {
	return d.channelDeleteChan
}
func (d *Dispatcher) ChannelPinsUpdateChan() <-chan *ChannelPinsUpdateBox {
	return d.channelPinsUpdateChan
}
func (d *Dispatcher) GuildCreateChan() <-chan *GuildCreateBox {
	return d.guildCreateChan
}
func (d *Dispatcher) GuildUpdateChan() <-chan *GuildUpdateBox {
	return d.guildUpdateChan
}
func (d *Dispatcher) GuildDeleteChan() <-chan *GuildDeleteBox {
	return d.guildDeleteChan
}
func (d *Dispatcher) GuildBanAddChan() <-chan *GuildBanAddBox {
	return d.guildBanAddChan
}
func (d *Dispatcher) GuildBanRemoveChan() <-chan *GuildBanRemoveBox {
	return d.guildBanRemoveChan
}
func (d *Dispatcher) GuildEmojisUpdateChan() <-chan *GuildEmojisUpdateBox {
	return d.guildEmojisUpdateChan
}
func (d *Dispatcher) GuildIntegrationsUpdateChan() <-chan *GuildIntegrationsUpdateBox {
	return d.guildIntegrationsUpdateChan
}
func (d *Dispatcher) GuildMemberAddChan() <-chan *GuildMemberAddBox { return d.guildMemberAddChan }
func (d *Dispatcher) GuildMemberRemoveChan() <-chan *GuildMemberRemoveBox {
	return d.guildMemberRemoveChan
}
func (d *Dispatcher) GuildMemberUpdateChan() <-chan *GuildMemberUpdateBox {
	return d.guildMemberUpdateChan
}
func (d *Dispatcher) GuildMembersChunkChan() <-chan *GuildMembersChunkBox {
	return d.guildMembersChunkChan
}
func (d *Dispatcher) GuildRoleUpdateChan() <-chan *GuildRoleUpdateBox {
	return d.guildRoleUpdateChan
}
func (d *Dispatcher) GuildRoleCreateChan() <-chan *GuildRoleCreateBox {
	return d.guildRoleCreateChan
}
func (d *Dispatcher) GuildRoleDeleteChan() <-chan *GuildRoleDeleteBox {
	return d.guildRoleDeleteChan
}
func (d *Dispatcher) MessageCreateChan() <-chan *MessageCreateBox {
	return d.messageCreateChan
}
func (d *Dispatcher) MessageUpdateChan() <-chan *MessageUpdateBox {
	return d.messageUpdateChan
}
func (d *Dispatcher) MessageDeleteChan() <-chan *MessageDeleteBox {
	return d.messageDeleteChan
}
func (d *Dispatcher) MessageDeleteBulkChan() <-chan *MessageDeleteBulkBox {
	return d.messageDeleteBulkChan
}
func (d *Dispatcher) MessageReactionAddChan() <-chan *MessageReactionAddBox {
	return d.messageReactionAddChan
}
func (d *Dispatcher) MessageReactionRemoveChan() <-chan *MessageReactionRemoveBox {
	return d.messageReactionRemoveChan
}
func (d *Dispatcher) MessageReactionRemoveAllChan() <-chan *MessageReactionRemoveAllBox {
	return d.messageReactionRemoveAllChan
}
func (d *Dispatcher) PresenceUpdateChan() <-chan *PresenceUpdateBox {
	return d.presenceUpdateChan
}
func (d *Dispatcher) TypingStartChan() <-chan *TypingStartBox {
	return d.typingStartChan
}
func (d *Dispatcher) UserUpdateChan() <-chan *UserUpdateBox {
	return d.userUpdateChan
}
func (d *Dispatcher) VoiceStateUpdateChan() <-chan *VoiceStateUpdateBox {
	return d.voiceStateUpdateChan
}
func (d *Dispatcher) VoiceServerUpdateChan() <-chan *VoiceServerUpdateBox {
	return d.voiceServerUpdateChan
}
func (d *Dispatcher) WebhooksUpdateChan() <-chan *WebhooksUpdateBox {
	return d.webhooksUpdateChan
}

// wtf is this
func Unmarshal(data []byte, box interface{}) {
	err := json.Unmarshal(data, box)
	if err != nil {
		panic(err)
	}
}
