package disgord

import (
	"context"
	"fmt"

	. "github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/resource"
	"github.com/andersfylling/disgord/state"
	"github.com/sirupsen/logrus"
)

// eventHandler Takes a incoming event from the discordws package, parses it, and sends
// trigger requests to the event dispatcher and state cacher.
func (c *Client) eventHandler() {
	for {
		select {
		case evt, alive := <-c.socketEvtChan:
			if !alive {
				logrus.Error("Event channel is dead!")
				break
			}

			ctx := context.Background()
			evtName := evt.Name()
			session := c
			data := evt.Data()

			switch evtName {
			case ReadyKey:
				box := &ReadyBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)

				// cache
				for _, partialGuild := range box.Guilds {
					c.state.ProcessGuild(&state.GuildDetail{
						Guild:  resource.NewGuildFromUnavailable(partialGuild),
						Dirty:  true,
						Action: evtName,
					})
				}
				// TODO-caching: c.state.Myself()
			case ResumedKey:
				box := &ResumedBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case ChannelCreateKey, ChannelUpdateKey, ChannelDeleteKey:
				chanContent := &resource.Channel{}
				Unmarshal(data, chanContent)

				switch evtName { // internal switch statement for ChannelEvt
				case ChannelCreateKey:
					box := &ChannelCreateBox{Channel: chanContent, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case ChannelUpdateKey:
					box := &ChannelUpdateBox{Channel: chanContent, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case ChannelDeleteKey:
					box := &ChannelDeleteBox{Channel: chanContent, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				} // END internal switch statement for ChannelEvt

				// cache channel
				c.state.ProcessChannel(&state.ChannelDetail{
					Channel: chanContent,
					Dirty:   true,
					Action:  evtName,
				})
			case ChannelPinsUpdateKey:
				box := &ChannelPinsUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)

				// cache what?
			case GuildCreateKey, GuildUpdateKey, GuildDeleteKey:
				g := &resource.Guild{}
				Unmarshal(data, g)

				switch evtName { // internal switch statement for guild events
				case GuildCreateKey:
					box := &GuildCreateBox{Guild: g, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case GuildUpdateKey:
					box := &GuildUpdateBox{Guild: g, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case GuildDeleteKey:
					unavailGuild := resource.NewGuildUnavailable(g.ID)
					box := &GuildDeleteBox{UnavailableGuild: unavailGuild, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				} // END internal switch statement for guild events

				// cache
				// TODO-caching: channels, users on guild create / update
				c.state.ProcessGuild(&state.GuildDetail{
					Guild:  g,
					Dirty:  true,
					Action: evtName,
				})
			case GuildBanAddKey:
				box := &GuildBanAddBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)

				// cache
				c.state.ProcessUser(&state.UserDetail{
					User:  box.User,
					Dirty: true,
				})
			case GuildBanRemoveKey:
				box := &GuildBanRemoveBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case GuildEmojisUpdateKey:
				box := &GuildEmojisUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)

				// TODO-caching: emoji
			case GuildIntegrationsUpdateKey:
				box := &GuildIntegrationsUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case GuildMemberAddKey:
				box := &GuildMemberAddBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)

				// TODO-caching: caching members
			case GuildMemberRemoveKey:
				box := &GuildMemberRemoveBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				// TODO-caching: remove cached members
			case GuildMemberUpdateKey:
				box := &GuildMemberUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				// TODO-caching: update a member
			case GuildMembersChunkKey:
				box := &GuildMembersChunkBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				// TODO-caching: member chunk.. ?
			case GuildRoleCreateKey:
				box := &GuildRoleCreateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				//TODO-caching: guild role add
			case GuildRoleUpdateKey:
				box := &GuildRoleUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				//TODO-caching: guild role change
			case GuildRoleDeleteKey:
				box := &GuildRoleDeleteBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				//TODO-caching: remove guild role
			case MessageCreateKey, MessageUpdateKey, MessageDeleteKey:
				msg := resource.NewMessage()
				Unmarshal(data, msg)

				switch evtName { // internal switch statement for MessageEvt
				case MessageCreateKey:
					box := &MessageCreateBox{Message: msg, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case MessageUpdateKey:
					box := &MessageUpdateBox{Message: msg, Ctx: ctx}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				case MessageDeleteKey:
					box := &MessageDeleteBox{MessageID: msg.ID, ChannelID: msg.ChannelID}
					c.evtDispatch.triggerChan(evtName, session, ctx, box)
					c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				} // END internal switch statement for MessageEvt
			case MessageDeleteBulkKey:
				box := &MessageDeleteBulkBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case MessageReactionAddKey:
				box := &MessageReactionAddBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case MessageReactionRemoveKey:
				box := &MessageReactionRemoveBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case MessageReactionRemoveAllKey:
				box := &MessageReactionRemoveAllBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case PresenceUpdateKey:
				box := &PresenceUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case TypingStartKey:
				box := &TypingStartBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case UserUpdateKey:
				box := &UserUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
				//TODO-caching: user update, is this @me?
			case VoiceStateUpdateKey:
				box := &VoiceStateUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case VoiceServerUpdateKey:
				box := &VoiceServerUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			case WebhooksUpdateKey:
				box := &WebhooksUpdateBox{}
				box.Ctx = ctx
				Unmarshal(data, box)

				c.evtDispatch.triggerChan(evtName, session, ctx, box)
				c.evtDispatch.triggerCallbacks(evtName, session, ctx, box)
			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evtName, string(data))
			}
		}
	}
}
