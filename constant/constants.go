package constant

// DisgordTestBot bot token for a disgord testing bot
const DisgordTestBot = "DISGORD_TEST_BOT"

// DisgordTestGuildAdmin A guild snowflake id where the bot has administrator permission
const DisgordTestGuildAdmin = "DISGORD_TEST_GUILD_ADMIN"

// DisgordTestGuildDefault Guild snowflake id where the default discord
// permission assigned when joining a vanilla guild
const DisgordTestGuildDefault = "DISGORD_TEST_GUILD_DEFAULT"

// DisgordTestGuildDefaultEmojiSnowflake the default guild should have one custom
// emoji which can be retrieved during testing
const DisgordTestGuildDefaultEmojiSnowflake = "DISGORD_TEST_GUILD_DEFAULT_EMOJI_SNOWFLAKE"

// DisgordTestLive set to true to properly test the functionality against
// Discord before a release is drafted
const DisgordTestLive = "DISGORD_TEST_LIVE"

// GitHubURL repository url
const GitHubURL = "https://github.com/andersfylling/disgord"

// DiscordVersion API version
const DiscordVersion = 6

// Version project version
// TODO: git hook which creates a new git tag after a commit, given that
//        the version here has changed or does not exist as a git tag yet.
const Version = "v0.8.0"
