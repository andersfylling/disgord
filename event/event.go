package event

import "github.com/andersfylling/disgord/guild"

// EventDispatcher is an application-level type for handling discord requests.
// All callbacks are optional, and whether they are defined or not
// is used to determine whether the EventDispatcher will send events to them.
type DispatcherInterface interface {
	// current EventHook fields here

	// OnEvent is called for all events.
	// Handlers must typecast the event type manually, and ensure
	// that it can handle receiving the same event twice if a type-specific
	// callback also exists.
	//OnEvent func(ctx *Context, ev event.DiscordEvent) error

	// OnMessageEvent is called for every message-related event.
	//OnMessageEvent func(ctx *Context, ev event.MessageEvent) error

	// OnConnectionEvent ...
	//OnUserEvent(eventName string, listener func(user *user.User))
	//OnMemberEvent(eventName string, listener func(member *guild.Member))
	// OnChannelEvent ...
	// OnGuildEvent ...

	OnEvent(eventName Type, listener interface{})
	Trigger(eventName Type, params ...interface{})
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		listeners: make(map[Type]([]interface{})),
	}
}

type Dispatcher struct {
	listeners map[Type]([]interface{})
}

func (d *Dispatcher) OnEvent(eventName Type, listener interface{}) {
	d.listeners[eventName] = append(d.listeners[eventName], listener)
}

func (d *Dispatcher) Trigger(eventName Type, params ...interface{}) {
	if val, ok := d.listeners[eventName]; ok {
		for _, listener := range val {
			switch l := listener.(type) {
			case func(guild *guild.Guild):
				l(params[0].(*guild.Guild))
			}
		}
	}
}
