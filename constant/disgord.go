package constant

// GitHubURL repository url
const GitHubURL = "https://github.com/andersfylling/disgord"

const Name = "DisGord"

const UserAgent = "DiscordBot (" + GitHubURL + ", " + Version + ") " // "DiscordBot (%s, %s) %s"

// Version project version
// TODO: git hook which creates a new git tag after a commit, given that
//        the version here has changed or does not exist as a git tag yet.
const Version = "v0.11.0-rc2"
