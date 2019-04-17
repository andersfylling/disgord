package disgord

//
//type cache2 struct {
//	users    usersCache
//	channels channelsCache
//}
//
//
////////////////////////////////////////////////////////
////
//// USERS
////
////////////////////////////////////////////////////////
//
//type /**/ usersCache struct {
//	internal interfaces.CacheAlger
//	config   *CacheConfig
//}
//
//func (c *usersCache) initUser() *User {
//	// TODO: lock free pool. Especially for EvtGuildMembersChunk
//	return &User{}
//}
//
//// handleGatewayEvent is responsible for updating/adding a user to the cache.
//// Note! the data param must be a user only(!) json. It is expected that this is parsed and handled
//// before calling handleGatewayEvent.
//func (c *usersCache) handleGatewayEvent(evt string, data []byte, flags ...Flag) (updated interface{}, err error) {
//	// if the flag IgnoreCache is passed, the update object is not returned nor allocated & copied
//	ignoreUpdated := mergeFlags(flags).Ignorecache()
//
//	// event-type-check
//	// only some events contains user data
//	// the rest is ignored
//	eventsWithUserContent := [...]string{
//		EvtUserUpdate,
//
//		EvtReady,
//
//		EvtGuildCreate, EvtGuildBanAdd, EvtGuildBanRemove, EvtGuildMemberAdd, EvtGuildMemberRemove,
//		EvtGuildMemberUpdate, EvtGuildMembersChunk,
//
//		EvtChannelCreate, EvtChannelUpdate, EvtChannelDelete,
//
//		EvtPresenceUpdate,
//
//		EvtVoiceStateUpdate,
//
//		EvtMessageCreate, EvtMessageUpdate,
//	}
//
//	var ok bool
//	for _, str := range eventsWithUserContent {
//		if ok = evt == str; ok {
//			break
//		}
//	}
//	if !ok {
//		return nil, errCacheUnsupportedEvt
//	}
//	// end event-type-check
//
//	// skip-redundant-objects
//	// User object might be empty which should not lock the cache.
//	var nrOfKeys uint
//	if nrOfKeys = countNrOfObjectKeys(data); nrOfKeys == 0 {
//		return nil, errCacheJSONObjectTooSmall
//	}
//	// end skip-redundant-objects
//
//	// get-user-id
//	var id Snowflake
//	if id, err = getJSONObjectID(data, "id"); err != nil {
//		return nil, err
//	}
//	// end get-user-id
//
//	var usr *User
//	// no-content
//	// When only field/key is present in the json object
//	// it can only be the id. So there is no need to pretend to update an object
//	if nrOfKeys == 1 {
//		if !ignoreUpdated {
//			updated = c.Get(id)
//		}
//		return updated, nil
//	}
//	// end no-content
//
//	// update-user
//	if !ignoreUpdated {
//		updated = c.initUser()
//	}
//	c.internal.Lock()
//	defer c.internal.Unlock()
//	if item, exists := c.internal.Get(id); exists {
//		usr = item.Object().(*User)
//		err = Unmarshal(data, usr)
//	} else {
//		usr = c.initUser()
//		if err = Unmarshal(data, usr); err == nil {
//			c.internal.Set(id, c.internal.CreateCacheableItem(usr))
//		}
//	}
//
//	// let's be kind to those that ignores errors, lol.
//	if err != nil {
//		return nil, err
//	}
//	if !ignoreUpdated {
//		_ = usr.copyOverToCache(updated) // double locking
//	}
//	// end update-user
//
//	return updated, nil
//}
//func (c *usersCache) handleRESTResponse(obj interface{}) (err error) {
//	// don't checking if it's actually a user.
//	// panics here will only help us improve the data flow if this method was incorrectly used.
//	//user := obj.(*User)
//	//if user == nil {
//	//	return
//	//}
//	//
//	//c.internal.Lock()
//	//if item, exists := c.internal.Get(user.ID); exists {
//	//	err = user.copyOverToCache(item.Object().(*User))
//	//} else {
//	//	c.internal.Set(user.ID, c.internal.CreateCacheableItem(user))
//	//}
//	//c.internal.Unlock()
//
//	return err
//}
//func (c *usersCache) Del(userID Snowflake) {
//	c.internal.Lock()
//	c.internal.Delete(userID)
//	c.internal.Unlock()
//}
//func (c *usersCache) Get(userID Snowflake) (user interface{}) {
//	c.internal.RLock()
//	if item, exists := c.internal.Get(userID); exists {
//		user = item.Object().(*User).DeepCopy() // double lock
//	}
//	c.internal.RUnlock()
//
//	return user
//}
//func (c *usersCache) Size() (size uint) {
//	c.internal.RLock()
//	size = c.internal.Size()
//	c.internal.RUnlock()
//
//	return
//}
//func (c *usersCache) Cap() (cap uint) {
//	c.internal.RLock()
//	cap = c.internal.Cap()
//	c.internal.RUnlock()
//
//	return
//}
//func (c *usersCache) ListIDs() (list []snowflake.ID) {
//	c.internal.RLock()
//	list = c.internal.ListIDs()
//	c.internal.RUnlock()
//
//	return
//}
//
//// Foreach allows you iterate over the users. This is not blocking for the rest of the system
//// as it blocks only when it copies or extract data from one user at the time.
//// This is faster when you make the cache mutable, but then again that introduces higher
//// risk are then involved (race conditions, incorrect cache, etc).
//func (c *usersCache) Foreach(cb func(*User)) {
//	ids := c.ListIDs()
//
//	for i := range ids {
//		if user := c.Get(ids[i]); user != nil {
//			cb(user.(*User))
//		}
//	}
//}
//
//var _ gatewayCacher = (*usersCache)(nil)
//var _ restCacher = (*usersCache)(nil)
//var _ BasicCacheRepo = (*usersCache)(nil)
//
////////////////////////////////////////////////////////
////
//// PRESENCE
////
////////////////////////////////////////////////////////
//
//type cachedGuildPresences struct {
//	sync.RWMutex
//	presences []*UserPresence
//}
//
//func (c *cachedGuildPresences) add(data []byte) {
//	p := &UserPresence{}
//	if err := Unmarshal(data, p); err != nil {
//		return
//	}
//
//	c.Lock()
//	c.presences = append(c.presences, p)
//	c.Unlock()
//}
//func (c *cachedGuildPresences) del(userID Snowflake) {
//	// keep deleting nil, incorrect entries
//	// TODO: move all entries to the end instead and then resize - check with benchmark
//	for {
//		i := -1
//		c.Lock()
//		for j := range c.presences {
//			if c.presences[j] == nil || c.presences[j].User == nil || c.presences[j].User.ID == userID {
//				i = j
//				break
//			}
//		}
//		if i == -1 {
//			c.Unlock()
//			break
//		}
//
//		// remove entry
//		c.presences[i] = c.presences[len(c.presences)-1]
//		c.presences[len(c.presences)-1] = nil
//		c.presences = c.presences[:len(c.presences)-1]
//		c.Unlock()
//	}
//}
//func (c *cachedGuildPresences) update(data []byte) {
//	var userID Snowflake
//	var err error
//	if userID, err = getJSONObjectID(data, "user", "id"); err != nil {
//		return
//	}
//
//	var p *UserPresence
//	c.Lock()
//	for j := range c.presences {
//		if c.presences[j].User.ID == userID {
//			p = c.presences[j]
//			break
//		}
//	}
//	// if the presence/user has not yet been added, we skip the update step
//	if p == nil {
//		c.Unlock()
//		c.add(data)
//	}
//
//	// TODO: presence.Activities, can this cause issues?
//	_ = Unmarshal(data, p)
//	c.Unlock()
//}
//
//// the guild should nil their presence field, and fetch them from here on build
//type /**/ presencesCache struct {
//	internal interfaces.CacheAlger
//	users    *usersCache
//	config   *CacheConfig
//}
//
//// handleGatewayEvent is responsible for updating/adding a presence to the cache.
//// Note! the data param must be a presence only(!) json. It is expected that this is parsed and handled
//// before calling handleGatewayEvent.
//func (c *presencesCache) handleGatewayEvent(evt string, data []byte, flags ...Flag) (updated interface{}, err error) {
//	// event-type-check
//	// only some events contains presence data,
//	// the rest is ignored
//	eventsWithContent := [...]string{
//		EvtGuildCreate,
//		EvtPresenceUpdate,
//	}
//
//	var ok bool
//	for _, str := range eventsWithContent {
//		if ok = evt == str; ok {
//			break
//		}
//	}
//	if !ok {
//		return nil, errCacheUnsupportedEvt
//	}
//	// end event-type-check
//
//	// skip-redundant-objects
//	// User object might be empty which should not lock the cache.
//	var nrOfKeys uint
//	if nrOfKeys = countNrOfObjectKeys(data); nrOfKeys < 2 { //user.id+guild_id
//		return nil, errCacheJSONObjectTooSmall
//	}
//	// end skip-redundant-objects
//
//	// get-guild-id
//	var guildID Snowflake
//	if evt == EvtGuildDelete {
//		if guildID, err = getJSONObjectID(data, "id"); err != nil {
//			return nil, err
//		}
//	} else {
//		if guildID, err = getJSONObjectID(data, "guild_id"); err != nil {
//			return nil, err
//		}
//	}
//	// end get-guild-id
//
//	// if the flag IgnoreCache is passed, the update object is not returned nor allocated & copied
//	ignoreUpdated := mergeFlags(flags).Ignorecache()
//
//	// update-user
//	var presences *cachedGuildPresences
//	c.internal.Lock()
//	defer c.internal.Unlock()
//	if item, exists := c.internal.Get(guildID); exists {
//		presences = item.Object().(*cachedGuildPresences)
//		if evt == EvtPresenceUpdate {
//			presences.update(data)
//		} else if evt == EvtGuildMemberRemove {
//			var userID Snowflake
//			if userID, err = getJSONObjectID(data, "user", "id"); err != nil {
//				return nil, err
//			}
//			presences.del(userID)
//		} else if evt == EvtGuildDelete {
//			for i := range presences.presences {
//				presences.presences[i] = nil
//			}
//			c.internal.Delete(guildID)
//		}
//	} else {
//		if evt == EvtPresenceUpdate {
//			presences = &cachedGuildPresences{}
//			presences.add(data)
//			c.internal.Set(guildID, c.internal.CreateCacheableItem(presences))
//		}
//	}
//
//	// let's be kind to those that ignores errors, lol.
//	if err != nil {
//		return nil, err
//	}
//	if presences == nil {
//		return nil, nil
//	}
//	if !ignoreUpdated {
//		updated := make([]*UserPresence, 0, len(presences.presences))
//		for i := range presences.presences {
//			updated = append(updated, presences.presences[i].DeepCopy().(*UserPresence))
//		}
//	}
//	// end update-user
//
//	return updated, nil
//}
//func (c *presencesCache) handleRESTResponse(obj interface{}) (err error) {
//	return nil
//}
//func (c *presencesCache) Del(userID Snowflake) {
//	c.internal.Lock()
//	c.internal.Delete(userID)
//	c.internal.Unlock()
//}
//func (c *presencesCache) Get(userID Snowflake) (user interface{}) {
//	c.internal.RLock()
//	if item, exists := c.internal.Get(userID); exists {
//		user = item.Object().(*User).DeepCopy() // double lock
//	}
//	c.internal.RUnlock()
//
//	return user
//}
//func (c *presencesCache) Size() (size uint) {
//	c.internal.RLock()
//	size = c.internal.Size()
//	c.internal.RUnlock()
//
//	return
//}
//func (c *presencesCache) Cap() (cap uint) {
//	c.internal.RLock()
//	cap = c.internal.Cap()
//	c.internal.RUnlock()
//
//	return
//}
//func (c *presencesCache) ListIDs() (list []snowflake.ID) {
//	c.internal.RLock()
//	list = c.internal.ListIDs()
//	c.internal.RUnlock()
//
//	return
//}
//
//var _ gatewayCacher = (*presencesCache)(nil)
//var _ restCacher = (*presencesCache)(nil)
//var _ BasicCacheRepo = (*presencesCache)(nil)
//
////////////////////////////////////////////////////////
////
//// CHANNELS
////
////////////////////////////////////////////////////////
//
//type /**/ channelsCache struct {
//	internal interfaces.CacheAlger
//	users    *usersCache
//	config   *CacheConfig
//}
//
//func (c *channelsCache) handleGatewayEvent(evt string, data []byte) (updated interface{}, err error) {
//	idBytes, _, _, err := jp.Get(data, "id")
//	if err != nil {
//		return nil, err
//	}
//
//	var id Snowflake
//	if err = id.UnmarshalJSON(idBytes); err != nil {
//		return nil, err
//	}
//
//	// update users first
//	var participants []*User
//	const recipientsKey = "recipient"
//	err = jp.ObjectEach(data, func(key []byte, value []byte, dataType jp.ValueType, offset int) error {
//		if usr, err := c.users.handleGatewayEvent(evt, value); err == nil {
//			participants = append(participants, usr.(*User))
//		}
//		return nil
//	}, recipientsKey)
//	if len(participants) > 0 {
//		data = jp.Delete(data, recipientsKey)
//	}
//
//	alias := &Channel{}
//	c.internal.Lock()
//	var peek *Channel
//	if item, exists := c.internal.Get(id); exists {
//		peek = item.Object().(*Channel)
//		err = Unmarshal(data, peek)
//	} else {
//		peek = &Channel{}
//		if err = Unmarshal(data, peek); err == nil {
//			c.internal.Set(id, c.internal.CreateCacheableItem(peek))
//		}
//	}
//
//	// participants
//	// always reset, as the channel type may change.. discord.
//	peek.recipientsIDs = nil
//	if len(participants) > 0 {
//		peek.recipientsIDs = make([]Snowflake, 0, len(participants))
//		for i := range participants {
//			peek.recipientsIDs = append(peek.recipientsIDs, participants[i].ID)
//		}
//	}
//
//	// create a copy to be returned
//	*alias = *peek
//	peek = nil // avoid further changes
//	c.internal.Unlock()
//
//	// add participants/recipients
//	alias.recipientsIDs = nil
//	alias.Recipients = participants
//
//	// let's be kind to those that ignores errors, lol.
//	if err != nil {
//		return nil, err
//	}
//
//	return alias, nil
//}
//func (c *channelsCache) handleRESTResponse(obj interface{}) (err error) {
//
//	return err
//}
//func (c *channelsCache) Delete(channelID Snowflake) {
//	c.internal.Lock()
//	defer c.internal.Unlock()
//	c.internal.Delete(channelID)
//}
//func (c *channelsCache) Get(id Snowflake) (channel interface{}) {
//	c.internal.RLock()
//	defer c.internal.RUnlock()
//	if item, exists := c.internal.Get(id); exists {
//		channel = item.Object().(*Channel).DeepCopy()
//	}
//
//	return channel
//}
//func (c *channelsCache) Size() uint {
//	c.internal.RLock()
//	defer c.internal.RUnlock()
//	return c.internal.Size()
//}
//func (c *channelsCache) Cap() uint {
//	c.internal.RLock()
//	defer c.internal.RUnlock()
//	return c.internal.Cap()
//}
//func (c *channelsCache) ListIDs() []Snowflake {
//	c.internal.RLock()
//	defer c.internal.RUnlock()
//	return c.internal.ListIDs()
//}
//
//// Foreach allows you iterate over the users. This is not blocking for the rest of the system
//// as it blocks only when it copies or extract data from one user at the time.
//// This is faster when you make the cache mutable, but then again that introduces higher
//// risk are then involved (race conditions, incorrect cache, etc).
//func (c *channelsCache) Foreach(cb func(*Channel)) {
//	ids := c.ListIDs()
//
//	for i := range ids {
//		channel := c.Get(ids[i])
//		if channel != nil {
//			cb(channel.(*Channel))
//		}
//	}
//}
//
//var _ gatewayCacher = (*channelsCache)(nil)
//var _ restCacher = (*channelsCache)(nil)
//var _ BasicCacheRepo = (*channelsCache)(nil)
