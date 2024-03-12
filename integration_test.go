//go:build integration
// +build integration

package disgord

import (
	"github.com/andersfylling/disgord/internal/logger"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	BotTokenEnvKey       = "DISGORD_TOKEN_INTEGRATION_TEST"
	TestGuildsNamePrefix = "3643-disgord-test-guild-"
)

var sharedTestSession = CreateTestSession()

func CreateTestSession() *TestSession {
	session := &TestSession{}
	session.once.Do(func() {
		session.session = New(Config{
			BotToken: os.Getenv(BotTokenEnvKey),
			Logger:   &logger.FmtPrinter{},
		})

		user, err := session.session.CurrentUser().Get()
		if err != nil {
			panic(err)
		}

		session.botID = user.ID
	})

	return session
}

type TestSession struct {
	sync.Mutex
	once sync.Once

	botID   Snowflake
	session *Client
	servers []*Guild
}

func (ts *TestSession) UnixInDays() int64 {
	return time.Now().Unix() / 60 / 24 // seconds => hours => days
}

func (ts *TestSession) CreateGuildName() string {
	int64ToStr := func(n int64) string {
		return strconv.FormatInt(n, 10)
	}

	unixInDays := ts.UnixInDays()
	randomSessionID := int64(rand.Int31n(100))
	return TestGuildsNamePrefix + int64ToStr(unixInDays) + "-" + int64ToStr(randomSessionID)
}

func (ts *TestSession) GuildNameIsOlderThanOneDay(name string) bool {
	suffix := strings.TrimPrefix(name, TestGuildsNamePrefix)
	if len(suffix) < 2 {
		return false
	}

	segments := strings.Split(suffix, "-")
	if len(segments) < 2 {
		return false
	}

	daysStr := segments[0]
	days, err := strconv.ParseInt(daysStr, 10, 64)
	if err != nil {
		return false
	}

	if days == 0 {
		return false
	}

	return days+1 < ts.UnixInDays()
}

func (ts *TestSession) GuildIsForTesting(guild *Guild) bool {
	return strings.HasPrefix(guild.Name, TestGuildsNamePrefix)
}

func (ts *TestSession) Cleanup() {
	ts.Lock()
	defer ts.Unlock()

	s := ts.session

	// delete locally saved guilds
	for _, server := range ts.servers {
		_ = s.Guild(server.ID).Delete()
	}

	// detect generated test guilds and delete those as well
	// in case the test session exited too early previously
	for _, guildID := range s.GetConnectedGuilds() {
		guild, err := s.Guild(guildID).Get()
		if err != nil {
			continue
		}

		if !ts.GuildIsForTesting(guild) {
			continue
		}

		if !ts.GuildNameIsOlderThanOneDay(guild.Name) {
			continue
		}

		_ = s.Guild(guildID).Delete()
	}
}

// CreateNewServer creates a new server for testing, with a pre-determined server name prefix to easily cleanup later.
func (ts *TestSession) CreateNewServer() *Guild {
	ts.Lock()
	defer ts.Unlock()

	name := ts.CreateGuildName()
	server, err := ts.session.CreateGuild(name, &CreateGuild{
		Channels: []*PartialChannel{
			{
				Name: "first",
				Type: ChannelTypeGuildText,
			},
			{
				Name: "second",
				Type: ChannelTypeGuildText,
			},
			{
				Name: "first-voice",
				Type: ChannelTypeGuildVoice,
			},
			{
				Name: "second-voice",
				Type: ChannelTypeGuildVoice,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	ts.servers = append(ts.servers, server)
	return server
}
