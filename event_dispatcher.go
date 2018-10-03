package disgord

// On ... TODO
func (d *Dispatch) On(event string, handlers ...interface{}) {
	d.ws.RegisterEvent(event)

	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	for _, handler := range handlers {
		d.listeners[event] = append(d.listeners[event], handler)
	}
}

// Once ... TODO
func (d *Dispatch) Once(event string, handlers ...interface{}) {
	d.ws.RegisterEvent(event) // TODO: remove event after firing. unless there are more handlers

	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()
	for _, handler := range handlers {
		index := len(d.listeners[event])
		d.listeners[event] = append(d.listeners[event], handler)
		d.listenOnceOnly[event] = append(d.listenOnceOnly[event], index)
	}
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

// AllChan sends all event types
func (d *Dispatch) All() <-chan interface{} {
	return d.allChan
}

func prepareBox(evtName string, box interface{}) {
	switch evtName {
	case EventGuildCreate:
		for _, role := range (box.(*GuildCreate)).Guild.Roles {
			role.guildID = (box.(*GuildCreate)).Guild.ID
		}
	case EventGuildUpdate:
		for _, role := range (box.(*GuildCreate)).Guild.Roles {
			role.guildID = (box.(*GuildCreate)).Guild.ID
		}
	case EventGuildRoleCreate:
		(box.(*GuildRoleCreate)).Role.guildID = (box.(*GuildCreate)).Guild.ID
	case EventGuildRoleUpdate:
		(box.(*GuildRoleUpdate)).Role.guildID = (box.(*GuildRoleUpdate)).GuildID
	}
}
