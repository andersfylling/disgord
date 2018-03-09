package disgord

import (
	"context"
	"encoding/json"
	"fmt"

	"sync"

	. "github.com/andersfylling/disgord/event"
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

		listeners:      make(map[string][]interface{}),
		listenOnceOnly: make(map[string][]int),
	}

	// make sure every channel has a reciever to avoid deadlock
	// hack...
	dispatcher.alwaysListenToChans()

	return dispatcher
}

//
type EvtDispatcher interface {
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

	AddHandler(evtName string, listener interface{})
	AddHandlerOnce(evtName string, listener interface{})
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

	listeners      map[string][]interface{}
	listenOnceOnly map[string][]int

	listenersLock sync.RWMutex
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

func (d *Dispatch) triggerChan(evtName string, session Session, ctx context.Context, box interface{}) {
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

func (d *Dispatch) triggerCallbacks(evtName string, session Session, ctx context.Context, box interface{}) {
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

	// remove the run only once listeners
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	for _, index := range d.listenOnceOnly[evtName] {
		// https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
		d.listeners[evtName][index] = d.listeners[evtName][len(d.listeners[evtName])-1]
		d.listeners[evtName][len(d.listeners[evtName])-1] = nil
		d.listeners[evtName] = d.listeners[evtName][:len(d.listeners[evtName])-1]
	}

	// remove the once only register
	_, exists := d.listenOnceOnly[evtName]
	if exists {
		delete(d.listenOnceOnly, evtName)
	}
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
func (d *Dispatch) GuildMemberAddChan() <-chan *GuildMemberAddBox {
	return d.guildMemberAddChan
}
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

func (d *Dispatch) AddHandler(evtName string, listener interface{}) {
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	d.listeners[evtName] = append(d.listeners[evtName], listener)
}

func (d *Dispatch) AddHandlerOnce(evtName string, listener interface{}) {
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	index := len(d.listeners[evtName])
	d.listeners[evtName] = append(d.listeners[evtName], listener)
	d.listenOnceOnly[evtName] = append(d.listenOnceOnly[evtName], index)
}

// wtf is this
func Unmarshal(data []byte, box interface{}) {
	err := json.Unmarshal(data, box)
	if err != nil {
		panic(err)
	}
}
