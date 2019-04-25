package disgord

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
