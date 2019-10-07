package endpoint

import "strconv"

// Gateway ...
func Gateway(v int) string {
	return discordAPI + version + strconv.Itoa(v) + gateway
}
