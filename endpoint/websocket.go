package endpoint

import "strconv"

func Gateway(v int) string {
	return discordAPI + version + strconv.Itoa(v) + gateway
}
