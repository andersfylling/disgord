package disgord

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/disgordctx"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/guild"
	"github.com/sirupsen/logrus"
)

const (
	// APIComEncoding data encoding when communicating with the discord API
	APIComEncoding string = "json"

	// APIVersion desired API version to use
	APIVersion int = 6
)

type Config struct {
	Token      string
	HTTPClient *http.Client
	Debug      bool
}

// NewClient creates a new default disgord instance
func NewClient(conf *Config) (*Client, error) {

	if conf.HTTPClient == nil {
		// http client configuration
		conf.HTTPClient = &http.Client{
			Timeout: time.Second * 10,
		}
	}

	// Use discordws to keep the socket connection going
	dws, err := discordws.NewClient(&discordws.Config{
		// user settings
		Token:      conf.Token,
		HTTPClient: conf.HTTPClient,
		Debug:      conf.Debug,

		// lib specific
		DAPIVersion:  APIVersion,
		DAPIEncoding: APIComEncoding,
	})
	if err != nil {
		return nil, err
	}

	// create a disgord instance
	d := &Client{
		HTTPClient: conf.HTTPClient,
		ws:         dws,
		EventChan:  dws.GetEventChannel(),
		Token:      conf.Token,
		Dispatcher: event.NewDispatcher(),
		State:      NewStateCache(),
	}

	return d, nil
}

// NewRequiredClient same as NewDisgord, but exits program if an error occours
func NewRequiredClient(conf *Config) *Client {
	dg, err := NewClient(conf)
	if err != nil {
		logrus.Fatal(err)
	}

	return dg
}

type Client struct {
	sync.RWMutex

	Token string

	ws *discordws.Client

	HTTPClient *http.Client

	EventChan <-chan discordws.EventInterface

	// register listeners for events
	*event.Dispatcher

	// cache
	State StateCacher
}

func (c *Client) String() string {
	return c.ws.String()
}

func (c *Client) logInfo(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Info(msg)
}

func (c *Client) logErr(msg string) {
	logrus.WithFields(logrus.Fields{
		"lib": c.ws.String(),
	}).Error(msg)
}

// Connect establishes a websocket connection to the discord API
func (c *Client) Connect() (err error) {
	c.logInfo("Connecting to discord Gateway")
	err = c.ws.Connect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Connected")

	// setup event observer
	go c.eventObserver()

	return nil
}

// Disconnect closes the discord websocket connection
func (c *Client) Disconnect() (err error) {
	fmt.Println()
	c.logInfo("Closing Discord gateway connection")
	err = c.ws.Disconnect()
	if err != nil {
		c.logErr(err.Error())
		return
	}
	c.logInfo("Disconnected")

	return nil
}

