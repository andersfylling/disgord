package disgord

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/disgord/websocket"
)

//////////////////////////////////////////////////////
//
// Reactor: Consists of a basic demultiplexer, dispatcher and handlerSpecification.
//
// HandlerSpecification can hold one or more handlers, zero or more middlewares, and one controller.
//
//////////////////////////////////////////////////////

//////////////////////////////////////////////////////
//
// Demultiplexer
//
//////////////////////////////////////////////////////

func demultiplexer(d *dispatcher, read <-chan *websocket.Event, cache *Cache) {
	for {
		var evt *websocket.Event
		var alive bool

		select {
		case evt, alive = <-read:
			if !alive {
				return
			}
		case <-d.shutdown:
			return
		}

		var resource eventBox
		if resource = defineResource(evt.Name); resource == nil {
			fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name, string(evt.Data))
			continue // move on to next event
		}

		// populate resource
		ctx := context.Background()
		resource.registerContext(ctx)

		// first unmarshal to get identifiers
		//tmp := *resource

		// unmarshal into cacheLink
		//err := c.cacheEvent2(evtName, resource)

		if err := httd.Unmarshal(evt.Data, resource); err != nil {
			d.session.Logger().Error(err)
			continue // ignore event
			// TODO: if an event is ignored, should it not at least send a signal for listeners with no parameters?
		}
		executeInternalUpdater(evt)

		// TODO: updating internal states should be independent of the public reactor?
		//  But should the public handlers wait to be triggered until all the internals are updated?
		//executeInternalClientUpdater(c, evt)

		// cache
		if cache != nil {
			cacheEvent(cache, evt.Name, resource, evt.Data)
		}

		go d.dispatch(ctx, evt.Name, resource)
	}
}

//////////////////////////////////////////////////////
//
// Dispatcher
//
//////////////////////////////////////////////////////

// dispatcher holds all the channels and internal state for all registered
// handlers
type dispatcher struct {
	sync.RWMutex

	*dispatcherChans
	activateEventChannels bool

	// an event can have one or more handlers
	handlerSpecs map[string][]*handlerSpec

	// use session to allow mocking the client instance later on
	session  Session
	shutdown chan struct{}
}

func (d *dispatcher) addSessionInstance(s Session) {
	d.session = s
	if d.activateEventChannels {
		d.dispatcherChans.session = s
	}
}

// register registers handlers.
// Note! While the dispatcher handles registration in form of a method,
// deregistration is done automatically by checking the controller spec after each dispatch.
// See HandlerCtrl.
func (d *dispatcher) register(evt string, inputs ...interface{}) error {
	// detect middleware then handlers. Ordering is important.
	spec := &handlerSpec{}
	if err := spec.populate(inputs...); err != nil { // TODO: improve redundant checking
		return err // if the pattern is wrong: (event,[ ...middlewares,] ...handlers[, controller])
		// if you want to error check before you use the .On, you can use disgord.ValidateHandlerInputs(...)
	}

	// tip to users: Try to remember the handlers won't be added until the
	//  OnInsert(..) exits.
	err := spec.ctrl.OnInsert(d.session)
	if err != nil {
		d.session.Logger().Error(err)
	}

	d.Lock()
	d.handlerSpecs[evt] = append(d.handlerSpecs[evt], spec)
	d.Unlock()

	return nil
}

func (d *dispatcher) dispatch(ctx context.Context, evtName string, evt resource) {
	// channels
	if d.activateEventChannels {
		go d.dispatcherChans.trigger(ctx, evtName, evt)
	}

	// handlers
	d.RLock()
	specs := d.handlerSpecs[evtName]
	d.RUnlock()

	dead := make([]*handlerSpec, 0)

	for i := range specs {
		if alive := specs[i].next(); !alive {
			dead = append(dead, specs[i])
			continue
		}

		localEvt := specs[i].runMdlws(evt)
		if localEvt == nil {
			continue
		}

		for _, handler := range specs[i].handlers {
			d.trigger(handler, evt)
		}
	}

	// time to remove the dead
	if len(dead) == 0 {
		return
	}

	d.Lock()

	// make sure the dead has not already been removed, after all this is concurrent
	specs = d.handlerSpecs[evtName]
	for _, deadspec := range dead {
		for i, spec := range specs {
			if spec == deadspec { // compare pointers
				// delete the dead spec, but keep the ordering
				copy(specs[i:], specs[i+1:])
				specs[len(specs)-1] = nil // GC
				specs = specs[:len(specs)-1]
				break
			}
		}
	}

	// update the specs
	d.handlerSpecs[evtName] = specs
	d.Unlock()

	// notify specs
	go func(dead []*handlerSpec) {
		for i := range dead {
			if err := dead[i].ctrl.OnRemove(d.session); err != nil {
				d.session.Logger().Error(err)
			}
		}
	}(dead)
}

