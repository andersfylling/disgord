package disgord

import (
	"fmt"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/gateway"
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

func (c *Client) demultiplexer(d *dispatcher, read <-chan *gateway.Event) {
	for {
		var evt *gateway.Event
		var alive bool

		select {
		case evt, alive = <-read:
			if !alive {
				return
			}
		case <-d.shutdown:
			return
		}

		// var resource evtResource
		// if resource = defineResource(evt.Name); resource == nil {
		// 	fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name, string(evt.Data))
		// 	continue // move on to next event
		// }

		resourceI, err := cacheDispatcher(c.cache, evt.Name, evt.Data)
		if resourceI == nil {
			err = fmt.Errorf("cache did not instantiate object. Prev error: %w", err)
		}
		if err != nil {
			// skip unknown events
			var knownEvent bool
			evts := AllEvents()
			for _, e := range evts {
				if knownEvent = e == evt.Name; knownEvent {
					break
				}
			}
			if !knownEvent {
				continue
			}

			err = fmt.Errorf("demultiplexer{%s}: %w, data '%s'", evt.Name, err, string(evt.Data))
			d.session.Logger().Error(err)
			continue
		}
		resource := resourceI.(evtResource)
		resource.setShardID(evt.ShardID)

		go d.dispatch(evt.Name, resource)
	}
}

//////////////////////////////////////////////////////
//
// Dispatcher
//
//////////////////////////////////////////////////////

// dispatcher holds all the Channels and internal state for all registered
// handlers
type dispatcher struct {
	sync.RWMutex

	// an event can have one or more handlers
	handlerSpecs map[string][]*handlerSpec

	// use session to allow mocking the Client instance later on
	session  Session
	shutdown chan struct{}
}

func (d *dispatcher) addSessionInstance(s Session) {
	d.session = s
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

	// tip to Users: Try to remember the handlers won't be added until the
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

func (d *dispatcher) dispatch(evtName string, evt resource) {
	// handlers
	d.RLock()
	specs := d.handlerSpecs[evtName]
	d.RUnlock()

	dead := make([]*handlerSpec, 0)

	for _, spec := range specs {
		// faster. But somewhat weird to check death before running the handler
		// this can be used if we find a different way to write the Client.Ready
		// logic.
		//if alive := spec.next(); !alive {
		//	dead = append(dead, spec)
		//	continue
		//}
		spec.Lock()
		if dead := spec.ctrl.IsDead(); !dead {
			localEvt := spec.runMdlws(evt)
			if localEvt == nil {
				spec.Unlock()
				continue
			}

			for _, handler := range spec.handlers {
				d.trigger(handler, localEvt)
			}

			spec.ctrl.Update()
		}

		if spec.ctrl.IsDead() {
			dead = append(dead, spec)
		}
		spec.Unlock()
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
				//specs[len(specs)-1] = nil // GC, setting entries to nil requires locking
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
		if handler, ok := inputs[i].(Handler); ok && isHandler(handler) {
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
		err = fmt.Errorf(format, i, len(inputs))
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
//
//	// register only the first 6 votes
//	Client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Runs: 6})
//
//	// Allow voting for only 10 minutes
//	Client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Duration: 10*time.Second})
//
//	// Allow voting until the month is over
//	Client.On("MESSAGE_CREATE", filter.NonVotes, registerVoteHandler, &disgord.Ctrl{Until: time.Now().AddDate(0, 1, 0)})
type Ctrl struct {
	Runs     int
	Until    time.Time
	Duration time.Duration
	Channel  interface{}
}

var _ HandlerCtrl = (*Ctrl)(nil)

func (c *Ctrl) OnInsert(Session) error {
	if c.Channel != nil && !isHandler(c.Channel) {
		panic("Ctrl.Channel is not a valid Disgord event channel")
	}
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

// CloseChannel must be called instead of closing an event channel directly.
// This is to make sure Disgord does not go into a deadlock
func (c *Ctrl) CloseChannel() {
	c.Runs = 0
	closeChannel(c.Channel)
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

// rdyCtrl is used to trigger notify the user when all the websocket sessions have received their first READY event
type rdyCtrl struct {
	sync.Mutex
	shardReady    []bool
	localShardIDs []uint
	cb            func()
}

var _ HandlerCtrl = (*rdyCtrl)(nil)

func (c *rdyCtrl) OnInsert(s Session) error {
	return nil
}

func (c *rdyCtrl) OnRemove(s Session) error {
	c.Lock()
	defer c.Unlock()

	if c.cb != nil {
		go c.cb()
		c.cb = nil // such that it is only called once. See Client.GuildsReady(...)
	}
	return nil
}

func (c *rdyCtrl) IsDead() bool {
	c.Lock()
	defer c.Unlock()

	for _, id := range c.localShardIDs {
		if !c.shardReady[id] {
			return false
		}
	}

	return true
}

func (c *rdyCtrl) Update() {
	// handled in the handler
}

type guildsRdyCtrl struct {
	rdyCtrl
	status map[Snowflake]bool
}

func (c *guildsRdyCtrl) IsDead() bool {
	c.Lock()
	defer c.Unlock()

	for _, ready := range c.status {
		if !ready {
			return false
		}
	}

	return true
}