func (c *Client) eventObserver() {
	for {
		select {
		case evt, ok := <-c.EventChan:
			if !ok {
				logrus.Error("Event channel is dead!")
				break
			}

			session := &disgordctx.Session{} //disgord context
			ctx := context.Background()

			// TODO: parsing JSON uses panic and not logging on issues..

			eventName := evt.Name()
			data := evt.Data()

			switch eventName {
			case event.ReadyKey:
				r := &event.ReadyBox{}
				r.Ctx = ctx
				event.Unmarshal(data, r)

				go c.ReadyEvent.Trigger(session, r)

				// allocate cache mem for guilds
				// for _, gu := range r.Guilds {
				// 	g := guild.NewGuildFromUnavailable(gu)
				// 	c.State.AddGuild(g) // checks if the ID already exists
				// }

				// updated myself
				//c.State.UpdateMySelf(r.User)

			case event.ResumedKey:
				resumed := &event.ResumedBox{}
				resumed.Ctx = ctx
				event.Unmarshal(data, resumed)

				// no need to handle this as its done at the socket level ...
				go c.ResumedEvent.Trigger(session, resumed)
			case event.ChannelCreateKey, event.ChannelUpdateKey, event.ChannelDeleteKey:
				chanContent := &channel.Channel{}
				event.Unmarshal(data, chanContent)

				switch eventName {
				case event.ChannelCreateKey:
					go c.ChannelCreateEvent.Trigger(session, &event.ChannelCreateBox{
						Channel: chanContent,
						Ctx:     ctx,
					})
				case event.ChannelUpdateKey:
					go c.ChannelUpdateEvent.Trigger(session, &event.ChannelUpdateBox{
						Channel: chanContent,
						Ctx:     ctx,
					})
				case event.ChannelDeleteKey:
					go c.ChannelDeleteEvent.Trigger(session, &event.ChannelDeleteBox{
						Channel: chanContent,
						Ctx:     ctx,
					})
				}
			case event.ChannelPinsUpdateKey:
				cpu := &event.ChannelPinsUpdateBox{}
				cpu.Ctx = ctx
				event.Unmarshal(data, cpu)
				go c.ChannelPinsUpdateEvent.Trigger(session, cpu)
			case event.GuildCreateKey, event.GuildUpdateKey, event.GuildDeleteKey:
				g := &guild.Guild{}
				event.Unmarshal(data, g)

				switch eventName { // internal switch statement for guild events
				case event.GuildCreateKey:
					// notifify listeners
					go c.GuildCreateEvent.Trigger(session, &event.GuildCreateBox{
						Guild: g,
						Ctx:   ctx,
					})
					// add to cache
					//c.State.AddGuild(g)
				case event.GuildUpdateKey:
					// notifify listeners
					go c.GuildUpdateEvent.Trigger(session, &event.GuildUpdateBox{
						Guild: g,
						Ctx:   ctx,
					})
					// update cache
					//c.State.UpdateGuild(g)
				case event.GuildDeleteKey:
					// notify listeners
					unavailGuild := guild.NewGuildUnavailable(g.ID)
					go c.GuildDeleteEvent.Trigger(session, &event.GuildDeleteBox{
						UnavailableGuild: unavailGuild,
						Ctx:              ctx,
					})
					//
					// cachedGuild, err := c.State.Guild(g.ID)
					// if err != nil {
					// 	// guild has not been cached earlier for some reason..
					// } else {
					// 	// update instance with complete info.
					// 	// Assumption: The cached version has no outdated information.
					// 	g = nil
					// 	g = cachedGuild
					// 	// delete the guild object from the cache
					// 	c.State.DeleteGuild(g)
					// }
				} // END internal switch statement for guild events
			case event.GuildBanAddKey:
				gba := &event.GuildBanAddBox{}
				gba.Ctx = ctx
				event.Unmarshal(data, gba)

				go c.GuildBanAddEvent.Trigger(session, gba)
			case event.GuildBanRemoveKey:
				gbr := &event.GuildBanRemoveBox{}
				gbr.Ctx = ctx
				event.Unmarshal(data, gbr)

				go c.GuildBanRemoveEvent.Trigger(session, gbr)
			case event.GuildEmojisUpdateKey:
				geu := &event.GuildEmojisUpdateBox{}
				geu.Ctx = ctx
				event.Unmarshal(data, geu)

				go c.GuildEmojisUpdateEvent.Trigger(session, geu)
			case event.GuildIntegrationsUpdateKey:
				giu := &event.GuildIntegrationsUpdateBox{}
				giu.Ctx = ctx
				event.Unmarshal(data, giu)

				go c.GuildIntegrationsUpdateEvent.Trigger(session, giu)
			case event.GuildMemberAddKey:
				gma := &event.GuildMemberAddBox{}
				gma.Ctx = ctx
				event.Unmarshal(data, gma)

				go c.GuildMemberAddEvent.Trigger(session, gma)
			case event.GuildMemberRemoveKey:
				gmr := &event.GuildMemberRemoveBox{}
				gmr.Ctx = ctx
				event.Unmarshal(data, gmr)

				go c.GuildMemberRemoveEvent.Trigger(session, gmr)
			case event.GuildMemberUpdateKey:
				gmu := &event.GuildMemberUpdateBox{}
				gmu.Ctx = ctx
				event.Unmarshal(data, gmu)

				go c.GuildMemberUpdateEvent.Trigger(session, gmu)
			case event.GuildMembersChunkKey:
				gmc := &event.GuildMembersChunkBox{}
				gmc.Ctx = ctx
				event.Unmarshal(data, gmc)

				go c.GuildMembersChunkEvent.Trigger(session, gmc)
			case event.GuildRoleCreateKey:
				r := &event.GuildRoleCreateBox{}
				r.Ctx = ctx
				event.Unmarshal(data, r)

				go c.GuildRoleCreateEvent.Trigger(session, r)

				// add to cache
				// g, err := c.State.Guild(r.GuildID)
				// if err != nil {
				// 	panic("you haven't correctly cached all guilds you fool!")
				// }
				//g.Lock()
				//g.AddRole(r.Role)
				//g.Unlock()
			case event.GuildRoleUpdateKey:
				r := &event.GuildRoleUpdateBox{}
				r.Ctx = ctx
				event.Unmarshal(data, r)

				go c.GuildRoleUpdateEvent.Trigger(session, r)
				// CACHING
				// g, err := c.State.Guild(r.GuildID)
				// if err != nil {
				// 	panic("you haven't correctly cached all guilds you fool!")
				// }
				//
				// // add to cache
				// g.Lock()
				// g.UpdateRole(r.Role)
				// g.Unlock()
			case event.GuildRoleDeleteKey:
				r := &event.GuildRoleDeleteBox{}
				r.Ctx = ctx
				event.Unmarshal(data, r)

				go c.GuildRoleDeleteEvent.Trigger(session, r)
				//
				// g, err := c.State.Guild(r.GuildID)
				// if err != nil {
				// 	panic("you haven't correctly cached all guilds you fool!")
				// }
				// g.Lock()
				// g.DeleteRoleByID(r.RoleID)
				// g.Unlock()
				// // TODO: remove role from guild members...
			case event.MessageCreateKey, event.MessageUpdateKey, event.MessageDeleteKey:
				msg := channel.NewMessage()
				event.Unmarshal(data, msg)

				// TODO: should i cache msg?..
				switch eventName {
				case event.MessageCreateKey:
					go c.MessageCreateEvent.Trigger(session, &event.MessageCreateBox{
						Message: msg,
						Ctx:     ctx,
					})
				case event.MessageUpdateKey:
					go c.MessageUpdateEvent.Trigger(session, &event.MessageUpdateBox{
						Message: msg,
						Ctx:     ctx,
					})
				case event.MessageDeleteKey:
					go c.MessageDeleteEvent.Trigger(session, &event.MessageDeleteBox{
						MessageID: msg.ID,
						ChannelID: msg.ChannelID,
					})
				}
			case event.MessageDeleteBulkKey:
				mdb := &event.MessageDeleteBulkBox{}
				mdb.Ctx = ctx
				event.Unmarshal(data, mdb)

				go c.MessageDeleteBulkEvent.Trigger(session, mdb)
			case event.MessageReactionAddKey:
				mra := &event.MessageReactionAddBox{}
				mra.Ctx = ctx
				event.Unmarshal(data, mra)

				go c.MessageReactionAddEvent.Trigger(session, mra)
			case event.MessageReactionRemoveKey:
				mrr := &event.MessageReactionRemoveBox{}
				mrr.Ctx = ctx
				event.Unmarshal(data, mrr)

				go c.MessageReactionRemoveEvent.Trigger(session, mrr)
			case event.MessageReactionRemoveAllKey:
				mrra := &event.MessageReactionRemoveAllBox{}
				mrra.Ctx = ctx
				event.Unmarshal(data, mrra)

				go c.MessageReactionRemoveAllEvent.Trigger(session, mrra)
			case event.PresenceUpdateKey:
				pu := &event.PresenceUpdateBox{}
				pu.Ctx = ctx
				event.Unmarshal(data, pu)

				go c.PresenceUpdateEvent.Trigger(session, pu)
				//
				// g, err := c.State.Guild(pu.GuildID)
				// if err != nil {
				// 	panic("you haven't correctly cached all guilds you fool!")
				// }
				// presence := &discord.Presence{
				// 	User:   pu.User,
				// 	Roles:  pu.RoleIDs,
				// 	Game:   pu.Game,
				// 	Status: pu.Status,
				// }
				// g.UpdatePresence(presence)
			case event.TypingStartKey:
				ts := &event.TypingStartBox{}
				ts.Ctx = ctx
				event.Unmarshal(data, ts)

				go c.TypingStartEvent.Trigger(session, ts)
			case event.UserUpdateKey:
				u := &event.UserUpdateBox{}
				u.Ctx = ctx
				event.Unmarshal(data, u)

				// dispatch event
				go c.UserUpdateEvent.Trigger(session, u)

				// update cache
				//c.State.UpdateMySelf(u.User)
			case event.VoiceStateUpdateKey:
				vsu := &event.VoiceStateUpdateBox{}
				vsu.Ctx = ctx
				event.Unmarshal(data, vsu)

				go c.VoiceStateUpdateEvent.Trigger(session, vsu)
			case event.VoiceServerUpdateKey:
				vsu := &event.VoiceServerUpdateBox{}
				vsu.Ctx = ctx
				event.Unmarshal(data, vsu)

				go c.VoiceServerUpdateEvent.Trigger(session, vsu)
			case event.WebhooksUpdateKey:
				wsu := &event.WebhooksUpdateBox{}
				wsu.Ctx = ctx
				event.Unmarshal(data, wsu)

				go c.WebhooksUpdateEvent.Trigger(session, wsu)

			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name(), string(evt.Data()))
			}
		}
	}
}