//////////////////////////////////////////////////////
//
// Handler logic
//
//////////////////////////////////////////////////////

// HandlerCtrl used when inserting a handler to dictate whether or not the handler(s) should
// still be kept in the handlers list..
type HandlerCtrl interface {
	OnInsert(Session) error
	OnRemove(Session) error

	// IsDead does not need to be locked as the demultiplexer access it synchronously.
	IsDead() bool

	// Update For every time Update is called, it's internal trackers must be updated.
	// you should assume that .Update() means the handler was used.
	Update()
}

// these "simple" handler can be used, if you don't care about the actual event data
type SimplestHandler = func()
type SimpleHandler = func(Session)

// Handler needs to match one of the *Handler signatures
type Handler = interface{}

// Middleware allows you to manipulate data during the "stream"
type Middleware = func(interface{}) interface{}

// handlerSpec (handler specification) holds the details for executing the handler(s)
// think about the code flow as a whole:
// eg. mdlw1 -> midlw2 -> handler1 -> handler2 -> ctrl
//
// If any of the middlewares manipulates the data, the next middlewares or handlers in the
// chain will get said manipulated data.
type handlerSpec struct {
	sync.RWMutex
	middlewares []Middleware
	handlers    []Handler
	ctrl        HandlerCtrl
}

func (hs *handlerSpec) next() bool {
	hs.Lock()
	defer hs.Unlock()

	if hs.ctrl.IsDead() {
		return false
	}

	hs.ctrl.Update()
	return true
}

// populate is essentially the constructor for a handlerSpec
func (hs *handlerSpec) populate(inputs ...interface{}) (err error) {
	var i int

	// middlewares
	for ; i < len(inputs); i++ {
		if mdlw, ok := inputs[i].(Middleware); ok {
			hs.middlewares = append(hs.middlewares, mdlw)
		} else {
			break
		}
	}

	// handlers
	for ; i < len(inputs)-1; i++ {
		if handler, ok := inputs[i].(Handler); ok {
			hs.handlers = append(hs.handlers, handler)
		} else {
			break
		}
	}

	// check if last arg is a controller
	if i < len(inputs) {
		if ctrl, ok := inputs[i].(HandlerCtrl); ok {
			hs.ctrl = ctrl
			i++
		} else if handler, ok := inputs[i].(Handler); ok {
			hs.handlers = append(hs.handlers, handler)
			hs.ctrl = eternalCtrl
			i++
		}
	}

	if len(inputs) != i {
		format := "unable to add all handlers/middlewares (%d/%d). Are they in correct order? middlewares, then handlers"
		err = errors.New(fmt.Sprintf(format, i, len(inputs)))
	}

	return err
}

func (hs *handlerSpec) runMdlws(evt interface{}) interface{} {
	for i := range hs.middlewares {
		evt = hs.middlewares[i](evt) // note how the evt pointer is overwritten
		if evt == nil {
			break
		}
	}

	return evt
}

//////////////////////////////////////////////////////
//
// Handler Controller
//
//////////////////////////////////////////////////////

// Ctrl is a handler controller that supports lifetime and max number of execution for one or several handlers.
//  // register only the first 6 votes
//  client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Runs: 6})
//
//  // Allow voting for only 10 minutes
//  client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Duration: 10*time.Second})
//
//  // Allow voting until the month is over
//  client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Until: time.Now().AddDate(0, 1, 0)})
type Ctrl struct {
	Runs     int
	Until    time.Time
	Duration time.Duration
}

var _ HandlerCtrl = (*Ctrl)(nil)

func (c *Ctrl) OnInsert(Session) error {
	if c.Runs == 0 {
		c.Runs = -1
	}
	if c.Duration.Nanoseconds() > 0 {
		if c.Until.IsZero() {
			c.Until = time.Now()
		}
		c.Until = c.Until.Add(c.Duration)
	}
	if c.Until.IsZero() {
		snow := Snowflake(^uint64(0))
		c.Until = snow.Date() // until the snowflakes fall
	}

	return nil
}

func (c *Ctrl) OnRemove(Session) error {
	return nil
}

func (c *Ctrl) IsDead() bool {
	return c.Runs == 0 || time.Now().After(c.Until)
}

func (c *Ctrl) Update() {
	if c.Runs > 0 {
		c.Runs--
	}
}

//////////////////////////////////////////////////////
//
// Custom Controllers
//
//////////////////////////////////////////////////////

// eternalHandlersCtrl is used for handlers without a defined controller. Letting them live forever.
type eternalHandlersCtrl struct {
	Ctrl
}

var _ HandlerCtrl = (*eternalHandlersCtrl)(nil)

func (c *eternalHandlersCtrl) Update()      {}
func (c *eternalHandlersCtrl) IsDead() bool { return false }

// reused by handlers that have no ctrl defined
var eternalCtrl = &eternalHandlersCtrl{}
