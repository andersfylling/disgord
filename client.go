package disgord

import (
	"context"
	"encoding/json"
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
			switch eventName {
			case event.Hello:
			case event.Ready:
				r := &event.ReadyBox{}
				err := json.Unmarshal(evt.Data(), r)
				if err != nil {
					panic(err) // TODO: remove panic before merging to master
				}

				go c.ReadyEvent.Trigger(ctx, r)

				// allocate cache mem for guilds
				for _, gu := range r.Guilds {
					g := guild.NewGuildFromUnavailable(gu)
					c.State.AddGuild(g) // checks if the ID already exists
				}

				// updated myself
				c.State.UpdateMySelf(r.User)

			case event.Resumed:
				resumed := &event.ResumedBox{}
				err := json.Unmarshal(evt.Data(), resumed)
				if err != nil {
					panic(err)
				}

				// no need to handle this as its done at the socket level ...
				go c.ResumedEvent.Trigger(ctx, resumed)
			case event.InvalidSession:
			case event.ChannelCreate, event.ChannelUpdate, event.ChannelDelete:
				chanContent := &channel.Channel{}
				err := json.Unmarshal(evt.Data(), chanContent)
				if err != nil {
					panic(err)
				}

				switch eventName {
				case event.ChannelCreate:
					go c.ChannelCreateEvent.Trigger(ctx, &event.ChannelCreateBox{chanContent})
				case event.ChannelUpdate:
					go c.ChannelUpdateEvent.Trigger(ctx, &event.ChannelUpdateBox{chanContent})
				case event.ChannelDelete:
					go c.ChannelDeleteEvent.Trigger(ctx, &event.ChannelDeleteBox{chanContent})
				}
			//case event.ChannelPinsUpdate:
			case event.GuildCreate, event.GuildUpdate, event.GuildDelete:
				g := &guild.Guild{}
				err := json.Unmarshal(evt.Data(), g)
				if err != nil {
					panic(err) // TODO: remove panic before merging to master
				}

				switch eventName { // internal switch statement for guild events
				case event.GuildCreate:
					// notifify listeners
					go c.GuildCreateEvent.Trigger(ctx, &event.GuildCreateBox{g})
					// add to cache
					c.State.AddGuild(g)
				case event.GuildUpdate:
					// notifify listeners
					go c.GuildUpdateEvent.Trigger(ctx, &event.GuildUpdateBox{g})
					// update cache
					c.State.UpdateGuild(g)
				case event.GuildDelete:
					// notify listeners
					unavailGuild := discord.NewGuildUnavailable(g.ID)
					go c.GuildDeleteEvent.Trigger(ctx, &event.GuildDeleteBox{unavailGuild})

					cachedGuild, err := c.State.Guild(g.ID)
					if err != nil {
						// guild has not been cached earlier for some reason..
					} else {
						// update instance with complete info.
						// Assumption: The cached version has no outdated information.
						g = nil
						g = cachedGuild
						// delete the guild object from the cache
						c.State.DeleteGuild(g)
					}
				} // END internal switch statement for guild events
			//case event.GuildBanAdd:
			//case event.GuildBanRemove:
			//case event.GuildEmojisUpdate:
			//case event.GuildIntegrationsUpdate:
			//case event.GuildMemberAdd:
			//case event.GuildMemberRemove:
			//case event.GuildMemberUpdate:
			//case event.GuildMemberChunk:
			case event.GuildRoleCreate:
				r := &event.GuildRoleCreateBox{}
				err := json.Unmarshal(evt.Data(), r)
				if err != nil {
					panic(err)
				}

				g, err := c.State.Guild(r.GuildID)
				if err != nil {
					panic("you haven't correctly cached all guilds you fool!")
				}

				go c.GuildRoleCreateEvent.Trigger(ctx, r)

				// add to cache
				g.Lock()
				g.AddRole(r.Role)
				g.Unlock()
			case event.GuildRoleUpdate:
				r := &event.GuildRoleUpdateBox{}
				err := json.Unmarshal(evt.Data(), r)
				if err != nil {
					panic(err)
				}

				g, err := c.State.Guild(r.GuildID)
				if err != nil {
					panic("you haven't correctly cached all guilds you fool!")
				}

				go c.GuildRoleUpdateEvent.Trigger(ctx, r)

				// add to cache
				g.Lock()
				g.UpdateRole(r.Role)
				g.Unlock()
			case event.GuildRoleDelete:
				r := &event.GuildRoleDeleteBox{}
				err := json.Unmarshal(evt.Data(), r)
				if err != nil {
					panic(err)
				}

				g, err := c.State.Guild(r.GuildID)
				if err != nil {
					panic("you haven't correctly cached all guilds you fool!")
				}
				g.Lock()
				g.DeleteRoleByID(r.RoleID)
				g.Unlock()
				// TODO: remove role from guild members...
			case event.MessageCreate, event.MessageUpdate, event.MessageDelete:
				msg := channel.NewMessage()
				err := json.Unmarshal(evt.Data(), msg)
				if err != nil {
					panic(err) // TODO: remove panic before merging to master
				}

				// TODO: should i cache msg?..
				switch eventName {
				case event.MessageCreate:
					go c.MessageCreateEvent.Trigger(ctx, &event.MessageCreateBox{msg})
				case event.MessageUpdate:
					go c.MessageUpdateEvent.Trigger(ctx, &event.MessageUpdateBox{msg})
				case event.MessageDelete:
					deletedMsg := &event.MessageDeleteBox{
						MessageID: msg.ID,
						ChannelID: msg.ChannelID,
					}
					go c.MessageDeleteEvent.Trigger(ctx, deletedMsg)
				}
			//case event.MessageDeleteBulk:
			//case event.MessageReactionAdd:
			//case event.MessageReactionRemove:
			//case event.MessageReactionRemoveAll:
			case event.PresenceUpdate:
				pu := &event.PresenceUpdateBox{}
				err := json.Unmarshal(evt.Data(), pu)
				if err != nil {
					panic(err)
				}

				go c.PresenceUpdateEvent.Trigger(ctx, pu)

				g, err := c.State.Guild(pu.GuildID)
				if err != nil {
					panic("you haven't correctly cached all guilds you fool!")
				}
				presence := &discord.Presence{
					User:   pu.User,
					Roles:  pu.RoleIDs,
					Game:   pu.Game,
					Status: pu.Status,
				}
				g.UpdatePresence(presence)
			case event.TypingStart:
				ts := &event.TypingStartBox{}
				err := json.Unmarshal(evt.Data(), ts)
				if err != nil {
					panic(err)
				}

				go c.TypingStartEvent.Trigger(ctx, ts)
			case event.UserUpdate:
				u := &event.UserUpdateBox{}
				err := json.Unmarshal(evt.Data(), u.User)
				if err != nil {
					panic(err) // TODO: remove panic before merging to master
				}

				// dispatch event
				go c.UserUpdateEvent.Trigger(ctx, u)

				// update cache
				c.State.UpdateMySelf(u.User)
			//case event.VoiceStateUpdate:
			//case event.VoiceServerUpdate:
			//case event.WebhooksUpdate:

			default:
				fmt.Printf("------\nTODO\nImplement event handler for `%s`, data: \n%+v\n------\n\n", evt.Name(), string(evt.Data()))
			}
		}
	}
}
