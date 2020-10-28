package disgord

import (
	"context"

	"github.com/andersfylling/disgord/internal/disgorderr"
	"github.com/andersfylling/disgord/internal/gateway"
)

// Deprecated: use .Gateway() instead
func (c *Client) Event() SocketHandlerRegistrator {
	panic("deprecated")
}

// Link allows basic Discord connection control. Affects all shards
// Deprecated
type Link interface {
	// Connect establishes a websocket connection to the discord API
	Connect(ctx context.Context) error

	// Disconnect closes the discord websocket connection
	Disconnect() error
}

// SocketHandler all socket related logic
// Deprecated: use Gateway()
type SocketHandler interface {
	// Link controls the connection to the Discord API. Affects all shards.
	// Link

	// Disconnect closes the discord websocket connection
	Disconnect() error

	// Suspend temporary closes the socket connection, allowing resources to be
	// reused on reconnect
	Suspend() error

	OnSocketEventer

	// Event gives access to type safe event handler registration using the builder pattern
	Event() SocketHandlerRegistrator

	Emitter
}

// Deprecated
type OnSocketEventer interface {
	// On creates a specification to be executed on the given event. The specification
	// consists of, in order, 0 or more middlewares, 1 or more handlers, 0 or 1 controller.
	// On incorrect ordering, or types, the method will panic. See reactor.go for types.
	//
	// Each of the three sub-types of a specification is run in sequence, as well as the specifications
	// registered for a event. However, the slice of specifications are executed in a goroutine to avoid
	// blocking future events. The middlewares allows manipulating the event data before it reaches the
	// handlers. The handlers executes short-running logic based on the event data (use go routine if
	// you need a long running task). The controller dictates lifetime of the specification.
	//
	//  // a handler that is executed on every Ready event
	//  Client.On(EvtReady, onReady)
	//
	//  // a handler that runs only the first three times a READY event is fired
	//  Client.On(EvtReady, onReady, &Ctrl{Runs: 3})
	//
	//  // a handler that only runs for events within the first 10 minutes
	//  Client.On(EvtReady, onReady, &Ctrl{Duration: 10*time.Minute})
	On(event string, inputs ...interface{})
}

// On creates a specification to be executed on the given event. The specification
// consists of, in order, 0 or more middlewares, 1 or more handlers, 0 or 1 controller.
// On incorrect ordering, or types, the method will panic. See reactor.go for types.
//
// Each of the three sub-types of a specification is run in sequence, as well as the specifications
// registered for a event. However, the slice of specifications are executed in a goroutine to avoid
// blocking future events. The middlewares allows manipulating the event data before it reaches the
// handlers. The handlers executes short-running logic based on the event data (use go routine if
// you need a long running task). The controller dictates lifetime of the specification.
//
//  // a handler that is executed on every Ready event
//  Client.On(EvtReady, onReady)
//
//  // a handler that runs only the first three times a READY event is fired
//  Client.On(EvtReady, onReady, &Ctrl{Runs: 3})
//
//  // a handler that only runs for events within the first 10 minutes
//  Client.On(EvtReady, onReady, &Ctrl{Duration: 10*time.Minute})
//
// Another example is to create a voting system where you specify a deadline instead of a Runs counter:
//  On("MESSAGE_CREATE", mdlwHasMentions, handleMsgsWithMentions, saveVoteToDB, &Ctrl{Until:time.Now().Add(time.Hour)})
//
// You can use your own Ctrl struct, as long as it implements disgord.HandlerCtrl. Do not execute long running tasks
// in the methods. Use a go routine instead.
//
// If the HandlerCtrl.OnInsert returns an error, the related handlers are still added to the dispatcher.
// But the error is logged to the injected logger instance (log.Error).
//
// This ctrl feature was inspired by https://github.com/discordjs/discord.js
// Deprecated: use Client.Gateway().... which also provides compile time validity
func (c *Client) On(event string, inputs ...interface{}) {
	if err := ValidateHandlerInputs(inputs...); err != nil {
		panic(err)
	}

	if err := c.dispatcher.register(event, inputs...); err != nil {
		panic(err)
	}
}

