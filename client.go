package disgord

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/andersfylling/disgord/channel"
	"github.com/andersfylling/disgord/discord"
	"github.com/andersfylling/disgord/discordws"
	"github.com/andersfylling/disgord/endpoint"
	"github.com/andersfylling/disgord/event"
	"github.com/andersfylling/disgord/guild"
	"github.com/sirupsen/logrus"
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
		DAPIVersion:  endpoint.APIVersion,
		DAPIEncoding: endpoint.APIComEncoding,
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

// NewRequiredDisgord same as NewDisgord, but exits program if an error occours
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

			ctx := context.Background()

			// TODO: parsing JSON uses panic and not logging on issues..

			eventName := evt.Name()
			data := evt.Data()

			switch eventName {
			case event.ReadyKey:
				r := &event.ReadyBox{}
				event.Unmarshal(data, r)

				go c.ReadyEvent.Trigger(ctx, r)

				// allocate cache mem for guilds
				for _, gu := range r.Guilds {
					g := guild.NewGuildFromUnavailable(gu)
					c.State.AddGuild(g) // checks if the ID already exists
				}

				// updated myself
				//c.State.UpdateMySelf(r.User)

			case event.ResumedKey:
				resumed := &event.ResumedBox{}
				event.Unmarshal(data, resumed)

				// no need to handle this as its done at the socket level ...
				go c.ResumedEvent.Trigger(ctx, resumed)
			case event.ChannelCreateKey, event.ChannelUpdateKey, event.ChannelDeleteKey:
				chanContent := &channel.Channel{}
				event.Unmarshal(data, chanContent)

				switch eventName {
				case event.ChannelCreateKey:
					go c.ChannelCreateEvent.Trigger(ctx, &event.ChannelCreateBox{chanContent})
				case event.ChannelUpdateKey:
					go c.ChannelUpdateEvent.Trigger(ctx, &event.ChannelUpdateBox{chanContent})
				case event.ChannelDeleteKey:
					go c.ChannelDeleteEvent.Trigger(ctx, &event.ChannelDeleteBox{chanContent})
				}
			case event.ChannelPinsUpdateKey:
				cpu := &event.ChannelPinsUpdateBox{}
				event.Unmarshal(data, cpu)
				go c.ChannelPinsUpdateEvent.Trigger(ctx, cpu)
			case event.GuildCreateKey, event.GuildUpdateKey, event.GuildDeleteKey:
				g := &guild.Guild{}
				event.Unmarshal(data, g)

				switch eventName { // internal switch statement for guild events
				case event.GuildCreateKey:
					// notifify listeners
					go c.GuildCreateEvent.Trigger(ctx, &event.GuildCreateBox{g})
					// add to cache
					//c.State.AddGuild(g)
				case event.GuildUpdateKey:
					// notifify listeners
					go c.GuildUpdateEvent.Trigger(ctx, &event.GuildUpdateBox{g})
					// update cache
					//c.State.UpdateGuild(g)
				case event.GuildDeleteKey:
					// notify listeners
					unavailGuild := discord.NewGuildUnavailable(g.ID)
					go c.GuildDeleteEvent.Trigger(ctx, &event.GuildDeleteBox{unavailGuild})
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
				event.Unmarshal(data, gba)

				go c.GuildBanAddEvent.Trigger(ctx, gba)
			case event.GuildBanRemoveKey:
				gbr := &event.GuildBanRemoveBox{}
				event.Unmarshal(data, gbr)

				go c.GuildBanRemoveEvent.Trigger(ctx, gbr)
			case event.GuildEmojisUpdateKey:
				geu := &event.GuildEmojisUpdateBox{}
				event.Unmarshal(data, geu)

				go c.GuildEmojisUpdateEvent.Trigger(ctx, geu)
			case event.GuildIntegrationsUpdateKey:
				giu := &event.GuildIntegrationsUpdateBox{}
				event.Unmarshal(data, giu)

				go c.GuildIntegrationsUpdateEvent.Trigger(ctx, giu)
			case event.GuildMemberAddKey:
				gma := &event.GuildMemberAddBox{}
				event.Unmarshal(data, gma)

				go c.GuildMemberAddEvent.Trigger(ctx, gma)
			case event.GuildMemberRemoveKey:
				gmr := &event.GuildMemberRemoveBox{}
				event.Unmarshal(data, gmr)

				go c.GuildMemberRemoveEvent.Trigger(ctx, gmr)
			case event.GuildMemberUpdateKey:
				gmu := &event.GuildMemberUpdateBox{}
				event.Unmarshal(data, gmu)

				go c.GuildMemberUpdateEvent.Trigger(ctx, gmu)
			case event.GuildMembersChunkKey:
				gmc := &event.GuildMembersChunkBox{}
				event.Unmarshal(data, gmc)

				go c.GuildMembersChunkEvent.Trigger(ctx, gmc)
			case event.GuildRoleCreateKey:
				r := &event.GuildRoleCreateBox{}
				event.Unmarshal(data, r)

				go c.GuildRoleCreateEvent.Trigger(ctx, r)

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
				event.Unmarshal(data, r)

				go c.GuildRoleUpdateEvent.Trigger(ctx, r)
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
				event.Unmarshal(data, r)

				go c.GuildRoleDeleteEvent.Trigger(ctx, r)
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
					go c.MessageCreateEvent.Trigger(ctx, &event.MessageCreateBox{msg})
				case event.MessageUpdateKey:
					go c.MessageUpdateEvent.Trigger(ctx, &event.MessageUpdateBox{msg})
				case event.MessageDeleteKey:
					go c.MessageDeleteEvent.Trigger(ctx, &event.MessageDeleteBox{
						MessageID: msg.ID,
						ChannelID: msg.ChannelID,
					})
				}
			case event.MessageDeleteBulkKey:
				mdb := &event.MessageDeleteBulkBox{}
				event.Unmarshal(data, mdb)

				go c.MessageDeleteBulkEvent.Trigger(ctx, mdb)
			case event.MessageReactionAddKey:
				mdb := &event.MessageReactionAddBox{}
				event.Unmarshal(data, mdb)

				go c.MessageReactionAddEvent.Trigger(ctx, mdb)
			case event.MessageReactionRemoveKey:
				mdb := &event.MessageReactionRemoveBox{}
				event.Unmarshal(data, mdb)

				go c.MessageReactionRemoveEvent.Trigger(ctx, mdb)
			case event.MessageReactionRemoveAllKey:
				mrra := &event.MessageReactionRemoveAllBox{}
				event.Unmarshal(data, mrra)

				go c.MessageReactionRemoveAllEvent.Trigger(ctx, mrra)
			case event.PresenceUpdateKey:
				pu := &event.PresenceUpdateBox{}
				event.Unmarshal(data, pu)

				go c.PresenceUpdateEvent.Trigger(ctx, pu)
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
				event.Unmarshal(data, ts)

				go c.TypingStartEvent.Trigger(ctx, ts)
			case event.UserUpdateKey:
				u := &event.UserUpdateBox{}
				event.Unmarshal(data, u)

				// dispatch event
				go c.UserUpdateEvent.Trigger(ctx, u)

				// update cache
				//c.State.UpdateMySelf(u.User)
			case event.VoiceStateUpdateKey:
				vsu := &event.VoiceStateUpdateBox{}
				event.Unmarshal(data, vsu)

				go c.VoiceStateUpdateEvent.Trigger(ctx, vsu)
			case event.VoiceServerUpdateKey:
				vsu := &event.VoiceServerUpdateBox{}
				event.Unmarshal(data, vsu)

				go c.VoiceServerUpdateEvent.Trigger(ctx, vsu)
			case event.WebhooksUpdateKey:
				wsu := &event.WebhooksUpdateBox{}
				event.Unmarshal(data, wsu)

				go c.WebhooksUpdateEvent.Trigger(ctx, wsu)

			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name(), string(evt.Data()))
			}
		}
	}
}