// Deprecated
func ValidateHandlerInputs(inputs ...interface{}) (err error) {
	var i int
	var ok bool

	// make sure that middlewares are only at beginning
	for j := i; j < len(inputs); j++ {
		if _, ok = inputs[j].(Middleware); ok {
			if j != i {
				return disgorderr.NewHandlerSpecErr(
					disgorderr.HandlerSpecErrCodeUnexpectedMiddleware,
					"middlewares can only be in the beginning. Grouped together")
			}
			i++
		}
	}

	// there should now be N handlers, 0 < N.
	if len(inputs) <= i {
		return disgorderr.NewHandlerSpecErr(
			disgorderr.HandlerSpecErrCodeMissingHandler, "missing handler(s)")
	}

	for j := i; j < len(inputs); j++ {
		if _, ok = inputs[j].(HandlerCtrl); ok {
			// first element after middlewares and last in inputs
			if j == i && len(inputs)-1 == j {
				return disgorderr.NewHandlerSpecErr(
					disgorderr.HandlerSpecErrCodeMissingHandler, "missing handler(s)")
			}
			// not last
			if len(inputs)-1 != j {
				return disgorderr.NewHandlerSpecErr(
					disgorderr.HandlerSpecErrCodeUnexpectedCtrl,
					"a handlerCtrl's can only be at the end of the definition and only one")
			}
			break
		}
		if _, ok = inputs[j].(Ctrl); ok {
			return disgorderr.NewHandlerSpecErr(
				disgorderr.HandlerSpecErrCodeNotHandlerCtrlImpl,
				"does not implement disgord.HandlerCtrl. Try to use &disgord.Ctrl instead of disgord.Ctrl")
		}

		if !isHandler(inputs[j]) {
			return disgorderr.NewHandlerSpecErr(
				disgorderr.HandlerSpecErrCodeUnknownHandlerSignature,
				"invalid handler signature. General tip: no handlers can use the param type `*disgord.Session`, try `disgord.Session` instead")
		}
	}

	return nil
}

// Connect establishes a websocket connection to the discord API
// Deprecated: use .Gateway().Connect() instead
func (c *Client) Connect(ctx context.Context) (err error) {
	return c.Gateway().WithContext(ctx).Connect()
}

// Disconnect closes the discord websocket connection
// Deprecated: use .Gateway().Disconnect() instead
func (c *Client) Disconnect() (err error) {
	return c.Gateway().Disconnect()
}

// Suspend in case you want to temporary disconnect from the Gateway. But plan on
// connecting again without restarting your software/application, this should be used.
// Deprecated: will be removed, feel free to request this feature again through a PR
func (c *Client) Suspend() (err error) {
	c.log.Info("Closing Discord gateway connection")
	if err = c.shardManager.Disconnect(); err != nil {
		return err
	}
	c.log.Info("Suspended")

	return nil
}

// DisconnectOnInterrupt wait until a termination signal is detected
// Deprecated: use .Gateway().DisconnectOnInterrupt() instead
func (c *Client) DisconnectOnInterrupt() (err error) {
	return c.Gateway().DisconnectOnInterrupt()
}

// Deprecated: use .Gateway().StayConnectedUntilInterrupted() instead
func (c *Client) StayConnectedUntilInterrupted() (err error) {
	return c.Gateway().StayConnectedUntilInterrupted()
}

// Ready triggers a given callback when all shards has gotten their first Ready event
// Warning: Do not call Client.Connect before this.
// Deprecated: use .Gateway().BotReady instead
func (c *Client) Ready(cb func()) {
	c.Gateway().BotReady(cb)
}

// GuildsReady is triggered once all unavailable Guilds given in the READY event has loaded from their respective GUILD_CREATE events.
// Deprecated: use .Gateway().BotGuildsReady instead
func (c *Client) GuildsReady(cb func()) {
	c.Gateway().BotGuildsReady(cb)
}

// Emit sends a socket command directly to Discord.
// Deprecated: use .Gateway().Dispatch instead
func (c *Client) Emit(name gatewayCmdName, payload gateway.CmdPayload) (unchandledGuildIDs []Snowflake, err error) {
	return c.Gateway().Dispatch(name, payload)
}
